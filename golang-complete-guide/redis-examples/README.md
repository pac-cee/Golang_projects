# Redis with Go Examples

This directory contains examples of using Redis with Go, demonstrating various Redis data structures and patterns.

## Prerequisites

1. Install Redis on Windows:
   - Download Redis for Windows from [https://github.com/microsoftarchive/redis/releases](https://github.com/microsoftarchive/redis/releases)
   - Run the installer
   - Redis server will be installed as a Windows service

2. Install Go Redis client:
```bash
go mod init redis-examples
go get github.com/redis/go-redis/v9
```

## Examples

1. **Basic Operations** (`basic/main.go`):
   - Connection setup
   - Simple GET/SET operations
   - Key expiration

2. **Data Structures** (`data_structures/main.go`):
   - Lists (for queues/stacks)
   - Sets (for unique collections)
   - Hashes (for structured data)
   - Sorted Sets (for rankings)

3. **Patterns** (`patterns/main.go`):
   - Caching
   - Rate limiting
   - Session management
   - Pub/Sub messaging

4. **Advanced** (`advanced/main.go`):
   - Pipelining
   - Transactions
   - Lua scripting
   - Error handling

## Running the Examples

Each directory contains a standalone example that can be run with:
```bash
cd [example-directory]
go run main.go
```

## Common Redis Use Cases

1. **Caching**:
   - Temporary data storage
   - Frequently accessed data
   - Session storage

2. **Rate Limiting**:
   - API request limiting
   - User action throttling

3. **Real-time Features**:
   - Live feeds
   - Real-time analytics
   - Message queues

4. **Session Management**:
   - User sessions
   - Authentication tokens
   - Temporary state

5. **Leaderboards**:
   - Game scores
   - User rankings
   - Real-time analytics
