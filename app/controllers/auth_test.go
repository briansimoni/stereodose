package controllers

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gorilla/sessions"
// )

// func TestMiddleware(t *testing.T) {
// 	store = sessions.NewCookieStore()
// 	s := httptest.NewServer(nil)
// 	defer s.Close()

// 	f := func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 	}
// 	http.HandleFunc("/", Middleware(f))
// 	http.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 	})

// 	res, err := http.Get(s.URL + "/")
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
// 	if res.StatusCode != 200 {
// 		t.Error("Expected 200, Got:", res.StatusCode)
// 	}

// 	cookie := sessions.NewCookie(sessionName, "wat", &sessions.Options{})
// 	req, err := http.NewRequest(http.MethodGet, s.URL+"/", nil)
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
// 	req.AddCookie(cookie)
// 	res, err = http.DefaultClient.Do(req)
// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}
// 	if res.StatusCode != 500 {
// 		t.Error("Expected 500, Got:", res.StatusCode)
// 	}
// }
