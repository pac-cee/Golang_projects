package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "GraphQL endpoint placeholder")
    })
    fmt.Println("GraphQL API running on :8082")
    http.ListenAndServe(":8082", nil)
}
