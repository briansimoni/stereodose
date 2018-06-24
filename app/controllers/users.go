package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/briansimoni/stereodose/app/models"
)

type UsersController struct {
	DB *models.StereoDoseDB
}

func (u *UsersController) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("User").(models.User)
	if !ok {
		http.Error(w, "Unable to obtain user from session", http.StatusInternalServerError)
		return
	}
	temp, err := u.DB.Users.ByID(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.MarshalIndent(temp, " ", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(data))
}
