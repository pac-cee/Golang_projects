// Package main demonstrates functions and methods in Go
package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// Shape interface defines common methods for shapes
type Shape interface {
	Area() float64
	Perimeter() float64
	String() string
}

// Circle implements Shape
type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle(radius=%.2f)", c.radius)
}

// Rectangle implements Shape
type Rectangle struct {
	width, height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rectangle(width=%.2f, height=%.2f)", r.width, r.height)
}

// Logger interface for different logging implementations
type Logger interface {
	Log(message string)
	Error(message string)
}

// ConsoleLogger implements Logger
type ConsoleLogger struct {
	prefix string
}

func (l ConsoleLogger) Log(message string) {
	fmt.Printf("[%s] LOG: %s\n", l.prefix, message)
}

func (l ConsoleLogger) Error(message string) {
	fmt.Printf("[%s] ERROR: %s\n", l.prefix, message)
}

// Calculator demonstrates function types and closures
type Calculator struct {
	logger Logger
}

// Operation represents a mathematical operation
type Operation func(a, b float64) float64

// Basic operations
func add(a, b float64) float64      { return a + b }
func subtract(a, b float64) float64 { return a - b }
func multiply(a, b float64) float64 { return a * b }

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Calculate performs the operation and logs the result
func (c Calculator) Calculate(a, b float64, op Operation) float64 {
	result := op(a, b)
	c.logger.Log(fmt.Sprintf("Calculated result: %.2f", result))
	return result
}

// StringProcessor demonstrates variadic functions and string manipulation
type StringProcessor struct {
	logger Logger
}

func (sp StringProcessor) Join(sep string, parts ...string) string {
	result := strings.Join(parts, sep)
	sp.logger.Log(fmt.Sprintf("Joined strings: %s", result))
	return result
}

// Counter demonstrates closure
func createCounter(start int) func() int {
	count := start
	return func() int {
		count++
		return count
	}
}

// TimeIt is a higher-order function for timing operations
func TimeIt(name string, logger Logger) func() {
	start := time.Now()
	return func() {
		duration := time.Since(start)
		logger.Log(fmt.Sprintf("%s took %v", name, duration))
	}
}

func main() {
	// Create a logger
	logger := ConsoleLogger{prefix: time.Now().Format("15:04:05")}

	// Demonstrate shape interface
	shapes := []Shape{
		Circle{radius: 5},
		Rectangle{width: 4, height: 6},
	}

	for _, shape := range shapes {
		logger.Log(fmt.Sprintf("Processing %s", shape))
		logger.Log(fmt.Sprintf("Area: %.2f", shape.Area()))
		logger.Log(fmt.Sprintf("Perimeter: %.2f", shape.Perimeter()))
	}

	// Demonstrate calculator
	calc := Calculator{logger: logger}
	
	// Basic operations
	logger.Log("\nPerforming calculations:")
	fmt.Printf("Addition: %.2f\n", calc.Calculate(10, 5, add))
	fmt.Printf("Subtraction: %.2f\n", calc.Calculate(10, 5, subtract))
	fmt.Printf("Multiplication: %.2f\n", calc.Calculate(10, 5, multiply))

	// Division with error handling
	if result, err := divide(10, 5); err != nil {
		logger.Error(fmt.Sprintf("Division error: %v", err))
	} else {
		fmt.Printf("Division: %.2f\n", result)
	}

	// Demonstrate string processor
	sp := StringProcessor{logger: logger}
	
	logger.Log("\nProcessing strings:")
	result := sp.Join(", ", "apple", "banana", "cherry")
	fmt.Printf("Joined string: %s\n", result)

	// Demonstrate counter closure
	counter := createCounter(0)
	logger.Log("\nDemonstrating counter:")
	for i := 0; i < 3; i++ {
		fmt.Printf("Count: %d\n", counter())
	}

	// Demonstrate timing function
	logger.Log("\nDemonstrating time measurement:")
	defer TimeIt("Main function", logger)()

	// Simulate some work
	time.Sleep(100 * time.Millisecond)
}
