package util

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

const sessionName = "_stereodose-session"

// Session holds all of the data stored in the user's session cookie
type Session struct {
	AccessToken     string
	TokenExpiration string
	RefreshToken    string
	SpotifyUserID   string
}

// GetSessionInfo returns a convenient Session struct
func GetSessionInfo(store *sessions.CookieStore, r *http.Request) (*Session, error) {
	s, err := store.Get(r, sessionName)
	if err != nil {
		return nil, err
	}
	session := &Session{}
	accessToken, ok := s.Values["Access_Token"].(string)
	if !ok {
		return nil, errors.New("Unable to obtain access token from session")
	}
	session.AccessToken = accessToken

	expiration, ok := s.Values["Expiry"].(string)
	if !ok {
		return nil, errors.New("Unable to obtain access token expiration from session")
	}
	session.TokenExpiration = expiration

	refreshToken, ok := s.Values["Refresh_Token"].(string)
	if !ok {
		return nil, errors.New("Unable to obtain refresh token from session")
	}
	session.RefreshToken = refreshToken

	spotifyID, ok := s.Values["Spotify_UserID"].(string)
	if !ok {
		return nil, errors.New("Unable to obtain Spotify ID from session")
	}
	session.SpotifyUserID = spotifyID

	return session, nil
}
