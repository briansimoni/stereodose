package models

import (
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserService interface {
	Me(spotifyID string) (*User, error)
}

type StereodoseUserService struct {
	store *sessions.CookieStore
	db    *gorm.DB
}

type User struct {
	gorm.Model
	Birthdate    string
	DisplayName  string
	Email        string
	SpotifyID    string `gorm:"unique;not null"`
	RefreshToken string
	//Images      []string
}

// Me first checks to see if the user already exists
// if it doesn't it creates one, otherwise it returns a pointer to user
// TODO: probably rethink the name of this method
func (u *StereodoseUserService) Me(spotifyID string) (*User, error) {
	user := &User{
		SpotifyID: spotifyID,
	}

	err := u.db.FirstOrCreate(user, "spotify_id = ?", spotifyID).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *StereodoseUserService) FirstOrCreate(user *User) (*User, error) {
	err := u.db.FirstOrCreate(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
