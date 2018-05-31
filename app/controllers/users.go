package controllers

import (
	"fmt"
	"net/http"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/jinzhu/gorm"
)

func CreateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := models.User{}
		u.CreateNewUser(db, "test@test.com")
		fromDB := models.User{}
		db.First(&fromDB)
		fmt.Fprintf(w, "%+v", fromDB)
	}
}
