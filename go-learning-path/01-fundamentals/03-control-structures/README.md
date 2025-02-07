# Control Structures in Go 🔄

## 📖 Table of Contents
1. [If Statements](#if-statements)
2. [For Loops](#for-loops)
3. [Switch Statements](#switch-statements)
4. [Defer](#defer)
5. [Panic and Recover](#panic-and-recover)
6. [Best Practices](#best-practices)
7. [Exercises](#exercises)

## If Statements

### Basic If Statement
```go
if condition {
    // code
}
```

### If-Else Statement
```go
if condition {
    // code
} else {
    // code
}
```

### If with Initialization
```go
if value := getValue(); value > 10 {
    // code
}
```

### Multiple Conditions
```go
if condition1 && condition2 {
    // both true
} else if condition1 || condition2 {
    // at least one true
} else {
    // none true
}
```

## For Loops

### Basic For Loop
```go
for i := 0; i < 10; i++ {
    // code
}
```

### While-style Loop
```go
for condition {
    // code
}
```

### Infinite Loop
```go
for {
    // code
    if shouldBreak {
        break
    }
}
```

### Range Loop
```go
// Slice/Array
for index, value := range slice {
    // code
}

// Map
for key, value := range map {
    // code
}

// String
for index, char := range string {
    // code
}

// Channel
for value := range channel {
    // code
}
```

## Switch Statements

### Basic Switch
```go
switch value {
case 1:
    // code
case 2:
    // code
default:
    // code
}
```

### Switch with Initialization
```go
switch os := runtime.GOOS; os {
case "darwin":
    // code
case "linux":
    // code
default:
    // code
}
```

### Switch without Expression
```go
switch {
case condition1:
    // code
case condition2:
    // code
default:
    // code
}
```

### Type Switch
```go
switch v := interface{}.(type) {
case int:
    // code
case string:
    // code
default:
    // code
}
```

## Defer

### Basic Defer
```go
func example() {
    defer fmt.Println("last")
    fmt.Println("first")
}
```

### Multiple Defers
```go
func example() {
    defer fmt.Println("3")
    defer fmt.Println("2")
    defer fmt.Println("1")
    // Prints: 1, 2, 3
}
```

### Defer with Functions
```go
func example() {
    defer func() {
        // cleanup code
    }()
}
```

## Panic and Recover

### Panic
```go
func example() {
    panic("something went wrong")
}
```

### Recover
```go
func example() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered from: %v\n", r)
        }
    }()
    panic("something went wrong")
}
```

## Best Practices

### 1. If Statements
- Keep conditions simple and readable
- Use initialization when possible
- Avoid nested if statements when possible

### 2. For Loops
- Use range when iterating over collections
- Break long loops into smaller functions
- Use continue and break appropriately

### 3. Switch Statements
- Use switch instead of long if-else chains
- Take advantage of case grouping
- Consider type switches for interface handling

### 4. Defer
- Use defer for cleanup operations
- Remember LIFO (Last In, First Out) order
- Don't defer in loops

### 5. Panic and Recover
- Use panic only for unrecoverable errors
- Always recover in deferred functions
- Consider error returns instead of panic

## Exercises

### Exercise 1: Control Flow
```go
// Implement FizzBuzz
func fizzBuzz(n int) {
    for i := 1; i <= n; i++ {
        switch {
        case i%15 == 0:
            fmt.Println("FizzBuzz")
        case i%3 == 0:
            fmt.Println("Fizz")
        case i%5 == 0:
            fmt.Println("Buzz")
        default:
            fmt.Println(i)
        }
    }
}
```

### Exercise 2: Defer and Files
```go
func readFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    // Read file operations
    return nil
}
```

### Exercise 3: Panic Recovery
```go
func safeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("division error: %v", r)
        }
    }()

    if b == 0 {
        panic("division by zero")
    }
    return a / b, nil
}
```

## Common Patterns

### 1. Error Handling
```go
if err != nil {
    return fmt.Errorf("failed to process: %w", err)
}
```

### 2. Resource Cleanup
```go
func processResource() error {
    resource, err := acquireResource()
    if err != nil {
        return err
    }
    defer resource.Release()

    return resource.Process()
}
```

### 3. Graceful Shutdown
```go
func main() {
    defer cleanup()

    // Main program logic
}
```

## Next Steps
- Practice writing clean control structures
- Experiment with different loop patterns
- Understand panic and recovery
- Move on to Functions and Methods
