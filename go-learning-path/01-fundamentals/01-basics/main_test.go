package main

import (
	"strings"
	"testing"
)

// TestGreeting tests the Greeting function
func TestGreeting(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "normal name",
			input:    "Alice",
			contains: "Alice",
		},
		{
			name:     "empty string",
			input:    "",
			contains: "World",
		},
		{
			name:     "with spaces",
			input:    "  Bob  ",
			contains: "Bob",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Greeting(tt.input)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("Greeting(%q) = %q, want it to contain %q", 
					tt.input, result, tt.contains)
			}
		})
	}
}

// ExampleGreeting provides an example of using the Greeting function
func ExampleGreeting() {
	// Note: The output will vary based on time of day
	message := Greeting("Gopher")
	// Output will be one of:
	// Good morning, Gopher!
	// Good afternoon, Gopher!
	// Good evening, Gopher!
	_ = message
}

// BenchmarkGreeting benchmarks the Greeting function
func BenchmarkGreeting(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Greeting("Gopher")
	}
}
