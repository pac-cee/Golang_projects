package main

import (
	"math"
	"strings"
	"testing"
)

// MockLogger for testing
type MockLogger struct {
	logs   []string
	errors []string
}

func (m *MockLogger) Log(message string) {
	m.logs = append(m.logs, message)
}

func (m *MockLogger) Error(message string) {
	m.errors = append(m.errors, message)
}

// TestShapes tests the Shape interface implementations
func TestShapes(t *testing.T) {
	t.Run("circle", func(t *testing.T) {
		c := Circle{radius: 5}
		expectedArea := math.Pi * 25
		expectedPerimeter := 2 * math.Pi * 5

		if area := c.Area(); math.Abs(area-expectedArea) > 0.001 {
			t.Errorf("Circle.Area() = %v, want %v", area, expectedArea)
		}

		if perim := c.Perimeter(); math.Abs(perim-expectedPerimeter) > 0.001 {
			t.Errorf("Circle.Perimeter() = %v, want %v", perim, expectedPerimeter)
		}

		if str := c.String(); !strings.Contains(str, "radius=5") {
			t.Errorf("Circle.String() = %v, want to contain 'radius=5'", str)
		}
	})

	t.Run("rectangle", func(t *testing.T) {
		r := Rectangle{width: 4, height: 6}
		expectedArea := 24.0
		expectedPerimeter := 20.0

		if area := r.Area(); area != expectedArea {
			t.Errorf("Rectangle.Area() = %v, want %v", area, expectedArea)
		}

		if perim := r.Perimeter(); perim != expectedPerimeter {
			t.Errorf("Rectangle.Perimeter() = %v, want %v", perim, expectedPerimeter)
		}

		if str := r.String(); !strings.Contains(str, "width=4") || 
			!strings.Contains(str, "height=6") {
			t.Errorf("Rectangle.String() = %v, want to contain 'width=4' and 'height=6'", str)
		}
	})
}

// TestCalculator tests the Calculator type
func TestCalculator(t *testing.T) {
	mock := &MockLogger{}
	calc := Calculator{logger: mock}

	tests := []struct {
		name     string
		a, b     float64
		op       Operation
		expected float64
	}{
		{"addition", 10, 5, add, 15},
		{"subtraction", 10, 5, subtract, 5},
		{"multiplication", 10, 5, multiply, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.Calculate(tt.a, tt.b, tt.op)
			if result != tt.expected {
				t.Errorf("Calculate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestDivide tests the divide function
func TestDivide(t *testing.T) {
	tests := []struct {
		name        string
		a, b        float64
		expected    float64
		expectError bool
	}{
		{"valid division", 10, 2, 5, false},
		{"division by zero", 10, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := divide(tt.a, tt.b)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("divide() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

// TestStringProcessor tests the StringProcessor type
func TestStringProcessor(t *testing.T) {
	mock := &MockLogger{}
	sp := StringProcessor{logger: mock}

	tests := []struct {
		name     string
		sep      string
		parts    []string
		expected string
	}{
		{
			name:     "comma separated",
			sep:      ",",
			parts:    []string{"a", "b", "c"},
			expected: "a,b,c",
		},
		{
			name:     "space separated",
			sep:      " ",
			parts:    []string{"hello", "world"},
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sp.Join(tt.sep, tt.parts...)
			if result != tt.expected {
				t.Errorf("Join() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestCounter tests the counter closure
func TestCounter(t *testing.T) {
	counter := createCounter(0)
	expected := []int{1, 2, 3}

	for i, want := range expected {
		if got := counter(); got != want {
			t.Errorf("counter() iteration %d = %v, want %v", i, got, want)
		}
	}
}

// BenchmarkCalculator benchmarks the Calculator operations
func BenchmarkCalculator(b *testing.B) {
	mock := &MockLogger{}
	calc := Calculator{logger: mock}

	operations := []struct {
		name string
		op   Operation
	}{
		{"add", add},
		{"subtract", subtract},
		{"multiply", multiply},
	}

	for _, op := range operations {
		b.Run(op.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				calc.Calculate(10, 5, op.op)
			}
		})
	}
}
