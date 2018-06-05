package models

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserService interface {
	Me(req *http.Request) (*User, error)
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
func (u *StereodoseUserService) Me(req *http.Request) (*User, error) {
	// s, err := util.GetSessionInfo(u.store, req)
	// if err != nil {
	// 	return nil, err
	// }
	spotifyID, ok := req.Context().Value("SpotifyID").(string)
	if !ok {
		return nil, errors.New("Unable to obtain SpotifyID from context")
	}
	// user := &User{
	// 	SpotifyID: s.SpotifyUserID,
	// }
	user := &User{
		SpotifyID: spotifyID,
	}

	err := u.db.FirstOrCreate(user, "spotify_id = ?", spotifyID).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
