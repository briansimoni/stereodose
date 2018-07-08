package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

var brokenPlaylist uint = 9001

// fake playlist service
type fakePlaylistService struct {
}

func (f *fakePlaylistService) GetPlaylists() ([]models.Playlist, error) {
	return nil, nil
}
func (f *fakePlaylistService) GetByID(ID uint) (*models.Playlist, error) {
	if ID == 0 {
		return nil, errors.New("Playlist with ID 0 does not exist")
	}
	playlist := &models.Playlist{
		Name: "Test Playlist",
	}
	return playlist, nil
}

var controller = &PlaylistsController{
	DB: &models.StereoDoseDB{
		Playlists: &fakePlaylistService{},
	},
}

var testRouter = &util.AppRouter{mux.NewRouter()}

func TestPlaylistsController_GetPlaylistByID(t *testing.T) {
	tt := []struct {
		name       string
		value      string
		playlistID string
		status     int
	}{
		{name: "Valid Playlist ID", value: "1", status: http.StatusOK},
		{name: "Playlist ID too long to be a uint", value: "999999999999999999999", status: http.StatusInternalServerError},
		{name: "Invalid Playlist ID", value: "0", status: http.StatusInternalServerError},
		{name: "Negative Number", value: "-1", status: http.StatusNotFound},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			testRouter.AppHandler("/api/playlists/{id:[0-9]+}", controller.GetPlaylistByID)
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
