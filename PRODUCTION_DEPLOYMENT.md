# Production Deployment Guide

This guide covers the complete production deployment process for the LFG platform.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Infrastructure Setup](#infrastructure-setup)
3. [Database Setup](#database-setup)
4. [Application Deployment](#application-deployment)
5. [Monitoring Setup](#monitoring-setup)
6. [Security Hardening](#security-hardening)
7. [Maintenance](#maintenance)

## Prerequisites

### Required Tools

- Docker & Docker Compose (v2.20+)
- Kubernetes cluster (v1.27+) or kubectl
- PostgreSQL 15+
- Redis 7+
- golang-migrate
- make

### Recommended Infrastructure

- **Compute**:
  - Kubernetes cluster with at least 3 nodes
  - 4 CPU cores, 8GB RAM minimum per node

- **Database**:
  - Managed PostgreSQL instance with automated backups
  - 4 CPU cores, 16GB RAM, 100GB SSD minimum

- **Redis**:
  - Managed Redis instance
  - 2GB RAM minimum

## Infrastructure Setup

### 1. Environment Variables

Create production environment file:

```bash
cp .env.example .env
```

Edit `.env` with production values:

```env
# Database
POSTGRES_HOST=your-postgres-host
POSTGRES_PORT=5432
POSTGRES_DB=lfg_production
POSTGRES_USER=lfg_user
POSTGRES_PASSWORD=your-secure-password

# Redis
REDIS_HOST=your-redis-host
REDIS_PORT=6379
REDIS_PASSWORD=your-redis-password

# Security
JWT_SECRET=your-jwt-secret-key-at-least-32-chars
INTERNAL_API_KEY=your-internal-api-key
ADMIN_API_KEY=your-admin-api-key

# External APIs
COINMARKETCAP_API_KEY=your-coinmarketcap-key

# Application
ENVIRONMENT=production
CORS_ALLOWED_ORIGINS=https://yourdomain.com,https://www.yourdomain.com
```

### 2. SSL/TLS Certificates

Obtain SSL certificates for your domain:

```bash
# Using Let's Encrypt with Certbot
certbot certonly --standalone -d yourdomain.com -d www.yourdomain.com
```

### 3. DNS Configuration

Configure DNS records:

```
A     yourdomain.com          -> YOUR_SERVER_IP
A     www.yourdomain.com      -> YOUR_SERVER_IP
A     api.yourdomain.com      -> YOUR_SERVER_IP
CNAME admin.yourdomain.com    -> yourdomain.com
```

## Database Setup

### 1. Initialize Database

```bash
# Create database and user
psql -U postgres -h your-postgres-host << EOF
CREATE DATABASE lfg_production;
CREATE USER lfg_user WITH PASSWORD 'your-secure-password';
GRANT ALL PRIVILEGES ON DATABASE lfg_production TO lfg_user;
EOF
```

### 2. Run Migrations

```bash
# Set environment variables
export POSTGRES_HOST=your-postgres-host
export POSTGRES_PORT=5432
export POSTGRES_USER=lfg_user
export POSTGRES_DB=lfg_production
export POSTGRES_PASSWORD=your-secure-password

# Run migrations
make migrate-up
```

### 3. Verify Database

```bash
psql -U lfg_user -h your-postgres-host -d lfg_production -c "\dt"
```

## Application Deployment

### Option 1: Docker Compose (Simple Deployment)

#### Deploy with Docker Compose

```bash
# Build images
make docker-build

# Start production stack
make docker-prod

# Verify services
docker-compose -f docker-compose.prod.yml ps
```

#### Access Services

- Frontend: http://localhost:3000
- API Gateway: http://localhost:8000
- Admin Panel: http://localhost:3001

### Option 2: Kubernetes (Production-Grade)

#### 1. Create Kubernetes Secrets

```bash
# Create namespace
kubectl apply -f k8s/namespace.yaml

# Create secrets from environment variables
kubectl create secret generic lfg-secrets \
  --from-literal=database-url="postgres://lfg_user:password@postgres:5432/lfg_production" \
  --from-literal=jwt-secret="your-jwt-secret" \
  --from-literal=internal-api-key="your-internal-key" \
  --from-literal=admin-api-key="your-admin-key" \
  --from-literal=coinmarketcap-api-key="your-cmc-key" \
  -n lfg
```

#### 2. Create ConfigMap

```bash
kubectl apply -f k8s/configmap.yaml
```

#### 3. Deploy Infrastructure

```bash
# Deploy PostgreSQL
kubectl apply -f k8s/postgres-statefulset.yaml

# Deploy Redis
kubectl apply -f k8s/redis-deployment.yaml

# Wait for databases to be ready
kubectl wait --for=condition=ready pod -l app=postgres -n lfg --timeout=300s
kubectl wait --for=condition=ready pod -l app=redis -n lfg --timeout=300s
```

#### 4. Deploy Application Services

```bash
# Deploy all services
kubectl apply -f k8s/all-services.yaml

# Wait for deployments
kubectl rollout status deployment -n lfg --timeout=10m
```

#### 5. Deploy Monitoring Stack

```bash
kubectl apply -f k8s/monitoring.yaml

# Check monitoring services
kubectl get all -n monitoring
```

#### 6. Verify Deployment

```bash
# Check all pods
kubectl get pods -n lfg

# Check services
kubectl get svc -n lfg

# View logs
kubectl logs -f deployment/api-gateway -n lfg
```

#### 7. Access Services

```bash
# Get LoadBalancer IP
kubectl get svc api-gateway -n lfg

# Port forward for testing
kubectl port-forward svc/api-gateway 8000:8000 -n lfg
```

## Monitoring Setup

### 1. Access Grafana

```bash
# Get Grafana LoadBalancer IP
kubectl get svc grafana -n monitoring

# Or port forward
kubectl port-forward svc/grafana 3000:3000 -n monitoring
```

Default credentials:
- Username: admin
- Password: (check secret or use default: changeme123)

### 2. Configure Dashboards

1. Login to Grafana
2. Add Prometheus data source: http://prometheus.monitoring:9090
3. Import dashboards from grafana.com:
   - Go Application Dashboard (ID: 14061)
   - PostgreSQL Dashboard (ID: 9628)
   - Redis Dashboard (ID: 11835)

### 3. Set Up Alerts

Configure alerts in Prometheus for:
- High error rates
- Service downtime
- Database connection issues
- High memory/CPU usage

## Security Hardening

### 1. Network Security

```bash
# Configure firewall rules (example for UFW)
ufw allow 80/tcp
ufw allow 443/tcp
ufw deny 5432/tcp  # Database should not be publicly accessible
ufw deny 6379/tcp  # Redis should not be publicly accessible
ufw enable
```

### 2. Database Security

```sql
-- Revoke public access
REVOKE ALL ON DATABASE lfg_production FROM PUBLIC;

-- Set up read-only user for reporting
CREATE USER lfg_readonly WITH PASSWORD 'secure-password';
GRANT CONNECT ON DATABASE lfg_production TO lfg_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO lfg_readonly;
```

### 3. API Rate Limiting

Rate limiting is configured in the API Gateway:
- 10 requests per second per IP
- Configurable via Redis-based rate limiter

### 4. SSL/TLS Configuration

Configure nginx or load balancer with SSL:

```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    location / {
        proxy_pass http://localhost:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Maintenance

### Database Backups

#### Automated Backups

Backups run automatically via the backup service in docker-compose.prod.yml:
- Daily backups at midnight
- 7 days retention for daily backups
- 4 weeks retention for weekly backups
- 6 months retention for monthly backups

#### Manual Backup

```bash
make backup
```

#### Restore from Backup

```bash
make restore FILE=./backups/lfg_backup_20231108_120000.sql.gz
```

### Database Migrations

#### Create New Migration

```bash
make migrate-create NAME=add_new_feature
```

#### Apply Migrations

```bash
make migrate-up
```

#### Rollback Migration

```bash
make migrate-down
```

### Log Management

#### View Logs

```bash
# Docker Compose
docker-compose -f docker-compose.prod.yml logs -f service-name

# Kubernetes
kubectl logs -f deployment/service-name -n lfg
```

#### Log Rotation

Logs are automatically rotated with max size 10MB and max 3 files.

### Health Checks

All services expose health endpoints:

```bash
# Liveness check
curl http://localhost:8080/health/live

# Readiness check
curl http://localhost:8080/health/ready
```

### Scaling

#### Docker Compose

```bash
docker-compose -f docker-compose.prod.yml up -d --scale user-service=3
```

#### Kubernetes

```bash
kubectl scale deployment user-service --replicas=5 -n lfg
```

### Updates and Rollbacks

#### Update Service

```bash
# Build new image
make docker-build

# Update Kubernetes deployment
kubectl set image deployment/user-service user-service=lfg/user-service:v2.0 -n lfg
```

#### Rollback

```bash
kubectl rollout undo deployment/user-service -n lfg
```

## Troubleshooting

### Service Not Starting

```bash
# Check logs
kubectl logs deployment/service-name -n lfg

# Check events
kubectl get events -n lfg --sort-by='.lastTimestamp'

# Check pod status
kubectl describe pod pod-name -n lfg
```

### Database Connection Issues

```bash
# Test database connection
psql -U lfg_user -h postgres-host -d lfg_production

# Check database logs
kubectl logs statefulset/postgres -n lfg
```

### High Memory/CPU Usage

```bash
# Check resource usage
kubectl top nodes
kubectl top pods -n lfg

# Scale up if needed
kubectl scale deployment/service-name --replicas=5 -n lfg
```

## Production Checklist

Before going live:

- [ ] All environment variables configured
- [ ] SSL certificates installed and valid
- [ ] Database migrations applied
- [ ] Backups configured and tested
- [ ] Monitoring and alerting set up
- [ ] Security hardening completed
- [ ] Load testing performed
- [ ] Disaster recovery plan documented
- [ ] Team trained on operations
- [ ] Documentation updated

## Support

For production support:
- Email: support@lfg-platform.com
- Slack: #lfg-production
- On-call: +1-XXX-XXX-XXXX
