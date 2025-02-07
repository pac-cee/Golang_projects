// Package calculator provides basic and advanced mathematical operations.
// It demonstrates proper package organization, documentation, and testing.
package calculator

import (
	"errors"
	"math"
)

// Common errors that can be returned by the calculator package
var (
	ErrDivisionByZero = errors.New("division by zero")
	ErrNegativeNumber = errors.New("negative number not allowed")
)

// Operation represents a mathematical operation
type Operation func(float64, float64) (float64, error)

// Calculator provides mathematical operations with optional logging
type Calculator struct {
	// EnableLogging determines if operations should be logged
	EnableLogging bool
	
	// precision specifies decimal places for rounding
	precision int
}

// NewCalculator creates a new Calculator with the specified precision
func NewCalculator(precision int) *Calculator {
	return &Calculator{
		precision: precision,
	}
}

// Add returns the sum of two numbers
func (c *Calculator) Add(a, b float64) float64 {
	return c.round(a + b)
}

// Subtract returns the difference between two numbers
func (c *Calculator) Subtract(a, b float64) float64 {
	return c.round(a - b)
}

// Multiply returns the product of two numbers
func (c *Calculator) Multiply(a, b float64) float64 {
	return c.round(a * b)
}

// Divide returns the quotient of two numbers
// Returns an error if dividing by zero
func (c *Calculator) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return c.round(a / b), nil
}

// Power returns a raised to the power of b
func (c *Calculator) Power(a, b float64) float64 {
	return c.round(math.Pow(a, b))
}

// SquareRoot returns the square root of a number
// Returns an error if the number is negative
func (c *Calculator) SquareRoot(a float64) (float64, error) {
	if a < 0 {
		return 0, ErrNegativeNumber
	}
	return c.round(math.Sqrt(a)), nil
}

// round rounds a number to the calculator's precision
func (c *Calculator) round(n float64) float64 {
	ratio := math.Pow(10, float64(c.precision))
	return math.Round(n*ratio) / ratio
}

// WithOperation performs a custom operation with logging
func (c *Calculator) WithOperation(name string, op Operation, a, b float64) (float64, error) {
	result, err := op(a, b)
	if err != nil {
		return 0, err
	}
	return c.round(result), nil
}
