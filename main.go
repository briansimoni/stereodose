package main

import (
	"encoding/gob"
	"net/http"
	"os"

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
		log.Fatal(err.Error())
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
		log.Fatal(err.Error())
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	stereodose := app.InitApp(c, db)
	log.Println("Starting stereodose app on port", port)
	log.Fatal(http.ListenAndServe(":"+port, stereodose))
}

// Register the oauth2.Token type so we can store it in sessions later
// additionally set the logger to either JSON or plaintext output
func init() {
	gob.Register(oauth2.Token{})
	// TODO: configure logger based on environment variables
	logger := &log.JSONFormatter{}
	log.SetFormatter(logger)

}
