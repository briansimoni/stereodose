package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// PlaylistsController is a collection of RESTful Handlers for Playlists
type PlaylistsController struct {
	DB *models.StereoDoseDB
}

// NewPlaylistsController returns a pointer to PlaylistsController
func NewPlaylistsController(db *models.StereoDoseDB) *PlaylistsController {
	return &PlaylistsController{DB: db}
}

// GetPlaylists will return a subset of all the playlists in the DB
// either offset or limit are required parameters
func (p *PlaylistsController) GetPlaylists(w http.ResponseWriter, r *http.Request) error {
	queryValues := r.URL.Query()
	offset := queryValues.Get("offset")
	limit := queryValues.Get("limit")
	category := queryValues.Get("category")
	subcategory := queryValues.Get("subcategory")
	if offset == "" {
		offset = "0"
	}
	if limit == "" {
		limit = "10"
	}

	playlists, err := p.DB.Playlists.GetPlaylists(offset, limit, category, subcategory)
	if err != nil {
		return errors.WithStack(err)
	}
	err = util.JSON(w, playlists)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// GetPlaylistByID reads the id variable from the url path and sends a JSON response
func (p *PlaylistsController) GetPlaylistByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ID := vars["id"]
	playlist, err := p.DB.Playlists.GetByID(ID)
	if err != nil {
		return errors.WithStack(err)
	}
	err = util.JSON(w, playlist)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// GetMyPlaylists returns all of the playlists added to Stereodose that belong to the requesting user
func (p *PlaylistsController) GetMyPlaylists(w http.ResponseWriter, r *http.Request) error {
	user, ok := r.Context().Value("User").(models.User)
	if !ok {
		return errors.New("Unable to obtain user from session")
	}
	playlists, err := p.DB.Playlists.GetMyPlaylists(user)
	if err != nil {
		return err
	}
	err = util.JSON(w, playlists)
	if err != nil {
		return err
	}
	return nil
}

// CreatePlaylist reads the SpotifyID from the POST body
// It then calls the spotify API to get the full info and store in the local DB
// TODO: return 409 conflict instead of 500 error if playlist already exists
func (p *PlaylistsController) CreatePlaylist(w http.ResponseWriter, r *http.Request) error {
	type jsonBody struct {
		SpotifyID   string `json:"SpotifyID"`
		Category    string `json:"Category"`
		SubCategory string `json:"SubCategory"`
	}
	var data jsonBody
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return &statusError{
			Message: fmt.Sprintf("Error parsing JSON: %s", err.Error()),
			Code:    http.StatusBadRequest,
		}
	}
	valid := models.Categories.Valid(data.Category, data.SubCategory)
	if !valid {
		return &statusError{
			Message: fmt.Sprintf("Invalid Category/Subcategory: %s / %s", data.Category, data.SubCategory),
			Code:    http.StatusBadRequest,
		}
	}
	user, ok := r.Context().Value("User").(models.User)
	if !ok {
		return errors.New("Unable to obtain user from session")
	}

	_, err = p.DB.Playlists.CreatePlaylistBySpotifyID(user, data.SpotifyID, data.Category, data.SubCategory)
	if err != nil {
		return errors.WithStack(err)
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

// DeletePlaylist takes the id variable from the url path
// It performs a hard delete of the playlist from the DB
func (p *PlaylistsController) DeletePlaylist(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ID := vars["id"]
	err := p.DB.Playlists.DeletePlaylist(ID)
	if err != nil {
		if err.Error() == "Delete failed. Playlist Did not exist" {
			return &statusError{errors.WithStack(err), err.Error(), http.StatusBadRequest}
		}
		return errors.WithStack(err)
	}
	return nil
}
