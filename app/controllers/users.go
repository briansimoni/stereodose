package controllers

import (
	"net/http"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/pkg/errors"
)

// UsersController contains methods for reading user data
type UsersController struct {
	DB *models.StereoDoseDB
}

// NewUsersController returns a pointer to UsersController
func NewUsersController(db *models.StereoDoseDB) *UsersController {
	return &UsersController{DB: db}
}

// Me returns the requesting and authenticated user's data
func (u *UsersController) Me(w http.ResponseWriter, r *http.Request) error {
	user, ok := r.Context().Value("User").(models.User)
	if !ok {
		return errors.WithStack(errors.New("Unable to obtain user from request context"))
	}
	data, err := u.DB.Users.ByID(user.ID)
	if err != nil {
		return errors.WithStack(err)
	}
	err = util.JSON(w, data)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
