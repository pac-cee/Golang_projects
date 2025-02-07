// Package main demonstrates basic Go programming concepts
package main

import (
	"fmt"
	"strings"
	"time"
)

// Greeting returns a greeting message based on the time of day
func Greeting(name string) string {
	hour := time.Now().Hour()
	var timeOfDay string

	switch {
	case hour < 12:
		timeOfDay = "morning"
	case hour < 17:
		timeOfDay = "afternoon"
	default:
		timeOfDay = "evening"
	}

	// Demonstrate string manipulation
	name = strings.TrimSpace(name)
	if name == "" {
		name = "World"
	}

	return fmt.Sprintf("Good %s, %s!", timeOfDay, name)
}

// PrintInfo demonstrates different ways to output information
func PrintInfo() {
	// Using Println
	fmt.Println("Welcome to Go Programming!")

	// Using Printf for formatted output
	version := 1.21
	fmt.Printf("Go Version: %.2f\n", version)

	// Using multiple values
	language, year := "Go", 2009
	fmt.Printf("%s was created in %d\n", language, year)
}

func main() {
	// Demonstrate basic function calls
	message := Greeting("Gopher")
	fmt.Println(message)

	// Print additional information
	PrintInfo()

	// Demonstrate basic string operations
	text := "  Go is fun!  "
	fmt.Printf("Original: %q\n", text)
	fmt.Printf("Trimmed: %q\n", strings.TrimSpace(text))
	fmt.Printf("Contains 'fun': %v\n", strings.Contains(text, "fun"))
	fmt.Printf("Uppercase: %s\n", strings.ToUpper(text))
}
