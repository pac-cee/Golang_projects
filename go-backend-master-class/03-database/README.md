# Go Database Integration

This module demonstrates how to integrate a PostgreSQL database with Go using the standard `database/sql` package and the `lib/pq` driver.

## Code Explanation

The `main.go` file implements a complete database layer with the following components:

### 1. Data Model
```go
type User struct {
    ID        int
    Username  string
    Email     string
    CreatedAt time.Time
}
```
- Defines the User model structure
- Maps to the database table columns
- Uses appropriate Go types for each field

### 2. Database Configuration
```go
const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "postgres"
    dbname   = "gobackend"
)
```
- Centralizes database connection parameters
- Makes configuration easily modifiable
- Follows standard PostgreSQL connection format

### 3. Repository Pattern
```go
type UserRepository struct {
    db *sql.DB
}
```
- Implements the repository pattern for data access
- Encapsulates all database operations
- Provides clean separation of concerns

### 4. Database Initialization
```go
func initDB() (*sql.DB, error)
```
- Establishes database connection
- Uses connection pooling by default
- Verifies connection with `db.Ping()`
- Returns configured database instance

### 5. Schema Management
```go
func (r *UserRepository) createTables() error
```
- Creates necessary database tables
- Uses `IF NOT EXISTS` for idempotency
- Defines proper constraints (UNIQUE, NOT NULL)
- Sets up automatic timestamps

### 6. CRUD Operations
The repository implements full CRUD (Create, Read, Update, Delete) operations:

#### Create
```go
func (r *UserRepository) CreateUser(username, email string) (*User, error)
```
- Uses parameterized queries for safety
- Returns complete user object with ID
- Handles database errors properly

#### Read
```go
func (r *UserRepository) GetUser(id int) (*User, error)
```
- Retrieves user by ID
- Uses proper SQL scanning
- Returns structured data

#### Update
```go
func (r *UserRepository) UpdateUser(id int, email string) error
```
- Updates user information
- Validates existence before update
- Returns appropriate errors

#### Delete
```go
func (r *UserRepository) DeleteUser(id int) error
```
- Removes user records
- Handles non-existent records
- Returns operation status

## Best Practices Demonstrated
1. SQL injection prevention with parameterized queries
2. Proper error handling and propagation
3. Connection pooling and management
4. Repository pattern implementation
5. Clean code organization
6. Type safety with structs

## Database Setup
1. Install PostgreSQL
2. Create database:
```sql
CREATE DATABASE gobackend;
```
3. Configure connection parameters in the code

## Running the Code
```bash
# Install PostgreSQL driver
go get github.com/lib/pq

# Run the application
go run main.go
```

## Testing Database Operations
The main function includes examples of:
- Creating new users
- Retrieving user information
- Updating user data
- Deleting users
- Error handling for each operation
