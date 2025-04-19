package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Set DATABASE_URL env var to use Postgres, e.g.
// postgres://user:password@host:port/dbname?sslmode=disable

var db *gorm.DB

func getTodos(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.Header.Get("X-User-ID"))
	var todos []Todo
	db.Where("user_id = ?", userID).Find(&todos)
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.Header.Get("X-User-ID"))
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo.UserID = uint(userID)
	if err := db.Create(&todo).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.Header.Get("X-User-ID"))
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	var todo Todo
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		http.NotFound(w, r)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db.Save(&todo)
	json.NewEncoder(w).Encode(todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.Atoi(r.Header.Get("X-User-ID"))
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if err := db.Where("id = ? AND user_id = ?", id, userID).Delete(&Todo{}).Error; err != nil {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	var err error
	db = connectDB()
	if err != nil {
		log.Fatal("failed to connect database")
	}
	db.AutoMigrate(&User{}, &Todo{})

	r := mux.NewRouter()
	r.HandleFunc("/register", register(db)).Methods("POST")
	r.HandleFunc("/login", login(db)).Methods("POST")

	r.HandleFunc("/todos", jwtMiddleware(getTodos)).Methods("GET")
	r.HandleFunc("/todos", jwtMiddleware(createTodo)).Methods("POST")
	r.HandleFunc("/todos/{id}", jwtMiddleware(updateTodo)).Methods("PUT")
	r.HandleFunc("/todos/{id}", jwtMiddleware(deleteTodo)).Methods("DELETE")

	handler := cors.Default().Handler(r)
	log.Println("Backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
