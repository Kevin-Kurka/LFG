package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"

	"lfg/matching-engine/engine"
	pb "lfg/matching-engine/proto"
)

func main() {
	// Get port from environment or use default
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	// Connect to NATS
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	natsConn, err := nats.Connect(natsURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to NATS: %v", err)
		log.Println("Continuing without NATS (trade events will not be published)")
	} else {
		log.Printf("Connected to NATS at %s", natsURL)
		defer natsConn.Close()
	}

	// Initialize matching engine
	matchingEngine := engine.NewMatchingEngine(natsConn)
	log.Println("Matching engine initialized")

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterMatchingEngineServer(grpcServer, matchingEngine)

	// Start listening
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Start server in goroutine
	go func() {
		log.Printf("Matching engine gRPC server listening on port %s...\n", port)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down matching engine...")
	grpcServer.GracefulStop()
	fmt.Println("Matching engine exited")
}
