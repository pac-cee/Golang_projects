# Arrays and Slices in Go 📚

## 📖 Table of Contents
1. [Arrays](#arrays)
2. [Slices](#slices)
3. [Operations](#operations)
4. [Memory Management](#memory-management)
5. [Best Practices](#best-practices)
6. [Exercises](#exercises)

## Arrays

### Array Declaration
```go
// Fixed-size array
var numbers [5]int

// Array with initialization
scores := [3]int{90, 95, 100}

// Array with size inference
names := [...]string{"Alice", "Bob", "Charlie"}
```

### Array Properties
- Fixed length
- Zero-based indexing
- Contiguous memory
- Pass by value (copying)

### Multi-dimensional Arrays
```go
// 2D array
var matrix [3][3]int

// Initialize 2D array
grid := [2][3]int{
    {1, 2, 3},
    {4, 5, 6},
}
```

## Slices

### Slice Declaration
```go
// Empty slice
var numbers []int

// Slice with initialization
scores := []int{90, 95, 100}

// Make with length and capacity
buffer := make([]byte, 5, 10)
```

### Slice Properties
- Dynamic length
- Reference type
- Built on arrays
- Three components:
  - Pointer to array
  - Length
  - Capacity

### Slice Operations

#### Append
```go
numbers := []int{1, 2, 3}
numbers = append(numbers, 4)
numbers = append(numbers, 5, 6, 7)
```

#### Slicing
```go
// Slice syntax: slice[startIndex:endIndex]
numbers := []int{0, 1, 2, 3, 4, 5}
subset := numbers[2:4]  // [2, 3]
```

#### Copy
```go
src := []int{1, 2, 3}
dst := make([]int, len(src))
copy(dst, src)
```

## Operations

### Iteration
```go
// Using range
for index, value := range slice {
    fmt.Printf("Index: %d, Value: %d\n", index, value)
}

// Using traditional for loop
for i := 0; i < len(slice); i++ {
    fmt.Printf("Value at %d: %d\n", i, slice[i])
}
```

### Sorting
```go
// Sort integers
numbers := []int{3, 1, 4, 1, 5, 9}
sort.Ints(numbers)

// Sort strings
names := []string{"Charlie", "Alice", "Bob"}
sort.Strings(names)

// Custom sort
sort.Slice(items, func(i, j int) bool {
    return items[i].Value < items[j].Value
})
```

### Filtering
```go
// Filter even numbers
numbers := []int{1, 2, 3, 4, 5, 6}
evens := make([]int, 0)
for _, n := range numbers {
    if n%2 == 0 {
        evens = append(evens, n)
    }
}
```

### Mapping
```go
// Double all numbers
numbers := []int{1, 2, 3, 4, 5}
doubled := make([]int, len(numbers))
for i, n := range numbers {
    doubled[i] = n * 2
}
```

## Memory Management

### Capacity Growth
```go
// Preallocate with expected capacity
numbers := make([]int, 0, 1000)

// Append efficiently
for i := 0; i < 1000; i++ {
    numbers = append(numbers, i)
}
```

### Memory Leaks
```go
// Potential memory leak
original := []int{1, 2, 3, 4, 5}
subset := original[2:4]  // Keeps reference to original

// Better approach
subset = append([]int(nil), original[2:4]...)
```

### Slice Internals
```go
// Understanding slice header
type sliceHeader struct {
    Data uintptr
    Len  int
    Cap  int
}
```

## Best Practices

### 1. Initialization
```go
// Prefer make for slices with known size
slice := make([]int, 0, expectedSize)

// Use literal for small, known values
small := []int{1, 2, 3}
```

### 2. Capacity Management
```go
// Preallocate when size is known
users := make([]User, 0, len(records))

// Double capacity when unknown
if cap(slice) == len(slice) {
    newSlice := make([]int, len(slice), 2*cap(slice))
    copy(newSlice, slice)
    slice = newSlice
}
```

### 3. Slicing
```go
// Full slice expression
slice := array[2:4:4]  // Controls capacity

// Avoid memory leaks
trimmed := make([]byte, len(slice))
copy(trimmed, slice)
```

### 4. Append
```go
// Efficient append
var buffer []byte
buffer = append(buffer, data...)

// Append multiple values
slice = append(slice, 1, 2, 3)
```

## Exercises

### Exercise 1: Array Operations
```go
// Implement a function to reverse an array
func reverse(arr [5]int) [5]int {
    var reversed [5]int
    for i := 0; i < len(arr); i++ {
        reversed[i] = arr[len(arr)-1-i]
    }
    return reversed
}
```

### Exercise 2: Slice Operations
```go
// Implement a function to remove duplicates
func removeDuplicates(slice []int) []int {
    seen := make(map[int]bool)
    result := make([]int, 0)
    
    for _, value := range slice {
        if !seen[value] {
            seen[value] = true
            result = append(result, value)
        }
    }
    return result
}
```

### Exercise 3: Advanced Operations
```go
// Implement a sliding window
func slidingWindow(slice []int, windowSize int) [][]int {
    var windows [][]int
    for i := 0; i <= len(slice)-windowSize; i++ {
        window := make([]int, windowSize)
        copy(window, slice[i:i+windowSize])
        windows = append(windows, window)
    }
    return windows
}
```

## Common Patterns

### 1. Stack Implementation
```go
type Stack struct {
    items []interface{}
}

func (s *Stack) Push(item interface{}) {
    s.items = append(s.items, item)
}

func (s *Stack) Pop() interface{} {
    if len(s.items) == 0 {
        return nil
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item
}
```

### 2. Queue Implementation
```go
type Queue struct {
    items []interface{}
}

func (q *Queue) Enqueue(item interface{}) {
    q.items = append(q.items, item)
}

func (q *Queue) Dequeue() interface{} {
    if len(q.items) == 0 {
        return nil
    }
    item := q.items[0]
    q.items = q.items[1:]
    return item
}
```

### 3. Ring Buffer
```go
type RingBuffer struct {
    data     []interface{}
    size     int
    readPos  int
    writePos int
}

func NewRingBuffer(size int) *RingBuffer {
    return &RingBuffer{
        data: make([]interface{}, size),
        size: size,
    }
}
```

## Next Steps
- Practice array and slice operations
- Implement common data structures
- Study memory management
- Move on to Maps and Sets
