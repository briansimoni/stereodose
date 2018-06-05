package controllers

import (
	"fmt"
	"net/http"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/gorilla/sessions"
)

func GetUser(db *models.StereoDoseDB, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := db.Users.GetUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "%+v", user)
	}
}
