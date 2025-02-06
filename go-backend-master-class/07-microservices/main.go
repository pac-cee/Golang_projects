package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

// Service represents a microservice
type Service struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Status   string `json:"status"`
	LastPing time.Time `json:"last_ping"`
}

// ServiceRegistry keeps track of available services
type ServiceRegistry struct {
	Services map[string]Service
}

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		Services: make(map[string]Service),
	}
}

// RegisterService registers a new service
func (sr *ServiceRegistry) RegisterService(w http.ResponseWriter, r *http.Request) {
	var service Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		respondWithJSON(w, http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid request payload",
		})
		return
	}

	service.Status = "healthy"
	service.LastPing = time.Now()
	sr.Services[service.Name] = service

	respondWithJSON(w, http.StatusCreated, Response{
		Status:  "success",
		Message: "Service registered successfully",
		Data:    service,
	})
}

// GetServices returns all registered services
func (sr *ServiceRegistry) GetServices(w http.ResponseWriter, r *http.Request) {
	services := make([]Service, 0, len(sr.Services))
	for _, service := range sr.Services {
		services = append(services, service)
	}

	respondWithJSON(w, http.StatusOK, Response{
		Status: "success",
		Data:   services,
	})
}

// HealthCheck handles service health check
func (sr *ServiceRegistry) HealthCheck(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	serviceName := vars["name"]

	service, exists := sr.Services[serviceName]
	if !exists {
		respondWithJSON(w, http.StatusNotFound, Response{
			Status:  "error",
			Message: "Service not found",
		})
		return
	}

	// Update last ping
	service.LastPing = time.Now()
	sr.Services[serviceName] = service

	respondWithJSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: "Service is healthy",
		Data:    service,
	})
}

// ServiceClient represents a client for inter-service communication
type ServiceClient struct {
	registry *ServiceRegistry
	client   *http.Client
}

// NewServiceClient creates a new service client
func NewServiceClient(registry *ServiceRegistry) *ServiceClient {
	return &ServiceClient{
		registry: registry,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// CallService makes a request to another service
func (sc *ServiceClient) CallService(serviceName string, path string) (*http.Response, error) {
	service, exists := sc.registry.Services[serviceName]
	if !exists {
		return nil, fmt.Errorf("service %s not found", serviceName)
	}

	url := service.URL + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return sc.client.Do(req)
}

// CircuitBreaker implements basic circuit breaker pattern
type CircuitBreaker struct {
	failures  int
	threshold int
	timeout   time.Duration
	lastError time.Time
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		threshold: threshold,
		timeout:   timeout,
	}
}

// Execute runs the given function with circuit breaker protection
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if cb.failures >= cb.threshold {
		if time.Since(cb.lastError) > cb.timeout {
			// Try again after timeout
			cb.failures = 0
		} else {
			return fmt.Errorf("circuit breaker is open")
		}
	}

	if err := fn(); err != nil {
		cb.failures++
		cb.lastError = time.Now()
		return err
	}

	cb.failures = 0
	return nil
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	// Initialize service registry
	registry := NewServiceRegistry()
	client := NewServiceClient(registry)
	circuitBreaker := NewCircuitBreaker(3, time.Minute)

	router := mux.NewRouter()

	// Service registry endpoints
	router.HandleFunc("/services", registry.RegisterService).Methods("POST")
	router.HandleFunc("/services", registry.GetServices).Methods("GET")
	router.HandleFunc("/services/{name}/health", registry.HealthCheck).Methods("GET")

	// Example service endpoint with circuit breaker
	router.HandleFunc("/example", func(w http.ResponseWriter, r *http.Request) {
		err := circuitBreaker.Execute(func() error {
			// Example service call
			resp, err := client.CallService("example-service", "/api/data")
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// Process response...
			return nil
		})

		if err != nil {
			respondWithJSON(w, http.StatusServiceUnavailable, Response{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		respondWithJSON(w, http.StatusOK, Response{
			Status:  "success",
			Message: "Service call successful",
		})
	})

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Configure and start server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Service Registry starting on http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
