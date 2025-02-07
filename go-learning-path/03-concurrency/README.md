# Concurrency in Go 🔄

## 📖 Table of Contents
1. [Goroutines](#goroutines)
2. [Channels](#channels)
3. [Select Statement](#select-statement)
4. [Sync Package](#sync-package)
5. [Concurrency Patterns](#concurrency-patterns)
6. [Best Practices](#best-practices)
7. [Exercises](#exercises)

## Goroutines

### What are Goroutines?
- Lightweight threads managed by Go runtime
- Much smaller stack size than OS threads
- Multiplexed onto OS threads
- Non-blocking by default

### Basic Usage
```go
func main() {
    go myFunction()    // Start a new goroutine
    go func() {        // Anonymous function as goroutine
        // Do something
    }()
}
```

### Key Points
- Goroutines are cheap and lightweight
- Go runtime handles scheduling
- No direct communication between goroutines (use channels)
- Main function runs in the main goroutine

## Channels

### Channel Basics
```go
ch := make(chan int)        // Unbuffered channel
ch := make(chan int, 100)   // Buffered channel with capacity 100

// Send and receive
ch <- value    // Send value
value := <-ch  // Receive value
```

### Channel Types
1. **Unbuffered Channels**
   - Synchronous communication
   - Sender blocks until receiver is ready
   - Receiver blocks until sender sends

2. **Buffered Channels**
   - Asynchronous up to buffer size
   - Sender blocks only when buffer is full
   - Receiver blocks only when buffer is empty

3. **Directional Channels**
```go
chan<- int   // Send-only channel
<-chan int   // Receive-only channel
```

### Channel Operations
```go
// Close channel
close(ch)

// Check if channel is closed
value, ok := <-ch
if !ok {
    // Channel is closed
}

// Range over channel
for value := range ch {
    // Process value
}
```

## Select Statement

### Basic Select
```go
select {
case v1 := <-ch1:
    // Handle value from ch1
case v2 := <-ch2:
    // Handle value from ch2
case ch3 <- x:
    // Send x to ch3
default:
    // Optional default case
}
```

### Common Patterns
1. **Timeout**
```go
select {
case res := <-ch:
    // Handle result
case <-time.After(1 * time.Second):
    // Handle timeout
}
```

2. **Non-blocking Operations**
```go
select {
case ch <- x:
    // Sent successfully
default:
    // Would block
}
```

## Sync Package

### Mutex
```go
var mu sync.Mutex
mu.Lock()
// Critical section
mu.Unlock()
```

### RWMutex
```go
var mu sync.RWMutex
mu.RLock()    // Multiple readers
// Read operations
mu.RUnlock()

mu.Lock()     // Single writer
// Write operations
mu.Unlock()
```

### WaitGroup
```go
var wg sync.WaitGroup
wg.Add(n)     // Add n goroutines
go func() {
    defer wg.Done()
    // Work
}()
wg.Wait()     // Wait for all to finish
```

### Once
```go
var once sync.Once
once.Do(func() {
    // Execute only once
})
```

## Concurrency Patterns

### Worker Pool
```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        results <- process(j)
    }
}

// Create worker pool
for w := 1; w <= 3; w++ {
    go worker(w, jobs, results)
}
```

### Pipeline
```go
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

func sq(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}
```

### Fan-out, Fan-in
```go
func fanOut(ch <-chan int, n int) []<-chan int {
    channels := make([]<-chan int, n)
    for i := 0; i < n; i++ {
        channels[i] = process(ch)
    }
    return channels
}

func fanIn(channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)
    
    output := func(c <-chan int) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }
    
    wg.Add(len(channels))
    for _, c := range channels {
        go output(c)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

## Best Practices

### 1. Error Handling
```go
// Use error channels for error handling
type Result struct {
    Value int
    Err   error
}

// Return both result and error
resultCh := make(chan Result)
go func() {
    value, err := riskyOperation()
    resultCh <- Result{value, err}
}()
```

### 2. Context Usage
```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            // Do work
        }
    }
}
```

### 3. Channel Sizing
- Use unbuffered channels for synchronization
- Use buffered channels when you know the exact capacity needed
- Avoid using channels just for synchronization when sync.WaitGroup is more appropriate

### 4. Goroutine Lifecycle
- Always ensure goroutines can exit
- Use context for cancellation
- Clean up resources properly

## Exercises

### Exercise 1: Concurrent Counter
```go
// Implement a thread-safe counter using channels
type Counter struct {
    value chan int
}

func NewCounter() *Counter {
    c := &Counter{value: make(chan int)}
    go c.run()
    return c
}
```

### Exercise 2: Rate Limiter
```go
// Implement a rate limiter using time.Ticker
func RateLimiter(requests <-chan int, limit time.Duration) <-chan int {
    // Implementation
}
```

### Exercise 3: Parallel Processing
```go
// Process items in parallel with limited concurrency
func ProcessParallel(items []int, concurrency int) []Result {
    // Implementation
}
```

## Common Pitfalls

### 1. Race Conditions
- Always use proper synchronization
- Run tests with -race flag
- Use sync package appropriately

### 2. Goroutine Leaks
- Ensure proper cleanup
- Use context for cancellation
- Close channels when done

### 3. Deadlocks
- Avoid circular dependencies
- Use select with default case
- Implement proper timeout mechanisms

## Next Steps
- Practice implementing concurrent algorithms
- Study real-world concurrency patterns
- Learn about performance optimization
- Move on to Testing and Debugging
