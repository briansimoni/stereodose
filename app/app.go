package app

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/briansimoni/stereodose/app/controllers"
	"github.com/briansimoni/stereodose/app/models"
	"github.com/briansimoni/stereodose/app/util"
	"github.com/briansimoni/stereodose/config"
	"github.com/google/go-cloud/blob"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
)

var (
	store        *sessions.CookieStore
	db           *gorm.DB
	stereoDoseDB *models.StereoDoseDB
	cloudBucket  *blob.Bucket
	err          error
	fileCache    map[string][]byte
)

// InitApp puts together the Router to use as the app's main HTTP handler
func InitApp(c *config.Config) (*util.AppRouter, *models.StereoDoseDB) {
	authKey, err := base64.StdEncoding.DecodeString(c.AuthKey)
	if err != nil {
		log.Fatal("Unable to obtain auth key", err.Error())
	}
	encryptionKey, err := base64.StdEncoding.DecodeString(c.EncryptionKey)
	if err != nil {
		log.Fatal("Unable to obtain encryption key", err.Error())
	}
	store = sessions.NewCookieStore(authKey, encryptionKey)

	cloudBucket, err = setupBucket("aws", "stereodose")
	if err != nil {
		log.Fatal("Unable to setup cloud bucket storage", err.Error())
	}

	stereoDoseDB = models.NewStereodoseDB(c, store)
	return createRouter(c), stereoDoseDB
}

func createRouter(c *config.Config) *util.AppRouter {
	app := &util.AppRouter{Router: mux.NewRouter()}

	app.Use(handlers.ProxyHeaders)
	app.Use(util.RequestLogger)

	categories := controllers.NewCategoriesController()
	users := controllers.NewUsersController(stereoDoseDB)
	playlists := controllers.NewPlaylistsController(stereoDoseDB, cloudBucket)
	auth := controllers.NewAuthController(stereoDoseDB, store, c)
	health := controllers.NewHealthController(stereoDoseDB)
	feedback := controllers.NewFeedbackController(stereoDoseDB, store)

	// Serve all of the static files
	fs := http.StripPrefix("/public/", http.FileServer(http.Dir("app/views/build/")))
	app.PathPrefix("/public/").Handler(fs)

	app.HandleFunc("/robots.txt", serveFile(fileCache["/robots.txt"], nil))
	app.HandleFunc("/manifest.json", serveFile(fileCache["/manifest.json"], nil))
	app.HandleFunc("/sw.js", serveFile(fileCache["/sw.js"], map[string]string{"Content-Type": "application/javascript"}))
	app.HandleFunc("/terms-and-conditions", serveFile(fileCache["/terms-and-conditions.html"], nil))
	app.HandleFunc("/privacy-policy", serveFile(fileCache["/privacy-policy.html"], nil))
	app.HandleFunc("/ads.txt", serveFile(fileCache["/ads.txt"], nil))

	healthRouter := util.AppRouter{Router: app.PathPrefix("/api/health").Subrouter()}
	healthRouter.AppHandler("/", health.CheckHealth).Methods(http.MethodGet)

	authRouter := util.AppRouter{Router: app.PathPrefix("/auth").Subrouter()}
	authRouter.AppHandler("/login", auth.Login).Methods(http.MethodGet)
	authRouter.AppHandler("/logout", auth.Logout).Methods(http.MethodGet)
	authRouter.AppHandler("/callback", auth.Callback).Methods(http.MethodGet)
	authRouter.AppHandler("/refresh", auth.Refresh).Methods(http.MethodGet)
	authRouter.AppHandler("/token", auth.GetMyAccessToken).Methods(http.MethodGet)
	authRouter.AppHandler("/token-swap", auth.TokenSwap).Methods(http.MethodPost)
	authRouter.AppHandler("/mobile-login", auth.MobileLogin).Methods(http.MethodPost)

	protectedUserRouter := util.AppRouter{Router: app.PathPrefix("/api/users/").Subrouter()}
	protectedUserRouter.Use(UserContextMiddleware)
	protectedUserRouter.AppHandler("/me", users.Me).Methods(http.MethodGet)

	userRouter := util.AppRouter{Router: app.PathPrefix("/api/users/").Subrouter()}
	userRouter.AppHandler("/{id}", users.GetByID).Methods(http.MethodGet)
	userRouter.AppHandler("/{id}/likes", users.GetUserLikes).Methods(http.MethodGet)

	// The order that the routes are registered does matter
	// protectedPlaylistsRouter contains endpoints that require an authenticated user
	protectedPlaylistsRouter := util.AppRouter{Router: app.PathPrefix("/api/playlists").Subrouter()}
	protectedPlaylistsRouter.Use(UserContextMiddleware)
	playlistsRouter := util.AppRouter{Router: app.PathPrefix("/api/playlists").Subrouter()}

	playlistsRouter.AppHandler("/", playlists.GetPlaylists).Methods(http.MethodGet)
	playlistsRouter.AppHandler("/", playlists.GetPlaylists).
		Queries(
			"offset", "{offset:[0-9]}",
			"limit", "{limit:[0-9]{2}}",
			"category", "{category:[a-zA-Z]+}",
			"subcategory", "{subcategory:[a-zA-Z]*}",
			"sort-key", "{sort-key:[a-zA-Z_]",
			"order", "{order:[a-zA-Z]").
		Methods(http.MethodGet)

	playlistsRouter.AppHandler("/random", playlists.GetRandomPlaylist).
		Queries(
			"category", "{category:.+}",
			"subcategory", "{subcategory:.+}",
		).Methods(http.MethodGet)

	protectedPlaylistsRouter.AppHandler("/me", playlists.GetMyPlaylists).Methods(http.MethodGet)
	playlistsRouter.AppHandler("/{id}", playlists.GetPlaylistByID).Methods(http.MethodGet)
	protectedPlaylistsRouter.AppHandler("/", playlists.CreatePlaylist).Methods(http.MethodPost)
	protectedPlaylistsRouter.AppHandler("/image", playlists.UploadImage).Methods(http.MethodPost)
	protectedPlaylistsRouter.AppHandler("/{id}", playlists.DeletePlaylist).Methods(http.MethodDelete)
	protectedPlaylistsRouter.AppHandler("/{id}/comments", playlists.Comment).Methods(http.MethodPost)
	protectedPlaylistsRouter.AppHandler("/{playlistID}/comments/{commentID}", playlists.DeleteComment).Methods(http.MethodDelete)
	protectedPlaylistsRouter.AppHandler("/{id}/likes", playlists.Like).Methods(http.MethodPost)
	protectedPlaylistsRouter.AppHandler("/{playlistID}/likes/{likeID}", playlists.Unlike).Methods(http.MethodDelete)

	categoriesRouter := util.AppRouter{Router: app.PathPrefix("/api/categories").Subrouter()}
	categoriesRouter.AppHandler("/", categories.GetAvailableCategories).Methods(http.MethodGet)

	feedbackRouter := util.AppRouter{Router: app.PathPrefix("/api/feedback").Subrouter()}
	feedbackRouter.AppHandler("/", feedback.CreateFeedback).Methods(http.MethodPost)

	app.HandleFunc("/", serveFile(fileCache["/index.html"], nil))
	// Serving the React app on 404's enables the use of arbitrary routes with react browser-router
	// Otherwise a request to /some/arbitrary/path from a different origin would simply 404
	// Could use the hash router for a looser coupling but /#/some/path is ugly
	app.HandleFunc("/{page1}", serveFile(fileCache["/index.html"], nil))
	app.HandleFunc("/{page1}/{page2}", serveFile(fileCache["/index.html"], nil))
	app.HandleFunc("/{page1}/{page2}/type", serveFile(fileCache["/index.html"], nil))
	app.HandleFunc("/{page1}/{page2}/{playlistID:[A-Za-z0-9]{22}}", dynamicStatusCodeHandler)
	// app.NotFoundHandler = util.RequestLogger(app.NotFoundHandler)
	app.NotFoundHandler = util.RequestLogger(http.HandlerFunc(serveReactApp404))

	return app
}

func dynamicStatusCodeHandler(w http.ResponseWriter, r *http.Request) {
	pathVars := mux.Vars(r)
	playlistID := pathVars["playlistID"]
	playlist, _ := stereoDoseDB.Playlists.GetByID(playlistID)

	if playlist == nil {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, string(fileCache["/index.html"]))
}

func serveReactApp404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, string(fileCache["/index.html"]))
}

// serve file takes file data and optionally headers and returns an http.Handler function
func serveFile(data []byte, headers map[string]string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for key, value := range headers {
			w.Header().Set(key, value)
		}
		fmt.Fprint(w, string(data))
	}
}

// loadFile adds hard-coded files to a cache which can be used later
// the fileCache map uses /filename.extension as the key
// for example, the key for ./app/views/build/index.html is simply /index.html
func loadFile(filePath string) error {
	split := strings.Split(filePath, "/")
	name := "/" + split[len(split)-1]
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fileCache[name] = data
	return nil
}

// load the contents of certain static files into memory only when the app starts up
// instead of on each request
func init() {
	fileCache = make(map[string][]byte, 0)
	files := []string{
		"./app/views/build/index.html",
		"./app/views/public/robots.txt",
		"./app/views/public/manifest.json",
		"./app/views/build/sw.js",
		"./app/views/public/terms-and-conditions.html",
		"./app/views/public/privacy-policy.html",
		"./app/views/public/ads.txt",
	}

	for _, file := range files {
		err := loadFile(file)
		if err != nil {
			log.Fatalf("Unable to load file: %s. %s", file, err.Error())
		}
	}
}
