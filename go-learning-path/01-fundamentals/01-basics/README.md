# Go Basics: Getting Started 🚀

Welcome to the first module of our Go programming journey! Let's start with the fundamentals.

## 📖 Table of Contents
1. [Program Structure](#program-structure)
2. [Basic Syntax](#basic-syntax)
3. [Comments and Documentation](#comments-and-documentation)
4. [First Program](#first-program)
5. [Running Go Programs](#running-go-programs)
6. [Go Modules](#go-modules)
7. [Exercises](#exercises)

## Program Structure

Every Go program follows a basic structure:

```go
// Package declaration
package main

// Import statements
import (
    "fmt"
    "strings"
)

// Main function - entry point
func main() {
    // Your code here
}
```

### Key Components:
1. **Package Declaration**: Every Go file starts with a package declaration
2. **Import Statements**: Required packages are imported
3. **Functions**: Code is organized into functions
4. **Main Package**: Executable programs must have a main package and main function

## Basic Syntax

### 1. Statements
- No semicolons required (inserted automatically)
- One statement per line (preferred)
```go
fmt.Println("Hello")
fmt.Println("World")
```

### 2. Code Blocks
- Defined by curly braces
- Opening brace must be on the same line
```go
func example() {
    // This is correct
}

func incorrect() 
{
    // This will cause an error
}
```

### 3. Naming Conventions
```go
// Package names: lowercase, single word
package userservice

// Variable names: camelCase
var userName string

// Constants: camelCase or UPPERCASE for exports
const MaxLength = 50
const defaultTimeout = 30

// Function names: camelCase
func getUserData() {}

// Interface names: usually ends with -er
type Reader interface {}
```

## Comments and Documentation

### 1. Single Line Comments
```go
// This is a single line comment
```

### 2. Multi-line Comments
```go
/* This is a multi-line comment
   It can span multiple lines
   Used for longer explanations
*/
```

### 3. Package Documentation
```go
// Package userservice provides user management functionality.
// It handles user creation, authentication, and authorization.
package userservice
```

### 4. Function Documentation
```go
// GetUser retrieves a user by their ID.
// It returns an error if the user is not found.
func GetUser(id string) (*User, error) {
    // Implementation
}
```

## First Program

Let's create your first Go program:

```go
// hello.go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

This program demonstrates:
1. Package declaration (`main`)
2. Import statement (`fmt` package)
3. Main function
4. Function call (`fmt.Println`)

## Running Go Programs

### 1. Direct Execution
```bash
# Run directly
go run hello.go

# Build and execute
go build hello.go
./hello  # or hello.exe on Windows
```

### 2. Module-based Execution
```bash
# Initialize module
go mod init hello

# Run program
go run .
```

### 3. Common Commands
```bash
go run    # Run program
go build  # Build executable
go fmt    # Format code
go test   # Run tests
go mod    # Module management
go get    # Download dependencies
```

## Go Modules

Modern Go uses modules for dependency management:

### 1. Creating a Module
```bash
mkdir myproject
cd myproject
go mod init myproject
```

### 2. go.mod File
```go
module myproject

go 1.21

require (
    github.com/some/dependency v1.2.3
)
```

### 3. Adding Dependencies
```bash
go get github.com/some/dependency
```

### 4. Updating Dependencies
```bash
go get -u        # Update all dependencies
go mod tidy      # Clean up dependencies
```

## Exercises

### Exercise 1: Hello World
Create a program that prints "Hello, [your name]!"
```go
// Solution
package main

import "fmt"

func main() {
    name := "Alice"
    fmt.Printf("Hello, %s!\n", name)
}
```

### Exercise 2: Multiple Prints
Create a program that prints three different lines using different print methods
```go
// Solution
package main

import "fmt"

func main() {
    // Using Println
    fmt.Println("Line 1")
    
    // Using Printf
    fmt.Printf("Line %d\n", 2)
    
    // Using Print
    fmt.Print("Line 3\n")
}
```

### Exercise 3: Documentation
Write a properly documented function that adds two numbers
```go
// Solution
package main

import "fmt"

// Add takes two integers and returns their sum.
// It demonstrates proper function documentation.
func Add(a, b int) int {
    return a + b
}

func main() {
    result := Add(5, 3)
    fmt.Printf("5 + 3 = %d\n", result)
}
```

## Best Practices

1. **Code Organization**
   - One package per directory
   - Clear package names
   - Logical file organization

2. **Documentation**
   - Document all exported items
   - Write clear, concise comments
   - Include examples in docs

3. **Module Management**
   - Use go modules
   - Keep dependencies updated
   - Clean up unused dependencies

4. **Code Style**
   - Follow Go conventions
   - Use `go fmt`
   - Keep functions small and focused

## Next Steps
- Practice writing simple programs
- Experiment with different packages
- Read Go documentation
- Move on to Variables and Types

## Additional Resources
- [Go Tour](https://tour.golang.org/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Style Guide](https://github.com/golang/go/wiki/CodeReviewComments)
