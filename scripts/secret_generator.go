package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	authBytes := make([]byte, 64)
	_, err := rand.Read(authBytes)
	if err != nil {
		log.Fatal(err.Error())
	}
	authKey := base64.StdEncoding.EncodeToString(authBytes)

	encBytes := make([]byte, 32)
	_, err = rand.Read(encBytes)
	if err != nil {
		log.Fatal(err.Error())
	}
	encKey := base64.StdEncoding.EncodeToString(encBytes)

	fmt.Println("Auth Key:")
	fmt.Println(authKey)
	fmt.Println("\nEncryption Key:")
	fmt.Println(encKey)
}
