package main

import (
	"reflect"
	"testing"
)

// TestDataProcessor tests the DataProcessor methods
func TestDataProcessor(t *testing.T) {
	// Test data
	data := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}

	t.Run("NewDataProcessor", func(t *testing.T) {
		processor := NewDataProcessor(data)
		if !reflect.DeepEqual(processor.data, data) {
			t.Errorf("NewDataProcessor() = %v, want %v", processor.data, data)
		}

		// Verify it's a copy
		data[0] = 100
		if processor.data[0] == 100 {
			t.Error("NewDataProcessor() did not create a copy of the data")
		}
	})

	t.Run("Sort", func(t *testing.T) {
		processor := NewDataProcessor(data)
		processor.Sort()
		expected := []int{1, 1, 2, 3, 3, 4, 5, 5, 5, 6, 9}
		if !reflect.DeepEqual(processor.data, expected) {
			t.Errorf("Sort() = %v, want %v", processor.data, expected)
		}
	})

	t.Run("Reverse", func(t *testing.T) {
		processor := NewDataProcessor([]int{1, 2, 3, 4, 5})
		processor.Reverse()
		expected := []int{5, 4, 3, 2, 1}
		if !reflect.DeepEqual(processor.data, expected) {
			t.Errorf("Reverse() = %v, want %v", processor.data, expected)
		}
	})

	t.Run("Filter", func(t *testing.T) {
		processor := NewDataProcessor(data)
		evens := processor.Filter(func(x int) bool {
			return x%2 == 0
		})
		expected := []int{4, 2, 6}
		if !reflect.DeepEqual(evens, expected) {
			t.Errorf("Filter(evens) = %v, want %v", evens, expected)
		}
	})

	t.Run("Map", func(t *testing.T) {
		processor := NewDataProcessor([]int{1, 2, 3})
		doubled := processor.Map(func(x int) int {
			return x * 2
		})
		expected := []int{2, 4, 6}
		if !reflect.DeepEqual(doubled, expected) {
			t.Errorf("Map(double) = %v, want %v", doubled, expected)
		}
	})

	t.Run("RemoveDuplicates", func(t *testing.T) {
		processor := NewDataProcessor(data)
		processor.RemoveDuplicates()
		expected := []int{3, 1, 4, 5, 9, 2, 6}
		if !reflect.DeepEqual(processor.data, expected) {
			t.Errorf("RemoveDuplicates() = %v, want %v", processor.data, expected)
		}
	})

	t.Run("SlidingWindow", func(t *testing.T) {
		processor := NewDataProcessor([]int{1, 2, 3, 4, 5})
		windows := processor.SlidingWindow(3)
		expected := [][]int{
			{1, 2, 3},
			{2, 3, 4},
			{3, 4, 5},
		}
		if !reflect.DeepEqual(windows, expected) {
			t.Errorf("SlidingWindow(3) = %v, want %v", windows, expected)
		}
	})

	t.Run("Chunk", func(t *testing.T) {
		processor := NewDataProcessor([]int{1, 2, 3, 4, 5})
		chunks := processor.Chunk(2)
		expected := [][]int{
			{1, 2},
			{3, 4},
			{5},
		}
		if !reflect.DeepEqual(chunks, expected) {
			t.Errorf("Chunk(2) = %v, want %v", chunks, expected)
		}
	})

	t.Run("Rotate", func(t *testing.T) {
		tests := []struct {
			name      string
			input     []int
			positions int
			expected  []int
		}{
			{
				name:      "rotate right by 2",
				input:     []int{1, 2, 3, 4, 5},
				positions: 2,
				expected: []int{4, 5, 1, 2, 3},
			},
			{
				name:      "rotate left by 2",
				input:     []int{1, 2, 3, 4, 5},
				positions: -2,
				expected: []int{3, 4, 5, 1, 2},
			},
			{
				name:      "rotate by length",
				input:     []int{1, 2, 3},
				positions: 3,
				expected: []int{1, 2, 3},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				processor := NewDataProcessor(tt.input)
				processor.Rotate(tt.positions)
				if !reflect.DeepEqual(processor.data, tt.expected) {
					t.Errorf("Rotate(%d) = %v, want %v", 
						tt.positions, processor.data, tt.expected)
				}
			})
		}
	})
}

// TestEdgeCases tests edge cases for DataProcessor methods
func TestEdgeCases(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		processor := NewDataProcessor([]int{})

		// Test all methods with empty slice
		processor.Sort()
		processor.Reverse()
		processor.RemoveDuplicates()
		processor.Rotate(5)

		if len(processor.data) != 0 {
			t.Error("Operations on empty slice should maintain empty slice")
		}

		// Test window operations
		if windows := processor.SlidingWindow(1); windows != nil {
			t.Error("SlidingWindow on empty slice should return nil")
		}

		if chunks := processor.Chunk(1); chunks != nil {
			t.Error("Chunk on empty slice should return nil")
		}
	})

	t.Run("invalid parameters", func(t *testing.T) {
		processor := NewDataProcessor([]int{1, 2, 3})

		// Test invalid window size
		if windows := processor.SlidingWindow(0); windows != nil {
			t.Error("SlidingWindow with size 0 should return nil")
		}

		if windows := processor.SlidingWindow(4); windows != nil {
			t.Error("SlidingWindow with size > len should return nil")
		}

		// Test invalid chunk size
		if chunks := processor.Chunk(0); chunks != nil {
			t.Error("Chunk with size 0 should return nil")
		}
	})
}

// BenchmarkDataProcessor benchmarks the DataProcessor methods
func BenchmarkDataProcessor(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	b.Run("Sort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			processor := NewDataProcessor(data)
			processor.Sort()
		}
	})

	b.Run("RemoveDuplicates", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			processor := NewDataProcessor(data)
			processor.RemoveDuplicates()
		}
	})

	b.Run("SlidingWindow", func(b *testing.B) {
		processor := NewDataProcessor(data)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			processor.SlidingWindow(10)
		}
	})

	b.Run("Rotate", func(b *testing.B) {
		processor := NewDataProcessor(data)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			processor.Rotate(10)
		}
	})
}
