package main

import (
	"fmt"
	"sync"
	"testing"
)

// TestCalculator tests the Calculator type
func TestCalculator(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		calc := &Calculator{}

		// Test Add
		if got := calc.Add(5); got != 5 {
			t.Errorf("Add(5) = %.2f; want 5", got)
		}

		// Test Subtract
		if got := calc.Subtract(3); got != 2 {
			t.Errorf("Subtract(3) = %.2f; want 2", got)
		}

		// Test Multiply
		if got := calc.Multiply(4); got != 8 {
			t.Errorf("Multiply(4) = %.2f; want 8", got)
		}

		// Test Divide
		if got, err := calc.Divide(2); err != nil || got != 4 {
			t.Errorf("Divide(2) = %.2f, %v; want 4, nil", got, err)
		}
	})

	t.Run("division by zero", func(t *testing.T) {
		calc := &Calculator{}
		calc.Add(10)
		if _, err := calc.Divide(0); err == nil {
			t.Error("Divide(0) should return error")
		}
	})

	t.Run("concurrent operations", func(t *testing.T) {
		calc := &Calculator{}
		var wg sync.WaitGroup
		n := 100

		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				calc.Add(1)
			}()
		}
		wg.Wait()

		if got := calc.GetMemory(); got != float64(n) {
			t.Errorf("After %d concurrent adds, memory = %.2f; want %d", n, got, n)
		}
	})
}

// TestStringProcessor tests the StringProcessor type
func TestStringProcessor(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"single character", "a", "a"},
		{"simple string", "hello", "olleh"},
		{"unicode string", "世界", "界世"},
	}

	sp := NewStringProcessor()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sp.Reverse(tt.input); got != tt.want {
				t.Errorf("Reverse(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("palindrome check", func(t *testing.T) {
		palindromes := []struct {
			input string
			want  bool
		}{
			{"", true},
			{"a", true},
			{"race a car", false},
			{"A man a plan a canal Panama", true},
			{"Was it a car or a cat I saw", true},
		}

		for _, tt := range palindromes {
			t.Run(tt.input, func(t *testing.T) {
				if got := sp.IsPalindrome(tt.input); got != tt.want {
					t.Errorf("IsPalindrome(%q) = %v; want %v", tt.input, got, tt.want)
				}
			})
		}
	})
}

// TestCounter tests the Counter type
func TestCounter(t *testing.T) {
	t.Run("concurrent increments", func(t *testing.T) {
		counter := &Counter{}
		var wg sync.WaitGroup
		n := 1000

		// Launch n goroutines to increment
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				counter.Increment()
			}()
		}
		wg.Wait()

		if got := counter.GetValue(); got != int64(n) {
			t.Errorf("After %d increments, counter = %d; want %d", n, got, n)
		}
	})
}

// TestDataProcessor tests the DataProcessor type
func TestDataProcessor(t *testing.T) {
	t.Run("basic processing", func(t *testing.T) {
		dp := NewDataProcessor(5)
		if err := dp.Process(10); err != nil {
			t.Errorf("Process(10) returned error: %v", err)
		}

		want := 50 // 5 elements * 10
		if got := dp.GetSum(); got != want {
			t.Errorf("GetSum() = %d; want %d", got, want)
		}
	})

	t.Run("negative value", func(t *testing.T) {
		dp := NewDataProcessor(5)
		if err := dp.Process(-1); err == nil {
			t.Error("Process(-1) should return error")
		}
	})

	t.Run("concurrent processing", func(t *testing.T) {
		dp := NewDataProcessor(100)
		var wg sync.WaitGroup
		n := 10

		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				dp.Process(val)
			}(i)
		}
		wg.Wait()

		// Just verify no race conditions occurred
		dp.GetSum()
	})
}

// TestFibonacci tests the Fibonacci function
func TestFibonacci(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{6, 8},
		{7, 13},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("n=%d", tt.n), func(t *testing.T) {
			if got := Fibonacci(tt.n); got != tt.want {
				t.Errorf("Fibonacci(%d) = %d; want %d", tt.n, got, tt.want)
			}
		})
	}
}

// TestConcatStrings tests string concatenation methods
func TestConcatStrings(t *testing.T) {
	strs := []string{"Hello", ", ", "World", "!"}
	want := "Hello, World!"

	methods := []string{"plus", "builder", "join"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			if got := ConcatStrings(strs, method); got != want {
				t.Errorf("ConcatStrings(strs, %q) = %q; want %q", method, got, want)
			}
		})
	}
}

// Benchmark tests
func BenchmarkCalculator(b *testing.B) {
	calc := &Calculator{}
	b.Run("Add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			calc.Add(1)
		}
	})

	b.Run("AddConcurrent", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				calc.Add(1)
			}
		})
	})
}

func BenchmarkStringProcessor(b *testing.B) {
	sp := NewStringProcessor()
	input := "Hello, World!"

	b.Run("Reverse", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sp.Reverse(input)
		}
	})

	b.Run("IsPalindrome", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sp.IsPalindrome(input)
		}
	})
}

func BenchmarkCounter(b *testing.B) {
	counter := &Counter{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Increment()
		}
	})
}

func BenchmarkFibonacci(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fibonacci(10)
	}
}

func BenchmarkConcatStrings(b *testing.B) {
	strs := []string{"Hello", ", ", "World", "!"}

	b.Run("plus", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ConcatStrings(strs, "plus")
		}
	})

	b.Run("builder", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ConcatStrings(strs, "builder")
		}
	})

	b.Run("join", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ConcatStrings(strs, "join")
		}
	})
}
