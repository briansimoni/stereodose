package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/briansimoni/stereodose/app/util"
	"github.com/pkg/errors"
)

const sessionName = "stereodose_session"

// StatusError is handy for when you want to return something other than 500 internal server error
type statusError struct {
	error
	Message string
	Code    int
}

func (e *statusError) Error() string {
	return e.Message
}

// Status returns the http response code
func (e *statusError) Status() int {
	return e.Code
}

// UserContextMiddleware inspects the cookie and adds the user to the context
// For this middleware to work, the user must be authenticated
func UserContextMiddleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) error {
		_, err := r.Cookie(sessionName)
		if err != nil {
			return &statusError{
				Message: fmt.Sprintf("unauthorized request: %s", err.Error()),
				Code:    http.StatusUnauthorized,
			}
		}
		s, err := store.Get(r, sessionName)
		if err != nil {
			return errors.WithStack(err)
		}

		ID, ok := s.Values["User_ID"].(uint)
		if !ok {
			return errors.New("Unable to obtain User_ID from session")
		}
		user, err := stereoDoseDB.Users.ByID(ID)
		if err != nil {
			return &statusError{
				Message: fmt.Sprintf("unauthorized request: %s", err.Error()),
				Code:    http.StatusUnauthorized,
			}
		}
		// We dereference the pointer and store the value in the context
		// instead of storing a pointer to the user
		ctx := context.WithValue(r.Context(), "User", *user)
		next.ServeHTTP(w, r.WithContext(ctx))
		return nil
	}
	return util.HandlerFunc(f)
}
