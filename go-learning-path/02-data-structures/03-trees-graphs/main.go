// Package main demonstrates trees and graphs in Go
package main

import (
	"fmt"
	"math"
)

// BinaryNode represents a node in a binary tree
type BinaryNode[T any] struct {
	Value       T
	Left, Right *BinaryNode[T]
}

// BinarySearchTree represents a binary search tree
type BinarySearchTree struct {
	Root *BinaryNode[int]
}

// NewBST creates a new binary search tree
func NewBST() *BinarySearchTree {
	return &BinarySearchTree{}
}

// Insert adds a new value to the BST
func (bst *BinarySearchTree) Insert(value int) {
	if bst.Root == nil {
		bst.Root = &BinaryNode[int]{Value: value}
		return
	}
	bst.insertNode(bst.Root, value)
}

func (bst *BinarySearchTree) insertNode(node *BinaryNode[int], value int) {
	if value <= node.Value {
		if node.Left == nil {
			node.Left = &BinaryNode[int]{Value: value}
		} else {
			bst.insertNode(node.Left, value)
		}
	} else {
		if node.Right == nil {
			node.Right = &BinaryNode[int]{Value: value}
		} else {
			bst.insertNode(node.Right, value)
		}
	}
}

// Search looks for a value in the BST
func (bst *BinarySearchTree) Search(value int) bool {
	return bst.searchNode(bst.Root, value)
}

func (bst *BinarySearchTree) searchNode(node *BinaryNode[int], value int) bool {
	if node == nil {
		return false
	}
	if node.Value == value {
		return true
	}
	if value < node.Value {
		return bst.searchNode(node.Left, value)
	}
	return bst.searchNode(node.Right, value)
}

// InOrderTraversal performs in-order traversal of the BST
func (bst *BinarySearchTree) InOrderTraversal() []int {
	var result []int
	bst.inOrder(bst.Root, &result)
	return result
}

func (bst *BinarySearchTree) inOrder(node *BinaryNode[int], result *[]int) {
	if node != nil {
		bst.inOrder(node.Left, result)
		*result = append(*result, node.Value)
		bst.inOrder(node.Right, result)
	}
}

// AVLNode represents a node in an AVL tree
type AVLNode struct {
	Value       int
	Left, Right *AVLNode
	Height      int
}

// AVLTree represents an AVL tree
type AVLTree struct {
	Root *AVLNode
}

// NewAVL creates a new AVL tree
func NewAVL() *AVLTree {
	return &AVLTree{}
}

// Height returns the height of a node
func (n *AVLNode) height() int {
	if n == nil {
		return 0
	}
	return n.Height
}

// UpdateHeight updates the height of a node
func (n *AVLNode) updateHeight() {
	n.Height = 1 + max(getHeight(n.Left), getHeight(n.Right))
}

func getHeight(n *AVLNode) int {
	if n == nil {
		return 0
	}
	return n.Height
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// GetBalance returns the balance factor of a node
func (n *AVLNode) getBalance() int {
	if n == nil {
		return 0
	}
	return getHeight(n.Left) - getHeight(n.Right)
}

// RotateRight performs a right rotation
func (avl *AVLTree) rotateRight(y *AVLNode) *AVLNode {
	x := y.Left
	T2 := x.Right

	x.Right = y
	y.Left = T2

	y.updateHeight()
	x.updateHeight()

	return x
}

// RotateLeft performs a left rotation
func (avl *AVLTree) rotateLeft(x *AVLNode) *AVLNode {
	y := x.Right
	T2 := y.Left

	y.Left = x
	x.Right = T2

	x.updateHeight()
	y.updateHeight()

	return y
}

// Insert adds a new value to the AVL tree
func (avl *AVLTree) Insert(value int) {
	avl.Root = avl.insertNode(avl.Root, value)
}

func (avl *AVLTree) insertNode(node *AVLNode, value int) *AVLNode {
	if node == nil {
		return &AVLNode{Value: value, Height: 1}
	}

	if value < node.Value {
		node.Left = avl.insertNode(node.Left, value)
	} else if value > node.Value {
		node.Right = avl.insertNode(node.Right, value)
	} else {
		return node // Duplicate values not allowed
	}

	node.updateHeight()

	balance := node.getBalance()

	// Left Left Case
	if balance > 1 && value < node.Left.Value {
		return avl.rotateRight(node)
	}

	// Right Right Case
	if balance < -1 && value > node.Right.Value {
		return avl.rotateLeft(node)
	}

	// Left Right Case
	if balance > 1 && value > node.Left.Value {
		node.Left = avl.rotateLeft(node.Left)
		return avl.rotateRight(node)
	}

	// Right Left Case
	if balance < -1 && value < node.Right.Value {
		node.Right = avl.rotateRight(node.Right)
		return avl.rotateLeft(node)
	}

	return node
}

// Graph represents a weighted directed graph
type Graph struct {
	vertices map[int]map[int]int // adjacency list with weights
}

// NewGraph creates a new graph
func NewGraph() *Graph {
	return &Graph{
		vertices: make(map[int]map[int]int),
	}
}

// AddVertex adds a new vertex to the graph
func (g *Graph) AddVertex(vertex int) {
	if _, exists := g.vertices[vertex]; !exists {
		g.vertices[vertex] = make(map[int]int)
	}
}

// AddEdge adds a weighted edge to the graph
func (g *Graph) AddEdge(from, to, weight int) {
	g.AddVertex(from)
	g.AddVertex(to)
	g.vertices[from][to] = weight
}

// DFS performs depth-first search starting from a vertex
func (g *Graph) DFS(start int) []int {
	visited := make(map[int]bool)
	var result []int
	g.dfsUtil(start, visited, &result)
	return result
}

func (g *Graph) dfsUtil(vertex int, visited map[int]bool, result *[]int) {
	visited[vertex] = true
	*result = append(*result, vertex)

	for neighbor := range g.vertices[vertex] {
		if !visited[neighbor] {
			g.dfsUtil(neighbor, visited, result)
		}
	}
}

// BFS performs breadth-first search starting from a vertex
func (g *Graph) BFS(start int) []int {
	visited := make(map[int]bool)
	var result []int
	queue := []int{start}
	visited[start] = true

	for len(queue) > 0 {
		vertex := queue[0]
		queue = queue[1:]
		result = append(result, vertex)

		for neighbor := range g.vertices[vertex] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return result
}

// Dijkstra finds shortest paths from start vertex to all other vertices
func (g *Graph) Dijkstra(start int) map[int]int {
	distances := make(map[int]int)
	visited := make(map[int]bool)

	// Initialize distances
	for vertex := range g.vertices {
		distances[vertex] = math.MaxInt32
	}
	distances[start] = 0

	for len(visited) < len(g.vertices) {
		// Find vertex with minimum distance
		u := -1
		minDist := math.MaxInt32
		for vertex := range g.vertices {
			if !visited[vertex] && distances[vertex] < minDist {
				u = vertex
				minDist = distances[vertex]
			}
		}

		if u == -1 {
			break
		}

		visited[u] = true

		// Update distances to neighbors
		for v, weight := range g.vertices[u] {
			if !visited[v] {
				newDist := distances[u] + weight
				if newDist < distances[v] {
					distances[v] = newDist
				}
			}
		}
	}

	return distances
}

func main() {
	// Demonstrate BST operations
	bst := NewBST()
	values := []int{5, 3, 7, 2, 4, 6, 8}
	fmt.Println("Building BST with values:", values)
	for _, v := range values {
		bst.Insert(v)
	}

	fmt.Println("BST In-order traversal:", bst.InOrderTraversal())
	fmt.Println("Search for 4:", bst.Search(4))
	fmt.Println("Search for 9:", bst.Search(9))

	// Demonstrate AVL Tree operations
	avl := NewAVL()
	fmt.Println("\nBuilding AVL tree with values:", values)
	for _, v := range values {
		avl.Insert(v)
	}

	// Demonstrate Graph operations
	graph := NewGraph()
	fmt.Println("\nBuilding weighted graph...")
	graph.AddEdge(0, 1, 4)
	graph.AddEdge(0, 2, 2)
	graph.AddEdge(1, 2, 1)
	graph.AddEdge(1, 3, 5)
	graph.AddEdge(2, 3, 8)
	graph.AddEdge(2, 4, 10)
	graph.AddEdge(3, 4, 2)

	fmt.Println("DFS starting from vertex 0:", graph.DFS(0))
	fmt.Println("BFS starting from vertex 0:", graph.BFS(0))
	
	distances := graph.Dijkstra(0)
	fmt.Println("Shortest distances from vertex 0:", distances)
}
