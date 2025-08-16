# Go Web Server Makefile

# Variables
BINARY_NAME=go-starter-web
MAIN_PATH=./cmd/web-server
BUILD_DIR=./bin
VERSION=1.0.0
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Default target
.PHONY: all
all: clean build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for production with optimizations
.PHONY: build-prod
build-prod:
	@echo "Building $(BINARY_NAME) for production..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -a -installsuffix cgo $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Production build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the development server
.PHONY: run
run: build
	@echo "Starting web server..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run in development mode with auto-reload
.PHONY: dev
dev:
	@echo "Starting development server with hot reload..."
	@command -v air >/dev/null 2>&1 || { echo "Installing air..."; go install github.com/cosmtrek/air@latest; }
	air -c .air.toml

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Lint the code
.PHONY: lint
lint:
	@echo "Running linter..."
	@command -v golangci-lint >/dev/null 2>&1 || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	golangci-lint run

# Format the code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed"

# Generate dependencies
.PHONY: generate
generate:
	@echo "Running go generate..."
	go generate ./...
	@echo "✅ Generation complete"

# Build frontend assets
.PHONY: build-frontend
build-frontend:
	@echo "Building frontend assets..."
	npm run build
	@echo "✅ Frontend build complete"

# Install frontend dependencies
.PHONY: install-frontend
install-frontend:
	@echo "Installing frontend dependencies..."
	npm install
	@echo "✅ Frontend dependencies installed"

# Full development setup
.PHONY: setup
setup: install-frontend deps
	@echo "Development setup complete"

# Full build (frontend + backend)
.PHONY: build-all
build-all: build-frontend build
	@echo "✅ Full build complete"

# Production build (frontend + backend)
.PHONY: build-all-prod
build-all-prod: build-frontend build-prod
	@echo "✅ Production build complete"

# Docker build
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .
	@echo "✅ Docker image built: $(BINARY_NAME):$(VERSION)"

# Docker run
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(BINARY_NAME):$(VERSION)

# Check if required tools are installed
.PHONY: check-tools
check-tools:
	@echo "Checking required tools..."
	@command -v go >/dev/null 2>&1 || { echo "❌ Go is not installed"; exit 1; }
	@command -v node >/dev/null 2>&1 || { echo "❌ Node.js is not installed"; exit 1; }
	@command -v npm >/dev/null 2>&1 || { echo "❌ npm is not installed"; exit 1; }
	@echo "✅ All required tools are installed"

# Show help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build           - Build the web server binary"
	@echo "  build-prod      - Build optimized production binary"
	@echo "  run             - Build and run the web server"
	@echo "  dev             - Run development server with hot reload"
	@echo "  clean           - Clean build artifacts"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  lint            - Run code linter"
	@echo "  fmt             - Format code"
	@echo "  deps            - Install Go dependencies"
	@echo "  build-frontend  - Build frontend assets"
	@echo "  install-frontend- Install frontend dependencies"
	@echo "  setup           - Full development setup"
	@echo "  build-all       - Build frontend and backend"
	@echo "  build-all-prod  - Production build of frontend and backend"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  check-tools     - Check if required tools are installed"
	@echo "  help            - Show this help message"

# Set default target
.DEFAULT_GOAL := help