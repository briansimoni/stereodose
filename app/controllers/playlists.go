package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type PlaylistsController struct {
	DB *models.StereoDoseDB
}

func (p *PlaylistsController) GetPlaylists(w http.ResponseWriter, r *http.Request) error {
	// The router will run some regex to make sure that they are [0-9]+
	// no need to check again here
	queryValues := r.URL.Query()
	offset := queryValues.Get("offset")
	limit := queryValues.Get("limit")

	playlists, err := p.DB.Playlists.GetPlaylists(offset, limit)
	if err != nil {
		return errors.WithStack(err)
	}
	err = util.JSON(w, playlists)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

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
		SpotifyID string `json:"SpotifyID"`
	}
	var data jsonBody
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return errors.WithStack(err)
	}
	user, ok := r.Context().Value("User").(models.User)
	if !ok {
		return errors.New("Unable to obtain user from session")
	}
	log.Println("ID", data.SpotifyID)

	_, err = p.DB.Playlists.CreatePlaylistBySpotifyID(user, data.SpotifyID)
	if err != nil {
		return errors.WithStack(err)
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

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
