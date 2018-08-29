package app

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/briansimoni/stereodose/app/controllers"
	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/briansimoni/stereodose/config"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
)

const sessionName = "_stereodose-session"

var store *sessions.CookieStore
var db *gorm.DB
var stereoDoseDB *models.StereoDoseDB
var err error

// InitApp puts together the Router to use as the app's main HTTP handler
func InitApp(c *config.Config, db *gorm.DB) *util.AppRouter {
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
	return createRouter(c)
}

func createRouter(c *config.Config) *util.AppRouter {
	app := &util.AppRouter{mux.NewRouter()}
	app.Use(func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, next)
	})

	categories := controllers.NewCategoriesController()
	users := controllers.NewUsersController(stereoDoseDB)
	playlists := controllers.NewPlaylistsController(stereoDoseDB)
	auth := controllers.NewAuthController(stereoDoseDB, store, c)

	// Serve all of the static files
	fs := http.StripPrefix("/public/", http.FileServer(http.Dir("app/views/build/")))
	app.PathPrefix("/public/").Handler(fs)

	app.Handle("/test", auth.Middleware(webPlayerTest))

	notFound := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Need to add a 404 page")
	}
	app.NotFoundHandler = http.HandlerFunc(notFound)

	authRouter := util.AppRouter{app.PathPrefix("/auth").Subrouter()}
	authRouter.AppHandler("/login", auth.Login).Methods(http.MethodGet)
	authRouter.AppHandler("/callback", auth.Callback).Methods(http.MethodGet)
	authRouter.AppHandler("/refresh", auth.Refresh).Methods(http.MethodGet)

	usersRouter := util.AppRouter{app.PathPrefix("/api/users/").Subrouter()}
	usersRouter.Use(UserContextMiddleware)
	usersRouter.Handle("/me", auth.Middleware(users.Me)).Methods(http.MethodGet)

	// The order that the routes are registered does matter
	playlistsRouter := util.AppRouter{app.PathPrefix("/api/playlists").Subrouter()}
	playlistsRouter.Use(UserContextMiddleware)
	playlistsRouter.AppHandler("/", playlists.GetPlaylists).Methods(http.MethodGet)
	playlistsRouter.AppHandler("/", playlists.GetPlaylists).
		Queries(
			"offset", "{offset:[0-9]+}",
			"limit", "{limit:[0-9]+}",
			"category", "{category:[a-zA-Z]+}",
			"subcategory", "{subcategory:[a-zA-Z]+}").
		Methods(http.MethodGet)
	playlistsRouter.AppHandler("/me", playlists.GetMyPlaylists).Methods(http.MethodGet)
	playlistsRouter.AppHandler("/{id}", playlists.GetPlaylistByID).Methods(http.MethodGet)
	playlistsRouter.AppHandler("/", playlists.CreatePlaylist).Methods(http.MethodPost)
	playlistsRouter.AppHandler("/{id}", playlists.DeletePlaylist).Methods(http.MethodDelete)

	categoriesRouter := util.AppRouter{app.PathPrefix("/api/categories").Subrouter()}
	categoriesRouter.AppHandler("/", categories.GetAvailableCategories).Methods(http.MethodGet)

	app.HandleFunc("/", webPlayerTest)

	return app
}

// could do this on a subrouter to handle auth for all routes
// app.Use(func(next http.Handler) http.Handler {
// 	return auth.Middleware(next.ServeHTTP)
// })

// app.Use(func(next http.Handler) http.Handler {
// 	return auth.Middleware(next)
// })
