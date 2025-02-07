# Go API Design

This module demonstrates how to design and implement a RESTful API in Go using the `gorilla/mux` router, following REST principles and best practices.

## Code Explanation

The `main.go` file implements a complete REST API for product management with the following components:

### 1. Data Models
```go
type Product struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    CreatedAt   time.Time `json:"created_at"`
}

type APIResponse struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```
- Clean data structures with JSON tags
- Standardized API response format
- Proper type definitions for all fields

### 2. In-Memory Store
```go
type ProductStore struct {
    sync.RWMutex
    products map[string]Product
}
```
- Thread-safe implementation using `sync.RWMutex`
- Efficient map-based storage
- Demonstrates concurrent access patterns

### 3. API Structure
```go
type API struct {
    store *ProductStore
}
```
- Clean separation of concerns
- Dependency injection pattern
- Extensible design

### 4. REST Endpoints

#### Create Product (POST /products)
```go
func (api *API) CreateProduct(w http.ResponseWriter, r *http.Request)
```
- Validates request payload
- Generates unique ID
- Thread-safe storage
- Returns created product

#### Get Product (GET /products/{id})
```go
func (api *API) GetProduct(w http.ResponseWriter, r *http.Request)
```
- URL parameter handling
- Not found error handling
- Concurrent read access
- JSON response

#### List Products (GET /products)
```go
func (api *API) ListProducts(w http.ResponseWriter, r *http.Request)
```
- Returns all products
- Efficient data retrieval
- Consistent response format

#### Update Product (PUT /products/{id})
```go
func (api *API) UpdateProduct(w http.ResponseWriter, r *http.Request)
```
- Full product update
- Existence validation
- Concurrent write access
- Error handling

#### Delete Product (DELETE /products/{id})
```go
func (api *API) DeleteProduct(w http.ResponseWriter, r *http.Request)
```
- Resource removal
- Status code handling
- Thread-safe deletion

### 5. Helper Functions
```go
func respondWithJSON(w http.ResponseWriter, code int, payload interface{})
```
- Centralizes JSON response handling
- Sets proper headers
- Handles HTTP status codes

## Best Practices Demonstrated
1. RESTful endpoint design
2. Proper HTTP method usage
3. Consistent error handling
4. Thread-safe operations
5. Clean code organization
6. Standard response format

## API Endpoints
- `POST /products`: Create a new product
- `GET /products/{id}`: Retrieve a specific product
- `GET /products`: List all products
- `PUT /products/{id}`: Update a product
- `DELETE /products/{id}`: Delete a product

## Running the API
```bash
# Install dependencies
go get github.com/gorilla/mux

# Run the server
go run main.go
```

## Testing the API
```bash
# Create a product
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Product 1","description":"Description","price":29.99}'

# Get a product
curl http://localhost:8080/products/{id}

# List all products
curl http://localhost:8080/products

# Update a product
curl -X PUT http://localhost:8080/products/{id} \
  -H "Content-Type: application/json" \
  -d '{"name":"Updated Product","description":"New Description","price":39.99}'

# Delete a product
curl -X DELETE http://localhost:8080/products/{id}
```
