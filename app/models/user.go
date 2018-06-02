package models

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserService interface {
	GetUser(ID string) (*User, error)
	CreateUser(spotifyID string) error
}

type StereodoseUserService struct {
	db *gorm.DB
}

// spotifyUser represents the json returend from the /me endpoint
type spotifyUser struct {
	Birthdate    string      `json:"birthdate"`
	Country      string      `json:"country"`
	DisplayName  interface{} `json:"display_name"`
	Email        string      `json:"email"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  interface{} `json:"href"`
		Total int         `json:"total"`
	} `json:"followers"`
	Href    string        `json:"href"`
	ID      string        `json:"id"`
	Images  []interface{} `json:"images"`
	Product string        `json:"product"`
	Type    string        `json:"type"`
	URI     string        `json:"uri"`
}

type User struct {
	gorm.Model
	Birthdate   string
	DisplayName string
	Email       string
	SpotifyID   string
	//Images      []string
}

// GetUser will retrieve the user from the DB by ID
// If the user does not exist, it will be created
func (u *StereodoseUserService) GetUser(id string) (*User, error) {
	var user *User
	u.db.Where(&User{SpotifyID: id}).First(user)
	return user, nil
}

func (u *StereodoseUserService) CreateUser(spotifyID string) error {
	user := &User{
		SpotifyID: spotifyID,
	}
	u.db.Create(user)
	var retval User
	u.db.First(&retval)
	log.Println(retval)
	return nil
}
