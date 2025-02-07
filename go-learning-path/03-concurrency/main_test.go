package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

// TestSafeCounter tests the thread-safe counter
func TestSafeCounter(t *testing.T) {
	t.Run("concurrent increments", func(t *testing.T) {
		counter := NewSafeCounter()
		var wg sync.WaitGroup
		n := 100

		// Launch n goroutines to increment
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				counter.Increment()
			}()
		}

		wg.Wait()
		value := counter.Value()
		if value != n {
			t.Errorf("Counter = %d; want %d", value, n)
		}
	})
}

// TestRateLimiter tests the rate limiter
func TestRateLimiter(t *testing.T) {
	t.Run("rate limiting", func(t *testing.T) {
		limiter := NewRateLimiter(100 * time.Millisecond)
		defer limiter.Stop()

		start := time.Now()
		count := 0
		timeout := time.After(550 * time.Millisecond)

		// Should allow ~5 actions in 500ms
		for {
			select {
			case <-timeout:
				if count < 4 || count > 6 {
					t.Errorf("Got %d actions; want 5 ±1", count)
				}
				return
			case <-limiter.Allow():
				count++
			}
		}
	})
}

// TestWorker tests the worker pool
func TestWorker(t *testing.T) {
	t.Run("worker processing", func(t *testing.T) {
		jobs := make(chan int, 5)
		results := make(chan int, 5)

		// Start a worker
		worker := NewWorker(1, jobs, results)
		worker.Start()

		// Send test jobs
		testJobs := []int{1, 2, 3, 4, 5}
		for _, j := range testJobs {
			jobs <- j
		}
		close(jobs)

		// Collect and verify results
		resultMap := make(map[int]bool)
		for i := 0; i < len(testJobs); i++ {
			result := <-results
			resultMap[result] = true
		}

		// Verify all expected results are present
		for _, j := range testJobs {
			expected := j * 2
			if !resultMap[expected] {
				t.Errorf("Missing result: %d", expected)
			}
		}
	})
}

// TestPipeline tests the pipeline pattern
func TestPipeline(t *testing.T) {
	t.Run("pipeline processing", func(t *testing.T) {
		done := make(chan struct{})
		defer close(done)

		pipeline := NewPipeline(done)
		numbers := pipeline.Generator(5)
		squares := pipeline.Square(numbers)
		filtered := pipeline.Filter(squares)

		// Collect results
		var results []int
		for n := range filtered {
			results = append(results, n)
		}

		// Verify results
		expected := []int{4, 16} // Only even squares of 1-5
		if len(results) != len(expected) {
			t.Errorf("Got %d results; want %d", len(results), len(expected))
		}

		resultMap := make(map[int]bool)
		for _, r := range results {
			resultMap[r] = true
		}

		for _, e := range expected {
			if !resultMap[e] {
				t.Errorf("Missing expected value: %d", e)
			}
		}
	})
}

// TestFanOutFanIn tests the fan-out, fan-in pattern
func TestFanOutFanIn(t *testing.T) {
	t.Run("fan-out fan-in", func(t *testing.T) {
		done := make(chan struct{})
		defer close(done)

		// Create input channel
		input := make(chan int)
		go func() {
			defer close(input)
			for i := 1; i <= 5; i++ {
				input <- i
			}
		}()

		// Fan out to 3 channels
		channels := FanOut(input, 3)

		// Fan in results
		merged := FanIn(done, channels...)

		// Collect and verify results
		results := make(map[int]bool)
		for i := 0; i < 5; i++ {
			result := <-merged
			results[result] = true
		}

		// Verify all squares are present
		for i := 1; i <= 5; i++ {
			square := i * i
			if !results[square] {
				t.Errorf("Missing square: %d", square)
			}
		}
	})
}

// TestAsyncOperation tests async operations with error handling
func TestAsyncOperation(t *testing.T) {
	t.Run("async operation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		result := <-AsyncOperation(ctx)
		if result.Err != nil {
			// Error is acceptable
			return
		}

		// If no error, value should be in range [0, 100)
		if result.Value < 0 || result.Value >= 100 {
			t.Errorf("Value %d out of expected range [0, 100)", result.Value)
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		result := <-AsyncOperation(ctx)
		if result.Err != context.Canceled {
			t.Errorf("Expected context.Canceled error, got: %v", result.Err)
		}
	})
}

// Benchmark tests
func BenchmarkSafeCounter(b *testing.B) {
	counter := NewSafeCounter()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Increment()
		}
	})
}

func BenchmarkWorkerPool(b *testing.B) {
	jobs := make(chan int, b.N)
	results := make(chan int, b.N)

	// Start 3 workers
	for w := 1; w <= 3; w++ {
		worker := NewWorker(w, jobs, results)
		worker.Start()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		jobs <- i
	}
	close(jobs)

	for i := 0; i < b.N; i++ {
		<-results
	}
}

func BenchmarkPipeline(b *testing.B) {
	done := make(chan struct{})
	defer close(done)

	b.ResetTimer()
	pipeline := NewPipeline(done)
	numbers := pipeline.Generator(b.N)
	squares := pipeline.Square(numbers)
	filtered := pipeline.Filter(squares)

	for range filtered {
		// Drain the pipeline
	}
}

func BenchmarkFanOutFanIn(b *testing.B) {
	done := make(chan struct{})
	defer close(done)

	input := make(chan int)
	go func() {
		defer close(input)
		for i := 0; i < b.N; i++ {
			input <- i
		}
	}()

	channels := FanOut(input, 3)
	merged := FanIn(done, channels...)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-merged
	}
}
