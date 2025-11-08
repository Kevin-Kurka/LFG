# LFG Platform - Quick Start Production Guide

## Prerequisites Checklist

- [ ] Docker 24.0+ installed
- [ ] Docker Compose 2.20+ installed
- [ ] kubectl installed (for Kubernetes)
- [ ] 8GB RAM minimum
- [ ] 50GB disk space
- [ ] Production secrets generated

## 5-Minute Quick Start (Docker Compose)

```bash
# 1. Clone and navigate
cd /home/user/LFG

# 2. Copy and edit environment
cp .env.example .env
nano .env  # Add your production secrets

# 3. Start production stack
make docker-prod

# 4. Verify all services are running
docker-compose -f docker-compose.prod.yml ps

# 5. Access the platform
# Frontend:    http://localhost:3000
# Admin Panel: http://localhost:3001
# API Gateway: http://localhost:8000
# Grafana:     Configure via monitoring service
```

## Production Secrets Generation

```bash
# Generate JWT Secret (64 chars)
openssl rand -base64 64

# Generate Encryption Key (32 chars)
openssl rand -base64 32

# Generate Internal API Key
openssl rand -hex 32

# Generate Admin API Key
openssl rand -hex 32

# Add these to .env file
```

## Kubernetes Quick Deploy

```bash
# 1. Create namespace
kubectl create namespace lfg

# 2. Create secrets
kubectl create secret generic lfg-secrets \
  --from-literal=database-url="postgres://user:pass@host:5432/db" \
  --from-literal=jwt-secret="your-jwt-secret" \
  --from-literal=internal-api-key="your-internal-key" \
  --from-literal=admin-api-key="your-admin-key" \
  --from-literal=coinmarketcap-api-key="your-cmc-key" \
  -n lfg

# 3. Deploy all services
make k8s-deploy

# 4. Check status
make k8s-status

# 5. Get API Gateway URL
kubectl get svc api-gateway -n lfg
```

## Common Make Commands

```bash
make help              # Show all commands
make docker-up         # Start development
make docker-prod       # Start production
make k8s-deploy        # Deploy to Kubernetes
make migrate-up        # Run database migrations
make backup            # Backup database
make test              # Run tests
make docker-logs       # View logs
```

## Health Check URLs

After deployment, verify services:

```bash
# API Gateway
curl http://localhost:8000/health

# User Service
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready

# Wallet Service
curl http://localhost:8081/health/live

# Market Service  
curl http://localhost:8082/health/live

# All services have /health/live and /health/ready endpoints
```

## Database Setup

```bash
# Run migrations
make migrate-up

# Create new migration
make migrate-create NAME=add_new_table

# Rollback migration
make migrate-down

# Manual migration
./database/migrate.sh up
```

## Backup & Restore

```bash
# Backup database
make backup

# Backups are stored in ./backups/

# Restore from backup
make restore FILE=./backups/lfg_backup_20231108_120000.sql.gz
```

## Monitoring Access

### Prometheus
```bash
kubectl port-forward svc/prometheus 9090:9090 -n monitoring
# Access: http://localhost:9090
```

### Grafana
```bash
kubectl port-forward svc/grafana 3000:3000 -n monitoring
# Access: http://localhost:3000
# Default: admin / changeme123
```

## Troubleshooting

### Service won't start
```bash
# Check logs
docker-compose -f docker-compose.prod.yml logs service-name

# Or for Kubernetes
kubectl logs deployment/service-name -n lfg

# Check health
curl http://localhost:PORT/health/live
```

### Database connection error
```bash
# Verify database is running
docker-compose -f docker-compose.prod.yml ps postgres

# Check database logs
docker-compose -f docker-compose.prod.yml logs postgres

# Test connection
psql -U lfg_user -h localhost -d lfg_production
```

### Can't access frontend
```bash
# Verify containers are running
docker-compose -f docker-compose.prod.yml ps

# Check API Gateway
curl http://localhost:8000/health

# Check frontend logs
docker-compose -f docker-compose.prod.yml logs frontend
```

## Production Checklist

Before going live:

- [ ] All environment variables configured
- [ ] Production secrets generated and stored securely
- [ ] SSL/TLS certificates obtained and configured
- [ ] DNS records configured
- [ ] Firewall rules configured
- [ ] Database migrations applied
- [ ] Backups tested and verified
- [ ] Monitoring configured and alerts set
- [ ] Health checks verified
- [ ] Load testing completed
- [ ] Security audit completed
- [ ] Legal compliance verified (see LEGAL_DISCLAIMER.md)
- [ ] Documentation reviewed
- [ ] Team trained on operations

## Next Steps

1. Review [PRODUCTION_DEPLOYMENT.md](./PRODUCTION_DEPLOYMENT.md) for complete guide
2. Review [README.md](./README.md) for platform overview
3. Review [LEGAL_DISCLAIMER.md](./LEGAL_DISCLAIMER.md) - **REQUIRED**
4. Review [api-docs.yaml](./api-docs.yaml) for API reference

## Support

For issues or questions:
- Check PRODUCTION_DEPLOYMENT.md troubleshooting section
- Review service logs
- Check health endpoints
- Review Kubernetes events: `kubectl get events -n lfg`

## Important Links

- Main README: [README.md](./README.md)
- Production Guide: [PRODUCTION_DEPLOYMENT.md](./PRODUCTION_DEPLOYMENT.md)
- Implementation Summary: [IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md)
- API Documentation: [api-docs.yaml](./api-docs.yaml)
- Legal Disclaimer: [LEGAL_DISCLAIMER.md](./LEGAL_DISCLAIMER.md)
- Makefile Commands: Run `make help`

---

**Status**: Production Ready
**Last Updated**: November 8, 2025
