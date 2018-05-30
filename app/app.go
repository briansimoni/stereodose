package app

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/briansimoni/stereodose/app/auth"
	"github.com/briansimoni/stereodose/app/models"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

func loggedIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "login success")
}

const sessionName = "_stereodose-session"

var store *sessions.CookieStore

func InitApp(c *auth.Config) *mux.Router {
	authKey, err := base64.StdEncoding.DecodeString(c.AuthKey)
	if err != nil {
		log.Fatal("Unable to obtain auth key", err.Error())
	}
	// encryptionKey, err := base64.StdEncoding.DecodeString(os.Getenv("STEREODOSE_ENCRYPTION_KEY"))
	// if err != nil {
	// 	log.Fatal("Unable to obtain encryption key", err.Error())
	// }
	store = sessions.NewCookieStore(authKey)

	app := mux.NewRouter()
	app.Use(func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, next)
	})

	authRouter := app.PathPrefix("/auth").Subrouter()
	auth.RegisterHandlers(c, store, authRouter)

	app.HandleFunc("/", index)
	app.HandleFunc("/test", auth.Middleware(webPlayerTest))

	notFound := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Need to add a 404 page")
	}
	app.NotFoundHandler = http.HandlerFunc(notFound)

	app.HandleFunc("/other", auth.Middleware(loggedIn))

	app.HandleFunc("/gorm", func(w http.ResponseWriter, r *http.Request) {
		users.HelloWorld()
		fmt.Fprint(w, "hi")
	})
	return app
}
