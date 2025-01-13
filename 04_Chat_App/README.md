# Real-time Chat Application

A real-time chat application built with Go and WebSocket, demonstrating advanced web development concepts.

## Features

- Real-time messaging using WebSocket
- Multiple chat rooms
- User presence (join/leave notifications)
- Clean and modern UI with Tailwind CSS
- Responsive design
- Message history
- Room management

## Tech Stack

- **Backend**:
  - Go
  - Gorilla WebSocket
  - Gorilla Mux
  - UUID generation
  
- **Frontend**:
  - HTML5
  - CSS3 (Tailwind CSS)
  - JavaScript (Vanilla)
  - WebSocket API

## Project Structure

```
04_Chat_App/
├── handlers/
│   └── chat_handler.go    # HTTP and WebSocket handlers
├── websocket/
│   ├── hub.go            # WebSocket hub for managing connections
│   └── client.go         # WebSocket client implementation
├── static/
│   ├── index.html        # Main HTML file
│   ├── style.css         # Custom styles
│   └── app.js            # Frontend JavaScript
├── main.go               # Application entry point
├── go.mod               # Go module file
├── .env.example         # Environment variables template
└── README.md           # Project documentation
```

## Prerequisites

- Go 1.21 or higher
- Modern web browser with WebSocket support

## Setup Instructions

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd 04_Chat_App
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Copy environment file and update if needed:
   ```bash
   cp .env.example .env
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

5. Open your browser and navigate to:
   ```
   http://localhost:8080
   ```

## Usage

1. Enter your username and click "Join"
2. Create a new room or join an existing one
3. Start chatting!

## Features Demonstrated

1. **WebSocket Communication**
   - Real-time bidirectional communication
   - Connection management
   - Message broadcasting

2. **Concurrent Programming**
   - Goroutines for handling connections
   - Channel-based communication
   - Thread-safe operations

3. **Frontend Development**
   - Modern UI with Tailwind CSS
   - Real-time updates
   - Responsive design

4. **Clean Architecture**
   - Separation of concerns
   - Modular design
   - Error handling

## API Endpoints

- `GET /` - Serves the chat application
- `GET /ws/{roomID}` - WebSocket endpoint
- `GET /api/rooms` - Get list of active rooms
- `POST /api/rooms` - Create a new room
- `GET /api/rooms/{roomID}` - Get room details
- `GET /api/rooms/{roomID}/messages` - Get room messages

## WebSocket Events

- `message` - New chat message
- `join` - User joined the room
- `leave` - User left the room

## Security Features

- Input sanitization
- Origin checking for WebSocket connections
- Rate limiting (TODO)
- Message size limits

## Future Improvements

1. Add persistent storage for messages
2. Implement user authentication
3. Add private messaging
4. Add file sharing
5. Implement typing indicators
6. Add emoji support
7. Add message reactions
8. Implement room moderation
9. Add rate limiting
10. Add message search

This project demonstrates intermediate to advanced Go concepts and serves as a foundation for building more complex real-time applications.
