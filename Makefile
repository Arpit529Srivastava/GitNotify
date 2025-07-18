# GitNotify Makefile

.PHONY: build test clean run help

# Default target
all: build

# Build the application
build:
	@echo "Building GitNotify..."
	go build -o gitnotify cmd/gitnotify/main.go
	@echo "Build complete!"

# Run tests
test:
	@echo "Running tests..."
	go test ./...
	@echo "Tests complete!"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...
	@echo "Test coverage complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f gitnotify
	@echo "Clean complete!"

# Run the application (requires config.yml)
run: build
	@echo "Starting GitNotify..."
	./gitnotify

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	@echo "Dependencies installed!"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatted!"

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run
	@echo "Linting complete!"

# Show help
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Build and run the application"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code (requires golangci-lint)"
	@echo "  help         - Show this help message" 