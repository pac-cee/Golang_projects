# Container Orchestration Tool

A simple container orchestration system built in Go, demonstrating container management and monitoring capabilities.

## Features

- Container lifecycle management (create, start, stop, remove)
- Container monitoring (CPU, memory usage)
- REST API interface
- Docker integration
- Real-time container statistics
- Environment variable support
- Port mapping
- Container health checks

## Prerequisites

1. Go 1.19 or later
2. Docker Engine installed and running
3. Docker API access

## API Endpoints

### Create Container
```http
POST /containers
Content-Type: application/json

{
    "name": "my-container",
    "image": "nginx:latest",
    "env": {
        "KEY": "value"
    },
    "ports": {
        "80/tcp": "8080"
    }
}
```

### List Containers
```http
GET /containers
```

### Stop Container
```http
POST /containers/{id}
```

### Remove Container
```http
DELETE /containers/{id}
```

## Implementation Details

### Container Management
- Uses Docker Engine API
- Supports container configuration
- Handles port mappings
- Manages environment variables
- Implements restart policies

### Monitoring
- Real-time CPU usage tracking
- Memory usage monitoring
- Container health status
- Automatic stats collection

### Data Structures

1. **Container**
```go
type Container struct {
    ID           string
    Name         string
    Image        string
    Status       string
    CreatedAt    time.Time
    Health       string
    CPU          float64
    Memory       int64
    RestartCount int
}
```

2. **ContainerRequest**
```go
type ContainerRequest struct {
    Name  string
    Image string
    Env   map[string]string
    Ports map[string]string
}
```

## How to Run

1. Start the orchestrator:
   ```bash
   go run main.go
   ```

2. Create a container:
   ```bash
   curl -X POST http://localhost:8080/containers \
     -H "Content-Type: application/json" \
     -d '{
       "name": "web-server",
       "image": "nginx:latest",
       "ports": {
         "80/tcp": "8080"
       }
     }'
   ```

3. List containers:
   ```bash
   curl http://localhost:8080/containers
   ```

4. Stop a container:
   ```bash
   curl -X POST http://localhost:8080/containers/{container-id}
   ```

5. Remove a container:
   ```bash
   curl -X DELETE http://localhost:8080/containers/{container-id}
   ```

## Architecture

1. **Core Components**
   - Container Manager
   - Monitoring System
   - REST API Server
   - Stats Collector

2. **Docker Integration**
   - Direct Docker Engine API usage
   - Container lifecycle management
   - Resource monitoring

3. **Monitoring System**
   - Periodic stats collection
   - Resource usage tracking
   - Health monitoring

## Security Considerations

1. API Security
   - No authentication (should be added)
   - No HTTPS (should be added)
   - No rate limiting

2. Container Security
   - Default Docker security
   - No resource limits
   - No network isolation

## Performance

1. Resource Usage
   - Lightweight monitoring
   - Periodic stats collection
   - In-memory container tracking

2. Scalability
   - Single node only
   - No clustering support
   - Limited by Docker host

## Next Steps

1. Add authentication and authorization
2. Implement HTTPS support
3. Add resource limits
4. Implement container networking
5. Add volume management
6. Implement clustering
7. Add logging system
8. Implement backup/restore
9. Add container templates
10. Implement rolling updates
11. Add service discovery
12. Implement load balancing
13. Add monitoring dashboard
14. Implement alerts
15. Add container logs API
