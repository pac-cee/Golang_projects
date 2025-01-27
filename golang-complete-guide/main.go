package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"
	"net/http"
	"github.com/redis/go-redis/v9"
)

// BASIC CONCEPTS
// 1. Variables and Data Types
func basicTypes() {
	// Integer types
	var i int = 42
	var i8 int8 = 127
	var ui uint = 123

	// Float types
	var f32 float32 = 3.14
	var f64 float64 = 3.141592653589793

	// Boolean
	var isTrue bool = true

	// String
	var str string = "Hello, Go!"

	// Type inference
	auto := "Type inferred automatically"

	fmt.Printf("Types: %T %T %T %T %T %T %T %T\n", i, i8, ui, f32, f64, isTrue, str, auto)
}

// 2. Constants
const (
	Pi       = 3.14159
	Username = "admin"
	// iota example
	Sunday = iota
	Monday
	Tuesday
)

// 3. Arrays and Slices
func arraysAndSlices() {
	// Array
	var arr [5]int = [5]int{1, 2, 3, 4, 5}

	// Slice
	slice := []int{1, 2, 3}
	slice = append(slice, 4)

	// Slice operations
	subSlice := slice[1:3]
	
	fmt.Println(arr, slice, subSlice)
}

// 4. Maps
func maps() {
	// Declaration and initialization
	scores := map[string]int{
		"Alice": 95,
		"Bob":   89,
	}

	// Adding and accessing
	scores["Charlie"] = 85
	fmt.Println(scores["Alice"])

	// Check if key exists
	if score, exists := scores["David"]; exists {
		fmt.Println(score)
	}
}

// INTERMEDIATE CONCEPTS
// 1. Structs and Methods
type Person struct {
	Name string
	Age  int
}

// Method
func (p Person) Greet() string {
	return fmt.Sprintf("Hello, my name is %s and I'm %d years old", p.Name, p.Age)
}

// 2. Interfaces
type Animal interface {
	Speak() string
}

type Dog struct {
	Name string
}

func (d Dog) Speak() string {
	return fmt.Sprintf("%s says Woof!", d.Name)
}

// 3. Error Handling
type CustomError struct {
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, &CustomError{Message: "division by zero"}
	}
	return a / b, nil
}

// 4. Goroutines and Channels
func goroutineExample() {
	ch := make(chan string)
	
	go func() {
		ch <- "Message from goroutine"
	}()

	msg := <-ch
	fmt.Println(msg)
}

// ADVANCED CONCEPTS
// 1. Reflection
func reflectionExample() {
	p := Person{Name: "John", Age: 30}
	t := reflect.TypeOf(p)
	v := reflect.ValueOf(p)

	fmt.Printf("Type: %v\n", t)
	fmt.Printf("Fields: %v\n", v.Field(0))
}

// 2. Context
func contextExample() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	select {
	case <-time.After(3 * time.Second):
		fmt.Println("Overslept")
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	}
}

// 3. Concurrency Patterns
// Worker Pool
func workerPool() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Start workers
	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Send jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	close(results)
}

func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("worker %d processing job %d\n", id, j)
		time.Sleep(time.Second)
		results <- j * 2
	}
}

// 4. Generics (Go 1.18+)
type Number interface {
	~int | ~float64
}

func Min[T Number](x, y T) T {
	if x < y {
		return x
	}
	return y
}

// 5. Thread-safe Counter with Mutex
type SafeCounter struct {
	mu sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) GetCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

// 6. Basic HTTP Server
func startHTTPServer() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Web!")
	})
	
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()
}

// 7. JSON Handling
type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Features []string `json:"features"`
}

func jsonExample() {
	config := Config{
		Name:     "MyApp",
		Version:  "1.0.0",
		Features: []string{"auth", "logging", "api"},
	}
	
	// Marshal (encode) to JSON
	jsonData, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("JSON marshaling error: %v\n", err)
		return
	}
	
	// Unmarshal (decode) from JSON
	var decodedConfig Config
	err = json.Unmarshal(jsonData, &decodedConfig)
	if err != nil {
		fmt.Printf("JSON unmarshaling error: %v\n", err)
		return
	}
	
	fmt.Printf("Decoded config: %+v\n", decodedConfig)
}

// 8. Panic Recovery
func recoverExample() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()
	
	// Simulating a panic
	panic("something went wrong")
}

// Redis Client wrapper
type RedisService struct {
	client *redis.Client
}

// Create new Redis service
func NewRedisService() (*RedisService, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // no password set
		DB:       0,                // use default DB
	})

	// Test the connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis connection failed: %v", err)
	}

	return &RedisService{client: client}, nil
}

// Basic String Operations
func (rs *RedisService) StringOperations(ctx context.Context) error {
	// SET operation
	err := rs.client.Set(ctx, "user:1:name", "John Doe", 0).Err()
	if err != nil {
		return fmt.Errorf("SET failed: %v", err)
	}

	// GET operation
	name, err := rs.client.Get(ctx, "user:1:name").Result()
	if err != nil {
		return fmt.Errorf("GET failed: %v", err)
	}
	fmt.Printf("user:1:name = %v\n", name)

	// SETEX (SET with expiration)
	err = rs.client.SetEx(ctx, "session:123", "active", 1*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("SETEX failed: %v", err)
	}

	// TTL (time to live)
	ttl, err := rs.client.TTL(ctx, "session:123").Result()
	if err != nil {
		return fmt.Errorf("TTL failed: %v", err)
	}
	fmt.Printf("session:123 TTL = %v\n", ttl)

	return nil
}

// Hash Operations
func (rs *RedisService) HashOperations(ctx context.Context) error {
	// HSET (set multiple hash fields)
	err := rs.client.HSet(ctx, "user:1", map[string]interface{}{
		"name":  "John Doe",
		"email": "john@example.com",
		"age":   30,
	}).Err()
	if err != nil {
		return fmt.Errorf("HSET failed: %v", err)
	}

	// HGET (get specific hash field)
	email, err := rs.client.HGet(ctx, "user:1", "email").Result()
	if err != nil {
		return fmt.Errorf("HGET failed: %v", err)
	}
	fmt.Printf("user:1 email = %v\n", email)

	// HGETALL (get all hash fields)
	userData, err := rs.client.HGetAll(ctx, "user:1").Result()
	if err != nil {
		return fmt.Errorf("HGETALL failed: %v", err)
	}
	fmt.Printf("user:1 data = %v\n", userData)

	return nil
}

// List Operations
func (rs *RedisService) ListOperations(ctx context.Context) error {
	// LPUSH (push to left/start of list)
	err := rs.client.LPush(ctx, "tasks", "task1", "task2", "task3").Err()
	if err != nil {
		return fmt.Errorf("LPUSH failed: %v", err)
	}

	// RPOP (pop from right/end of list)
	task, err := rs.client.RPop(ctx, "tasks").Result()
	if err != nil {
		return fmt.Errorf("RPOP failed: %v", err)
	}
	fmt.Printf("Popped task = %v\n", task)

	// LRANGE (get range of list elements)
	tasks, err := rs.client.LRange(ctx, "tasks", 0, -1).Result()
	if err != nil {
		return fmt.Errorf("LRANGE failed: %v", err)
	}
	fmt.Printf("All tasks = %v\n", tasks)

	return nil
}

// Set Operations
func (rs *RedisService) SetOperations(ctx context.Context) error {
	// SADD (add members to set)
	err := rs.client.SAdd(ctx, "online_users", "user1", "user2", "user3").Err()
	if err != nil {
		return fmt.Errorf("SADD failed: %v", err)
	}

	// SISMEMBER (check if member exists in set)
	exists, err := rs.client.SIsMember(ctx, "online_users", "user1").Result()
	if err != nil {
		return fmt.Errorf("SISMEMBER failed: %v", err)
	}
	fmt.Printf("Is user1 online? %v\n", exists)

	// SMEMBERS (get all set members)
	users, err := rs.client.SMembers(ctx, "online_users").Result()
	if err != nil {
		return fmt.Errorf("SMEMBERS failed: %v", err)
	}
	fmt.Printf("Online users = %v\n", users)

	return nil
}

// Sorted Set Operations
func (rs *RedisService) SortedSetOperations(ctx context.Context) error {
	// ZADD (add members with scores)
	users := []redis.Z{
		{Score: 100, Member: "user1"},
		{Score: 200, Member: "user2"},
		{Score: 150, Member: "user3"},
	}
	err := rs.client.ZAdd(ctx, "scores", users...).Err()
	if err != nil {
		return fmt.Errorf("ZADD failed: %v", err)
	}

	// ZRANGE (get range by index with scores)
	scores, err := rs.client.ZRangeWithScores(ctx, "scores", 0, -1).Result()
	if err != nil {
		return fmt.Errorf("ZRANGE failed: %v", err)
	}
	fmt.Printf("User scores = %v\n", scores)

	return nil
}

// Pipeline Example (execute multiple commands in one round trip)
func (rs *RedisService) PipelineExample(ctx context.Context) error {
	pipe := rs.client.Pipeline()
	
	// Queue multiple commands
	incr := pipe.Incr(ctx, "counter")
	pipe.Expire(ctx, "counter", time.Hour)
	
	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("pipeline execution failed: %v", err)
	}
	
	fmt.Printf("Counter value = %v\n", incr.Val())
	return nil
}

func main() {
	fmt.Println("=== Basic Concepts ===")
	basicTypes()
	arraysAndSlices()
	maps()

	fmt.Println("\n=== Intermediate Concepts ===")
	p := Person{Name: "Alice", Age: 25}
	fmt.Println(p.Greet())

	d := Dog{Name: "Rex"}
	fmt.Println(d.Speak())

	result, err := divide(10, 2)
	fmt.Printf("Division result: %v, error: %v\n", result, err)

	fmt.Println("\n=== Advanced Concepts ===")
	reflectionExample()
	contextExample()
	workerPool()

	// Generics example
	fmt.Printf("Min of 10 and 5: %v\n", Min(10, 5))
	fmt.Printf("Min of 3.14 and 2.71: %v\n", Min(3.14, 2.71))
	
	fmt.Println("\n=== Additional Advanced Concepts ===")
	// SafeCounter example
	counter := &SafeCounter{}
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}
	wg.Wait()
	fmt.Printf("Final count: %d\n", counter.GetCount())
	
	// Start HTTP server
	startHTTPServer()
	fmt.Println("HTTP server started on :8080")
	
	// JSON example
	jsonExample()
	
	// Panic recovery example
	fmt.Println("\nTesting panic recovery:")
	recoverExample()

	// Create Redis service
	redisService, err := NewRedisService()
	if err != nil {
		fmt.Printf("Failed to create Redis service: %v\n", err)
		return
	}

	ctx := context.Background()

	fmt.Println("\n=== String Operations ===")
	if err := redisService.StringOperations(ctx); err != nil {
		fmt.Printf("String operations failed: %v\n", err)
	}

	fmt.Println("\n=== Hash Operations ===")
	if err := redisService.HashOperations(ctx); err != nil {
		fmt.Printf("Hash operations failed: %v\n", err)
	}

	fmt.Println("\n=== List Operations ===")
	if err := redisService.ListOperations(ctx); err != nil {
		fmt.Printf("List operations failed: %v\n", err)
	}

	fmt.Println("\n=== Set Operations ===")
	if err := redisService.SetOperations(ctx); err != nil {
		fmt.Printf("Set operations failed: %v\n", err)
	}

	fmt.Println("\n=== Sorted Set Operations ===")
	if err := redisService.SortedSetOperations(ctx); err != nil {
		fmt.Printf("Sorted set operations failed: %v\n", err)
	}

	fmt.Println("\n=== Pipeline Example ===")
	if err := redisService.PipelineExample(ctx); err != nil {
		fmt.Printf("Pipeline example failed: %v\n", err)
	}
}
