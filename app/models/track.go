package models

import (
	"time"
)

// Track is a data structure representing a particular song from Spotify
type Track struct {
	SpotifyID        string    `json:"spotifyID" gorm:"primary_key:true"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	Name             string    `json:"name"`
	Duration         int       `json:"duration"`
	PreviewURL       string    `json:"previewURL"`
	TrackNumber      int       `json:"trackNumber"`
	URI              string    `json:"URI"`
	Artists          string    `json:"artists"`
	SpotifyArtistIDs string    `json:"spotifyArtistIDs"`
}
