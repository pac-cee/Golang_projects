package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

// Product represents a product in our system
type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	LastUpdated time.Time `json:"last_updated"`
}

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Source  string      `json:"source,omitempty"` // "cache" or "database"
}

// Cache interface defines our caching behavior
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

// RedisCache implements Cache interface using Redis
type RedisCache struct {
	client *redis.Client
}

// InMemoryCache implements Cache interface using a map
type InMemoryCache struct {
	sync.RWMutex
	data map[string]string
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(addr string) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{client: client}
}

// NewInMemoryCache creates a new in-memory cache
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]string),
	}
}

// Redis Cache implementation
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// In-Memory Cache implementation
func (c *InMemoryCache) Get(ctx context.Context, key string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	if val, ok := c.data[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key not found")
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	c.Lock()
	defer c.Unlock()
	str, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.data[key] = string(str)
	return nil
}

func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	c.Lock()
	defer c.Unlock()
	delete(c.data, key)
	return nil
}

// API represents our REST API
type API struct {
	cache  Cache
	store  map[string]Product // Simulated database
	mutex  sync.RWMutex
	ctx    context.Context
}

// NewAPI creates a new API instance
func NewAPI(cache Cache) *API {
	return &API{
		cache:  cache,
		store:  make(map[string]Product),
		ctx:    context.Background(),
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// CreateProduct handles product creation
func (api *API) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	product.LastUpdated = time.Now()

	// Store in "database"
	api.mutex.Lock()
	api.store[product.ID] = product
	api.mutex.Unlock()

	// Cache the product
	if err := api.cache.Set(api.ctx, fmt.Sprintf("product:%s", product.ID), product, time.Hour); err != nil {
		log.Printf("Error caching product: %v", err)
	}

	respondWithJSON(w, http.StatusCreated, Response{
		Status:  "success",
		Message: "Product created successfully",
		Data:    product,
	})
}

// GetProduct handles getting a single product
func (api *API) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	cacheKey := fmt.Sprintf("product:%s", id)

	// Try to get from cache first
	if cached, err := api.cache.Get(api.ctx, cacheKey); err == nil {
		var product Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			respondWithJSON(w, http.StatusOK, Response{
				Status: "success",
				Data:   product,
				Source: "cache",
			})
			return
		}
	}

	// If not in cache, get from "database"
	api.mutex.RLock()
	product, exists := api.store[id]
	api.mutex.RUnlock()

	if !exists {
		respondWithJSON(w, http.StatusNotFound, Response{
			Status:  "error",
			Message: "Product not found",
		})
		return
	}

	// Cache the product for future requests
	if err := api.cache.Set(api.ctx, cacheKey, product, time.Hour); err != nil {
		log.Printf("Error caching product: %v", err)
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status: "success",
		Data:   product,
		Source: "database",
	})
}

// UpdateProduct handles updating a product
func (api *API) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	api.mutex.Lock()
	if _, exists := api.store[id]; !exists {
		api.mutex.Unlock()
		respondWithJSON(w, http.StatusNotFound, Response{
			Status:  "error",
			Message: "Product not found",
		})
		return
	}

	product.ID = id
	product.LastUpdated = time.Now()
	api.store[id] = product
	api.mutex.Unlock()

	// Update cache
	cacheKey := fmt.Sprintf("product:%s", id)
	if err := api.cache.Set(api.ctx, cacheKey, product, time.Hour); err != nil {
		log.Printf("Error updating cache: %v", err)
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Product updated successfully",
		Data:    product,
	})
}

func main() {
	// Initialize cache (using in-memory cache for demonstration)
	cache := NewInMemoryCache()
	api := NewAPI(cache)

	router := mux.NewRouter()

	// Register routes
	router.HandleFunc("/products", api.CreateProduct).Methods("POST")
	router.HandleFunc("/products/{id}", api.GetProduct).Methods("GET")
	router.HandleFunc("/products/{id}", api.UpdateProduct).Methods("PUT")

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
