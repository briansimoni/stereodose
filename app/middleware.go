package app

import (
	"context"
	"net/http"

	"github.com/briansimoni/stereodose/app/util"
)

const sessionName = "_stereodose-session"

// UserContextMiddleware inspects the cookie and adds the user to the context
// For this middleware to work, the user must be authenticated
// TODO: well if the user must be authenticated, might as well return a 403 if they're not
func UserContextMiddleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		s, err := store.Get(r, sessionName)
		if err != nil {
			return err
		}

		ID, ok := s.Values["User_ID"].(uint)
		if !ok {
			next.ServeHTTP(w, r)
			return nil
		}
		user, err := stereoDoseDB.Users.ByID(ID)
		if err != nil {
			return err
		}
		// We dereference the pointer and store the value in the context
		// instead of storing a pointer to the user
		ctx := context.WithValue(r.Context(), "User", *user)
		next.ServeHTTP(w, r.WithContext(ctx))
		return nil
	}
	return util.HandlerFunc(f)
}
