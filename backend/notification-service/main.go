package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"lfg/notification-service/handlers"
)

func main() {
	// Initialize WebSocket hub
	hub := handlers.NewHub()
	log.Println("WebSocket hub initialized")

	// Start hub in background
	go hub.Run()
	log.Println("WebSocket hub running")

	// Connect to NATS and subscribe to trade events
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	natsConn, err := nats.Connect(natsURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to NATS: %v", err)
		log.Println("Continuing without NATS (trade notifications will not work)")
	} else {
		log.Printf("Connected to NATS at %s", natsURL)
		defer natsConn.Close()

		// Subscribe to trades topic
		_, err := natsConn.Subscribe("trades", func(msg *nats.Msg) {
			var tradeEvent map[string]interface{}
			if err := json.Unmarshal(msg.Data, &tradeEvent); err != nil {
				log.Printf("Failed to unmarshal trade event: %v", err)
				return
			}

			log.Printf("Received trade event: %v", tradeEvent)

			// Extract user IDs
			makerUserID, _ := tradeEvent["maker_user_id"].(string)
			takerUserID, _ := tradeEvent["taker_user_id"].(string)

			// Prepare notification message
			notification := map[string]interface{}{
				"type":  "trade",
				"event": tradeEvent,
			}

			notificationJSON, err := json.Marshal(notification)
			if err != nil {
				log.Printf("Failed to marshal notification: %v", err)
				return
			}

			// Broadcast to both users involved in the trade
			if makerUserID != "" {
				hub.BroadcastToUser(makerUserID, notificationJSON)
				log.Printf("Sent trade notification to maker user: %s", makerUserID)
			}

			if takerUserID != "" {
				hub.BroadcastToUser(takerUserID, notificationJSON)
				log.Printf("Sent trade notification to taker user: %s", takerUserID)
			}
		})

		if err != nil {
			log.Printf("Failed to subscribe to trades: %v", err)
		} else {
			log.Println("Subscribed to NATS trades topic")
		}
	}

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/ws", handlers.HandleWebSocket(hub))

	// Create HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("Notification service (WebSocket) listening on port %s...\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server exited")
}
