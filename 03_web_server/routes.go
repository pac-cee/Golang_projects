package main

import "net/http"

func registerRoutes() {
    http.HandleFunc("/about", aboutHandler)
    http.HandleFunc("/contact", contactHandler)
}
