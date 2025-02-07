# Go Web Basics

This module demonstrates fundamental concepts of web development in Go using the standard `net/http` package.

## Code Explanation

The `main.go` file implements a basic HTTP server with the following components:

### 1. Response Structure
```go
type Response struct {
    Status  string      `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```
- Defines a standard API response structure
- Uses JSON tags for proper serialization
- `omitempty` tag skips empty Data field in response

### 2. Middleware Implementation
```go
func loggingMiddleware(next http.Handler) http.Handler
```
- Demonstrates the middleware pattern in Go
- Logs request start and completion times
- Shows how to chain HTTP handlers
- Uses `time.Since()` for duration measurement

### 3. JSON Response Helper
```go
func respondWithJSON(w http.ResponseWriter, code int, payload interface{})
```
- Centralizes JSON response handling
- Sets proper Content-Type header
- Handles HTTP status codes
- Marshals payload to JSON

### 4. Route Handlers
```go
func HomeHandler(w http.ResponseWriter, r *http.Request)
func HealthCheckHandler(w http.ResponseWriter, r *http.Request)
```
- `HomeHandler`: Basic route handling with path validation
- `HealthCheckHandler`: System health status endpoint
- Demonstrates proper error handling
- Shows how to structure API responses

### 5. Server Configuration
```go
server := &http.Server{
    Addr:         ":8080",
    Handler:      handler,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
}
```
- Configures server timeouts
- Sets up middleware chain
- Uses proper error handling for server startup

## Best Practices Demonstrated
1. Structured API responses
2. Middleware for cross-cutting concerns
3. Proper timeout configuration
4. Error handling in HTTP context
5. Clean route organization
6. JSON response handling

## API Endpoints
1. `GET /`: Home endpoint
   - Returns welcome message
   - Handles 404 for invalid paths
   
2. `GET /health`: Health check endpoint
   - Returns service status
   - Includes version and uptime information

## Running the Server
```bash
go run main.go
```
The server will start on `http://localhost:8080`

## Testing the API
```bash
# Test home endpoint
curl http://localhost:8080/

# Test health endpoint
curl http://localhost:8080/health
```
