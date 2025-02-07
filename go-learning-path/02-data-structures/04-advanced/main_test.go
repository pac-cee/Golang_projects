package main

import (
	"reflect"
	"testing"
)

// TestMinHeap tests heap operations
func TestMinHeap(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		heap := NewMinHeap()
		values := []int{5, 3, 7, 1, 4}

		// Test Insert
		for _, v := range values {
			heap.Insert(v)
		}

		// Test ExtractMin
		expected := []int{1, 3, 4, 5, 7}
		for _, want := range expected {
			got, err := heap.ExtractMin()
			if err != nil {
				t.Errorf("ExtractMin() error = %v", err)
			}
			if got != want {
				t.Errorf("ExtractMin() = %v, want %v", got, want)
			}
		}

		// Test empty heap
		_, err := heap.ExtractMin()
		if err == nil {
			t.Error("ExtractMin() on empty heap should return error")
		}
	})
}

// TestTrie tests trie operations
func TestTrie(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		trie := NewTrie()
		words := []string{"apple", "app", "apricot", "banana"}

		// Test Insert and Search
		for _, word := range words {
			trie.Insert(word)
		}

		tests := []struct {
			word     string
			exists   bool
			message  string
		}{
			{"apple", true, "existing word"},
			{"app", true, "existing word"},
			{"apt", false, "non-existing word"},
			{"", true, "empty string"},
		}

		for _, tt := range tests {
			t.Run(tt.message, func(t *testing.T) {
				if got := trie.Search(tt.word); got != tt.exists {
					t.Errorf("Search(%q) = %v, want %v", tt.word, got, tt.exists)
				}
			})
		}

		// Test StartsWith
		prefixes := []struct {
			prefix   string
			exists   bool
			message  string
		}{
			{"ap", true, "existing prefix"},
			{"ban", true, "existing prefix"},
			{"cat", false, "non-existing prefix"},
		}

		for _, tt := range prefixes {
			t.Run(tt.message, func(t *testing.T) {
				if got := trie.StartsWith(tt.prefix); got != tt.exists {
					t.Errorf("StartsWith(%q) = %v, want %v", tt.prefix, got, tt.exists)
				}
			})
		}
	})
}

// TestDisjointSet tests union-find operations
func TestDisjointSet(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		ds := NewDisjointSet(5)

		// Test initial state
		for i := 0; i < 5; i++ {
			if got := ds.Find(i); got != i {
				t.Errorf("Find(%d) = %d, want %d", i, got, i)
			}
		}

		// Test Union
		unions := [][2]int{{0, 1}, {2, 3}, {1, 2}}
		for _, u := range unions {
			ds.Union(u[0], u[1])
		}

		// Test connected components
		tests := []struct {
			x, y     int
			connected bool
			message  string
		}{
			{0, 3, true, "should be connected"},
			{1, 4, false, "should not be connected"},
		}

		for _, tt := range tests {
			t.Run(tt.message, func(t *testing.T) {
				if got := ds.Find(tt.x) == ds.Find(tt.y); got != tt.connected {
					t.Errorf("Find(%d) == Find(%d) = %v, want %v", 
						tt.x, tt.y, got, tt.connected)
				}
			})
		}
	})
}

// TestBloomFilter tests bloom filter operations
func TestBloomFilter(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		bf := NewBloomFilter(100, 3)
		items := []string{"apple", "banana", "orange"}

		// Test Add and Contains
		for _, item := range items {
			bf.Add(item)
		}

		// Test existing items
		for _, item := range items {
			if !bf.Contains(item) {
				t.Errorf("Contains(%q) = false, want true", item)
			}
		}

		// Test false positives
		falsePositives := 0
		nonExistingItems := []string{"grape", "melon", "peach"}
		for _, item := range nonExistingItems {
			if bf.Contains(item) {
				falsePositives++
			}
		}

		// Check false positive rate (should be reasonable)
		falsePositiveRate := float64(falsePositives) / float64(len(nonExistingItems))
		if falsePositiveRate > 0.1 {
			t.Errorf("False positive rate too high: %v", falsePositiveRate)
		}
	})
}

// TestLRUCache tests LRU cache operations
func TestLRUCache(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		cache := NewLRUCache(3)

		// Test Put and Get
		cache.Put("A", 1)
		cache.Put("B", 2)
		cache.Put("C", 3)

		tests := []struct {
			key      string
			want     interface{}
			exists   bool
			message  string
		}{
			{"A", 1, true, "existing key"},
			{"B", 2, true, "existing key"},
			{"D", nil, false, "non-existing key"},
		}

		for _, tt := range tests {
			t.Run(tt.message, func(t *testing.T) {
				got, exists := cache.Get(tt.key)
				if exists != tt.exists {
					t.Errorf("Get(%q) exists = %v, want %v", tt.key, exists, tt.exists)
				}
				if exists && !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Get(%q) = %v, want %v", tt.key, got, tt.want)
				}
			})
		}

		// Test eviction
		cache.Put("D", 4) // Should evict A
		if _, exists := cache.Get("A"); exists {
			t.Error("Key 'A' should have been evicted")
		}

		// Test update existing
		cache.Put("B", 20)
		if val, _ := cache.Get("B"); val != 20 {
			t.Errorf("Get(B) = %v, want 20", val)
		}
	})
}

// Benchmark tests
func BenchmarkMinHeap(b *testing.B) {
	heap := NewMinHeap()
	
	b.Run("Insert", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			heap.Insert(i)
		}
	})

	b.Run("ExtractMin", func(b *testing.B) {
		heap := NewMinHeap()
		for i := 0; i < b.N; i++ {
			heap.Insert(i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			heap.ExtractMin()
		}
	})
}

func BenchmarkTrie(b *testing.B) {
	words := []string{"apple", "application", "append", "banana", "ball"}
	
	b.Run("Insert", func(b *testing.B) {
		trie := NewTrie()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			trie.Insert(words[i%len(words)])
		}
	})

	b.Run("Search", func(b *testing.B) {
		trie := NewTrie()
		for _, word := range words {
			trie.Insert(word)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			trie.Search(words[i%len(words)])
		}
	})
}

func BenchmarkDisjointSet(b *testing.B) {
	size := 1000
	
	b.Run("Union", func(b *testing.B) {
		ds := NewDisjointSet(size)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ds.Union(i%size, (i+1)%size)
		}
	})

	b.Run("Find", func(b *testing.B) {
		ds := NewDisjointSet(size)
		for i := 0; i < size-1; i++ {
			ds.Union(i, i+1)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ds.Find(i % size)
		}
	})
}

func BenchmarkBloomFilter(b *testing.B) {
	items := []string{"apple", "banana", "orange", "grape", "melon"}
	
	b.Run("Add", func(b *testing.B) {
		bf := NewBloomFilter(1000, 3)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bf.Add(items[i%len(items)])
		}
	})

	b.Run("Contains", func(b *testing.B) {
		bf := NewBloomFilter(1000, 3)
		for _, item := range items {
			bf.Add(item)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bf.Contains(items[i%len(items)])
		}
	})
}

func BenchmarkLRUCache(b *testing.B) {
	cache := NewLRUCache(1000)
	
	b.Run("Put", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cache.Put(string(rune(i%26+'A')), i)
		}
	})

	b.Run("Get", func(b *testing.B) {
		cache := NewLRUCache(1000)
		for i := 0; i < 1000; i++ {
			cache.Put(string(rune(i%26+'A')), i)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.Get(string(rune(i%26+'A')))
		}
	})
}
