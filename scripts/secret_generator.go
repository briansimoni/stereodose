package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

func main() {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err.Error())
	}
	secret := base64.StdEncoding.EncodeToString(b)
	fmt.Println(secret)
}
