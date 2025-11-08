package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Kevin-Kurka/LFG/backend/common/auth"
	"github.com/Kevin-Kurka/LFG/backend/common/database"
	"github.com/Kevin-Kurka/LFG/backend/sportsbook-service/handlers"
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

	log.Printf("Sportsbook service listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, corsHandler.Handler(router)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
