package models

import (
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/zmb3/spotify"
)

const (
	// Weed ...
	Weed = "Weed"
	// LSD ...
	LSD = "LSD"
	// Shrooms ...
	Shrooms = "Shrooms"
	// Ecstasy ...
	Ecstasy = "Ectasy"
)

type PlaylistService interface {
	GetPlaylists() ([]Playlist, error)
	GetByID(ID uint) (*Playlist, error)
}

type Playlist struct {
	gorm.Model
	Category      string `json:"Category"`
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

// TODO: narrow this down to the specific category
// TODO: add paging
func (s *StereodosePlaylistService) GetPlaylists() ([]Playlist, error) {
	playlists := []Playlist{}
	err := s.db.Find(&playlists).Error
	if err != nil {
		return nil, err
	}
	return playlists, nil
}

func (s *StereodosePlaylistService) GetByID(ID uint) (*Playlist, error) {
	playlist := &Playlist{}
	err := s.db.Preload("Tracks").Find(playlist, "id = ?", ID).Error
	if err != nil {
		return nil, err
	}
	return playlist, nil
}
