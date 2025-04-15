# Project 1: Simple REST API (To-Do List)

## Overview
This project demonstrates how to build a basic RESTful API in Go for managing a to-do list. You'll learn how to use Go's standard library to create HTTP servers, handle routes, encode/decode JSON, and manage in-memory data.

## Learning Goals
- Understand Go's `net/http` package for web servers
- Learn about HTTP handlers, routing, and request/response lifecycle
- Practice using structs, slices, and JSON encoding/decoding
- Implement basic CRUD (Create, Read, Update, Delete) operations

## Endpoints
- `GET /todos` — List all to-do items
- (You can extend this project with POST, PUT, DELETE for full CRUD)

## How to Run
```sh
# In the project directory
# Run the server
 go run main.go
```

Visit [http://localhost:8080/todos](http://localhost:8080/todos) in your browser or use curl/Postman.

## Why Go?
- Go's standard library makes it easy to build fast, lightweight APIs
- Static typing and simple syntax help prevent bugs
- Built-in concurrency (not used here, but great for scaling APIs)

## Sample Response
```json
[
  {"id": 1, "task": "Learn Go basics", "done": false},
  {"id": 2, "task": "Build a REST API", "done": false}
]
```

## Next Steps
- Add endpoints for creating, updating, and deleting tasks
- Persist data with a database (see later projects)
- Add authentication (see Project 2)

---

This project is your entry point to Go backend development. Experiment, extend, and have fun!
