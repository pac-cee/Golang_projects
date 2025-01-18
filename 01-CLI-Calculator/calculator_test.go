package main

import (
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"5 + 3", 8, false},
		{"10 - 5", 5, false},
		{"4 * 3", 12, false},
		{"15 / 3", 5, false},
		{"10 / 0", 0, true},     // Division by zero
		{"abc + 3", 0, true},    // Invalid number
		{"5 & 3", 0, true},      // Invalid operator
		{"5", 0, true},          // Invalid format
		{"5 + 5 + 5", 0, true}, // Too many operands
	}

	for _, test := range tests {
		result, err := calculate(test.input)
		
		// Check error cases
		if test.hasError && err == nil {
			t.Errorf("Expected error for input %s, but got none", test.input)
			continue
		}
		if !test.hasError && err != nil {
			t.Errorf("Unexpected error for input %s: %v", test.input, err)
			continue
		}

		// Check results for non-error cases
		if !test.hasError && result != test.expected {
			t.Errorf("For input %s, expected %.2f but got %.2f", 
				test.input, test.expected, result)
		}
	}
}
