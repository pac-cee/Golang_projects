# Go Backend Master Tutorial

A comprehensive backend application demonstrating production-ready practices in Go.

## Features

- RESTful API with Gin framework
- Clean Architecture
- JWT Authentication
- Role-based Authorization
- PostgreSQL Database
- Redis Caching
- Swagger Documentation
- Unit & Integration Testing
- Docker Support
- CI/CD Pipeline
- Logging & Monitoring
- Rate Limiting
- Request Validation
- Error Handling
- Database Migrations
- Environment Configuration
- Graceful Shutdown

## Project Structure

```
07_Backend_Master/
├── api/
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── user.go
│   │   └── product.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── rate_limit.go
│   └── routes/
│       └── routes.go
├── internal/
│   ├── service/
│   │   ├── auth.go
│   │   ├── user.go
│   │   └── product.go
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── user.go
│   │   │   └── product.go
│   │   └── redis/
│   │       └── cache.go
│   └── model/
│       ├── user.go
│       └── product.go
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   ├── database/
│   │   ├── postgres.go
│   │   └── redis.go
│   └── auth/
│       └── jwt.go
├── config/
│   ├── config.go
│   └── config.yaml
├── migrations/
│   └── 000001_init.up.sql
├── docs/
│   └── swagger.json
├── tests/
│   ├── integration/
│   └── unit/
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── main.go
```

## Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Redis 6+
- Docker & Docker Compose
- Make (optional)

## Getting Started

1. Clone the repository
2. Copy `.env.example` to `.env` and configure
3. Start dependencies:
```bash
docker-compose up -d postgres redis
```

4. Run migrations:
```bash
make migrate-up
```

5. Start the server:
```bash
make run
```

## API Documentation

Access Swagger documentation at: `http://localhost:8080/swagger/index.html`

### Key Endpoints

- Authentication:
  - POST /api/v1/auth/register
  - POST /api/v1/auth/login
  - POST /api/v1/auth/refresh

- Users:
  - GET /api/v1/users
  - GET /api/v1/users/:id
  - PUT /api/v1/users/:id
  - DELETE /api/v1/users/:id

- Products:
  - GET /api/v1/products
  - POST /api/v1/products
  - GET /api/v1/products/:id
  - PUT /api/v1/products/:id
  - DELETE /api/v1/products/:id

## Development

### Running Tests
```bash
make test        # Run unit tests
make test-int    # Run integration tests
make test-cov    # Generate coverage report
```

### Database Migrations
```bash
make migrate-up      # Apply migrations
make migrate-down    # Rollback migrations
make migrate-create  # Create new migration
```

### Code Quality
```bash
make lint       # Run linters
make fmt        # Format code
make vet        # Run Go vet
```

## Deployment

### Docker
```bash
# Build image
docker build -t backend-master .

# Run container
docker run -p 8080:8080 backend-master
```

### Docker Compose
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app
```

## Architecture

The project follows Clean Architecture principles:

1. API Layer (Handlers)
   - HTTP request handling
   - Input validation
   - Response formatting

2. Service Layer
   - Business logic
   - Use case implementation
   - Transaction management

3. Repository Layer
   - Data access
   - Cache management
   - Database operations

4. Domain Layer
   - Business entities
   - Domain logic
   - Interfaces

## Best Practices

1. Error Handling
   - Custom error types
   - Error wrapping
   - Proper error responses

2. Logging
   - Structured logging
   - Log levels
   - Request ID tracking

3. Security
   - JWT authentication
   - Password hashing
   - Input sanitization
   - Rate limiting

4. Performance
   - Connection pooling
   - Caching
   - Query optimization

5. Testing
   - Unit tests
   - Integration tests
   - Mocking
   - Test coverage

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License - See LICENSE file for details
