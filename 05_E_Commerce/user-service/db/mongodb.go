package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var (
	client *mongo.Client
	db     *mongo.Database
)

// User represents a user in the system
type User struct {
	ID           string    `bson:"_id,omitempty"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password_hash"`
	FirstName    string    `bson:"first_name"`
	LastName     string    `bson:"last_name"`
	CreatedAt    time.Time `bson:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at"`
}

// InitMongoDB initializes the MongoDB connection
func InitMongoDB() error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return fmt.Errorf("MONGODB_URI environment variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "ecommerce"
	}
	db = client.Database(dbName)

	// Create unique index on email
	collection := GetUsersCollection()
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create index: %v", err)
	}

	log.Println("Successfully connected to MongoDB")
	return nil
}

// GetUsersCollection returns the users collection
func GetUsersCollection() *mongo.Collection {
	return db.Collection("users")
}

// CreateUser creates a new user
func CreateUser(ctx context.Context, user *User, password string) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user.PasswordHash = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	collection := GetUsersCollection()
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// GetUser retrieves a user by ID
func GetUser(ctx context.Context, userID string) (*User, error) {
	collection := GetUsersCollection()
	var user User
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	collection := GetUsersCollection()
	var user User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user information
func UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	collection := GetUsersCollection()
	_, err := collection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": updates})
	return err
}

// ValidatePassword validates a user's password
func ValidatePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Close closes the MongoDB connection
func Close() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}
}
