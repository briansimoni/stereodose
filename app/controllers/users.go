package controllers

import (
	"fmt"
	"net/http"

	"github.com/briansimoni/stereodose/app/models"
)

type UsersController struct {
	DB *models.StereoDoseDB
}

func (u *UsersController) Me() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value("User").(models.User)
		if !ok {
			http.Error(w, "Unable to obtain user from session", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%+v", user)
	}
}
