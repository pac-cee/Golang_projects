package main

import (
    "encoding/json"
    "net/http"
)

type Message struct {
    Text string `json:"text"`
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(Message{Text: "Hello from Go REST API!"})
}

func main() {
    http.HandleFunc("/api/message", apiHandler)
    http.ListenAndServe(":8081", nil)
}
