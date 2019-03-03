package controllers

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/google/go-cloud/blob"
	"github.com/google/go-cloud/blob/driver"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// io.Reader that always errors
type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

type fakeWriteCloser struct {
}

func (w fakeWriteCloser) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, errors.New("unable to write to blob storage")
	}
	return len(p), nil
}
func (w fakeWriteCloser) Close() error {
	return nil
}

// fakeBucket is used to compose the fake controller
type fakeBucket struct{}

func (f fakeBucket) Attributes(ctx context.Context, key string) (driver.Attributes, error) {
	return driver.Attributes{}, errors.New("Attributes not implemented")
}
func (f fakeBucket) NewRangeReader(ctx context.Context, key string, offset, length int64) (driver.Reader, error) {
	return nil, errors.New("NewRangeReader not implemented")
}
func (f fakeBucket) NewTypedWriter(ctx context.Context, key string, contentType string, opt *driver.WriterOptions) (driver.Writer, error) {
	fakeWriter := &fakeWriteCloser{}
	return fakeWriter, nil
}
func (f fakeBucket) Delete(ctx context.Context, key string) error {
	return nil
}
func (f fakeBucket) WriteAll(ctx context.Context, key string, data []byte, options *blob.WriterOptions) error {
	return errors.New("WriteAll not implemented")
}

// fake playlist service
type fakePlaylistService struct {
}

func (f *fakePlaylistService) GetPlaylists(offset, limit, category, subcategory string) ([]models.Playlist, error) {
	off, _ := strconv.Atoi(offset)
	lim, _ := strconv.Atoi(limit)
	if off < 0 || lim < 0 {
		return nil, errors.New("Negative offset or limit")
	}
	return nil, nil
}
func (f *fakePlaylistService) GetByID(ID string) (*models.Playlist, error) {
	if ID == "" {
		return nil, errors.New("Playlist with empty string does not exist")
	}
	if ID == "error-condition" {
		return nil, errors.New("Error reading playlist from DB")
	}
	if ID == "9000" {
		return nil, nil
	}
	playlist := &models.Playlist{
		Name: "Test Playlist",
	}
	return playlist, nil
}

func (f *fakePlaylistService) CreatePlaylistBySpotifyID(user models.User, spotifyID, category, subcategory, image, thumbnailImage string) (*models.Playlist, error) {
	if spotifyID == "alreadyExists" {
		return nil, errors.New("Playlist with this id already exists")
	}
	return nil, nil
}
func (f *fakePlaylistService) GetMyPlaylists(user models.User) ([]models.Playlist, error) {
	if user.DisplayName == "BadTestCase" {
		return nil, errors.New("Unable to obtain playlists because reasons")
	}
	if user.DisplayName == "HasPlaylistsUser1" && user.ID == 1 {
		playlists := []models.Playlist{
			models.Playlist{SpotifyID: "10"},
		}
		return playlists, nil
	}
	if user.DisplayName == "HasPlaylistsUser2" && user.ID == 2 {
		playlists := []models.Playlist{
			models.Playlist{SpotifyID: "20"},
		}
		return playlists, nil
	}
	return nil, nil
}
func (f *fakePlaylistService) DeletePlaylist(id string) error {
	return nil
}

func (f *fakePlaylistService) Comment(playlistID, text string, user models.User) (*models.Comment, error) {
	if text == "leet hacks" {
		return nil, errors.New("wow something broke")
	}
	return nil, nil
}

var controller = &PlaylistsController{
	DB: &models.StereoDoseDB{
		Playlists: &fakePlaylistService{},
	},
	Bucket: blob.NewBucket(fakeBucket{}),
}

func TestPlaylistsController_GetPlaylistByID(t *testing.T) {
	var testRouter = &util.AppRouter{mux.NewRouter()}
	tt := []struct {
		name       string
		value      string
		playlistID string
		status     int
	}{
		{name: "Valid Playlist ID", value: "1", status: http.StatusOK},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			testRouter.AppHandler("/api/playlists/{id}", controller.GetPlaylistByID)
			req, err := http.NewRequest(http.MethodGet, "/api/playlists/"+tc.value, nil)
			if err != nil {
				t.Fatal("Failed to create a request")
			}
			recorder := httptest.NewRecorder()

			// act
			testRouter.ServeHTTP(recorder, req)
			result := recorder.Result()

			// assert
			if result.StatusCode != tc.status {
				t.Errorf("Expected status: %v; Got: %v", tc.status, result.Status)
			}
		})
	}
}

func TestPlaylistsController_GetPlaylists(t *testing.T) {
	var testRouter = &util.AppRouter{mux.NewRouter()}
	tt := []struct {
		name   string
		limit  string
		offset string
		status int
	}{
		{name: "Valid limit and offset", limit: "10", offset: "10", status: http.StatusOK},
		{name: "Invalid limit and offset", limit: "-4", offset: "-9000", status: http.StatusInternalServerError},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			testRouter.AppHandler("/api/playlists/", controller.GetPlaylists).
				Queries("limit", "", "offset", "")

			path := fmt.Sprintf("/api/playlists/?limit=%s&offset=%s", tc.limit, tc.offset)
			req, err := http.NewRequest(http.MethodGet, path, nil)
			if err != nil {
				t.Fatal("Failed to create a request")
			}
			recorder := httptest.NewRecorder()

			// act
			testRouter.ServeHTTP(recorder, req)
			result := recorder.Result()

			// assert
			if result.StatusCode != tc.status {
				t.Errorf("Expected status: %v; Got: %v", tc.status, result.Status)
			}
		})
	}
}

func TestPlaylistsController_CreatePlaylist(t *testing.T) {

	type postBody struct {
		SpotifyID   string
		Category    string
		SubCategory string
	}

	validData := postBody{
		SpotifyID:   "test",
		Category:    "weed",
		SubCategory: "chill",
	}
	var testRouter = &util.AppRouter{mux.NewRouter()}
	tt := []struct {
		name   string
		status int
		user   interface{}
		data   interface{}
	}{
		{name: "Valid ID", status: 201, user: models.User{}, data: validData},
		{name: "Invalid Categories", status: 400, user: nil, data: postBody{"test", "Fake", "Category"}},
		{name: "Invalid User Context", status: 500, user: nil, data: validData},
		{name: "Invalid POST body", status: 400, user: models.User{}, data: 69},
		{name: "Database Error", status: 500, user: models.User{}, data: postBody{
			SpotifyID:   "alreadyExists",
			Category:    "weed",
			SubCategory: "chill",
		},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			testRouter.AppHandler("/api/playlists/", controller.CreatePlaylist).Methods(http.MethodPost)
			body, _ := json.Marshal(tc.data)
			req, err := http.NewRequest(http.MethodPost, "/api/playlists/", bytes.NewBuffer(body))
			if err != nil {
				t.Error(err.Error())
			}
			recorder := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), "User", tc.user)
			testRouter.ServeHTTP(recorder, req.WithContext(ctx))
			result := recorder.Result()

			if result.StatusCode != tc.status {
				t.Errorf("Expected status: %v; Got: %v", tc.status, result.Status)
			}
		})
	}
}

func TestPlaylistsController_GetMyPlaylists(t *testing.T) {
	var testRouter = &util.AppRouter{mux.NewRouter()}
	tt := []struct {
		name   string
		status int
		user   *models.User
	}{
		{name: "Valid Test", status: 200, user: &models.User{}},
		{name: "User Missing", status: 500, user: nil},
		{name: "Database Error", status: 500, user: &models.User{DisplayName: "BadTestCase"}},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			testRouter.AppHandler("/api/playlists/me", controller.GetMyPlaylists).Methods(http.MethodGet)
			req, err := http.NewRequest(http.MethodGet, "/api/playlists/me", nil)
			if err != nil {
				t.Error(err.Error())
			}

			user := tc.user
			var ctx context.Context
			if tc.user != nil {
				ctx = context.WithValue(req.Context(), "User", *user)
			} else {
				ctx = req.Context()
			}
			recorder := httptest.NewRecorder()

			testRouter.ServeHTTP(recorder, req.WithContext(ctx))
			result := recorder.Result()

			if result.StatusCode != tc.status {
				t.Errorf("Expected status: %v; Got: %v", tc.status, result.Status)
			}
		})
	}
}

func TestPlaylistsController_DeletePlaylist(t *testing.T) {
	var testRouter = &util.AppRouter{mux.NewRouter()}

	user1 := models.User{}
	user1.ID = 1
	user1.DisplayName = "HasPlaylistsUser1"

	user2 := models.User{}
	user2.ID = 2
	tests := []struct {
		name       string
		user       interface{}
		playlistID string
		statusCode int
	}{
		{name: "authorized delete", user: user1, playlistID: "10", statusCode: 200},
		{name: "unauthorized delete", user: user1, playlistID: "20", statusCode: 401},
		{name: "noexistent playlist", user: user1, playlistID: "9000", statusCode: 404},
		{name: "bad session cookie", user: nil, playlistID: "10", statusCode: 500},
		{name: "empty playlist id", user: user1, playlistID: "error-condition", statusCode: 500},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testRouter.AppHandler("/api/playlists/{id}", controller.DeletePlaylist).Methods(http.MethodDelete)
			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodDelete, "/api/playlists/"+tc.playlistID, nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(context.Background(), "User", tc.user)
			req = req.WithContext(ctx)
			testRouter.ServeHTTP(recorder, req)
			result := recorder.Result()

			if result.StatusCode != tc.statusCode {
				t.Errorf("Expected status code: %d, Got: %d", tc.statusCode, result.StatusCode)
			}
		})
	}
}

func TestPlaylistsController_uploadImage(t *testing.T) {
	var fakeImageData = make([]byte, 100)
	_, err := rand.Read(fakeImageData)
	if err != nil {
		t.Fatal("Failed to create random fake image ", err.Error())
	}

	type args struct {
		img       []byte
		imageName string
	}
	tests := []struct {
		name    string
		p       *PlaylistsController
		args    args
		wantErr bool
	}{
		{name: "Normal Image", p: controller, args: args{img: fakeImageData, imageName: "playlist-image.jpeg"}, wantErr: false},
		{name: "No image data", p: controller, args: args{img: []byte{}, imageName: ""}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.uploadImage(tt.args.img, tt.args.imageName); (err != nil) != tt.wantErr {
				t.Errorf("PlaylistsController.uploadImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPlaylistsController_UploadImage(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	req, err := newMultiPartUploadRequest()
	if err != nil {
		t.Fatal("Error creating multi-part upload request", err.Error())
	}

	tests := []struct {
		name    string
		p       *PlaylistsController
		args    args
		wantErr bool
	}{
		{name: "normal run", p: controller, args: args{w: httptest.NewRecorder(), r: req}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.UploadImage(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("PlaylistsController.UploadImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func newMultiPartUploadRequest() (*http.Request, error) {
	// Create an 100 x 50 image
	img := image.NewRGBA(image.Rect(0, 0, 100, 50))
	xMax := img.Bounds().Dx()
	yMax := img.Bounds().Dy()
	for x := 0; x < xMax; x++ {
		for y := 0; y < yMax; y++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	defer writer.Close()
	formFile, err := writer.CreateFormFile("playlist-image", "image.jpeg")
	jpeg.Encode(formFile, img, nil)

	req, err := http.NewRequest(http.MethodPost, "/api/playlists/1/image", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func TestPlaylistsController_Comment(t *testing.T) {

	testUser := models.User{}
	testUser.ID = 1
	testBody1 := bytes.NewReader([]byte(`{"text": "cool playlist bro"}`))
	request1, err := http.NewRequest(http.MethodPost, "/api/playlists/1234/comment", testBody1)
	if err != nil {
		t.Fatal(err.Error())
	}
	ctx1 := context.WithValue(request1.Context(), "User", testUser)
	request1 = request1.WithContext(ctx1)

	request2, err := http.NewRequest(http.MethodPost, "/api/playlists/1234/comment", testBody1)
	if err != nil {
		t.Fatal(err.Error())
	}

	request3, err := http.NewRequest(http.MethodPost, "/api/playlists/1234/comment", errReader(0))
	if err != nil {
		t.Fatal(err.Error())
	}
	ctx3 := context.WithValue(request3.Context(), "User", testUser)
	request3 = request3.WithContext(ctx3)

	testBody4 := bytes.NewReader([]byte(`not valid json{} [] paisdjfpoj ""`))
	request4, err := http.NewRequest(http.MethodPost, "/api/playlists/1234/comment", testBody4)
	if err != nil {
		t.Fatal(err.Error())
	}
	ctx4 := context.WithValue(request4.Context(), "User", testUser)
	request4 = request4.WithContext(ctx4)

	testBody5 := bytes.NewReader([]byte(`{"text": ""}`))
	request5, err := http.NewRequest(http.MethodPost, "/api/playlists/1234/comment", testBody5)
	if err != nil {
		t.Fatal(err.Error())
	}
	ctx5 := context.WithValue(request5.Context(), "User", testUser)
	request5 = request5.WithContext(ctx5)

	testBody6 := bytes.NewReader([]byte(`{"text": "leet hacks"}`))
	request6, err := http.NewRequest(http.MethodPost, "/api/playlists/1234/comment", testBody6)
	if err != nil {
		t.Fatal(err.Error())
	}
	ctx6 := context.WithValue(request6.Context(), "User", testUser)
	request6 = request6.WithContext(ctx6)

	mux.SetURLVars(request1, map[string]string{"id": "1234"})

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	tests := []struct {
		name    string
		p       *PlaylistsController
		args    args
		wantErr bool
	}{
		{name: "normal test", p: controller, args: args{w: httptest.NewRecorder(), r: request1}, wantErr: false},
		{name: "no authenticated user", p: controller, args: args{w: httptest.NewRecorder(), r: request2}, wantErr: true},
		{name: "something wrong with the body", p: controller, args: args{w: httptest.NewRecorder(), r: request3}, wantErr: true},
		{name: "invalid json", p: controller, args: args{w: httptest.NewRecorder(), r: request4}, wantErr: true},
		{name: "empty comment", p: controller, args: args{w: httptest.NewRecorder(), r: request5}, wantErr: true},
		{name: "database error", p: controller, args: args{w: httptest.NewRecorder(), r: request6}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.Comment(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("PlaylistsController.Comment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
