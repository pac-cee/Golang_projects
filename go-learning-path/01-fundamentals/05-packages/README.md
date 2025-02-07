# Packages and Modules in Go 📦

## 📖 Table of Contents
1. [Package Basics](#package-basics)
2. [Module System](#module-system)
3. [Package Organization](#package-organization)
4. [Import Patterns](#import-patterns)
5. [Dependency Management](#dependency-management)
6. [Best Practices](#best-practices)
7. [Exercises](#exercises)

## Package Basics

### Package Declaration
```go
// Single package declaration at the top of each file
package mypackage

// Main package for executables
package main
```

### Package Naming Conventions
```go
// Good package names
package user
package database
package httputil

// Avoid package names like
package userUtils     // mixed caps
package user_utils    // underscores
```

### Package Documentation
```go
// Package user provides user management functionality.
// It handles user creation, authentication, and authorization.
package user
```

## Module System

### Creating a New Module
```bash
# Initialize a new module
go mod init example.com/myproject

# Add dependencies
go get github.com/some/dependency

# Update dependencies
go get -u ./...

# Clean up dependencies
go mod tidy
```

### Module File (go.mod)
```go
module example.com/myproject

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/go-sql-driver/mysql v1.7.1
)
```

### Workspace File (go.work)
```go
go 1.21

use (
    ./pkg1
    ./pkg2
    ./cmd/server
)
```

## Package Organization

### Standard Layout
```
myproject/
├── cmd/                    # Command applications
│   └── server/
│       └── main.go
├── internal/               # Private packages
│   ├── auth/
│   └── database/
├── pkg/                    # Public packages
│   ├── models/
│   └── utils/
├── api/                    # API definitions
│   └── openapi.yaml
├── web/                    # Web assets
│   ├── templates/
│   └── static/
├── configs/                # Configuration files
├── docs/                   # Documentation
├── scripts/                # Build scripts
├── test/                   # Additional test files
├── go.mod
└── README.md
```

### Package Visibility
```go
// Exported (public) names start with uppercase
func ExportedFunction() {}
type ExportedType struct {}

// Unexported (private) names start with lowercase
func unexportedFunction() {}
type unexportedType struct {}
```

## Import Patterns

### Basic Import
```go
import "fmt"
import "strings"

// Or grouped
import (
    "fmt"
    "strings"
)
```

### Named Imports
```go
import (
    "fmt"
    maths "math"
    . "strings"      // dot import (avoid)
    _ "image/png"    // blank import
)
```

### Relative Imports
```go
// Inside example.com/myproject/cmd/server/main.go
import (
    "example.com/myproject/internal/auth"
    "example.com/myproject/pkg/models"
)
```

## Dependency Management

### Adding Dependencies
```bash
# Add a specific version
go get github.com/pkg/errors@v0.9.1

# Add latest version
go get github.com/pkg/errors

# Add from specific branch
go get github.com/pkg/errors@master
```

### Updating Dependencies
```bash
# Update all dependencies
go get -u ./...

# Update specific dependency
go get -u github.com/pkg/errors

# Update patch releases only
go get -u=patch ./...
```

### Versioning
```bash
# Tag a new version
git tag v1.0.0
git push origin v1.0.0

# Use semantic versioning
v1.0.0  # Major.Minor.Patch
```

## Best Practices

### 1. Package Organization
- One package per directory
- Package name matches directory name
- Keep packages focused and cohesive
- Use internal/ for private packages

### 2. Import Organization
```go
import (
    // Standard library
    "fmt"
    "strings"
    
    // Third party
    "github.com/gin-gonic/gin"
    "github.com/go-sql-driver/mysql"
    
    // Local packages
    "example.com/myproject/internal/auth"
    "example.com/myproject/pkg/models"
)
```

### 3. Module Management
- Use go.mod for dependency management
- Keep dependencies up to date
- Use go mod tidy regularly
- Pin dependency versions

### 4. Documentation
```go
// Package level documentation
package mypackage // import "example.com/myproject/pkg/mypackage"

// Type documentation
// User represents a system user with authentication information.
type User struct {
    // Name is the user's full name
    Name string
    
    // Email must be a valid email address
    Email string
}

// Function documentation
// CreateUser creates a new user with the given name and email.
// It returns an error if the email is invalid.
func CreateUser(name, email string) (*User, error) {
    // ...
}
```

## Exercises

### Exercise 1: Create a Module
Create a new module with multiple packages:

```bash
# Create module structure
mkdir -p myapp/{cmd/server,internal/auth,pkg/models}
cd myapp

# Initialize module
go mod init example.com/myapp

# Create main package
cat > cmd/server/main.go << EOF
package main

import (
    "example.com/myapp/internal/auth"
    "example.com/myapp/pkg/models"
)

func main() {
    // Use imported packages
}
EOF
```

### Exercise 2: Package Organization
Create a well-organized package structure:

```go
// pkg/models/user.go
package models

type User struct {
    ID    int
    Name  string
    Email string
}

// internal/auth/auth.go
package auth

import "example.com/myapp/pkg/models"

func Authenticate(user *models.User) bool {
    // Implementation
}
```

### Exercise 3: Documentation
Write well-documented packages:

```go
// Package calculator provides basic mathematical operations.
package calculator

// Add returns the sum of two integers.
// It demonstrates proper function documentation.
func Add(a, b int) int {
    return a + b
}
```

## Common Patterns

### 1. Factory Pattern
```go
// pkg/database/database.go
package database

type DB interface {
    Connect() error
    Close() error
}

func NewDB(driver string) DB {
    switch driver {
    case "postgres":
        return &PostgresDB{}
    case "mysql":
        return &MySQLDB{}
    default:
        return &SQLiteDB{}
    }
}
```

### 2. Options Pattern
```go
// pkg/server/server.go
package server

type ServerOption func(*Server)

func WithPort(port int) ServerOption {
    return func(s *Server) {
        s.port = port
    }
}

func NewServer(options ...ServerOption) *Server {
    s := &Server{
        port: 8080, // default
    }
    for _, opt := range options {
        opt(s)
    }
    return s
}
```

### 3. Builder Pattern
```go
// pkg/query/builder.go
package query

type QueryBuilder struct {
    table  string
    where  string
    limit  int
}

func (q *QueryBuilder) From(table string) *QueryBuilder {
    q.table = table
    return q
}
```

## Next Steps
- Create a multi-package project
- Practice package organization
- Write package documentation
- Use dependency management
- Move on to Data Structures and Algorithms
