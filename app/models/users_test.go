package models

import (
	"net/http"
	"testing"
)

func TestMe(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/me", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	req.AddCookie(&http.Cookie{})
}
