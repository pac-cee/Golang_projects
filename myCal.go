package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nSimple Calculator")
		fmt.Println("1. Add")
		fmt.Println("2. Subtract")
		fmt.Println("3. Multiply")
		fmt.Println("4. Divide")
		fmt.Println("5. Exit")
		fmt.Print("Choose operation (1-5): ")

		scanner.Scan()
		choice := scanner.Text()

		if choice == "5" {
			fmt.Println("Goodbye!")
			break
		}

		if !isValidChoice(choice) {
			fmt.Println("Invalid choice! Please choose 1-5")
			continue
		}

		fmt.Print("Enter first number: ")
		scanner.Scan()
		num1, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			fmt.Println("Invalid number!")
			continue
		}

		fmt.Print("Enter second number: ")
		scanner.Scan()
		num2, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			fmt.Println("Invalid number!")
			continue
		}

		result := calculate(choice, num1, num2)
		fmt.Printf("Result: %.2f\n", result)
	}
}

func isValidChoice(choice string) bool {
	validChoices := []string{"1", "2", "3", "4"}
	for _, c := range validChoices {
		if choice == c {
			return true
		}
	}
	return false
}

func calculate(operation string, a, b float64) float64 {
	// Using map for operations
	ops := map[string]func(float64, float64) float64{
		"1": func(a, b float64) float64 { return a + b },
		"2": func(a, b float64) float64 { return a - b },
		"3": func(a, b float64) float64 { return a * b },
		"4": func(a, b float64) float64 {
			if b == 0 {
				fmt.Println("Error: Division by zero!")
				return 0
			}
			return a / b
		},
	}

	if fn, ok := ops[operation]; ok {
		return fn(a, b)
	}
	return 0
}
