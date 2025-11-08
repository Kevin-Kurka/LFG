package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"lfg/backend/common/auth"
	"lfg/backend/crypto-service/exchange"
	"lfg/backend/crypto-service/handlers"
	"lfg/backend/crypto-service/workers"
)

func main() {
	log.Println("Starting Crypto Service...")

	// Connect to database
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database successfully")

	// Initialize exchange rate service
	rateUpdateInterval := 60 * time.Second // Update rates every 60 seconds
	rateService := exchange.NewRateService(db, rateUpdateInterval)
	rateService.Start()
	defer rateService.Stop()
	log.Println("Exchange rate service started")

	// Initialize deposit monitor
	depositScanInterval := 30 * time.Second // Scan for deposits every 30 seconds
	depositMonitor := workers.NewDepositMonitor(db, rateService, depositScanInterval)
	go depositMonitor.Start()
	defer depositMonitor.Stop()
	log.Println("Deposit monitor started")

	// Initialize withdrawal processor
	withdrawalProcessInterval := 60 * time.Second // Process withdrawals every 60 seconds
	withdrawalProcessor := workers.NewWithdrawalProcessor(db, withdrawalProcessInterval)
	go withdrawalProcessor.Start()
	defer withdrawalProcessor.Stop()
	log.Println("Withdrawal processor started")

	// Initialize handlers
	cryptoHandler := handlers.NewCryptoHandler(db, rateService)

	// Setup router
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// API routes (protected with JWT middleware)
	api := router.PathPrefix("/crypto").Subrouter()
	api.Use(auth.JWTMiddleware)

	// Wallet endpoints
	api.HandleFunc("/wallets", cryptoHandler.GetWallets).Methods("GET")
	api.HandleFunc("/wallets/{currency}", cryptoHandler.CreateWallet).Methods("POST")

	// Deposit endpoints
	api.HandleFunc("/deposits", cryptoHandler.GetDeposits).Methods("GET")
	api.HandleFunc("/deposits/pending", cryptoHandler.GetPendingDeposits).Methods("GET")
	api.HandleFunc("/deposits/simulate", cryptoHandler.SimulateDeposit).Methods("POST")

	// Withdrawal endpoints
	api.HandleFunc("/withdraw", cryptoHandler.RequestWithdrawal).Methods("POST")
	api.HandleFunc("/withdrawals", cryptoHandler.GetWithdrawals).Methods("GET")

	// Exchange rate endpoints
	api.HandleFunc("/rates", cryptoHandler.GetExchangeRates).Methods("GET")
	api.HandleFunc("/convert", cryptoHandler.ConvertAmount).Methods("GET")

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Crypto Service listening on port %s\n", port)
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
