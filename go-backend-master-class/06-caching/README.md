# Go Caching Implementation

This module demonstrates how to implement caching in Go using both Redis and in-memory caching strategies.

## Code Explanation

The `main.go` file implements a complete caching system with the following components:

### 1. Data Models
```go
type Product struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Price       float64   `json:"price"`
    LastUpdated time.Time `json:"last_updated"`
}

type Response struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Source  string      `json:"source,omitempty"` // "cache" or "database"
}
```
- Product model with cache-relevant fields
- API response with cache source tracking

### 2. Cache Interface
```go
type Cache interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    Delete(ctx context.Context, key string) error
}
```
- Generic cache interface
- Context-aware operations
- Support for expiration

### 3. Redis Cache Implementation
```go
type RedisCache struct {
    client *redis.Client
}
```
- Redis client integration
- Implements Cache interface
- Handles distributed caching
- Supports TTL

### 4. In-Memory Cache Implementation
```go
type InMemoryCache struct {
    sync.RWMutex
    data map[string]string
}
```
- Thread-safe implementation
- Map-based storage
- Local memory caching
- Lock-based synchronization

### 5. API Integration
```go
type API struct {
    cache Cache
    store map[string]Product
    mutex sync.RWMutex
    ctx   context.Context
}
```
- Cache-aware API design
- Strategy pattern for cache selection
- Thread-safe operations
- Context propagation

### 6. Cache Operations

#### Create/Update Cache
```go
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
```
- Atomic operations
- JSON serialization
- Expiration handling
- Error management

#### Read from Cache
```go
func (c *RedisCache) Get(ctx context.Context, key string) (string, error)
```
- Cache hit/miss handling
- Fallback strategies
- Error handling
- Performance optimization

#### Cache Invalidation
```go
func (c *RedisCache) Delete(ctx context.Context, key string) error
```
- Key removal
- Error handling
- Atomic operations

## Best Practices Demonstrated
1. Interface-based design
2. Thread-safe operations
3. Context propagation
4. Error handling
5. Cache invalidation
6. Performance optimization

## Caching Strategies
1. **Cache-Aside (Lazy Loading)**
   - Check cache first
   - Load from database on miss
   - Update cache with new data

2. **Write-Through**
   - Update cache and database together
   - Ensures consistency
   - Handles race conditions

3. **Cache Invalidation**
   - Delete on update
   - TTL-based expiration
   - Atomic operations

## Running the Application
```bash
# Install dependencies
go get github.com/go-redis/redis/v8
go get github.com/gorilla/mux

# Start Redis (required for Redis cache)
docker run --name redis -p 6379:6379 -d redis

# Run the application
go run main.go
```

## API Endpoints

### Products API
- `POST /products`: Create product
  ```json
  {
    "name": "Product 1",
    "price": 29.99
  }
  ```

- `GET /products/{id}`: Get product
  - Returns cache source in response
  ```json
  {
    "status": "success",
    "data": {
      "id": "123",
      "name": "Product 1",
      "price": 29.99
    },
    "source": "cache"
  }
  ```

- `PUT /products/{id}`: Update product
  - Invalidates cache on update

## Cache Configuration
1. Redis Cache
   - Default address: localhost:6379
   - Configurable expiration
   - Distributed caching support

2. In-Memory Cache
   - Local storage
   - Thread-safe operations
   - No automatic expiration

## Performance Considerations
1. Use appropriate cache strategy
2. Monitor cache hit/miss ratios
3. Set reasonable TTL values
4. Handle cache stampede
5. Implement circuit breakers
6. Monitor memory usage
