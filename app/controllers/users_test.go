package controllers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type fakeUserService struct{}

func (f *fakeUserService) ByID(ID uint) (*models.User, error) {

	like1 := models.Like{UserID: 1, PlaylistID: "2"}
	like2 := models.Like{UserID: 1, PlaylistID: "2"}
	like2.ID = 9000

	if ID == 1 {
		user := &models.User{
			DisplayName: "Test User",
			Likes: []models.Like{
				like1,
				like2,
			},
		}
		user.ID = 1
		return user, nil
	}
	if ID == 9999 {
		return nil, errors.New("Unable to read user from database")
	}
	return nil, nil
}

func (f *fakeUserService) BySpotifyID(ID string) (*models.User, error) {
	return nil, nil
}

func (f *fakeUserService) FirstOrCreate(user *models.User, tok *oauth2.Token) (*models.User, error) {
	return nil, nil
}

func (f *fakeUserService) Update(user *models.User) error {
	return nil
}

func (f *fakeUserService) UpdateAccessToken(user *models.User) error {
	return nil
}

var userTestDB = &models.StereoDoseDB{
	Users: &fakeUserService{},
}

func TestUsersController_Me(t *testing.T) {

	validRequest, err := http.NewRequest(http.MethodGet, "/api/users/me", nil)
	if err != nil {
		t.Fatal(err)
	}
	validUser := models.User{}
	validUser.ID = 1

	invalidUser := models.User{}
	invalidUser.ID = 9999

	tests := []struct {
		name             string
		user             interface{}
		wantErr          bool
		wantResponseBody bool
	}{
		{name: "valid request to /me", wantErr: false, wantResponseBody: true, user: validUser},
		{name: "nil user context", wantErr: true, user: nil},
		{name: "database error", wantErr: true, user: invalidUser},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UsersController{
				DB: userTestDB,
			}

			req, err := http.NewRequest(http.MethodGet, "/api/users/me", nil)
			if err != nil {
				t.Fatal(err)
			}
			ctx := context.WithValue(validRequest.Context(), "User", tt.user)
			req = req.WithContext(ctx)
			recorder := httptest.NewRecorder()
			if err := u.Me(recorder, req); (err != nil) != tt.wantErr {
				t.Errorf("UsersController.Me() error = %v, wantErr %v", err, tt.wantErr)
			}
			resp := recorder.Result()
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected response code to be %d Got: %d", http.StatusOK, resp.StatusCode)
			}

			if tt.wantResponseBody {
				var u models.User
				err := json.NewDecoder(resp.Body).Decode(&u)
				if err != nil {
					t.Error(err)
				}
				if u.ID != 1 {
					t.Errorf("Expected user ID: %d Got: %d", 1, u.ID)
				}
			}
		})
	}
}

func init() {
	log.SetOutput(ioutil.Discard)
}
