module lfg/order-service

go 1.24.3

require (
	lfg/shared v0.0.0
	lfg/matching-engine v0.0.0
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.2
	google.golang.org/grpc v1.68.1
	google.golang.org/protobuf v1.35.2
)

replace lfg/shared => ../shared
replace lfg/matching-engine => ../matching-engine
