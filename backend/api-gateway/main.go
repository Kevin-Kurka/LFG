package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// Define the target URLs for the microservices
	userServiceURL, _ := url.Parse("http://localhost:8080")
	walletServiceURL, _ := url.Parse("http://localhost:8081")
	orderServiceURL, _ := url.Parse("http://localhost:8082")
	marketServiceURL, _ := url.Parse("http://localhost:8083")
	creditExchangeURL, _ := url.Parse("http://localhost:8084")
	notificationServiceURL, _ := url.Parse("http://localhost:8085")

	// Create reverse proxies for each service
	userProxy := httputil.NewSingleHostReverseProxy(userServiceURL)
	walletProxy := httputil.NewSingleHostReverseProxy(walletServiceURL)
	orderProxy := httputil.NewSingleHostReverseProxy(orderServiceURL)
	marketProxy := httputil.NewSingleHostReverseProxy(marketServiceURL)
	creditExchangeProxy := httputil.NewSingleHostReverseProxy(creditExchangeURL)
	notificationProxy := httputil.NewSingleHostReverseProxy(notificationServiceURL)

	// Define routing rules
	http.Handle("/register", applyMiddleware(userProxy))
	http.Handle("/login", applyMiddleware(userProxy))
	http.Handle("/profile", applyMiddleware(userProxy))

	http.Handle("/balance", applyMiddleware(walletProxy))
	http.Handle("/transactions", applyMiddleware(walletProxy))

	http.Handle("/orders/", applyMiddleware(orderProxy))

	http.Handle("/markets", applyMiddleware(marketProxy))
	http.Handle("/markets/", applyMiddleware(marketProxy))

	http.Handle("/exchange/", applyMiddleware(creditExchangeProxy))

	http.Handle("/ws", applyMiddleware(notificationProxy))

	fmt.Println("API Gateway listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// applyMiddleware is a placeholder for JWT authentication and rate limiting.
func applyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Placeholder for Rate Limiting Middleware
		// log.Println("Rate limit check would go here")

		// 2. Placeholder for JWT Authentication Middleware
		// log.Println("JWT authentication would go here")

		// Forward the request to the next handler (the reverse proxy)
		next.ServeHTTP(w, r)
	})
}
