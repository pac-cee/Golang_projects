package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("your_secret_key")

// User represents the user model
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims represents the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Store represents our in-memory user store
var users = make(map[string]string)

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	if _, exists := users[user.Username]; exists {
		respondWithJSON(w, http.StatusConflict, Response{
			Status:  "error",
			Message: "Username already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Error while hashing password",
		})
		return
	}

	users[user.Username] = string(hashedPassword)

	respondWithJSON(w, http.StatusCreated, Response{
		Status:  "success",
		Message: "User registered successfully",
	})
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	storedPassword, exists := users[user.Username]
	if !exists {
		respondWithJSON(w, http.StatusUnauthorized, Response{
			Status:  "error",
			Message: "Invalid credentials",
		})
		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
		respondWithJSON(w, http.StatusUnauthorized, Response{
			Status:  "error",
			Message: "Invalid credentials",
		})
		return
	}

	// Create token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Error while generating token",
		})
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status: "success",
		Data: map[string]string{
			"token": tokenString,
		},
	})
}

// AuthMiddleware verifies the JWT token
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			respondWithJSON(w, http.StatusUnauthorized, Response{
				Status:  "error",
				Message: "Missing authorization token",
			})
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			respondWithJSON(w, http.StatusUnauthorized, Response{
				Status:  "error",
				Message: "Invalid or expired token",
			})
			return
		}

		next.ServeHTTP(w, r)
	}
}

// ProtectedHandler is an example of a protected endpoint
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "You have access to protected resource",
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/register", RegisterHandler).Methods("POST")
	router.HandleFunc("/login", LoginHandler).Methods("POST")

	// Protected routes
	router.HandleFunc("/protected", AuthMiddleware(ProtectedHandler)).Methods("GET")

	// Configure and start server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server starting on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
