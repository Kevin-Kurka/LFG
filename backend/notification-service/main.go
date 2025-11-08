package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Kevin-Kurka/LFG/backend/notification-service/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Start WebSocket hub
	go handlers.GlobalHub.Run()

	router := mux.NewRouter()
	router.HandleFunc("/ws", handlers.HandleWebSocket)

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
		port = "8087"
	}

	log.Printf("Notification service (WebSocket) listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, corsHandler.Handler(router)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
