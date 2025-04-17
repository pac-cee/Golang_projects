package main

import (
    "fmt"
    "sync"
)

func main() {
    var once sync.Once
    once.Do(func() {
        fmt.Println("This will only print once!")
    })
    once.Do(func() {
        fmt.Println("This will NOT print!")
    })
}
