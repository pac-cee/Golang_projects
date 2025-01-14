package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// SECTION 1: Basic Types and Structs
// ---------------------------------

// Person represents a basic struct with methods
type Person struct {
	Name string
	Age  int
}

// Speak demonstrates a method implementation
func (p Person) Speak() string {
	return fmt.Sprintf("Hello, I'm %s and I'm %d years old", p.Name, p.Age)
}

// SECTION 2: Interfaces
// -------------------

// Speaker defines behavior for types that can speak
type Speaker interface {
	Speak() string
}

// SECTION 3: Error Handling
// -----------------------

// CustomError demonstrates custom error types
type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("error %d: %s", e.Code, e.Message)
}

// SECTION 4: Basic Functions
// ------------------------

func basicTypes() {
	// Variable declarations
	var name string = "John"
	age := 30 // Short declaration

	// Constants
	const Pi = 3.14159

	// Multiple variables
	var (
		isActive bool   = true
		count    int    = 42
		message  string = "Hello, Go!"
	)

	// Arrays and Slices
	numbers := []int{1, 2, 3, 4, 5}
	
	// Maps
	person := map[string]string{
		"name": "Alice",
		"city": "New York",
	}

	fmt.Printf("Basic types example: %v, %v, %v, %v, %v, %v, %v\n",
		name, age, Pi, isActive, count, message, numbers)
	fmt.Printf("Map example: %v\n", person)
}

// SECTION 5: Control Flow
// ---------------------

func controlFlow() {
	// If-else
	x := 10
	if x > 5 {
		fmt.Println("x is greater than 5")
	} else {
		fmt.Println("x is less than or equal to 5")
	}

	// Switch
	day := time.Now().Weekday()
	switch day {
	case time.Saturday, time.Sunday:
		fmt.Println("It's the weekend!")
	default:
		fmt.Println("It's a weekday")
	}

	// Loops
	for i := 0; i < 3; i++ {
		fmt.Printf("Loop %d\n", i)
	}

	// Range loop
	fruits := []string{"apple", "banana", "orange"}
	for index, value := range fruits {
		fmt.Printf("Index: %d, Value: %s\n", index, value)
	}
}

// SECTION 6: Error Handling Functions
// --------------------------------

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, &CustomError{
			Code:    400,
			Message: "division by zero",
		}
	}
	return a / b, nil
}

func errorHandling() {
	// Basic error handling
	result, err := divide(10, 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %f\n", result)
	}

	// Custom error
	customErr := &CustomError{
		Code:    404,
		Message: "Not Found",
	}
	fmt.Printf("Custom error: %v\n", customErr)
}

// SECTION 7: Concurrency
// --------------------

func concurrencyExample() {
	// Basic channel communication
	ch := make(chan string)
	done := make(chan bool)

	// Goroutine with channel
	go func() {
		for i := 0; i < 3; i++ {
			ch <- fmt.Sprintf("Message %d", i)
			time.Sleep(time.Millisecond * 500)
		}
		close(ch)
	}()

	// Receiving goroutine
	go func() {
		for msg := range ch {
			fmt.Println("Received:", msg)
		}
		done <- true
	}()

	<-done
}

// Worker Pool pattern
func workerPool(jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		results <- job * 2
		time.Sleep(time.Millisecond * 100)
	}
}

// SECTION 8: HTTP Server
// --------------------

type APIServer struct {
	router *http.ServeMux
}

func NewAPIServer() *APIServer {
	return &APIServer{
		router: http.NewServeMux(),
	}
}

func (s *APIServer) handleHello(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Hello, World!"}
	json.NewEncoder(w).Encode(response)
}

func (s *APIServer) Start() {
	s.router.HandleFunc("/hello", s.handleHello)
	log.Printf("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", s.router))
}

// SECTION 9: Context Usage
// ----------------------

func contextExample() {
	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	select {
	case <-time.After(3 * time.Second):
		fmt.Println("Work completed")
	case <-ctx.Done():
		fmt.Println("Work cancelled:", ctx.Err())
	}
}

// SECTION 10: Database Operations
// ----------------------------

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user Person) error {
	query := `INSERT INTO users (name, age) VALUES (?, ?)`
	_, err := r.db.Exec(query, user.Name, user.Age)
	return err
}

// SECTION 11: File Operations
// -------------------------

func fileOperations() error {
	// Writing to file
	data := []byte("Hello, Go!")
	if err := os.WriteFile("example.txt", data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Reading from file
	content, err := os.ReadFile("example.txt")
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	fmt.Printf("File content: %s\n", content)
	return nil
}

// Main function demonstrating all features
func main() {
	fmt.Println("\n=== Basic Types ===")
	basicTypes()

	fmt.Println("\n=== Control Flow ===")
	controlFlow()

	fmt.Println("\n=== Error Handling ===")
	errorHandling()

	fmt.Println("\n=== Concurrency ===")
	concurrencyExample()

	fmt.Println("\n=== Context Example ===")
	contextExample()

	fmt.Println("\n=== File Operations ===")
	if err := fileOperations(); err != nil {
		fmt.Printf("File operation error: %v\n", err)
	}

	// Create and use a Person
	person := Person{Name: "Alice", Age: 30}
	fmt.Println("\n=== Person Example ===")
	fmt.Println(person.Speak())

	// Demonstrate interface usage
	var speaker Speaker = person
	fmt.Println(speaker.Speak())

	fmt.Println("\n=== Server Example ===")
	fmt.Println("Uncomment the following line to start the HTTP server:")
	// server := NewAPIServer()
	// server.Start()
}
