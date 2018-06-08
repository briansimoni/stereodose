package models

import (
	"log"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
)

const connectionString = "postgresql://postgres:development@db:5432/stereodose?sslmode=disable"

type StereoDoseDB struct {
	db    *gorm.DB
	store *sessions.CookieStore
	Users UserService
}

// NewStereodoseDB takes a reference to gorm and returns
// an abstraction for use throughout the app
func NewStereodoseDB(db *gorm.DB, s *sessions.CookieStore) *StereoDoseDB {
	err := db.AutoMigrate(User{}).Error
	if err != nil {
		log.Fatal(err.Error())
	}
	database := &StereoDoseDB{}
	database.db = db
	database.store = s
	database.Users = &StereodoseUserService{db: db, store: s}
	return database
}
