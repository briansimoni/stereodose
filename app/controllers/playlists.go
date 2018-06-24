package controllers

import (
	"net/http"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
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
