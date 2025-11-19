module lfg/api-gateway

go 1.24.3

require lfg/shared v0.0.0

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
)

replace lfg/shared => ../shared
