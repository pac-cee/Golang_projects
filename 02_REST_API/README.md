# Book Management API

A RESTful API built with Go for managing a personal book collection. This project demonstrates intermediate-level Go concepts and best practices for building web services.

## Features
- User authentication with JWT
- CRUD operations for books
- PostgreSQL database integration
- Middleware for authentication and logging
- Input validation
- Error handling
- Environment configuration

## Tech Stack
- Go
- PostgreSQL
- Gorilla Mux (routing)
- JWT for authentication
- Validator for input validation
- Bcrypt for password hashing

## Project Structure
```
02_REST_API/
├── config/         # Configuration settings
├── database/       # Database connection and migrations
├── handlers/       # Request handlers
├── middleware/     # Custom middleware
├── models/         # Data models
├── utils/          # Utility functions
├── .env.example    # Environment variables template
├── go.mod         # Go module file
├── main.go        # Application entry point
└── README.md      # Project documentation
```

## Prerequisites
- Go 1.21 or higher
- PostgreSQL
- Git

## Setup Instructions

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd 02_REST_API
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up the database:
   - Install PostgreSQL
   - Create a new database
   - Copy `.env.example` to `.env` and update the values

4. Run the application:
   ```bash
   go run main.go
   ```

## API Endpoints

### Authentication
- POST /api/register - Register a new user
- POST /api/login - Login and get JWT token

### Books (Protected Routes)
- GET /api/books - Get all books
- GET /api/books/{id} - Get a specific book
- POST /api/books - Create a new book
- PUT /api/books/{id} - Update a book
- DELETE /api/books/{id} - Delete a book

## Request Examples

### Register User
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "secure_password",
    "email": "john@example.com"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "secure_password"
  }'
```

### Create Book (with JWT)
```bash
curl -X POST http://localhost:8080/api/books \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Go Programming Language",
    "author": "Alan A. A. Donovan",
    "isbn": "9780134190440",
    "description": "A comprehensive guide to Go",
    "published_year": 2015
  }'
```

## Error Handling
The API uses standard HTTP status codes:
- 200: Success
- 201: Created
- 400: Bad Request
- 401: Unauthorized
- 404: Not Found
- 500: Internal Server Error

## Security Features
- Password hashing with bcrypt
- JWT-based authentication
- Input validation
- SQL injection prevention
- CORS support

## Development Practices
- Clean code architecture
- Dependency injection
- Middleware pattern
- Error handling best practices
- Input validation
- Logging

## Next Steps
1. Add unit tests
2. Implement rate limiting
3. Add API documentation with Swagger
4. Add caching layer
5. Implement refresh tokens
6. Add pagination for book listing
7. Implement search functionality

This project demonstrates intermediate Go concepts and serves as a foundation for building more complex web services.
