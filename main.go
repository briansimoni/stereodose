package main

import (
	"log"
	"net/http"
	"os"

	"github.com/briansimoni/stereodose/app"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
	stereodose := app.InitApp()
	log.Println("Starting stereodose app on port", port)
	log.Fatal(http.ListenAndServe(":"+port, stereodose))
}
