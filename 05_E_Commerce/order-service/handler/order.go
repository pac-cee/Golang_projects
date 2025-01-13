package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"order-service/db"
	pb "order-service/proto"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	pb.UnimplementedOrderServiceServer
	productClient pb.ProductServiceClient
	userClient    pb.UserServiceClient
}

func NewOrderServer(productClient pb.ProductServiceClient, userClient pb.UserServiceClient) *OrderServer {
	return &OrderServer{
		productClient: productClient,
		userClient:    userClient,
	}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	// Validate user exists
	userResp, err := s.userClient.GetUser(ctx, &pb.GetUserRequest{UserId: req.UserId})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user: %v", err)
	}

	// Calculate total amount and validate products
	var totalAmount float64
	for _, productID := range req.ProductIds {
		product, err := s.productClient.GetProduct(ctx, &pb.GetProductRequest{ProductId: productID})
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid product %s: %v", productID, err)
		}
		totalAmount += product.Price
	}

	// Create order in database
	order := &db.Order{
		UserID:      userResp.Id,
		ProductIDs:  req.ProductIds,
		TotalAmount: totalAmount,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.CreateOrder(ctx, order); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	// Convert to protobuf response
	createdAt, err := ptypes.TimestampProto(order.CreatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert timestamp: %v", err)
	}

	updatedAt, err := ptypes.TimestampProto(order.UpdatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert timestamp: %v", err)
	}

	return &pb.Order{
		Id:          order.ID,
		UserId:      order.UserID,
		ProductIds:  order.ProductIDs,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	order, err := db.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	createdAt, err := ptypes.TimestampProto(order.CreatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert timestamp: %v", err)
	}

	updatedAt, err := ptypes.TimestampProto(order.UpdatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert timestamp: %v", err)
	}

	return &pb.Order{
		Id:          order.ID,
		UserId:      order.UserID,
		ProductIds:  order.ProductIDs,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.Order, error) {
	// Validate order exists
	order, err := db.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	// Update status
	if err := db.UpdateOrderStatus(ctx, req.OrderId, req.Status); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}

	order.Status = req.Status
	order.UpdatedAt = time.Now()

	createdAt, err := ptypes.TimestampProto(order.CreatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert timestamp: %v", err)
	}

	updatedAt, err := ptypes.TimestampProto(order.UpdatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert timestamp: %v", err)
	}

	return &pb.Order{
		Id:          order.ID,
		UserId:      order.UserID,
		ProductIds:  order.ProductIDs,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (s *OrderServer) GetUserOrders(ctx context.Context, req *pb.GetUserOrdersRequest) (*pb.GetUserOrdersResponse, error) {
	// Validate user exists
	_, err := s.userClient.GetUser(ctx, &pb.GetUserRequest{UserId: req.UserId})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user: %v", err)
	}

	// Get orders from database
	orders, err := db.GetUserOrders(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user orders: %v", err)
	}

	// Convert to protobuf response
	var pbOrders []*pb.Order
	for _, order := range orders {
		createdAt, err := ptypes.TimestampProto(order.CreatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert timestamp: %v", err)
		}

		updatedAt, err := ptypes.TimestampProto(order.UpdatedAt)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert timestamp: %v", err)
		}

		pbOrders = append(pbOrders, &pb.Order{
			Id:          order.ID,
			UserId:      order.UserID,
			ProductIds:  order.ProductIDs,
			TotalAmount: order.TotalAmount,
			Status:      order.Status,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}

	return &pb.GetUserOrdersResponse{
		Orders: pbOrders,
	}, nil
}
