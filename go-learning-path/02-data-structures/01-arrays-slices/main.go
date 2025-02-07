// Package main demonstrates array and slice operations in Go
package main

import (
	"fmt"
	"sort"
)

// DataProcessor demonstrates various array and slice operations
type DataProcessor struct {
	data []int
}

// NewDataProcessor creates a new DataProcessor with initial data
func NewDataProcessor(data []int) *DataProcessor {
	// Create a copy of the input data
	copiedData := make([]int, len(data))
	copy(copiedData, data)
	
	return &DataProcessor{
		data: copiedData,
	}
}

// Sort sorts the data in ascending order
func (dp *DataProcessor) Sort() {
	sort.Ints(dp.data)
}

// Reverse reverses the order of elements
func (dp *DataProcessor) Reverse() {
	for i := 0; i < len(dp.data)/2; i++ {
		j := len(dp.data) - 1 - i
		dp.data[i], dp.data[j] = dp.data[j], dp.data[i]
	}
}

// Filter returns elements that satisfy the predicate
func (dp *DataProcessor) Filter(predicate func(int) bool) []int {
	result := make([]int, 0, len(dp.data)) // Preallocate capacity
	for _, value := range dp.data {
		if predicate(value) {
			result = append(result, value)
		}
	}
	return result
}

// Map applies a transformation to each element
func (dp *DataProcessor) Map(transform func(int) int) []int {
	result := make([]int, len(dp.data))
	for i, value := range dp.data {
		result[i] = transform(value)
	}
	return result
}

// RemoveDuplicates removes duplicate elements while preserving order
func (dp *DataProcessor) RemoveDuplicates() {
	if len(dp.data) == 0 {
		return
	}

	// Use a map to track seen values
	seen := make(map[int]bool)
	result := make([]int, 0, len(dp.data))

	for _, value := range dp.data {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}

	dp.data = result
}

// SlidingWindow returns all possible windows of the specified size
func (dp *DataProcessor) SlidingWindow(windowSize int) [][]int {
	if windowSize <= 0 || windowSize > len(dp.data) {
		return nil
	}

	windows := make([][]int, 0, len(dp.data)-windowSize+1)
	
	for i := 0; i <= len(dp.data)-windowSize; i++ {
		// Create a new window with copied data
		window := make([]int, windowSize)
		copy(window, dp.data[i:i+windowSize])
		windows = append(windows, window)
	}

	return windows
}

// Chunk splits the data into chunks of the specified size
func (dp *DataProcessor) Chunk(size int) [][]int {
	if size <= 0 {
		return nil
	}

	chunks := make([][]int, 0, (len(dp.data)+size-1)/size)
	
	for i := 0; i < len(dp.data); i += size {
		end := i + size
		if end > len(dp.data) {
			end = len(dp.data)
		}
		
		chunk := make([]int, end-i)
		copy(chunk, dp.data[i:end])
		chunks = append(chunks, chunk)
	}

	return chunks
}

// Rotate rotates the elements by the specified positions
func (dp *DataProcessor) Rotate(positions int) {
	if len(dp.data) == 0 {
		return
	}

	// Normalize positions
	positions = positions % len(dp.data)
	if positions < 0 {
		positions += len(dp.data)
	}

	// Create a temporary slice with copied data
	temp := make([]int, len(dp.data))
	copy(temp, dp.data)

	// Perform rotation
	for i := 0; i < len(dp.data); i++ {
		newPos := (i + positions) % len(dp.data)
		dp.data[newPos] = temp[i]
	}
}

// String returns a string representation of the data
func (dp *DataProcessor) String() string {
	return fmt.Sprintf("%v", dp.data)
}

func main() {
	// Create sample data
	data := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	processor := NewDataProcessor(data)

	fmt.Printf("Original data: %v\n", processor)

	// Demonstrate sorting
	processor.Sort()
	fmt.Printf("Sorted: %v\n", processor)

	// Demonstrate reversing
	processor.Reverse()
	fmt.Printf("Reversed: %v\n", processor)

	// Demonstrate filtering (even numbers)
	evens := processor.Filter(func(x int) bool {
		return x%2 == 0
	})
	fmt.Printf("Even numbers: %v\n", evens)

	// Demonstrate mapping (double each number)
	doubled := processor.Map(func(x int) int {
		return x * 2
	})
	fmt.Printf("Doubled: %v\n", doubled)

	// Demonstrate duplicate removal
	processor.RemoveDuplicates()
	fmt.Printf("Without duplicates: %v\n", processor)

	// Demonstrate sliding window
	windows := processor.SlidingWindow(3)
	fmt.Printf("Sliding windows (size 3): %v\n", windows)

	// Demonstrate chunking
	chunks := processor.Chunk(2)
	fmt.Printf("Chunks (size 2): %v\n", chunks)

	// Demonstrate rotation
	processor.Rotate(2)
	fmt.Printf("Rotated by 2 positions: %v\n", processor)
}
