package util

import (
	"net/http"

	"github.com/gorilla/mux"
)

type AppRouter struct {
	*mux.Router
}

func (a *AppRouter) AppHandler(path string, f func(w http.ResponseWriter, r *http.Request) error) *mux.Route {
	return a.Handle(path, Handler{f})
}
