# Simple CLI Calculator

A basic command-line calculator implemented in Go to demonstrate fundamental programming concepts.

## Concepts Covered

- Basic Go syntax and structure
- Variables and data types
- User input/output
- Control structures (if/else, switch, loops)
- Functions
- Error handling
- Type conversion

## Features

- Basic arithmetic operations:
  - Addition
  - Subtraction
  - Multiplication
  - Division
- Input validation
- Error handling for:
  - Invalid operations
  - Invalid numbers
  - Division by zero

## How to Run

1. Make sure you have Go installed on your system
2. Navigate to the project directory:
   ```bash
   cd level-1/01-calculator
   ```
3. Run the program:
   ```bash
   go run main.go
   ```

## Usage Example

```
Simple Calculator
1. Add
2. Subtract
3. Multiply
4. Divide
5. Exit
Choose operation (1-5): 1
Enter first number: 5
Enter second number: 3
Result: 8.00
```

## Project Structure

```
01-calculator/
├── main.go      # Main program file
└── README.md    # Project documentation
```

## Learning Objectives

- Understanding basic Go program structure
- Working with the standard library (fmt, bufio, os, strconv)
- Implementing user input validation
- Basic error handling
- Creating modular code with functions

## Next Steps

To extend this project, you could:
1. Add more operations (square root, power, etc.)
2. Implement memory functions (store/recall results)
3. Add support for more complex expressions
4. Create a basic test suite
