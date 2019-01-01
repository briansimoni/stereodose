package app

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

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
	store         *sessions.CookieStore
	db            *gorm.DB
	stereoDoseDB  *models.StereoDoseDB
	cloudBucket   *blob.Bucket
	file          *os.File
	indexHTML     []byte
	robotsTXT     []byte
	manifest      []byte
	serviceWorker []byte
	err           error
)

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

	cloudBucket, err = setupBucket("aws", "stereodose")
	if err != nil {
		log.Fatal("Unable to setup cloud bucket storage", err.Error())
	}

	stereoDoseDB = models.NewStereodoseDB(db, store)
	return createRouter(c)
}

func createRouter(c *config.Config) *util.AppRouter {
	app := &util.AppRouter{mux.NewRouter()}
	app.Use(func(next http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, next)
	})

	// TODO: just move all the files that have to be served from "/" to their own handler funcs
	// // For a progressive web app to work with a service worker not served from
	// // the "/" directory we have to do this.
	// // See https://www.w3.org/TR/service-workers-1/#extended-http-headers
	// app.Use(func(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		if r.URL.Path == "/public/service-worker.js" {
	// 			w.Header().Add("Service-Worker-Allowed", "/")
	// 		}
	// 		next.ServeHTTP(w, r)
	// 	})
	// })

	categories := controllers.NewCategoriesController()
	users := controllers.NewUsersController(stereoDoseDB)
	playlists := controllers.NewPlaylistsController(stereoDoseDB, cloudBucket)
	auth := controllers.NewAuthController(stereoDoseDB, store, c)

	// Serve all of the static files
	fs := http.StripPrefix("/public/", http.FileServer(http.Dir("app/views/build/")))
	app.PathPrefix("/public/").Handler(fs)

	app.HandleFunc("/robots.txt", serveRobotsTxt)
	app.HandleFunc("/manifest.json", serveManifest)
	app.HandleFunc("/sw.js", serveServiceWorker)

	authRouter := util.AppRouter{app.PathPrefix("/auth").Subrouter()}
	authRouter.AppHandler("/login", auth.Login).Methods(http.MethodGet)
	authRouter.AppHandler("/logout", auth.Logout).Methods(http.MethodGet)
	authRouter.AppHandler("/callback", auth.Callback).Methods(http.MethodGet)
	authRouter.AppHandler("/refresh", auth.Refresh).Methods(http.MethodGet)
	authRouter.AppHandler("/token", auth.GetMyAccessToken).Methods(http.MethodGet)

	usersRouter := util.AppRouter{app.PathPrefix("/api/users/").Subrouter()}
	usersRouter.Use(UserContextMiddleware)
	usersRouter.AppHandler("/me", users.Me).Methods(http.MethodGet)

	// The order that the routes are registered does matter
	// authPlaylistsRouter contains endpoints that require an authenticated user
	authPlaylistsRouter := util.AppRouter{app.PathPrefix("/api/playlists").Subrouter()}
	authPlaylistsRouter.Use(UserContextMiddleware)
	playlistsRouter := util.AppRouter{app.PathPrefix("/api/playlists").Subrouter()}

	playlistsRouter.AppHandler("/", playlists.GetPlaylists).Methods(http.MethodGet)
	playlistsRouter.AppHandler("/", playlists.GetPlaylists).
		Queries(
			"offset", "{offset:[0-9]+}",
			"limit", "{limit:[0-9]+}",
			"category", "{category:[a-zA-Z]+}",
			"subcategory", "{subcategory:[a-zA-Z]+}").
		Methods(http.MethodGet)

	authPlaylistsRouter.AppHandler("/me", playlists.GetMyPlaylists).Methods(http.MethodGet)
	playlistsRouter.AppHandler("/{id}", playlists.GetPlaylistByID).Methods(http.MethodGet)
	authPlaylistsRouter.AppHandler("/", playlists.CreatePlaylist).Methods(http.MethodPost)
	authPlaylistsRouter.AppHandler("/{id}/image", playlists.UploadImage).Methods(http.MethodPost)
	authPlaylistsRouter.AppHandler("/{id}", playlists.DeletePlaylist).Methods(http.MethodDelete)

	categoriesRouter := util.AppRouter{app.PathPrefix("/api/categories").Subrouter()}
	categoriesRouter.AppHandler("/", categories.GetAvailableCategories).Methods(http.MethodGet)

	app.HandleFunc("/", serveReactApp)
	// Serving the React app on 404's enables the use of arbitrary routes with react browser-router
	// Otherwise a requet to /some/arbitrary/path from a different origin would simply 404
	// Could use the hash router for a looser coupling but /#/some/path is ugly
	app.HandleFunc("/{page1}", serveReactApp)
	app.HandleFunc("/{page1}/{page2}", serveReactApp)
	app.NotFoundHandler = http.HandlerFunc(serveReactApp404)

	return app
}

func serveReactApp404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, string(indexHTML))
}

func serveReactApp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, string(indexHTML))
}

func serveRobotsTxt(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, string(robotsTXT))
}

func serveManifest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, string(manifest))
}

func serveServiceWorker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprint(w, string(serviceWorker))
}

// load the contents of index.html and robots.txt into memory only when the app starts up
// instead of on each request
func init() {
	file, err = os.Open("./app/views/build/index.html")
	if err != nil {
		log.Fatalf("Unable to open index.html: %s", err.Error())
	}
	defer file.Close()

	indexHTML, err = ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Unable to read the contets of index.html, %s", err.Error())
	}

	robotsFile, err := os.Open("./app/views/public/robots.txt")
	if err != nil {
		log.Fatalf("Unable to open robots.txt: %s", err.Error())
	}
	defer robotsFile.Close()
	robotsTXT, err = ioutil.ReadAll(robotsFile)
	if err != nil {
		log.Fatalf("Unable to read contents of robots.txt %s", err.Error())
	}

	manifestFile, err := os.Open("./app/views/public/manifest.json")
	if err != nil {
		log.Fatalf("Unable to open manifest.json")
	}
	defer manifestFile.Close()
	manifest, err = ioutil.ReadAll(manifestFile)
	if err != nil {
		log.Fatalf("Unable to read contents of manifest.json %s", err.Error())
	}

	serviceWorkerFlie, err := os.Open("./app/views/build/sw.js")
	if err != nil {
		log.Fatalf("Unable to open manifest.json")
	}
	defer serviceWorkerFlie.Close()
	serviceWorker, err = ioutil.ReadAll(serviceWorkerFlie)
	if err != nil {
		log.Fatalf("Unable to read contents of serviceworker.json %s", err.Error())
	}
}

// could do this on a subrouter to handle auth for all routes
// app.Use(func(next http.Handler) http.Handler {
// 	return auth.Middleware(next.ServeHTTP)
// })

// app.Use(func(next http.Handler) http.Handler {
// 	return auth.Middleware(next)
// })
