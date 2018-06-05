package models

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/briansimoni/stereodose/app/util"
)

type UserService interface {
	GetUser(req *http.Request) (*User, error)
}

type StereodoseUserService struct {
	store *sessions.CookieStore
	db    *gorm.DB
}

type User struct {
	gorm.Model
	Birthdate   string
	DisplayName string
	Email       string
	SpotifyID   string `gorm:"unique;not null"`
	//Images      []string
}

// GetUser first checks to see if the user already exists
// if it doesn't it creates one, otherwise it returns a pointer to user
func (u *StereodoseUserService) GetUser(req *http.Request) (*User, error) {
	s, err := util.GetSessionInfo(u.store, req)
	if err != nil {
		return nil, err
	}
	user := &User{
		SpotifyID: s.SpotifyUserID,
	}

	err = u.db.FirstOrCreate(user, "spotify_id = ?", s.SpotifyUserID).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
