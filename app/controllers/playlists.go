package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	// register jpeg type
	_ "image/jpeg"
	// register png type
	_ "image/png"
	// register gif type
	_ "image/gif"
	// register webp type
	_ "golang.org/x/image/webp"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/google/go-cloud/blob"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// PlaylistsController is a collection of RESTful Handlers for Playlists
type PlaylistsController struct {
	DB     *models.StereoDoseDB
	Bucket *blob.Bucket
}

// NewPlaylistsController returns a pointer to PlaylistsController
func NewPlaylistsController(db *models.StereoDoseDB, b *blob.Bucket) *PlaylistsController {
	return &PlaylistsController{DB: db, Bucket: b}
}

// GetPlaylists will return a subset of all the playlists in the DB
// either offset or limit are required parameters
func (p *PlaylistsController) GetPlaylists(w http.ResponseWriter, r *http.Request) error {
	queryValues := r.URL.Query()
	offset := queryValues.Get("offset")
	limit := queryValues.Get("limit")
	category := queryValues.Get("category")
	subcategory := queryValues.Get("subcategory")
	if offset == "" {
		offset = "0"
	}
	if limit == "" {
		limit = "10"
	}

	playlists, err := p.DB.Playlists.GetPlaylists(offset, limit, category, subcategory)
	if err != nil {
		return errors.WithStack(err)
	}
	err = util.JSON(w, playlists)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// GetPlaylistByID reads the id variable from the url path and sends a JSON response
func (p *PlaylistsController) GetPlaylistByID(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ID := vars["id"]
	playlist, err := p.DB.Playlists.GetByID(ID)
	if err != nil {
		return errors.WithStack(err)
	}
	err = util.JSON(w, playlist)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// GetMyPlaylists returns all of the playlists added to Stereodose that belong to the requesting user
func (p *PlaylistsController) GetMyPlaylists(w http.ResponseWriter, r *http.Request) error {
	user, ok := r.Context().Value("User").(models.User)
	if !ok {
		return errors.New("Unable to obtain user from session")
	}
	playlists, err := p.DB.Playlists.GetMyPlaylists(user)
	if err != nil {
		return err
	}
	err = util.JSON(w, playlists)
	if err != nil {
		return err
	}
	return nil
}

// CreatePlaylist reads the SpotifyID from the POST body
// It then calls the spotify API to get the full info and store in the local DB
// TODO: return 409 conflict instead of 500 error if playlist already exists
func (p *PlaylistsController) CreatePlaylist(w http.ResponseWriter, r *http.Request) error {
	type jsonBody struct {
		SpotifyID   string `json:"SpotifyID"`
		Category    string `json:"Category"`
		SubCategory string `json:"SubCategory"`
		ImageURL    string `json:"ImageURL"`
	}
	var data jsonBody
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return &statusError{
			Message: fmt.Sprintf("Error parsing JSON: %s", err.Error()),
			Code:    http.StatusBadRequest,
		}
	}
	valid := models.Categories.Valid(data.Category, data.SubCategory)
	if !valid {
		return &statusError{
			Message: fmt.Sprintf("Invalid Category/Subcategory: %s / %s", data.Category, data.SubCategory),
			Code:    http.StatusBadRequest,
		}
	}
	user, ok := r.Context().Value("User").(models.User)
	if !ok {
		return errors.New("Unable to obtain user from session")
	}

	_, err = p.DB.Playlists.CreatePlaylistBySpotifyID(user, data.SpotifyID, data.Category, data.SubCategory, data.ImageURL)
	if err != nil {
		return errors.WithStack(err)
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

// DeletePlaylist takes the id variable from the url path
// It performs a hard delete of the playlist from the DB
func (p *PlaylistsController) DeletePlaylist(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ID := vars["id"]
	user, ok := r.Context().Value("User").(models.User)
	if !ok {
		return errors.New("Unable to obtain user from session")
	}
	targetPlaylist, err := p.DB.Playlists.GetByID(ID)
	if err != nil {
		return errors.WithStack(err)
	}
	if targetPlaylist == nil {
		return &statusError{
			Message: fmt.Sprintf("Playlist does not exist"),
			Code:    http.StatusNotFound,
		}
	}
	userPlaylists, err := p.DB.Playlists.GetMyPlaylists(user)
	if err != nil {
		return errors.WithStack(err)
	}
	authorized := false
	for _, playlist := range userPlaylists {
		if playlist.SpotifyID == ID {
			authorized = true
			break
		}
	}
	if !authorized {
		return &statusError{
			Message: fmt.Sprintf("Unauthorized to remove this playlist"),
			Code:    http.StatusUnauthorized,
		}
	}
	err = p.DB.Playlists.DeletePlaylist(ID)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UploadImage saves an image the corresponds to a playlist
// The actual data is saved to cloud bucket storage
// A permalink to the object is stored in the database and returned to the client
// TODO: resize to optimal dimensions
func (p *PlaylistsController) UploadImage(w http.ResponseWriter, r *http.Request) error {
	data, header, err := r.FormFile("playlist-image")
	if err != nil {
		return err
	}

	// Deny if greater than 4mb
	if header.Size > 4000000 {
		return &statusError{
			Message: "Image was too large",
			Code:    http.StatusRequestEntityTooLarge,
		}
	}

	opts := &blob.WriterOptions{}
	image, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}

	// Only allow web-safe image files
	actualContentType := http.DetectContentType(image)
	validContentTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}
	valid := false
	for _, contentType := range validContentTypes {
		if actualContentType == contentType {
			valid = true
		}
	}
	if !valid {
		return &statusError{
			Message: fmt.Sprintf("%s is an invalid file type. Try jpeg, jpg, png, or gif", actualContentType),
			Code:    http.StatusBadRequest,
		}
	}

	// I place a random uuid in the image name so that there are no naming collisions
	// The playlistID is in the name to simply relate the image back to the playlist
	id := uuid.New().String()
	playlistID := mux.Vars(r)["id"]
	suffix := strings.Split(actualContentType, "/")[1]

	// the images/ prefix is the target folder inside of the bucket
	imageName := fmt.Sprintf("images/%s-%s.%s", id, playlistID, suffix)
	ctx := context.Background()
	err = p.Bucket.WriteAll(ctx, imageName, image, opts)
	if err != nil {
		return &statusError{
			Message: fmt.Sprintf("Error uploading to S3 bucket: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}
	}

	// write some useful JSON back
	imageCreatedResponse := struct {
		Status   int    `json:"status"`
		ImageURL string `json:"imageURL"`
	}{
		Status: http.StatusCreated,
		// TODO: somehow not hardcode this
		ImageURL: "https://s3.amazonaws.com/stereodose/" + imageName,
	}

	w.WriteHeader(http.StatusCreated)
	util.JSON(w, &imageCreatedResponse)
	return nil
}

// uploadImage actually handles the call to S3 or GCP
func (p *PlaylistsController) uploadImage(img []byte, imageName string) error {
	ctx := context.Background()
	opts := &blob.WriterOptions{}
	err := p.Bucket.WriteAll(ctx, imageName, img, opts)
	if err != nil {
		return &statusError{
			Message: fmt.Sprintf("Error uploading to S3 bucket: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}
	}
	return nil
}
