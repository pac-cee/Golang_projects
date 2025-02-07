# Testing and Debugging in Go 🧪

## 📖 Table of Contents
1. [Unit Testing](#unit-testing)
2. [Table-Driven Tests](#table-driven-tests)
3. [Benchmarking](#benchmarking)
4. [Profiling](#profiling)
5. [Race Detection](#race-detection)
6. [Debugging](#debugging)
7. [Best Practices](#best-practices)
8. [Exercises](#exercises)

## Unit Testing

### Test File Structure
```go
package mypackage_test  // Use separate package for black-box testing

import (
    "testing"
    "github.com/your/package"
)

func TestMyFunction(t *testing.T) {
    // Test implementation
}
```

### Basic Test Function
```go
func TestAdd(t *testing.T) {
    got := Add(2, 3)
    want := 5
    if got != want {
        t.Errorf("Add(2, 3) = %d; want %d", got, want)
    }
}
```

### Test Helper Functions
```go
func TestComplex(t *testing.T) {
    t.Helper()  // Marks this as a helper function
    // Helper implementation
}
```

### Subtests
```go
func TestMath(t *testing.T) {
    t.Run("addition", func(t *testing.T) {
        if got := Add(2, 3); got != 5 {
            t.Errorf("Add(2, 3) = %d; want 5", got)
        }
    })

    t.Run("multiplication", func(t *testing.T) {
        if got := Multiply(2, 3); got != 6 {
            t.Errorf("Multiply(2, 3) = %d; want 6", got)
        }
    })
}
```

## Table-Driven Tests

### Basic Structure
```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        x, y int
        want int
    }{
        {"positive", 2, 3, 5},
        {"negative", -2, -3, -5},
        {"mixed", -2, 3, 1},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.x, tt.y)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.x, tt.y, got, tt.want)
            }
        })
    }
}
```

### Testing Errors
```go
func TestDivide(t *testing.T) {
    tests := []struct {
        name    string
        x, y    int
        want    int
        wantErr bool
    }{
        {"valid", 6, 2, 3, false},
        {"zero divisor", 6, 0, 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.x, tt.y)
            if (err != nil) != tt.wantErr {
                t.Errorf("Divide(%d, %d) error = %v; wantErr %v",
                    tt.x, tt.y, err, tt.wantErr)
                return
            }
            if !tt.wantErr && got != tt.want {
                t.Errorf("Divide(%d, %d) = %d; want %d",
                    tt.x, tt.y, got, tt.want)
            }
        })
    }
}
```

## Benchmarking

### Basic Benchmark
```go
func BenchmarkFibonacci(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Fibonacci(10)
    }
}
```

### Benchmark with Setup
```go
func BenchmarkComplexOperation(b *testing.B) {
    // Setup code here
    data := makeTestData()
    
    b.ResetTimer()  // Reset timer after setup
    for i := 0; i < b.N; i++ {
        ComplexOperation(data)
    }
}
```

### Parallel Benchmark
```go
func BenchmarkParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            // Code to benchmark
        }
    })
}
```

## Profiling

### CPU Profiling
```go
import "runtime/pprof"

func main() {
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // Your program here
}
```

### Memory Profiling
```go
import "runtime/pprof"

func main() {
    f, _ := os.Create("mem.prof")
    defer f.Close()
    
    // Your program here
    
    pprof.WriteHeapProfile(f)
}
```

### Using go test for Profiling
```bash
# CPU profile
go test -cpuprofile=cpu.prof -bench .

# Memory profile
go test -memprofile=mem.prof -bench .

# Block profile
go test -blockprofile=block.prof -bench .
```

## Race Detection

### Running Tests with Race Detector
```bash
go test -race ./...
```

### Common Race Conditions
```go
// Data race
var counter int
go func() { counter++ }()
go func() { counter++ }()

// Fix with mutex
var mu sync.Mutex
var counter int
go func() {
    mu.Lock()
    counter++
    mu.Unlock()
}()
```

## Debugging

### Using Delve Debugger
```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Start debugging
dlv debug

# Common commands
(dlv) break main.go:20
(dlv) continue
(dlv) next
(dlv) step
(dlv) print varName
```

### Logging
```go
import "log"

func Process() {
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    log.Printf("Processing started")
    // ...
    log.Printf("Processing completed")
}
```

## Best Practices

### 1. Test Organization
```go
// Group related tests
func TestUserOperations(t *testing.T) {
    t.Run("group=create", func(t *testing.T) {
        t.Run("valid user", testCreateValidUser)
        t.Run("invalid user", testCreateInvalidUser)
    })
    
    t.Run("group=update", func(t *testing.T) {
        t.Run("valid update", testValidUpdate)
        t.Run("invalid update", testInvalidUpdate)
    })
}
```

### 2. Test Coverage
```bash
# Run tests with coverage
go test -cover

# Generate coverage profile
go test -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out
```

### 3. Test Fixtures
```go
func setupTestCase(t *testing.T) func() {
    // Setup code
    return func() {
        // Teardown code
    }
}

func TestWithFixture(t *testing.T) {
    teardown := setupTestCase(t)
    defer teardown()
    
    // Test code
}
```

## Exercises

### Exercise 1: Unit Testing
```go
// Implement comprehensive tests for a Calculator type
type Calculator struct {
    memory float64
}

func (c *Calculator) Add(x float64) float64 {
    c.memory += x
    return c.memory
}

// Write tests covering:
// - Basic operations
// - Edge cases
// - Error conditions
```

### Exercise 2: Benchmark Optimization
```go
// Optimize and benchmark string operations
func ConcatStrings(n int) string {
    // Implementation
}

// Write benchmarks comparing:
// - String concatenation
// - StringBuilder
// - bytes.Buffer
```

### Exercise 3: Race Detection
```go
// Fix race conditions in concurrent counter
type Counter struct {
    value int
}

func (c *Counter) Increment() {
    c.value++
}

// Implement thread-safe version and test with -race
```

## Common Patterns

### 1. Test Setup
```go
type testCase struct {
    input    string
    want     string
    wantErr  bool
}

func runTestCases(t *testing.T, fn func(string) (string, error), cases []testCase) {
    for _, tc := range cases {
        got, err := fn(tc.input)
        if (err != nil) != tc.wantErr {
            t.Errorf("error = %v, wantErr %v", err, tc.wantErr)
            continue
        }
        if got != tc.want {
            t.Errorf("got %v, want %v", got, tc.want)
        }
    }
}
```

### 2. Mocking
```go
type DataStore interface {
    Get(key string) (string, error)
    Set(key, value string) error
}

type MockDataStore struct {
    data map[string]string
}

func (m *MockDataStore) Get(key string) (string, error) {
    if val, ok := m.data[key]; ok {
        return val, nil
    }
    return "", fmt.Errorf("not found")
}
```

### 3. Benchmarking Patterns
```go
func BenchmarkWithDifferentSizes(b *testing.B) {
    sizes := []int{100, 1000, 10000}
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                ProcessData(size)
            }
        })
    }
}
```

## Next Steps
- Practice writing comprehensive test suites
- Learn advanced debugging techniques
- Study performance optimization
- Move on to Web Development
