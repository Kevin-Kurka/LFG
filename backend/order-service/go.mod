module github.com/Kevin-Kurka/LFG/backend/order-service

go 1.21

replace github.com/Kevin-Kurka/LFG/backend/common => ../common

require (
	github.com/Kevin-Kurka/LFG/backend/common v0.0.0
	github.com/google/uuid v1.5.0
	github.com/gorilla/mux v1.8.1
	github.com/jackc/pgx/v5 v5.5.1
	github.com/rs/cors v1.10.1
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
