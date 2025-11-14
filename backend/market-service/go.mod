module lfg/market-service

go 1.24.3

require (
	lfg/shared v0.0.0
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.2
)

replace lfg/shared => ../shared
