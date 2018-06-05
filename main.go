package main

import (
	"log"
	"net/http"
	"os"

	"github.com/briansimoni/stereodose/app"
	"github.com/briansimoni/stereodose/config"
	"github.com/jinzhu/gorm"
)

func main() {
	connectionString := os.Getenv("STEREODOSE_DB_STRING")
	if connectionString == "" {
		// docker-compose default
		connectionString = "postgresql://postgres:development@db:5432/stereodose?sslmode=disable"
	}
	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		panic(err.Error())
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
		port = "3000"
	}
	stereodose := app.InitApp(c, db)
	log.Println("Starting stereodose app on port", port)
	log.Fatal(http.ListenAndServe(":"+port, stereodose))
}
