# CADUBI v2 Makefile

.PHONY: build test clean run fmt lint help

# Default target
all: build

# Build the application
build:
	@echo "Building CADUBI v2..."
	go build -o bin/cadubi cmd/cadubi/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Run the application
run: build
	@echo "Running CADUBI v2..."
	./bin/cadubi

# Run the application directly without building
dev:
	@echo "Running CADUBI v2 in development mode..."
	go run cmd/cadubi/main.go

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	go vet ./...

# Run JSON-driven tests specifically
test-json:
	@echo "Running JSON-driven command tests..."
	go test -v ./pkg/commands/ -run TestJSON

# Run headless tests
test-headless:
	@echo "Running headless tests..."
	go test -v ./cmd/cadubi/ -run TestHeadless

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Create release builds for multiple platforms
release:
	@echo "Building release binaries..."
	@mkdir -p bin/release
	GOOS=linux GOARCH=amd64 go build -o bin/release/cadubi-linux-amd64 cmd/cadubi/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/release/cadubi-darwin-amd64 cmd/cadubi/main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/release/cadubi-darwin-arm64 cmd/cadubi/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/release/cadubi-windows-amd64.exe cmd/cadubi/main.go
	@echo "Release binaries created in bin/release/"

# Help target
help:
	@echo "Available targets:"
	@echo "  build         Build the application"
	@echo "  test          Run all tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  test-json     Run JSON-driven command tests"
	@echo "  test-headless Run headless tests"
	@echo "  clean         Clean build artifacts"
	@echo "  run           Build and run the application"
	@echo "  dev           Run application in development mode"
	@echo "  fmt           Format code"
	@echo "  lint          Lint code"
	@echo "  deps          Install dependencies"
	@echo "  release       Build release binaries for multiple platforms"
	@echo "  help          Show this help message"