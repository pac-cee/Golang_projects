package main

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// BASIC CONCEPTS
// 1. Variables and Data Types
func basicTypes() {
	// Integer types
	var i int = 42
	var i8 int8 = 127
	var ui uint = 123

	// Float types
	var f32 float32 = 3.14
	var f64 float64 = 3.141592653589793

	// Boolean
	var isTrue bool = true

	// String
	var str string = "Hello, Go!"

	// Type inference
	auto := "Type inferred automatically"

	fmt.Printf("Types: %T %T %T %T %T %T %T %T\n", i, i8, ui, f32, f64, isTrue, str, auto)
}

// 2. Constants
const (
	Pi       = 3.14159
	Username = "admin"
	// iota example
	Sunday = iota
	Monday
	Tuesday
)

// 3. Arrays and Slices
func arraysAndSlices() {
	// Array
	var arr [5]int = [5]int{1, 2, 3, 4, 5}

	// Slice
	slice := []int{1, 2, 3}
	slice = append(slice, 4)

	// Slice operations
	subSlice := slice[1:3]
	
	fmt.Println(arr, slice, subSlice)
}

// 4. Maps
func maps() {
	// Declaration and initialization
	scores := map[string]int{
		"Alice": 95,
		"Bob":   89,
	}

	// Adding and accessing
	scores["Charlie"] = 85
	fmt.Println(scores["Alice"])

	// Check if key exists
	if score, exists := scores["David"]; exists {
		fmt.Println(score)
	}
}

// INTERMEDIATE CONCEPTS
// 1. Structs and Methods
type Person struct {
	Name string
	Age  int
}

// Method
func (p Person) Greet() string {
	return fmt.Sprintf("Hello, my name is %s and I'm %d years old", p.Name, p.Age)
}

// 2. Interfaces
type Animal interface {
	Speak() string
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return fmt.Sprintf("%s says Woof!", d.Name)
}

// 3. Error Handling
type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, &CustomError{Message: "division by zero"}
	}
	return a / b, nil
}

// 4. Goroutines and Channels
func goroutineExample() {
	ch := make(chan string)
	
	go func() {
		ch <- "Message from goroutine"
	}()

	msg := <-ch
	fmt.Println(msg)
}

// ADVANCED CONCEPTS
// 1. Reflection
func reflectionExample() {
	p := Person{Name: "John", Age: 30}
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)

	fmt.Printf("Type: %v\n", t)
	fmt.Printf("Fields: %v\n", v.Field(0))
}

// 2. Context
func contextExample() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	select {
	case <-time.After(3 * time.Second):
		fmt.Println("Overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}

// 3. Concurrency Patterns
// Worker Pool
func workerPool() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Start workers
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(results)
}

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("worker %d processing job %d\n", id, j)
		time.Sleep(time.Second)
		results <- j * 2
	}
}

// 4. Generics (Go 1.18+)
type Number interface {
	~int | ~float64
}

func Min[T Number](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// Main function demonstrating usage
func main() {
	fmt.Println("=== Basic Concepts ===")
	basicTypes()
	arraysAndSlices()
	maps()

	fmt.Println("\n=== Intermediate Concepts ===")
	p := Person{Name: "Alice", Age: 25}
	fmt.Println(p.Greet())

	d := Dog{Name: "Rex"}
	fmt.Println(d.Speak())

	result, err := divide(10, 2)
	fmt.Printf("Division result: %v, error: %v\n", result, err)

	fmt.Println("\n=== Advanced Concepts ===")
	reflectionExample()
	contextExample()
	workerPool()

	// Generics example
	fmt.Printf("Min of 10 and 5: %v\n", Min(10, 5))
	fmt.Printf("Min of 3.14 and 2.71: %v\n", Min(3.14, 2.71))
}
