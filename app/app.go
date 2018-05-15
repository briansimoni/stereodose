package app

import (
	"fmt"
	"net/http"

	"github.com/briansimoni/stereodose/app/auth"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello")
}

func loggedIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "login success")
}

func InitApp() *mux.Router {
	var store = sessions.NewCookieStore([]byte("something-very-secret"))

	app := mux.NewRouter()

	authRouter := app.PathPrefix("/auth").Subrouter()
	auth.RegisterHandlers(store, authRouter)

	app.HandleFunc("/", index)

	notFound := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Need to add a 404 page")
	}
	app.NotFoundHandler = http.HandlerFunc(notFound)

	app.HandleFunc("/other", auth.Middleware(loggedIn))
	return app
}
