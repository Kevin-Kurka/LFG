package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Kevin-Kurka/LFG/backend/common/database"
	"github.com/Kevin-Kurka/LFG/backend/market-service/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	router := mux.NewRouter()

	// Market routes
	router.HandleFunc("/markets", handlers.ListMarkets).Methods("GET")
	router.HandleFunc("/markets", handlers.CreateMarket).Methods("POST")
	router.HandleFunc("/markets/{id}", handlers.GetMarket).Methods("GET")
	router.HandleFunc("/markets/{id}", handlers.UpdateMarket).Methods("PUT")
	router.HandleFunc("/markets/{id}/resolve", handlers.ResolveMarket).Methods("POST")
	router.HandleFunc("/markets/{id}/orderbook", handlers.GetOrderBook).Methods("GET")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Market service listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, corsHandler.Handler(router)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
