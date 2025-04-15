package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Todo represents a single to-do item.
// Struct tags (e.g. `json:"id"`) tell Go how to encode/decode fields to/from JSON.
type Todo struct {
	ID   int    `json:"id"`   // Unique identifier
	Task string `json:"task"` // Task description
	Done bool   `json:"done"` // Completion status
}

// In-memory slice to store to-do items. In real apps, you'd use a database.
var todos = []Todo{
	{ID: 1, Task: "Learn Go basics", Done: false},
	{ID: 2, Task: "Build a REST API", Done: false},
}

// getTodos handles GET requests to /todos.
// It encodes the todos slice as JSON and writes it to the response.
func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Set response type
	json.NewEncoder(w).Encode(todos) // Encode todos as JSON
}

func main() {
	// Register the handler function for the /todos route.
	http.HandleFunc("/todos", getTodos)
	// Print a message so you know the server is running.
	fmt.Println("Server running at http://localhost:8080/")
	// Start the HTTP server. log.Fatal will print any errors if the server fails.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
