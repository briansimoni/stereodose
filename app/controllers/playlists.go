package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/google/go-cloud/blob"
	"github.com/google/go-cloud/blob/s3blob"
)

// PlaylistsController is a collection of RESTful Handlers for Playlists
type PlaylistsController struct {
	DB *models.StereoDoseDB
}

// NewPlaylistsController returns a pointer to PlaylistsController
func NewPlaylistsController(db *models.StereoDoseDB) *PlaylistsController {
	return &PlaylistsController{DB: db}
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

	_, err = p.DB.Playlists.CreatePlaylistBySpotifyID(user, data.SpotifyID, data.Category, data.SubCategory)
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
// A reference to the bucket is stored in the database and returned to the client
func (p *PlaylistsController) UploadImage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	log.Println(vars["id"])
	data, header, err := r.FormFile("playlist-image")
	if err != nil {
		return err
	}
	// file, err := os.OpenFile("upload.jpg", os.O_RDWR|os.O_CREATE, 755)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()
	log.Println("size", header.Size)
	log.Println("header", header.Header)
	log.Println("filename", header.Filename)
	// _, err = io.Copy(file, data)
	_, err = io.Copy(os.Stdout, data)
	if err != nil {
		log.Println("zomg err", err.Error())
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}

func gimmeBucket() {
	ctx := context.Background()
	// Open a connection to the bucket.
	var (
		b   *blob.Bucket
		err error
	)
	cloud := "aws"
	bucketName := "stereodose"
	switch cloud {
	case "gcp":
		b, err = setupGCP(ctx, bucketName)
	case "aws":
		// AWS is handled below in the next code sample.
		b, err = setupAWS(ctx, bucketName)
	default:
		log.Fatalf("Failed to recognize cloud. Want gcp or aws, got: %s", cloud)
	}
	if err != nil {
		log.Fatalf("Failed to setup bucket: %s", err)
	}
	log.Println(b)
}

func setupAWS(ctx context.Context, bucket string) (*blob.Bucket, error) {
	c := &aws.Config{
		// Either hard-code the region or use AWS_REGION.
		Region: aws.String("us-east-2"),
		// credentials.NewEnvCredentials assumes two environment variables are
		// present:
		// 1. AWS_ACCESS_KEY_ID, and
		// 2. AWS_SECRET_ACCESS_KEY.
		Credentials: credentials.NewEnvCredentials(),
	}
	s := session.Must(session.NewSession(c))
	return s3blob.OpenBucket(ctx, s, bucket)
}

func setupGCP(ctx context.Context, bucket string) (*blob.Bucket, error) {
	return nil, nil
}
