package app

import (
	"context"
	"log"
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
			return
		}

		log.Println(s.Values["User_ID"])
		ID, ok := s.Values["User_ID"].(uint)
		if !ok {
			log.Println("NOT OK")
			next.ServeHTTP(w, r)
			return
		}
		user, err := stereoDoseDB.Users.Me(ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(user)
		// We dereference the pointer and store the value in the context
		// instead of storing a pointer to the user
		ctx := context.WithValue(r.Context(), "User", *user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(f)
}
