.PHONY: build test lint install clean setup help run dev-build validate

# Default target
help: ## Show this help message
	@echo "go-starter development commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Build commands
build: ## Build the CLI binary
	@echo "Building go-starter..."
	go build -o bin/go-starter .
	@echo "✓ Built: bin/go-starter"

dev-build: ## Build with race detection for development
	@echo "Building go-starter with race detection..."
	go build -race -o bin/go-starter-dev .
	@echo "✓ Built: bin/go-starter-dev"

install: ## Install go-starter to $GOPATH/bin
	@echo "Installing go-starter..."
	go install .
	@echo "✓ Installed go-starter"

# Test commands
test: ## Run all tests with coverage
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Tests completed. Coverage report: coverage.html"

test-short: ## Run tests without coverage
	@echo "Running tests (short)..."
	go test -v ./...

# Code quality
lint: ## Run golangci-lint
	@echo "Running linter..."
	golangci-lint run --config .golangci.yml
	@echo "✓ Linting completed"

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...
	@echo "✓ Code formatted"

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...
	@echo "✓ Vet completed"

# Development setup
setup: ## Set up development environment
	@echo "Setting up development environment..."
	@echo "Installing development dependencies..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "Downloading Go modules..."
	go mod download
	@echo "✓ Development environment ready"

# Validation
validate: build test lint ## Run all validation checks
	@echo "✓ All validation checks passed"

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean
	@echo "✓ Cleanup completed"

# Run commands (for testing)
run: build ## Build and run with sample arguments
	@echo "Running go-starter with --help..."
	./bin/go-starter --help

# Development utilities
mod-tidy: ## Run go mod tidy
	@echo "Tidying modules..."
	go mod tidy
	@echo "✓ Modules tidied"

deps: ## Show dependency graph
	@echo "Dependency graph:"
	go mod graph

check-updates: ## Check for dependency updates
	@echo "Checking for dependency updates..."
	go list -u -m all

# Release preparation (for future use)
release-dry: ## Dry run release (requires goreleaser)
	@echo "Dry run release..."
	@if command -v goreleaser &> /dev/null; then \
		goreleaser release --snapshot --rm-dist; \
	else \
		echo "goreleaser not installed. Install with: go install github.com/goreleaser/goreleaser@latest"; \
	fi