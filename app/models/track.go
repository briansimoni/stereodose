package models

import "github.com/jinzhu/gorm"

type Track struct {
	gorm.Model
	SpotifyID   string `gorm:"not null;unique"`
	Name        string
	Duration    int
	PreviewURL  string
	TrackNumber int
	URI         string
	// PlaylistID  uint
}
