# Go Fundamentals

This section covers the essential concepts of Go programming language.

## Topics Covered
1. Basic Syntax
2. Data Types
3. Control Structures
4. Functions
5. Structs and Interfaces
6. Goroutines and Channels
7. Error Handling
8. Packages and Modules

## Code Explanation

The `main.go` file demonstrates several key Go concepts:

### 1. Structs and Interfaces
```go
type User struct {
    ID       int
    Username string
    Email    string
}

type UserInterface interface {
    GetEmail() string
    UpdateEmail(newEmail string)
}
```
- The `User` struct represents a basic user data structure with ID, username, and email fields
- `UserInterface` defines a contract for user-related operations
- This demonstrates Go's approach to interface implementation (implicit rather than explicit)

### 2. Method Receivers
```go
func (u *User) GetEmail() string
func (u *User) UpdateEmail(newEmail string)
```
- Methods are defined with pointer receivers (`*User`) to allow modification of the struct
- This implements the `UserInterface`, showing Go's interface satisfaction

### 3. Concurrent Programming
```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}
```
- `SafeCounter` demonstrates thread-safe operations using mutexes
- `sync.Mutex` ensures only one goroutine can access the counter at a time
- Methods `Increment()` and `GetCount()` show proper mutex usage

### 4. Error Handling
```go
func processUser(u *User) error {
    if u == nil {
        return fmt.Errorf("user cannot be nil")
    }
    // ...
}
```
- Shows Go's idiomatic error handling using the `error` interface
- Demonstrates error creation with `fmt.Errorf`
- Multiple error conditions are checked and returned

### 5. Goroutines and Channels
```go
func demonstrateChannel(done chan bool)
```
- Shows how to use channels for goroutine communication
- Demonstrates synchronization between concurrent operations
- Uses the channel to signal completion

### 6. Main Function Flow
The `main` function showcases:
- Struct initialization
- Interface usage
- Error handling patterns
- Goroutine execution
- Channel-based communication
- Concurrent counter operations with WaitGroups

## Best Practices Demonstrated
1. Proper error handling with meaningful messages
2. Thread-safe concurrent operations
3. Clean interface design
4. Effective use of Go's concurrency primitives
5. Clear struct organization
6. Proper method receiver usage

## Running the Code
```bash
go run main.go
```
This will demonstrate:
- Interface implementation
- Concurrent operations
- Channel communication
- Thread-safe counter operations
