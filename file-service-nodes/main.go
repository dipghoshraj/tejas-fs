package main

import (
	"fmt"
	"log"
	"net"

	"github.com/dipghoshraj/media-service/file-service-nodes/domain"
	"github.com/dipghoshraj/media-service/file-service-nodes/domain/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// func main() {
// 	fmt.Println("Hello, World! From file-service-nodes")
// }

func main() {
	// Set up the listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register the HelloService server
	proto.RegisterStorageBoxServer(s, &domain.StorageServer{})

	// Register reflection service on gRPC server (for testing/inspecting)
	reflection.Register(s)

	// Start the server
	fmt.Println("Server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
