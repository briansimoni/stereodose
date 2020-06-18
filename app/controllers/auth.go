package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/briansimoni/stereodose/config"

	log "github.com/sirupsen/logrus"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	endpoint "golang.org/x/oauth2/spotify"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

const spotifyURL = "https://accounts.spotify.com"

var sessionKeys = struct {
	SessionCookieName   string
	AuthStateCookieName string
	State               string
	Token               string
	UserID              string
	ReturnPath          string
}{
	SessionCookieName:   "stereodose_session",
	AuthStateCookieName: "stereodose_auth_state",
	State:               "State",
	Token:               "Token",
	UserID:              "User_ID",
	ReturnPath:          "ReturnPath",
}

// AuthController is a collection of RESTful Handlers for authentication
type AuthController struct {
	DB               *models.StereoDoseDB
	Store            *sessions.CookieStore
	Config           *oauth2.Config
	StereodoseConfig *config.Config
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
			// "user-read-birthdate",
			"user-read-email",
			"user-read-private",
			"playlist-read-collaborative",
			"playlist-read-private",
			"user-modify-playback-state",
		},
		Endpoint: endpoint.Endpoint,
	}
	a := &AuthController{
		DB:     db,
		Store:  store,
		Config: oauthConfig,
	}
	a.StereodoseConfig = config
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

	// stereodose_auth_state is a separate cookie from the actual session cookie
	// this allows for the presence of a stereodose_session cookie to be proof of authentication
	authState, err := a.Store.Get(r, sessionKeys.AuthStateCookieName)
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
	authState.Values[sessionKeys.ReturnPath] = r.URL.Query().Get("path")
	authState.Values[sessionKeys.State] = state
	authState.Save(r, w)

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
		http.Redirect(w, r, redir, http.StatusFound)
		return nil
	}

	redir := a.Config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	http.Redirect(w, r, redir, http.StatusFound)
	return nil
}

// Callback is a handler function that the user is redirected to in the OAuth flow
// In this step of authorization, we exchange a code for an access token
// and we query the user's profile on Spotify to get their identity
func (a *AuthController) Callback(w http.ResponseWriter, r *http.Request) error {
	log.Printf("%+v\n", r.URL.Query())
	thing := r.URL.Query()
	transactionID := r.Context().Value("TransactionID")
	log.WithFields(log.Fields{
		"TransactionID": transactionID,
		"Params":        thing,
	}).Info("query parameters from Callback URL")
	authState, err := a.Store.Get(r, sessionKeys.AuthStateCookieName)
	if err != nil {
		return errors.WithStack(err)
	}
	err = checkState(r, authState)
	if err != nil {
		return err
	}

	// per documentation, delete the session by setting Max Age less than 0
	authState.Options.MaxAge = -1
	err = authState.Save(r, w)
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

	s, err := a.Store.Get(r, sessionKeys.SessionCookieName)
	if err != nil {
		return errors.WithStack(err)
	}

	s.Values[sessionKeys.Token] = *tok
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

	s.Values[sessionKeys.UserID] = sdUser.ID
	err = s.Save(r, w)
	if err != nil {
		return errors.WithStack(err)
	}
	returnPath := authState.Values[sessionKeys.ReturnPath].(string)
	if returnPath == "" {
		returnPath = "/"
	}

	http.Redirect(w, r, returnPath, http.StatusFound)
	return nil
}

// TokenSwap was created to support the iOS app.
// The Spotify iOS documentation refers to a "token swap" API endpoint which is essentially
// the same as the OAuth callback or redirect URL.
// The difference here is that instead of 302 redirecting on the callback,
// we simply return a 200 response with the JSON returned from the Spotify code exchange
func (a *AuthController) TokenSwap(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return errors.New("unable to parse form data")
	}
	code := r.Form.Get("code")
	if code == "" {
		return &util.StatusError{
			Code:    http.StatusBadRequest,
			Message: "missing 'code' in form data",
		}
	}
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", a.StereodoseConfig.IOSRedirectURL)
	data.Set("code", code)
	body := strings.NewReader(data.Encode())
	request, err := http.NewRequest(http.MethodPost, a.Config.Endpoint.TokenURL, body)
	if err != nil {
		return errors.WithStack(err)
	}
	request.SetBasicAuth(a.Config.ClientID, a.Config.ClientSecret)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.WithStack(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		errorBody, _ := ioutil.ReadAll(response.Body)
		return &util.StatusError{
			Code:    response.StatusCode,
			Message: "Error exchanging token with Spotify: " + string(errorBody),
		}
	}
	type TokenSet struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}
	tokenSet := new(TokenSet)
	err = json.NewDecoder(response.Body).Decode(tokenSet)
	if err != nil {
		return errors.WithStack(err)
	}

	s, err := a.Store.Get(r, sessionKeys.SessionCookieName)
	if err != nil {
		return errors.WithStack(err)
	}

	oauth2Token := &oauth2.Token{
		AccessToken:  tokenSet.AccessToken,
		RefreshToken: tokenSet.RefreshToken,
		// Expiry: time.Now().Add(time.Duration(tokenSet.ExpiresIn) * time.Second),
	}
	log.Info("The token value: " + oauth2Token.AccessToken)

	s.Values[sessionKeys.Token] = *oauth2Token
	client := spotify.Authenticator{}.NewClient(oauth2Token)
	currentUser, err := client.CurrentUser()
	if err != nil {
		return errors.WithStack(err)
	}
	sdUser, err := a.saveUserData(oauth2Token, currentUser)
	if err != nil {
		return errors.WithStack(err)
	}
	// saveUserData doesn't do updates. This call makes sure that
	// the database has an up-to-date AccessToken
	sdUser.AccessToken = oauth2Token.AccessToken
	err = a.DB.Users.Update(sdUser)
	if err != nil {
		return errors.WithStack(err)
	}

	util.JSON(w, tokenSet)
	return nil
}

// MobileLogin is here to support the iOS app.
// Because Spotify basically constrains iOS developers to using only their SDK
// and not the WebAPI for authentication, we have to create this seperate endpoint
// It will take a Spotify access token, and then create a Stereodose session.
// The response is simply a 200 and a set-cookie header
func (a *AuthController) MobileLogin(w http.ResponseWriter, r *http.Request) error {
	type TokenBody struct {
		AccessToken string `json:"access_token"`
	}
	tokenBody := new(TokenBody)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(tokenBody)
	if err != nil {
		return errors.WithStack(err)
	}

	oauth2Token := &oauth2.Token{
		AccessToken: tokenBody.AccessToken,
	}

	s, err := a.Store.Get(r, sessionKeys.SessionCookieName)
	if err != nil {
		return errors.WithStack(err)
	}

	client := spotify.Authenticator{}.NewClient(oauth2Token)
	currentUser, err := client.CurrentUser()
	if err != nil {
		return errors.WithStack(err)
	}
	sdUser, err := a.DB.Users.BySpotifyID(currentUser.ID)
	if err != nil {
		return errors.WithStack(err)
	}

	oauth2Token.RefreshToken = sdUser.RefreshToken

	s.Values[sessionKeys.Token] = *oauth2Token
	s.Values[sessionKeys.UserID] = sdUser.ID
	err = s.Save(r, w)
	if err != nil {
		return errors.WithStack(err)
	}
	cookie := w.Header().Get("Set-Cookie")
	type CookieResponse struct {
		Cookie string `json:"cookie"`
	}
	cookieResponse := &CookieResponse{
		Cookie: cookie,
	}
	err = util.JSON(w, cookieResponse)
	if err != nil {
		return err
	}
	return nil
}

// GetMyAccessToken will return the Spotify Access token associated with the user's current session
// It probably would've have been better if I embedded this into a JWT and ditched cookies...
// TODO: perhaps it is better design to encapsulate the refresh logic here as well
func (a *AuthController) GetMyAccessToken(w http.ResponseWriter, r *http.Request) error {
	s, err := a.Store.Get(r, sessionKeys.SessionCookieName)
	if err != nil {
		return errors.WithStack(err)
	}
	tok, ok := s.Values["Token"].(oauth2.Token)
	if !ok {
		return errors.WithStack(errors.New("Failed to read access_token from session cookie"))
	}
	err = util.JSON(w, tok)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Refresh will update the Spotify API Access Token for the user's session
// TODO: check the refresh token and save it (it might be a new refresh token)
func (a *AuthController) Refresh(w http.ResponseWriter, r *http.Request) error {
	s, err := a.Store.Get(r, sessionKeys.SessionCookieName)
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
	s, err := a.Store.Get(r, sessionKeys.SessionCookieName)
	if err != nil {
		return err
	}
	// per documentation, delete the session by setting Max Age less than 0
	s.Options.MaxAge = -1
	err = s.Save(r, w)
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/", http.StatusFound)
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
		str, _ := ioutil.ReadAll(res.Body)
		return nil, errors.New("Response from spotify.com/api/token " + res.Status + " " + string(str))
	}
	defer res.Body.Close()
	var tok refreshTokenResponse
	err = json.NewDecoder(res.Body).Decode(&tok)
	if err != nil {
		return nil, err
	}
	return &tok, nil
}

func Test() string {
	return "lol"
}

func checkState(r *http.Request, s *sessions.Session) error {
	responseState := r.URL.Query().Get("state")
	if responseState == "" {
		return errors.New("Unable to obtain state from URL query params")
	}

	sessionState := s.Values[sessionKeys.State]

	if responseState != sessionState {
		return errors.New(fmt.Sprintf("State mismatch. responseState: %s sessionState: %s", responseState, sessionState))
	}
	return nil
}

func (a *AuthController) saveUserData(token *oauth2.Token, u *spotify.PrivateUser) (*models.User, error) {
	log.Println(token.RefreshToken)
	user := &models.User{
		Birthdate:    u.Birthdate,
		DisplayName:  u.DisplayName,
		Email:        u.Email,
		SpotifyID:    u.ID,
		Product:      u.Product,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}

	user, err := a.DB.Users.FirstOrCreate(user, token)
	if err != nil {
		return nil, err
	}

	// we need to load up the user images using the ByID method before we attempt updates
	user, err = a.DB.Users.ByID(user.ID)
	if err != nil {
		return nil, err
	}

	// now we add the user's profile images from spotify
	// first we make sure that it isn't already saved in the database
	for _, spotifyImage := range u.Images {
		newImage := true
		for _, userImage := range user.Images {
			if spotifyImage.URL == userImage.URL {
				newImage = false
				break
			}
		}
		if newImage {
			var image models.UserImage
			image.Height = spotifyImage.Height
			image.Width = spotifyImage.Width
			image.URL = spotifyImage.URL
			user.Images = append(user.Images, image)
		}

	}
	// make sure that the tokens are up to date
	user.RefreshToken = token.RefreshToken
	user.AccessToken = token.AccessToken

	err = a.DB.Users.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
