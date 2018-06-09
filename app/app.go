package app

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/briansimoni/stereodose/app/controllers"
	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/config"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hello world")
}

func loggedIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "login success")
}

const sessionName = "_stereodose-session"

var store *sessions.CookieStore
var db *gorm.DB
var stereoDoseDB *models.StereoDoseDB

// InitApp puts together the Router to use as the app's main HTTP handler
func InitApp(c *config.Config, db *gorm.DB) *mux.Router {
	var err error

	authKey, err := base64.StdEncoding.DecodeString(c.AuthKey)
	if err != nil {
		log.Fatal("Unable to obtain auth key", err.Error())
	}
	encryptionKey, err := base64.StdEncoding.DecodeString(c.EncryptionKey)
	if err != nil {
		log.Fatal("Unable to obtain encryption key", err.Error())
	}
	store = sessions.NewCookieStore(authKey, encryptionKey)

	stereoDoseDB = models.NewStereodoseDB(db, store)

	app := mux.NewRouter()
	app.Use(func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, next)
	})

	users := controllers.UsersController{
		DB: stereoDoseDB,
	}
	auth := controllers.AuthController{
		DB:    stereoDoseDB,
		Store: store,
	}

	app.Use(UserContextMiddleware)

	// authRouter := app.PathPrefix("/auth").Subrouter()
	// auth.RegisterHandlers(c, store, authRouter)

	app.HandleFunc("/", index)
	app.HandleFunc("/test", auth.Middleware(webPlayerTest))

	notFound := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Need to add a 404 page")
	}
	app.NotFoundHandler = http.HandlerFunc(notFound)

	app.HandleFunc("/other", auth.Middleware(loggedIn))

	authRouter := app.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", auth.Login).Methods(http.MethodGet)
	authRouter.HandleFunc("/callback", auth.Callback).Methods(http.MethodGet)
	authRouter.HandleFunc("/refresh", auth.Refresh).Methods(http.MethodGet)
	app.HandleFunc("/me", auth.Middleware(users.Me)).Methods(http.MethodGet)
	return app
}

// could do this on a subrouter to handle auth for all routes
// app.Use(func(next http.Handler) http.Handler {
// 	return auth.Middleware(next.ServeHTTP)
// })
