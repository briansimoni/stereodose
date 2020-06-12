package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	// register jpeg type
	"image/jpeg"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/disintegration/imaging"
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

// GetPlaylists will return a subset of all the playlists in the DB based on query parameters
func (p *PlaylistsController) GetPlaylists(w http.ResponseWriter, r *http.Request) error {
	params, err := createSearchParamsFromRequest(r)
	if err != nil {
		return &util.StatusError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	playlists, err := p.DB.Playlists.GetPlaylists(params)
	if err != nil {
		return errors.WithStack(err)
	}
	util.JSON(w, playlists)
	return nil
}

func createSearchParamsFromRequest(r *http.Request) (*models.PlaylistSearchParams, error) {
	queryValues := r.URL.Query()

	offset := queryValues.Get("offset")
	limit := queryValues.Get("limit")
	category := queryValues.Get("category")
	subcategory := queryValues.Get("subcategory")
	sortKey := queryValues.Get("sort-key")
	order := queryValues.Get("order")
	IDs := queryValues.Get("spotify-ids")

	if offset == "" {
		offset = "0"
	}
	if limit == "" {
		limit = "10"
	}
	if sortKey == "" {
		sortKey = "likes_count"
	}
	if sortKey == "name" && order == "" {
		order = "asc"
	}
	if order == "" {
		order = "desc"
	}

	sortKeyAllowed := false
	for _, allowedSortKey := range []string{
		"created_at",
		"updated_at",
		"name",
		"likes_count",
	} {
		if sortKey == allowedSortKey {
			sortKeyAllowed = true
			break
		}
	}
	if !sortKeyAllowed {
		return nil, errors.Errorf("sort key: '%s' is not allowed", sortKey)
	}

	orderValueAllowed := false
	for _, allowedOrderValue := range []string{
		"desc",
		"asc",
	} {
		if order == allowedOrderValue {
			orderValueAllowed = true
			break
		}
	}
	if !orderValueAllowed {
		return nil, errors.Errorf("order value: '%s' is not allowed", order)
	}

	spotifyIds := strings.Split(IDs, " ")
	// Split always returns a string even if its empty
	// In the case that it is empty, we'll set the slice to be nil
	if spotifyIds[0] == "" {
		spotifyIds = nil
	}

	return &models.PlaylistSearchParams{
		Category:    category,
		Subcategory: subcategory,
		Offset:      offset,
		Limit:       limit,
		SortKey:     sortKey,
		Order:       order,
		SpotifyIDs:  spotifyIds,
	}, nil
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

// GetRandomPlaylist will inspect the category and subcategory query parameters
// and return a Playlist filled with randomly selected tracks
func (p *PlaylistsController) GetRandomPlaylist(w http.ResponseWriter, r *http.Request) error {
	queryValues := r.URL.Query()
	category := queryValues.Get("category")
	subcategory := queryValues.Get("subcategory")

	if !models.Categories.Valid(category, subcategory) {
		return &util.StatusError{
			Message: fmt.Sprintf("Category: %s Subcategory: %s is invalid", category, subcategory),
			Code:    http.StatusBadRequest,
		}
	}

	playlist, err := p.DB.Playlists.GetRandomPlaylist(category, subcategory)
	if err != nil {
		return err
	}

	err = util.JSON(w, playlist)
	if err != nil {
		return err
	}
	return nil
}

// CreatePlaylist reads the SpotifyID from the POST body
// It then calls the spotify API to get the full info and store in the local DB
func (p *PlaylistsController) CreatePlaylist(w http.ResponseWriter, r *http.Request) error {
	type jsonBody struct {
		SpotifyID    string `json:"SpotifyID"`
		Category     string `json:"Category"`
		SubCategory  string `json:"SubCategory"`
		ImageURL     string `json:"ImageURL"`
		ThumbnailURL string `json:"ThumbnailURL"`
	}
	var data jsonBody
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return &util.StatusError{
			Message: fmt.Sprintf("Error parsing JSON: %s", err.Error()),
			Code:    http.StatusBadRequest,
		}
	}
	valid := models.Categories.Valid(data.Category, data.SubCategory)
	if !valid {
		return &util.StatusError{
			Message: fmt.Sprintf("Invalid Category/Subcategory: %s / %s", data.Category, data.SubCategory),
			Code:    http.StatusBadRequest,
		}
	}
	user, ok := r.Context().Value("User").(models.User)
	if !ok {
		return errors.New("Unable to obtain user from session")
	}

	existingPlaylist, _ := p.DB.Playlists.GetByID(data.SpotifyID)
	if existingPlaylist != nil {
		return &util.StatusError{
			Message: fmt.Sprintf("Playlist with ID %s already exists", existingPlaylist.SpotifyID),
			Code:    http.StatusConflict,
		}
	}

	playlist, err := p.DB.Playlists.CreatePlaylistBySpotifyID(user, data.SpotifyID, data.Category, data.SubCategory, data.ImageURL, data.ThumbnailURL)
	if err != nil {
		return errors.WithStack(err)
	}
	log.WithFields(logrus.Fields{
		"Type":          "AppLog",
		"TransactionID": r.Context().Value(util.TransactionIDKey),
		"Name":          playlist.Name,
		"Category":      playlist.Category,
		"SubCategory":   playlist.SubCategory,
		"UserID":        playlist.UserID,
	}).Info("Successfully Created New Playlist")
	w.WriteHeader(http.StatusCreated)
	return nil
}

// DeletePlaylist takes the id variable from the url path
// It performs a hard delete of the playlist from the DB, but a soft delete on likes and comments
// Only admins may perform this operation
func (p *PlaylistsController) DeletePlaylist(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	ID := vars["id"]
	user, ok := r.Context().Value("User").(models.User)
	if !ok {
		return errors.New("Unable to obtain user from session")
	}
	if !user.Admin {
		return &util.StatusError{
			Message: fmt.Sprintf("Unauthorized"),
			Code:    http.StatusUnauthorized,
		}
	}
	targetPlaylist, err := p.DB.Playlists.GetByID(ID)
	if err != nil {
		return errors.WithStack(err)
	}
	if targetPlaylist == nil {
		return &util.StatusError{
			Message: fmt.Sprintf("Playlist does not exist"),
			Code:    http.StatusNotFound,
		}
	}

	err = p.DB.Playlists.DeletePlaylist(ID)
	if err != nil {
		return errors.WithStack(err)
	}

	// not really a big deal if the images don't delete
	// thus, these functions run concurrently and are not checked for errors
	go p.Bucket.Delete(context.Background(), getImageKey(targetPlaylist.BucketImageURL))
	go p.Bucket.Delete(context.Background(), getImageKey(targetPlaylist.BucketThumbnailURL))
	return nil
}

// getImageKey takes an S3 (or generic bucket URL) and returns just the key
func getImageKey(url string) string {
	if url == "" {
		return ""
	}
	split := strings.Split(url, "/")
	return split[len(split)-2] + "/" + split[len(split)-1]
}

// UploadImage saves an image the corresponds to a playlist
// The actual data is saved to cloud bucket storage
// Permalinks are returned to the client. The following request (CreatePlaylist) saves the playlist
// along with the permalinks to the database.
func (p *PlaylistsController) UploadImage(w http.ResponseWriter, r *http.Request) error {
	multipartFile, header, err := r.FormFile("playlist-image")
	if err != nil {
		return errors.WithStack(err)
	}

	// Deny if greater than 8mb
	if header.Size > 8000000 {
		return &util.StatusError{
			Message: "Image was too large",
			Code:    http.StatusRequestEntityTooLarge,
		}
	}

	// here we create a space in memory to copy the image
	buffer := new(bytes.Buffer)
	// we use tee reader so I can ioutil.ReadAll, and then read again from buffer later
	reader := io.TeeReader(multipartFile, buffer)
	imageCopy, err := jpeg.Decode(reader)
	if err != nil {
		return errors.WithStack(err)
	}
	resizedImage := imaging.Resize(imageCopy, 250, 250, imaging.Lanczos)
	imageDataCopy := new(bytes.Buffer)
	err = jpeg.Encode(imageDataCopy, resizedImage, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	imageData, err := ioutil.ReadAll(buffer)
	if err != nil {
		return errors.WithStack(err)
	}

	// Only allow web-safe image files
	actualContentType := http.DetectContentType(imageData)
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
		return &util.StatusError{
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
	thumbNailName := fmt.Sprintf("images/%s-%s-thumbnail.%s", id, playlistID, suffix)

	// upload the images in parallel for a small performance boost
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		err = p.uploadImage(imageData, imageName)
		wg.Done()
	}()
	go func() {
		err = p.uploadImage(imageDataCopy.Bytes(), thumbNailName)
		wg.Done()
	}()
	wg.Wait()
	if err != nil {
		return err
	}

	// write some useful JSON back
	imageCreatedResponse := struct {
		Status       int    `json:"status"`
		ImageURL     string `json:"imageURL"`
		ThumbNailURL string `json:"thumbnailURL"`
	}{
		Status: http.StatusCreated,
		// TODO: somehow not hardcode the base URL
		// The "images" part of the imageName gets removed here since we want the cloudfront URL
		// cloudfront base URL starts in the images directory of S3
		// For example s3.stereodose.com/images/image1 === images.stereodose.app/image1
		// It is like this because cloudfront wasn't originally part of the architecture
		ImageURL:     "https://images.stereodose.app/" + strings.Replace(imageName, "images/", "", 1),
		ThumbNailURL: "https://images.stereodose.app/" + strings.Replace(thumbNailName, "images/", "", 1),
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
		return &util.StatusError{
			Message: fmt.Sprintf("Error uploading to S3 bucket: %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}
	}
	return nil
}

// Comment parses the JSON body and saves a user comment to the database
// the JSON body looks like this: { "text": "wow this playlist is cool" }
func (p *PlaylistsController) Comment(w http.ResponseWriter, r *http.Request) error {
	user, ok := r.Context().Value("User").(models.User)
	if !ok {
		return errors.New("Unable to obtain user from session")
	}
	vars := mux.Vars(r)
	playlistID := vars["id"]

	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return &util.StatusError{
			Message: fmt.Sprintf("Error parsing body %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}
	}

	type model struct {
		Text string
	}
	var m = new(model)
	err = json.Unmarshal(data, m)
	if err != nil {
		return &util.StatusError{
			Message: fmt.Sprintf("Error parsing body %s", err.Error()),
			Code:    http.StatusInternalServerError,
		}
	}

	if m.Text == "" {
		return &util.StatusError{
			Message: "Cannot upload empty comment",
			Code:    http.StatusBadRequest,
		}
	}

	// escape user data to avoid XSS attacks
	// ...or maybe react is handling this
	// escapedText := html.EscapeString(m.Text)

	comment, err := p.DB.Playlists.Comment(playlistID, m.Text, user)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	err = util.JSON(w, comment)
	return err
}

// DeleteComment removes a comment from a playlist and soft deletes in the database
func (p *PlaylistsController) DeleteComment(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["commentID"])
	if err != nil {
		return &util.StatusError{
			Message: "Unable to parse comment ID: " + err.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	user, ok := r.Context().Value("User").(models.User)
	if !ok {
		return &util.StatusError{
			Message: "Unable to obtain user from session",
			Code:    http.StatusUnauthorized,
		}
	}

	comment, err := p.DB.Comments.ByID(uint(commentID))
	if err != nil {
		return &util.StatusError{
			Message: "Unable to read comment from database: " + err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}

	if comment.UserID != user.ID {
		return &util.StatusError{
			Message: "Not authorized - unable to delete other users' playlists",
			Code:    http.StatusForbidden,
		}
	}

	err = p.DB.Playlists.DeleteComment(uint(commentID))
	if err != nil {
		return &util.StatusError{
			Message: "Error deleting comment from database: " + err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	return nil
}

// Like will add a like to the playlist in the database
// Like checks to see if the user had already liked the playlist
func (p *PlaylistsController) Like(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	playlistID := vars["id"]

	u, ok := r.Context().Value("User").(models.User)
	if !ok {
		return &util.StatusError{
			Message: "Unable to obtain user from session",
			Code:    http.StatusUnauthorized,
		}
	}

	// lets make sure that we have the full and updated data set for the user
	user, err := p.DB.Users.ByID(u.ID)
	if err != nil {
		return err
	}

	for _, userLike := range user.Likes {
		if userLike.PlaylistID == playlistID {
			return &util.StatusError{
				Message: "The user has already liked this playlist",
				Code:    http.StatusConflict,
			}
		}
	}

	like, err := p.DB.Playlists.Like(playlistID, *user)
	if err != nil {
		return &util.StatusError{
			Message: "Error writing to database: " + err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}

	w.WriteHeader(http.StatusCreated)
	err = util.JSON(w, like)
	if err != nil {
		return &util.StatusError{
			Message: "Failed to write JSON " + err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	return nil

}

// Unlike removes a like from a playlist
// TODO: could improve the performance of this by looking up the Like by ID instead of
// searching through all of the user's likes
func (p *PlaylistsController) Unlike(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	playlistID := vars["playlistID"]
	likeID, err := strconv.Atoi(vars["likeID"])
	if err != nil {
		return &util.StatusError{
			Message: "Playlist ID is not an integer. Error: " + err.Error(),
			Code:    http.StatusBadRequest,
		}
	}

	u, ok := r.Context().Value("User").(models.User)
	if !ok {
		return &util.StatusError{
			Message: "Unable to obtain user from session",
			Code:    http.StatusUnauthorized,
		}
	}

	// make sure we have the full and updated user data
	user, err := p.DB.Users.ByID(u.ID)
	if err != nil {
		return err
	}

	authorized := false
	for _, like := range user.Likes {
		if like.ID == uint(likeID) {
			authorized = true
			break
		}
	}

	if !authorized {
		return &util.StatusError{
			Message: "The user does not own this like",
			Code:    http.StatusForbidden,
		}
	}

	err = p.DB.Playlists.Unlike(playlistID, uint(likeID))
	if err != nil {
		return &util.StatusError{
			Message: "Database Error " + err.Error(),
			Code:    http.StatusInternalServerError,
		}
	}
	return nil
}
