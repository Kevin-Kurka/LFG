package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Kevin-Kurka/LFG/backend/common/auth"
	"github.com/Kevin-Kurka/LFG/backend/common/database"
	"github.com/Kevin-Kurka/LFG/backend/wallet-service/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Create router
	router := mux.NewRouter()

	// Protected routes (require authentication)
	router.Use(auth.AuthMiddleware)
	router.HandleFunc("/balance", handlers.GetBalance).Methods("GET")
	router.HandleFunc("/transactions", handlers.GetTransactions).Methods("GET")

	// Internal route (requires internal API key for service-to-service calls)
	internalRouter := mux.NewRouter()
	internalRouter.Use(auth.InternalAPIKeyMiddleware)
	internalRouter.HandleFunc("/internal/transactions", handlers.CreateTransaction).Methods("POST")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	// CORS configuration
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	origins := []string{"http://localhost:3000"} // Default for development
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Internal-API-Key"},
		AllowCredentials: true,
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Wallet service listening on port %s...", port)

	// Combine routers
	http.Handle("/internal/", http.StripPrefix("/internal", internalRouter))
	http.Handle("/", corsHandler.Handler(router))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
