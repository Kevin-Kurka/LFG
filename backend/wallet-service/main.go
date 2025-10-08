package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Placeholder for API handlers
	http.HandleFunc("/balance", balanceHandler)
	http.HandleFunc("/transactions", transactionsHandler)

	fmt.Println("Wallet service listening on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func balanceHandler(w http.ResponseWriter, r *http.Request) {
	// Logic for retrieving wallet balance
	fmt.Fprintln(w, "Wallet balance endpoint")
}

func transactionsHandler(w http.ResponseWriter, r *http.Request) {
	// Logic for retrieving transaction history
	fmt.Fprintln(w, "Transaction history endpoint")
}
