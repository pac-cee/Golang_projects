package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to your first Go web server!")
}

func servicesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Services page")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/services", servicesHandler)
	http.ListenAndServe(":8080", nil)
}
