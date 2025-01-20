package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type Gateway struct {
	productServiceURL string
	orderServiceURL   string
	userServiceURL    string
}

func NewGateway() *Gateway {
	return &Gateway{
		productServiceURL: os.Getenv("PRODUCT_SERVICE_URL"),
		orderServiceURL:   os.Getenv("ORDER_SERVICE_URL"),
		userServiceURL:    os.Getenv("USER_SERVICE_URL"),
	}
}

func (g *Gateway) createProxy(targetURL string) *httputil.ReverseProxy {
	url, err := url.Parse(targetURL)
	if err != nil {
		log.Fatal(err)
	}
	return httputil.NewSingleHostReverseProxy(url)
}

func (g *Gateway) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/auth") {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (g *Gateway) setupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Product Service Routes
	productProxy := g.createProxy(g.productServiceURL)
	r.PathPrefix("/api/products").Handler(g.authMiddleware(productProxy))

	// Order Service Routes
	orderProxy := g.createProxy(g.orderServiceURL)
	r.PathPrefix("/api/orders").Handler(g.authMiddleware(orderProxy))

	// User Service Routes
	userProxy := g.createProxy(g.userServiceURL)
	r.PathPrefix("/api/auth").Handler(userProxy)
	r.PathPrefix("/api/users").Handler(g.authMiddleware(userProxy))

	// Health Check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func main() {
	gateway := NewGateway()
	router := gateway.setupRoutes()

	// Add middleware
	handler := corsMiddleware(loggingMiddleware(router))

	port := "8080"
	fmt.Printf("API Gateway starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
