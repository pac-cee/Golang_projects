package main

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/blog_test?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create migration manager and run migrations
	migrationManager := NewMigrationManager(db)
	if err := migrationManager.Migrate(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db, func() {
		// Clean up database
		db.Exec("DROP TABLE comments")
		db.Exec("DROP TABLE posts")
		db.Exec("DROP TABLE schema_migrations")
		db.Close()
	}
}

func TestBlogStore(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	store := NewBlogStore(db)
	ctx := context.Background()

	t.Run("create and get post", func(t *testing.T) {
		post := &Post{
			Title:   "Test Post",
			Content: "Test Content",
		}

		// Create post
		err := store.CreatePost(ctx, post)
		if err != nil {
			t.Fatalf("Failed to create post: %v", err)
		}
		if post.ID == 0 {
			t.Error("Post ID should not be 0")
		}

		// Get post
		retrieved, err := store.GetPost(ctx, post.ID)
		if err != nil {
			t.Fatalf("Failed to get post: %v", err)
		}
		if retrieved.Title != post.Title {
			t.Errorf("Expected title %q, got %q", post.Title, retrieved.Title)
		}
	})

	t.Run("update post", func(t *testing.T) {
		post := &Post{
			Title:   "Original Title",
			Content: "Original Content",
		}

		// Create post
		err := store.CreatePost(ctx, post)
		if err != nil {
			t.Fatalf("Failed to create post: %v", err)
		}

		// Update post
		post.Title = "Updated Title"
		err = store.UpdatePost(ctx, post)
		if err != nil {
			t.Fatalf("Failed to update post: %v", err)
		}

		// Verify update
		retrieved, err := store.GetPost(ctx, post.ID)
		if err != nil {
			t.Fatalf("Failed to get post: %v", err)
		}
		if retrieved.Title != "Updated Title" {
			t.Errorf("Expected title %q, got %q", "Updated Title", retrieved.Title)
		}
	})

	t.Run("delete post", func(t *testing.T) {
		post := &Post{
			Title:   "To Be Deleted",
			Content: "This post will be deleted",
		}

		// Create post
		err := store.CreatePost(ctx, post)
		if err != nil {
			t.Fatalf("Failed to create post: %v", err)
		}

		// Delete post
		err = store.DeletePost(ctx, post.ID)
		if err != nil {
			t.Fatalf("Failed to delete post: %v", err)
		}

		// Verify deletion
		_, err = store.GetPost(ctx, post.ID)
		if err == nil {
			t.Error("Expected error when getting deleted post")
		}
	})

	t.Run("add and get comments", func(t *testing.T) {
		post := &Post{
			Title:   "Post with Comments",
			Content: "This post will have comments",
		}

		// Create post
		err := store.CreatePost(ctx, post)
		if err != nil {
			t.Fatalf("Failed to create post: %v", err)
		}

		// Add comments
		comment1 := &Comment{
			PostID:  post.ID,
			Content: "First comment",
		}
		err = store.AddComment(ctx, comment1)
		if err != nil {
			t.Fatalf("Failed to add comment: %v", err)
		}

		comment2 := &Comment{
			PostID:  post.ID,
			Content: "Second comment",
		}
		err = store.AddComment(ctx, comment2)
		if err != nil {
			t.Fatalf("Failed to add comment: %v", err)
		}

		// Get comments
		comments, err := store.GetPostComments(ctx, post.ID)
		if err != nil {
			t.Fatalf("Failed to get comments: %v", err)
		}
		if len(comments) != 2 {
			t.Errorf("Expected 2 comments, got %d", len(comments))
		}
	})
}

func TestDBPool(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/blog_test?sslmode=disable"
	}

	t.Run("create pool", func(t *testing.T) {
		pool, err := NewDBPool(dsn, 5)
		if err != nil {
			t.Fatalf("Failed to create pool: %v", err)
		}
		defer pool.db.Close()

		if pool.maxConns != 5 {
			t.Errorf("Expected max connections 5, got %d", pool.maxConns)
		}
	})

	t.Run("get and release connection", func(t *testing.T) {
		pool, err := NewDBPool(dsn, 5)
		if err != nil {
			t.Fatalf("Failed to create pool: %v", err)
		}
		defer pool.db.Close()

		ctx := context.Background()
		conn, err := pool.GetConnection(ctx)
		if err != nil {
			t.Fatalf("Failed to get connection: %v", err)
		}

		err = pool.ReleaseConnection(conn)
		if err != nil {
			t.Fatalf("Failed to release connection: %v", err)
		}
	})

	t.Run("pool metrics", func(t *testing.T) {
		pool, err := NewDBPool(dsn, 5)
		if err != nil {
			t.Fatalf("Failed to create pool: %v", err)
		}
		defer pool.db.Close()

		metrics := pool.GetMetrics()
		if metrics.MaxConns != 5 {
			t.Errorf("Expected max connections 5, got %d", metrics.MaxConns)
		}
	})
}

func TestMigrationManager(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	manager := NewMigrationManager(db)

	t.Run("get current version", func(t *testing.T) {
		version, err := manager.GetCurrentVersion()
		if err != nil {
			t.Fatalf("Failed to get current version: %v", err)
		}
		if version != 2 {
			t.Errorf("Expected version 2, got %d", version)
		}
	})

	t.Run("rollback migration", func(t *testing.T) {
		err := manager.Rollback()
		if err != nil {
			t.Fatalf("Failed to rollback migration: %v", err)
		}

		version, err := manager.GetCurrentVersion()
		if err != nil {
			t.Fatalf("Failed to get current version: %v", err)
		}
		if version != 1 {
			t.Errorf("Expected version 1 after rollback, got %d", version)
		}
	})
}

func BenchmarkBlogStore(b *testing.B) {
	db, cleanup := setupTestDB(b)
	defer cleanup()

	store := NewBlogStore(db)
	ctx := context.Background()

	b.Run("create post", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			post := &Post{
				Title:   "Benchmark Post",
				Content: "This is a benchmark post",
			}
			err := store.CreatePost(ctx, post)
			if err != nil {
				b.Fatalf("Failed to create post: %v", err)
			}
		}
	})

	b.Run("get post", func(b *testing.B) {
		post := &Post{
			Title:   "Benchmark Post",
			Content: "This is a benchmark post",
		}
		err := store.CreatePost(ctx, post)
		if err != nil {
			b.Fatalf("Failed to create post: %v", err)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := store.GetPost(ctx, post.ID)
			if err != nil {
				b.Fatalf("Failed to get post: %v", err)
			}
		}
	})
}

func BenchmarkDBPool(b *testing.B) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/blog_test?sslmode=disable"
	}

	pool, err := NewDBPool(dsn, 10)
	if err != nil {
		b.Fatalf("Failed to create pool: %v", err)
	}
	defer pool.db.Close()

	ctx := context.Background()

	b.Run("get and release connection", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			conn, err := pool.GetConnection(ctx)
			if err != nil {
				b.Fatalf("Failed to get connection: %v", err)
			}
			err = pool.ReleaseConnection(conn)
			if err != nil {
				b.Fatalf("Failed to release connection: %v", err)
			}
		}
	})
}
