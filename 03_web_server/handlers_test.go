package main

import (
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

func TestServicesHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/services", nil)
	w := httptest.NewRecorder()
	servicesHandler(w, req)
	if w.Body.String() != "Our services include web development, cloud hosting, and API integration." {
		t.Errorf("Unexpected response: %s", w.Body.String())
	}
}
