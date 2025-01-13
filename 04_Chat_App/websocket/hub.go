package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"chat-app/database"
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients by room
	rooms map[string]map[*Client]bool

	// Messages to be broadcast to the clients in a room
	broadcast chan *Message

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for concurrent access to rooms
	mu sync.RWMutex
}

type Message struct {
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	Type      string    `json:"type"` // "message", "join", "leave"
	Timestamp time.Time `json:"timestamp"`
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.rooms[client.roomID]; !ok {
				h.rooms[client.roomID] = make(map[*Client]bool)
			}
			h.rooms[client.roomID][client] = true
			h.mu.Unlock()

			// Create join message
			joinMsg := &Message{
				RoomID:    client.roomID,
				UserID:    client.userID,
				Username:  client.username,
				Type:      "join",
				Timestamp: time.Now(),
				Content:   client.username + " joined the room",
			}

			// Save join message to database
			if err := database.SaveMessage((*database.Message)(joinMsg)); err != nil {
				log.Printf("Error saving join message: %v", err)
			}

			// Broadcast join message
			h.broadcast <- joinMsg

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.rooms[client.roomID]; ok {
				if _, ok := h.rooms[client.roomID][client]; ok {
					delete(h.rooms[client.roomID], client)
					close(client.send)
					// If room is empty, delete it
					if len(h.rooms[client.roomID]) == 0 {
						delete(h.rooms, client.roomID)
					}
				}
			}
			h.mu.Unlock()

			// Create leave message
			leaveMsg := &Message{
				RoomID:    client.roomID,
				UserID:    client.userID,
				Username:  client.username,
				Type:      "leave",
				Timestamp: time.Now(),
				Content:   client.username + " left the room",
			}

			// Save leave message to database
			if err := database.SaveMessage((*database.Message)(leaveMsg)); err != nil {
				log.Printf("Error saving leave message: %v", err)
			}

			// Broadcast leave message
			h.broadcast <- leaveMsg

		case message := <-h.broadcast:
			// Save message to database if it's a chat message
			if message.Type == "message" {
				if err := database.SaveMessage((*database.Message)(message)); err != nil {
					log.Printf("Error saving message: %v", err)
				}
			}

			h.mu.RLock()
			if clients, ok := h.rooms[message.RoomID]; ok {
				// Marshal the message
				data, err := json.Marshal(message)
				if err != nil {
					log.Printf("Error marshaling message: %v", err)
					h.mu.RUnlock()
					continue
				}

				// Send to all clients in the room
				for client := range clients {
					select {
					case client.send <- data:
					default:
						close(client.send)
						delete(clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// GetRoomClients returns a list of clients in a room
func (h *Hub) GetRoomClients(roomID string) []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var clients []string
	if room, ok := h.rooms[roomID]; ok {
		for client := range room {
			clients = append(clients, client.username)
		}
	}
	return clients
}

// GetRooms returns a list of active rooms
func (h *Hub) GetRooms() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	rooms := make([]string, 0, len(h.rooms))
	for room := range h.rooms {
		rooms = append(rooms, room)
	}
	return rooms
}
