package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type CacheEntry struct {
	Value      interface{} `json:"value"`
	ExpiresAt  time.Time   `json:"expires_at"`
	ReplicaIDs []string    `json:"replica_ids,omitempty"`
}

type Message struct {
	Type    string      `json:"type"`
	Key     string      `json:"key,omitempty"`
	Value   interface{} `json:"value,omitempty"`
	TTL     int         `json:"ttl,omitempty"`
	NodeID  string      `json:"node_id,omitempty"`
	Success bool        `json:"success,omitempty"`
}

type CacheNode struct {
	ID            string
	cache         map[string]CacheEntry
	peers         map[string]net.Conn
	listener      net.Listener
	mutex         sync.RWMutex
	replicaCount  int
	cleanupTicker *time.Ticker
}

func NewCacheNode(id string, port string, replicaCount int) (*CacheNode, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}

	node := &CacheNode{
		ID:           id,
		cache:        make(map[string]CacheEntry),
		peers:        make(map[string]net.Conn),
		listener:     listener,
		replicaCount: replicaCount,
	}

	// Start cleanup routine
	node.cleanupTicker = time.NewTicker(1 * time.Minute)
	go node.cleanup()

	return node, nil
}

func (n *CacheNode) Set(key string, value interface{}, ttl int) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	expiresAt := time.Now().Add(time.Duration(ttl) * time.Second)
	entry := CacheEntry{
		Value:     value,
		ExpiresAt: expiresAt,
	}

	n.cache[key] = entry

	// Replicate to peers
	msg := Message{
		Type:   "set",
		Key:    key,
		Value:  value,
		TTL:    ttl,
		NodeID: n.ID,
	}

	return n.replicateToPeers(msg)
}

func (n *CacheNode) Get(key string) (interface{}, bool) {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	entry, exists := n.cache[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		delete(n.cache, key)
		return nil, false
	}

	return entry.Value, true
}

func (n *CacheNode) Delete(key string) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	delete(n.cache, key)

	msg := Message{
		Type:   "delete",
		Key:    key,
		NodeID: n.ID,
	}

	return n.replicateToPeers(msg)
}

func (n *CacheNode) cleanup() {
	for range n.cleanupTicker.C {
		n.mutex.Lock()
		now := time.Now()
		for key, entry := range n.cache {
			if now.After(entry.ExpiresAt) {
				delete(n.cache, key)
			}
		}
		n.mutex.Unlock()
	}
}

func (n *CacheNode) handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			log.Printf("Error decoding message: %v", err)
			return
		}

		switch msg.Type {
		case "set":
			n.mutex.Lock()
			n.cache[msg.Key] = CacheEntry{
				Value:     msg.Value,
				ExpiresAt: time.Now().Add(time.Duration(msg.TTL) * time.Second),
			}
			n.mutex.Unlock()
			encoder.Encode(Message{Type: "ack", Success: true})

		case "delete":
			n.mutex.Lock()
			delete(n.cache, msg.Key)
			n.mutex.Unlock()
			encoder.Encode(Message{Type: "ack", Success: true})

		case "get":
			value, exists := n.Get(msg.Key)
			encoder.Encode(Message{
				Type:    "response",
				Value:   value,
				Success: exists,
			})

		case "join":
			n.mutex.Lock()
			n.peers[msg.NodeID] = conn
			n.mutex.Unlock()
			encoder.Encode(Message{Type: "ack", Success: true})

		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}
}

func (n *CacheNode) replicateToPeers(msg Message) error {
	for _, conn := range n.peers {
		encoder := json.NewEncoder(conn)
		if err := encoder.Encode(msg); err != nil {
			return err
		}

		// Wait for acknowledgment
		decoder := json.NewDecoder(conn)
		var response Message
		if err := decoder.Decode(&response); err != nil {
			return err
		}

		if !response.Success {
			return fmt.Errorf("replication failed")
		}
	}
	return nil
}

func (n *CacheNode) Start() error {
	fmt.Printf("Cache node %s listening on %s\n", n.ID, n.listener.Addr())

	for {
		conn, err := n.listener.Accept()
		if err != nil {
			return err
		}

		go n.handleConnection(conn)
	}
}

func (n *CacheNode) ConnectToPeer(peerID, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	// Send join message
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(Message{
		Type:   "join",
		NodeID: n.ID,
	}); err != nil {
		return err
	}

	// Wait for acknowledgment
	decoder := json.NewDecoder(conn)
	var response Message
	if err := decoder.Decode(&response); err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("failed to join peer")
	}

	n.mutex.Lock()
	n.peers[peerID] = conn
	n.mutex.Unlock()

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: cache-system <node-id> <port> [peer-addresses...]")
		os.Exit(1)
	}

	nodeID := os.Args[1]
	port := os.Args[2]
	peerAddresses := os.Args[3:]

	node, err := NewCacheNode(nodeID, port, 2)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to peers
	for i, addr := range peerAddresses {
		peerID := fmt.Sprintf("node-%d", i+1)
		if err := node.ConnectToPeer(peerID, addr); err != nil {
			log.Printf("Failed to connect to peer %s: %v", addr, err)
		}
	}

	if err := node.Start(); err != nil {
		log.Fatal(err)
	}
}
