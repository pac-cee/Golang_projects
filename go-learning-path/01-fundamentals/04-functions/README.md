# Functions and Methods in Go 🔧

## 📖 Table of Contents
1. [Function Basics](#function-basics)
2. [Methods](#methods)
3. [Interfaces](#interfaces)
4. [Function Types](#function-types)
5. [Error Handling](#error-handling)
6. [Best Practices](#best-practices)
7. [Exercises](#exercises)

## Function Basics

### Basic Function Declaration
```go
func functionName(param1 Type1, param2 Type2) ReturnType {
    // function body
    return value
}
```

### Multiple Return Values
```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

### Named Return Values
```go
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return  // naked return
}
```

### Variadic Functions
```go
func sum(nums ...int) int {
    total := 0
    for _, num := range nums {
        total += num
    }
    return total
}
```

## Methods

### Method Declaration
```go
type Rectangle struct {
    width, height float64
}

func (r Rectangle) Area() float64 {
    return r.width * r.height
}
```

### Pointer Receivers
```go
func (r *Rectangle) Scale(factor float64) {
    r.width *= factor
    r.height *= factor
}
```

### Value vs Pointer Receivers
```go
// Value receiver - gets a copy
func (r Rectangle) Area() float64 {
    return r.width * r.height
}

// Pointer receiver - can modify the original
func (r *Rectangle) SetWidth(w float64) {
    r.width = w
}
```

## Interfaces

### Interface Declaration
```go
type Shape interface {
    Area() float64
    Perimeter() float64
}
```

### Implementing Interfaces
```go
type Circle struct {
    radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.radius
}
```

### Empty Interface
```go
func PrintAnything(v interface{}) {
    fmt.Printf("Type: %T, Value: %v\n", v, v)
}
```

## Function Types

### Function as a Type
```go
type Operation func(a, b int) int

func Calculate(a, b int, op Operation) int {
    return op(a, b)
}
```

### Anonymous Functions
```go
func main() {
    sum := func(a, b int) int {
        return a + b
    }
    
    result := sum(3, 4)
}
```

### Closures
```go
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}
```

## Error Handling

### Error Return Pattern
```go
func processFile(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()
    
    return nil
}
```

### Custom Errors
```go
type ValidationError struct {
    Field string
    Error string
}

func (v *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", v.Field, v.Error)
}
```

### Error Wrapping
```go
func processData() error {
    err := doSomething()
    if err != nil {
        return fmt.Errorf("processing failed: %w", err)
    }
    return nil
}
```

## Best Practices

### 1. Function Design
- Keep functions focused and small
- Use meaningful parameter names
- Return early for errors
- Use named returns judiciously

### 2. Method Design
- Choose receiver type appropriately
- Be consistent with receiver names
- Use pointer receivers for mutations
- Use value receivers for transformations

### 3. Interface Design
- Keep interfaces small
- Follow the interface segregation principle
- Design for behavior, not objects
- Use composition over inheritance

### 4. Error Handling
- Always check errors
- Wrap errors with context
- Use custom error types when needed
- Don't panic in libraries

## Exercises

### Exercise 1: Function Basics
```go
// Implement a function that takes a slice of integers
// and returns the sum and average
func calculateStats(numbers []int) (sum int, avg float64) {
    for _, num := range numbers {
        sum += num
    }
    avg = float64(sum) / float64(len(numbers))
    return
}
```

### Exercise 2: Methods
```go
// Implement a BankAccount type with methods
type BankAccount struct {
    balance float64
}

func (b *BankAccount) Deposit(amount float64) error {
    if amount <= 0 {
        return errors.New("invalid deposit amount")
    }
    b.balance += amount
    return nil
}
```

### Exercise 3: Interfaces
```go
// Implement a Logger interface
type Logger interface {
    Log(message string)
    Error(message string)
}

type ConsoleLogger struct{}

func (l ConsoleLogger) Log(message string) {
    fmt.Printf("LOG: %s\n", message)
}

func (l ConsoleLogger) Error(message string) {
    fmt.Printf("ERROR: %s\n", message)
}
```

## Common Patterns

### 1. Options Pattern
```go
type ServerOption func(*Server)

func WithPort(port int) ServerOption {
    return func(s *Server) {
        s.port = port
    }
}

func NewServer(options ...ServerOption) *Server {
    s := &Server{}
    for _, opt := range options {
        opt(s)
    }
    return s
}
```

### 2. Builder Pattern
```go
type QueryBuilder struct {
    table string
    where string
    limit int
}

func (q *QueryBuilder) From(table string) *QueryBuilder {
    q.table = table
    return q
}

func (q *QueryBuilder) Where(condition string) *QueryBuilder {
    q.where = condition
    return q
}
```

### 3. Factory Pattern
```go
type PaymentMethod interface {
    Pay(amount float64) error
}

func NewPaymentMethod(method string) PaymentMethod {
    switch method {
    case "credit":
        return &CreditCardPayment{}
    case "debit":
        return &DebitCardPayment{}
    default:
        return &CashPayment{}
    }
}
```

## Next Steps
- Practice writing clean functions
- Implement common interfaces
- Understand error handling patterns
- Move on to Packages and Modules
