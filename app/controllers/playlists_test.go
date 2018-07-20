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

func (f *fakePlaylistService) GetPlaylists(offset, limit string) ([]models.Playlist, error) {
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
		Category:    "Weed",
		SubCategory: "Chill",
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
			Category:    "Weed",
			SubCategory: "Chill",
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
