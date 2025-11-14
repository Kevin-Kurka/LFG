.PHONY: help build test lint fmt check docker-up docker-down docker-restart docker-logs docker-clean db-migrate db-rollback db-seed db-reset db-shell

# Default target
help:
	@echo "LFG Platform - Make Commands"
	@echo ""
	@echo "Development:"
	@echo "  make build          - Build all Go services"
	@echo "  make test           - Run all tests"
	@echo "  make lint           - Run linters"
	@echo "  make fmt            - Format Go code"
	@echo "  make check          - Run tests + lint + fmt"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up      - Start all services in Docker"
	@echo "  make docker-down    - Stop all services"
	@echo "  make docker-restart - Restart all services"
	@echo "  make docker-logs    - View logs from all services"
	@echo "  make docker-clean   - Remove all containers and volumes"
	@echo ""
	@echo "Database:"
	@echo "  make db-migrate     - Run database migrations"
	@echo "  make db-rollback    - Rollback last migration"
	@echo "  make db-seed        - Seed database with test data"
	@echo "  make db-reset       - Reset database (drop and recreate)"
	@echo "  make db-shell       - Connect to database shell"

# Build all Go services
build:
	@echo "Building all services..."
	@cd backend/api-gateway && go build -o ../../bin/api-gateway
	@cd backend/user-service && go build -o ../../bin/user-service
	@cd backend/wallet-service && go build -o ../../bin/wallet-service
	@cd backend/order-service && go build -o ../../bin/order-service
	@cd backend/market-service && go build -o ../../bin/market-service
	@cd backend/credit-exchange-service && go build -o ../../bin/credit-exchange
	@cd backend/notification-service && go build -o ../../bin/notification-service
	@cd backend/matching-engine && go build -o ../../bin/matching-engine
	@echo "Build complete!"

# Run all tests
test:
	@echo "Running tests..."
	@cd backend/shared && go test -v -cover ./...
	@cd backend/api-gateway && go test -v -cover ./...
	@cd backend/user-service && go test -v -cover ./...
	@cd backend/wallet-service && go test -v -cover ./...
	@cd backend/order-service && go test -v -cover ./...
	@cd backend/market-service && go test -v -cover ./...
	@cd backend/credit-exchange-service && go test -v -cover ./...
	@cd backend/notification-service && go test -v -cover ./...
	@cd backend/matching-engine && go test -v -cover ./...
	@echo "Tests complete!"

# Run linters
lint:
	@echo "Running linters..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; exit 1; }
	@cd backend/shared && golangci-lint run
	@cd backend/api-gateway && golangci-lint run
	@cd backend/user-service && golangci-lint run
	@cd backend/wallet-service && golangci-lint run
	@cd backend/order-service && golangci-lint run
	@cd backend/market-service && golangci-lint run
	@cd backend/credit-exchange-service && golangci-lint run
	@cd backend/notification-service && golangci-lint run
	@cd backend/matching-engine && golangci-lint run
	@echo "Linting complete!"

# Format Go code
fmt:
	@echo "Formatting code..."
	@cd backend && gofmt -s -w .
	@echo "Formatting complete!"

# Run all checks
check: fmt lint test

# Docker commands
docker-up:
	@echo "Starting all services..."
	@docker-compose up -d
	@echo "Services started! Check status with: docker-compose ps"

docker-down:
	@echo "Stopping all services..."
	@docker-compose down
	@echo "Services stopped!"

docker-restart:
	@echo "Restarting all services..."
	@docker-compose restart
	@echo "Services restarted!"

docker-logs:
	@docker-compose logs -f

docker-clean:
	@echo "Cleaning up Docker resources..."
	@docker-compose down -v --remove-orphans
	@echo "Cleanup complete!"

# Database commands
db-migrate:
	@echo "Running database migrations..."
	@docker exec -i lfg-postgres psql -U lfg -d lfg < database/migrations/001_enhanced_schema.up.sql
	@echo "Migrations complete!"

db-rollback:
	@echo "Rolling back last migration..."
	@docker exec -i lfg-postgres psql -U lfg -d lfg < database/migrations/001_enhanced_schema.down.sql
	@echo "Rollback complete!"

db-seed:
	@echo "Seeding database..."
	@docker exec -i lfg-postgres psql -U lfg -d lfg < database/seed.sql
	@echo "Seed complete!"

db-reset: db-rollback db-migrate db-seed
	@echo "Database reset complete!"

db-shell:
	@docker exec -it lfg-postgres psql -U lfg -d lfg

# Install dependencies
deps:
	@echo "Installing Go dependencies..."
	@cd backend/shared && go mod download
	@cd backend/api-gateway && go mod download
	@cd backend/user-service && go mod download
	@cd backend/wallet-service && go mod download
	@cd backend/order-service && go mod download
	@cd backend/market-service && go mod download
	@cd backend/credit-exchange-service && go mod download
	@cd backend/notification-service && go mod download
	@cd backend/matching-engine && go mod download
	@echo "Dependencies installed!"

# Generate code (protobuf, etc.)
generate:
	@echo "Generating code..."
	@cd backend/matching-engine && protoc --go_out=. --go-grpc_out=. proto/*.proto
	@echo "Code generation complete!"
