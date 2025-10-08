package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Placeholder for gRPC client to connect to the Matching Engine

	// Placeholder for NATS subscriber to listen for TradeExecuted events
	// for conditional order (Stop, Stop-Limit) activation.

	// API handlers
	http.HandleFunc("/orders/place", placeOrderHandler)
	http.HandleFunc("/orders/cancel", cancelOrderHandler)
	http.HandleFunc("/orders/status", statusOrderHandler)

	fmt.Println("Order service listening on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func placeOrderHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Decode and validate the incoming order request.
	// 2. Check user's wallet for sufficient funds (gRPC call to Wallet Service).
	// 3. If Market or Limit order, forward to Matching Engine (gRPC call).
	// 4. If Stop or Stop-Limit order, store it locally and wait for trigger.
	fmt.Fprintln(w, "Place order endpoint")
}

func cancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Logic for cancelling an open order.
	fmt.Fprintln(w, "Cancel order endpoint")
}

func statusOrderHandler(w http.ResponseWriter, r *http.Request) {
	// Logic for retrieving the status of an order.
	fmt.Fprintln(w, "Order status endpoint")
}
