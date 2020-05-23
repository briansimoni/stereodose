package models

import (
	"log"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
)

// StereoDoseDB is a layer on top of Gorm
// It allows callers to easily use the structs relevant to the rest of the app
type StereoDoseDB struct {
	DB        *gorm.DB
	store     *sessions.CookieStore
	Users     UserService
	Playlists PlaylistService
	Comments  CommentService
	Feedback  FeedbackService
}

// NewStereodoseDB takes a reference to gorm and returns
// an abstraction for use throughout the app
func NewStereodoseDB(db *gorm.DB, s *sessions.CookieStore) *StereoDoseDB {
	// db = db.Debug()
	err := db.AutoMigrate(
		User{},
		Playlist{},
		UserImage{},
		PlaylistImage{},
		Track{},
		Comment{},
		Like{},
		Feedback{},
	).Error

	if err != nil {
		log.Fatal(err.Error())
	}
	database := &StereoDoseDB{}
	database.DB = db
	database.store = s
	database.Users = &StereodoseUserService{db: db}
	database.Playlists = &StereodosePlaylistService{db: db}
	database.Comments = &StereodoseCommentService{db: db}
	database.Feedback = NewFeedbackService(db)
	//test
	// u, _ := database.Users.ByID(1)
	// _, err = database.Playlists.Like("6DRd1s2Hx7VEWWV85GYx6S", *u)
	// err = database.Playlists.Unlike("6DRd1s2Hx7VEWWV85GYx6S", 1)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	return database
}
