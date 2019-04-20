package models

import (
	"github.com/jinzhu/gorm"
)

// Like is a struct that contains data about how a user "liked" a  playlist
type Like struct {
	gorm.Model
	UserID       uint   `json:"userID"`
	PlaylistID   string `json:"playlistID"`
	PlaylistName string `json:"playlistName"`
	Permalink    string `json:"permalink"`
}
