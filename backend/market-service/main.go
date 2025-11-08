package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Kevin-Kurka/LFG/backend/common/auth"
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

	// Public market routes (read-only)
	router.HandleFunc("/markets", handlers.ListMarkets).Methods("GET")
	router.HandleFunc("/markets/{id}", handlers.GetMarket).Methods("GET")
	router.HandleFunc("/markets/{id}/orderbook", handlers.GetOrderBook).Methods("GET")

	// Admin routes (require admin API key)
	adminRouter := router.PathPrefix("/markets").Subrouter()
	adminRouter.Use(auth.AdminAPIKeyMiddleware)
	adminRouter.HandleFunc("", handlers.CreateMarket).Methods("POST")
	adminRouter.HandleFunc("/{id}", handlers.UpdateMarket).Methods("PUT")
	adminRouter.HandleFunc("/{id}/resolve", handlers.ResolveMarket).Methods("POST")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	// Get allowed origins from environment
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	origins := []string{"http://localhost:3000"} // Default for development
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Admin-API-Key"},
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
