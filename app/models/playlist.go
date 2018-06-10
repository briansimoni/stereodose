package models

import (
	"github.com/jinzhu/gorm"
)

type PlaylistService interface {
	CreatePlaylist(*Playlist)
}

type Playlist struct {
	gorm.Model
	Href   string
	UserID uint
}

type myPlaylistsResponse struct {
	Href     string         `json:"href"`
	Items    []PlaylistJSON `json:"items"`
	Limit    int            `json:"limit"`
	Next     interface{}    `json:"next"`
	Offset   int            `json:"offset"`
	Previous interface{}    `json:"previous"`
	Total    int            `json:"total"`
}

type PlaylistJSON struct {
	Collaborative bool `json:"collaborative"`
	ExternalUrls  struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href   string `json:"href"`
	ID     string `json:"id"`
	Images []struct {
		Height int    `json:"height"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
	} `json:"images"`
	Name  string `json:"name"`
	Owner struct {
		DisplayName  interface{} `json:"display_name"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
		Href string `json:"href"`
		ID   string `json:"id"`
		Type string `json:"type"`
		URI  string `json:"uri"`
	} `json:"owner"`
	PrimaryColor interface{} `json:"primary_color"`
	Public       bool        `json:"public"`
	SnapshotID   string      `json:"snapshot_id"`
	Tracks       struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"tracks"`
	Type string `json:"type"`
	URI  string `json:"uri"`
}
