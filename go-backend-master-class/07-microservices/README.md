# Go Microservices Implementation

This module demonstrates how to implement a microservices architecture in Go, including service discovery, health checks, and circuit breaker patterns.

## Code Explanation

The `main.go` file implements a complete microservices system with the following components:

### 1. Service Model
```go
type Service struct {
    Name     string    `json:"name"`
    URL      string    `json:"url"`
    Status   string    `json:"status"`
    LastPing time.Time `json:"last_ping"`
}

type ServiceRegistry struct {
    Services map[string]Service
}
```
- Service definition
- Registry for service discovery
- Health status tracking
- Last ping monitoring

### 2. Service Registry
```go
func (sr *ServiceRegistry) RegisterService(w http.ResponseWriter, r *http.Request)
func (sr *ServiceRegistry) GetServices(w http.ResponseWriter, r *http.Request)
func (sr *ServiceRegistry) HealthCheck(w http.ResponseWriter, r *http.Request)
```
- Service registration
- Service discovery
- Health monitoring
- Status updates

### 3. Service Client
```go
type ServiceClient struct {
    registry *ServiceRegistry
    client   *http.Client
}
```
- Inter-service communication
- HTTP client management
- Service discovery integration
- Error handling

### 4. Circuit Breaker
```go
type CircuitBreaker struct {
    failures  int
    threshold int
    timeout   time.Duration
    lastError time.Time
}
```
- Failure threshold tracking
- Timeout management
- Error tracking
- Circuit state management

### 5. Circuit Breaker Implementation
```go
func (cb *CircuitBreaker) Execute(fn func() error) error
```
- Failure counting
- Circuit state transitions
- Timeout handling
- Error propagation

## Microservices Patterns

### 1. Service Discovery
- Dynamic service registration
- Health-based service filtering
- Last ping tracking
- Service status updates

### 2. Health Checks
- Regular health monitoring
- Status updates
- Ping timestamp tracking
- Failure detection

### 3. Circuit Breaker
- Prevents cascade failures
- Automatic recovery
- Configurable thresholds
- Timeout management

### 4. Inter-Service Communication
- HTTP-based communication
- Error handling
- Retry logic
- Circuit breaker protection

## API Endpoints

### Service Registry
- `POST /services`: Register a new service
  ```json
  {
    "name": "user-service",
    "url": "http://localhost:8081"
  }
  ```

- `GET /services`: List all services
- `GET /services/{name}/health`: Check service health

## Best Practices Demonstrated
1. Service isolation
2. Failure handling
3. Health monitoring
4. Circuit breaking
5. Service discovery
6. Error propagation

## Running the Services
```bash
# Install dependencies
go get github.com/gorilla/mux
go get google.golang.org/grpc

# Run the service registry
go run main.go

# Environment Variables
export SERVICE_PORT=8080      # Service registry port
export SERVICE_NAME=registry  # Service name
```

## Service Configuration
1. Circuit Breaker
   ```go
   breaker := NewCircuitBreaker(
       3,                  // Failure threshold
       time.Second * 30,   // Reset timeout
   )
   ```

2. Service Client
   ```go
   client := NewServiceClient(registry)
   ```

## Testing Services
```bash
# Register a service
curl -X POST http://localhost:8080/services \
  -H "Content-Type: application/json" \
  -d '{"name":"user-service","url":"http://localhost:8081"}'

# List services
curl http://localhost:8080/services

# Check service health
curl http://localhost:8080/services/user-service/health
```

## Monitoring and Maintenance
1. Health Check Configuration
   - Regular interval: 30 seconds
   - Timeout: 5 seconds
   - Failure threshold: 3 attempts

2. Circuit Breaker States
   - Closed: Normal operation
   - Open: Failing, no requests allowed
   - Half-Open: Testing recovery

3. Service States
   - Healthy: Responding normally
   - Unhealthy: Failed health checks
   - Unknown: No recent health data

## Error Handling
1. Service not found
2. Communication failures
3. Timeout errors
4. Circuit breaker trips
5. Health check failures

## Scalability Considerations
1. Service registration capacity
2. Health check frequency
3. Circuit breaker thresholds
4. Network timeout values
5. Recovery strategies
