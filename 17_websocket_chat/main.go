package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "WebSocket endpoint placeholder")
    })
    fmt.Println("WebSocket chat server running on :8083")
    http.ListenAndServe(":8083", nil)
}
