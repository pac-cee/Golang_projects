package main

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
)

type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/users", GetUsers).Methods("GET")
    r.HandleFunc("/users/{id}", GetUser).Methods("GET")

    log.Fatal(http.ListenAndServe(":8081", r))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
    users := []User{
        {ID: "1", Name: "John Doe", Email: "john@example.com"},
    }
    json.NewEncoder(w).Encode(users)
}