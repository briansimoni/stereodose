package models

import "github.com/jinzhu/gorm"

type Track struct {
	gorm.Model
	SpotifyID   string
	Name        string
	Duration    int
	PreviewURL  string
	TrackNumber int
	URI         string
	// PlaylistID  uint
}
