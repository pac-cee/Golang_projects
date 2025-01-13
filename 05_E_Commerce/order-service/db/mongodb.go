package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	db     *mongo.Database
)

// Order represents an order in the system
type Order struct {
	ID          string    `bson:"_id,omitempty"`
	UserID      string    `bson:"user_id"`
	ProductIDs  []string  `bson:"product_ids"`
	TotalAmount float64   `bson:"total_amount"`
	Status      string    `bson:"status"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
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

	log.Println("Successfully connected to MongoDB")
	return nil
}

// GetOrdersCollection returns the orders collection
func GetOrdersCollection() *mongo.Collection {
	return db.Collection("orders")
}

// CreateOrder creates a new order
func CreateOrder(ctx context.Context, order *Order) error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	
	collection := GetOrdersCollection()
	_, err := collection.InsertOne(ctx, order)
	return err
}

// GetOrder retrieves an order by ID
func GetOrder(ctx context.Context, orderID string) (*Order, error) {
	collection := GetOrdersCollection()
	var order Order
	err := collection.FindOne(ctx, map[string]string{"_id": orderID}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateOrderStatus updates the status of an order
func UpdateOrderStatus(ctx context.Context, orderID, status string) error {
	collection := GetOrdersCollection()
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		},
	}
	_, err := collection.UpdateOne(ctx, map[string]string{"_id": orderID}, update)
	return err
}

// GetUserOrders retrieves all orders for a user
func GetUserOrders(ctx context.Context, userID string) ([]*Order, error) {
	collection := GetOrdersCollection()
	cursor, err := collection.Find(ctx, map[string]string{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []*Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
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
