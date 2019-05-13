package controllers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"

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

func TestAuthController_Logout(t *testing.T) {
	testStore := sessions.NewCookieStore([]byte("something-very-secret"))
	testDB := &models.StereoDoseDB{
		Users:     &fakeUserService{},
		Playlists: &fakePlaylistService{},
	}
	req1, _ := http.NewRequest(http.MethodGet, "/auth/logout", nil)
	testCookie, err := generateFakeCookie()
	if err != nil {
		t.Fatal(err.Error())
	}
	res1 := httptest.NewRecorder()
	req1.AddCookie(testCookie)

	type fields struct {
		DB     *models.StereoDoseDB
		Store  *sessions.CookieStore
		Config *oauth2.Config
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "test 1", fields: fields{DB: testDB, Store: testStore}, args: args{w: res1, r: req1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &AuthController{
				DB:     tt.fields.DB,
				Store:  tt.fields.Store,
				Config: tt.fields.Config,
			}
			if err := a.Logout(tt.args.w, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("AuthController.Logout() error = %v, wantErr %v", err, tt.wantErr)
			}
			result := res1.Result()
			cookieHeader := result.Header.Get("Set-Cookie")
			// TODO: parse the data and make sure it's going to be unset
			if len(cookieHeader) == 0 {
				t.Error("Expected Logout to respond with a Set-Cookie header")
			}
		})
	}
}

// literally goes through an entirely fake http cycle to have gorilla sessions
// create a Set-Cookie response
// then parse the header, return a new http.Cookie for testing
// The cookie needs to be created with the same secret, otherwise the digital signatures won't match
func generateFakeCookie() (*http.Cookie, error) {
	testStore := sessions.NewCookieStore([]byte("something-very-secret"))
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		return nil, err
	}
	sess, err := testStore.Get(req, sessionName)
	if err != nil {
		return nil, err
	}
	recorder := httptest.NewRecorder()
	sess.Save(req, recorder)
	response := recorder.Result()
	cookieValue := response.Header.Get("Set-Cookie")
	// cookieValue is in the form: name=value; path=path; expires, date; Max-Age=maxage
	value := strings.Split(cookieValue, " ")[0]
	value = strings.Split(value, "=")[1]
	value = strings.Trim(value, ";")
	log.Println(cookieValue)
	log.Println(value)
	c := &http.Cookie{
		Name:  sessionName,
		Value: value,
	}
	return c, nil
}

func init() {
	log.SetOutput(ioutil.Discard)
}
