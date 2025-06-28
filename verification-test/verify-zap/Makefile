.PHONY: build run test lint clean dev docker-build docker-run help

# Variables
BINARY_NAME=verify-zap
MAIN_PATH=./cmd/server
BUILD_DIR=./bin
DOCKER_IMAGE=verify-zap:latest

# Default target
.DEFAULT_GOAL := help

## Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✓ Built: $(BUILD_DIR)/$(BINARY_NAME)"

## Run the application
run: build
	@echo "Starting $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

## Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Tests completed. Coverage report: coverage.html"

## Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run
	@echo "✓ Linting completed"

## Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✓ Cleaned"

## Run in development mode with hot reload
dev:
	@echo "Starting development server..."
	@go run $(MAIN_PATH)

## Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"

## Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "✓ Dependencies tidied"

## Update dependencies
update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✓ Dependencies updated"

## Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "✓ Docker image built: $(DOCKER_IMAGE)"

## Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(DOCKER_IMAGE)

## Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✓ Development tools installed"

## Security check
security:
	@echo "Running security checks..."
	@govulncheck ./...
	@echo "✓ Security check completed"

## Generate API documentation
docs:
	@echo "Generating API documentation..."
	@swag init -g $(MAIN_PATH)/main.go -o ./docs
	@echo "✓ API documentation generated"

## Help
help:
	@echo "verify-zap - Available commands:"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""