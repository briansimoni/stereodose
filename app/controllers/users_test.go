package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

// type UserService interface {
// 	ByID(ID uint) (*User, error)
// 	FirstOrCreate(user *User, tok *oauth2.Token) (*User, error)
// 	Update(user *User) error
// }

type fakeUserService struct{}

func (f *fakeUserService) ByID(ID uint) (*models.User, error) {
	if ID == 1 {
		user := &models.User{
			DisplayName: "Test User",
		}
		user.ID = 1
		return user, nil
	}
	if ID == 9999 {
		return nil, errors.New("Unable to read user from database")
	}
	return nil, nil
}

func (f *fakeUserService) FirstOrCreate(user *models.User, tok *oauth2.Token) (*models.User, error) {
	return nil, nil
}

func (f *fakeUserService) Update(user *models.User) error {
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