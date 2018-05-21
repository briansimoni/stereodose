package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/briansimoni/stereodose/app"
)

func main() {
	fmt.Println(os.Getenv("STEREODOSE_CLIENT_ID"))
	fmt.Println(os.Getenv("STEREODOSE_CLIENT_SECRET"))
	fmt.Println(os.Getenv("STEREODOSE_AUTH_KEY"))
	fmt.Println(os.Getenv("STEREODOSE_REDIRECT_URL"))
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	stereodose := app.InitApp()
	log.Println("Starting stereodose app on port", port)
	log.Fatal(http.ListenAndServe(":"+port, stereodose))
}
