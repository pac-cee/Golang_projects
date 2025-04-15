package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// User struct holds user credentials. In real apps, store users in a database.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// users map stores username to hashed password for demo purposes.
var users = map[string]string{} // username:hashedPassword
// jwtKey is the secret key used to sign JWT tokens. Keep this safe in real apps!
var jwtKey = []byte("supersecretkey")

// Claims struct is used for JWT token payload.
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// signup handles POST /signup. It hashes the password and stores the user.
func signup(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	// Hash the password with bcrypt (never store plain passwords!)
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	users[user.Username] = string(hash)
	w.WriteHeader(http.StatusCreated)
}

// login handles POST /login. It checks credentials and returns a JWT token.
func login(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	hash, ok := users[user.Username]
	if !ok || bcrypt.CompareHashAndPassword([]byte(hash), []byte(user.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Create JWT claims, set expiry
	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Return token as JSON
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// protected handles GET /protected. It checks for a valid JWT in the Authorization header.
func protected(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	claims := &Claims{}
	// Parse and validate the JWT token
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "Welcome, %s!", claims.Username)
}

func main() {
	// Register handlers for endpoints
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/protected", protected)
	fmt.Println("Auth service running at http://localhost:8081/")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
