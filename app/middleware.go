package app

import (
	"context"
	"net/http"
)

// SpotifyIDMiddleware inspects the cookie and adds the spotify ID to the context
// If an error occurs, it just continues to the next handler
// which basically means, functions that require the spotify ID
// also need to be behind auth middleware
func SpotifyIDMiddleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		s, err := store.Get(r, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		spotifyID, ok := s.Values["Spotify_UserID"].(string)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "SpotifyID", spotifyID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(f)
}
