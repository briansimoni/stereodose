package main

import (
	"log"
	"net/http"
	"os"

	"github.com/briansimoni/stereodose/app"
	"github.com/briansimoni/stereodose/config"
)

func main() {
	dbString := os.Getenv("STEREODOSE_DB_STRING")
	if dbString == "" {
		// docker-compose default
		dbString = "postgresql://postgres:development@db:5432/stereodose?sslmode=disable"
	}
	c := &config.Config{
		ClientID:           os.Getenv("STEREODOSE_CLIENT_ID"),
		ClientSecret:       os.Getenv("STEREODOSE_CLIENT_SECRET"),
		AuthKey:            os.Getenv("STEREODOSE_AUTH_KEY"),
		RedirectURL:        os.Getenv("STEREODOSE_REDIRECT_URL"),
		DBConnectionString: dbString,
	}
	err := c.Verify()
	if err != nil {
		log.Fatal(err.Error())
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	stereodose := app.InitApp(c)
	log.Println("Starting stereodose app on port", port)
	log.Fatal(http.ListenAndServe(":"+port, stereodose))
}
