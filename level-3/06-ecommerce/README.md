# Microservice-based E-commerce System

A scalable e-commerce system built with microservices architecture using Go, MongoDB, and Docker.

## Architecture Overview

The system consists of four main services:

1. **API Gateway (Port 8080)**
   - Routes requests to appropriate services
   - Handles authentication
   - Implements CORS and logging

2. **Product Service (Port 8081)**
   - Manages product catalog
   - Handles inventory
   - CRUD operations for products

3. **Order Service (Port 8082)**
   - Processes orders
   - Manages order status
   - Validates product availability

4. **User Service (Port 8083)**
   - User authentication
   - Profile management
   - JWT token generation

## Technologies Used

- Go (Golang) 1.19
- MongoDB
- Docker & Docker Compose
- JWT for authentication
- Gorilla Mux for routing
- WebSocket for real-time updates

## Prerequisites

1. Docker and Docker Compose installed
2. Go 1.19 or later (for development)
3. MongoDB (handled by Docker)

## Project Structure

```
06-ecommerce/
├── api-gateway/
│   ├── main.go
│   └── Dockerfile
├── services/
│   ├── product/
│   │   ├── main.go
│   │   ├── Dockerfile
│   │   └── go.mod
│   ├── order/
│   │   ├── main.go
│   │   ├── Dockerfile
│   │   └── go.mod
│   └── user/
│       ├── main.go
│       ├── Dockerfile
│       └── go.mod
├── docker-compose.yml
└── README.md
```

## How to Run

1. Clone the repository
2. Navigate to the project directory:
   ```bash
   cd level-3/06-ecommerce
   ```

3. Start the services:
   ```bash
   docker-compose up --build
   ```

4. Access the API at `http://localhost:8080`

## API Endpoints

### Product Service
- `POST /api/products` - Create product
- `GET /api/products` - List products
- `GET /api/products/{id}` - Get product
- `PUT /api/products/{id}` - Update product
- `DELETE /api/products/{id}` - Delete product

### Order Service
- `POST /api/orders` - Create order
- `GET /api/orders` - List orders
- `GET /api/orders/{id}` - Get order
- `PUT /api/orders/{id}/status` - Update order status

### User Service
- `POST /api/auth/register` - Register user
- `POST /api/auth/login` - Login user
- `GET /api/users/{id}` - Get user profile
- `PUT /api/users/{id}` - Update user profile

## Authentication

The system uses JWT tokens for authentication. To access protected endpoints:

1. Register a user or login to get a JWT token
2. Include the token in the Authorization header:
   ```
   Authorization: Bearer <your-token>
   ```

## Development

Each service can be developed and tested independently:

1. Navigate to a service directory:
   ```bash
   cd services/product
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the service:
   ```bash
   go run main.go
   ```

## Testing

To test the services:

1. Register a user:
   ```bash
   curl -X POST http://localhost:8080/api/auth/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password","name":"Test User"}'
   ```

2. Login to get a token:
   ```bash
   curl -X POST http://localhost:8080/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password"}'
   ```

3. Use the token to create a product:
   ```bash
   curl -X POST http://localhost:8080/api/products \
     -H "Authorization: Bearer <your-token>" \
     -H "Content-Type: application/json" \
     -d '{"name":"Test Product","price":29.99,"stock":100}'
   ```

## Next Steps

1. Add payment service integration
2. Implement caching with Redis
3. Add message queues for async operations
4. Implement service discovery
5. Add monitoring and logging
6. Implement CI/CD pipeline
7. Add unit and integration tests
8. Implement rate limiting
9. Add data backup and recovery
10. Implement search functionality
