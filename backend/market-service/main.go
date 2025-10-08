package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// API handlers
	http.HandleFunc("/markets", listMarketsHandler)
	http.HandleFunc("/markets/detail", marketDetailHandler)
	http.HandleFunc("/markets/orderbook", orderbookHandler)

	fmt.Println("Market service listening on port 8083...")
	log.Fatal(http.ListenAndServe(":8083", nil))
}

func listMarketsHandler(w http.ResponseWriter, r *http.Request) {
	// Logic for listing all available markets, with filtering.
	fmt.Fprintln(w, "List markets endpoint")
}

func marketDetailHandler(w http.ResponseWriter, r *http.Request) {
	// Logic for retrieving the detailed information for a single market.
	fmt.Fprintln(w, "Market detail endpoint")
}

func orderbookHandler(w http.ResponseWriter, r *http.Request) {
	// Logic for retrieving the current order book for a market.
	fmt.Fprintln(w, "Market order book endpoint")
}
