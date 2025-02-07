// Package main demonstrates testing and debugging concepts in Go
package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Calculator represents a basic calculator with memory
type Calculator struct {
	memory float64
	mu     sync.RWMutex
}

// Add adds a number to memory and returns the result
func (c *Calculator) Add(x float64) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.memory += x
	return c.memory
}

// Subtract subtracts a number from memory and returns the result
func (c *Calculator) Subtract(x float64) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.memory -= x
	return c.memory
}

// Multiply multiplies memory by a number and returns the result
func (c *Calculator) Multiply(x float64) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.memory *= x
	return c.memory
}

// Divide divides memory by a number and returns the result
func (c *Calculator) Divide(x float64) (float64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if x == 0 {
		return 0, errors.New("division by zero")
	}
	c.memory /= x
	return c.memory, nil
}

// GetMemory returns the current value in memory
func (c *Calculator) GetMemory() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.memory
}

// Clear resets the memory to zero
func (c *Calculator) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.memory = 0
}

// StringProcessor provides string manipulation operations
type StringProcessor struct {
	cache map[string]string
	mu    sync.RWMutex
}

// NewStringProcessor creates a new StringProcessor
func NewStringProcessor() *StringProcessor {
	return &StringProcessor{
		cache: make(map[string]string),
	}
}

// Reverse reverses a string
func (sp *StringProcessor) Reverse(s string) string {
	// Check cache first
	sp.mu.RLock()
	if cached, ok := sp.cache[s]; ok {
		sp.mu.RUnlock()
		return cached
	}
	sp.mu.RUnlock()

	// Process if not in cache
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	result := string(runes)

	// Cache the result
	sp.mu.Lock()
	sp.cache[s] = result
	sp.mu.Unlock()

	return result
}

// IsPalindrome checks if a string is a palindrome
func (sp *StringProcessor) IsPalindrome(s string) bool {
	s = strings.ToLower(strings.ReplaceAll(s, " ", ""))
	return s == sp.Reverse(s)
}

// Counter represents a thread-safe counter
type Counter struct {
	value int64
	mu    sync.Mutex
}

// Increment increases the counter by one
func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// GetValue returns the current counter value
func (c *Counter) GetValue() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// DataProcessor handles data processing operations
type DataProcessor struct {
	data []int
	mu   sync.RWMutex
}

// NewDataProcessor creates a new DataProcessor
func NewDataProcessor(size int) *DataProcessor {
	return &DataProcessor{
		data: make([]int, size),
	}
}

// Process simulates data processing with artificial delay
func (dp *DataProcessor) Process(value int) error {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	if value < 0 {
		return errors.New("negative values not allowed")
	}

	// Simulate processing time
	time.Sleep(time.Millisecond * time.Duration(value%10))

	for i := range dp.data {
		dp.data[i] = value
	}

	return nil
}

// GetSum returns the sum of all processed data
func (dp *DataProcessor) GetSum() int {
	dp.mu.RLock()
	defer dp.mu.RUnlock()

	sum := 0
	for _, v := range dp.data {
		sum += v
	}
	return sum
}

// Fibonacci calculates the nth Fibonacci number
func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// ConcatStrings concatenates strings in different ways
func ConcatStrings(strs []string, method string) string {
	switch method {
	case "plus":
		result := ""
		for _, s := range strs {
			result += s
		}
		return result

	case "builder":
		var builder strings.Builder
		for _, s := range strs {
			builder.WriteString(s)
		}
		return builder.String()

	default:
		return strings.Join(strs, "")
	}
}

func main() {
	// Example usage of Calculator
	calc := &Calculator{}
	calc.Add(5)
	calc.Multiply(2)
	result, _ := calc.Divide(2)
	fmt.Printf("Calculator result: %.2f\n", result)

	// Example usage of StringProcessor
	sp := NewStringProcessor()
	text := "A man a plan a canal Panama"
	fmt.Printf("'%s' is palindrome: %v\n", text, sp.IsPalindrome(text))

	// Example usage of Counter
	counter := &Counter{}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	fmt.Printf("Counter value: %d\n", counter.GetValue())

	// Example usage of DataProcessor
	dp := NewDataProcessor(5)
	dp.Process(10)
	fmt.Printf("Data sum: %d\n", dp.GetSum())

	// Example of Fibonacci
	n := 10
	fmt.Printf("Fibonacci(%d) = %d\n", n, Fibonacci(n))

	// Example of string concatenation
	strs := []string{"Hello", ", ", "World", "!"}
	fmt.Printf("Concatenated string: %s\n", ConcatStrings(strs, "builder"))
}
