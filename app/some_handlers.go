package app

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

func webPlayerTest(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./app/views/index.gohtml")
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

	t, err := template.New("test").Parse(string(data))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s, err := store.Get(r, sessionName)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tok, ok := s.Values["Token"].(oauth2.Token)
	if !ok {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d := struct {
		AccessToken string
	}{
		AccessToken: tok.AccessToken,
	}
	err = t.Execute(w, d)
	if err != nil {
		log.Fatal(err.Error())
	}

}
