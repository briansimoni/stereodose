package models

import (
	"log"

	"github.com/gorilla/sessions"
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
	err := user.getMyPlaylists(tok)
	if err != nil {
		return nil, err
	}
	err = u.db.Debug().FirstOrCreate(user).Error
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

func (u *User) getMyPlaylists(tok *oauth2.Token) error {
	c := spotify.Authenticator{}.NewClient(tok)
	result, err := c.CurrentUsersPlaylists()
	if err != nil {
		return err
	}
	log.Println(result.Playlists[0].Tracks.Endpoint)
	for _, playlist := range result.Playlists {
		for _, oldPlaylists := range u.Playlists {
			if string(playlist.ID) == oldPlaylists.SpotifyID {
				break
			}
		}
		tracks, err := c.GetPlaylistTracks(u.SpotifyID, playlist.ID)
		if err != nil {
			return err
		}
		playlistToAdd := Playlist{
			UserID:        u.ID,
			SpotifyID:     playlist.ID.String(),
			Collaborative: playlist.Collaborative,
			Endpoint:      playlist.Endpoint,
			Name:          playlist.Name,
			IsPublic:      playlist.IsPublic,
			SnapshotID:    playlist.SnapshotID,
			URI:           string(playlist.URI),
		}
		for _, track := range tracks.Tracks {
			trackToAdd := Track{
				PlaylistID:  playlistToAdd.ID,
				SpotifyID:   string(track.Track.ID),
				Name:        track.Track.Name,
				Duration:    track.Track.Duration,
				PreviewURL:  track.Track.PreviewURL,
				TrackNumber: track.Track.TrackNumber,
			}
			playlistToAdd.Tracks = append(playlistToAdd.Tracks, trackToAdd)
		}
		u.Playlists = append(u.Playlists, playlistToAdd)
	}
	return nil
}
