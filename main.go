package main

import (
	"encoding/gob"
	"net/http"
	"os/signal"
	"syscall"
	"os"
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/briansimoni/stereodose/app"
	"github.com/briansimoni/stereodose/config"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
)

func main() {
	connectionString := os.Getenv("STEREODOSE_DB_STRING")
	if connectionString == "" {
		// docker-compose default
		connectionString = "postgresql://postgres:development@db:5432/stereodose?sslmode=disable"
		// localhost default
		// connectionString = "postgresql://postgres:development@127.0.0.1:5432/stereodose?sslmode=disable"
	}
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.WithFields(log.Fields{
			"Type": "AppLog",
		}).Fatal("Unable to connect to the database", err.Error())
	}
	defer db.Close()

	c := &config.Config{
		ClientID:      os.Getenv("STEREODOSE_CLIENT_ID"),
		ClientSecret:  os.Getenv("STEREODOSE_CLIENT_SECRET"),
		AuthKey:       os.Getenv("STEREODOSE_AUTH_KEY"),
		RedirectURL:   os.Getenv("STEREODOSE_REDIRECT_URL"),
		EncryptionKey: os.Getenv("STEREODOSE_ENCRYPTION_KEY"),
	}
	err = c.Verify()

	if err != nil {
		log.WithFields(log.Fields{
			"Type": "AppLog",
		}).Fatal("Incorrect or missing configuration", err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	stereodose := app.InitApp(c, db)
	server := http.Server{
		Addr: ":"+port,
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
