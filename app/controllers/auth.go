package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/briansimoni/stereodose/config"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	endpoint "golang.org/x/oauth2/spotify"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

const (
	spotifyURL = "https://accounts.spotify.com"

	// session keys
	sessionName = "_stereodose-session"
	token       = "Token"
	state       = "State"
	userID      = "User_ID"
	returnPath  = "return_path"
)

// AuthController is a collection of RESTful Handlers for authentication
type AuthController struct {
	DB     *models.StereoDoseDB
	Store  *sessions.CookieStore
	Config *oauth2.Config
}

// NewAuthController takes a StereodoseDB, CookieStore, and App Config
// returns an AuthController
func NewAuthController(db *models.StereoDoseDB, store *sessions.CookieStore, config *config.Config) *AuthController {
	oauthConfig := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectURL,
		Scopes: []string{
			"playlist-modify-public",
			"streaming",
			"user-read-birthdate",
			"user-read-email",
			"user-read-private",
			"playlist-read-collaborative",
			//"playlist-read-private",
			"user-modify-playback-state",
		},
		Endpoint: endpoint.Endpoint,
	}
	a := &AuthController{
		DB:     db,
		Store:  store,
		Config: oauthConfig,
	}
	return a
}

type refreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// Login is the handler that you can send the user to initiate an authorization code flow
func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) error {

	s, err := a.Store.Get(r, sessionName)
	if err != nil {
		return errors.WithStack(err)
	}
	// build a cryptographically random state to prevent csrf attacks
	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return errors.WithStack(err)
	}
	state := base64.StdEncoding.EncodeToString(b)
	s.Values["State"] = state
	s.Save(r, w)

	// If we are behind a proxy, we dynamically grab the port based on the X-Forwarded-Port header
	// This support more diverse cloud deployments without having to add more configuration
	if r.Header.Get("X-Forwarded-Port") != "" && r.Header.Get("X-Forwarded-Port") != "443" {
		port := r.Header.Get("X-Forwarded-Port")
		redirectURL, err := url.Parse(a.Config.RedirectURL)
		if err != nil {
			return err
		}
		redirect := fmt.Sprintf("%s://%s:%s%s", redirectURL.Scheme, redirectURL.Host, port, redirectURL.RequestURI())
		copiedConfig := &oauth2.Config{}
		err = copier.Copy(copiedConfig, a.Config)
		if err != nil {
			return err
		}
		copiedConfig.RedirectURL = redirect

		redir := copiedConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
		http.Redirect(w, r, redir, http.StatusTemporaryRedirect)
		return nil
	}

	redir := a.Config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, redir, http.StatusTemporaryRedirect)
	return nil
}

// Callback is a handler function that the user is redirected to in the OAuth flow
// In this step of authorization, we exchange a code for an access token
// and we query the user's profile on Spotify to get their identity
func (a *AuthController) Callback(w http.ResponseWriter, r *http.Request) error {
	s, err := a.Store.Get(r, sessionName)
	if err != nil {
		return errors.WithStack(err)
	}
	err = checkState(r, s)
	if err != nil {
		return err
	}

	// if behind load balancer, dynamically check the port so we build the correct redirect uri
	var tok *oauth2.Token
	if r.Header.Get("X-Forwarded-Port") != "" && r.Header.Get("X-Forwarded-Port") != "443" {
		port := r.Header.Get("X-Forwarded-Port")
		redirectURL, err := url.Parse(a.Config.RedirectURL)
		if err != nil {
			return err
		}
		redirect := fmt.Sprintf("%s://%s:%s%s", redirectURL.Scheme, redirectURL.Host, port, redirectURL.RequestURI())
		copiedConfig := &oauth2.Config{}
		err = copier.Copy(copiedConfig, a.Config)
		if err != nil {
			return err
		}
		copiedConfig.RedirectURL = redirect

		tok, err = copiedConfig.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			return errors.New("Error obtaining token from Spotify: " + err.Error())
		}
	} else {
		tok, err = a.Config.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			return errors.New("Error obtaining token from Spotify: " + err.Error())
		}
	}

	s.Values[token] = *tok
	client := spotify.Authenticator{}.NewClient(tok)
	currentUser, err := client.CurrentUser()
	if err != nil {
		return errors.WithStack(err)
	}
	sdUser, err := a.saveUserData(tok, currentUser)
	if err != nil {
		return errors.WithStack(err)
	}
	// saveUserData doesn't do updates. This call makes sure that
	// the database has an up-to-date AccessToken
	sdUser.AccessToken = tok.AccessToken
	err = a.DB.Users.Update(sdUser)
	if err != nil {
		return errors.WithStack(err)
	}

	s.Values[userID] = sdUser.ID
	err = s.Save(r, w)
	if err != nil {
		return errors.WithStack(err)
	}
	returnPath, ok := s.Values[returnPath].(string)
	if !ok {
		returnPath = "/"
	}
	http.Redirect(w, r, returnPath, http.StatusTemporaryRedirect)
	return nil
}

// GetMyAccessToken will return the Spotify Access token associated with the user's current session
// It probably would've have been better if I embedded this into a JWT and ditched cookies...
// TODO: perhaps it is better design to encapsulate the refresh logic here as well
func (a *AuthController) GetMyAccessToken(w http.ResponseWriter, r *http.Request) error {
	s, err := a.Store.Get(r, sessionName)
	if err != nil {
		return errors.WithStack(err)
	}
	tok, ok := s.Values["Token"].(oauth2.Token)
	if !ok {
		return errors.WithStack(errors.New("Failed to read access_token from session cookie"))
	}
	err = util.JSON(w, tok)
	if err != nil {
		return err
	}
	return nil
}

// Refresh will update the Spotify API Access Token for the user's session
// TODO: check the refresh token and save it (it might be a new refresh token)
func (a *AuthController) Refresh(w http.ResponseWriter, r *http.Request) error {
	s, err := a.Store.Get(r, sessionName)
	if err != nil {
		return errors.WithStack(err)
	}
	ID, ok := s.Values["User_ID"].(uint)
	if !ok {
		return errors.New("unable to obtain user from session data")
	}
	user, err := a.DB.Users.ByID(ID)
	if err != nil {
		return errors.WithStack(err)
	}
	tok, err := refreshToken(a.Config, user.RefreshToken)
	if err != nil {
		return errors.WithStack(err)
	}
	user.AccessToken = tok.AccessToken
	err = a.DB.Users.Update(user)
	if err != nil {
		return errors.WithStack(err)
	}
	sessionToken, ok := s.Values["Token"].(oauth2.Token)
	if !ok {
		return errors.New("Error reading OAuth Token from Session")
	}
	sessionToken.AccessToken = tok.AccessToken
	now := time.Now()
	sessionToken.Expiry = now.Add(time.Duration(tok.ExpiresIn) * time.Second)
	s.Values["Token"] = sessionToken
	err = s.Save(r, w)
	if err != nil {
		return errors.WithStack(err)
	}
	err = util.JSON(w, &tok)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Logout ultimately deletes the user's session
// Since these sessions are stateless, all we have to do is set the max-age to less than 0
func (a *AuthController) Logout(w http.ResponseWriter, r *http.Request) error {
	s, err := a.Store.Get(r, sessionName)
	if err != nil {
		return err
	}
	// per documentation, delete the session by setting Max Age less than 0
	s.Options.MaxAge = -1
	err = s.Save(r, w)
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return nil
}

func refreshToken(c *oauth2.Config, refreshToken string) (*refreshTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	body := strings.NewReader(data.Encode())
	req, err := http.NewRequest(http.MethodPost, spotifyURL+"/api/token", body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.ClientID, c.ClientSecret)
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
		Product:      u.Product,
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
			return errors.WithStack(err)
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
