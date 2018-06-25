package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	endpoint "golang.org/x/oauth2/spotify"

	"github.com/gorilla/sessions"
)

const sessionName = "_stereodose-session"

var spotifyURL = "https://accounts.spotify.com"

type AuthController struct {
	DB    *models.StereoDoseDB
	Store *sessions.CookieStore
}

type refreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// TODO: get this from the app config struct and not os env
var conf = &oauth2.Config{
	ClientID:     os.Getenv("STEREODOSE_CLIENT_ID"),
	ClientSecret: os.Getenv("STEREODOSE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("STEREODOSE_REDIRECT_URL"),
	Scopes: []string{
		"playlist-modify-public",
		"streaming",
		"user-read-birthdate",
		"user-read-email",
		"user-read-private",
		"playlist-read-collaborative",
		"user-modify-playback-state",
	},
	Endpoint: endpoint.Endpoint,
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) error {
	s, err := a.Store.Get(r, sessionName)
	if err != nil {
		return err
	}

	if s.Values["Access_Token"] == nil {
		// user is not logged in. send to authorization code flow
		// Redirect user to consent page to ask for permission
		// for the specified scopes.

		b := make([]byte, 32)
		_, err = rand.Read(b)
		if err != nil {
			return err
		}
		state := base64.StdEncoding.EncodeToString(b)
		s.Values["State"] = state
		s.Save(r, w)

		url := conf.AuthCodeURL(state, oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return nil
	}
	tok, ok := s.Values["Token"].(oauth2.Token)
	if !ok {
		return errors.New("Unable to obtain token from session")
	}
	if !tok.Valid() {
		http.Redirect(w, r, "/auth/refresh", http.StatusTemporaryRedirect)
		return nil
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return nil

}

func (a *AuthController) Callback(w http.ResponseWriter, r *http.Request) error {
	s, err := a.Store.Get(r, sessionName)
	if err != nil {
		return err
	}
	err = checkState(r, s)
	if err != nil {
		return err
	}

	tok, err := conf.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		return errors.New("Error obtaining token from Spotify: " + err.Error())
	}

	s.Values["Token"] = *tok
	client := spotify.Authenticator{}.NewClient(tok)
	currentUser, err := client.CurrentUser()
	if err != nil {
		return err
	}
	sdUser, err := a.saveUserData(tok, currentUser)
	if err != nil {
		return err
	}

	s.Values["User_ID"] = sdUser.ID
	err = s.Save(r, w)
	if err != nil {
		return err
	}
	returnPath, ok := s.Values["return_path"].(string)
	if !ok {
		returnPath = "/"
	}
	http.Redirect(w, r, returnPath, http.StatusTemporaryRedirect)
	return nil
}

// Refresh will update the Spotify API Access Token for the user's session
// TODO: check the refresh token and save it (it might be a new refresh token)
func (a *AuthController) Refresh(w http.ResponseWriter, r *http.Request) error {
	user, ok := r.Context().Value("User").(models.User)
	if !ok {
		return errors.New("unable to obtain user from context")
	}
	s, err := a.Store.Get(r, sessionName)
	if err != nil {
		return err
	}
	tok, err := refreshToken(user.RefreshToken)
	if err != nil {
		return err
	}
	user.AccessToken = tok.AccessToken
	err = a.DB.Users.Update(&user)
	if err != nil {
		return err
	}
	sessionToken, ok := s.Values["Token"].(oauth2.Token)
	if !ok {
		return errors.New("Error reading OAuth Token from Session")
	}
	sessionToken.AccessToken = tok.AccessToken
	sessionToken.Expiry = sessionToken.Expiry.Add(time.Duration(tok.ExpiresIn) * time.Second)
	s.Values["Token"] = sessionToken
	log.Println(sessionToken)
	err = s.Save(r, w)
	if err != nil {
		return err
	}
	j, err := json.MarshalIndent(&tok, " ", " ")
	if err != nil {
		return err
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		return err
	}
	return nil
}

func refreshToken(refreshToken string) (*refreshTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	body := strings.NewReader(data.Encode())
	creds := conf.ClientID + ":" + conf.ClientSecret
	basicAuth := base64.StdEncoding.EncodeToString([]byte(creds))
	req, err := http.NewRequest(http.MethodPost, spotifyURL+"/api/token", body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+basicAuth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("Response from spotify.com/api/token " + res.Status)
	}
	defer res.Body.Close()
	var tok refreshTokenResponse
	err = json.NewDecoder(res.Body).Decode(&tok)
	if err != nil {
		return nil, err
	}
	return &tok, nil
}

func checkState(r *http.Request, s *sessions.Session) error {
	responseState := r.URL.Query().Get("state")
	if responseState == "" {
		return errors.New("Unable to obtain state from URL query params")
	}

	state := s.Values["State"]
	if state == "" {
		return errors.New("Unable to obtain state from session")
	}

	if r.URL.Query().Get("state") != state {
		return errors.New("State from query params did not match session state")
	}
	return nil
}

func (a *AuthController) saveUserData(token *oauth2.Token, u *spotify.PrivateUser) (*models.User, error) {
	user := &models.User{
		Birthdate:    u.Birthdate,
		DisplayName:  u.DisplayName,
		Email:        u.Email,
		SpotifyID:    u.ID,
		Images:       u.Images,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	user, err := a.DB.Users.FirstOrCreate(user, token)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Middleware checks to see if the user is logged in before
// allowing the request to continue
func (a *AuthController) Middleware(next http.HandlerFunc) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		s, err := a.Store.Get(r, sessionName)
		if err != nil {
			return err
		}
		tok, ok := s.Values["Token"].(oauth2.Token)
		if !ok || !tok.Valid() {
			s.Values["return_path"] = r.URL.Path
			s.Save(r, w)
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return nil
		}
		next.ServeHTTP(w, r)
		return nil
	}
	return util.Handler{f}
}
