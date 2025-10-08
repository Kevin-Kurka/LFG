package main

import (
	"fmt"
	"log"
	"net/http"
	// "github.com/gorilla/websocket"
	// "github.com/nats-io/nats.go"
)

// Placeholder for the WebSocket upgrader
// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true // Allow all connections for now
// 	},
// }

func main() {
	// 1. Placeholder for NATS Connection
	// Connect to NATS and subscribe to topics like "trades.executed" and "orders.updated".
	// When a message is received, it will be forwarded to the relevant WebSocket clients.
	fmt.Println("Placeholder for NATS subscriber setup.")

	// 2. Setup WebSocket handler
	http.HandleFunc("/ws", handleWebSocket)

	// 3. Start the server
	fmt.Println("Notification service (WebSocket) listening on port 8085...")
	log.Fatal(http.ListenAndServe(":8085", nil))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// This is where the HTTP connection is upgraded to a WebSocket connection.
	// conn, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	log.Println("Failed to upgrade connection:", err)
	// 	return
	// }
	// defer conn.Close()

	// Logic for managing client connections, subscriptions to specific market data,
	// and pushing messages received from NATS to the client.
	fmt.Fprintln(w, "WebSocket connection endpoint")
}
