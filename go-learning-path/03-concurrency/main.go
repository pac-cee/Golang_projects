// Package main demonstrates concurrency patterns in Go
package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// SafeCounter is a thread-safe counter using channels
type SafeCounter struct {
	value chan int
	read  chan chan int
}

// NewSafeCounter creates a new thread-safe counter
func NewSafeCounter() *SafeCounter {
	c := &SafeCounter{
		value: make(chan int),
		read:  make(chan chan int),
	}
	go c.run()
	return c
}

func (c *SafeCounter) run() {
	value := 0
	for {
		select {
		case <-c.value:
			value++
		case reader := <-c.read:
			reader <- value
		}
	}
}

// Increment increases the counter by one
func (c *SafeCounter) Increment() {
	c.value <- 1
}

// Value returns the current value of the counter
func (c *SafeCounter) Value() int {
	reader := make(chan int)
	c.read <- reader
	return <-reader
}

// RateLimiter implements a simple rate limiter
type RateLimiter struct {
	ticker *time.Ticker
	stop   chan struct{}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit time.Duration) *RateLimiter {
	return &RateLimiter{
		ticker: time.NewTicker(limit),
		stop:   make(chan struct{}),
	}
}

// Allow returns a channel that receives when action is allowed
func (r *RateLimiter) Allow() <-chan time.Time {
	return r.ticker.C
}

// Stop stops the rate limiter
func (r *RateLimiter) Stop() {
	r.ticker.Stop()
	close(r.stop)
}

// Worker represents a worker pool worker
type Worker struct {
	id      int
	jobs    <-chan int
	results chan<- int
}

// NewWorker creates a new worker
func NewWorker(id int, jobs <-chan int, results chan<- int) *Worker {
	return &Worker{
		id:      id,
		jobs:    jobs,
		results: results,
	}
}

// Start starts the worker
func (w *Worker) Start() {
	go func() {
		for j := range w.jobs {
			fmt.Printf("worker %d processing job %d\n", w.id, j)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			w.results <- j * 2
		}
	}()
}

// Pipeline demonstrates a simple pipeline pattern
type Pipeline struct {
	done <-chan struct{}
}

// NewPipeline creates a new pipeline
func NewPipeline(done <-chan struct{}) *Pipeline {
	return &Pipeline{done: done}
}

// Generator generates integers from 1 to n
func (p *Pipeline) Generator(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			select {
			case <-p.done:
				return
			case out <- i:
			}
		}
	}()
	return out
}

// Square squares numbers from in channel
func (p *Pipeline) Square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case <-p.done:
				return
			case out <- n * n:
			}
		}
	}()
	return out
}

// Filter filters out odd numbers
func (p *Pipeline) Filter(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if n%2 == 0 {
				select {
				case <-p.done:
					return
				case out <- n:
				}
			}
		}
	}()
	return out
}

// FanOut demonstrates the fan-out pattern
func FanOut(in <-chan int, n int) []<-chan int {
	channels := make([]<-chan int, n)
	for i := 0; i < n; i++ {
		channels[i] = processChannel(in)
	}
	return channels
}

func processChannel(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// FanIn combines multiple channels into one
func FanIn(done <-chan struct{}, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	multiplexedStream := make(chan int)

	multiplex := func(c <-chan int) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

// Result represents a computation result with possible error
type Result struct {
	Value int
	Err   error
}

// AsyncOperation demonstrates error handling in concurrent operations
func AsyncOperation(ctx context.Context) <-chan Result {
	results := make(chan Result)
	go func() {
		defer close(results)
		select {
		case <-ctx.Done():
			results <- Result{Err: ctx.Err()}
			return
		case <-time.After(time.Duration(rand.Intn(1000)) * time.Millisecond):
			if rand.Float32() < 0.3 { // 30% chance of error
				results <- Result{Err: fmt.Errorf("random error occurred")}
				return
			}
			results <- Result{Value: rand.Intn(100)}
		}
	}()
	return results
}

func main() {
	// Demonstrate SafeCounter
	counter := NewSafeCounter()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	fmt.Printf("Counter value: %d\n", counter.Value())

	// Demonstrate RateLimiter
	limiter := NewRateLimiter(200 * time.Millisecond)
	go func() {
		for i := 1; i <= 5; i++ {
			<-limiter.Allow()
			fmt.Printf("Action %d allowed\n", i)
		}
		limiter.Stop()
	}()

	// Demonstrate Worker Pool
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// Start 3 workers
	for w := 1; w <= 3; w++ {
		worker := NewWorker(w, jobs, results)
		worker.Start()
	}

	// Send 5 jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect results
	for a := 1; a <= 5; a++ {
		<-results
	}

	// Demonstrate Pipeline
	done := make(chan struct{})
	defer close(done)

	pipeline := NewPipeline(done)
	numbers := pipeline.Generator(10)
	squares := pipeline.Square(numbers)
	filtered := pipeline.Filter(squares)

	fmt.Println("\nPipeline results:")
	for n := range filtered {
		fmt.Printf("%d ", n)
	}
	fmt.Println()

	// Demonstrate Fan-out, Fan-in
	generator := func() <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for i := 1; i <= 10; i++ {
				out <- i
			}
		}()
		return out
	}

	channels := FanOut(generator(), 3)
	merged := FanIn(done, channels...)

	fmt.Println("\nFan-out, Fan-in results:")
	for n := range merged {
		fmt.Printf("%d ", n)
	}
	fmt.Println()

	// Demonstrate error handling with context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	fmt.Println("\nAsync operations results:")
	for i := 0; i < 5; i++ {
		result := <-AsyncOperation(ctx)
		if result.Err != nil {
			fmt.Printf("Operation %d failed: %v\n", i+1, result.Err)
		} else {
			fmt.Printf("Operation %d succeeded: %d\n", i+1, result.Value)
		}
	}

	// Allow some time for goroutines to finish
	time.Sleep(time.Second)
}
