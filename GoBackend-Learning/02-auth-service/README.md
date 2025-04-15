# Project 2: Auth Service (JWT-based)

## Overview
This project demonstrates how to build a simple authentication service in Go using JWT (JSON Web Tokens) for stateless authentication and bcrypt for secure password storage.

## Learning Goals
- Build and use HTTP middleware in Go
- Securely hash and check passwords with bcrypt
- Generate and validate JWT tokens for authentication
- Understand stateless session management

## Endpoints
- `POST /signup` — Register a new user. Body: `{ "username": "user", "password": "pass" }`
- `POST /login` — Authenticate and receive a JWT. Body: `{ "username": "user", "password": "pass" }`
- `GET /protected` — Example protected route. Requires `Authorization` header with JWT token.

## How to Run
```sh
# In the project directory
 go run main.go
```

## Why Go?
- Go's standard library offers strong support for cryptography and HTTP
- Middleware patterns are simple and powerful
- Fast, efficient, and secure backend development

## Example Usage
1. Signup:
   ```sh
   curl -X POST -d '{"username":"alice","password":"mypw"}' \
     -H "Content-Type: application/json" http://localhost:8081/signup
   ```
2. Login:
   ```sh
   curl -X POST -d '{"username":"alice","password":"mypw"}' \
     -H "Content-Type: application/json" http://localhost:8081/login
   # Response: { "token": "...jwt..." }
   ```
3. Access protected:
   ```sh
   curl -H "Authorization: <token>" http://localhost:8081/protected
   ```

## Next Steps
- Add persistent storage (database)
- Implement token refresh/expiry
- Add user roles/permissions

---

This project shows how Go makes secure authentication easy and robust!
