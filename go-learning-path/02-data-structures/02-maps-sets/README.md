# Maps and Sets in Go 🗺️

## 📖 Table of Contents
1. [Maps](#maps)
2. [Sets](#sets)
3. [Operations](#operations)
4. [Concurrency](#concurrency)
5. [Best Practices](#best-practices)
6. [Exercises](#exercises)

## Maps

### Map Declaration
```go
// Empty map
var scores map[string]int

// Map with make
users := make(map[string]User)

// Map literal
capitals := map[string]string{
    "USA": "Washington",
    "UK":  "London",
    "FR":  "Paris",
}
```

### Map Properties
- Key-value pairs
- Unordered collection
- Reference type
- Fast lookups (O(1) average)
- Dynamic size

### Key Types
Valid key types:
- Comparable types
- Numbers
- Strings
- Pointers
- Interfaces
- Structs
- Arrays

Invalid key types:
- Slices
- Maps
- Functions

## Operations

### Basic Operations
```go
// Insert or update
scores["Alice"] = 100

// Retrieve value
score := scores["Bob"]

// Check existence
score, exists := scores["Charlie"]
if !exists {
    fmt.Println("Charlie not found")
}

// Delete
delete(scores, "David")

// Length
size := len(scores)
```

### Iteration
```go
// Range over map
for key, value := range scores {
    fmt.Printf("%s: %d\n", key, value)
}

// Keys only
for key := range scores {
    fmt.Println(key)
}

// Values only
for _, value := range scores {
    fmt.Println(value)
}
```

### Map of Maps
```go
// Nested maps
graph := map[string]map[string]int{
    "A": {"B": 5, "C": 3},
    "B": {"A": 5, "C": 2},
    "C": {"A": 3, "B": 2},
}
```

## Sets

### Set Implementation
```go
// Set using map
type Set map[string]struct{}

// Create new set
func NewSet() Set {
    return make(Set)
}

// Add element
func (s Set) Add(item string) {
    s[item] = struct{}{}
}

// Remove element
func (s Set) Remove(item string) {
    delete(s, item)
}

// Contains element
func (s Set) Contains(item string) bool {
    _, exists := s[item]
    return exists
}
```

### Set Operations
```go
// Union
func (s Set) Union(other Set) Set {
    result := NewSet()
    for item := range s {
        result.Add(item)
    }
    for item := range other {
        result.Add(item)
    }
    return result
}

// Intersection
func (s Set) Intersection(other Set) Set {
    result := NewSet()
    for item := range s {
        if other.Contains(item) {
            result.Add(item)
        }
    }
    return result
}

// Difference
func (s Set) Difference(other Set) Set {
    result := NewSet()
    for item := range s {
        if !other.Contains(item) {
            result.Add(item)
        }
    }
    return result
}
```

## Concurrency

### Sync.Map
```go
// Thread-safe map
var cache sync.Map

// Store
cache.Store("key", value)

// Load
value, ok := cache.Load("key")

// Delete
cache.Delete("key")

// Load or Store
value, loaded := cache.LoadOrStore("key", defaultValue)
```

### Map Mutex
```go
type SafeMap struct {
    sync.RWMutex
    data map[string]interface{}
}

func (m *SafeMap) Set(key string, value interface{}) {
    m.Lock()
    defer m.Unlock()
    m.data[key] = value
}

func (m *SafeMap) Get(key string) (interface{}, bool) {
    m.RLock()
    defer m.RUnlock()
    value, ok := m.data[key]
    return value, ok
}
```

## Best Practices

### 1. Initialization
```go
// Always initialize maps
scores := make(map[string]int)

// Or use literal for known values
scores := map[string]int{
    "Alice": 100,
    "Bob":   95,
}
```

### 2. Zero Values
```go
// Check zero value
value, exists := scores["Alice"]
if !exists {
    // Key doesn't exist
}

// Or use zero value directly
score := scores["Alice"] // 0 if not found
```

### 3. Memory Management
```go
// Clear map
for key := range scores {
    delete(scores, key)
}

// Or reassign
scores = make(map[string]int)
```

### 4. Key Design
```go
// Good key design
type UserKey struct {
    ID      int
    Version int
}

cache := make(map[UserKey]User)
```

## Exercises

### Exercise 1: Frequency Counter
```go
func frequency(text string) map[rune]int {
    freq := make(map[rune]int)
    for _, char := range text {
        freq[char]++
    }
    return freq
}
```

### Exercise 2: Graph Implementation
```go
type Graph struct {
    edges map[string]map[string]int
}

func (g *Graph) AddEdge(from, to string, weight int) {
    if g.edges[from] == nil {
        g.edges[from] = make(map[string]int)
    }
    g.edges[from][to] = weight
}
```

### Exercise 3: Cache Implementation
```go
type Cache struct {
    sync.RWMutex
    data    map[string]interface{}
    expires map[string]time.Time
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
    c.Lock()
    defer c.Unlock()
    c.data[key] = value
    c.expires[key] = time.Now().Add(ttl)
}
```

## Common Patterns

### 1. Default Values
```go
type Config struct {
    defaults map[string]interface{}
    custom   map[string]interface{}
}

func (c *Config) Get(key string) interface{} {
    if value, ok := c.custom[key]; ok {
        return value
    }
    return c.defaults[key]
}
```

### 2. Counting/Grouping
```go
// Count occurrences
counts := make(map[string]int)
for _, item := range items {
    counts[item]++
}

// Group by property
groups := make(map[string][]Item)
for _, item := range items {
    key := item.Category
    groups[key] = append(groups[key], item)
}
```

### 3. Memoization
```go
type Memoizer struct {
    cache map[string]interface{}
}

func (m *Memoizer) Compute(key string, fn func() interface{}) interface{} {
    if value, ok := m.cache[key]; ok {
        return value
    }
    value := fn()
    m.cache[key] = value
    return value
}
```

## Next Steps
- Practice map operations
- Implement custom set types
- Study concurrent map usage
- Move on to Trees and Graphs
