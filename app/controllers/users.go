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
		user, err := u.DB.Users.Me(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "%+v", user)
	}
}
