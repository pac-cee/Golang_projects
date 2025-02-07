# Go Deployment and Monitoring

This module demonstrates how to deploy a Go application using Docker and Kubernetes, with built-in monitoring using Prometheus metrics.

## Code Explanation

### 1. Application Components

#### Metrics Implementation
```go
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
)
```
- Request counting
- Duration tracking
- Label-based metrics
- Prometheus integration

#### Middleware
```go
func metricsMiddleware(next http.HandlerFunc) http.HandlerFunc
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc
```
- Request metrics
- Response timing
- Logging
- Error tracking

#### Health Checks
```go
func healthHandler(w http.ResponseWriter, r *http.Request)
func readinessHandler(w http.ResponseWriter, r *http.Request)
```
- Liveness probe
- Readiness probe
- Version info
- Uptime tracking

### 2. Docker Configuration

```dockerfile
# Multi-stage build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```
- Multi-stage build
- Minimal final image
- Dependency management
- Port configuration

### 3. Kubernetes Deployment

#### Deployment Configuration
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-backend-app
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: go-backend
        image: go-backend:latest
        resources:
          limits:
            cpu: "500m"
            memory: "512Mi"
```
- Replica management
- Resource limits
- Container configuration
- Health probes

#### Service Configuration
```yaml
apiVersion: v1
kind: Service
metadata:
  name: go-backend-service
spec:
  type: LoadBalancer
  ports:
    - port: 80
      targetPort: 8080
```
- Load balancing
- Port mapping
- Service discovery
- External access

## Deployment Process

### 1. Building the Docker Image
```bash
# Build the image
docker build -t go-backend:latest .

# Run locally
docker run -p 8080:8080 go-backend:latest
```

### 2. Kubernetes Deployment
```bash
# Apply the deployment
kubectl apply -f k8s/deployment.yaml

# Check status
kubectl get deployments
kubectl get pods
kubectl get services
```

### 3. Monitoring Setup
```bash
# Access Prometheus metrics
curl http://localhost:8080/metrics

# Common metrics:
- http_requests_total
- http_request_duration_seconds
- go_goroutines
- go_threads
```

## Best Practices

### 1. Docker Best Practices
- Multi-stage builds
- Minimal base images
- Layer optimization
- Security considerations
- Environment variables

### 2. Kubernetes Best Practices
- Resource limits
- Health checks
- Rolling updates
- Pod anti-affinity
- Service configuration

### 3. Monitoring Best Practices
- Key metrics
- Alert thresholds
- Dashboard setup
- Log aggregation
- Performance tracking

## Configuration Options

### 1. Environment Variables
```bash
APP_VERSION=1.0.0
PORT=8080
LOG_LEVEL=info
```

### 2. Resource Limits
```yaml
resources:
  limits:
    cpu: "500m"
    memory: "512Mi"
  requests:
    cpu: "200m"
    memory: "256Mi"
```

### 3. Health Checks
```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10
```

## Monitoring and Metrics

### 1. Request Metrics
- Total requests
- Request duration
- Status codes
- Endpoint usage

### 2. System Metrics
- CPU usage
- Memory usage
- Goroutine count
- GC statistics

### 3. Custom Metrics
- Business metrics
- Error rates
- Cache hits/misses
- Queue lengths

## Scaling Considerations

### 1. Horizontal Scaling
- Pod replicas
- Load balancing
- Session handling
- Data consistency

### 2. Resource Management
- CPU allocation
- Memory limits
- Network resources
- Storage requirements

### 3. Performance Optimization
- Cache usage
- Connection pooling
- Request batching
- Async processing

## Security Considerations

### 1. Container Security
- Minimal base images
- Non-root users
- Read-only filesystems
- Security scanning

### 2. Network Security
- Service isolation
- Network policies
- TLS configuration
- API authentication

### 3. Resource Protection
- Rate limiting
- Resource quotas
- Pod security policies
- Secret management
