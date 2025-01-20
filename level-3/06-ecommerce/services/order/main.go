package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderItem struct {
	ProductID string  `json:"product_id" bson:"product_id"`
	Quantity  int     `json:"quantity" bson:"quantity"`
	Price     float64 `json:"price" bson:"price"`
}

type Order struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    string            `json:"user_id" bson:"user_id"`
	Items     []OrderItem       `json:"items" bson:"items"`
	Total     float64           `json:"total" bson:"total"`
	Status    string            `json:"status" bson:"status"`
	CreatedAt time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" bson:"updated_at"`
}

type OrderService struct {
	collection     *mongo.Collection
	productService string
}

func NewOrderService() (*OrderService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		return nil, err
	}

	db := client.Database(os.Getenv("DB_NAME"))
	collection := db.Collection("orders")

	return &OrderService{
		collection:     collection,
		productService: os.Getenv("PRODUCT_SERVICE_URL"),
	}, nil
}

func (s *OrderService) validateProducts(items []OrderItem) error {
	for _, item := range items {
		resp, err := http.Get(fmt.Sprintf("%s/api/products/%s", s.productService, item.ProductID))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("product %s not found", item.ProductID)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("error validating product %s", item.ProductID)
		}

		var product struct {
			Stock int     `json:"stock"`
			Price float64 `json:"price"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
			return err
		}

		if product.Stock < item.Quantity {
			return fmt.Errorf("insufficient stock for product %s", item.ProductID)
		}
		item.Price = product.Price
	}
	return nil
}

func (s *OrderService) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.validateProducts(order.Items); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.Status = "pending"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// Calculate total
	var total float64
	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}
	order.Total = total

	result, err := s.collection.InsertOne(context.Background(), order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	order.ID = result.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(order)
}

func (s *OrderService) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var order Order
	err = s.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (s *OrderService) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var status struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status.Status,
			"updated_at": time.Now(),
		},
	}

	var order Order
	result := s.collection.FindOneAndUpdate(
		context.Background(),
		bson.M{"_id": id},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if err := result.Decode(&order); err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (s *OrderService) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	filter := bson.M{}
	if userID != "" {
		filter["user_id"] = userID
	}

	cursor, err := s.collection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.Background())

	var orders []Order
	if err := cursor.All(context.Background(), &orders); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}

func main() {
	service, err := NewOrderService()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/orders", service.CreateOrder).Methods("POST")
	r.HandleFunc("/api/orders", service.ListOrders).Methods("GET")
	r.HandleFunc("/api/orders/{id}", service.GetOrder).Methods("GET")
	r.HandleFunc("/api/orders/{id}/status", service.UpdateOrderStatus).Methods("PUT")

	port := "8082"
	fmt.Printf("Order service starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
