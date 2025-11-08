# LFG Platform - Complete Deployment Guide

**Production-Ready Deployment Instructions**

---

## 📋 Table of Contents

1. [Prerequisites](#prerequisites)
2. [Local Development Setup](#local-development-setup)
3. [Docker Deployment](#docker-deployment)
4. [Production Deployment](#production-deployment)
5. [Environment Configuration](#environment-configuration)
6. [Database Setup](#database-setup)
7. [Service Health Checks](#service-health-checks)
8. [Monitoring & Logging](#monitoring--logging)
9. [Troubleshooting](#troubleshooting)
10. [Security Hardening](#security-hardening)

---

## 🔧 Prerequisites

### System Requirements

**Minimum:**
- 4 CPU cores
- 8 GB RAM
- 50 GB disk space
- Linux/macOS/Windows with WSL2

**Recommended Production:**
- 8+ CPU cores
- 16+ GB RAM
- 100+ GB SSD storage
- Ubuntu 22.04 LTS or similar

### Required Software

```bash
# Docker & Docker Compose
docker --version  # >= 24.0
docker-compose --version  # >= 2.20

# Go (for development)
go version  # >= 1.24

# Node.js (for frontend development)
node --version  # >= 18.0
npm --version   # >= 9.0

# PostgreSQL Client (optional, for database access)
psql --version  # >= 15.0
```

### Installation

**Ubuntu/Debian:**
```bash
# Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Docker Compose
sudo apt-get update
sudo apt-get install docker-compose-plugin

# Go
wget https://go.dev/dl/go1.24.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

**macOS:**
```bash
# Install Homebrew first: https://brew.sh

# Docker Desktop
brew install --cask docker

# Go
brew install go

# Node.js
brew install node
```

---

## 💻 Local Development Setup

### Quick Start

```bash
# 1. Clone repository
git clone <repository-url>
cd LFG

# 2. Set up environment
cp .env.example .env
# Edit .env with your configuration

# 3. Generate secure keys
# JWT Secret (random 64-character string)
openssl rand -base64 64

# Encryption Key (32 bytes for AES-256)
openssl rand -base64 32

# Update .env with these values

# 4. Start all services
docker-compose up -d

# 5. Wait for services to be healthy
docker-compose ps

# 6. Access the platform
# Frontend: http://localhost:3000
# Admin Panel: http://localhost:3001
# API Gateway: http://localhost:8000
```

### Verify Installation

```bash
# Check all containers are running
docker-compose ps

# Expected output: All services should show "Up" or "healthy"
# - lfg-postgres (healthy)
# - lfg-api-gateway (Up)
# - lfg-user-service (Up)
# - lfg-wallet-service (Up)
# - lfg-order-service (Up)
# - lfg-market-service (Up)
# - lfg-credit-exchange-service (Up)
# - lfg-notification-service (Up)
# - lfg-sportsbook-service (Up)
# - lfg-matching-engine (Up)
# - lfg-frontend (Up)
# - lfg-admin-panel (Up)

# Check logs
docker-compose logs -f api-gateway

# Test API Gateway
curl http://localhost:8000/health

# Test database connection
docker exec -it lfg-postgres psql -U lfguser -d lfg -c "SELECT COUNT(*) FROM users;"
```

---

## 🐳 Docker Deployment

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Load Balancer (nginx)                    │
│                    (Production only)                         │
└─────────────┬──────────────────────────┬────────────────────┘
              │                          │
    ┌─────────▼─────────┐      ┌────────▼────────┐
    │   Frontend:3000   │      │ Admin Panel:3001 │
    └───────────────────┘      └─────────────────┘
              │                          │
              └──────────┬───────────────┘
                         │
              ┌──────────▼──────────┐
              │  API Gateway:8000   │
              └──────────┬──────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
   ┌────▼─────┐   ┌─────▼──────┐  ┌─────▼──────┐
   │  User    │   │  Wallet    │  │  Market    │
   │ :8080    │   │  :8081     │  │  :8083     │
   └────┬─────┘   └─────┬──────┘  └─────┬──────┘
        │               │               │
   ┌────▼─────┐   ┌─────▼──────┐  ┌─────▼──────┐
   │  Order   │   │  Credit    │  │Sportsbook  │
   │ :8082    │   │ Exchange   │  │  :8086     │
   └────┬─────┘   │  :8084     │  └─────┬──────┘
        │         └─────┬──────┘        │
   ┌────▼─────┐   ┌─────▼──────┐  ┌─────▼──────┐
   │ Matching │   │Notification│  │            │
   │ Engine   │   │  :8085     │  │            │
   │ :9000    │   └────────────┘  │            │
   └────┬─────┘                   │            │
        │                         │            │
        └─────────────────────────┴────────────┘
                      │
              ┌───────▼────────┐
              │  PostgreSQL    │
              │    :5432       │
              └────────────────┘
```

### Build All Services

```bash
# Build all Docker images
docker-compose build

# Build specific service
docker-compose build user-service

# Build with no cache (clean build)
docker-compose build --no-cache

# View built images
docker images | grep lfg
```

### Start Services

```bash
# Start all services
docker-compose up -d

# Start specific services
docker-compose up -d postgres user-service wallet-service

# Start with build
docker-compose up -d --build

# View logs
docker-compose logs -f

# View logs for specific service
docker-compose logs -f sportsbook-service

# Follow logs in real-time
docker-compose logs -f --tail=100
```

### Stop Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: Deletes database data)
docker-compose down -v

# Stop without removing containers
docker-compose stop

# Restart services
docker-compose restart

# Restart specific service
docker-compose restart api-gateway
```

---

## 🚀 Production Deployment

### Production Checklist

- [ ] Change all default passwords and secrets
- [ ] Generate strong JWT_SECRET (64+ characters)
- [ ] Generate strong ENCRYPTION_KEY (32 bytes)
- [ ] Configure production database with backups
- [ ] Set up SSL/TLS certificates
- [ ] Configure firewall rules
- [ ] Set up monitoring and alerting
- [ ] Configure log aggregation
- [ ] Set up automated backups
- [ ] Review and sign LEGAL_DISCLAIMER.md
- [ ] Implement rate limiting
- [ ] Configure CORS properly
- [ ] Set up CDN for static assets
- [ ] Perform security audit
- [ ] Load testing
- [ ] Disaster recovery plan

### Environment Variables

Create production `.env` file:

```bash
# Production .env
# NEVER commit this file to version control

# Security (CHANGE THESE!)
JWT_SECRET=<generate-with-openssl-rand-base64-64>
ENCRYPTION_KEY=<generate-with-openssl-rand-base64-32>

# Database (Use managed PostgreSQL in production)
DATABASE_URL=postgresql://lfguser:<strong-password>@db.production.com:5432/lfg?sslmode=require

# API URLs (Use your domain)
API_GATEWAY_URL=https://api.yourdomain.com
FRONTEND_URL=https://yourdomain.com
ADMIN_URL=https://admin.yourdomain.com

# Service URLs (Internal, if using Kubernetes/Docker Swarm)
USER_SERVICE_URL=http://user-service:8080
WALLET_SERVICE_URL=http://wallet-service:8081
ORDER_SERVICE_URL=http://order-service:8082
MARKET_SERVICE_URL=http://market-service:8083
CREDIT_EXCHANGE_URL=http://credit-exchange-service:8084
NOTIFICATION_SERVICE_URL=http://notification-service:8085
SPORTSBOOK_SERVICE_URL=http://sportsbook-service:8086
MATCHING_ENGINE_URL=http://matching-engine:9000

# Monitoring (optional)
SENTRY_DSN=<your-sentry-dsn>
LOG_LEVEL=info

# CORS (restrict in production)
ALLOWED_ORIGINS=https://yourdomain.com,https://admin.yourdomain.com
```

### Deployment Options

#### Option 1: Docker Compose (Simple)

```bash
# On production server
git clone <repository>
cd LFG

# Set up production environment
cp .env.example .env
nano .env  # Configure with production values

# Deploy
docker-compose -f docker-compose.prod.yml up -d

# Monitor
docker-compose logs -f
```

#### Option 2: Kubernetes (Scalable)

```bash
# Create namespace
kubectl create namespace lfg-prod

# Create secrets
kubectl create secret generic lfg-secrets \
  --from-literal=jwt-secret=$JWT_SECRET \
  --from-literal=encryption-key=$ENCRYPTION_KEY \
  --from-literal=db-password=$DB_PASSWORD \
  -n lfg-prod

# Apply configurations
kubectl apply -f k8s/postgres.yaml -n lfg-prod
kubectl apply -f k8s/services.yaml -n lfg-prod
kubectl apply -f k8s/ingress.yaml -n lfg-prod

# Check status
kubectl get pods -n lfg-prod
kubectl get services -n lfg-prod
```

#### Option 3: Cloud Deployment

**AWS:**
- Use ECS/EKS for container orchestration
- RDS PostgreSQL for database
- ALB for load balancing
- CloudFront for CDN
- S3 for static assets
- CloudWatch for monitoring

**GCP:**
- Use GKE for Kubernetes
- Cloud SQL for PostgreSQL
- Cloud Load Balancing
- Cloud CDN
- Cloud Storage
- Cloud Monitoring

**Azure:**
- Use AKS for Kubernetes
- Azure Database for PostgreSQL
- Azure Load Balancer
- Azure CDN
- Blob Storage
- Azure Monitor

---

## ⚙️ Environment Configuration

### Complete Environment Variables

```bash
# ===== SECURITY =====
JWT_SECRET=your-super-secret-jwt-key-at-least-64-characters-long
ENCRYPTION_KEY=your-32-byte-aes-encryption-key-base64-encoded
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001

# ===== DATABASE =====
DATABASE_URL=postgresql://lfguser:lfgpassword@postgres:5432/lfg?sslmode=disable
DB_HOST=postgres
DB_PORT=5432
DB_NAME=lfg
DB_USER=lfguser
DB_PASSWORD=lfgpassword
DB_SSLMODE=disable  # Use 'require' in production

# PostgreSQL settings (for postgres container)
POSTGRES_DB=lfg
POSTGRES_USER=lfguser
POSTGRES_PASSWORD=lfgpassword

# Connection pool settings
DB_MAX_CONNS=25
DB_MIN_CONNS=5
DB_MAX_CONN_LIFETIME=3600
DB_MAX_CONN_IDLE_TIME=1800

# ===== SERVICE URLS (Internal) =====
USER_SERVICE_URL=http://user-service:8080
WALLET_SERVICE_URL=http://wallet-service:8081
ORDER_SERVICE_URL=http://order-service:8082
MARKET_SERVICE_URL=http://market-service:8083
CREDIT_EXCHANGE_URL=http://credit-exchange-service:8084
NOTIFICATION_SERVICE_URL=http://notification-service:8085
SPORTSBOOK_SERVICE_URL=http://sportsbook-service:8086
MATCHING_ENGINE_URL=http://matching-engine:9000

# ===== PUBLIC URLS =====
API_GATEWAY_URL=http://localhost:8000
FRONTEND_URL=http://localhost:3000
ADMIN_URL=http://localhost:3001

# ===== REACT APPS =====
# Frontend
REACT_APP_API_URL=http://localhost:8000
REACT_APP_WS_URL=ws://localhost:8000

# Admin Panel
REACT_APP_API_BASE_URL=http://localhost:8000

# ===== LOGGING =====
LOG_LEVEL=info  # debug, info, warn, error
LOG_FORMAT=json  # json or text

# ===== RATE LIMITING =====
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60  # seconds

# ===== MONITORING (Optional) =====
SENTRY_DSN=
DATADOG_API_KEY=
NEW_RELIC_LICENSE_KEY=

# ===== FEATURE FLAGS =====
ENABLE_REGISTRATION=true
ENABLE_SPORTSBOOK=true
ENABLE_PREDICTION_MARKETS=true
MAINTENANCE_MODE=false
```

### Generating Secure Secrets

```bash
# JWT Secret (64 characters minimum)
openssl rand -base64 64

# Encryption Key (exactly 32 bytes for AES-256)
openssl rand -base64 32

# Database Password
openssl rand -base64 32 | tr -d "=+/" | cut -c1-32

# UUID for secrets
uuidgen
```

---

## 🗄️ Database Setup

### Initialize Database

Database is automatically initialized on first run using:
1. `/database/init.sql` - Schema creation
2. `/database/seed.sql` - Seed data

### Manual Database Access

```bash
# Access PostgreSQL container
docker exec -it lfg-postgres psql -U lfguser -d lfg

# Common queries
\dt  # List all tables
\d users  # Describe users table
SELECT * FROM sportsbook_providers;
SELECT * FROM markets;

# Exit
\q
```

### Database Backup

```bash
# Backup database
docker exec lfg-postgres pg_dump -U lfguser lfg > backup_$(date +%Y%m%d_%H%M%S).sql

# Backup with compression
docker exec lfg-postgres pg_dump -U lfguser lfg | gzip > backup_$(date +%Y%m%d_%H%M%S).sql.gz

# Automated daily backups (add to crontab)
0 2 * * * docker exec lfg-postgres pg_dump -U lfguser lfg | gzip > /backups/lfg_$(date +\%Y\%m\%d).sql.gz
```

### Database Restore

```bash
# Restore from backup
docker exec -i lfg-postgres psql -U lfguser -d lfg < backup.sql

# Restore from compressed backup
gunzip < backup.sql.gz | docker exec -i lfg-postgres psql -U lfguser -d lfg
```

### Database Migrations

```bash
# Create new migration
# migrations/YYYYMMDD_HHMMSS_description.sql

# Apply migration manually
docker exec -i lfg-postgres psql -U lfguser -d lfg < migrations/20250107_120000_add_new_column.sql
```

---

## 🏥 Service Health Checks

### Check All Services

```bash
# Docker health status
docker-compose ps

# API Gateway health
curl http://localhost:8000/health

# Individual service health checks
curl http://localhost:8080/health  # User Service
curl http://localhost:8081/health  # Wallet Service
curl http://localhost:8083/health  # Market Service
curl http://localhost:8082/health  # Order Service
curl http://localhost:8084/health  # Matching Engine
curl http://localhost:8086/health  # Credit Exchange
curl http://localhost:8087/health  # Notification Service
curl http://localhost:8088/health  # Sportsbook Service

# Database health
docker exec lfg-postgres pg_isready -U lfguser

# Frontend health
curl http://localhost:3000
curl http://localhost:3001  # Admin Panel
```

### Automated Health Check Script

```bash
#!/bin/bash
# health-check.sh

services=(
  "http://localhost:8000"  # API Gateway
  "http://localhost:8080"  # User
  "http://localhost:8081"  # Wallet
  "http://localhost:8083"  # Market
  "http://localhost:8082"  # Order
  "http://localhost:8084"  # Matching
  "http://localhost:8086"  # Credit Exchange
  "http://localhost:8087"  # Notification
  "http://localhost:8088"  # Sportsbook
)

for service in "${services[@]}"; do
  if curl -f -s "${service}/health" > /dev/null; then
    echo "✓ ${service} - healthy"
  else
    echo "✗ ${service} - unhealthy"
  fi
done
```

---

## 📊 Monitoring & Logging

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f sportsbook-service

# Last 100 lines
docker-compose logs --tail=100 api-gateway

# Since timestamp
docker-compose logs --since 2025-01-07T10:00:00

# Export logs
docker-compose logs > full-logs-$(date +%Y%m%d).txt
```

### Log Aggregation

For production, use:
- **ELK Stack** (Elasticsearch, Logstash, Kibana)
- **Loki + Grafana**
- **Datadog**
- **Splunk**
- **CloudWatch** (AWS)

### Monitoring Stack (Optional)

```yaml
# Add to docker-compose.yml
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3002:3000"
    depends_on:
      - prometheus
```

---

## 🔧 Troubleshooting

### Common Issues

**Issue: Services fail to start**
```bash
# Check logs
docker-compose logs

# Check disk space
df -h

# Check memory
free -h

# Recreate containers
docker-compose down
docker-compose up -d --force-recreate
```

**Issue: Database connection errors**
```bash
# Check database is running
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Test connection
docker exec -it lfg-postgres psql -U lfguser -d lfg

# Restart database
docker-compose restart postgres
```

**Issue: Port already in use**
```bash
# Find process using port 8000
sudo lsof -i :8000

# Kill process
sudo kill -9 <PID>

# Or change port in docker-compose.yml
```

**Issue: Frontend can't connect to API**
```bash
# Check CORS settings
# Verify REACT_APP_API_URL in .env
# Check browser console for errors
# Verify API Gateway is accessible

curl http://localhost:8000/health
```

### Reset Everything

```bash
# Nuclear option: delete everything and start fresh
docker-compose down -v --remove-orphans
docker system prune -a --volumes -f

# Rebuild and restart
docker-compose up -d --build
```

---

## 🔒 Security Hardening

### Production Security Checklist

- [ ] Change all default credentials
- [ ] Use strong, unique JWT_SECRET (64+ chars)
- [ ] Use strong ENCRYPTION_KEY (32 bytes)
- [ ] Enable SSL/TLS (HTTPS)
- [ ] Configure firewall (ufw/iptables)
- [ ] Restrict database access to internal network
- [ ] Enable PostgreSQL SSL connections
- [ ] Implement rate limiting on all APIs
- [ ] Set up CORS with specific origins (no *)
- [ ] Enable security headers (HSTS, CSP, X-Frame-Options)
- [ ] Disable directory listing on web servers
- [ ] Implement input validation on all endpoints
- [ ] Use prepared statements (prevent SQL injection)
- [ ] Encrypt sensitive data at rest
- [ ] Set up intrusion detection (fail2ban)
- [ ] Regular security updates
- [ ] Implement DDOS protection
- [ ] Set up Web Application Firewall (WAF)
- [ ] Enable audit logging
- [ ] Regular penetration testing

### SSL/TLS Setup

```bash
# Install certbot
sudo apt-get install certbot

# Get Let's Encrypt certificate
sudo certbot certonly --standalone -d api.yourdomain.com
sudo certbot certonly --standalone -d yourdomain.com
sudo certbot certonly --standalone -d admin.yourdomain.com

# Auto-renewal
sudo crontab -e
0 0 1 * * certbot renew --quiet
```

### Firewall Configuration

```bash
# Allow SSH
sudo ufw allow 22

# Allow HTTP/HTTPS
sudo ufw allow 80
sudo ufw allow 443

# Block all other external access to service ports
sudo ufw deny 8000:9000/tcp

# Enable firewall
sudo ufw enable
```

---

## 📞 Support & Maintenance

### Update Services

```bash
# Pull latest images
docker-compose pull

# Rebuild services
docker-compose up -d --build

# Rolling update (zero downtime)
docker-compose up -d --no-deps --build service-name
```

### Performance Tuning

```bash
# PostgreSQL tuning
# Edit /var/lib/postgresql/data/postgresql.conf

shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 4MB
min_wal_size = 1GB
max_wal_size = 4GB
```

### Scaling

```bash
# Horizontal scaling with Docker Swarm
docker swarm init
docker stack deploy -c docker-compose.yml lfg

# Scale specific services
docker service scale lfg_user-service=3
docker service scale lfg_market-service=3
```

---

## 📝 Maintenance Tasks

### Daily
- Check service health
- Monitor error logs
- Check disk space
- Verify backups completed

### Weekly
- Review security logs
- Update dependencies
- Performance metrics review
- Database optimization (VACUUM, ANALYZE)

### Monthly
- Security patches
- Load testing
- Disaster recovery drill
- Capacity planning

---

**For additional support, consult the README.md and LEGAL_DISCLAIMER.md files.**

**Last Updated: January 2025**
