# Project 3: Real-time Chat Server (WebSockets)

## Overview
This project demonstrates a real-time chat server using Go's goroutines, channels, and the WebSocket protocol. It allows multiple clients to connect and exchange messages instantly.

## Learning Goals
- Use goroutines and channels for concurrency
- Implement the WebSocket protocol for real-time communication
- Manage multiple client connections efficiently

## Endpoints
- `GET /ws` — WebSocket endpoint for chat clients

## How to Run
```sh
# In the project directory
 go run main.go
```

Connect with a WebSocket client (e.g., browser, Postman, or `websocat`):
- `ws://localhost:8082/ws`

## Why Go?
- Goroutines and channels make handling many clients easy and efficient
- The `gorilla/websocket` package simplifies WebSocket handling

## Example Usage
- Open two browser tabs or WebSocket clients and connect to `ws://localhost:8082/ws`
- Send a JSON message: `{ "username": "alice", "content": "Hello!" }`
- All connected clients receive the message in real time

## Next Steps
- Add user authentication
- Store chat history
- Support private messaging or chat rooms

---

This project shows Go's strength in real-time, concurrent backend systems!
