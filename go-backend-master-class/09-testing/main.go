package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
)

// User represents a user in our system
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserStore interface defines the methods for user storage
type UserStore interface {
	Create(user User) error
	Get(id int) (User, error)
	Update(user User) error
	Delete(id int) error
}

// InMemoryUserStore implements UserStore interface
type InMemoryUserStore struct {
	sync.RWMutex
	users map[int]User
}

// NewInMemoryUserStore creates a new in-memory user store
func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[int]User),
	}
}

// Create adds a new user
func (s *InMemoryUserStore) Create(user User) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.users[user.ID]; exists {
		return errors.New("user already exists")
	}

	s.users[user.ID] = user
	return nil
}

// Get retrieves a user by ID
func (s *InMemoryUserStore) Get(id int) (User, error) {
	s.RLock()
	defer s.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

// Update modifies an existing user
func (s *InMemoryUserStore) Update(user User) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.users[user.ID]; !exists {
		return errors.New("user not found")
	}

	s.users[user.ID] = user
	return nil
}

// Delete removes a user
func (s *InMemoryUserStore) Delete(id int) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(s.users, id)
	return nil
}

// UserService handles user-related business logic
type UserService struct {
	store UserStore
}

// NewUserService creates a new user service
func NewUserService(store UserStore) *UserService {
	return &UserService{store: store}
}

// CreateUser creates a new user with validation
func (s *UserService) CreateUser(user User) error {
	if err := validateUser(user); err != nil {
		return err
	}
	return s.store.Create(user)
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id int) (User, error) {
	return s.store.Get(id)
}

// UpdateUser updates a user with validation
func (s *UserService) UpdateUser(user User) error {
	if err := validateUser(user); err != nil {
		return err
	}
	return s.store.Update(user)
}

// DeleteUser removes a user
func (s *UserService) DeleteUser(id int) error {
	return s.store.Delete(id)
}

// validateUser performs user validation
func validateUser(user User) error {
	if user.Username == "" {
		return errors.New("username is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

// UserHandler handles HTTP requests for users
type UserHandler struct {
	service *UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser handles user creation requests
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser handles user retrieval requests
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// In a real application, you would parse the ID from the URL
	id := 1 // Example ID

	user, err := h.service.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func main() {
	// Initialize dependencies
	store := NewInMemoryUserStore()
	service := NewUserService(store)
	handler := NewUserHandler(service)

	// Setup routes
	http.HandleFunc("/users", handler.CreateUser)
	http.HandleFunc("/users/1", handler.GetUser)

	// Start server
	http.ListenAndServe(":8080", nil)
}
