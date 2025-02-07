package main

import (
	"reflect"
	"testing"
)

// TestBinarySearchTree tests BST operations
func TestBinarySearchTree(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		bst := NewBST()
		values := []int{5, 3, 7, 2, 4, 6, 8}

		// Test Insert
		for _, v := range values {
			bst.Insert(v)
		}

		// Test Search
		for _, v := range values {
			if !bst.Search(v) {
				t.Errorf("Search(%d) = false; want true", v)
			}
		}

		// Test non-existent value
		if bst.Search(9) {
			t.Error("Search(9) = true; want false")
		}

		// Test InOrderTraversal
		expected := []int{2, 3, 4, 5, 6, 7, 8}
		result := bst.InOrderTraversal()
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("InOrderTraversal() = %v; want %v", result, expected)
		}
	})
}

// TestAVLTree tests AVL tree operations
func TestAVLTree(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		avl := NewAVL()
		values := []int{5, 3, 7, 2, 4, 6, 8}

		// Test Insert and balance
		for _, v := range values {
			avl.Insert(v)
			if !isBalanced(avl.Root) {
				t.Errorf("Tree not balanced after inserting %d", v)
			}
		}

		// Test height property
		if height := avl.Root.height(); height > 4 {
			t.Errorf("Tree height = %d; want <= 4", height)
		}
	})

	t.Run("rotations", func(t *testing.T) {
		avl := NewAVL()
		
		// Test Left-Left case
		values := []int{30, 20, 10}
		for _, v := range values {
			avl.Insert(v)
		}
		if avl.Root.Value != 20 {
			t.Errorf("Root value after LL rotation = %d; want 20", avl.Root.Value)
		}

		// Test Right-Right case
		avl = NewAVL()
		values = []int{10, 20, 30}
		for _, v := range values {
			avl.Insert(v)
		}
		if avl.Root.Value != 20 {
			t.Errorf("Root value after RR rotation = %d; want 20", avl.Root.Value)
		}
	})
}

// Helper function to check if AVL tree is balanced
func isBalanced(node *AVLNode) bool {
	if node == nil {
		return true
	}

	balance := node.getBalance()
	if balance < -1 || balance > 1 {
		return false
	}

	return isBalanced(node.Left) && isBalanced(node.Right)
}

// TestGraph tests graph operations
func TestGraph(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		graph := NewGraph()

		// Test AddVertex and AddEdge
		graph.AddEdge(0, 1, 4)
		graph.AddEdge(0, 2, 2)
		graph.AddEdge(1, 2, 1)
		graph.AddEdge(1, 3, 5)

		// Verify edges exist
		if weight := graph.vertices[0][1]; weight != 4 {
			t.Errorf("Edge weight (0->1) = %d; want 4", weight)
		}
	})

	t.Run("traversal", func(t *testing.T) {
		graph := NewGraph()
		edges := []struct{ from, to, weight int }{
			{0, 1, 4}, {0, 2, 2}, {1, 2, 1}, {1, 3, 5}, {2, 3, 8},
		}

		for _, e := range edges {
			graph.AddEdge(e.from, e.to, e.weight)
		}

		// Test DFS
		dfs := graph.DFS(0)
		if len(dfs) != 4 {
			t.Errorf("DFS visited %d vertices; want 4", len(dfs))
		}

		// Test BFS
		bfs := graph.BFS(0)
		if len(bfs) != 4 {
			t.Errorf("BFS visited %d vertices; want 4", len(bfs))
		}
	})

	t.Run("dijkstra", func(t *testing.T) {
		graph := NewGraph()
		edges := []struct{ from, to, weight int }{
			{0, 1, 4}, {0, 2, 2}, {1, 2, 1}, {1, 3, 5}, {2, 3, 8},
		}

		for _, e := range edges {
			graph.AddEdge(e.from, e.to, e.weight)
		}

		distances := graph.Dijkstra(0)
		expectedDistances := map[int]int{
			0: 0,  // Distance to self
			1: 4,  // Direct edge
			2: 2,  // Direct edge
			3: 9,  // Path: 0->1->3
		}

		for vertex, distance := range expectedDistances {
			if distances[vertex] != distance {
				t.Errorf("Shortest distance to %d = %d; want %d",
					vertex, distances[vertex], distance)
			}
		}
	})
}

// Benchmark tests
func BenchmarkBST(b *testing.B) {
	values := []int{5, 3, 7, 2, 4, 6, 8}
	
	b.Run("Insert", func(b *testing.B) {
		bst := NewBST()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bst.Insert(values[i%len(values)])
		}
	})

	b.Run("Search", func(b *testing.B) {
		bst := NewBST()
		for _, v := range values {
			bst.Insert(v)
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bst.Search(values[i%len(values)])
		}
	})
}

func BenchmarkAVL(b *testing.B) {
	values := []int{5, 3, 7, 2, 4, 6, 8}
	
	b.Run("Insert", func(b *testing.B) {
		avl := NewAVL()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			avl.Insert(values[i%len(values)])
		}
	})
}

func BenchmarkGraph(b *testing.B) {
	graph := NewGraph()
	edges := []struct{ from, to, weight int }{
		{0, 1, 4}, {0, 2, 2}, {1, 2, 1}, {1, 3, 5}, {2, 3, 8},
	}
	for _, e := range edges {
		graph.AddEdge(e.from, e.to, e.weight)
	}

	b.Run("DFS", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			graph.DFS(0)
		}
	})

	b.Run("BFS", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			graph.BFS(0)
		}
	})

	b.Run("Dijkstra", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			graph.Dijkstra(0)
		}
	})
}
