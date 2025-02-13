package database

import (
    "context"
    "fmt"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type PostgresConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg PostgresConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
    )

    // Configure GORM logger
    gormConfig := &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    }

    // Connect to database with retry mechanism
    var db *gorm.DB
    var err error
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        db, err = gorm.Open(postgres.Open(dsn), gormConfig)
        if err == nil {
            break
        }
        time.Sleep(time.Second * time.Duration(i+1))
    }

    if err != nil {
        return nil, fmt.Errorf("failed to connect to database after %d retries: %v", maxRetries, err)
    }

    // Configure connection pool
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get database instance: %v", err)
    }

    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := sqlDB.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %v", err)
    }

    return db, nil
}

// CloseDB closes the database connection
func CloseDB(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return fmt.Errorf("failed to get database instance: %v", err)
    }
    return sqlDB.Close()
}

// Health checks the database connection
func Health(db *gorm.DB) error {
    sqlDB, err := db.DB()
    if err != nil {
        return fmt.Errorf("failed to get database instance: %v", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    return sqlDB.PingContext(ctx)
}