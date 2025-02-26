package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hosue-of-segments/domin-segment/proto"
	"hosue-of-segments/monitor"

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

func periodicalHealthCheck(ctx context.Context) {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	broker := []string{"localhost:9092"}
	producer, err := monitor.NewProducer(broker)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer producer.Close()

	// 	// Start an infinite loop that calls the API every time the ticker ticks
	for {
		select {
		case <-ticker.C:
			err = producer.SendMessage()
			if err != nil {
				log.Printf("Error sending message: %v", err)
			}

		case <-ctx.Done():
			log.Println("Shutting down health check...")
			return
		}
	}
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go periodicalHealthCheck(ctx)

	gracefulShutdown(s)
}
