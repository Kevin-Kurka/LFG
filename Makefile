.PHONY: help build test clean run stop migrate-up migrate-down backup restore docker-build docker-up docker-down k8s-deploy k8s-delete lint format

# Default target
help:
	@echo "LFG Platform - Makefile Commands"
	@echo ""
	@echo "Development:"
	@echo "  make run              - Start all services locally"
	@echo "  make stop             - Stop all services"
	@echo "  make test             - Run all tests"
	@echo "  make lint             - Run linting"
	@echo "  make format           - Format code"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build     - Build all Docker images"
	@echo "  make docker-up        - Start Docker Compose stack"
	@echo "  make docker-down      - Stop Docker Compose stack"
	@echo "  make docker-prod      - Start production Docker stack"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-up       - Run database migrations"
	@echo "  make migrate-down     - Rollback last migration"
	@echo "  make migrate-create   - Create new migration (NAME=migration_name)"
	@echo "  make backup           - Backup database"
	@echo "  make restore          - Restore database (FILE=backup_file)"
	@echo ""
	@echo "Kubernetes:"
	@echo "  make k8s-deploy       - Deploy to Kubernetes"
	@echo "  make k8s-delete       - Delete Kubernetes resources"
	@echo "  make k8s-validate     - Validate Kubernetes manifests"
	@echo ""
	@echo "Other:"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make deps             - Install dependencies"

# Build all services
build:
	@echo "Building all services..."
	cd backend && go build -o ../bin/api-gateway ./api-gateway
	cd backend && go build -o ../bin/user-service ./user-service
	cd backend && go build -o ../bin/wallet-service ./wallet-service
	cd backend && go build -o ../bin/market-service ./market-service
	cd backend && go build -o ../bin/order-service ./order-service
	cd backend && go build -o ../bin/matching-engine ./matching-engine
	cd backend && go build -o ../bin/credit-exchange-service ./credit-exchange-service
	cd backend && go build -o ../bin/notification-service ./notification-service
	cd backend && go build -o ../bin/sportsbook-service ./sportsbook-service
	cd backend && go build -o ../bin/crypto-service ./crypto-service
	@echo "Build complete!"

# Run tests
test:
	@echo "Running tests..."
	cd backend && go test -v -race -coverprofile=coverage.out ./...
	@echo "Tests complete!"

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	cd backend && go test -v -race -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: backend/coverage.html"

# Run linting
lint:
	@echo "Running linting..."
	cd backend && golangci-lint run ./...
	@echo "Linting complete!"

# Format code
format:
	@echo "Formatting code..."
	cd backend && gofmt -s -w .
	cd backend && goimports -w .
	@echo "Formatting complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f backend/coverage.out backend/coverage.html
	@echo "Clean complete!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	cd backend && go mod download
	cd backend && go mod verify
	@echo "Dependencies installed!"

# Database migrations
migrate-up:
	@echo "Running migrations..."
	./database/migrate.sh up

migrate-down:
	@echo "Rolling back migration..."
	./database/migrate.sh down

migrate-create:
ifndef NAME
	$(error NAME is required. Usage: make migrate-create NAME=migration_name)
endif
	@echo "Creating migration: $(NAME)"
	./database/migrate.sh create $(NAME)

# Database backup and restore
backup:
	@echo "Backing up database..."
	./scripts/backup-db.sh
	@echo "Backup complete!"

restore:
ifndef FILE
	$(error FILE is required. Usage: make restore FILE=backup_file)
endif
	@echo "Restoring database from $(FILE)..."
	./scripts/restore-db.sh $(FILE)

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker-compose build
	@echo "Docker build complete!"

docker-up:
	@echo "Starting Docker Compose stack..."
	docker-compose up -d
	@echo "Docker stack started!"
	@echo "Services:"
	@echo "  - API Gateway: http://localhost:8000"
	@echo "  - Frontend: http://localhost:3000"
	@echo "  - Admin Panel: http://localhost:3001"

docker-down:
	@echo "Stopping Docker Compose stack..."
	docker-compose down
	@echo "Docker stack stopped!"

docker-prod:
	@echo "Starting production Docker stack..."
	docker-compose -f docker-compose.prod.yml up -d
	@echo "Production stack started!"

docker-logs:
	docker-compose logs -f

# Kubernetes commands
k8s-deploy:
	@echo "Deploying to Kubernetes..."
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/configmap.yaml
	kubectl apply -f k8s/secrets.yaml
	kubectl apply -f k8s/postgres-statefulset.yaml
	kubectl apply -f k8s/redis-deployment.yaml
	kubectl apply -f k8s/all-services.yaml
	kubectl apply -f k8s/monitoring.yaml
	@echo "Deployment complete!"

k8s-delete:
	@echo "Deleting Kubernetes resources..."
	kubectl delete -f k8s/all-services.yaml --ignore-not-found
	kubectl delete -f k8s/monitoring.yaml --ignore-not-found
	kubectl delete -f k8s/redis-deployment.yaml --ignore-not-found
	kubectl delete -f k8s/postgres-statefulset.yaml --ignore-not-found
	@echo "Resources deleted!"

k8s-validate:
	@echo "Validating Kubernetes manifests..."
	kubectl apply -f k8s/ --dry-run=client
	@echo "Validation complete!"

k8s-status:
	@echo "Kubernetes cluster status..."
	kubectl get all -n lfg
	kubectl get all -n monitoring

# Run locally
run:
	@echo "Starting services locally..."
	docker-compose up -d postgres redis
	@echo "Waiting for database..."
	sleep 5
	make migrate-up
	@echo "Starting services..."
	./bin/api-gateway &
	./bin/user-service &
	./bin/wallet-service &
	./bin/market-service &
	@echo "Services started!"

stop:
	@echo "Stopping services..."
	pkill -f "bin/api-gateway" || true
	pkill -f "bin/user-service" || true
	pkill -f "bin/wallet-service" || true
	pkill -f "bin/market-service" || true
	docker-compose down
	@echo "Services stopped!"

# Development helpers
dev-setup:
	@echo "Setting up development environment..."
	make deps
	cp .env.example .env
	@echo "Please edit .env file with your configuration"
	make docker-up
	sleep 10
	make migrate-up
	@echo "Development setup complete!"

# Production deployment
prod-deploy:
	@echo "Deploying to production..."
	@echo "Building images..."
	make docker-build
	@echo "Pushing images..."
	docker-compose -f docker-compose.prod.yml push
	@echo "Deploying to Kubernetes..."
	make k8s-deploy
	@echo "Production deployment complete!"
