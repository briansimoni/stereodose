package models

import "github.com/jinzhu/gorm"

type PlaylistService interface {
	CreatePlaylist(*Playlist)
}

type Playlist struct {
	gorm.Model
}
