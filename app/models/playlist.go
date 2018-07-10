package models

import (
	"log"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

const (
	// Weed ...
	Weed = "Weed"
	// LSD ...
	LSD = "LSD"
	// Shrooms ...
	Shrooms = "Shrooms"
	// Ecstasy ...
	Ecstasy = "Ecstasy"
)

type PlaylistService interface {
	GetPlaylists(offset, limit string) ([]Playlist, error)
	GetByID(ID uint) (*Playlist, error)
	CreatePlaylistBySpotifyID(user User, spotifyID string) (*Playlist, error)
}

type Playlist struct {
	gorm.Model
	Category      string `json:"Category"`
	SpotifyID     string `gorm:"not null;unique"`
	Collaborative bool   `json:"collaborative"`
	//ExternalURLs  map[string]string `json:"external_urls"`
	Endpoint   string          `json:"href"`
	Images     []PlaylistImage `json:"images"`
	Name       string          `json:"name"`
	IsPublic   bool            `json:"public"`
	SnapshotID string          `json:"snapshot_id"`
	Tracks     []Track         `json:"tracks" gorm:"many2many:playlist_track"`
	URI        string          `json:"uri"`
	UserID     uint
}

type PlaylistImage struct {
	gorm.Model
	spotify.Image
	PlaylistID uint
}

type StereodosePlaylistService struct {
	store *sessions.CookieStore
	db    *gorm.DB
}

// TODO: narrow this down to the specific category
func (s *StereodosePlaylistService) GetPlaylists(offset, limit string) ([]Playlist, error) {
	playlists := []Playlist{}
	err := s.db.Debug().Offset(offset).Limit(limit).Find(&playlists).Error
	if err != nil {
		return nil, err
	}
	return playlists, nil
}

func (s *StereodosePlaylistService) GetByID(ID uint) (*Playlist, error) {
	playlist := &Playlist{}
	err := s.db.Preload("Tracks").Find(playlist, "id = ?", ID).Error
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

// CreatePlaylist is given a user and playlistID
// It uses the information to call the Spotify API and save the information to the local db
func (s *StereodosePlaylistService) CreatePlaylistBySpotifyID(user User, playlistID string) (*Playlist, error) {
	// 1. get the tracks for the playlist
	// 2. create playlist, add tracks
	// 3. add to db
	tok := &oauth2.Token{AccessToken: user.AccessToken}
	c := spotify.Authenticator{}.NewClient(tok)

	list, err := c.GetPlaylist(user.SpotifyID, spotify.ID(playlistID))
	if err != nil {
		return nil, err
	}
	log.Println("NAME", list.Name)
	playlist := &Playlist{
		SpotifyID:     string(list.ID),
		Collaborative: list.Collaborative,
		Endpoint:      list.Endpoint,
		Name:          list.Name,
		IsPublic:      list.IsPublic,
		SnapshotID:    list.SnapshotID,
		URI:           string(list.URI),
		UserID:        user.ID,
	}
	for _, image := range list.Images {
		playlist.Images = append(playlist.Images, PlaylistImage{Image: image})
	}
	tracksPage, err := c.GetPlaylistTracks(user.SpotifyID, spotify.ID(playlist.SpotifyID))
	if err != nil {
		return nil, err
	}
	for _, trk := range tracksPage.Tracks {
		track := trk.Track
		log.Println(track.Name)
		trackToAdd := Track{
			SpotifyID:   string(track.ID),
			Name:        track.Name,
			Duration:    track.Duration,
			PreviewURL:  track.PreviewURL,
			TrackNumber: track.TrackNumber,
			URI:         string(track.URI),
		}
		playlist.Tracks = append(playlist.Tracks, trackToAdd)
	}

	err = s.db.Debug().Create(playlist).Error
	if err != nil {
		return nil, err
	}
	return playlist, nil
}
