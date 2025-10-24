.PHONY: help run build test clean docker-build docker-run

# Default target
help:
	@echo "Stori Backend - Available Commands:"
	@echo "  make run           - Run the server locally"
	@echo "  make build         - Build the server binary"
	@echo "  make test          - Run all tests"
	@echo "  make test-coverage - Run tests with coverage"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"
	@echo "  make lint          - Run linter"

# Run the server
run:
	@echo "🚀 Starting server..."
	go run main.go

# Build the server binary
build:
	@echo "🔨 Building server..."
	go build -o server main.go
	@echo "✅ Build complete: ./server"

# Run all tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f server coverage.out coverage.html
	@echo "✅ Clean complete"

# Build Docker image
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t stori-backend:latest .
	@echo "✅ Docker image built: stori-backend:latest"

# Run Docker container
docker-run:
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 \
		-e PORT=8080 \
		-e CORS_ALLOWED_ORIGINS="http://localhost:5173,http://localhost:3000" \
		stori-backend:latest

# Run linter
lint:
	@echo "🔍 Running linter..."
	go vet ./...
	@echo "✅ Linting complete"


