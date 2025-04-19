package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"encoding/json"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&User{}, &Todo{})
	return db
}

func TestRegisterLoginAndTodoCRUD(t *testing.T) {
	db := setupTestDB()
	// Register
	r := httptest.NewRequest("POST", "/register", strings.NewReader(`{"username":"testuser","password":"testpass"}`))
	w := httptest.NewRecorder()
	register(db)(w, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	// Login
	r = httptest.NewRequest("POST", "/login", strings.NewReader(`{"username":"testuser","password":"testpass"}`))
	w = httptest.NewRecorder()
	login(db)(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp AuthResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp.Token == "" {
		t.Fatal("expected token")
	}
	// Create Todo
	todoBody := `{"task":"Test Todo","done":false}`
	r = httptest.NewRequest("POST", "/todos", strings.NewReader(todoBody))
	r.Header.Set("Authorization", "Bearer "+resp.Token)
	w = httptest.NewRecorder()
	createTodo := jwtMiddleware(createTodo)
	createTodo(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var todo Todo
	_ = json.NewDecoder(w.Body).Decode(&todo)
	if todo.Task != "Test Todo" {
		t.Fatalf("expected task 'Test Todo', got '%s'", todo.Task)
	}
	// Get Todos
	r = httptest.NewRequest("GET", "/todos", nil)
	r.Header.Set("Authorization", "Bearer "+resp.Token)
	w = httptest.NewRecorder()
	getTodos := jwtMiddleware(getTodos)
	getTodos(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var todos []Todo
	_ = json.NewDecoder(w.Body).Decode(&todos)
	if len(todos) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todos))
	}
	// Update Todo
	updateBody := `{"task":"Updated Todo","done":true}`
	r = httptest.NewRequest("PUT", "/todos/1", strings.NewReader(updateBody))
	r.Header.Set("Authorization", "Bearer "+resp.Token)
	w = httptest.NewRecorder()
	updateTodo := jwtMiddleware(updateTodo)
	updateTodo(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var updated Todo
	_ = json.NewDecoder(w.Body).Decode(&updated)
	if !updated.Done {
		t.Fatal("expected todo to be done")
	}
	// Delete Todo
	r = httptest.NewRequest("DELETE", "/todos/1", nil)
	r.Header.Set("Authorization", "Bearer "+resp.Token)
	w = httptest.NewRecorder()
	deleteTodo := jwtMiddleware(deleteTodo)
	deleteTodo(w, r)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}
