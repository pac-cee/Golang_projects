package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"chat-app/database"
	"chat-app/models"
	"chat-app/websocket"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

type ChatHandler struct {
	hub *websocket.Hub
}

func NewChatHandler(hub *websocket.Hub) *ChatHandler {
	return &ChatHandler{hub: hub}
}

// ServeWs handles websocket requests from clients
func (h *ChatHandler) ServeWs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomID"]

	// Get user information from query parameters (in a real app, this would come from authentication)
	userID := r.URL.Query().Get("user_id")
	username := r.URL.Query().Get("username")

	if userID == "" {
		userID = uuid.New().String()
	}
	if username == "" {
		username = "Anonymous-" + userID[:6]
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &websocket.Client{
		Hub:      h.hub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
		RoomID:   roomID,
		UserID:   userID,
		Username: username,
	}

	h.hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}

// GetRooms returns a list of active chat rooms
func (h *ChatHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms := h.hub.GetRooms()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"rooms": rooms,
	})
}

// CreateRoom creates a new chat room
func (h *ChatHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	roomID := uuid.New().String()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"room_id": roomID,
		"name":    input.Name,
	})
}

// GetRoom returns information about a specific room
func (h *ChatHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomID"]

	clients := h.hub.GetRoomClients(roomID)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"room_id":  roomID,
		"clients":  clients,
		"messages": []interface{}{}, // In a real app, you'd fetch messages from a database
	})
}

// GetMessages returns the message history for a room
func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["roomID"]

	// Get messages from database (limit to last 50 messages)
	messages, err := database.GetMessages(roomID, 50)
	if err != nil {
		http.Error(w, "Error retrieving messages", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": messages,
	})
}

// CorsMiddleware handles CORS
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
