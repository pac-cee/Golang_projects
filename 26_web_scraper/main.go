package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
)

func main() {
    resp, err := http.Get("https://example.com")
    if err != nil {
        fmt.Println("Error fetching page:", err)
        return
    }
    defer resp.Body.Close()
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Printf("Fetched %d bytes\n", len(body))
}
