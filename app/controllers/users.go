package controllers

import (
	"fmt"
	"net/http"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/gorilla/sessions"
)

func CreateUser(db *models.StereoDoseDB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// s, err := store.Get(r, "_stereodose_session")
		// if err != nil {
		// 	log.Println(err.Error())
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		db.Users.CreateUser(r)
		// fmt.Fprintf(w, "%+v", interface{})
		fmt.Fprintf(w, "maybe created")
	}
}
