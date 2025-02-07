// Package main demonstrates control structures in Go
package main

import (
	"fmt"
	"time"
)

// UserRole represents user access levels
type UserRole int

const (
	Guest UserRole = iota
	User
	Moderator
	Admin
)

// User represents a system user
type User struct {
	ID       int
	Name     string
	Role     UserRole
	Active   bool
	JoinDate time.Time
}

// checkUserAccess demonstrates if statements and switch
func checkUserAccess(u User, resource string) {
	// If with initialization
	if lastLogin := time.Since(u.JoinDate); lastLogin.Hours() < 24 {
		fmt.Printf("New user detected: %s (joined %v ago)\n", 
			u.Name, lastLogin.Round(time.Minute))
	}

	// Switch on user role
	switch u.Role {
	case Admin:
		fmt.Printf("Admin %s has full access to %s\n", u.Name, resource)
	case Moderator:
		fmt.Printf("Moderator %s has write access to %s\n", u.Name, resource)
	case User:
		fmt.Printf("User %s has read access to %s\n", u.Name, resource)
	default:
		fmt.Printf("Guest has no access to %s\n", resource)
	}
}

// processUsers demonstrates for loops and range
func processUsers(users []User) {
	// Range over slice
	for i, user := range users {
		fmt.Printf("\nProcessing user %d:\n", i+1)
		
		// If-else statement
		if !user.Active {
			fmt.Printf("Skipping inactive user: %s\n", user.Name)
			continue
		}

		checkUserAccess(user, "system")
	}
}

// safeOperation demonstrates defer, panic, and recover
func safeOperation(operation string) (err error) {
	// Defer recovery
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("operation %s failed: %v", operation, r)
		}
	}()

	// Simulate operations that might panic
	switch operation {
	case "dangerous":
		panic("dangerous operation failed")
	case "risky":
		fmt.Println("Risky operation completed")
	default:
		fmt.Println("Safe operation completed")
	}

	return nil
}

// demonstrateDefer shows defer order
func demonstrateDefer() {
	fmt.Println("\nDemonstrating defer:")
	
	// Defers are executed in LIFO order
	defer fmt.Println("First defer (executed last)")
	defer fmt.Println("Second defer (executed second)")
	defer fmt.Println("Third defer (executed first)")
	
	fmt.Println("Main function body")
}

func main() {
	// Create test users
	users := []User{
		{
			ID:       1,
			Name:     "Alice",
			Role:     Admin,
			Active:   true,
			JoinDate: time.Now().Add(-48 * time.Hour),
		},
		{
			ID:       2,
			Name:     "Bob",
			Role:     User,
			Active:   true,
			JoinDate: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:       3,
			Name:     "Charlie",
			Role:     Moderator,
			Active:   false,
			JoinDate: time.Now().Add(-24 * time.Hour),
		},
	}

	// Process users
	processUsers(users)

	// Demonstrate defer
	demonstrateDefer()

	// Demonstrate safe operations
	fmt.Println("\nDemonstrating safe operations:")
	operations := []string{"safe", "risky", "dangerous"}
	
	for _, op := range operations {
		if err := safeOperation(op); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}
