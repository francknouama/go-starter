# verify-cli Makefile

.PHONY: build install clean test lint run help

# Build variables
BINARY_NAME=verify-cli
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse HEAD)
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

## help: Show this help message
help:
	@echo "Available commands:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: Build the binary
build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

## install: Install the binary
install:
	go install $(LDFLAGS) .

## run: Run the application
run:
	go run . --help

## test: Run tests
test:
	go test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

## lint: Run linter
lint:
	golangci-lint run

## clean: Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

## deps: Install dependencies
deps:
	go mod download
	go mod tidy

## docker-build: Build Docker image
docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .

## docker-run: Run Docker container
docker-run:
	docker run --rm -it $(BINARY_NAME):$(VERSION)

.DEFAULT_GOAL := help