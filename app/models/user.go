package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // This enables the postgres driver for gorm
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// UserService is an interface used to describe all of the behavior
// of some kind of service that a "User Service" should offer
// Useful for mocks/fakes when unit testing
type UserService interface {
	ByID(ID uint) (*User, error)
	BySpotifyID(ID string) (*User, error)
	FirstOrCreate(user *User, tok *oauth2.Token) (*User, error)
	Update(user *User) error
}

// User is the data structure that contains user metadata from Spotify
// It additionally a relation to playlists Stereodose
type User struct {
	gorm.Model
	Birthdate   string `json:"birthDate"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	// TODO: may want to change this to not unique to handle soft delete cases
	SpotifyID    string      `json:"spotifyID" gorm:"unique;not null"`
	RefreshToken string      `json:"-"` // Hide the RefreshToken in json responses
	AccessToken  string      `json:"accessToken"`
	Images       []UserImage `json:"images"`
	Playlists    []Playlist  `json:"playlists"`
	Comments     []Comment   `json:"comments"`
	Likes        []Like      `json:"likes"`
	// Product is the user's subscription level: "premium, free etc..."
	Product string `json:"product"`
}

// UserImage should contain a URL or reference to an image
// It originally comes from Spotify, thus the embedded type
type UserImage struct {
	gorm.Model
	spotify.Image
	UserID uint
}

// StereodoseUserService contains a db and several methods
// for acting on users in the local database
type StereodoseUserService struct {
	db *gorm.DB
}

// ByID first checks to see if the user already exists
// if it doesn't it creates one, otherwise it returns a pointer to user
func (u *StereodoseUserService) ByID(ID uint) (*User, error) {
	user := &User{}
	err := u.db.Preload("Images").Preload("Playlists").Preload("Comments").Preload("Likes").Find(user, "id = ?", ID).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// BySpotifyID searches by the SpotifyID and returns a User
func (u *StereodoseUserService) BySpotifyID(ID string) (*User, error) {
	user := &User{}
	err := u.db.Preload("Images").Preload("Playlists").Preload("Comments").Preload("Likes").Find(user, "spotify_id = ?", ID).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FirstOrCreate finds the first matched User or creates a new one
func (u *StereodoseUserService) FirstOrCreate(user *User, tok *oauth2.Token) (*User, error) {
	err := u.db.FirstOrCreate(user, "spotify_id = ?", user.SpotifyID).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update runs a User update
func (u *StereodoseUserService) Update(user *User) error {
	err := u.db.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser attempts to soft delete a User
func (u *StereodoseUserService) DeleteUser(user *User) error {
	err := u.db.Delete(user).Error
	if err != nil {
		return err
	}
	return nil
}
