# Database Integration in Go 🗄️

## 📖 Table of Contents
1. [Database Basics](#database-basics)
2. [SQL Databases](#sql-databases)
3. [NoSQL Databases](#nosql-databases)
4. [ORM Usage](#orm-usage)
5. [Migration Management](#migration-management)
6. [Best Practices](#best-practices)
7. [Exercises](#exercises)

## Database Basics

### Connection Management
```go
import (
    "database/sql"
    _ "github.com/lib/pq"
)

db, err := sql.Open("postgres", "postgres://user:pass@localhost/dbname?sslmode=disable")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Set connection pool settings
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)
```

### Connection Pool Best Practices
1. Set appropriate pool size
2. Monitor connection usage
3. Use connection timeouts
4. Handle connection errors gracefully

## SQL Databases

### Basic CRUD Operations
```go
// Create
result, err := db.Exec(
    "INSERT INTO users (name, email) VALUES ($1, $2)",
    "John Doe",
    "john@example.com",
)

// Read
rows, err := db.Query("SELECT * FROM users WHERE active = $1", true)
for rows.Next() {
    var user User
    err := rows.Scan(&user.ID, &user.Name, &user.Email)
}

// Update
result, err := db.Exec(
    "UPDATE users SET name = $1 WHERE id = $2",
    "Jane Doe",
    1,
)

// Delete
result, err := db.Exec("DELETE FROM users WHERE id = $1", 1)
```

### Prepared Statements
```go
stmt, err := db.Prepare("INSERT INTO users(name, email) VALUES($1, $2)")
defer stmt.Close()

for _, user := range users {
    _, err := stmt.Exec(user.Name, user.Email)
}
```

### Transactions
```go
tx, err := db.Begin()
if err != nil {
    return err
}
defer tx.Rollback()

// Perform multiple operations
_, err = tx.Exec(query1, args1...)
_, err = tx.Exec(query2, args2...)

return tx.Commit()
```

## NoSQL Databases

### MongoDB Example
```go
import "go.mongodb.org/mongo-driver/mongo"

client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
defer client.Disconnect(ctx)

collection := client.Database("test").Collection("users")

// Insert
result, err := collection.InsertOne(ctx, bson.D{
    {"name", "John"},
    {"email", "john@example.com"},
})

// Find
var user User
err = collection.FindOne(ctx, bson.D{{"name", "John"}}).Decode(&user)
```

### Redis Example
```go
import "github.com/go-redis/redis/v8"

rdb := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// Set value
err := rdb.Set(ctx, "key", "value", 0).Err()

// Get value
val, err := rdb.Get(ctx, "key").Result()
```

## ORM Usage

### GORM Example
```go
import "gorm.io/gorm"

type User struct {
    gorm.Model
    Name  string
    Email string
}

db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// Create
user := User{Name: "John", Email: "john@example.com"}
db.Create(&user)

// Read
var user User
db.First(&user, 1)
db.Where("name = ?", "John").First(&user)

// Update
db.Model(&user).Update("Name", "Jane")

// Delete
db.Delete(&user)
```

## Migration Management

### Using golang-migrate
```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration
migrate create -ext sql -dir migrations -seq create_users_table

# Run migrations
migrate -path migrations -database "postgres://localhost:5432/dbname?sslmode=disable" up

# Rollback
migrate -path migrations -database "postgres://localhost:5432/dbname?sslmode=disable" down
```

### Migration Files
```sql
-- 000001_create_users_table.up.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- 000001_create_users_table.down.sql
DROP TABLE users;
```

## Best Practices

### 1. Connection Management
- Use connection pooling
- Set appropriate timeouts
- Handle connection errors gracefully
- Monitor connection pool metrics

### 2. Query Optimization
- Use prepared statements for repeated queries
- Index frequently queried columns
- Use EXPLAIN ANALYZE to understand query performance
- Avoid N+1 query problems

### 3. Transaction Management
- Use transactions for atomic operations
- Keep transactions short
- Handle deadlocks gracefully
- Use appropriate isolation levels

### 4. Error Handling
- Check for specific database errors
- Implement retry mechanisms
- Log database errors appropriately
- Use context for cancellation

### 5. Security
- Use parameterized queries to prevent SQL injection
- Encrypt sensitive data
- Use TLS for database connections
- Implement proper access control

## Exercises

### Exercise 1: Basic CRUD
Implement a basic CRUD application for a blog system with posts and comments.

```go
type Post struct {
    ID        int
    Title     string
    Content   string
    CreatedAt time.Time
}

type Comment struct {
    ID        int
    PostID    int
    Content   string
    CreatedAt time.Time
}

// Implement:
// - CreatePost(post *Post) error
// - GetPost(id int) (*Post, error)
// - UpdatePost(post *Post) error
// - DeletePost(id int) error
// - AddComment(comment *Comment) error
// - GetPostComments(postID int) ([]Comment, error)
```

### Exercise 2: Connection Pool
Create a connection pool manager with monitoring capabilities.

```go
type DBPool struct {
    db          *sql.DB
    maxConns    int
    activeConns int
    metrics     *Metrics
}

// Implement:
// - NewDBPool(dsn string, maxConns int) (*DBPool, error)
// - GetConnection() (*sql.Conn, error)
// - ReleaseConnection(*sql.Conn) error
// - GetMetrics() *Metrics
```

### Exercise 3: Migration Tool
Create a simple database migration tool.

```go
type Migration struct {
    Version     int
    Description string
    UpSQL       string
    DownSQL     string
}

// Implement:
// - LoadMigrations(dir string) ([]Migration, error)
// - ApplyMigration(db *sql.DB, m Migration) error
// - RollbackMigration(db *sql.DB, m Migration) error
// - GetCurrentVersion(db *sql.DB) (int, error)
```

## Next Steps
- Study database performance optimization
- Learn about database replication and sharding
- Explore NoSQL database patterns
- Practice database design patterns
