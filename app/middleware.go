package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// UserContextMiddleware inspects the cookie and adds the user to the context
// For this middleware to work, the user must be authenticated
func UserContextMiddleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		s, err := store.Get(r, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ID, ok := s.Values["User_ID"].(uint)
		if !ok {
			// this doesn't work
			// http.Error(w, "Unable to find User_ID in session", http.StatusInternalServerError)
			// return
			next.ServeHTTP(w, r)
			return
		}
		user, err := stereoDoseDB.Users.ByID(ID)
		fmt.Println(user.Playlists)
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
