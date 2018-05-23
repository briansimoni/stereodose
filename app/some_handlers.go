package app

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func webPlayerTest(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./app/templates/index.html")
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

	tok, ok := s.Values["Access_Token"].(string)
	if !ok {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	d := struct {
		AccessToken string
	}{
		AccessToken: tok,
	}
	t.Execute(w, d)

}
