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

// SECTION 1: Basic Concepts
// ------------------------

// Custom type definitions
type Person struct {
	Name string
	Age  int
}

// Interface definition
type Speaker interface {
	Speak() string
}

// Method implementation
func (p Person) Speak() string {
	return fmt.Sprintf("Hello, I'm %s and I'm %d years old", p.Name, p.Age)
}

// SECTION 2: Variables and Data Types
// ---------------------------------

func basicTypes() {
	// Basic variable declarations
	var name string = "John"
	age := 30 // Short declaration

	// Constants
	const Pi = 3.14159

	// Multiple variable declaration
	var (
		isActive bool   = true
		count    int    = 42
		message  string = "Hello"
	)

	// Array and Slice
	numbers := []int{1, 2, 3, 4, 5}
	
	// Map
	person := map[string]string{
		"name": "Alice",
		"city": "New York",
	}

	fmt.Printf("Basic types example: %v, %v, %v, %v, %v, %v\n",
		name, age, Pi, isActive, numbers, person)
}

// SECTION 3: Control Flow
// ----------------------

func controlFlow() {
	// If statement
	x := 10
	if x > 5 {
		fmt.Println("x is greater than 5")
	} else if x < 5 {
		fmt.Println("x is less than 5")
	} else {
		fmt.Println("x equals 5")
	}

	// Switch statement
	switch day := time.Now().Weekday(); day {
	case time.Saturday, time.Sunday:
		fmt.Println("It's the weekend")
	default:
		fmt.Println("It's a weekday")
	}

	// For loop
	for i := 0; i < 5; i++ {
		fmt.Printf("Iteration %d\n", i)
	}

	// Range loop
	fruits := []string{"apple", "banana", "orange"}
	for index, value := range fruits {
		fmt.Printf("Index: %d, Value: %s\n", index, value)
	}
}

// SECTION 4: Functions
// -------------------

// Function with multiple return values
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}

// Variadic function
func sum(numbers ...int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total
}

// SECTION 5: Error Handling
// ------------------------

type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("error %d: %s", e.Code, e.Message)
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
	err = &CustomError{
		Code:    404,
		Message: "Not Found",
	}
	fmt.Printf("Custom error: %v\n", err)
}

// SECTION 6: Concurrency
// ---------------------

func concurrencyExample() {
	// Channels
	ch := make(chan string)
	done := make(chan bool)

	// Goroutine with channel
	go func() {
		for i := 0; i < 5; i++ {
			ch <- fmt.Sprintf("message %d", i)
			time.Sleep(time.Millisecond * 500)
		}
		close(ch)
	}()

	// Receive from channel
	go func() {
		for msg := range ch {
			fmt.Printf("Received: %s\n", msg)
		}
		done <- true
	}()

	<-done // Wait for completion
}

// Worker Pool pattern
func workerPool(jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		results <- job * 2 // Simulate work
		time.Sleep(time.Millisecond * 100)
	}
}

// SECTION 7: HTTP Server
// ---------------------

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
	log.Fatal(http.ListenAndServe(":8080", s.router))
}

// SECTION 8: Context Usage
// ----------------------

func contextExample() {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Simulate work with context
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("Work completed")
	case <-ctx.Done():
		fmt.Println("Work cancelled:", ctx.Err())
	}
}

// SECTION 9: Database Operations
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

// SECTION 10: File Operations
// -------------------------

func fileOperations() error {
	// Write to file
	data := []byte("Hello, Go!")
	if err := os.WriteFile("example.txt", data, 0644); err != nil {
		return err
	}

	// Read from file
	content, err := os.ReadFile("example.txt")
	if err != nil {
		return err
	}
	fmt.Printf("File content: %s\n", content)
	return nil
}

// Main function demonstrating usage
func main() {
	fmt.Println("Go Programming Tutorial")
	fmt.Println("----------------------")

	// Basic concepts
	basicTypes()
	controlFlow()

	// Error handling
	errorHandling()

	// Concurrency
	concurrencyExample()

	// Worker pool example
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go workerPool(jobs, results, &wg)
	}

	// Send jobs
	for i := 0; i < numJobs; i++ {
		jobs <- i
	}
	close(jobs)

	// Wait for workers
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for result := range results {
		fmt.Printf("Result: %d\n", result)
	}

	// Context example
	contextExample()

	// File operations
	if err := fileOperations(); err != nil {
		fmt.Printf("File operation error: %v\n", err)
	}

	// Start HTTP server
	server := NewAPIServer()
	fmt.Println("Starting HTTP server on :8080")
	server.Start()
}
