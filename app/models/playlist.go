package models

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/briansimoni/stereodose/app/util"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

// PlaylistService is an interface used to describe all of the behavior
// of some kind of service that a "Playlist Service" should offer
// Useful for mocks/fakes when unit testing
type PlaylistService interface {
	GetPlaylists(params *PlaylistSearchParams) ([]Playlist, error)
	GetByID(ID string) (*Playlist, error)
	GetMyPlaylists(user User) ([]Playlist, error)
	GetRandomPlaylist(category, subcategory string) (*Playlist, error)
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
	SpotifyID           string          `json:"spotifyID" gorm:"primary_key:true"`
	CreatedAt           time.Time       `json:"createdAt"`
	UpdatedAt           time.Time       `json:"updatedAt"`
	Category            string          `json:"category"`
	CategoryDisplayName string          `json:"categoryDisplayName"`
	SubCategory         string          `json:"subCategory"`
	Collaborative       bool            `json:"collaborative"`
	Endpoint            string          `json:"href"`
	Images              []PlaylistImage `json:"images"`
	Name                string          `json:"name"`
	IsPublic            bool            `json:"public"`
	SnapshotID          string          `json:"snapshot_id"`
	Tracks              []Track         `json:"tracks" gorm:"many2many:playlist_tracks"`
	Comments            []Comment       `json:"comments" gorm:"ForeignKey:PlaylistID;AssociationForeignKey:spotify_id"`
	Likes               []Like          `json:"likes" gorm:"ForeignKey:PlaylistID;AssociationForeignKey:spotify_id"`
	LikesCount          uint            `json:"likesCount"`
	URI                 string          `json:"URI"`
	UserID              uint            `json:"userID"`
	BucketImageURL      string          `json:"bucketImageURL"`
	BucketThumbnailURL  string          `json:"bucketThumbnailURL"`
	Permalink           string          `json:"permalink"`
	TotalTracks         int             `json:"totalTracks"`
}

// PlaylistSearchParams can be created using URL query parameters
type PlaylistSearchParams struct {
	Offset      string
	Limit       string
	Category    string
	Subcategory string
	SortKey     string
	Order       string
	SpotifyIDs  []string
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
// callers of this method are responsible for checking the limit and offset
func (s *StereodosePlaylistService) GetPlaylists(params *PlaylistSearchParams) ([]Playlist, error) {
	playlists := []Playlist{}

	db := s.db.
		Offset(params.Offset).
		Limit(params.Limit)

	if params.Category == "" && params.Subcategory != "" {
		db = db.Where("category = ?", params.Category)
	}

	if params.Category != "" && params.Subcategory == "" {
		db = db.Where("category = ?", params.Category)
	}

	if params.Category != "" && params.Subcategory != "" {
		db = db.Where("category = ? AND sub_category = ?", params.Category, params.Subcategory)
	}

	if len(params.SpotifyIDs) > 0 {
		db = db.Where("spotify_id IN(?)", params.SpotifyIDs)
	}

	err := db.Order(fmt.Sprintf("%s %s", params.SortKey, params.Order)).Find(&playlists).Error

	if err != nil {
		return nil, err
	}
	return playlists, nil
}

// GetByID returns a playlist populated with all of its tracks
func (s *StereodosePlaylistService) GetByID(ID string) (*Playlist, error) {
	playlist := &Playlist{}
	err := s.db.Preload("Tracks").Preload("Comments.Playlist").Preload("Likes.Playlist").Find(playlist, "spotify_id = ?", ID).Error
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

// GetRandomPlaylist tells the database to grab a random set of playlists from the selected category
// then a random set of tracks is selected across those playlists to get a completely new playlist
// made up of randomly selected tracks
// Using the gorm.Expr somewhat breaks the compatibility with other databases
// as the random() function is supported in Postgres but it's rand() in MySQL
func (s *StereodosePlaylistService) GetRandomPlaylist(category, subcategory string) (*Playlist, error) {
	playlists := []Playlist{}
	err := s.db.Where("category = ? AND sub_category = ?", category, subcategory).Preload("Tracks").Order(gorm.Expr("random()")).Limit(10).Find(&playlists).Error
	if err != nil {
		return nil, err
	}
	if len(playlists) == 0 {
		return nil, fmt.Errorf("Unable to create random playlist. No results were found for %s %s", category, subcategory)
	}

	// create a random playlist of length 20
	randomPlaylist := randomPlaylistFromSet(playlists, 20)
	return randomPlaylist, nil
}

// randomPlaylistFromSet can be given a list of playlists
// it will iterate through the list and select a random track to append to a new playlist
// the length parameter is how many tracks you want in the random playlist
// Since this function doesn't make any calls on outside resources, it is easy to unit test
func randomPlaylistFromSet(playlists []Playlist, length int) *Playlist {
	randomPlaylist := new(Playlist)
	randomPlaylist.Tracks = make([]Track, 0)
	i := 0
	for len(randomPlaylist.Tracks) < length {
		if i > len(playlists)-1 {
			i = 0
		}
		playlist := playlists[i]
		// if the playlist is empty for whatever reason, skip to the next one
		if len(playlist.Tracks) == 0 {
			i++
			continue
		}

		// if the length of the playlist is only 1, then the only song to pick is the first one
		var trackIndex int
		if len(playlist.Tracks) == 1 {
			trackIndex = 0
		} else {
			// otherwise, we randomly grab a track from the playlist
			trackIndex = rand.Intn(len(playlist.Tracks))
		}

		randomPlaylist.Tracks = append(randomPlaylist.Tracks, playlist.Tracks[trackIndex])
		if (i + 1) == len(playlists) {
			i = 0
		} else {
			i++
		}
	}
	return randomPlaylist
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
	categoryObj, err := getCategoryFromName(category)
	if err != nil {
		return nil, err
	}

	permalink := fmt.Sprintf("/%s/%s/%s", category, subCategory, playlistID)
	playlist := &Playlist{
		SpotifyID:           string(list.ID),
		Collaborative:       list.Collaborative,
		Endpoint:            list.Endpoint,
		Name:                list.Name,
		IsPublic:            list.IsPublic,
		SnapshotID:          list.SnapshotID,
		URI:                 string(list.URI),
		UserID:              user.ID,
		Category:            category,
		CategoryDisplayName: categoryObj.DisplayName,
		SubCategory:         subCategory,
		BucketImageURL:      image,
		BucketThumbnailURL:  thumbnail,
		Permalink:           permalink,
	}
	for _, image := range list.Images {
		playlist.Images = append(playlist.Images, PlaylistImage{Image: image})
	}
	tracks, err := getAllPlaylistTracks(c, spotify.ID(playlist.SpotifyID))
	if err != nil {
		return nil, err
	}
	playlist.TotalTracks = len(tracks)

	if playlist.TotalTracks < 5 {
		return nil, &util.StatusError{
			Code:    http.StatusBadRequest,
			Message: "your playlist needs 5 or more songs",
		}
	}

	for _, trk := range tracks {
		track := trk.Track

		// Apparently there are very rare cases where some tracks obtained through
		// the Spotify API don't have an ID. This is required by my database
		// so we just skip the ones that would've otherwise violated this constraint.
		if string(track.ID) == "" {
			// TODO: add transactionID to this log statement
			log.WithFields(logrus.Fields{
				"User":       user.ID,
				"PlaylistID": playlist.SpotifyID,
				"TrackID":    string(track.ID),
				"TrackName":  track.Name,
			}).Warn("This track was skipped during playlist creation")
			continue
		}
		trackToAdd := Track{
			SpotifyID:        string(track.ID),
			Name:             track.Name,
			Duration:         track.Duration,
			PreviewURL:       track.PreviewURL,
			TrackNumber:      track.TrackNumber,
			URI:              string(track.URI),
			Artists:          simpleArtistsToString(track.Artists),
			SpotifyArtistIDs: simpleArtistIdsToString(track.Artists),
		}
		playlist.Tracks = append(playlist.Tracks, trackToAdd)
	}

	err = s.db.Save(playlist).Error
	if err != nil {
		return nil, err
	}

	playlist, err = s.GetByID(playlist.SpotifyID)
	if err != nil {
		return nil, err
	}
	if len(playlist.Likes) > 0 {
		playlist.LikesCount = uint(len(playlist.Likes))
		err = s.db.Save(playlist).Error
		if err != nil {
			return nil, err
		}
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

// simpleArtistIdsToString is pretty much the same as the above function.
// Just note that I removed the space
func simpleArtistIdsToString(artists []spotify.SimpleArtist) string {
	data := make([]string, 0)
	for _, artist := range artists {
		data = append(data, artist.ID.String())
	}
	return strings.Join(data, ",")
}

// DeletePlaylist hard deletes the playlist (only from the StereodoseDB)
// TODO: find the related likes/comments and delete those too
func (s *StereodosePlaylistService) DeletePlaylist(spotifyID string) error {
	if spotifyID == "" {
		return errors.New("spotifyID was empty string")
	}
	playlist, err := s.GetByID(spotifyID)
	if err != nil {
		return err
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

	// grab the permalink from the playlist struct
	playlist := &Playlist{}
	err := s.db.Find(playlist, "spotify_id = ?", playlistID).Error
	if err != nil {
		return nil, err
	}

	comment := &Comment{
		Content:     text,
		UserID:      user.ID,
		PlaylistID:  playlistID,
		DisplayName: user.DisplayName,
		Permalink:   playlist.Permalink,
	}

	err = s.db.Create(comment).Error
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
// TODO: need to refactor this so the Playlist struct "knows" about who owns the likes
// otherwise, if a playlist gets deleted/created again, the likes count can drop to negative numbers
// comments works like this
// could be problematic for very large number of likes
func (s *StereodosePlaylistService) Like(playlistID string, user User) (*Like, error) {
	if playlistID == "" {
		return nil, errors.New("spotifyID was empty string")
	}

	like := &Like{
		PlaylistID: playlistID,
		UserID:     user.ID,
	}

	tx := s.db.Begin()

	var playlist Playlist
	err := tx.Where("spotify_id = ?", playlistID).Find(&playlist).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	like.Permalink = playlist.Permalink
	like.PlaylistName = playlist.Name

	err = tx.Create(like).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Model(&playlist).Update("likes_count", playlist.LikesCount+1).Error
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

	playlist := &Playlist{}
	err := s.db.Preload("Likes").Find(playlist, "spotify_id = ?", playlistID).Error
	if err != nil {
		return err
	}

	authorized := false
	for _, like := range playlist.Likes {
		if like.PlaylistID == playlistID {
			authorized = true
			break
		}
	}
	if !authorized {
		return errors.New("This like does not belong to this playlist")
	}

	like := &Like{}
	like.ID = likeID

	// idk I was drunk and this fixed some bug. I don't even know how
	// It just creates a new array of likes. It does not include the like
	// that we are deleting
	playlist.Likes = func() []Like {
		newList := make([]Like, 0)
		for _, l := range playlist.Likes {
			if l.PlaylistID != playlistID {
				newList = append(newList, l)
			}
		}
		return newList
	}()

	tx := s.db.Begin()
	err = tx.Delete(like).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Model(&playlist).Update("likes_count", playlist.LikesCount-1).Error
	if err != nil {
		return err
	}
	return tx.Commit().Error
}
