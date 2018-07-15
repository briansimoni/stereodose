package models

import "time"

type Track struct {
	// ID          uint   `gorm:"primary_key:true"`
	SpotifyID   string `gorm:"primary_key:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Name        string
	Duration    int
	PreviewURL  string
	TrackNumber int
	URI         string
}
