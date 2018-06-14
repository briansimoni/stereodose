package controllers

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

type Handler struct {
	H func(w http.ResponseWriter, r *http.Request) error
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.H(w, r)
	switch e := err.(type) {
	case Error:
		log.Println("[Error]", e.Status(), e.Error())
		http.Error(w, e.Error(), e.Status())
	default:
		log.Println("[Error]", e.Error())
		http.Error(w, e.Error(), http.StatusInternalServerError)
	}
}
