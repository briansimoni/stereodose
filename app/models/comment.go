package models

import (
	"github.com/jinzhu/gorm"
)

// Comment is a struct that contains the text content for a comment a user made on a playlist
// There is a one-to-many relationship between playlists and comments
// There is a one-to-many relationship between users and comments
type Comment struct {
	gorm.Model
	Content     string   `json:"content" gorm:"type:varchar(10000);"`
	UserID      uint     `json:"userID"`
	PlaylistID  string   `json:"playlistID"`
	Playlist    Playlist `json:"playlist" gorm:"foreignkey:PlaylistID"`
	DisplayName string   `json:"displayName"`
	Permalink   string   `json:"permalink"`
}

// CommentService is an interface for directly performing actions on the comments table
type CommentService interface {
	ByID(id uint) (*Comment, error)
}

// StereodoseCommentService contains a db and several methods
// for acting on comments in the local database
type StereodoseCommentService struct {
	db *gorm.DB
}

// ByID searches the comments table for a comment matching the given ID
func (s *StereodoseCommentService) ByID(id uint) (*Comment, error) {
	comment := new(Comment)
	comment.ID = id
	err := s.db.Find(comment).Error
	if err != nil {
		return nil, err
	}
	return comment, nil
}
