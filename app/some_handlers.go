package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func webPlayerTest(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./app/views/build/index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(data))

}
