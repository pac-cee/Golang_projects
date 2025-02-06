package main

import (
	"fmt"
	"sync"
	"time"
)

// User represents a basic user structure
type User struct {
	ID       int
	Username string
	Email    string
}

// UserInterface defines behavior for user operations
type UserInterface interface {
	GetEmail() string
	UpdateEmail(newEmail string)
}

// GetEmail returns user's email
func (u *User) GetEmail() string {
	return u.Email
}

// UpdateEmail updates user's email
func (u *User) UpdateEmail(newEmail string) {
	u.Email = newEmail
}

// SafeCounter is a thread-safe counter using mutex
type SafeCounter struct {
	mu    sync.Mutex
	count int
}

// Increment increases the counter by 1
func (sc *SafeCounter) Increment() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.count++
}

// GetCount returns the current count
func (sc *SafeCounter) GetCount() int {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.count
}

// processUser demonstrates error handling
func processUser(u *User) error {
	if u == nil {
		return fmt.Errorf("user cannot be nil")
	}
	if u.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	return nil
}

// demonstrateChannel shows channel usage
func demonstrateChannel(done chan bool) {
	fmt.Println("Processing...")
	time.Sleep(time.Second) // Simulate work
	done <- true
}

func main() {
	// 1. Basic variable declaration and structs
	user := &User{
		ID:       1,
		Username: "johndoe",
		Email:    "john@example.com",
	}

	// 2. Interface usage
	var userInterface UserInterface = user
	userInterface.UpdateEmail("newemail@example.com")
	fmt.Printf("Updated email: %s\n", userInterface.GetEmail())

	// 3. Error handling
	if err := processUser(user); err != nil {
		fmt.Printf("Error processing user: %v\n", err)
	}

	// 4. Goroutines and channels
	done := make(chan bool)
	go demonstrateChannel(done)
	<-done
	fmt.Println("Channel demonstration completed")

	// 5. Mutex and concurrent access
	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// Launch multiple goroutines to increment counter
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Final count: %d\n", counter.GetCount())
}
