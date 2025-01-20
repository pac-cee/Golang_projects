package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Operation represents a calculation operation
type Operation struct {
	FirstNumber  float64
	SecondNumber float64
	Operator     string
	Result       float64
}

// History stores the last 5 calculations
var history []Operation

func main() {
	fmt.Println("Welcome to the Go Calculator!")
	fmt.Println("Enter calculations in the format: number operator number")
	fmt.Println("Example: 5 + 3")
	fmt.Println("Available operators: +, -, *, /")
	fmt.Println("Type 'exit' to quit or 'history' to see past calculations")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nEnter calculation: ")
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			break
		}

		if input == "history" {
			showHistory()
			continue
		}

		result, err := calculate(input)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		fmt.Printf("Result: %.2f\n", result)
	}
}

func calculate(input string) (float64, error) {
	parts := strings.Fields(input)
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid format: please use 'number operator number'")
	}

	num1, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid first number")
	}

	num2, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid second number")
	}

	var result float64
	switch parts[1] {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "*":
		result = num1 * num2
	case "/":
		if num2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		result = num1 / num2
	default:
		return 0, fmt.Errorf("invalid operator: use +, -, *, /")
	}

	// Store calculation in history
	operation := Operation{
		FirstNumber:  num1,
		SecondNumber: num2,
		Operator:     parts[1],
		Result:       result,
	}
	
	if len(history) >= 5 {
		history = history[1:]
	}
	history = append(history, operation)

	return result, nil
}

func showHistory() {
	if len(history) == 0 {
		fmt.Println("No calculations in history")
		return
	}

	fmt.Println("\nCalculation History:")
	for i, op := range history {
		fmt.Printf("%d: %.2f %s %.2f = %.2f\n",
			i+1, op.FirstNumber, op.Operator, op.SecondNumber, op.Result)
	}
}
