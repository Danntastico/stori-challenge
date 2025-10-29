# Implementation Roadmap - Quick Reference

## Pre-Flight Checklist

- [x] User stories defined and prioritized
- [x] Technical spec completed
- [x] Architecture decisions documented
- [x] Scope aligned with challenge requirements
- [x] GitHub repository created
- [ ] OpenAI API key obtained (or plan for mock)
- [ ] AWS account access verified
- [ ] Development environment ready

---

## Day 1: Backend Foundation (6-8 hours)

### Morning Session (3-4 hours)

#### 1. Project Setup (30 min)
```bash
mkdir stori-challenge && cd stori-challenge
mkdir backend && cd backend
go mod init github.com/yourusername/stori-backend
go get github.com/go-chi/chi/v5

# Create directory structure
mkdir -p cmd/server internal/{domain,repository,service,handlers,middleware} data
```

#### 2. Domain Models (30 min)
Create `internal/domain/transaction.go`:
- Transaction struct
- Helper types (Period, CategoryDetail, etc.)
- Basic validation methods

#### 3. JSON Repository (1 hour)
Create `internal/repository/`:
- `repository.go` - Interface definition
- `json_repository.go` - Implementation
- Load embedded JSON data
- Implement GetAll() and GetByDateRange()

**Test checkpoint:** Can load and parse JSON data

#### 4. Analytics Service (1.5 hours)
Create `internal/service/analytics_service.go`:
- Calculate category summaries
- Calculate timeline aggregations
- Percentage calculations
- Date range filtering logic

**Test checkpoint:** Unit test aggregations with sample data

---

### Afternoon Session (3-4 hours)

#### 5. HTTP Handlers (2 hours)
Create `internal/handlers/`:
- `health_handler.go` - Simple health check
- `transaction_handler.go` - GET /api/transactions
- `summary_handler.go` - Category & timeline endpoints
- Wire up service dependencies

#### 6. Middleware (30 min)
Create `internal/middleware/`:
- `cors.go` - CORS handling
- `logger.go` - Request logging
- `recovery.go` - Panic recovery

#### 7. Main Server (1 hour)
Create `cmd/server/main.go`:
- Initialize chi router
- Register middleware
- Register routes
- Embed JSON data
- Environment variable handling
- Graceful shutdown

#### 8. Dockerfile (30 min)
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

**End of Day 1 Deliverable:**
- Working API locally
- All endpoints testable with curl
- Docker container builds and runs
- Ready for frontend integration

**Testing:**
```bash
# Test endpoints
curl http://localhost:8080/api/health
curl http://localhost:8080/api/transactions
curl http://localhost:8080/api/summary/categories
curl http://localhost:8080/api/summary/timeline
```

---

## Day 2: Frontend + AI Integration (6-8 hours)

### Morning Session (3-4 hours)

#### 1. Frontend Setup (30 min)
```bash
cd ..
npm create vite@latest frontend -- --template react
cd frontend
npm install
npm install recharts axios tailwindcss autoprefixer postcss
npx tailwindcss init -p
```

Configure Tailwind in `tailwind.config.js`

#### 2. API Service Layer (30 min)
Create `src/services/api.js`:
- Axios instance with base URL
- API methods for all endpoints
- Error handling

#### 3. Financial Overview Component (1 hour)
Create `src/components/Dashboard/FinancialOverview.jsx`:
- Fetch category summary
- Display total income, expenses, savings
- Card-based layout
- Responsive grid

**Test checkpoint:** See data cards rendering

#### 4. Category Chart (1.5 hours)
Create `src/components/Charts/CategoryChart.jsx`:
- Fetch category summary
- Recharts Pie/Donut chart
- Format expense categories
- Color coding
- Responsive sizing
- Legend

**Test checkpoint:** Chart renders with real data

---

### Afternoon Session (3-4 hours)

#### 5. Timeline Chart (1.5 hours)
Create `src/components/Charts/TimelineChart.jsx`:
- Fetch timeline data
- Recharts Line/Area chart
- Income vs Expenses lines
- Tooltips with formatting
- Responsive sizing

**Test checkpoint:** Timeline shows 10 months of data

#### 6. AI Service Integration (1 hour)
Create `internal/service/ai_service.go` (backend):
- OpenAI API integration
- Prompt construction with real data
- HTTP client for API calls
- Error handling & fallback mock

Create `internal/handlers/advice_handler.go`:
- POST /api/advice endpoint
- Call AI service
- Return formatted advice

#### 7. AI Advisor Component (1.5 hours)
Create `src/components/AI/AIAdvisor.jsx`:
- Button to request advice
- Loading state
- Display advice text
- Insights and recommendations formatting
- Error handling

#### 8. Main App Layout (1 hour)
Create `src/App.jsx`:
- Responsive layout
- Navbar/header
- Dashboard grid
- Mobile-friendly spacing
- Error boundary

**End of Day 2 Deliverable:**
- Complete working application locally
- All features functional
- Mobile responsive
- Docker Compose setup
- Ready for deployment

**Docker Compose:**
```yaml
version: '3.8'
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - OPENAI_API_KEY=${OPENAI_API_KEY}
  
  frontend:
    build: ./frontend
    ports:
      - "5173:5173"
    environment:
      - VITE_API_BASE_URL=http://localhost:8080/api
```

---

## Day 3: AWS Deployment (4-6 hours)

### Frontend Deployment (2 hours)

#### 1. Build & Upload to S3
```bash
cd frontend
npm run build

# Create S3 bucket
aws s3 mb s3://stori-challenge-yourname

# Enable static website hosting
aws s3 website s3://stori-challenge-yourname \
  --index-document index.html

# Upload build
aws s3 sync dist/ s3://stori-challenge-yourname --acl public-read

# Set bucket policy for public access
```

#### 2. CloudFront Setup
- Create CloudFront distribution
- Point to S3 origin
- Configure default root object
- Enable HTTPS
- Wait for deployment (~15 min)

**Test checkpoint:** Frontend accessible via CloudFront URL

---

### Backend Deployment (2-3 hours)

#### Option A: EC2 Deployment (Recommended for reliability)

```bash
# Launch EC2 instance (t3.micro)
# Security group: Allow 22 (SSH), 8080 (API), 443 (HTTPS)

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o server cmd/server/main.go

# SCP binary to EC2
scp -i key.pem server ec2-user@<instance-ip>:/home/ec2-user/

# SSH and setup
ssh -i key.pem ec2-user@<instance-ip>
chmod +x server

# Create systemd service
sudo nano /etc/systemd/system/stori-backend.service
sudo systemctl enable stori-backend
sudo systemctl start stori-backend

# Setup nginx reverse proxy (optional for SSL)
```

#### Option B: ECS Deployment (If time permits)

```bash
# Push Docker image to ECR
aws ecr create-repository --repository-name stori-backend
docker tag stori-backend:latest <ecr-url>/stori-backend:latest
docker push <ecr-url>/stori-backend:latest

# Create ECS cluster, task definition, service
# Configure ALB
```

---

### Integration & Testing (1 hour)

#### 1. Update Frontend API URL
Update `VITE_API_BASE_URL` to point to EC2/ECS endpoint

#### 2. Rebuild & Redeploy Frontend
```bash
npm run build
aws s3 sync dist/ s3://stori-challenge-yourname
aws cloudfront create-invalidation --distribution-id XXX --paths "/*"
```

#### 3. End-to-End Testing
- [ ] Frontend loads via CloudFront
- [ ] API calls work cross-origin
- [ ] Charts render with data
- [ ] AI advice generates
- [ ] Mobile responsive on actual device
- [ ] All error states work

**End of Day 3 Deliverable:**
- Live URL for demo
- Full stack running in AWS
- SSL/HTTPS configured
- Everything tested end-to-end

---

## Day 4: Documentation & Polish (4-6 hours)

### Documentation (3-4 hours)

#### 1. Main README.md (1 hour)
Create `/README.md`:
```markdown
# Stori Financial Tracker

## Overview
[One paragraph description]

## Live Demo
- Frontend: https://xxx.cloudfront.net
- API: https://ec2-xxx.amazonaws.com

## Features
- Financial overview dashboard
- Category spending breakdown
- Income/expense timeline
- AI-powered financial advice

## Tech Stack
[List with rationale]

## Local Development
[Docker Compose instructions]

## Architecture
See docs/architecture.md

## Testing
See docs/testing-strategy.md
```

#### 2. Architecture Document (1.5 hours)
Create `docs/architecture.md`:
- System architecture diagram (mermaid)
- Component interaction flow
- Technology choices with rationale
- Design tradeoffs discussion
- Strengths and weaknesses
- Scalability considerations

**Include diagrams:**
```mermaid
# MVP Architecture
# Data flow
# Deployment architecture
```

#### 3. Testing Strategy (1 hour)
Create `docs/testing-strategy.md`:
- Unit testing approach
- Integration testing approach
- E2E testing approach
- Manual testing checklist
- What would be tested in production

Option: Write actual tests if time permits (more impressive)

#### 4. Future Improvements (30 min)
Add section to architecture.md:
- Database migration path
- CRUD operations implementation
- Multi-account support
- Budget tracking
- Advanced features

---

### Polish (1-2 hours if time permits)

#### Bonus Features (Pick One)
- [ ] Dark mode toggle
- [ ] Loading skeleton states
- [ ] Date range filtering
- [ ] Export CSV functionality
- [ ] GitHub Actions CI/CD
- [ ] Terraform IaC scripts

#### Code Quality
- [ ] Add comments to complex functions
- [ ] Clean up console.logs
- [ ] Verify error handling
- [ ] Check mobile responsiveness
- [ ] Test on different browsers

---

## Final Submission Checklist

### Required Deliverables
- [x] Working Code
  - [x] Public GitHub repository
  - [x] README.md with documentation
  - [x] All code committed
- [x] Evaluation
  - [x] Testing strategy documented OR tests written
- [x] System Design Overview
  - [x] Architecture description (1-2 paragraphs)
  - [x] Architectural diagram
  - [x] Design tradeoffs explained
  - [x] Strengths and weaknesses discussed
  - [x] Future improvements outlined

### Quality Checks
- [ ] Live demo URL works
- [ ] README has clear setup instructions
- [ ] All challenge requirements met
- [ ] Mobile-friendly verified
- [ ] Code is clean and commented
- [ ] No obvious bugs
- [ ] Professional presentation

---

## Emergency Fallback Plan

### If AWS Deployment Fails (Day 3)
- Have Docker Compose working perfectly
- Record demo video of local setup
- Document deployment steps attempted
- Explain issues in README

### If Time Runs Short
**Minimum Viable Submission:**
- Day 1 + Day 2 complete (working locally)
- Basic README
- Architecture diagram (even if simple)
- Testing strategy writeup (don't need actual tests)

The MVP is already a complete submission!

---

## Quick Command Reference

### Backend
```bash
# Run locally
go run cmd/server/main.go

# Test
go test ./...

# Build
go build -o server cmd/server/main.go

# Docker
docker build -t stori-backend .
docker run -p 8080:8080 -e OPENAI_API_KEY=$OPENAI_API_KEY stori-backend
```

### Frontend
```bash
# Run dev
npm run dev

# Build
npm run build

# Preview build
npm run preview

# Docker
docker build -t stori-frontend .
docker run -p 5173:5173 stori-frontend
```

### AWS
```bash
# S3 deploy
aws s3 sync dist/ s3://bucket-name

# CloudFront invalidate
aws cloudfront create-invalidation --distribution-id XXX --paths "/*"

# EC2 deploy
scp -i key.pem server ec2-user@ip:/home/ec2-user/
ssh -i key.pem ec2-user@ip 'sudo systemctl restart stori-backend'
```

---

## Success Criteria

After 4 days, you should have:

1. ✅ Working live demo in AWS
2. ✅ All 3 challenge requirements met:
   - Category spending visualization
   - Income/expense timeline
   - AI financial advice
3. ✅ Mobile-responsive design
4. ✅ Professional documentation
5. ✅ Clean, extensible architecture
6. ✅ Testing strategy
7. ✅ Future improvements articulated

**You're ready to start Day 1! 🚀**

