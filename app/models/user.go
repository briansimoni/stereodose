package models

import (
	"log"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserService interface {
	GetUser(ID string) (*User, error)
	CreateUser(spotifyID string) error
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
	SpotifyID   string `gorm:"AUTO_INCREMENT"`
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
