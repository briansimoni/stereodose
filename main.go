package main

import (
	"context"
	"encoding/gob"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/briansimoni/stereodose/app"
	"github.com/briansimoni/stereodose/config"
	"golang.org/x/oauth2"
)

func main() {
	c := &config.Config{
		ClientID:       os.Getenv("STEREODOSE_CLIENT_ID"),
		ClientSecret:   os.Getenv("STEREODOSE_CLIENT_SECRET"),
		AuthKey:        os.Getenv("STEREODOSE_AUTH_KEY"),
		RedirectURL:    os.Getenv("STEREODOSE_REDIRECT_URL"),
		IOSRedirectURL: os.Getenv("STEREODOSE_IOS_REDIRECT_URL"),
		EncryptionKey:  os.Getenv("STEREODOSE_ENCRYPTION_KEY"),
	}
	err := c.Verify()

	if err != nil {
		log.WithFields(log.Fields{
			"Type": "AppLog",
		}).Fatal("Incorrect or missing configuration", err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	stereodose, db := app.InitApp(c)
	server := http.Server{
		Addr:    ":" + port,
		Handler: stereodose,
	}

	log.WithFields(log.Fields{
		"Type": "AppLog",
	}).Info("Starting Stereodose on port: " + port)

	go func() {
		err = http.ListenAndServe(":"+port, stereodose)
		if err != nil && err != http.ErrServerClosed {
			log.WithFields(log.Fields{
				"Type": "AppLog",
			}).Fatal("The server encountered a fatal error", err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	// 1 is SIGHUP (hangup)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quit
	log.WithFields(log.Fields{
		"Type": "AppLog",
	}).Info("Shutdown Signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.WithFields(log.Fields{
			"Type": "AppLog",
		}).Fatal("Server shutdown", err.Error())
	}

	db.DB.Close()

	log.WithFields(log.Fields{
		"Type": "AppLog",
	}).Info("process exiting without error")
}

// Register the oauth2.Token type so we can store it in sessions later
// additionally set the logger to either JSON or plaintext output
func init() {
	gob.Register(oauth2.Token{})
	logger := &log.JSONFormatter{}
	log.SetFormatter(logger)

}
