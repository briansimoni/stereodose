package util

import (
	"log"
	"net/http"
)

type Error interface {
	error
	Status() int
}

type StatusError struct {
	Code int
	Err  error
}

func (e StatusError) Error() string {
	return e.Error()
}

// Handler struct allows for functions to return errors and still implement
// the the http.Handler interface
type Handler struct {
	H func(w http.ResponseWriter, r *http.Request) error
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(w, r)
	if err != nil {
		switch e := err.(type) {
		case Error:
			log.Println("[Error]", e.Status(), e.Error())
			http.Error(w, e.Error(), e.Status())
		default:
			log.Println("[Error]", e.Error())
			http.Error(w, "error: "+e.Error(), http.StatusInternalServerError)
		}
	}
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Handler{h}.ServeHTTP(w, r)
}
