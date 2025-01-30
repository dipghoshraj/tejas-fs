package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"hosue-of-segments/domin-segment/proto"

	domainsegment "hosue-of-segments/domin-segment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func gracefulShutdown(server *grpc.Server) {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	// Attempt graceful shutdown
	server.Stop()

}

func main() {
	fmt.Println("Welcome to unit commander")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Register the HelloService server
	proto.RegisterStorageBoxServiceServer(s, &domainsegment.StorageServer{})
	reflection.Register(s)

	// Start the server
	fmt.Println("Server is running on port 50051")
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	gracefulShutdown(s)
}
