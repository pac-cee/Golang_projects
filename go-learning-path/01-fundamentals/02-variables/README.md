# Variables, Types, and Constants in Go 🔄

## 📖 Table of Contents
1. [Basic Types](#basic-types)
2. [Variable Declaration](#variable-declaration)
3. [Type Conversion](#type-conversion)
4. [Constants](#constants)
5. [Type Inference](#type-inference)
6. [Zero Values](#zero-values)
7. [Composite Types](#composite-types)
8. [Exercises](#exercises)

## Basic Types

Go has several basic types:

```go
// Numeric Types
int8, int16, int32, int64
uint8, uint16, uint32, uint64
int  // platform dependent, 32 or 64 bit
uint // platform dependent, 32 or 64 bit
float32, float64
complex64, complex128

// Text Types
string
rune   // alias for int32, represents a Unicode code point
byte   // alias for uint8

// Boolean Type
bool

// Error Type
error
```

### Size and Range Examples
```go
// Integer ranges
int8:   -128 to 127
int16:  -32768 to 32767
int32:  -2147483648 to 2147483647
int64:  -9223372036854775808 to 9223372036854775807

// Unsigned integer ranges
uint8:  0 to 255
uint16: 0 to 65535
uint32: 0 to 4294967295
uint64: 0 to 18446744073709551615
```

## Variable Declaration

There are several ways to declare variables in Go:

### 1. Standard Declaration
```go
var name string
var age int
var isActive bool

// Multiple declarations
var (
    firstName string
    lastName  string
    score     int
)
```

### 2. Declaration with Initial Value
```go
var name string = "Gopher"
var age int = 25
var isActive = true  // Type inference
```

### 3. Short Declaration (Inside Functions)
```go
func example() {
    name := "Gopher"    // string
    age := 25          // int
    score := 95.5      // float64
}
```

### 4. Multiple Assignments
```go
var x, y int = 10, 20
name, age := "Gopher", 25
```

## Type Conversion

Go requires explicit type conversion:

```go
var i int = 42
var f float64 = float64(i)
var u uint = uint(f)

// String conversions
str := strconv.Itoa(i)      // int to string
num, err := strconv.Atoi(str) // string to int

// Float conversions
s := fmt.Sprintf("%.2f", f) // float to string
f, err = strconv.ParseFloat(s, 64) // string to float
```

## Constants

Constants are declared using the `const` keyword:

```go
// Single constant
const Pi = 3.14159

// Multiple constants
const (
    StatusOK    = 200
    StatusError = 500
)

// iota for enumerated constants
const (
    Sunday = iota  // 0
    Monday         // 1
    Tuesday        // 2
    Wednesday      // 3
    Thursday       // 4
    Friday        // 5
    Saturday      // 6
)
```

## Type Inference

Go can infer types based on values:

```go
var x = 42       // int
var f = 3.14     // float64
var s = "hello"  // string
var b = true     // bool

// Short declaration
num := 42        // int
pi := 3.14       // float64
name := "Gopher" // string
```

## Zero Values

Variables declared without an explicit initial value are given their zero value:

```go
var (
    i int     // 0
    f float64 // 0.0
    s string  // ""
    b bool    // false
    p *int    // nil
)
```

## Composite Types

### Arrays
```go
var numbers [5]int          // Array of 5 integers
scores := [3]int{90, 95, 100} // Array with initial values
```

### Slices
```go
var numbers []int           // Slice of integers
scores := []int{90, 95, 100} // Slice with initial values
```

### Maps
```go
var m map[string]int       // Map declaration
scores := map[string]int{  // Map with initial values
    "Alice": 95,
    "Bob":   89,
}
```

## Exercises

### Exercise 1: Variable Declaration
Create variables of different types and print them:

```go
package main

import "fmt"

func main() {
    // Your solution here
    var name string = "Gopher"
    age := 25
    height := 1.75
    isStudent := true

    fmt.Printf("Name: %s\n", name)
    fmt.Printf("Age: %d\n", age)
    fmt.Printf("Height: %.2f\n", height)
    fmt.Printf("Is Student: %v\n", isStudent)
}
```

### Exercise 2: Type Conversion
Convert between different numeric types:

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    // Integer to Float
    x := 42
    f := float64(x)
    fmt.Printf("Integer %d to Float: %.2f\n", x, f)

    // Float to Integer
    y := 3.14
    i := int(y)
    fmt.Printf("Float %.2f to Integer: %d\n", y, i)

    // Number to String
    str := strconv.Itoa(x)
    fmt.Printf("Integer %d to String: %s\n", x, str)
}
```

### Exercise 3: Constants and Iota
Create an enumerated type using constants:

```go
package main

import "fmt"

const (
    January = iota + 1
    February
    March
    April
    May
)

func main() {
    fmt.Printf("January: %d\n", January)
    fmt.Printf("February: %d\n", February)
    fmt.Printf("March: %d\n", March)
}
```

## Best Practices

1. **Variable Naming**
   - Use camelCase for variable names
   - Use short, descriptive names
   - Use i, j, k for loop indices

2. **Type Selection**
   - Use int for integers unless you need a specific size
   - Use float64 for floating-point numbers
   - Use rune for characters

3. **Constants**
   - Use constants for fixed values
   - Use iota for related constants
   - Name constants in MixedCaps or ALL_CAPS

4. **Type Conversion**
   - Always handle errors in string conversions
   - Be aware of potential data loss in numeric conversions
   - Use explicit conversions for clarity

## Common Pitfalls

1. **Unused Variables**
   ```go
   // This will cause a compilation error
   var x int  // declared but not used
   ```

2. **Type Mismatch**
   ```go
   var x int = "hello"  // Type mismatch
   y := 42
   var f float64 = y    // Need explicit conversion
   ```

3. **Constant Overflow**
   ```go
   const x int8 = 1000  // Overflow: constant 1000 overflows int8
   ```

## Next Steps
- Practice type conversions
- Experiment with different numeric types
- Create complex constant declarations
- Move on to Control Structures
