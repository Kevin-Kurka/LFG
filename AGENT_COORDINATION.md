# Agent Coordination Guide

## Overview

This document defines how 6 parallel agents will coordinate to deliver the LFG Platform in 3-4 weeks.

---

## Agent Roster & Responsibilities

### Agent 1: INFRA (DevOps Specialist)
**Primary Focus**: Infrastructure, CI/CD, Observability

**Responsibilities**:
- Docker containerization of all services
- CI/CD pipeline setup and maintenance
- Monitoring and logging infrastructure
- Deployment automation
- Performance optimization

**Key Deliverables**:
- Fully containerized development environment
- Automated CI/CD pipeline
- Prometheus + Grafana observability stack
- Kubernetes deployment manifests

**Communication Priority**: Works with all agents on deployment needs

---

### Agent 2: DATA (Database Architect)
**Primary Focus**: Database layer, models, data integrity

**Responsibilities**:
- Database schema design and migrations
- Repository pattern implementation
- Data models and validation
- Database optimization and indexing
- Test infrastructure

**Key Deliverables**:
- Production-ready database schema with indexes
- Repository interfaces for all entities
- Shared models package
- Database migration system
- Test data fixtures

**Communication Priority**: Critical handoff to ACCOUNTS and TRADING agents

---

### Agent 3: AUTH (Security Engineer)
**Primary Focus**: Authentication, authorization, security

**Responsibilities**:
- JWT authentication system
- API gateway security
- Rate limiting and CORS
- Input validation and sanitization
- Security audits and penetration testing
- Compliance (GDPR, audit logging)

**Key Deliverables**:
- Secure API gateway with auth middleware
- JWT package for all services
- RBAC authorization system
- Security-hardened platform

**Communication Priority**: Handoff JWT package to ACCOUNTS early in Week 1

---

### Agent 4: ACCOUNTS (Backend Developer - User Domain)
**Primary Focus**: User service, wallet service

**Responsibilities**:
- User registration and authentication
- User profile management
- Wallet balance and transactions
- Credit exchange service
- NATS event publishing for user/wallet events

**Key Deliverables**:
- Fully functional user service
- Fully functional wallet service
- Integration with AUTH for secure registration/login
- Event-driven wallet updates

**Communication Priority**: Coordinates with TRADING on wallet balance checks

---

### Agent 5: TRADING (Backend Developer - Trading Domain)
**Primary Focus**: Markets, orders, matching engine

**Responsibilities**:
- Market service (listing, details, orderbook)
- Order service (place, cancel, status)
- Matching engine (core trading logic)
- WebSocket notification service
- gRPC inter-service communication

**Key Deliverables**:
- Production-grade matching engine
- Complete order lifecycle management
- Real-time trading notifications
- Market creation and resolution

**Communication Priority**: Most complex stream, needs early start on matching engine

---

### Agent 6: FRONTEND (Full-Stack Developer - UI/UX)
**Primary Focus**: Mobile app (Flutter) and admin panel (React)

**Responsibilities**:
- Flutter mobile application
- React admin panel
- API integration
- Real-time WebSocket updates
- Mobile app store submission

**Key Deliverables**:
- Production-ready mobile app
- Admin panel for market management
- Published apps on App Store and Google Play

**Communication Priority**: Can work independently on UI, integrates with APIs as they become available

---

## Handoff Protocol

### Handoff 1: Shared Models Package
**From**: DATA (Agent 2)
**To**: All agents
**Timing**: End of Day 1
**Deliverable**: `backend/shared/models` package

**Acceptance Criteria**:
- [ ] All models defined with JSON tags
- [ ] Validation tags added
- [ ] go.mod updated with dependencies
- [ ] Models compile without errors
- [ ] Basic unit tests for model validation

**Notification**: Post in #lfg-general when complete

---

### Handoff 2: JWT Authentication Package
**From**: AUTH (Agent 3)
**To**: ACCOUNTS (Agent 4), TRADING (Agent 5)
**Timing**: End of Day 1
**Deliverable**: `backend/shared/auth` package

**Acceptance Criteria**:
- [ ] GenerateToken(userID, email) function
- [ ] ValidateToken(tokenString) function
- [ ] RefreshToken functionality
- [ ] Unit tests with >90% coverage
- [ ] README with usage examples

**Notification**: Tag @ACCOUNTS and @TRADING in #lfg-backend

---

### Handoff 3: Database Connection Package
**From**: DATA (Agent 2)
**To**: ACCOUNTS (Agent 4), TRADING (Agent 5)
**Timing**: End of Day 2
**Deliverable**: `backend/shared/db` package + all repositories

**Acceptance Criteria**:
- [ ] Connection pool configured
- [ ] Repository interfaces documented
- [ ] Example usage in README
- [ ] Integration tests passing
- [ ] Seed data populated

**Notification**: Post in #lfg-backend with migration instructions

---

### Handoff 4: Docker Compose Environment
**From**: INFRA (Agent 1)
**To**: All agents
**Timing**: End of Day 1
**Deliverable**: `docker-compose.yml` + Makefile

**Acceptance Criteria**:
- [ ] All services start with `make docker-up`
- [ ] Database initializes with schema
- [ ] NATS server running
- [ ] Health checks working
- [ ] Documentation in QUICKSTART.md updated

**Notification**: Post in #lfg-general when ready for team to use

---

### Handoff 5: User Registration/Login APIs
**From**: ACCOUNTS (Agent 4)
**To**: FRONTEND (Agent 6)
**Timing**: End of Day 3
**Deliverable**: Working `/register` and `/login` endpoints

**Acceptance Criteria**:
- [ ] Registration returns JWT token
- [ ] Login validates credentials
- [ ] Password hashing with bcrypt
- [ ] Integration tests passing
- [ ] Postman collection provided

**Notification**: Tag @FRONTEND in #lfg-backend with API docs

---

### Handoff 6: Market Listing API
**From**: TRADING (Agent 5)
**To**: FRONTEND (Agent 6)
**Timing**: End of Day 3
**Deliverable**: `/markets` endpoint with filtering

**Acceptance Criteria**:
- [ ] Returns paginated market list
- [ ] Filters work (status, search)
- [ ] Mock data available
- [ ] Response matches contract
- [ ] Postman collection provided

**Notification**: Tag @FRONTEND in #lfg-backend with API docs

---

### Handoff 7: Order Placement API
**From**: TRADING (Agent 5)
**To**: FRONTEND (Agent 6)
**Timing**: End of Day 4
**Deliverable**: `/orders/place` endpoint

**Acceptance Criteria**:
- [ ] Accepts LIMIT and MARKET orders
- [ ] Validates user balance
- [ ] Returns order ID
- [ ] WebSocket notification on fill
- [ ] Error handling documented

**Notification**: Tag @FRONTEND in #lfg-backend with API docs

---

### Handoff 8: Matching Engine gRPC Service
**From**: TRADING (Agent 5)
**To**: ACCOUNTS (Agent 4) for integration testing
**Timing**: End of Week 2
**Deliverable**: Fully functional matching engine

**Acceptance Criteria**:
- [ ] gRPC server running
- [ ] Order matching working (price-time priority)
- [ ] Trade events published to NATS
- [ ] Performance: >1000 orders/sec
- [ ] Integration tests provided

**Notification**: Post in #lfg-backend for integration testing

---

## Daily Coordination

### Daily Standup (9:00 AM - 15 minutes)

**Format** (each agent, 2 minutes max):
1. What I completed yesterday
2. What I'm working on today
3. Blockers or dependencies needed

**Example**:
> Agent 2 (DATA): Yesterday I completed the database schema migrations (DATA-1.1, DATA-1.2). Today I'm implementing the repository pattern (DATA-1.4). No blockers, but AUTH will need the models package by EOD.

### Integration Sync (3:00 PM - 30 minutes)

**Purpose**: Coordinate on cross-service integration

**Attendees**: All agents with integration dependencies that day

**Agenda**:
- Review handoffs completed today
- Test integrations together
- Resolve any interface mismatches
- Plan next day's integrations

**Example Topics**:
- Day 2: AUTH and ACCOUNTS test JWT integration
- Day 3: TRADING and FRONTEND test market listing API
- Day 4: ACCOUNTS and TRADING test wallet balance check

---

## Communication Channels

### Slack Channels

**#lfg-general**
- Announcements
- Cross-team updates
- Handoff notifications
- Daily standup notes

**#lfg-infra**
- Docker/deployment issues
- CI/CD pipeline
- Environment configuration

**#lfg-backend**
- API design discussions
- Database questions
- gRPC/NATS integration
- Backend handoffs

**#lfg-frontend**
- UI/UX discussions
- API integration help
- Mobile app issues

**#lfg-blockers**
- Critical blockers only
- Tag relevant agents
- Requires <2 hour response time

### Slack Naming Conventions

Tag agents with `@AGENTNAME`:
- @INFRA, @DATA, @AUTH, @ACCOUNTS, @TRADING, @FRONTEND

Urgent issues: Use `🚨 BLOCKER:`
Example: "🚨 BLOCKER: @DATA database migrations failing in CI"

### GitHub

**Branch Naming**:
- `infra/INFRA-1.2`
- `data/DATA-1.4`
- `auth/AUTH-1.1`
- etc.

**PR Naming**:
- `feat(infra): INFRA-1.2 - Add docker-compose configuration`
- `fix(auth): AUTH-1.3 - Fix JWT expiration validation`

**PR Labels**:
- `stream:infra`, `stream:data`, `stream:auth`, etc.
- `priority:critical`, `priority:high`, `priority:medium`
- `status:blocked`, `status:review`, `status:merged`

---

## Conflict Resolution

### Code Conflicts

**Scenario**: Two agents modify the same file

**Resolution**:
1. Agent who pushed second rebases on main
2. Resolve conflicts locally
3. Request review from first agent
4. Merge after approval

**Prevention**: Use shared packages for common code, coordinate in daily standup

---

### Design Disagreements

**Scenario**: Agents disagree on interface design

**Resolution**:
1. Discuss in #lfg-backend or relevant channel
2. If no consensus in 30 minutes, escalate to team vote
3. Majority wins, document decision
4. Revisit in retrospective if needed

---

### Blocking Dependencies

**Scenario**: Agent A blocks Agent B's work

**Resolution**:
1. Agent B posts in #lfg-blockers tagging Agent A
2. Agent A provides:
   - Mock interface immediately (<1 hour)
   - ETA for real implementation
   - Ping when unblocked
3. Agent B proceeds with mock
4. Integration test when both sides ready

**Example**:
```
@ACCOUNTS blocked on @AUTH JWT package for user registration (ACCOUNTS-1.1)

@AUTH response:
- Mock JWT package pushed to shared/auth/mock
- Real implementation ETA: EOD today
- Will ping when AUTH-1.1 merged
```

---

## Integration Testing Strategy

### Service-to-Service Testing

**When**: After each handoff

**Who**: Both agents involved in handoff

**How**:
1. Agent A completes feature, pushes to branch
2. Agent B pulls branch and tests locally
3. Both join 15-minute integration session
4. Test happy path and error cases
5. Document any issues found
6. Fix issues before merging

**Example**:
- ACCOUNTS and AUTH test user registration with JWT
- TRADING and ACCOUNTS test wallet balance check
- FRONTEND and TRADING test order placement

---

### End-to-End Testing

**When**: End of each phase

**Who**: All agents

**Duration**: 1-2 hours

**Scenarios**:
1. User registers and logs in
2. User browses markets
3. User places order
4. Order executes (matched)
5. User receives notification
6. User checks balance (updated)

**Success Criteria**: Full flow works without errors

---

## Weekly Rituals

### Monday Planning (10:00 AM - 1 hour)

**Agenda**:
1. Review last week's velocity
2. Review this week's tasks
3. Identify risks and dependencies
4. Assign tasks if not already assigned
5. Set weekly goals

**Output**: Updated TASK_TRACKING.md with assignments

---

### Friday Retrospective (4:00 PM - 1 hour)

**Agenda**:
1. What went well this week?
2. What could be improved?
3. Action items for next week
4. Celebrate wins

**Format**: Start/Stop/Continue

**Example**:
- Start: Writing integration tests earlier
- Stop: Pushing directly to main (use PRs)
- Continue: Daily standups are helpful

---

## Emergency Protocols

### Production Incident (Post-Launch)

**Severity Levels**:
- **P0 (Critical)**: System down, all hands on deck
- **P1 (High)**: Major feature broken, fix within 4 hours
- **P2 (Medium)**: Minor feature broken, fix within 1 day
- **P3 (Low)**: Cosmetic issue, fix in next sprint

**On-Call Rotation**: TBD after launch

---

### Blocking Bug in Development

**Process**:
1. Post in #lfg-blockers with `🚨 BLOCKER:`
2. Tag relevant agent(s)
3. Pause other work if critical path
4. Pair program to resolve quickly
5. Document root cause
6. Add regression test

---

## Knowledge Sharing

### Documentation

**Each agent must document**:
- API contracts (OpenAPI specs)
- Database schemas and migrations
- Architecture decisions (ADRs)
- Setup instructions
- Troubleshooting guides

**Location**: `/docs` folder in repository

---

### Code Reviews

**Process**:
1. All PRs require 1 approval
2. Critical path PRs require 2 approvals
3. Review within 4 hours during work hours
4. Focus on:
   - Correctness
   - Security
   - Performance
   - Maintainability

**Review Checklist**:
- [ ] Tests included and passing
- [ ] No security vulnerabilities
- [ ] Follows code style
- [ ] Documentation updated
- [ ] No secrets committed

---

### Pair Programming

**When**:
- Complex features (e.g., matching engine)
- Debugging tricky issues
- Onboarding to new area of codebase

**Tools**: VS Code Live Share, Tuple, Zoom

---

## Metrics & Tracking

### Velocity Tracking

**Measured Daily**:
- Tasks completed (count)
- Hours logged (actual vs. estimated)
- Blockers encountered (count)
- PRs merged (count)

**Dashboard**: GitHub Projects board + TASK_TRACKING.md

---

### Quality Metrics

**Tracked Weekly**:
- Test coverage (target: >80%)
- Code review turnaround time (target: <4 hours)
- Bug count (by severity)
- CI build success rate (target: >95%)

---

### Team Health

**Tracked Weekly** (in retrospective):
- Team satisfaction (1-5 scale)
- Workload balance (hours per agent)
- Blocker frequency
- Communication effectiveness

---

## Success Criteria

### Phase 1 Complete (End of Week 1)

**Technical**:
- [ ] All services running in Docker
- [ ] User can register and login
- [ ] Markets can be listed
- [ ] Database fully seeded
- [ ] CI pipeline green

**Process**:
- [ ] All agents on track (within 10% of estimates)
- [ ] Zero critical blockers
- [ ] Daily standups happening
- [ ] Handoffs smooth

---

### Phase 2 Complete (End of Week 2)

**Technical**:
- [ ] Matching engine functional
- [ ] End-to-end trading flow works
- [ ] WebSocket notifications working
- [ ] Test coverage >70%
- [ ] Security audit passed

**Process**:
- [ ] Integration tests running daily
- [ ] All handoffs on time
- [ ] Code review turnaround <4 hours
- [ ] Team velocity stable

---

### Phase 3 Complete (End of Week 3-4)

**Technical**:
- [ ] Load testing passed (10k concurrent users)
- [ ] Production deployment successful
- [ ] Mobile apps in beta
- [ ] Observability stack deployed
- [ ] Documentation complete

**Process**:
- [ ] Team working well together
- [ ] Minimal technical debt
- [ ] Ready for production launch
- [ ] Post-MVP roadmap defined

---

## Contact & Escalation

### Agent Leads

- **INFRA Lead**: TBD
- **DATA Lead**: TBD
- **AUTH Lead**: TBD
- **ACCOUNTS Lead**: TBD
- **TRADING Lead**: TBD
- **FRONTEND Lead**: TBD

### Project Manager

- **Name**: TBD
- **Slack**: @pm
- **Email**: pm@lfg.example.com
- **Availability**: Mon-Fri 9am-6pm

### Escalation Path

1. **Level 1**: Ask in stream channel
2. **Level 2**: Post in #lfg-blockers
3. **Level 3**: Tag agent lead
4. **Level 4**: Tag project manager
5. **Level 5**: Emergency call (P0 only)

---

## Quick Reference

### Key Documents
- `IMPLEMENTATION_PLAN.md` - Full project plan
- `TASK_TRACKING.md` - Task status
- `QUICKSTART.md` - Developer setup
- `AGENT_COORDINATION.md` - This document

### Key Commands
- `make docker-up` - Start all services
- `make test` - Run all tests
- `make db-migrate` - Run database migrations

### Key Times
- **9:00 AM** - Daily standup
- **3:00 PM** - Integration sync
- **10:00 AM Monday** - Weekly planning
- **4:00 PM Friday** - Retrospective

### Key Channels
- `#lfg-general` - Announcements
- `#lfg-backend` - Backend discussion
- `#lfg-frontend` - Frontend discussion
- `#lfg-blockers` - Critical blockers

---

**Let's collaborate and ship this together! 🚀**

Last Updated: 2025-11-14
Questions? Ask in #lfg-general
