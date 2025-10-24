# Stori Financial Tracker - Backend

Go-based REST API for financial transaction tracking and AI-powered insights.

## 🚀 Quick Start

```bash
# Run the server
make run

# Or directly with Go
go run main.go

# Server starts on http://localhost:8080
```

## 📋 Tech Stack

- **Go 1.22+** with chi router (minimal, idiomatic)
- **Embedded JSON data** (no external database for MVP)
- **RESTful API** with comprehensive error handling
- **Docker** ready with multi-stage builds

## 📁 Project Structure

```
backend/
├── main.go                 # Application entry point
├── internal/
│   ├── domain/            # Business entities & validation
│   ├── repository/        # Data access layer (JSON)
│   ├── service/          # Business logic & calculations
│   ├── handlers/         # HTTP request handlers
│   └── middleware/       # CORS, logging, recovery
├── data/
│   └── transactions.json  # Embedded transaction data (112 records)
├── Dockerfile            # Multi-stage Docker build
└── Makefile             # Development commands
```

## 🎯 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | API info & available endpoints |
| `/api/health` | GET | Health check |
| `/api/transactions` | GET | All transactions (supports date filters) |
| `/api/summary/categories` | GET | Spending breakdown by category |
| `/api/summary/timeline` | GET | Monthly income vs expenses |

## 🔧 Development

### Prerequisites

- Go 1.22 or higher
- Docker (optional)
- Make (optional, but recommended)

### Available Commands

```bash
make run            # Run server locally
make build          # Build binary
make test           # Run all tests
make test-coverage  # Generate coverage report  
make clean          # Clean build artifacts
make docker-build   # Build Docker image
make docker-run     # Run in Docker
make lint           # Run linter
```

### Running Tests

```bash
# All tests
make test

# With coverage
make test-coverage

# Specific package
go test -v ./internal/service/

# Single test
go test -v ./internal/domain/ -run TestTransaction_Validate
```

### Environment Variables

```bash
# Server configuration
PORT=8080                    # Default: 8080

# CORS configuration  
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# Logging
LOG_LEVEL=info              # Default: info
```

## 🐳 Docker

### Build Image

```bash
make docker-build

# Or manually
docker build -t stori-backend:latest .
```

### Run Container

```bash
make docker-run

# Or manually with custom config
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e CORS_ALLOWED_ORIGINS="*" \
  stori-backend:latest
```

## 🧪 Testing the API

### Using curl

```bash
# Health check
curl http://localhost:8080/api/health

# Get all transactions
curl http://localhost:8080/api/transactions

# Get category summary
curl http://localhost:8080/api/summary/categories

# Get timeline
curl http://localhost:8080/api/summary/timeline

# Filter by date range
curl "http://localhost:8080/api/transactions?startDate=2024-01-01&endDate=2024-01-31"
```

### Using Test Script

```bash
chmod +x test-api.sh
./test-api.sh
```

## 🏗️ Architecture

### Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`)
   - Pure business logic
   - No external dependencies
   - Validation rules
   - Error definitions

2. **Repository Layer** (`internal/repository/`)
   - Data access interface
   - JSON implementation (MVP)
   - Ready for database swap

3. **Service Layer** (`internal/service/`)
   - Business logic
   - Financial calculations
   - Data aggregations

4. **Handler Layer** (`internal/handlers/`)
   - HTTP request/response
   - Input validation
   - Error mapping

5. **Middleware** (`internal/middleware/`)
   - CORS
   - Logging
   - Panic recovery

### Design Patterns

- **Repository Pattern**: Abstracts data access
- **Dependency Injection**: Handlers depend on interfaces
- **Middleware Chain**: Composable request processing
- **Embedded Assets**: Binary includes all data

## 📊 Data

The backend uses embedded JSON data with 112 transactions:
- **Period**: January 2024 - October 2024 (10 months)
- **Income**: 20 transactions (bi-weekly salary)
- **Expenses**: 92 transactions across 9 categories
- **Categories**: rent, groceries, utilities, dining, transportation, entertainment, shopping, healthcare

## 🔒 Security Features

- Non-root user in Docker
- Configurable CORS origins
- Request timeout (60s)
- Panic recovery
- Input validation

## 📈 Performance

- **Binary size**: ~15MB (with embedded data)
- **Docker image**: ~20MB (Alpine-based)
- **Startup time**: <100ms
- **Memory usage**: ~10MB idle

## 🚢 Deployment

### Local Binary

```bash
go build -o server main.go
./server
```

### Docker

```bash
docker build -t stori-backend .
docker run -p 8080:8080 stori-backend
```

### AWS EC2 (Production)

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o server main.go

# Deploy
scp server user@instance:/app/
ssh user@instance 'sudo systemctl restart stori-backend'
```

## 🧪 Test Coverage

```
Domain:      95%+ (transaction models, validation)
Repository:  90%+ (data access, filtering)
Service:     95%+ (calculations, aggregations)
Handlers:    90%+ (HTTP endpoints, error handling)
Middleware:  95%+ (CORS, logging, recovery)
```

## 📝 API Response Examples

### Health Check
```json
{
  "status": "healthy",
  "timestamp": "2024-10-24T10:00:00Z"
}
```

### Category Summary
```json
{
  "income": {
    "salary": {"total": 56000, "count": 20, "percentage": 100}
  },
  "expenses": {
    "rent": {"total": 12000, "count": 10, "percentage": 25.5},
    "groceries": {"total": 10240, "count": 23, "percentage": 21.8}
  },
  "summary": {
    "total_income": 56000,
    "total_expenses": 47000,
    "net_savings": 9000,
    "savings_rate": 16.1
  }
}
```

## 🤝 Contributing

This project follows:
- Go best practices
- Clean architecture principles
- RESTful API design
- Comprehensive testing

## 📚 Additional Documentation

See project root `/docs` for:
- System architecture diagrams
- API specifications
- Architecture decision records (ADRs)
- Testing strategy

