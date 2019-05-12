package util

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

type stackTracerError interface {
	StackTrace() errors.StackTrace
	error
}

type statusError interface {
	Status() int
	error
}

// Handler struct allows for functions to return errors and still implement
// the the http.Handler interface
type Handler struct {
	H func(w http.ResponseWriter, r *http.Request) error
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(w, r)
	if err != nil {
		transactionID := r.Context().Value(TransactionIDKey)
		switch e := err.(type) {
		case statusError:
			log.WithFields(logrus.Fields{
				"Type":          "AppLog",
				"TransactionID": transactionID,
				"ErrorType":     fmt.Sprintf("%T", e),
			}).Error(e.Error())
			http.Error(w, "error: "+err.Error(), e.Status())
		case stackTracerError:
			st := e.StackTrace()
			prettyStackTrace := strings.Split(strings.Replace(fmt.Sprintf("%+v", st), "\t", "    ", -1), "\n")
			log.WithFields(logrus.Fields{
				"Type":          "AppLog",
				"TransactionID": transactionID,
				"ErrorType":     fmt.Sprintf("%T", e),
				"StackTrace":    prettyStackTrace,
			}).Error(e.Error())
			http.Error(w, "error: "+err.Error(), http.StatusInternalServerError)
		default:
			log.WithFields(logrus.Fields{
				"Type":          "AppLog",
				"TransactionID": transactionID,
				"ErrorType":     fmt.Sprintf("%T", e),
			}).Error(e.Error())
			http.Error(w, "error: "+e.Error(), http.StatusInternalServerError)
		}
	}
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Handler{h}.ServeHTTP(w, r)
}
