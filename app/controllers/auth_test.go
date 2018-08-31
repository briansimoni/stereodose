package controllers

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/config"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func TestNewAuthController(t *testing.T) {
	testDB := &models.StereoDoseDB{}
	store := &sessions.CookieStore{}
	config := &config.Config{
		ClientID:     "testclient",
		ClientSecret: "secret",
	}

	got := NewAuthController(testDB, store, config)

	if got.Config.ClientID != "testclient" {
		t.Errorf("Expected client id: %s. Got: %s", config.ClientID, got.Config.ClientID)
	}
	if got.Config.ClientSecret != "secret" {
		t.Errorf("Expected client secret: %s. Got: %s", config.ClientSecret, got.Config.ClientSecret)
	}
	if !reflect.DeepEqual(testDB, got.DB) {
		t.Error("The controller's database was not equivalent to what was passed in")
	}
	if !reflect.DeepEqual(store, got.Store) {
		t.Error("The controller's cookie store was not equivalent to what was passed in")
	}
}

var testAuthController = &AuthController{
	// will need to provide fake services to the DB for testing later
	DB:    &models.StereoDoseDB{},
	Store: sessions.NewCookieStore([]byte("something-very-secret")),
	Config: &oauth2.Config{
		ClientID:     "test-client",
		ClientSecret: "secret",
		RedirectURL:  "http://localhost:4000/auth/callback",
		Scopes: []string{
			"some-fake-scope",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://fake-endpoint.com/authorization",
			TokenURL: "https://fake-endpoint.com/token",
		},
	},
}

func TestAuthController_Login(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	firstTestRequest, err := http.NewRequest(http.MethodGet, "/auth/login", nil)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Normal Login Request", wantErr: false, args: args{httptest.NewRecorder(), firstTestRequest}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testAuthController.Login(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("AuthController.Login() error = %v, wantErr %v", err, tt.wantErr)
			}

			recorder, _ := tt.args.w.(*httptest.ResponseRecorder)
			response := recorder.Result()
			if response.StatusCode != http.StatusTemporaryRedirect {
				t.Errorf("Expected %d, got %d", http.StatusTemporaryRedirect, response.StatusCode)
			}

			location, err := response.Location()
			if err != nil {
				t.Fatal(err)
			}
			if location.Path != "/authorization" {
				t.Errorf("Expected redirect location to be: /authorization Got: %s", location.Path)
			}
		})
	}
}
