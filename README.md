# LFG - Prediction Marketplace & Sportsbook Platform

> **⚠️ CRITICAL LEGAL WARNING**: This software is for demonstration and educational purposes only. Operating prediction markets and sportsbook integration services requires extensive legal licenses and regulatory compliance. **DO NOT use in production without proper legal counsel.** See [LEGAL_DISCLAIMER.md](./LEGAL_DISCLAIMER.md) for details.

---

## 🎯 Overview

**LFG** is a full-featured prediction marketplace and sportsbook aggregation platform that combines:

1. **Kalshi-Style Prediction Markets** - Binary outcome markets with order book trading
2. **Sportsbook Integration** - Multi-account management and odds comparison across major sportsbooks
3. **Arbitrage Detection** - Automatic identification of guaranteed profit opportunities
4. **Hedge Calculator** - Optimal hedge calculations for existing positions

### Key Features

- ✅ **Prediction Markets**: Create and trade on YES/NO outcome markets
- ✅ **Order Book Trading**: Limit and market orders with price-time priority matching
- ✅ **Multi-Sportsbook**: FanDuel, DraftKings, BetMGM, Caesars, ESPN BET, Fanatics
- ✅ **Odds Comparison**: Side-by-side odds from all major sportsbooks
- ✅ **Arbitrage Detection**: Real-time profit opportunity identification
- ✅ **Hedge Calculator**: Calculate optimal hedge stakes for guaranteed profit
- ✅ **Bet Tracking**: Track bets across all sportsbooks with P&L
- ✅ **Real-Time Updates**: WebSocket notifications for live odds and trades
- ✅ **Wallet System**: Credit-based internal currency with transactions
- ✅ **Admin Panel**: Full market and sportsbook management interface

---

## 🏗️ Architecture

### Microservices Backend (Go)

```
├── API Gateway (8000)         - Reverse proxy, auth, rate limiting
├── User Service (8080)        - Authentication, user management
├── Wallet Service (8081)      - Balance, transactions
├── Order Service (8082)       - Order placement, cancellation
├── Market Service (8083)      - Market CRUD, order books
├── Matching Engine (8084)     - Order matching, trade execution
├── Credit Exchange (8086)     - Buy/sell credits
├── Notification Service (8087) - WebSocket real-time updates
└── Sportsbook Service (8088)  - Odds aggregation, arbitrage, hedging
```

### Frontend (React + TypeScript)

```
├── Web App (3000)             - User-facing prediction market & sportsbook UI
└── Admin Panel (3001)         - Market management, analytics, moderation
```

### Database

- **PostgreSQL 15** with 20+ tables
- Prediction markets, orders, trades
- Sports events, odds, sportsbook accounts
- Arbitrage and hedge opportunities
- User bets and wallet transactions

---

## 🚀 Quick Start

### Prerequisites

- **Docker** >= 24.0
- **Docker Compose** >= 2.20
- **Make** (optional, for convenience commands)
- **kubectl** (optional, for Kubernetes deployment)
- 8GB RAM minimum
- 50GB disk space

### Development Installation

```bash
# 1. Clone repository
git clone <repository-url>
cd LFG

# 2. Configure environment
cp .env.example .env

# 3. Generate secure secrets
# JWT Secret
openssl rand -base64 64

# Encryption Key
openssl rand -base64 32

# 4. Update .env with generated secrets
nano .env

# 5. Start all services (using Make)
make docker-up

# OR manually
docker-compose up -d

# 6. Wait 30 seconds for services to initialize

# 7. Access the platform
# Frontend:    http://localhost:3000
# Admin Panel: http://localhost:3001
# API Gateway: http://localhost:8000
```

### Production Installation

```bash
# 1. Set up production environment
make dev-setup

# 2. Run database migrations
make migrate-up

# 3. Start production stack
make docker-prod

# For complete production deployment guide, see:
# - PRODUCTION_DEPLOYMENT.md
```

### Verify Installation

```bash
# Check all containers are running
docker-compose ps

# Test API Gateway
curl http://localhost:8000/health

# View logs
docker-compose logs -f
```

---

## 📚 Documentation

- **[PRODUCTION_DEPLOYMENT.md](./PRODUCTION_DEPLOYMENT.md)** - Production deployment guide with Kubernetes
- **[DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md)** - General deployment instructions
- **[LEGAL_DISCLAIMER.md](./LEGAL_DISCLAIMER.md)** - **REQUIRED READING** - Legal warnings and compliance requirements
- **[api-docs.yaml](./api-docs.yaml)** - OpenAPI/Swagger API specification
- **[Makefile](./Makefile)** - Available make commands for development and deployment
- **[backend/README.md](./backend/README.md)** - Backend architecture and API docs
- **[frontend-web/README.md](./frontend-web/README.md)** - Frontend setup and features
- **[admin-panel/README.md](./admin-panel/README.md)** - Admin panel documentation

---

## 🎮 Usage

### User Workflow

1. **Register Account**
   - Navigate to http://localhost:3000
   - Click "Register" and create account
   - Login with credentials

2. **Prediction Markets**
   - Browse markets at `/markets`
   - View market details and order book
   - Place limit or market orders
   - Monitor positions in dashboard

3. **Sportsbook Integration**
   - Link sportsbook accounts at `/link-account`
   - View odds comparison at `/sportsbook`
   - Check arbitrage opportunities at `/arbitrage`
   - Calculate hedges at `/hedges`
   - Track bets at `/bets`

### Admin Workflow

1. **Access Admin Panel**
   - Navigate to http://localhost:3001
   - Login with admin credentials

2. **Create Markets**
   - Go to "Markets" → "Create Market"
   - Define question, outcomes, expiry
   - Publish market

3. **Resolve Markets**
   - View market details
   - Click "Resolve Market"
   - Select winning outcome
   - Confirm resolution

4. **Manage Sportsbooks**
   - View connected sportsbooks
   - Monitor odds data
   - View arbitrage opportunities

---

## 🔧 Development

### Backend Development

```bash
# Navigate to service directory
cd backend/user-service

# Install dependencies
go mod tidy

# Run service locally
DATABASE_URL="postgresql://lfguser:lfgpassword@localhost:5432/lfg" \
JWT_SECRET="dev-secret" \
go run main.go

# Run tests
go test ./...

# Build
go build -o user-service
```

### Frontend Development

```bash
# Navigate to frontend
cd frontend-web

# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Run tests
npm test
```

### Admin Panel Development

```bash
# Navigate to admin panel
cd admin-panel

# Install dependencies
npm install

# Start development server
npm start

# Build for production
npm run build
```

---

## 📊 Technology Stack

### Backend
- **Language**: Go 1.24+
- **Database**: PostgreSQL 15 with pgx driver
- **Auth**: JWT (golang-jwt/jwt)
- **WebSocket**: gorilla/websocket
- **Password**: bcrypt
- **Encryption**: AES-256-GCM for credentials

### Frontend
- **Framework**: React 18 + TypeScript
- **Build Tool**: Vite
- **Styling**: Tailwind CSS
- **Routing**: React Router v6
- **State**: React Query
- **HTTP Client**: Axios
- **WebSocket**: Native WebSocket API

### Infrastructure
- **Containers**: Docker + Docker Compose
- **Orchestration**: Kubernetes with full manifests
- **Cache**: Redis 7+ for rate limiting and caching
- **Monitoring**: Prometheus + Grafana
- **Database Migrations**: golang-migrate
- **CI/CD**: GitHub Actions workflows
- **Reverse Proxy**: nginx (in production)

---

## 🔒 Security

### Implemented Security Measures

- ✅ JWT-based authentication with secure tokens
- ✅ bcrypt password hashing (cost 12)
- ✅ AES-256-GCM encryption for sportsbook credentials
- ✅ Input validation on all endpoints
- ✅ SQL injection prevention (prepared statements)
- ✅ CORS configuration
- ✅ Redis-based rate limiting (10 req/sec per IP)
- ✅ Request ID tracking for audit trails
- ✅ Structured logging with zap
- ✅ Health checks (liveness and readiness)
- ✅ Graceful shutdown for all services
- ✅ Automated database backups
- ✅ Secure headers (HSTS, CSP, X-Frame-Options)

### Required for Production

- [ ] SSL/TLS certificates (Let's Encrypt)
- [ ] Web Application Firewall (WAF)
- [ ] DDoS protection
- [ ] Regular security audits
- [ ] Penetration testing
- [ ] Intrusion detection (fail2ban)
- [ ] Security monitoring (SIEM)

See [LEGAL_DISCLAIMER.md](./LEGAL_DISCLAIMER.md) for complete security requirements.

---

## 📈 API Documentation

### Authentication Endpoints

```bash
# Register
POST /register
{
  "email": "user@example.com",
  "password": "securepassword"
}

# Login
POST /login
{
  "email": "user@example.com",
  "password": "securepassword"
}
Response: {"token": "jwt-token"}

# Get Profile (requires JWT)
GET /profile
Headers: Authorization: Bearer <token>
```

### Prediction Market Endpoints

```bash
# List Markets
GET /markets?status=OPEN&search=chiefs

# Get Market Details
GET /markets/:id

# Create Market (admin)
POST /markets
{
  "ticker": "CHIEFS-SB-2025",
  "question": "Will Chiefs win Super Bowl LIX?",
  "rules": "Resolves YES if Chiefs win...",
  "resolution_source": "Official NFL records",
  "expires_at": "2025-02-09T23:59:59Z"
}

# Place Order
POST /orders/place
{
  "contract_id": "uuid",
  "type": "LIMIT",
  "side": "BUY",
  "quantity": 100,
  "limit_price": 0.65
}

# Get Order Book
GET /markets/:id/orderbook
```

### Sportsbook Endpoints

```bash
# Link Sportsbook Account
POST /sportsbook/accounts
{
  "sportsbook_id": "uuid",
  "username": "my_username",
  "password": "my_password"
}

# Get Events with Odds
GET /sportsbook/events?sport=nfl&date=2025-01-10

# Get Arbitrage Opportunities
GET /sportsbook/arbitrage

# Get Hedge Opportunities
GET /sportsbook/hedges

# Track Bet
POST /sportsbook/bets
{
  "sportsbook_account_id": "uuid",
  "event_id": "uuid",
  "bet_type": "MONEYLINE",
  "selection": "Kansas City Chiefs",
  "stake": 100.00,
  "odds_american": -140
}
```

For complete API documentation, see [backend/README.md](./backend/README.md).

---

## 🧪 Testing

### Automated Tests

```bash
# Backend tests
cd backend/user-service
go test ./...

# Frontend tests
cd frontend-web
npm test

# E2E tests (if implemented)
npm run test:e2e
```

### Manual Testing Checklist

- [ ] User registration and login
- [ ] Market creation and listing
- [ ] Order placement (limit and market)
- [ ] Order matching and trade execution
- [ ] Wallet balance updates
- [ ] Sportsbook account linking
- [ ] Odds display and comparison
- [ ] Arbitrage detection
- [ ] Hedge calculator
- [ ] Bet tracking
- [ ] WebSocket real-time updates
- [ ] Admin market management
- [ ] Market resolution

---

## 📦 Deployment

### Local Development
```bash
make docker-up
# OR
docker-compose up -d
```

### Production Deployment

#### Using Docker Compose
```bash
make docker-prod
```

#### Using Kubernetes
```bash
# Deploy all services
make k8s-deploy

# Validate manifests
make k8s-validate

# Check status
make k8s-status
```

### Database Management
```bash
# Run migrations
make migrate-up

# Create new migration
make migrate-create NAME=add_new_table

# Backup database
make backup

# Restore database
make restore FILE=./backups/backup.sql.gz
```

### Common Commands
```bash
# See all available commands
make help

# Run tests
make test

# Build services
make build

# View logs
make docker-logs
```

See [PRODUCTION_DEPLOYMENT.md](./PRODUCTION_DEPLOYMENT.md) for complete production deployment guide including:
- Infrastructure setup
- Kubernetes deployment
- SSL/TLS configuration
- Monitoring setup (Prometheus + Grafana)
- Database migrations
- Backup and restore procedures
- Security hardening
- Troubleshooting

---

## 🎨 Screenshots

### Frontend
- **Home**: Landing page with platform overview
- **Markets**: Browse prediction markets
- **Trading**: Order book and trade execution
- **Sportsbook**: Multi-book odds comparison
- **Arbitrage**: Profit opportunities dashboard
- **Dashboard**: User portfolio overview

### Admin Panel
- **Dashboard**: Platform statistics
- **Markets**: Create and manage markets
- **Users**: User management
- **Sportsbooks**: Provider configuration
- **Analytics**: Betting analytics and insights

---

## 🗺️ Roadmap

### Phase 1: MVP (Complete)
- ✅ Prediction market infrastructure
- ✅ Order book trading
- ✅ Sportsbook integration
- ✅ Arbitrage detection
- ✅ Hedge calculator
- ✅ Web and admin interfaces

### Phase 2: Enhancement
- [ ] Mobile apps (Flutter iOS/Android)
- [ ] Advanced charting and analytics
- [ ] Social features (comments, sharing)
- [ ] Portfolio analytics
- [ ] Tax reporting
- [ ] API for third-party integrations

### Phase 3: Scaling
- [ ] Blockchain integration (optional)
- [ ] Decentralized oracle integration
- [ ] Advanced market types (scalar, categorical)
- [ ] Automated market makers (AMM)
- [ ] Liquidity mining programs
- [ ] DAO governance

---

## ⚠️ Legal & Compliance

### **CRITICAL WARNINGS**

This platform operates in heavily regulated industries:

1. **Prediction Markets**: Regulated by CFTC in the US
2. **Sports Betting**: State-by-state licensing required
3. **Money Transmission**: FinCEN and state licenses required
4. **Credential Storage**: Violates sportsbook Terms of Service

### Required Before Production

- [ ] CFTC registration (if applicable)
- [ ] State gambling licenses
- [ ] Money transmitter licenses (48+ jurisdictions)
- [ ] Legal review by gaming law attorney
- [ ] KYC/AML compliance program
- [ ] Age and geolocation verification
- [ ] Responsible gaming features
- [ ] Problem gambling resources
- [ ] Terms of Service
- [ ] Privacy Policy
- [ ] Proper insurance and bonding

### Recommended Approach

**Operate as odds comparison tool ONLY:**
- Display publicly available odds
- NO credential storage
- NO automated sportsbook access
- NO bet placement through platform
- User-initiated actions only

**See [LEGAL_DISCLAIMER.md](./LEGAL_DISCLAIMER.md) for complete legal requirements.**

---

## 🤝 Contributing

This is a demonstration project. Contributions welcome:

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

### Development Guidelines

- Follow Go best practices
- Write tests for new features
- Update documentation
- Follow existing code style
- Ensure all tests pass

---

## 📄 License

[Specify your license here - e.g., MIT, Apache 2.0]

**This software is provided "AS IS" without warranty of any kind. See [LEGAL_DISCLAIMER.md](./LEGAL_DISCLAIMER.md) for complete disclaimer.**

---

## 🆘 Support

### Issues
Report bugs and feature requests via GitHub Issues.

### Questions
- Check [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md)
- Review [backend/README.md](./backend/README.md)
- Check service-specific documentation

### Security Issues
Report security vulnerabilities privately to: [security email]

---

## 📞 Contact

- **Project Lead**: [Your Name]
- **Email**: [contact email]
- **Website**: [project website]

---

## 🙏 Acknowledgments

- Kalshi for prediction market inspiration
- Major sportsbooks for odds data structure reference
- Open source community for amazing tools

---

## ⚡ Quick Links

- [LEGAL DISCLAIMER](./LEGAL_DISCLAIMER.md) - **READ THIS FIRST**
- [Deployment Guide](./DEPLOYMENT_GUIDE.md)
- [Backend Documentation](./backend/README.md)
- [Frontend Documentation](./frontend-web/README.md)
- [Admin Panel Documentation](./admin-panel/README.md)

---

**Built with ❤️ for demonstration and educational purposes**

**Last Updated: January 2025**

---

## 📊 Project Statistics

- **Lines of Code**: 15,000+
- **Services**: 9 microservices
- **API Endpoints**: 50+
- **Database Tables**: 20+
- **React Components**: 30+
- **Languages**: Go, TypeScript, SQL
- **Test Coverage**: [Add coverage %]

---

**Remember: This is demonstration software. Do not use in production without proper legal compliance. See LEGAL_DISCLAIMER.md.**
