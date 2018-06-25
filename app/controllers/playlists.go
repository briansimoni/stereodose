package controllers

import (
	"net/http"
	"strconv"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/gorilla/mux"
)

type PlaylistsController struct {
	DB *models.StereoDoseDB
}

func (p *PlaylistsController) GetPlaylists(w http.ResponseWriter, r *http.Request) error {
	playlists, err := p.DB.Playlists.GetPlaylists()
	if err != nil {
		return err
	}
	err = util.JSON(w, playlists)
	if err != nil {
		return err
	}
	return nil
}

func (p *PlaylistsController) GetPlaylistByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err
	}
	playlist, err := p.DB.Playlists.GetByID(uint(ID))
	if err != nil {
		return err
	}
	err = util.JSON(w, playlist)
	if err != nil {
		return err
	}
	return nil
}
