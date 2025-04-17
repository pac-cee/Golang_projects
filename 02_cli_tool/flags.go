package main

import (
    "flag"
    "fmt"
)

func PrintFlags() {
    age := flag.Int("age", 0, "Your age")
    location := flag.String("location", "", "Your location")
    flag.Parse()
    fmt.Printf("Age: %d, Location: %s\n", *age, *location)
}
