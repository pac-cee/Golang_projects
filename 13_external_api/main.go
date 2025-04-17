package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

type Joke struct {
    Setup   string `json:"setup"`
    Punchline string `json:"punchline"`
}

func main() {
    resp, err := http.Get("https://official-joke-api.appspot.com/random_joke")
    if err != nil {
        fmt.Println("Error fetching joke:", err)
        return
    }
    defer resp.Body.Close()
    var joke Joke
    json.NewDecoder(resp.Body).Decode(&joke)
    fmt.Printf("Joke: %s - %s\n", joke.Setup, joke.Punchline)
}
