package main

import (
    "fmt"
    "net/http"
)

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "About page. This is a Go web server.")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Contact us at go@example.com")
}
