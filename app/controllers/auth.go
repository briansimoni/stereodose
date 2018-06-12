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

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"

	"github.com/gorilla/sessions"
)

const sessionName = "_stereodose-session"

var spotifyURL = "https://accounts.spotify.com"

type AuthController struct {
	DB    *models.StereoDoseDB
	Store *sessions.CookieStore
}

// spotifyUser struct is used when querying the /me API endpoint
type spotifyUser struct {
	Birthdate    string      `json:"birthdate"`
	Country      string      `json:"country"`
	DisplayName  interface{} `json:"display_name"`
	Email        string      `json:"email"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Followers struct {
		Href  interface{} `json:"href"`
		Total int         `json:"total"`
	} `json:"followers"`
	Href    string        `json:"href"`
	ID      string        `json:"id"`
	Images  []interface{} `json:"images"`
	Product string        `json:"product"`
	Type    string        `json:"type"`
	URI     string        `json:"uri"`
}

type refreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// RegisterHandlers adds the routes and handlers to a router
// that are needed for authentication purposes
// func RegisterHandlers(c *config.Config, cookieStore *sessions.CookieStore, r *mux.Router) {
// 	store = cookieStore
// 	r.HandleFunc("/login", login).Methods(http.MethodGet)
// 	r.HandleFunc("/callback", callback).Methods(http.MethodGet)
// }

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
	},
	Endpoint: spotify.Endpoint,
}

func (a *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	s, err := a.Store.Get(r, sessionName)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if s.Values["Access_Token"] == nil {
		// user is not logged in. send to authorization code flow
		// Redirect user to consent page to ask for permission
		// for the specified scopes.

		b := make([]byte, 32)
		_, err = rand.Read(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		state := base64.StdEncoding.EncodeToString(b)
		s.Values["State"] = state
		s.Save(r, w)

		url := conf.AuthCodeURL(state, oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

func (a *AuthController) Callback(w http.ResponseWriter, r *http.Request) {
	s, err := a.Store.Get(r, sessionName)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = checkState(r, s)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tok, err := conf.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Error obtaining getting token from Spotify"+err.Error(), http.StatusInternalServerError)
		return
	}
	s.Values["Access_Token"] = tok.AccessToken
	s.Values["Expiry"] = tok.Expiry.Format(time.RFC822)
	user, err := GetUserData(tok.AccessToken)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sdUser, err := a.saveUserData(tok, user)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.Values["User_ID"] = sdUser.ID
	err = s.Save(r, w)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	returnPath, ok := s.Values["return_path"].(string)
	if !ok {
		returnPath = "/"
	}
	http.Redirect(w, r, returnPath, http.StatusTemporaryRedirect)
}

// Refresh will update the Spotify API Access Token for the user's session
// TODO: check the refresh token and save it (it might be a new refresh token)
func (a *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("User").(models.User)
	if !ok {
		http.Error(w, "Unable to obtain user from context", http.StatusInternalServerError)
		return
	}
	s, err := a.Store.Get(r, sessionName)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tok, err := refreshToken(user.RefreshToken)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.AccessToken = tok.AccessToken
	err = a.DB.Users.Update(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.Values["Access_Token"] = tok.AccessToken
	// TODO: fix this
	s.Values["Expiry"] = tokenExpirationDate(tok.ExpiresIn)
	err = s.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	j, err := json.MarshalIndent(&tok, " ", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

// TODO: figure out if display name is deterministic or something
func (a *AuthController) saveUserData(token *oauth2.Token, u *spotifyUser) (*models.User, error) {
	var displayName string
	displayName, ok := u.DisplayName.(string)
	if !ok {
		displayName = u.ID
	}

	user := &models.User{
		Birthdate:    u.Birthdate,
		DisplayName:  displayName,
		Email:        u.Email,
		SpotifyID:    u.ID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	user, err := a.DB.Users.FirstOrCreate(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserData(accessToken string) (*spotifyUser, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	// TODO: do not use default client
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var u spotifyUser
	err = json.NewDecoder(res.Body).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Middleware checks to see if the user is logged in before
// allowing the request to continue
func (a *AuthController) Middleware(next http.HandlerFunc) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		s, err := a.Store.Get(r, sessionName)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		expireTime, ok := s.Values["Expiry"].(string)
		if !ok {
			s.Values["return_path"] = r.URL.Path
			s.Save(r, w)
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return
		}
		expired, err := isExpired(expireTime)
		if err != nil {
			log.Println("[ERROR]", err.Error())
			s.Values["return_path"] = r.URL.Path
			s.Save(r, w)
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return
		}
		if expired {
			s.Values["return_path"] = r.URL.Path
			s.Save(r, w)
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return
		}

		log.Println("[INFO] EXPIRES:", s.Values["Expiry"])
		log.Printf("%T", s.Values["Expiry"])
		if s.Values["Access_Token"] == nil {
			s.Values["return_path"] = r.URL.Path
			s.Save(r, w)
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

// takes the expires time from Spotify API response and converts it to a date
func tokenExpirationDate(expires int) string {
	t := time.Now()
	expiresTime := t.Add(time.Second * time.Duration(expires))
	// convert this to rfc822
	return expiresTime.Format(time.RFC822)
}

// TODO: use better time format
// need to write unit tests for these. This function is working right now
func isExpired(date string) (bool, error) {
	t, err := time.Parse(time.RFC822, date)
	if err != nil {
		return true, err
	}
	if time.Now().After(t) {
		return true, nil
	}
	return false, nil
}
