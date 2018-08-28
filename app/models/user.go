package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

type UserService interface {
	ByID(ID uint) (*User, error)
	FirstOrCreate(user *User, tok *oauth2.Token) (*User, error)
	Update(user *User) error
}

type StereodoseUserService struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	Birthdate   string
	DisplayName string
	Email       string
	// TODO: may want to change this to not unique to handle soft delete cases
	SpotifyID    string `gorm:"unique;not null"`
	RefreshToken string `json:"-"` // Hide the RefreshToken in json responses
	AccessToken  string
	Images       []spotify.Image
	Playlists    []Playlist
}

type UserImage struct {
	gorm.Model
	spotify.Image
	UserID uint
}

// ByID first checks to see if the user already exists
// if it doesn't it creates one, otherwise it returns a pointer to user
func (u *StereodoseUserService) ByID(ID uint) (*User, error) {
	user := &User{}
	err := u.db.Debug().Preload("Playlists").Find(user, "id = ?", ID).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *StereodoseUserService) FirstOrCreate(user *User, tok *oauth2.Token) (*User, error) {
	err := u.db.Debug().FirstOrCreate(user, "spotify_id = ?", user.SpotifyID).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *StereodoseUserService) Update(user *User) error {
	err := u.db.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *StereodoseUserService) DeleteUser(user *User) error {
	err := u.db.Delete(user).Error
	if err != nil {
		return err
	}
	return nil
}
