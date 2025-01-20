# Real-time Chat Application

A web-based chat application demonstrating real-time communication using WebSocket protocol in Go.

## Concepts Covered

- WebSocket protocol
- Goroutines and channels
- Concurrent programming
- JSON encoding/decoding
- Web server creation
- Client-server communication
- Event handling
- Frontend development

## Features

- Real-time messaging
- User join/leave notifications
- Message timestamps
- Clean and responsive UI
- Username support
- System messages
- Auto-scroll to latest messages

## Prerequisites

1. Go installed on your system
2. Install required Go package:
   ```bash
   go get github.com/gorilla/websocket
   ```

## How to Run

1. Navigate to the project directory:
   ```bash
   cd level-2/05-chat-app
   ```

2. Run the program:
   ```bash
   go run main.go
   ```

3. Open your browser and visit:
   ```
   http://localhost:8080
   ```

4. Open multiple browser windows to simulate different users.

## Project Structure

```
05-chat-app/
├── main.go          # Main server file with WebSocket handling
├── README.md        # Project documentation
└── static/
    ├── index.html   # Frontend HTML
    └── style.css    # CSS styles
```

## Components

### Backend
- WebSocket server implementation
- Client connection management
- Message broadcasting system
- Concurrent message handling

### Frontend
- Real-time WebSocket client
- Dynamic message display
- User interface components
- Message input handling

## Message Types

```go
type Message struct {
    Type     string `json:"type"`     // "message" or "system"
    Content  string `json:"content"`   // Message content
    Username string `json:"username"`  // Sender's username
    Time     string `json:"time"`      // Message timestamp
}
```

## Learning Objectives

- Understanding WebSocket protocol
- Managing concurrent connections
- Implementing pub/sub pattern
- Handling real-time events
- Frontend-backend integration
- Error handling in real-time systems
- UI/UX design principles

## Next Steps

To extend this project, you could:
1. Add private messaging
2. Implement chat rooms
3. Add user authentication
4. Store chat history
5. Add file sharing
6. Implement typing indicators
7. Add user profiles
8. Add emoji support
9. Implement message editing
10. Add unit tests
11. Deploy to a cloud platform
12. Add message encryption
13. Implement rate limiting
