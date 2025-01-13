package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgreSQL represents a PostgreSQL database connection pool
type PostgreSQL struct {
	Pool *pgxpool.Pool
}

// PostgresConfig holds the configuration for PostgreSQL connection
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	Timeout  time.Duration
}

// NewPostgreSQL creates a new PostgreSQL connection pool
func NewPostgreSQL(ctx context.Context, cfg PostgresConfig) (*PostgreSQL, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	// Set connection pool configuration
	poolConfig.MaxConns = int32(cfg.MaxConns)
	poolConfig.ConnConfig.ConnectTimeout = cfg.Timeout

	// Create the connection pool
	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &PostgreSQL{Pool: pool}, nil
}

// Close closes the database connection pool
func (p *PostgreSQL) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

// HealthCheck performs a health check on the database
func (p *PostgreSQL) HealthCheck(ctx context.Context) error {
	return p.Pool.Ping(ctx)
}

// Transaction executes a function within a transaction
func (p *PostgreSQL) Transaction(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("tx failed: %v, rb failed: %v", err, rbErr)
			}
			return
		}
		err = tx.Commit(ctx)
	}()

	return fn(tx)
}

// ExecuteInTransaction executes multiple queries in a single transaction
func (p *PostgreSQL) ExecuteInTransaction(ctx context.Context, queries []string) error {
	return p.Transaction(ctx, func(tx pgx.Tx) error {
		for _, query := range queries {
			if _, err := tx.Exec(ctx, query); err != nil {
				return fmt.Errorf("failed to execute query: %w", err)
			}
		}
		return nil
	})
}
