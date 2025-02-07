// Package main demonstrates the usage of the calculator package
package main

import (
	"fmt"
	"log"

	"go-learning-path/01-fundamentals/05-packages/pkg/calculator"
)

func main() {
	// Create a new calculator with 2 decimal places precision
	calc := calculator.NewCalculator(2)

	// Enable logging
	calc.EnableLogging = true

	// Basic operations
	a, b := 10.0, 5.0
	fmt.Printf("Numbers: a = %.2f, b = %.2f\n", a, b)
	fmt.Printf("Addition: %.2f\n", calc.Add(a, b))
	fmt.Printf("Subtraction: %.2f\n", calc.Subtract(a, b))
	fmt.Printf("Multiplication: %.2f\n", calc.Multiply(a, b))

	// Division with error handling
	if result, err := calc.Divide(a, b); err != nil {
		log.Printf("Division error: %v\n", err)
	} else {
		fmt.Printf("Division: %.2f\n", result)
	}

	// Power and square root
	fmt.Printf("Power (a^b): %.2f\n", calc.Power(a, b))
	
	if sqrt, err := calc.SquareRoot(a); err != nil {
		log.Printf("Square root error: %v\n", err)
	} else {
		fmt.Printf("Square root of %.2f: %.2f\n", a, sqrt)
	}

	// Custom operation using WithOperation
	// Define a custom operation that calculates percentage
	percentage := func(value, total float64) (float64, error) {
		if total == 0 {
			return 0, calculator.ErrDivisionByZero
		}
		return (value / total) * 100, nil
	}

	// Calculate what percentage a is of b
	if result, err := calc.WithOperation("percentage", percentage, a, b); err != nil {
		log.Printf("Percentage calculation error: %v\n", err)
	} else {
		fmt.Printf("%.2f is %.2f%% of %.2f\n", a, result, b)
	}

	// Error handling examples
	fmt.Println("\nError handling examples:")

	// Division by zero
	if _, err := calc.Divide(a, 0); err != nil {
		fmt.Printf("Division by zero error: %v\n", err)
	}

	// Square root of negative number
	if _, err := calc.SquareRoot(-a); err != nil {
		fmt.Printf("Negative square root error: %v\n", err)
	}
}
