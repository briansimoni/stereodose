package models

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/briansimoni/stereodose/app/util"
)

type UserService interface {
	GetUser(ID string) (*User, error)
	CreateUser(req *http.Request) (*User, error)
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

// GetUser will retrieve the user from the DB by ID
// If the user does not exist, it will be created
func (u *StereodoseUserService) GetUser(id string) (*User, error) {
	var user *User
	u.db.Where(&User{SpotifyID: id}).First(user)
	return user, nil
}

// CreateUser first checks to see if the user already exists
// if it doesn't it creates one, otherwise it returns a pointer to user
func (u *StereodoseUserService) CreateUser(req *http.Request) (*User, error) {
	s, err := util.GetSessionInfo(u.store, req)
	if err != nil {
		return nil, err
	}
	user := &User{
		SpotifyID: s.SpotifyUserID,
	}

	u.db.FirstOrCreate(user, "spotify_id = ?", s.SpotifyUserID)
	log.Println("logging the USER", user)

	var retval User
	u.db.Where("spotify_id = ?", s.SpotifyUserID).First(user)
	log.Println(retval)
	return nil, nil
}
