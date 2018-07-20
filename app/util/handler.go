package util

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
)

type stackTracerError interface {
	StackTrace() errors.StackTrace
}

type statusError interface {
	Status() int
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
		case statusError:
			log.Printf("statusStackTracer %T\n", e)
			log.Printf("[ERROR] %d %s\n", e.Status(), err.Error())
			http.Error(w, "error: "+err.Error(), e.Status())
		case stackTracerError:
			log.Printf("%T\n", e)
			st := e.StackTrace()
			log.Printf("[ERROR] %s\n%+v", err.Error(), st[0])
			http.Error(w, "error: "+err.Error(), http.StatusInternalServerError)
		default:
			log.Printf("%T\n", e)
			log.Println("[ERROR]", e.Error())
			http.Error(w, "error: "+e.Error(), http.StatusInternalServerError)
		}
	}
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Handler{h}.ServeHTTP(w, r)
}
