.PHONY: build test lint install clean setup help run dev-build validate

# Default target
help: ## Show this help message
	@echo "go-starter development commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# Build commands
build: build-cli ## Build the CLI binary (alias for build-cli)

build-cli: ## Build the CLI binary
	@echo "Building go-starter CLI..."
	go build -o bin/go-starter ./cmd/go-starter
	@echo "✓ Built: bin/go-starter"

build-web: ## Build the production web server
	@echo "Building go-starter web server (production)..."
	cd web && go build -o ../bin/go-starter-web ./cmd/web-server
	@echo "✓ Built: bin/go-starter-web"

build-dev: ## Build the development web server
	@echo "Building go-starter web server (development)..."
	go build -o bin/go-starter-dev ./cmd/go-starter-dev
	@echo "✓ Built: bin/go-starter-dev"

build-all: ## Build all binaries
	@echo "Building all go-starter binaries..."
	$(MAKE) build-cli
	$(MAKE) build-web
	$(MAKE) build-dev
	@echo "✓ Built all binaries"

dev-build: ## Build CLI with race detection for development
	@echo "Building go-starter CLI with race detection..."
	go build -race -o bin/go-starter-race ./cmd/go-starter
	@echo "✓ Built: bin/go-starter-race"

install: ## Install go-starter CLI to $GOPATH/bin
	@echo "Installing go-starter CLI..."
	go install ./cmd/go-starter
	@echo "✓ Installed go-starter"

# Legacy build (with deprecation warning)
build-legacy: ## Build using legacy root main.go (deprecated)
	@echo "⚠️  WARNING: Building from root main.go is deprecated"
	@echo "   Use 'make build-cli' instead"
	go build -o bin/go-starter-legacy .
	@echo "✓ Built: bin/go-starter-legacy"

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

setup-dev: ## Full developer environment setup using setup script
	@echo "Running full developer setup..."
	./scripts/setup.sh

# Validation
validate: build test lint ## Run all validation checks
	@echo "✓ All validation checks passed"

validate-all: ## Validate all template combinations
	@echo "Validating all template combinations..."
	./scripts/test_all_combinations.sh

validate-templates: ## Quick template validation for CI/CD
	@echo "Validating template generation..."
	./scripts/validate_template_generation.sh

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning up..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean
	@echo "✓ Cleanup completed"

# Run commands (for testing)
run: build-cli ## Build and run CLI with sample arguments
	@echo "Running go-starter with --help..."
	./bin/go-starter --help

run-web: build-web ## Build and run production web server
	@echo "Running go-starter web server (production)..."
	./bin/go-starter-web

run-dev: build-dev ## Build and run development web server
	@echo "Running go-starter web server (development)..."
	./bin/go-starter-dev

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

release: ## Create a new release
	@echo "Creating new release..."
	./scripts/release.sh