package calculator

import (
	"math"
	"testing"
)

// TestCalculator tests basic calculator operations
func TestCalculator(t *testing.T) {
	// Create calculator with 2 decimal places precision
	calc := NewCalculator(2)

	// Test cases for basic operations
	tests := []struct {
		name     string
		a, b     float64
		op       func(float64, float64) float64
		expected float64
	}{
		{"addition", 3.14159, 2.0, calc.Add, 5.14},
		{"subtraction", 5.0, 3.14159, calc.Subtract, 1.86},
		{"multiplication", 2.0, 3.14159, calc.Multiply, 6.28},
		{"power", 2.0, 3.0, calc.Power, 8.00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.op(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("%s(%v, %v) = %v; want %v",
					tt.name, tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestDivision tests division including error cases
func TestDivision(t *testing.T) {
	calc := NewCalculator(2)

	tests := []struct {
		name        string
		a, b        float64
		expected    float64
		expectError bool
	}{
		{"normal division", 6.0, 2.0, 3.00, false},
		{"division by zero", 6.0, 0.0, 0.00, true},
		{"decimal division", 5.0, 2.0, 2.50, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Divide(tt.a, tt.b)

			// Check error cases
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				if err != ErrDivisionByZero {
					t.Errorf("expected ErrDivisionByZero but got %v", err)
				}
				return
			}

			// Check successful cases
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Divide(%v, %v) = %v; want %v",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestSquareRoot tests square root including error cases
func TestSquareRoot(t *testing.T) {
	calc := NewCalculator(2)

	tests := []struct {
		name        string
		input       float64
		expected    float64
		expectError bool
	}{
		{"perfect square", 9.0, 3.00, false},
		{"decimal number", 2.0, 1.41, false},
		{"negative number", -1.0, 0.00, true},
		{"zero", 0.0, 0.00, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.SquareRoot(tt.input)

			// Check error cases
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				if err != ErrNegativeNumber {
					t.Errorf("expected ErrNegativeNumber but got %v", err)
				}
				return
			}

			// Check successful cases
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("SquareRoot(%v) = %v; want %v",
					tt.input, result, tt.expected)
			}
		})
	}
}

// TestCustomOperation tests the WithOperation function
func TestCustomOperation(t *testing.T) {
	calc := NewCalculator(2)

	// Define a custom operation
	modulo := func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return math.Mod(a, b), nil
	}

	tests := []struct {
		name        string
		a, b        float64
		expected    float64
		expectError bool
	}{
		{"normal modulo", 5.0, 2.0, 1.00, false},
		{"zero modulo", 5.0, 0.0, 0.00, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.WithOperation("modulo", modulo, tt.a, tt.b)

			// Check error cases
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			// Check successful cases
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("modulo(%v, %v) = %v; want %v",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestPrecision tests different precision levels
func TestPrecision(t *testing.T) {
	tests := []struct {
		name      string
		precision int
		input     float64
		expected  float64
	}{
		{"zero precision", 0, 3.14159, 3.0},
		{"one decimal", 1, 3.14159, 3.1},
		{"two decimals", 2, 3.14159, 3.14},
		{"three decimals", 3, 3.14159, 3.142},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := NewCalculator(tt.precision)
			result := calc.round(tt.input)
			if result != tt.expected {
				t.Errorf("round(%v) with precision %d = %v; want %v",
					tt.input, tt.precision, result, tt.expected)
			}
		})
	}
}

// BenchmarkCalculator benchmarks basic calculator operations
func BenchmarkCalculator(b *testing.B) {
	calc := NewCalculator(2)

	operations := []struct {
		name string
		fn   func()
	}{
		{"Add", func() { calc.Add(3.14159, 2.0) }},
		{"Subtract", func() { calc.Subtract(5.0, 3.14159) }},
		{"Multiply", func() { calc.Multiply(2.0, 3.14159) }},
		{"Divide", func() { calc.Divide(6.0, 2.0) }},
		{"Power", func() { calc.Power(2.0, 3.0) }},
		{"SquareRoot", func() { calc.SquareRoot(9.0) }},
	}

	for _, op := range operations {
		b.Run(op.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				op.fn()
			}
		})
	}
}
