package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Kevin-Kurka/LFG/backend/common/auth"
	"github.com/Kevin-Kurka/LFG/backend/common/database"
	"github.com/Kevin-Kurka/LFG/backend/order-service/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	router := mux.NewRouter()
	router.Use(auth.AuthMiddleware)

	router.HandleFunc("/orders/place", handlers.PlaceOrder).Methods("POST")
	router.HandleFunc("/orders/cancel", handlers.CancelOrder).Methods("POST")
	router.HandleFunc("/orders", handlers.GetOrders).Methods("GET")
	router.HandleFunc("/orders/{id}", handlers.GetOrder).Methods("GET")

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
		port = "8085"
	}

	log.Printf("Order service listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, corsHandler.Handler(router)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
