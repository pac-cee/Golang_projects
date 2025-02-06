package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Message represents a message in our system
type Message struct {
	ID        string      `json:"id"`
	Topic     string      `json:"topic"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// Subscription represents a subscription to a topic
type Subscription struct {
	Topic    string
	Channel  chan Message
	ClientID string
}

// MessageBroker handles pub/sub messaging
type MessageBroker struct {
	sync.RWMutex
	subscriptions map[string][]Subscription
	topics        map[string][]Message
}

// WebSocketClient represents a connected websocket client
type WebSocketClient struct {
	ID       string
	Conn     *websocket.Conn
	Broker   *MessageBroker
	Topics   []string
	SendChan chan Message
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for demo
	},
}

// NewMessageBroker creates a new message broker
func NewMessageBroker() *MessageBroker {
	return &MessageBroker{
		subscriptions: make(map[string][]Subscription),
		topics:       make(map[string][]Message),
	}
}

// Subscribe adds a new subscription for a topic
func (mb *MessageBroker) Subscribe(topic, clientID string) chan Message {
	mb.Lock()
	defer mb.Unlock()

	ch := make(chan Message, 100)
	sub := Subscription{
		Topic:    topic,
		Channel:  ch,
		ClientID: clientID,
	}

	mb.subscriptions[topic] = append(mb.subscriptions[topic], sub)
	return ch
}

// Unsubscribe removes a subscription
func (mb *MessageBroker) Unsubscribe(topic, clientID string) {
	mb.Lock()
	defer mb.Unlock()

	subs := mb.subscriptions[topic]
	for i, sub := range subs {
		if sub.ClientID == clientID {
			close(sub.Channel)
			mb.subscriptions[topic] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}

// Publish sends a message to all subscribers of a topic
func (mb *MessageBroker) Publish(msg Message) {
	mb.Lock()
	defer mb.Unlock()

	// Store message in topic history
	mb.topics[msg.Topic] = append(mb.topics[msg.Topic], msg)

	// Send to all subscribers
	for _, sub := range mb.subscriptions[msg.Topic] {
		select {
		case sub.Channel <- msg:
		default:
			// Channel is full, skip this subscriber
			log.Printf("Warning: Subscriber %s's channel is full", sub.ClientID)
		}
	}
}

// GetTopicHistory returns the message history for a topic
func (mb *MessageBroker) GetTopicHistory(topic string) []Message {
	mb.RLock()
	defer mb.RUnlock()
	return mb.topics[topic]
}

// NewWebSocketClient creates a new websocket client
func NewWebSocketClient(conn *websocket.Conn, broker *MessageBroker) *WebSocketClient {
	return &WebSocketClient{
		ID:       fmt.Sprintf("client-%d", time.Now().UnixNano()),
		Conn:     conn,
		Broker:   broker,
		Topics:   make([]string, 0),
		SendChan: make(chan Message, 100),
	}
}

// HandleWebSocket upgrades HTTP connection to WebSocket
func (mb *MessageBroker) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := NewWebSocketClient(conn, mb)
	go client.WritePump()
	go client.ReadPump()
}

// ReadPump handles incoming WebSocket messages
func (c *WebSocketClient) ReadPump() {
	defer func() {
		c.Conn.Close()
		// Cleanup subscriptions
		for _, topic := range c.Topics {
			c.Broker.Unsubscribe(topic, c.ID)
		}
	}()

	for {
		var msg struct {
			Action string `json:"action"`
			Topic  string `json:"topic"`
			Data   string `json:"data,omitempty"`
		}

		if err := c.Conn.ReadJSON(&msg); err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		switch msg.Action {
		case "subscribe":
			ch := c.Broker.Subscribe(msg.Topic, c.ID)
			c.Topics = append(c.Topics, msg.Topic)
			go func() {
				for message := range ch {
					c.SendChan <- message
				}
			}()

		case "publish":
			message := Message{
				ID:        fmt.Sprintf("msg-%d", time.Now().UnixNano()),
				Topic:     msg.Topic,
				Data:      msg.Data,
				Timestamp: time.Now(),
			}
			c.Broker.Publish(message)
		}
	}
}

// WritePump sends messages to the WebSocket client
func (c *WebSocketClient) WritePump() {
	ticker := time.NewTicker(time.Second * 30)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.SendChan:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("Error writing message: %v", err)
				return
			}

		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// HTTP handler for publishing messages
func (mb *MessageBroker) PublishHandler(w http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	msg.ID = fmt.Sprintf("msg-%d", time.Now().UnixNano())
	msg.Timestamp = time.Now()
	mb.Publish(msg)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Message published successfully",
	})
}

// HTTP handler for getting topic history
func (mb *MessageBroker) GetTopicHistoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topic := vars["topic"]

	messages := mb.GetTopicHistory(topic)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   messages,
	})
}

func main() {
	broker := NewMessageBroker()
	router := mux.NewRouter()

	// REST endpoints
	router.HandleFunc("/publish", broker.PublishHandler).Methods("POST")
	router.HandleFunc("/topics/{topic}/history", broker.GetTopicHistoryHandler).Methods("GET")

	// WebSocket endpoint
	router.HandleFunc("/ws", broker.HandleWebSocket)

	// Start server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Message broker starting on http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
