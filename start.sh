#!/bin/bash

# LFG Platform Quick Start Script
# This script sets up and starts the entire LFG platform

set -e  # Exit on error

echo "🚀 LFG Platform - Quick Start"
echo "=============================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker is not installed. Please install Docker first.${NC}"
    echo "Visit: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}❌ Docker Compose is not installed. Please install Docker Compose first.${NC}"
    echo "Visit: https://docs.docker.com/compose/install/"
    exit 1
fi

echo -e "${GREEN}✓ Docker and Docker Compose are installed${NC}"
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${YELLOW}⚠️  .env file not found. Creating from .env.example...${NC}"
    cp .env.example .env

    echo ""
    echo -e "${YELLOW}⚠️  IMPORTANT: You should generate secure secrets!${NC}"
    echo ""
    echo "Generate JWT Secret:"
    echo "  openssl rand -base64 64"
    echo ""
    echo "Generate Encryption Key:"
    echo "  openssl rand -base64 32"
    echo ""
    echo "Then update .env file with these values."
    echo ""
    read -p "Press Enter to continue with default values (NOT recommended for production) or Ctrl+C to exit and configure..."
fi

echo -e "${GREEN}✓ Environment file exists${NC}"
echo ""

# Stop any existing containers
echo "🛑 Stopping any existing containers..."
docker-compose down 2>/dev/null || true
echo ""

# Build and start services
echo "🏗️  Building Docker images..."
docker-compose build --parallel

echo ""
echo "🚀 Starting all services..."
docker-compose up -d

echo ""
echo "⏳ Waiting for services to be healthy..."
sleep 10

# Check service health
echo ""
echo "🏥 Checking service health..."

services=(
  "http://localhost:8000|API Gateway"
  "http://localhost:8080|User Service"
  "http://localhost:8081|Wallet Service"
  "http://localhost:8083|Market Service"
)

all_healthy=true

for service_info in "${services[@]}"; do
  IFS='|' read -r url name <<< "$service_info"
  if curl -f -s "${url}/health" > /dev/null 2>&1; then
    echo -e "${GREEN}✓ ${name} - healthy${NC}"
  else
    echo -e "${RED}✗ ${name} - unhealthy${NC}"
    all_healthy=false
  fi
done

echo ""

if [ "$all_healthy" = true ]; then
  echo -e "${GREEN}✅ All services are healthy!${NC}"
else
  echo -e "${YELLOW}⚠️  Some services are not healthy yet. They may still be starting up.${NC}"
  echo "Check logs with: docker-compose logs -f"
fi

echo ""
echo "=============================="
echo -e "${GREEN}🎉 LFG Platform is running!${NC}"
echo "=============================="
echo ""
echo "Access URLs:"
echo "  📱 Frontend:     http://localhost:3000"
echo "  ⚙️  Admin Panel:  http://localhost:3001"
echo "  🔌 API Gateway:  http://localhost:8000"
echo ""
echo "Useful Commands:"
echo "  View logs:       docker-compose logs -f"
echo "  Stop services:   docker-compose down"
echo "  Restart:         docker-compose restart"
echo "  Database:        docker exec -it lfg-postgres psql -U lfguser -d lfg"
echo ""
echo -e "${RED}⚠️  LEGAL WARNING:${NC}"
echo "This software is for demonstration purposes only."
echo "Read LEGAL_DISCLAIMER.md before using in any production capacity."
echo ""
echo "Happy trading! 🎲"
