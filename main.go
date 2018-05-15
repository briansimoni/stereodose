package main

import (
	"log"
	"net/http"

	"github.com/briansimoni/stereodose/app"
)

func main() {
	stereodose := app.InitApp()
	log.Fatal(http.ListenAndServe(":4000", stereodose))
}
