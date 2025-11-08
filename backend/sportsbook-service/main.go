package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Kevin-Kurka/LFG/backend/common/auth"
	"github.com/Kevin-Kurka/LFG/backend/common/database"
	"github.com/Kevin-Kurka/LFG/backend/sportsbook-service/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	log.Println("Starting Sportsbook Service...")

	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	router := mux.NewRouter()
	router.Use(auth.AuthMiddleware)

	// Account management
	router.HandleFunc("/sportsbook/accounts", handlers.LinkAccount).Methods("POST")
	router.HandleFunc("/sportsbook/accounts", handlers.GetLinkedAccounts).Methods("GET")
	router.HandleFunc("/sportsbook/accounts/{id}", handlers.DeleteLinkedAccount).Methods("DELETE")

	// Sports events and odds
	router.HandleFunc("/sportsbook/events", handlers.GetSportsEvents).Methods("GET")
	router.HandleFunc("/sportsbook/events/{id}", handlers.GetEventDetails).Methods("GET")

	// Arbitrage and hedge opportunities
	router.HandleFunc("/sportsbook/arbitrage", handlers.GetArbitrageOpportunities).Methods("GET")
	router.HandleFunc("/sportsbook/hedges", handlers.GetHedgeOpportunities).Methods("GET")

	// Bet tracking
	router.HandleFunc("/sportsbook/bets", handlers.TrackBet).Methods("POST")
	router.HandleFunc("/sportsbook/bets", handlers.GetUserBets).Methods("GET")

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
		port = "8088"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      corsHandler.Handler(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Sportsbook Service listening on port %s\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
