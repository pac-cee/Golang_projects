package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

// DatabaseType represents the type of database being used
type DatabaseType string

const (
	MongoDB    DatabaseType = "mongodb"
	PostgreSQL DatabaseType = "postgresql"
	MySQL      DatabaseType = "mysql"
)

// Database interface defines common methods that all database implementations must support
type Database interface {
	Close() error
	Ping() error
	GetType() DatabaseType
}

// MongoDatabase wraps MongoDB connection
type MongoDatabase struct {
	Client *mongo.Client
	DB     *mongo.Database
}

// PostgresDB wraps PostgreSQL connection
type PostgresDB struct {
	DB *sql.DB
}

// MySQLDB wraps MySQL connection
type MySQLDB struct {
	DB *sql.DB
}

var (
	activeDB Database
	dbType   DatabaseType
)

// ConnectDB establishes a connection to the specified database type
func ConnectDB(dbType DatabaseType) (Database, error) {
	switch dbType {
	case MongoDB:
		return connectMongoDB()
	case PostgreSQL:
		return NewPostgresDB()
	case MySQL:
		return NewMySQLDB()
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}

// GetDB returns the active database connection
func GetDB() Database {
	return activeDB
}

// GetDBType returns the current database type
func GetDBType() DatabaseType {
	return dbType
}

// connectMongoDB establishes a connection to MongoDB
func connectMongoDB() (*MongoDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	db := client.Database(os.Getenv("DB_NAME"))
	log.Printf("Connected to MongoDB: %s\n", os.Getenv("DB_NAME"))

	// Create indexes
	if err := createMongoIndexes(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to create MongoDB indexes: %v", err)
	}

	return &MongoDatabase{
		Client: client,
		DB:     db,
	}, nil
}

// createMongoIndexes creates necessary indexes for MongoDB
func createMongoIndexes(ctx context.Context, db *mongo.Database) error {
	// User indexes
	userColl := db.Collection("users")
	_, err := userColl.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	// Expense indexes
	expenseColl := db.Collection("expenses")
	_, err = expenseColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "userId", Value: 1},
				{Key: "date", Value: -1},
			},
		},
		{
			Keys: bson.D{{Key: "category", Value: 1}},
		},
	})
	return err
}

// MongoDatabase methods
func (mdb *MongoDatabase) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mdb.Client.Disconnect(ctx)
}

func (mdb *MongoDatabase) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return mdb.Client.Ping(ctx, nil)
}

func (mdb *MongoDatabase) GetType() DatabaseType {
	return MongoDB
}

// NewPostgresDB establishes a connection to PostgreSQL
func NewPostgresDB() (*PostgresDB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL")

	return &PostgresDB{DB: db}, nil
}

// NewMySQLDB establishes a connection to MySQL
func NewMySQLDB() (*MySQLDB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DB"),
	)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %v", err)
	}

	log.Println("Connected to MySQL")

	return &MySQLDB{DB: db}, nil
}

// PostgresDB methods
func (pdb *PostgresDB) Close() error {
	return pdb.DB.Close()
}

func (pdb *PostgresDB) Ping() error {
	return pdb.DB.Ping()
}

func (pdb *PostgresDB) GetType() DatabaseType {
	return PostgreSQL
}

// MySQLDB methods
func (mdb *MySQLDB) Close() error {
	return mdb.DB.Close()
}

func (mdb *MySQLDB) Ping() error {
	return mdb.DB.Ping()
}

func (mdb *MySQLDB) GetType() DatabaseType {
	return MySQL
}
