# Architecture Decision Records (ADR)

## ADR-001: Use JSON File for Data Storage (MVP)

**Status:** Accepted  
**Date:** 2025-10-24

### Context
Challenge provides static JSON data. Need to decide between using as-is vs migrating to database.

### Decision
Use embedded JSON file for MVP, design with repository pattern for future database migration.

### Rationale
- Challenge provides JSON data (signal of expected scope)
- 3-4 day timeline is tight for database setup
- Shows good judgment ("right tool for the problem")
- Reduces deployment risk
- Repository pattern enables future swap

### Consequences
- **Positive:** Fast development, zero infrastructure overhead, version-controlled data
- **Negative:** Read-only operations, not suitable for production scale
- **Mitigation:** Document database migration path in architecture writeup

---

## ADR-002: Use Go + chi Router for Backend

**Status:** Accepted  
**Date:** 2025-10-24

### Context
Challenge offers bonus points for Go. Need to balance learning curve vs speed.

### Decision
Use Go 1.22+ with chi router (single dependency).

### Rationale
- Earns bonus points mentioned in challenge
- chi is minimal, idiomatic Go
- Shows restraint (1 dependency vs heavy framework)
- Demonstrates good architectural judgment
- Standard library handles JSON, HTTP well

### Consequences
- **Positive:** Bonus points, shows Go competency, minimal deps
- **Negative:** Slightly slower than Node.js for rapid prototyping
- **Mitigation:** Keep handlers simple, clear structure

---

## ADR-003: Use React + Vite for Frontend

**Status:** Accepted  
**Date:** 2025-10-24

### Context
Need fast development while showcasing modern frontend skills.

### Decision
React 18 with Vite build tool.

### Rationale
- User is strongest in frontend development
- Vite offers fastest dev experience
- React has best ecosystem for charts/visualization
- Industry standard, familiar to reviewers

### Consequences
- **Positive:** Fast development, hot reload, excellent DX
- **Negative:** Bundle size larger than Svelte
- **Mitigation:** Code splitting if needed

---

## ADR-004: Use Recharts for Visualizations

**Status:** Accepted  
**Date:** 2025-10-24

### Context
Need chart library that's simple, responsive, and React-friendly.

### Decision
Use Recharts library for all visualizations.

### Rationale
- React-native (declarative API)
- Responsive by default
- Good balance of simplicity vs features
- Mobile-friendly out of box

### Alternatives Considered
- D3.js: Too complex, time-intensive
- Chart.js: Not React-native
- ApexCharts: More complex than needed

### Consequences
- **Positive:** Fast implementation, responsive, good docs
- **Negative:** Less customization than D3
- **Mitigation:** Sufficient for MVP needs

---

## ADR-005: Use TailwindCSS for Styling

**Status:** Accepted  
**Date:** 2025-10-24

### Context
Need fast, mobile-first styling solution.

### Decision
Use TailwindCSS utility-first framework.

### Rationale
- Mobile-first by default
- Rapid development
- No CSS naming conflicts
- Small production bundle (purged unused styles)
- Modern, professional look achievable quickly

### Consequences
- **Positive:** Fast styling, responsive utilities, modern aesthetic
- **Negative:** Learning curve if unfamiliar
- **Mitigation:** Good documentation, common patterns

---

## ADR-006: Server-Side Data Aggregation

**Status:** Accepted  
**Date:** 2025-10-24

### Context
With only 114 transactions, could aggregate on frontend or backend.

### Decision
Implement aggregation logic in backend service layer.

### Rationale
- Proper separation of concerns
- Scalable pattern (works with 1M transactions)
- Reusable API endpoints
- Shows backend competency
- Frontend stays thin

### Consequences
- **Positive:** Scalable, clean architecture, reusable APIs
- **Negative:** More backend code
- **Mitigation:** Worth the trade for architectural quality

---

## ADR-007: OpenAI for AI Advisor

**Status:** Accepted  
**Date:** 2025-10-24

### Context
Need LLM integration for financial advice feature.

### Decision
Use OpenAI GPT-3.5-turbo API with direct HTTP calls.

### Rationale
- Best quality responses
- Simplest API
- No SDK needed (just HTTP)
- Challenge explicitly allows any LLM provider

### Alternatives Considered
- Google Gemini: Free tier, but more complex API
- Anthropic Claude: Great for finance, requires API key

### Consequences
- **Positive:** High quality advice, simple integration
- **Negative:** Requires API key, costs money
- **Mitigation:** Implement fallback mock response if no key

---

## ADR-008: AWS EC2 for Backend Deployment

**Status:** Accepted  
**Date:** 2025-10-24

### Context
Challenge requires AWS hosting. Choose between EC2, Lambda, ECS.

### Decision
Use single EC2 instance for MVP, document serverless alternative.

### Rationale
- Simplest AWS option
- Reliable for demo (no cold starts)
- Easy to demonstrate
- Less infrastructure complexity

### Alternatives Considered
- Lambda + API Gateway: More impressive, but complex deployment
- ECS/Fargate: Production-like, but heavier infrastructure

### Consequences
- **Positive:** Simple, reliable, easy to troubleshoot
- **Negative:** Not auto-scaling, manual management
- **Mitigation:** Document serverless path in architecture writeup

---

## ADR-009: S3 + CloudFront for Frontend

**Status:** Accepted  
**Date:** 2025-10-24

### Context
Need to host static React build in AWS.

### Decision
Use S3 for storage, CloudFront for global CDN.

### Rationale
- Standard AWS pattern for static sites
- CloudFront provides SSL, caching, global distribution
- Cost-effective
- Professional setup

### Consequences
- **Positive:** Fast load times, SSL, scalable, cheap
- **Negative:** CloudFront cache invalidation needed on deploys
- **Mitigation:** Automated invalidation in deployment script

---

## ADR-010: Repository Pattern for Data Access

**Status:** Accepted  
**Date:** 2025-10-24

### Context
Need to design for extensibility while keeping MVP simple.

### Decision
Implement repository pattern with interface and JSON implementation.

### Rationale
- Enables future database swap without changing handlers
- Shows architectural maturity
- Clean separation of concerns
- Industry best practice

### Code Structure
```go
type TransactionRepository interface {
    GetAll() ([]Transaction, error)
    GetByDateRange(start, end time.Time) ([]Transaction, error)
}

// MVP
type JSONRepository struct { ... }

// Future (documented)
type PostgresRepository struct { ... }
```

### Consequences
- **Positive:** Extensible, testable, professional architecture
- **Negative:** Slight over-engineering for 114 records
- **Mitigation:** Worth it to demonstrate design thinking

---

## ADR-011: Docker Compose for Local Development

**Status:** Accepted  
**Date:** 2025-10-24

### Context
Need easy way to run full stack locally and demonstrate DevOps thinking.

### Decision
Provide Docker Compose setup for backend + frontend.

### Rationale
- Shows DevOps competency
- Easy for reviewers to run locally
- Consistent environment
- Minimal effort to implement

### Consequences
- **Positive:** Professional setup, easy demo, DevOps points
- **Negative:** Another config to maintain
- **Mitigation:** Simple compose file, well worth the effort

---

## ADR-012: No Authentication for MVP

**Status:** Accepted  
**Date:** 2025-10-24

### Context
Challenge doesn't mention auth. Data is read-only static JSON.

### Decision
No authentication/authorization in MVP.

### Rationale
- Not in requirements
- Data is sample/demo (not sensitive)
- Time better spent on required features
- Can document as Phase 3 feature

### Consequences
- **Positive:** Faster development, simpler architecture
- **Negative:** Not production-ready
- **Mitigation:** Document auth strategy (Cognito) in future improvements

---

## ADR-013: Write Tests or Document Testing Strategy

**Status:** Pending  
**Date:** 2025-10-24

### Context
Challenge allows either writing tests OR explaining testing strategy.

### Decision
TBD based on Day 4 time availability.

### Options
1. Write actual tests (backend unit + integration, frontend component)
2. Write comprehensive testing-strategy.md document

### Leaning Toward
Option 1 if time permits - more impressive, shows actual capability.

---

## Summary: MVP Architecture

```
Frontend: React + Vite + TailwindCSS + Recharts
         ↓ (axios)
Backend:  Go + chi router + stdlib
         ↓ (repository pattern)
Data:     Embedded JSON file
         ↓ (HTTP)
AI:       OpenAI GPT-3.5-turbo API

Deployment:
- Frontend: S3 + CloudFront
- Backend: EC2 (t3.micro)
- Local: Docker Compose

Total External Dependencies:
- Backend: 1 (chi)
- Frontend: 4 (react, recharts, axios, tailwind)
```

This architecture prioritizes:
1. ✅ Meeting all challenge requirements
2. ✅ Demonstrating good judgment (scope, tools)
3. ✅ Showing breadth (frontend, backend, DevOps, AI)
4. ✅ Extensibility (repository pattern, clean architecture)
5. ✅ Reliability (simple, proven patterns)

