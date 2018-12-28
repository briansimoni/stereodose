package models

import (
	"log"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
)

type StereoDoseDB struct {
	db        *gorm.DB
	store     *sessions.CookieStore
	Users     UserService
	Playlists PlaylistService
}

// NewStereodoseDB takes a reference to gorm and returns
// an abstraction for use throughout the app
func NewStereodoseDB(db *gorm.DB, s *sessions.CookieStore) *StereoDoseDB {
	// db = db.Debug()
	// db.Debug().DropTable(User{}, Playlist{}, UserImage{}, PlaylistImage{}, Track{})
	err := db.AutoMigrate(User{}, Playlist{}, UserImage{}, PlaylistImage{}, Track{}).Error
	if err != nil {
		log.Fatal(err.Error())
	}
	database := &StereoDoseDB{}
	database.db = db
	database.store = s
	database.Users = &StereodoseUserService{db: db}
	database.Playlists = &StereodosePlaylistService{db: db}
	return database
}
