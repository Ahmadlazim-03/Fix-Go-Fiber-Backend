.PHONY: build run test clean dev help

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/server ./cmd/server

# Run the application
run: build
	@echo "Running application..."
	@./bin/server

# Run in development mode with hot reload (requires air)
dev:
	@echo "Running in development mode..."
	@air

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Run database migrations (if you have migrate tool)
migrate-up:
	@echo "Running database migrations..."
	@migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	@echo "Rolling back database migrations..."
	@migrate -path migrations -database "$(DATABASE_URL)" down

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Generate mocks (requires mockgen)
generate-mocks:
	@echo "Generating mocks..."
	@go generate ./...

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/golangci-lint/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang/mock/mockgen@latest

# Help
help:
	@echo "Available commands:"
	@echo "  build           - Build the application"
	@echo "  run             - Build and run the application"
	@echo "  dev             - Run with hot reload (requires air)"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  clean           - Clean build artifacts"
	@echo "  deps            - Install dependencies"
	@echo "  migrate-up      - Run database migrations"
	@echo "  migrate-down    - Rollback database migrations"
	@echo "  fmt             - Format code"
	@echo "  lint            - Run linter"
	@echo "  generate-mocks  - Generate mocks"
	@echo "  install-tools   - Install development tools"
	@echo "  help            - Show this help"