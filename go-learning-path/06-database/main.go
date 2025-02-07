// Package main demonstrates database integration concepts in Go
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// Post represents a blog post
type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// BlogStore handles database operations for the blog
type BlogStore struct {
	db *sql.DB
}

// NewBlogStore creates a new BlogStore
func NewBlogStore(db *sql.DB) *BlogStore {
	return &BlogStore{db: db}
}

// CreatePost creates a new blog post
func (s *BlogStore) CreatePost(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (title, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now

	return s.db.QueryRowContext(
		ctx, query,
		post.Title,
		post.Content,
		post.CreatedAt,
		post.UpdatedAt,
	).Scan(&post.ID)
}

// GetPost retrieves a post by ID
func (s *BlogStore) GetPost(ctx context.Context, id int) (*Post, error) {
	post := &Post{}
	query := `
		SELECT id, title, content, created_at, updated_at
		FROM posts
		WHERE id = $1`

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("post not found: %d", id)
	}
	return post, err
}

// UpdatePost updates an existing post
func (s *BlogStore) UpdatePost(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, updated_at = $3
		WHERE id = $4`

	post.UpdatedAt = time.Now()
	result, err := s.db.ExecContext(
		ctx, query,
		post.Title,
		post.Content,
		post.UpdatedAt,
		post.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("post not found: %d", post.ID)
	}
	return nil
}

// DeletePost deletes a post and its comments
func (s *BlogStore) DeletePost(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete comments first
	_, err = tx.ExecContext(ctx, "DELETE FROM comments WHERE post_id = $1", id)
	if err != nil {
		return err
	}

	// Delete the post
	result, err := tx.ExecContext(ctx, "DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("post not found: %d", id)
	}

	return tx.Commit()
}

// AddComment adds a comment to a post
func (s *BlogStore) AddComment(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (post_id, content, created_at)
		VALUES ($1, $2, $3)
		RETURNING id`

	comment.CreatedAt = time.Now()

	return s.db.QueryRowContext(
		ctx, query,
		comment.PostID,
		comment.Content,
		comment.CreatedAt,
	).Scan(&comment.ID)
}

// GetPostComments retrieves all comments for a post
func (s *BlogStore) GetPostComments(ctx context.Context, postID int) ([]Comment, error) {
	query := `
		SELECT id, post_id, content, created_at
		FROM comments
		WHERE post_id = $1
		ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, rows.Err()
}

// DBPool represents a database connection pool with monitoring
type DBPool struct {
	db          *sql.DB
	maxConns    int
	activeConns int
	mu          sync.RWMutex
	metrics     *Metrics
}

// Metrics stores database pool metrics
type Metrics struct {
	MaxConns     int
	ActiveConns  int
	IdleConns    int
	WaitCount    int64
	WaitDuration time.Duration
	MaxIdleTime  time.Duration
}

// NewDBPool creates a new database connection pool
func NewDBPool(dsn string, maxConns int) (*DBPool, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(maxConns)
	db.SetConnMaxLifetime(5 * time.Minute)

	pool := &DBPool{
		db:       db,
		maxConns: maxConns,
		metrics:  &Metrics{MaxConns: maxConns},
	}

	// Start metrics collection
	go pool.collectMetrics()

	return pool, nil
}

// GetConnection gets a connection from the pool
func (p *DBPool) GetConnection(ctx context.Context) (*sql.Conn, error) {
	p.mu.Lock()
	if p.activeConns >= p.maxConns {
		p.mu.Unlock()
		return nil, errors.New("connection pool exhausted")
	}
	p.activeConns++
	p.mu.Unlock()

	conn, err := p.db.Conn(ctx)
	if err != nil {
		p.mu.Lock()
		p.activeConns--
		p.mu.Unlock()
		return nil, err
	}

	return conn, nil
}

// ReleaseConnection releases a connection back to the pool
func (p *DBPool) ReleaseConnection(conn *sql.Conn) error {
	p.mu.Lock()
	p.activeConns--
	p.mu.Unlock()

	return conn.Close()
}

// GetMetrics returns current pool metrics
func (p *DBPool) GetMetrics() *Metrics {
	p.mu.RLock()
	defer p.mu.RUnlock()

	stats := p.db.Stats()
	return &Metrics{
		MaxConns:     p.maxConns,
		ActiveConns:  p.activeConns,
		IdleConns:    stats.Idle,
		WaitCount:    stats.WaitCount,
		WaitDuration: stats.WaitDuration,
		MaxIdleTime:  stats.MaxIdleClosed,
	}
}

// collectMetrics periodically collects pool metrics
func (p *DBPool) collectMetrics() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		metrics := p.GetMetrics()
		log.Printf("DB Pool Metrics: Active=%d, Idle=%d, Wait=%d",
			metrics.ActiveConns,
			metrics.IdleConns,
			metrics.WaitCount,
		)
	}
}

// Migration represents a database migration
type Migration struct {
	Version     int
	Description string
	UpSQL       string
	DownSQL     string
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db         *sql.DB
	migrations []Migration
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sql.DB) *MigrationManager {
	return &MigrationManager{
		db: db,
		migrations: []Migration{
			{
				Version:     1,
				Description: "Create posts table",
				UpSQL: `
					CREATE TABLE posts (
						id SERIAL PRIMARY KEY,
						title VARCHAR(255) NOT NULL,
						content TEXT NOT NULL,
						created_at TIMESTAMP NOT NULL,
						updated_at TIMESTAMP NOT NULL
					)`,
				DownSQL: "DROP TABLE posts",
			},
			{
				Version:     2,
				Description: "Create comments table",
				UpSQL: `
					CREATE TABLE comments (
						id SERIAL PRIMARY KEY,
						post_id INTEGER REFERENCES posts(id),
						content TEXT NOT NULL,
						created_at TIMESTAMP NOT NULL
					)`,
				DownSQL: "DROP TABLE comments",
			},
		},
	}
}

// GetCurrentVersion gets the current migration version
func (m *MigrationManager) GetCurrentVersion() (int, error) {
	var version int
	err := m.db.QueryRow("SELECT version FROM schema_migrations").Scan(&version)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return version, err
}

// Migrate runs all pending migrations
func (m *MigrationManager) Migrate() error {
	// Create migrations table if it doesn't exist
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL
		)`)
	if err != nil {
		return err
	}

	currentVersion, err := m.GetCurrentVersion()
	if err != nil {
		return err
	}

	for _, migration := range m.migrations {
		if migration.Version > currentVersion {
			tx, err := m.db.Begin()
			if err != nil {
				return err
			}

			log.Printf("Applying migration %d: %s", migration.Version, migration.Description)

			_, err = tx.Exec(migration.UpSQL)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("migration %d failed: %v", migration.Version, err)
			}

			_, err = tx.Exec(
				"INSERT INTO schema_migrations (version, applied_at) VALUES ($1, $2)",
				migration.Version,
				time.Now(),
			)
			if err != nil {
				tx.Rollback()
				return err
			}

			err = tx.Commit()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Rollback rolls back the last migration
func (m *MigrationManager) Rollback() error {
	currentVersion, err := m.GetCurrentVersion()
	if err != nil {
		return err
	}

	if currentVersion == 0 {
		return errors.New("no migrations to roll back")
	}

	var migration Migration
	for i := len(m.migrations) - 1; i >= 0; i-- {
		if m.migrations[i].Version == currentVersion {
			migration = m.migrations[i]
			break
		}
	}

	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	log.Printf("Rolling back migration %d: %s", migration.Version, migration.Description)

	_, err = tx.Exec(migration.DownSQL)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("rollback of migration %d failed: %v", migration.Version, err)
	}

	_, err = tx.Exec("DELETE FROM schema_migrations WHERE version = $1", migration.Version)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func main() {
	// Get database connection string from environment
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/blog?sslmode=disable"
	}

	// Create connection pool
	pool, err := NewDBPool(dsn, 10)
	if err != nil {
		log.Fatal(err)
	}

	// Create migration manager and run migrations
	migrationManager := NewMigrationManager(pool.db)
	if err := migrationManager.Migrate(); err != nil {
		log.Fatal(err)
	}

	// Create blog store
	store := NewBlogStore(pool.db)

	// Example usage
	ctx := context.Background()

	// Create a post
	post := &Post{
		Title:   "Hello, Database!",
		Content: "This is a test post demonstrating database operations in Go.",
	}
	if err := store.CreatePost(ctx, post); err != nil {
		log.Fatal(err)
	}
	log.Printf("Created post: %d", post.ID)

	// Add a comment
	comment := &Comment{
		PostID:  post.ID,
		Content: "Great post about database operations!",
	}
	if err := store.AddComment(ctx, comment); err != nil {
		log.Fatal(err)
	}
	log.Printf("Added comment: %d", comment.ID)

	// Get post with comments
	retrievedPost, err := store.GetPost(ctx, post.ID)
	if err != nil {
		log.Fatal(err)
	}

	comments, err := store.GetPostComments(ctx, post.ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Retrieved post: %s", retrievedPost.Title)
	log.Printf("Number of comments: %d", len(comments))
}
