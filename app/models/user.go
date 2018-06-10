package models

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserService interface {
	ByID(ID uint) (*User, error)
	FirstOrCreate(user *User) (*User, error)
	Update(user *User) error
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
	// TODO: may want to change this to not unique to handle soft delete cases
	SpotifyID    string `gorm:"unique;not null"`
	RefreshToken string `json:"-"` // Hide the RefreshToken in json responses
	AccessToken  string `json:"-"`
	//Images      []string
	Playlists []Playlist
}

// Me first checks to see if the user already exists
// if it doesn't it creates one, otherwise it returns a pointer to user
func (u *StereodoseUserService) ByID(ID uint) (*User, error) {
	user := &User{}

	err := u.db.Find(user, "id = ?", ID).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *StereodoseUserService) FirstOrCreate(user *User) (*User, error) {
	err := user.getMyPlaylists()
	if err != nil {
		return nil, err
	}
	err = u.db.FirstOrCreate(user).Error
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

func (u *User) getMyPlaylists() error {
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me/playlists", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+u.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("Response was " + res.Status)
	}
	defer res.Body.Close()
	var playlists myPlaylistsResponse
	err = json.NewDecoder(res.Body).Decode(&playlists)
	if err != nil {
		return err
	}
	// probably am going to need to compare and only add if its not there
	for _, playlist := range playlists.Items {
		myPlaylist := Playlist{
			Href: playlist.Href,
		}
		u.Playlists = append(u.Playlists, myPlaylist)
	}
	return nil
}
