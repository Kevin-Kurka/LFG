module github.com/Kevin-Kurka/LFG/backend/credit-exchange-service

go 1.21

replace github.com/Kevin-Kurka/LFG/backend/common => ../common

require (
	github.com/Kevin-Kurka/LFG/backend/common v0.0.0
	github.com/google/uuid v1.5.0
	github.com/gorilla/mux v1.8.1
	github.com/rs/cors v1.10.1
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
)
