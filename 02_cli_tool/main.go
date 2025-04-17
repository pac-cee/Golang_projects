package main

import (
    "flag"
    "fmt"
)

func main() {
    name := flag.String("name", "Go Developer", "Your name")
    flag.Parse()
    fmt.Printf("Hello, %s!\n", *name)
}
