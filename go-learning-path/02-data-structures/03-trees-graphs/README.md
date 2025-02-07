# Trees and Graphs in Go 🌳

## 📖 Table of Contents
1. [Trees](#trees)
2. [Binary Trees](#binary-trees)
3. [Binary Search Trees](#binary-search-trees)
4. [Graphs](#graphs)
5. [Common Algorithms](#common-algorithms)
6. [Best Practices](#best-practices)
7. [Exercises](#exercises)

## Trees

### Tree Node
```go
type Node struct {
    Value    interface{}
    Children []*Node
}

// Generic tree node
type TreeNode[T any] struct {
    Value    T
    Children []*TreeNode[T]
}
```

### Tree Operations
```go
// Add child
func (n *Node) AddChild(value interface{}) *Node {
    child := &Node{Value: value}
    n.Children = append(n.Children, child)
    return child
}

// Remove child
func (n *Node) RemoveChild(child *Node) {
    for i, c := range n.Children {
        if c == child {
            n.Children = append(n.Children[:i], n.Children[i+1:]...)
            return
        }
    }
}
```

## Binary Trees

### Binary Tree Node
```go
type BinaryNode struct {
    Value       interface{}
    Left, Right *BinaryNode
}

// Generic binary node
type BinaryNode[T any] struct {
    Value       T
    Left, Right *BinaryNode[T]
}
```

### Traversal Methods
```go
// Pre-order traversal
func PreOrder(node *BinaryNode) {
    if node == nil {
        return
    }
    fmt.Println(node.Value)
    PreOrder(node.Left)
    PreOrder(node.Right)
}

// In-order traversal
func InOrder(node *BinaryNode) {
    if node == nil {
        return
    }
    InOrder(node.Left)
    fmt.Println(node.Value)
    InOrder(node.Right)
}

// Post-order traversal
func PostOrder(node *BinaryNode) {
    if node == nil {
        return
    }
    PostOrder(node.Left)
    PostOrder(node.Right)
    fmt.Println(node.Value)
}

// Level-order traversal
func LevelOrder(root *BinaryNode) {
    if root == nil {
        return
    }
    queue := []*BinaryNode{root}
    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        fmt.Println(node.Value)
        if node.Left != nil {
            queue = append(queue, node.Left)
        }
        if node.Right != nil {
            queue = append(queue, node.Right)
        }
    }
}
```

## Binary Search Trees

### BST Node
```go
type BSTNode struct {
    Value       int
    Left, Right *BSTNode
}

// Insert value
func (n *BSTNode) Insert(value int) {
    if value <= n.Value {
        if n.Left == nil {
            n.Left = &BSTNode{Value: value}
        } else {
            n.Left.Insert(value)
        }
    } else {
        if n.Right == nil {
            n.Right = &BSTNode{Value: value}
        } else {
            n.Right.Insert(value)
        }
    }
}

// Search value
func (n *BSTNode) Search(value int) bool {
    if n == nil {
        return false
    }
    if n.Value == value {
        return true
    }
    if value < n.Value {
        return n.Left.Search(value)
    }
    return n.Right.Search(value)
}
```

### BST Operations
```go
// Delete value
func (n *BSTNode) Delete(value int) *BSTNode {
    if n == nil {
        return nil
    }
    
    if value < n.Value {
        n.Left = n.Left.Delete(value)
    } else if value > n.Value {
        n.Right = n.Right.Delete(value)
    } else {
        // Node to delete found
        if n.Left == nil {
            return n.Right
        }
        if n.Right == nil {
            return n.Left
        }
        
        // Node has two children
        minRight := n.Right.findMin()
        n.Value = minRight
        n.Right = n.Right.Delete(minRight)
    }
    return n
}

// Find minimum value
func (n *BSTNode) findMin() int {
    current := n
    for current.Left != nil {
        current = current.Left
    }
    return current.Value
}
```

## Graphs

### Graph Types
```go
// Adjacency List
type Graph struct {
    vertices map[int][]int
}

// Adjacency Matrix
type GraphMatrix struct {
    vertices int
    matrix   [][]bool
}

// Weighted Graph
type WeightedGraph struct {
    vertices map[int]map[int]int
}
```

### Graph Operations
```go
// Add edge (Adjacency List)
func (g *Graph) AddEdge(from, to int) {
    g.vertices[from] = append(g.vertices[from], to)
}

// Add edge (Adjacency Matrix)
func (g *GraphMatrix) AddEdge(from, to int) {
    g.matrix[from][to] = true
}

// Add weighted edge
func (g *WeightedGraph) AddEdge(from, to, weight int) {
    if g.vertices[from] == nil {
        g.vertices[from] = make(map[int]int)
    }
    g.vertices[from][to] = weight
}
```

## Common Algorithms

### Depth-First Search (DFS)
```go
func (g *Graph) DFS(start int, visited map[int]bool) {
    if visited == nil {
        visited = make(map[int]bool)
    }
    
    visited[start] = true
    fmt.Println(start)
    
    for _, neighbor := range g.vertices[start] {
        if !visited[neighbor] {
            g.DFS(neighbor, visited)
        }
    }
}
```

### Breadth-First Search (BFS)
```go
func (g *Graph) BFS(start int) {
    visited := make(map[int]bool)
    queue := []int{start}
    visited[start] = true
    
    for len(queue) > 0 {
        vertex := queue[0]
        queue = queue[1:]
        fmt.Println(vertex)
        
        for _, neighbor := range g.vertices[vertex] {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
}
```

### Dijkstra's Algorithm
```go
func (g *WeightedGraph) Dijkstra(start int) map[int]int {
    distances := make(map[int]int)
    visited := make(map[int]bool)
    
    // Initialize distances
    for vertex := range g.vertices {
        distances[vertex] = math.MaxInt32
    }
    distances[start] = 0
    
    for len(visited) < len(g.vertices) {
        // Find minimum distance vertex
        u := minDistance(distances, visited)
        visited[u] = true
        
        // Update distances to neighbors
        for v, weight := range g.vertices[u] {
            if !visited[v] && 
               distances[u] != math.MaxInt32 &&
               distances[u]+weight < distances[v] {
                distances[v] = distances[u] + weight
            }
        }
    }
    
    return distances
}
```

## Best Practices

### 1. Tree Implementation
```go
// Use interfaces for flexibility
type Tree interface {
    Insert(value interface{})
    Delete(value interface{})
    Search(value interface{}) bool
}

// Use generics for type safety
type BinaryTree[T comparable] struct {
    Root *BinaryNode[T]
}
```

### 2. Graph Implementation
```go
// Use appropriate representation
type Graph interface {
    AddEdge(from, to int)
    RemoveEdge(from, to int)
    HasEdge(from, to int) bool
    Neighbors(vertex int) []int
}

// Consider memory vs performance
type SparseGraph struct {
    vertices map[int][]int    // Good for sparse graphs
}

type DenseGraph struct {
    matrix [][]bool           // Good for dense graphs
}
```

### 3. Algorithm Implementation
```go
// Use efficient data structures
type PriorityQueue struct {
    items    []int
    priority map[int]int
}

// Consider time/space complexity
func (g *Graph) DFSIterative(start int) {
    stack := []int{start}
    visited := make(map[int]bool)
    
    for len(stack) > 0 {
        vertex := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        
        if !visited[vertex] {
            visited[vertex] = true
            fmt.Println(vertex)
            
            // Add neighbors in reverse order
            neighbors := g.vertices[vertex]
            for i := len(neighbors)-1; i >= 0; i-- {
                if !visited[neighbors[i]] {
                    stack = append(stack, neighbors[i])
                }
            }
        }
    }
}
```

## Exercises

### Exercise 1: Binary Tree Height
```go
func (n *BinaryNode) Height() int {
    if n == nil {
        return 0
    }
    leftHeight := n.Left.Height()
    rightHeight := n.Right.Height()
    if leftHeight > rightHeight {
        return leftHeight + 1
    }
    return rightHeight + 1
}
```

### Exercise 2: Check BST
```go
func (n *BSTNode) IsBST(min, max int) bool {
    if n == nil {
        return true
    }
    if n.Value <= min || n.Value >= max {
        return false
    }
    return n.Left.IsBST(min, n.Value) &&
           n.Right.IsBST(n.Value, max)
}
```

### Exercise 3: Graph Cycle Detection
```go
func (g *Graph) HasCycle() bool {
    visited := make(map[int]bool)
    recStack := make(map[int]bool)
    
    for vertex := range g.vertices {
        if !visited[vertex] {
            if g.hasCycleUtil(vertex, visited, recStack) {
                return true
            }
        }
    }
    return false
}

func (g *Graph) hasCycleUtil(vertex int, visited, recStack map[int]bool) bool {
    visited[vertex] = true
    recStack[vertex] = true
    
    for _, neighbor := range g.vertices[vertex] {
        if !visited[neighbor] {
            if g.hasCycleUtil(neighbor, visited, recStack) {
                return true
            }
        } else if recStack[neighbor] {
            return true
        }
    }
    
    recStack[vertex] = false
    return false
}
```

## Common Patterns

### 1. Tree Balancing
```go
// AVL Tree Node
type AVLNode struct {
    Value       int
    Left, Right *AVLNode
    Height      int
}

func (n *AVLNode) balance() {
    balance := n.getBalance()
    if balance > 1 {
        if n.Left.getBalance() < 0 {
            n.Left.rotateLeft()
        }
        n.rotateRight()
    } else if balance < -1 {
        if n.Right.getBalance() > 0 {
            n.Right.rotateRight()
        }
        n.rotateLeft()
    }
}
```

### 2. Graph Traversal Patterns
```go
// Topological Sort
func (g *Graph) TopologicalSort() []int {
    visited := make(map[int]bool)
    stack := make([]int, 0)
    
    for vertex := range g.vertices {
        if !visited[vertex] {
            g.topologicalSortUtil(vertex, visited, &stack)
        }
    }
    
    // Reverse stack for correct order
    for i := 0; i < len(stack)/2; i++ {
        j := len(stack) - 1 - i
        stack[i], stack[j] = stack[j], stack[i]
    }
    
    return stack
}
```

### 3. Path Finding
```go
// A* Algorithm
type PathNode struct {
    vertex    int
    gScore    int // Cost from start
    hScore    int // Estimated cost to goal
    fScore    int // gScore + hScore
    parent    *PathNode
}

func (g *WeightedGraph) AStar(start, goal int, heuristic func(int) int) []int {
    // Implementation of A* algorithm
}
```

## Next Steps
- Practice tree and graph implementations
- Study advanced algorithms
- Implement real-world applications
- Move on to Advanced Data Structures
