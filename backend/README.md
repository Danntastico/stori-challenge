# Stori Financial Tracker - Backend

Go-based REST API for financial transaction tracking and AI-powered advice.

## Tech Stack

- **Go 1.22+**
- **chi router** - Minimal, idiomatic routing
- **Standard library** - JSON, HTTP, embedding
- **OpenAI API** - Financial advice generation

## Project Structure

```
backend/
├── cmd/server/          # Application entry point
├── internal/
│   ├── domain/          # Business entities and models
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic
│   ├── handlers/        # HTTP handlers
│   └── middleware/      # HTTP middleware
├── data/                # Embedded JSON data
└── go.mod
```

## Getting Started

### Prerequisites

- Go 1.22 or higher
- OpenAI API key (optional - will use mock if not provided)

### Installation

```bash
# Install dependencies
go mod download

# Copy environment variables
cp .env.example .env

# Add your OpenAI API key to .env (optional)
```

### Running Locally

```bash
# Run the server
go run cmd/server/main.go

# Server will start on http://localhost:8080
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Test specific package
go test ./internal/service/
```

### Building

```bash
# Build binary
go build -o server cmd/server/main.go

# Run binary
./server
```

### Docker

```bash
# Build image
docker build -t stori-backend .

# Run container
docker run -p 8080:8080 \
  -e OPENAI_API_KEY=your-key \
  stori-backend
```

## API Endpoints

- `GET /api/health` - Health check
- `GET /api/transactions` - Get all transactions (with optional filters)
- `GET /api/summary/categories` - Category spending breakdown
- `GET /api/summary/timeline` - Monthly income/expense timeline
- `POST /api/advice` - Get AI-powered financial advice

## Architecture

This backend follows clean architecture principles:

- **Domain Layer**: Core business models (Transaction, Category, etc.)
- **Repository Layer**: Data access abstraction (currently JSON, designed for DB swap)
- **Service Layer**: Business logic (aggregations, AI integration)
- **Handler Layer**: HTTP request/response handling
- **Middleware**: Cross-cutting concerns (CORS, logging, recovery)

## Development

See `/docs` in the project root for:
- Full API specification
- Architecture decisions
- Testing strategy

