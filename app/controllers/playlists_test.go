package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

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

func (f *fakePlaylistService) CreatePlaylistBySpotifyID(user models.User, spotifyID, category, subcategory string) (*models.Playlist, error) {
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

var controller = &PlaylistsController{
	DB: &models.StereoDoseDB{
		Playlists: &fakePlaylistService{},
	},
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
