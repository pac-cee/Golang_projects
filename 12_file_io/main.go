package main

import (
    "fmt"
    "io/ioutil"
    "os"
)

func main() {
    data := []byte("Hello, File!")
    err := ioutil.WriteFile("test.txt", data, 0644)
    if err != nil {
        fmt.Println("Error writing file:", err)
        return
    }
    content, err := ioutil.ReadFile("test.txt")
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }
    fmt.Println("File content:", string(content))
}
