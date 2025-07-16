.PHONY: build run test clean docker-build docker-run help

# Variables
BINARY_NAME=mcp-server
PORT=8082
DATABASE_URL=postgres://smartsec:smartsec123@localhost:5432/smartsec?sslmode=disable

# Build the application
build:
	go build -o $(BINARY_NAME) .

# Run the application
run:
	@echo "Starting MCP Server on port $(PORT)..."
	@export PORT=$(PORT) && \
	export DATABASE_URL="$(DATABASE_URL)" && \
	export LOG_LEVEL=debug && \
	./$(BINARY_NAME)

# Run with go run
dev:
	@echo "Starting MCP Server in development mode..."
	@export PORT=$(PORT) && \
	export DATABASE_URL="$(DATABASE_URL)" && \
	export LOG_LEVEL=debug && \
	go run main.go

# Test the application
test:
	go test -v ./...

# Test with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy

# Build Docker image
docker-build:
	docker build -t smartsec-mcp-server .

# Run Docker container
docker-run:
	docker run -p $(PORT):$(PORT) \
		-e DATABASE_URL="$(DATABASE_URL)" \
		-e PORT=$(PORT) \
		-e LOG_LEVEL=info \
		smartsec-mcp-server

# Run with docker-compose
docker-compose-up:
	docker-compose up --build

# Stop docker-compose
docker-compose-down:
	docker-compose down

# Generate API documentation
docs:
	@echo "API Documentation available at: http://localhost:$(PORT)/"
	@echo "Schema endpoint: http://localhost:$(PORT)/mcp/schema"
	@echo "Examples endpoint: http://localhost:$(PORT)/mcp/examples"

# Check health
health:
	curl -f http://localhost:$(PORT)/health || echo "Service is not healthy"

# Example queries
examples:
	@echo "Example: Get all devices"
	curl -X GET http://localhost:$(PORT)/mcp/entities/devices
	@echo "\n\nExample: Query Linux devices"
	curl -X POST http://localhost:$(PORT)/mcp/query \
		-H "Content-Type: application/json" \
		-d '{"entity":"devices","filters":[{"field":"os","operator":"eq","value":"Linux"}],"limit":10}'

# Install dependencies
deps:
	go mod download

# Check if required tools are installed
check-tools:
	@command -v go >/dev/null 2>&1 || { echo "Go is required but not installed. Aborting." >&2; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo "Docker is required but not installed. Aborting." >&2; exit 1; }

# Help
help:
	@echo "Available targets:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  dev             - Run in development mode"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  clean           - Clean build artifacts"
	@echo "  fmt             - Format code"
	@echo "  lint            - Lint code"
	@echo "  tidy            - Tidy dependencies"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  docker-compose-up   - Run with docker-compose"
	@echo "  docker-compose-down - Stop docker-compose"
	@echo "  docs            - Show API documentation URLs"
	@echo "  health          - Check service health"
	@echo "  examples        - Show example API calls"
	@echo "  deps            - Install dependencies"
	@echo "  check-tools     - Check if required tools are installed"
	@echo "  help            - Show this help message"

# Default target
all: check-tools deps build
