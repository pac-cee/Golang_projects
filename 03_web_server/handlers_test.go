package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestAboutHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/about", nil)
    w := httptest.NewRecorder()
    aboutHandler(w, req)
    if w.Body.String() != "About page. This is a Go web server." {
        t.Errorf("Unexpected response: %s", w.Body.String())
    }
}
