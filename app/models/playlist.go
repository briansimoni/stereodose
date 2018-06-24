package models

import (
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/zmb3/spotify"
)

type PlaylistService interface {
	GetPlaylists() ([]Playlist, error)
}

type Playlist struct {
	gorm.Model
	SpotifyID     string
	Collaborative bool `json:"collaborative"`
	//ExternalURLs  map[string]string `json:"external_urls"`
	Endpoint   string          `json:"href"`
	Images     []PlaylistImage `json:"images"`
	Name       string          `json:"name"`
	IsPublic   bool            `json:"public"`
	SnapshotID string          `json:"snapshot_id"`
	Tracks     []Track         `json:"tracks"`
	URI        string          `json:"uri"`
	UserID     uint
}

type PlaylistImage struct {
	gorm.Model
	spotify.Image
	PlaylistID uint
}

type StereodosePlaylistService struct {
	store *sessions.CookieStore
	db    *gorm.DB
}

func (s *StereodosePlaylistService) GetPlaylists() ([]Playlist, error) {
	playlists := []Playlist{}
	err := s.db.Find(&playlists).Error
	if err != nil {
		return nil, err
	}
	return playlists, nil

}
