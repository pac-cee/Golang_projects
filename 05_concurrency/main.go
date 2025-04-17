package main

import (
    "fmt"
    "time"
)

func printNumbers() {
    for i := 1; i <= 5; i++ {
        fmt.Println(i)
        time.Sleep(500 * time.Millisecond)
    }
}

func main() {
    go printNumbers()
    fmt.Println("Goroutine started!")
    time.Sleep(3 * time.Second)
}
