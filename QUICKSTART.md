# LFG Platform - Quick Start Guide

## For New Developers Joining the Project

Welcome to the LFG Platform development team! This guide will get you up and running in under 30 minutes.

---

## Prerequisites

Before you begin, ensure you have the following installed:

### Required Software
- **Git**: Version 2.30+
- **Go**: Version 1.24.3 (exact version)
- **Docker**: Version 20.10+
- **Docker Compose**: Version 2.0+
- **Node.js**: Version 18+ (for admin panel)
- **Flutter**: Version 3.0+ (for mobile app)
- **Make**: GNU Make (usually pre-installed on macOS/Linux)

### Optional (Recommended)
- **VSCode** with Go and Flutter extensions
- **Postman** or **Insomnia** for API testing
- **pgAdmin** or **TablePlus** for database inspection
- **k6** for load testing

---

## Initial Setup (One-time)

### 1. Clone the Repository

```bash
git clone https://github.com/Kevin-Kurka/LFG.git
cd LFG
```

### 2. Checkout Your Development Branch

```bash
# Your branch was already created
git checkout claude/code-review-status-011BB6FtFnEDkr7u4Wkffp7F

# Or create your own feature branch
git checkout -b your-stream-name/your-task-id
```

### 3. Review Your Assigned Stream

Open and read:
1. `IMPLEMENTATION_PLAN.md` - Full project plan
2. `TASK_TRACKING.md` - Your specific tasks
3. Find your stream assignment (INFRA, DATA, AUTH, ACCOUNTS, TRADING, or FRONTEND)

### 4. Set Up Environment Variables

```bash
# Copy example env file (will be created by INFRA team)
cp .env.example .env

# Edit .env with your local settings
# For now, defaults are fine for local development
```

---

## Development Environment Setup

### Option A: Docker-Based Development (Recommended)

**Prerequisites**: Docker tasks from INFRA stream must be completed

```bash
# Start all services
make docker-up

# Or manually:
docker-compose up -d

# Check all services are running
docker-compose ps

# View logs
docker-compose logs -f
```

**Services will be available at**:
- API Gateway: http://localhost:8000
- User Service: http://localhost:8080
- Wallet Service: http://localhost:8081
- Order Service: http://localhost:8082
- Market Service: http://localhost:8083
- Credit Exchange: http://localhost:8084
- Notification Service: http://localhost:8085
- PostgreSQL: localhost:5432
- NATS: localhost:4222
- Admin Panel: http://localhost:3000

### Option B: Local Development (For Early Development)

**Backend Services**:

```bash
# Terminal 1: Start PostgreSQL
docker run -d \
  --name lfg-postgres \
  -e POSTGRES_DB=lfg \
  -e POSTGRES_USER=lfg \
  -e POSTGRES_PASSWORD=lfg_dev_password \
  -p 5432:5432 \
  postgres:15

# Terminal 2: Run database migrations
make db-migrate

# Terminal 3: Start NATS
docker run -d \
  --name lfg-nats \
  -p 4222:4222 \
  -p 8222:8222 \
  nats:latest

# Terminal 4+: Start each service
cd backend/user-service && go run main.go
cd backend/api-gateway && go run main.go
# ... etc for other services
```

**Frontend Development**:

```bash
# Admin Panel
cd admin-panel
npm install
npm start

# Flutter App
cd frontend/lfg_app
flutter pub get
flutter run -d chrome  # For web
flutter run  # For mobile (simulator must be running)
```

---

## Your First Task

### 1. Claim a Task

1. Go to `TASK_TRACKING.md`
2. Find your stream's first task
3. Change status from 🔴 TODO to 🟡 IN PROGRESS
4. Add your name to Assignee column

Example:
```markdown
| INFRA-1.2 | Create docker-compose.yml | Your Name | 🟡 IN PROGRESS | 4h | - |
```

### 2. Create a Feature Branch

```bash
git checkout -b stream-name/task-id
# Example: git checkout -b infra/INFRA-1.2
```

### 3. Implement the Task

Follow the detailed specification in `IMPLEMENTATION_PLAN.md` for your task.

**General Guidelines**:
- Write tests alongside your code
- Follow Go conventions (use `gofmt`, `golint`)
- Add comments for complex logic
- Update documentation if needed

### 4. Test Your Changes

```bash
# Run tests for your service
cd backend/your-service
go test ./...

# Run linting
golangci-lint run

# Test manually with curl or Postman
curl http://localhost:8080/your-endpoint
```

### 5. Commit and Push

```bash
# Stage your changes
git add .

# Commit with descriptive message
git commit -m "feat(stream): Task ID - Brief description

- Detailed change 1
- Detailed change 2

Closes TASK-ID"

# Push to remote
git push -u origin stream-name/task-id
```

### 6. Create Pull Request

1. Go to GitHub
2. Create PR from your branch to `main` (or designated base branch)
3. Fill in PR template:
   - Link to task in TASK_TRACKING.md
   - Describe what you implemented
   - List any breaking changes
   - Add screenshots if UI changes
4. Request review from stream lead or relevant agent

### 7. Update Task Status

After PR is merged:
1. Update `TASK_TRACKING.md`
2. Change status to 🟢 DONE
3. Move to next task

---

## Stream-Specific Setup

### INFRA Stream

**Focus**: DevOps, infrastructure, CI/CD

```bash
# Install additional tools
brew install golangci-lint  # macOS
# or
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Test Docker setup
docker --version
docker-compose --version
```

**Key Files to Modify**:
- `Dockerfile` (in each service directory)
- `docker-compose.yml`
- `.github/workflows/ci.yml`
- `Makefile`

---

### DATA Stream

**Focus**: Database, models, repositories

```bash
# Install database tools
brew install postgresql  # For psql client
brew install golang-migrate

# Connect to database
psql -h localhost -U lfg -d lfg

# Run migrations
make db-migrate

# Seed data
make db-seed
```

**Key Files to Modify**:
- `database/migrations/*.sql`
- `database/init.sql`
- `backend/shared/models/*.go`
- `backend/shared/db/*.go`
- `backend/shared/repository/*.go`

---

### AUTH Stream

**Focus**: Authentication, authorization, security

```bash
# Install security scanning tools
brew install trivy
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Generate JWT secret
openssl rand -base64 32
```

**Key Files to Modify**:
- `backend/shared/auth/*.go`
- `backend/api-gateway/middleware/*.go`
- `backend/api-gateway/main.go`

**Test Authentication**:
```bash
# Register a user
curl -X POST http://localhost:8000/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"securepass123"}'

# Should return JWT token
```

---

### ACCOUNTS Stream

**Focus**: User service, wallet service

```bash
# No special tools needed
# Just Go and database access
```

**Key Files to Modify**:
- `backend/user-service/main.go`
- `backend/user-service/handlers/*.go`
- `backend/wallet-service/main.go`
- `backend/wallet-service/handlers/*.go`

**Test User Registration**:
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

---

### TRADING Stream

**Focus**: Markets, orders, matching engine

```bash
# Install gRPC tools
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install protoc compiler
brew install protobuf

# Generate gRPC code from proto
cd backend/matching-engine
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

**Key Files to Modify**:
- `backend/matching-engine/engine/*.go`
- `backend/matching-engine/proto/*.proto`
- `backend/order-service/main.go`
- `backend/market-service/main.go`

**Test Order Placement**:
```bash
# Place an order
curl -X POST http://localhost:8082/orders/place \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"contract_id":"...","type":"LIMIT","quantity":10,"limit_price":0.55}'
```

---

### FRONTEND Stream

**Focus**: Flutter app, React admin panel

```bash
# Flutter setup
flutter doctor  # Check Flutter installation
flutter devices  # List available devices

# React setup
cd admin-panel
npm install

# Start development servers
npm start  # React (admin panel)
flutter run  # Flutter (mobile app)
```

**Key Files to Modify**:
- `frontend/lfg_app/lib/*.dart`
- `admin-panel/src/*.tsx`

**Test Mobile App**:
```bash
# Run on web browser
cd frontend/lfg_app
flutter run -d chrome

# Run on iOS simulator (macOS only)
flutter run -d iPhone

# Run on Android emulator
flutter run -d emulator
```

---

## Common Commands

### Database

```bash
# Run migrations
make db-migrate

# Rollback migration
make db-rollback

# Seed database
make db-seed

# Reset database (drop and recreate)
make db-reset

# Connect to database
make db-shell
```

### Development

```bash
# Build all services
make build

# Run tests
make test

# Run linting
make lint

# Format code
make fmt

# Run all checks (test + lint + fmt)
make check
```

### Docker

```bash
# Start all services
make docker-up

# Stop all services
make docker-down

# Rebuild and restart
make docker-restart

# View logs
make docker-logs

# Clean everything (remove volumes)
make docker-clean
```

---

## Testing Guidelines

### Unit Tests

```go
// backend/user-service/handlers/register_test.go
package handlers_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestRegisterHandler(t *testing.T) {
    // Arrange
    // ... setup

    // Act
    // ... call function

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### Integration Tests

```go
// backend/user-service/integration_test.go
func TestRegistrationFlow(t *testing.T) {
    // Use testcontainers for database
    // Test full flow: register -> login -> get profile
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestRegisterHandler ./...
```

---

## Debugging

### Backend Services

```bash
# Run with debugger (Delve)
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./backend/user-service

# View service logs
docker-compose logs -f user-service

# Check service health
curl http://localhost:8080/health
```

### Database Issues

```bash
# Connect to database
docker exec -it lfg-postgres psql -U lfg -d lfg

# Check tables
\dt

# Describe table
\d users

# Check connections
SELECT * FROM pg_stat_activity;
```

### NATS Issues

```bash
# Check NATS server
curl http://localhost:8222/varz

# Monitor NATS messages (install nats CLI)
brew install nats-io/nats-tools/nats
nats sub ">"  # Subscribe to all messages
```

---

## Code Style & Conventions

### Go

- Use `gofmt` for formatting (automatic)
- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use meaningful variable names
- Add comments for exported functions
- Error messages should be lowercase, no punctuation

```go
// Good
func GetUser(id uuid.UUID) (*User, error) {
    if id == uuid.Nil {
        return nil, errors.New("user id cannot be nil")
    }
    // ...
}

// Bad
func get_user(ID uuid.UUID) (*User, error) {
    // Missing comment
    // Snake case instead of camel case
}
```

### TypeScript/React

- Use functional components
- Use TypeScript strict mode
- Follow Airbnb style guide
- Use meaningful component names

```typescript
// Good
interface UserProps {
  user: User;
  onUpdate: (user: User) => void;
}

const UserProfile: React.FC<UserProps> = ({ user, onUpdate }) => {
  // ...
};
```

### Flutter/Dart

- Follow official Dart style guide
- Use StatelessWidget where possible
- Organize by feature, not by type

```dart
// Good
class MarketListScreen extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    // ...
  }
}
```

---

## Getting Help

### Documentation
- `IMPLEMENTATION_PLAN.md` - Detailed task specifications
- `TASK_TRACKING.md` - Task status and assignments
- `README.md` - Project overview
- Code comments - In-line documentation

### Communication
- **Slack**: #lfg-general for general questions
- **Slack**: #lfg-your-stream for stream-specific questions
- **GitHub Issues**: For bugs and feature requests
- **PR Comments**: For code-specific questions

### Daily Standup (9:00 AM)
- What did you complete yesterday?
- What are you working on today?
- Any blockers?

### Debugging Assistance
1. Check service logs
2. Check database state
3. Ask in stream channel
4. Tag relevant agent if blocked
5. Escalate to daily standup if critical

---

## FAQ

**Q: I can't start my service, port is already in use?**
A: Another service is using the port. Stop it: `lsof -ti:8080 | xargs kill -9`

**Q: Database migrations failing?**
A: Check if database is running: `docker ps | grep postgres`

**Q: Tests are failing locally but not in CI?**
A: Ensure you have the same Go version: `go version`

**Q: How do I know which task to work on next?**
A: Check `TASK_TRACKING.md` for your stream, pick the next unblocked task

**Q: My PR has merge conflicts?**
A: Rebase on main: `git fetch origin && git rebase origin/main`

**Q: I need to share code with another agent?**
A: Put shared code in `backend/shared/*` and coordinate in daily standup

**Q: Frontend needs API that's not ready yet?**
A: Use mock data or create a mock API response until backend is ready

---

## Useful Links

- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [NATS Documentation](https://docs.nats.io/)
- [Flutter Documentation](https://flutter.dev/docs)
- [React Documentation](https://react.dev/)
- [Docker Documentation](https://docs.docker.com/)

---

## Ready to Start?

1. ✅ Clone repository
2. ✅ Review your stream in IMPLEMENTATION_PLAN.md
3. ✅ Set up development environment
4. ✅ Claim your first task in TASK_TRACKING.md
5. ✅ Create feature branch
6. ✅ Start coding!

**Welcome to the team! Let's build something amazing! 🚀**

---

Last Updated: 2025-11-14
Questions? Ask in #lfg-general on Slack
