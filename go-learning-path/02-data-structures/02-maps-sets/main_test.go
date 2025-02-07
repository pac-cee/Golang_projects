package main

import (
	"reflect"
	"testing"
	"time"
)

// TestSet tests the Set data structure
func TestSet(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		set := NewSet()

		// Test Add and Contains
		set.Add("apple")
		if !set.Contains("apple") {
			t.Error("set should contain 'apple'")
		}

		// Test Size
		if set.Size() != 1 {
			t.Errorf("set size = %d; want 1", set.Size())
		}

		// Test Remove
		set.Remove("apple")
		if set.Contains("apple") {
			t.Error("set should not contain 'apple' after removal")
		}

		// Test Items
		set.Add("banana")
		set.Add("cherry")
		items := set.Items()
		if len(items) != 2 {
			t.Errorf("items length = %d; want 2", len(items))
		}
	})

	t.Run("set operations", func(t *testing.T) {
		set1 := NewSet()
		set1.Add("apple")
		set1.Add("banana")
		set1.Add("cherry")

		set2 := NewSet()
		set2.Add("banana")
		set2.Add("cherry")
		set2.Add("date")

		// Test Union
		union := set1.Union(set2)
		expectedUnion := []string{"apple", "banana", "cherry", "date"}
		if !sameElements(union.Items(), expectedUnion) {
			t.Errorf("Union = %v; want %v", union.Items(), expectedUnion)
		}

		// Test Intersection
		intersection := set1.Intersection(set2)
		expectedIntersection := []string{"banana", "cherry"}
		if !sameElements(intersection.Items(), expectedIntersection) {
			t.Errorf("Intersection = %v; want %v", intersection.Items(), expectedIntersection)
		}

		// Test Difference
		difference := set1.Difference(set2)
		expectedDifference := []string{"apple"}
		if !sameElements(difference.Items(), expectedDifference) {
			t.Errorf("Difference = %v; want %v", difference.Items(), expectedDifference)
		}
	})
}

// TestCache tests the Cache data structure
func TestCache(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		cache := NewCache()

		// Test Set and Get
		cache.Set("key1", "value1", time.Minute)
		if value, exists := cache.Get("key1"); !exists || value != "value1" {
			t.Errorf("cache.Get(key1) = %v, %v; want value1, true", value, exists)
		}

		// Test non-existent key
		if _, exists := cache.Get("key2"); exists {
			t.Error("cache should not contain key2")
		}

		// Test Delete
		cache.Delete("key1")
		if _, exists := cache.Get("key1"); exists {
			t.Error("cache should not contain key1 after deletion")
		}
	})

	t.Run("expiration", func(t *testing.T) {
		cache := NewCache()

		// Add item with short TTL
		cache.Set("key1", "value1", time.Millisecond*100)

		// Wait for expiration
		time.Sleep(time.Millisecond * 200)

		// Check if item has expired
		if _, exists := cache.Get("key1"); exists {
			t.Error("cache item should have expired")
		}
	})
}

// TestFrequencyCounter tests the FrequencyCounter
func TestFrequencyCounter(t *testing.T) {
	t.Run("basic counting", func(t *testing.T) {
		counter := NewFrequencyCounter()

		// Add words
		words := []string{"apple", "banana", "apple", "cherry", "apple", "date"}
		for _, word := range words {
			counter.Add(word)
		}

		// Test counts
		tests := []struct {
			word     string
			expected int
		}{
			{"apple", 3},
			{"banana", 1},
			{"cherry", 1},
			{"date", 1},
			{"fig", 0}, // non-existent word
		}

		for _, tt := range tests {
			if count := counter.Count(tt.word); count != tt.expected {
				t.Errorf("Count(%s) = %d; want %d", tt.word, count, tt.expected)
			}
		}
	})

	t.Run("top N", func(t *testing.T) {
		counter := NewFrequencyCounter()

		// Add words with different frequencies
		words := []string{
			"apple", "apple", "apple",
			"banana", "banana",
			"cherry",
		}
		for _, word := range words {
			counter.Add(word)
		}

		// Test TopN
		top2 := counter.TopN(2)
		if len(top2) != 2 {
			t.Errorf("TopN(2) returned %d items; want 2", len(top2))
		}
		if top2["apple"] != 3 {
			t.Errorf("TopN(2)['apple'] = %d; want 3", top2["apple"])
		}
		if top2["banana"] != 2 {
			t.Errorf("TopN(2)['banana'] = %d; want 2", top2["banana"])
		}
	})
}

// Benchmark tests
func BenchmarkSet(b *testing.B) {
	set := NewSet()
	items := []string{"apple", "banana", "cherry", "date"}

	b.Run("Add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			set.Add(items[i%len(items)])
		}
	})

	b.Run("Contains", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			set.Contains(items[i%len(items)])
		}
	})

	b.Run("Union", func(b *testing.B) {
		set1 := NewSet()
		set2 := NewSet()
		for _, item := range items[:2] {
			set1.Add(item)
		}
		for _, item := range items[2:] {
			set2.Add(item)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			set1.Union(set2)
		}
	})
}

func BenchmarkCache(b *testing.B) {
	cache := NewCache()

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key%d", i)
			cache.Set(key, i, time.Minute)
		}
	})

	b.Run("Get", func(b *testing.B) {
		key := "testKey"
		cache.Set(key, "value", time.Minute)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cache.Get(key)
		}
	})
}

func BenchmarkFrequencyCounter(b *testing.B) {
	counter := NewFrequencyCounter()
	words := []string{"apple", "banana", "cherry", "date"}

	b.Run("Add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			counter.Add(words[i%len(words)])
		}
	})

	b.Run("TopN", func(b *testing.B) {
		// Pre-populate counter
		for _, word := range words {
			counter.Add(word)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			counter.TopN(3)
		}
	})
}

// Helper function to compare slices regardless of order
func sameElements(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	aMap := make(map[string]int)
	for _, item := range a {
		aMap[item]++
	}

	bMap := make(map[string]int)
	for _, item := range b {
		bMap[item]++
	}

	return reflect.DeepEqual(aMap, bMap)
}
