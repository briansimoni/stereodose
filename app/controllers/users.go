package controllers

import (
	"net/http"
	"strconv"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/gorilla/mux"
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

// GetByID grabs the user ID from the path parameter
// fetches from the database and returns JSON to the client
func (u *UsersController) GetByID(w http.ResponseWriter, r *http.Request) error {
	pathVars := mux.Vars(r)
	userID, err := strconv.Atoi(pathVars["id"])
	if err != nil {
		return &util.StatusError{
			Message: "Unable to get the UserID from the path parameter",
			Code:    http.StatusBadRequest,
		}
	}
	user, err := u.DB.Users.ByID(uint(userID))
	if err != nil {
		return errors.WithStack(err)
	}
	// since this endpoint is available to any "authenticated user"
	// the spotify access token needs to be removed to maintain security
	user.AccessToken = ""
	util.JSON(w, user)
	return nil
}
