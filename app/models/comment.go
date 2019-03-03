package models

import (
	"github.com/jinzhu/gorm"
)

// Comment is a struct that contains the text content for a comment a user made on a playlist
// There is a one-to-many relationship between playlists and comments
// There is a one-to-many relationship between users and comments
type Comment struct {
	gorm.Model
	Content     string `json:"content" gorm:"type:varchar(10000);"`
	UserID      uint   `json:"userID"`
	PlaylistID  string `json:"playlistID"`
	DisplayName string `json:"displayName"`
}
