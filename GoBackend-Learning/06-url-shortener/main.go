package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

var (
	store = make(map[string]string)
	mu    sync.Mutex
)

func generateShortCode() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var req ShortenRequest
	json.NewDecoder(r.Body).Decode(&req)
	code := generateShortCode()
	mu.Lock()
	store[code] = req.URL
	mu.Unlock()
	json.NewEncoder(w).Encode(ShortenResponse{ShortURL: "http://localhost:8086/" + code})
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:]
	mu.Lock()
	url, ok := store[code]
	mu.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)
	fmt.Println("URL Shortener running at http://localhost:8086/")
	log.Fatal(http.ListenAndServe(":8086", nil))
}
