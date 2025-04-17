package main

import (
    "log"
    "net"
    "google.golang.org/grpc"
)

// This is a skeleton for a gRPC microservice. Define your proto and implement your service here.

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    // Register your service here
    log.Println("gRPC server listening on :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
