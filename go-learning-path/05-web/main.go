// Package main demonstrates web development concepts in Go
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Todo represents a todo item
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// TodoStore manages todo items
type TodoStore struct {
	sync.RWMutex
	todos map[int]*Todo
	seq   int
}

// NewTodoStore creates a new todo store
func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos: make(map[int]*Todo),
	}
}

// Create adds a new todo
func (s *TodoStore) Create(title string) *Todo {
	s.Lock()
	defer s.Unlock()

	s.seq++
	todo := &Todo{
		ID:        s.seq,
		Title:     title,
		CreatedAt: time.Now(),
	}
	s.todos[todo.ID] = todo
	return todo
}

// Get retrieves a todo by ID
func (s *TodoStore) Get(id int) (*Todo, bool) {
	s.RLock()
	defer s.RUnlock()
	todo, ok := s.todos[id]
	return todo, ok
}

// List returns all todos
func (s *TodoStore) List() []*Todo {
	s.RLock()
	defer s.RUnlock()
	todos := make([]*Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}
	return todos
}

// Update modifies a todo
func (s *TodoStore) Update(id int, completed bool) bool {
	s.Lock()
	defer s.Unlock()
	if todo, ok := s.todos[id]; ok {
		todo.Completed = completed
		return true
	}
	return false
}

// Delete removes a todo
func (s *TodoStore) Delete(id int) bool {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.todos[id]; ok {
		delete(s.todos, id)
		return true
	}
	return false
}

// Server represents the web server
type Server struct {
	router    *mux.Router
	todos     *TodoStore
	templates *template.Template
	upgrader  websocket.Upgrader
	clients   map[*websocket.Conn]bool
	broadcast chan Message
}

// Message represents a WebSocket message
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	User    string `json:"user"`
}

// NewServer creates a new web server
func NewServer() *Server {
	s := &Server{
		router:  mux.NewRouter(),
		todos:   NewTodoStore(),
		clients: make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // In production, check origin
			},
		},
	}

	// Parse templates
	var err error
	s.templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	// Setup routes
	s.routes()

	return s
}

func (s *Server) routes() {
	// Static files
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Web pages
	s.router.HandleFunc("/", s.handleHome).Methods("GET")
	s.router.HandleFunc("/chat", s.handleChat).Methods("GET")

	// API endpoints
	api := s.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/todos", s.handleGetTodos).Methods("GET")
	api.HandleFunc("/todos", s.handleCreateTodo).Methods("POST")
	api.HandleFunc("/todos/{id:[0-9]+}", s.handleUpdateTodo).Methods("PUT")
	api.HandleFunc("/todos/{id:[0-9]+}", s.handleDeleteTodo).Methods("DELETE")

	// WebSocket
	s.router.HandleFunc("/ws", s.handleWebSocket)
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
		Todos []*Todo
	}{
		Title: "Todo App",
		Todos: s.todos.List(),
	}
	s.templates.ExecuteTemplate(w, "index.html", data)
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	s.templates.ExecuteTemplate(w, "chat.html", nil)
}

func (s *Server) handleGetTodos(w http.ResponseWriter, r *http.Request) {
	todos := s.todos.List()
	json.NewEncoder(w).Encode(todos)
}

func (s *Server) handleCreateTodo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo := s.todos.Create(req.Title)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (s *Server) handleUpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := 0
	fmt.Sscanf(vars["id"], "%d", &id)

	var req struct {
		Completed bool `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !s.todos.Update(id, req.Completed) {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleDeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := 0
	fmt.Sscanf(vars["id"], "%d", &id)

	if !s.todos.Delete(id) {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	s.clients[conn] = true

	// Handle WebSocket messages
	go func() {
		defer func() {
			conn.Close()
			delete(s.clients, conn)
		}()

		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
				break
			}
			s.broadcast <- msg
		}
	}()
}

func (s *Server) handleMessages() {
	for msg := range s.broadcast {
		for client := range s.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(s.clients, client)
			}
		}
	}
}

// Start runs the web server
func (s *Server) Start(port string) error {
	// Start WebSocket message handler
	go s.handleMessages()

	// Start HTTP server
	log.Printf("Server starting on port %s", port)
	return http.ListenAndServe(port, s.router)
}

func main() {
	// Create server
	server := NewServer()

	// Add some sample todos
	server.todos.Create("Learn Go")
	server.todos.Create("Build a web app")
	server.todos.Create("Write tests")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	// Start server
	log.Fatal(server.Start(port))
}
