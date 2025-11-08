package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Kevin-Kurka/LFG/backend/common/database"
	"github.com/Kevin-Kurka/LFG/backend/matching-engine/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	router := mux.NewRouter()

	// Internal endpoints for order service
	router.HandleFunc("/submit", handlers.SubmitOrder).Methods("POST")
	router.HandleFunc("/cancel", handlers.CancelOrder).Methods("POST")
	router.HandleFunc("/orderbook/{contractId}", handlers.GetOrderBook).Methods("GET")

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
		port = "8084"
	}

	log.Printf("Matching engine listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, corsHandler.Handler(router)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
