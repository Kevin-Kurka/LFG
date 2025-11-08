module github.com/Kevin-Kurka/LFG/backend/notification-service

go 1.21

replace github.com/Kevin-Kurka/LFG/backend/common => ../common

require (
	github.com/google/uuid v1.5.0
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.1
	github.com/rs/cors v1.10.1
)

require golang.org/x/net v0.17.0 // indirect
