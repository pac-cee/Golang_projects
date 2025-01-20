package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type URLEntry struct {
	ID        int
	LongURL   string
	ShortCode string
	CreatedAt time.Time
	Clicks    int
}

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "urls.db")
	if err != nil {
		return err
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		long_url TEXT NOT NULL,
		short_code TEXT UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		clicks INTEGER DEFAULT 0
	);`

	_, err = db.Exec(createTable)
	return err
}

func generateShortCode() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:6], nil
}

func createURL(longURL string) (*URLEntry, error) {
	// Check if URL already exists
	var entry URLEntry
	err := db.QueryRow("SELECT id, long_url, short_code, created_at, clicks FROM urls WHERE long_url = ?",
		longURL).Scan(&entry.ID, &entry.LongURL, &entry.ShortCode, &entry.CreatedAt, &entry.Clicks)
	if err == nil {
		return &entry, nil
	}

	// Generate new short code
	shortCode, err := generateShortCode()
	if err != nil {
		return nil, err
	}

	// Insert new URL
	result, err := db.Exec("INSERT INTO urls (long_url, short_code) VALUES (?, ?)",
		longURL, shortCode)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &URLEntry{
		ID:        int(id),
		LongURL:   longURL,
		ShortCode: shortCode,
		CreatedAt: time.Now(),
		Clicks:    0,
	}, nil
}

func getURL(shortCode string) (*URLEntry, error) {
	var entry URLEntry
	err := db.QueryRow("SELECT id, long_url, short_code, created_at, clicks FROM urls WHERE short_code = ?",
		shortCode).Scan(&entry.ID, &entry.LongURL, &entry.ShortCode, &entry.CreatedAt, &entry.Clicks)
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func incrementClicks(shortCode string) error {
	_, err := db.Exec("UPDATE urls SET clicks = clicks + 1 WHERE short_code = ?", shortCode)
	return err
}

func getAllURLs() ([]URLEntry, error) {
	rows, err := db.Query("SELECT id, long_url, short_code, created_at, clicks FROM urls ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []URLEntry
	for rows.Next() {
		var entry URLEntry
		err := rows.Scan(&entry.ID, &entry.LongURL, &entry.ShortCode, &entry.CreatedAt, &entry.Clicks)
		if err != nil {
			return nil, err
		}
		urls = append(urls, entry)
	}
	return urls, nil
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	urls, err := getAllURLs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, urls)
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	entry, err := createURL(longURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.URL.Path[1:]
	entry, err := getURL(shortCode)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = incrementClicks(shortCode)
	if err != nil {
		log.Printf("Error incrementing clicks: %v", err)
	}

	http.Redirect(w, r, entry.LongURL, http.StatusMovedPermanently)
}

func main() {
	if err := initDB(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/shorten", handleShorten)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
