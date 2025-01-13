package handler

import (
	"context"
	"fmt"
	"time"

	"product-service/db"
	pb "product-service/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductHandler struct {
	pb.UnimplementedProductServiceServer
	mongodb *db.MongoDB
	redis   *db.RedisClient
}

func NewProductHandler(mongodb *db.MongoDB, redis *db.RedisClient) *ProductHandler {
	return &ProductHandler{
		mongodb: mongodb,
		redis:   redis,
	}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	// Validate request
	if req.Name == "" || req.Price <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid product details")
	}

	// Create product document
	product := bson.M{
		"name":        req.Name,
		"description": req.Description,
		"price":       req.Price,
		"stock":       req.Stock,
		"category":    req.Category,
		"images":      req.Images,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}

	// Insert into MongoDB
	result, err := h.mongodb.Collection("products").InsertOne(ctx, product)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	// Get the inserted product
	var created bson.M
	err = h.mongodb.Collection("products").FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&created)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch created product: %v", err)
	}

	// Convert to protobuf response
	return &pb.ProductResponse{
		Product: convertToProtoProduct(created),
	}, nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	// Try to get from cache first
	if product, err := h.getFromCache(req.Id); err == nil {
		return &pb.ProductResponse{Product: product}, nil
	}

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID")
	}

	// Get from MongoDB
	var product bson.M
	err = h.mongodb.Collection("products").FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}

	// Convert to protobuf
	protoProduct := convertToProtoProduct(product)

	// Cache the result
	h.cacheProduct(protoProduct)

	return &pb.ProductResponse{Product: protoProduct}, nil
}

func (h *ProductHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	// Set default values
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// Build query
	filter := bson.M{}
	if req.Category != "" {
		filter["category"] = req.Category
	}

	// Calculate skip
	skip := (req.Page - 1) * req.Limit

	// Get total count
	total, err := h.mongodb.Collection("products").CountDocuments(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to count products: %v", err)
	}

	// Get products
	cursor, err := h.mongodb.Collection("products").Find(ctx, filter,
		&options.FindOptions{
			Skip:  &skip,
			Limit: &req.Limit,
			Sort:  bson.M{"created_at": -1},
		})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch products: %v", err)
	}
	defer cursor.Close(ctx)

	// Decode products
	var products []bson.M
	if err := cursor.All(ctx, &products); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to decode products: %v", err)
	}

	// Convert to protobuf
	protoProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = convertToProtoProduct(product)
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(total),
	}, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID")
	}

	// Build update document
	update := bson.M{
		"$set": bson.M{
			"name":        req.Name,
			"description": req.Description,
			"price":       req.Price,
			"stock":       req.Stock,
			"category":    req.Category,
			"images":      req.Images,
			"updated_at":  time.Now(),
		},
	}

	// Update in MongoDB
	result := h.mongodb.Collection("products").FindOneAndUpdate(
		ctx,
		bson.M{"_id": objectID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	// Check if product exists
	var updated bson.M
	if err := result.Decode(&updated); err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}

	// Convert to protobuf
	protoProduct := convertToProtoProduct(updated)

	// Update cache
	h.cacheProduct(protoProduct)

	return &pb.ProductResponse{Product: protoProduct}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID")
	}

	// Delete from MongoDB
	result, err := h.mongodb.Collection("products").DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}

	if result.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}

	// Delete from cache
	h.redis.Del(ctx, fmt.Sprintf("product:%s", req.Id))

	return &pb.DeleteProductResponse{Success: true}, nil
}

// Helper functions for caching
func (h *ProductHandler) getFromCache(id string) (*pb.Product, error) {
	data, err := h.redis.Get(context.Background(), fmt.Sprintf("product:%s", id))
	if err != nil {
		return nil, err
	}

	var product pb.Product
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (h *ProductHandler) cacheProduct(product *pb.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	return h.redis.Set(context.Background(),
		fmt.Sprintf("product:%s", product.Id),
		string(data),
		time.Hour*24,
	)
}

// Helper function to convert MongoDB document to protobuf
func convertToProtoProduct(doc bson.M) *pb.Product {
	return &pb.Product{
		Id:          doc["_id"].(primitive.ObjectID).Hex(),
		Name:        doc["name"].(string),
		Description: doc["description"].(string),
		Price:       doc["price"].(float64),
		Stock:       int32(doc["stock"].(int64)),
		Category:    doc["category"].(string),
		Images:      convertToStringSlice(doc["images"]),
		CreatedAt:   doc["created_at"].(time.Time).Format(time.RFC3339),
		UpdatedAt:   doc["updated_at"].(time.Time).Format(time.RFC3339),
	}
}

func convertToStringSlice(v interface{}) []string {
	if v == nil {
		return []string{}
	}
	raw := v.(primitive.A)
	result := make([]string, len(raw))
	for i, v := range raw {
		result[i] = v.(string)
	}
	return result
}
