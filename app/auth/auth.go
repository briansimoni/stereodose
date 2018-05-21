package auth

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"

	"github.com/gorilla/sessions"
)

const sessionName = "_stereodose-session"

var store *sessions.CookieStore

// RegisterHandlers adds the routes and handlers to a router
// that are needed for authentication purposes
func RegisterHandlers(cookieStore *sessions.CookieStore, r *mux.Router) {
	store = cookieStore
	r.HandleFunc("/login", login)
	r.HandleFunc("/callback", callback)
}

var conf = &oauth2.Config{
	ClientID:     os.Getenv("STEREODOSE_CLIENT_ID"),
	ClientSecret: os.Getenv("STEREODOSE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("STEREODOSE_REDIRECT_URL"),
	Scopes:       []string{"playlist-modify-public"},
	Endpoint:     spotify.Endpoint,
}

func login(w http.ResponseWriter, r *http.Request) {
	s, err := store.Get(r, sessionName)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if s.Values["Access_Token"] == nil {
		// user is not logged in. send to authorization code flow
		// Redirect user to consent page to ask for permission
		// for the specified scopes.

		b := make([]byte, 32)
		_, err = rand.Read(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		state := base64.StdEncoding.EncodeToString(b)
		s.Values["State"] = state
		s.Save(r, w)

		url := conf.AuthCodeURL(state, oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

}

func callback(w http.ResponseWriter, r *http.Request) {

	s, err := store.Get(r, sessionName)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseState := r.URL.Query().Get("state")
	if responseState == "" {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := s.Values["State"]
	if state == "" {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != state {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tok, err := conf.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.Values["Access_Token"] = tok.AccessToken
	s.Values["Expiry"] = tok.Expiry.String()
	s.Values["Refresh_Token"] = tok.RefreshToken
	err = s.Save(r, w)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	returnPath, ok := s.Values["return_path"].(string)
	if !ok {
		log.Println("not okay")
		returnPath = "/"
	}
	http.Redirect(w, r, returnPath, http.StatusTemporaryRedirect)

}

// Middleware checks to see if the user is logged in before
// allowing the request to continue
func Middleware(next http.HandlerFunc) http.HandlerFunc {
	f := func(w http.ResponseWriter, r *http.Request) {
		s, err := store.Get(r, sessionName)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if s.Values["Access_Token"] == nil {
			s.Values["return_path"] = r.URL.Path
			s.Save(r, w)
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}
