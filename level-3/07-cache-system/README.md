# Distributed Cache System

A distributed in-memory cache system built in Go, featuring peer-to-peer replication and automatic data expiration.

## Features

- In-memory key-value storage
- Distributed architecture with peer-to-peer replication
- TTL (Time-To-Live) support for cache entries
- Automatic cleanup of expired entries
- TCP-based communication between nodes
- Thread-safe operations
- JSON-based protocol for inter-node communication

## Architecture

The system implements a distributed cache with the following components:

1. **Cache Node**
   - Stores key-value pairs in memory
   - Handles peer connections
   - Manages data replication
   - Performs automatic cleanup

2. **Cache Entry**
   - Stores value with expiration time
   - Tracks replica locations

3. **Message Protocol**
   - Set operations
   - Get operations
   - Delete operations
   - Node join/leave operations

## Prerequisites

- Go 1.19 or later

## How to Run

1. Start the first node:
   ```bash
   go run main.go node1 8081
   ```

2. Start additional nodes:
   ```bash
   go run main.go node2 8082 localhost:8081
   go run main.go node3 8083 localhost:8081 localhost:8082
   ```

## Protocol Specification

### Message Types

1. **Set**
   ```json
   {
     "type": "set",
     "key": "example-key",
     "value": "example-value",
     "ttl": 3600,
     "node_id": "node1"
   }
   ```

2. **Get**
   ```json
   {
     "type": "get",
     "key": "example-key"
   }
   ```

3. **Delete**
   ```json
   {
     "type": "delete",
     "key": "example-key",
     "node_id": "node1"
   }
   ```

4. **Join**
   ```json
   {
     "type": "join",
     "node_id": "node2"
   }
   ```

### Response Format
```json
{
  "type": "response",
  "value": "example-value",
  "success": true
}
```

## Implementation Details

### Data Structures

1. **CacheEntry**
   - Value: The stored data
   - ExpiresAt: Expiration timestamp
   - ReplicaIDs: List of nodes containing replicas

2. **CacheNode**
   - ID: Unique node identifier
   - Cache: Thread-safe map of key-value pairs
   - Peers: Connected peer nodes
   - ReplicaCount: Number of replicas to maintain

### Key Components

1. **Concurrency Control**
   - Uses sync.RWMutex for thread-safe operations
   - Read/Write locking for cache access
   - Connection handling in separate goroutines

2. **Data Replication**
   - Automatic replication to peer nodes
   - Configurable replica count
   - Acknowledgment-based consistency

3. **Cleanup**
   - Periodic cleanup of expired entries
   - Configurable cleanup interval
   - Thread-safe cleanup operation

## Testing

To test the system:

1. Start multiple nodes as described above

2. Use netcat or telnet to send commands:
   ```bash
   echo '{"type":"set","key":"test","value":"hello","ttl":300}' | nc localhost 8081
   echo '{"type":"get","key":"test"}' | nc localhost 8082
   ```

## Performance Considerations

1. **Memory Usage**
   - In-memory storage
   - Automatic cleanup of expired entries
   - No persistence to disk

2. **Network**
   - TCP-based communication
   - JSON serialization overhead
   - Connection pooling for peers

3. **Concurrency**
   - Read/Write mutex for thread safety
   - Separate goroutines for connections
   - Lock contention under high load

## Next Steps

1. Add persistence layer
2. Implement consistent hashing
3. Add compression for network traffic
4. Implement failure detection
5. Add monitoring and metrics
6. Implement cache eviction policies
7. Add authentication and encryption
8. Implement request batching
9. Add support for complex data types
10. Implement cache warming
