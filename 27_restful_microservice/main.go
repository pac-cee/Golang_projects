package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "OK")
    })
    fmt.Println("RESTful microservice running on :8090")
    http.ListenAndServe(":8090", nil)
}
