package models

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// PlaylistService is an interface used to describe all of the behavior
// of some kind of service that a "Playlist Service" should offer
// Useful for mocks/fakes when unit testing
type PlaylistService interface {
	GetPlaylists(offset, limit, category, subcategory string) ([]Playlist, error)
	GetByID(ID string) (*Playlist, error)
	GetMyPlaylists(user User) ([]Playlist, error)
	// TODO: refactor this to take a Playlist struct instead of a ton of strings
	CreatePlaylistBySpotifyID(user User, playlistID, category, subCategory, image, thumbnail string) (*Playlist, error)
	DeletePlaylist(spotifyID string) error
	Comment(playlistID, text string, user User) (*Comment, error)
	DeleteComment(commentID uint) error
	Like(playlistID string, user User) (*Like, error)
	Unlike(playlistID string, likeID uint) error
}

// Playlist is the data structure that contains playlist metadata from Spotify
// It additionally has relations to users and tracks on Stereodose
type Playlist struct {
	//gorm.Model
	SpotifyID     string          `json:"spotifyID" gorm:"primary_key:true"`
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
	Category      string          `json:"category"`
	SubCategory   string          `json:"subCategory"`
	Collaborative bool            `json:"collaborative"`
	Endpoint      string          `json:"href"`
	Images        []PlaylistImage `json:"images"`
	Name          string          `json:"name"`
	IsPublic      bool            `json:"public"`
	SnapshotID    string          `json:"snapshot_id"`
	Tracks        []Track         `json:"tracks" gorm:"many2many:playlist_tracks"`
	Comments      []Comment       `json:"comments" gorm:"ForeignKey:PlaylistID;AssociationForeignKey:spotify_id"`
	// Likes is an int representing the total number of likes
	Likes              int    `json:"likes"`
	URI                string `json:"URI"`
	UserID             uint   `json:"userID"`
	BucketImageURL     string `json:"bucketImageURL"`
	BucketThumbnailURL string `json:"bucketThumbnailURL"`
}

// PlaylistImage should contain a URL or reference to an image
// It originally comes from Spotify
type PlaylistImage struct {
	gorm.Model
	spotify.Image
	PlaylistID uint
}

// StereodosePlaylistService contains a db and several methods
// for acting on playlists in the local database
type StereodosePlaylistService struct {
	db *gorm.DB
}

// GetPlaylists takes search parameters and returns a subset of playlists
func (s *StereodosePlaylistService) GetPlaylists(offset, limit, category, subcategory string) ([]Playlist, error) {
	playlists := []Playlist{}

	err := s.db.
		Offset(offset).
		Limit(limit).
		Where("category = ? AND sub_category = ?", category, subcategory).
		Find(&playlists).Error

	if err != nil {
		return nil, err
	}
	return playlists, nil
}

// GetByID returns a playlist populated with all of its tracks
func (s *StereodosePlaylistService) GetByID(ID string) (*Playlist, error) {
	playlist := &Playlist{}
	err := s.db.Preload("Tracks").Preload("Comments").Find(playlist, "spotify_id = ?", ID).Error
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

// GetMyPlaylists returns all of the playlists that belong to a particular User
func (s *StereodosePlaylistService) GetMyPlaylists(user User) ([]Playlist, error) {
	playlists := []Playlist{}
	err := s.db.Find(&playlists, "user_id = ?", user.ID).Error
	if err != nil {
		return nil, err
	}
	return playlists, nil
}

// CreatePlaylistBySpotifyID is given a user and playlistID
// It uses the information to call the Spotify API and save the information to the local db
func (s *StereodosePlaylistService) CreatePlaylistBySpotifyID(user User, playlistID, category, subCategory, image, thumbnail string) (*Playlist, error) {
	// 1. get the tracks for the playlist
	// 2. create playlist, add tracks
	// 3. add to db

	// first we want to make sure the playlist isn't already in the db
	p := &Playlist{}
	s.db.Take(p, "spotify_id = ?", playlistID)
	if p.SpotifyID != "" {
		return nil, errors.New("Playlist already exists")
	}

	tok := &oauth2.Token{AccessToken: user.AccessToken}
	c := spotify.Authenticator{}.NewClient(tok)

	list, err := c.GetPlaylist(spotify.ID(playlistID))
	if err != nil {
		return nil, err
	}
	playlist := &Playlist{
		SpotifyID:          string(list.ID),
		Collaborative:      list.Collaborative,
		Endpoint:           list.Endpoint,
		Name:               list.Name,
		IsPublic:           list.IsPublic,
		SnapshotID:         list.SnapshotID,
		URI:                string(list.URI),
		UserID:             user.ID,
		Category:           category,
		SubCategory:        subCategory,
		BucketImageURL:     image,
		BucketThumbnailURL: thumbnail,
	}
	for _, image := range list.Images {
		playlist.Images = append(playlist.Images, PlaylistImage{Image: image})
	}
	tracks, err := getAllPlaylistTracks(c, spotify.ID(playlist.SpotifyID))
	if err != nil {
		return nil, err
	}
	for i, trk := range tracks {
		track := trk.Track
		log.Println(i, track.Name)
		trackToAdd := Track{
			SpotifyID:   string(track.ID),
			Name:        track.Name,
			Duration:    track.Duration,
			PreviewURL:  track.PreviewURL,
			TrackNumber: track.TrackNumber,
			URI:         string(track.URI),
			Artists:     simpleArtistsToString(track.Artists),
		}
		playlist.Tracks = append(playlist.Tracks, trackToAdd)
	}

	err = s.db.Save(playlist).Error
	if err != nil {
		return nil, err
	}
	return playlist, nil
}

// getAllPlaylistTracks will go through all the pages and build a giant list
// the spotify returns a maximum of 100 tracks per page
// it will probably need to make requests synchronously, so it may be slow
// it would be be best to indicate that a long operation is happening to the end user
func getAllPlaylistTracks(c spotify.Client, ID spotify.ID) ([]spotify.PlaylistTrack, error) {
	tracks := make([]spotify.PlaylistTrack, 0)
	limit := new(int)
	offset := new(int)
	*limit = 100
	*offset = 0

	moreTracks := true

	for moreTracks {
		opts := &spotify.Options{
			Limit:  limit,
			Offset: offset,
		}
		page, err := c.GetPlaylistTracksOpt(ID, opts, "")
		if err != nil {
			return nil, err
		}
		if len(page.Tracks) < 100 {
			moreTracks = false
		}
		*offset = *offset + 100
		for _, track := range page.Tracks {
			tracks = append(tracks, track)
		}
	}
	return tracks, nil
}

// simpleArtistToString is a converts an array of SimpleArtists, to one string.
// more convenient that creating yet more tables and requiring more joins
func simpleArtistsToString(artists []spotify.SimpleArtist) string {
	data := make([]string, 0)
	for _, artist := range artists {
		data = append(data, artist.Name)
	}
	return strings.Join(data, ", ")
}

// DeletePlaylist hard deletes the playlist (only from the StereodoseDB)
func (s *StereodosePlaylistService) DeletePlaylist(spotifyID string) error {
	if spotifyID == "" {
		return errors.New("spotifyID was empty string")
	}
	playlist := &Playlist{
		SpotifyID: spotifyID,
	}
	result := s.db.Delete(playlist)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("Delete failed. Playlist Did not exist")
	}
	return nil
}

// Comment will save a comment to the specified playlist
func (s *StereodosePlaylistService) Comment(playlistID, text string, user User) (*Comment, error) {
	if playlistID == "" {
		return nil, errors.New("spotifyID was empty string")
	}

	comment := &Comment{
		Content:     text,
		UserID:      user.ID,
		PlaylistID:  playlistID,
		DisplayName: user.DisplayName,
	}

	err := s.db.Create(comment).Error
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// DeleteComment will soft delete a comment from a playlist
func (s *StereodosePlaylistService) DeleteComment(commentID uint) error {
	comment := new(Comment)
	comment.ID = commentID
	err := s.db.Delete(comment).Error
	if err != nil {
		return err
	}
	return nil
}

// Like will increment the like column for the respective playlist
// it also adds an entry in the likes table
// it is the responsibility of the caller to make sure the user has not liked the playlist already
// this method by itself is effectively Medium's "claps"
// TODO: refactor likes to how comments works
func (s *StereodosePlaylistService) Like(playlistID string, user User) (*Like, error) {
	if playlistID == "" {
		return nil, errors.New("spotifyID was empty string")
	}

	playlist := new(Playlist)
	err := s.db.Find(playlist, "spotify_id = ?", playlistID).Error
	if err != nil {
		return nil, err
	}

	like := &Like{
		UserID:     user.ID,
		PlaylistID: playlistID,
	}

	// we use a database transaction for an "all or nothing set of database transactions"
	tx := s.db.Begin()
	if err = tx.Create(like).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Debug().Model(&Playlist{}).Where("spotify_id = ?", playlist.SpotifyID).Update("likes", playlist.Likes+1).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}
	return like, nil
}

// Unlike will soft delete a like from a playlist
// basically the same thing as Like but in reverse
func (s *StereodosePlaylistService) Unlike(playlistID string, likeID uint) error {

	if playlistID == "" {
		return errors.New("spotifyID was empty string")
	}

	playlist := new(Playlist)
	err := s.db.Find(playlist, "spotify_id = ?", playlistID).Error
	if err != nil {
		return err
	}

	like := new(Like)
	like.ID = likeID

	// we use a database transaction for an "all or nothing set of database transactions"
	tx := s.db.Begin()
	if err = tx.Delete(like).Error; err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&Playlist{}).Where("spotify_id = ?", playlist.SpotifyID).Update("likes", playlist.Likes-1).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}
