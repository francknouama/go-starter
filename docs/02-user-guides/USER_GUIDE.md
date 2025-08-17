# go-starter User Guide

Complete guide to using go-starter for generating production-ready Go projects with modern best practices.

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Basic Usage](#basic-usage)
- [Advanced Features](#advanced-features)
- [Project Types Deep Dive](#project-types-deep-dive)
- [Configuration Management](#configuration-management)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)
- [Migration and Upgrades](#migration-and-upgrades)

## Overview

go-starter is a comprehensive Go project generator that combines the simplicity of create-react-app with the flexibility of Spring Initializr. It generates production-ready Go projects with:

### ✨ Key Features

- **12 Project Types**: CLI, Web API, Lambda, Microservice, Library, Monolith, Workspace
- **Multiple Architectures**: Standard, Clean Architecture, DDD, Hexagonal, Event-driven
- **Progressive Disclosure**: Adapts interface to user experience level
- **Simplified Logger System**: 60-90% code reduction across all project types
- **Production Hardening**: Advanced AST operations, automated test generation
- **Cross-platform**: Windows, macOS, Linux support

### 🎯 Philosophy

- **Start Simple, Grow Organically**: Begin with minimal structure, add complexity as needed
- **Progressive Learning**: Clear path from beginner to expert
- **Best Practices by Default**: Generate code following Go community standards
- **Flexibility**: Extensive customization options for specific needs

## Installation

### Prerequisites

- **Go 1.21+**: Required for building and running go-starter
- **Git**: For version control and module management
- **Make**: Optional, for using generated Makefiles

### Installation Methods

#### 1. Go Install (Recommended)
```bash
go install github.com/francknouama/go-starter@latest
```

**Pros**: Always gets latest version, works on all platforms
**Cons**: Requires Go to be installed

#### 2. Homebrew (macOS/Linux)
```bash
brew tap francknouama/go-starter
brew install go-starter
```

**Pros**: Easy updates, no Go required
**Cons**: macOS/Linux only

#### 3. Package Managers

##### Windows
```powershell
# Chocolatey
choco install go-starter

# Scoop
scoop bucket add go-starter https://github.com/francknouama/scoop-go-starter
scoop install go-starter

# WinGet (coming soon)
winget install go-starter
```

##### Linux
```bash
# Snap
sudo snap install go-starter

# APT (Ubuntu/Debian)
sudo add-apt-repository ppa:francknouama/go-starter
sudo apt update
sudo apt install go-starter

# YUM/DNF (RedHat/Fedora)
sudo dnf install go-starter
```

#### 4. Direct Download
Download binaries from [GitHub Releases](https://github.com/francknouama/go-starter/releases).

#### 5. Quick Install Script
```bash
curl -sSL https://raw.githubusercontent.com/francknouama/go-starter/main/scripts/install.sh | bash
```

### Verification

```bash
go-starter version
go-starter --help
```

## Basic Usage

### Command Structure

```
go-starter <command> [arguments] [flags]
```

### Core Commands

#### 1. `new` - Generate New Project

```bash
# Interactive mode (beginner-friendly)
go-starter new

# Direct generation
go-starter new <project-name> --type=<type> [options]

# Advanced mode with all options
go-starter new --advanced
```

#### 2. `list` - Show Available Options

```bash
# List all project types
go-starter list

# List with descriptions
go-starter list --verbose

# List specific category
go-starter list --category=web
```

#### 3. `version` - Show Version Information

```bash
go-starter version
```

### Essential Flags

#### Basic Mode Flags (14 total)
- `--name`: Project name
- `--type`: Project type (cli, web-api, lambda, etc.)
- `--module`: Go module path
- `--framework`: Framework choice (gin, echo, cobra, etc.)
- `--logger`: Logger type (slog, zap, logrus, zerolog)
- `--go-version`: Go version (default: 1.21)
- `--output`: Output directory
- `--complexity`: Complexity level (simple, standard, advanced, expert)
- `--dry-run`: Preview without creating files
- `--quiet`: Minimal output
- `--no-git`: Skip git initialization
- `--random-name`: Generate random project name
- `--help`: Show help
- `--advanced`: Enable advanced mode

#### Advanced Mode Additional Flags (18+ total)
- `--architecture`: Architecture pattern (clean, ddd, hexagonal)
- `--database-driver`: Database driver (postgres, mysql, mongodb, sqlite)
- `--database-orm`: ORM choice (gorm, sqlx, ent)
- `--auth-type`: Authentication type (jwt, oauth2, session)
- `--no-banner`: Disable ASCII banner
- `--banner-style`: Banner style choice

### Progressive Disclosure System

go-starter adapts its interface based on user experience:

#### Basic Mode (Default)
- Shows 14 essential flags
- Beginner-friendly descriptions
- Smart defaults for most options
- Guided experience

```bash
# See basic help
go-starter new --help
```

#### Advanced Mode
- Shows all 18+ flags
- Expert-level options
- Full customization control
- Power user features

```bash
# See advanced help  
go-starter new --advanced --help

# Use advanced generation
go-starter new my-project --type=web-api --advanced
```

## Advanced Features

### Complexity Levels

go-starter uses complexity levels to generate appropriate project structures:

#### Simple (Beginner)
- **File Count**: 8-15 files
- **Use Case**: Learning, prototypes, simple utilities
- **Features**: Minimal structure, basic functionality
- **Example**: 
  ```bash
  go-starter new my-tool --type=cli --complexity=simple
  ```

#### Standard (Intermediate)
- **File Count**: 25-35 files
- **Use Case**: Production applications, team development
- **Features**: Complete structure, testing, documentation
- **Example**: 
  ```bash
  go-starter new my-api --type=web-api --complexity=standard
  ```

#### Advanced (Expert)
- **File Count**: 40-60 files
- **Use Case**: Enterprise applications, complex business logic
- **Features**: Advanced patterns, monitoring, security
- **Example**: 
  ```bash
  go-starter new my-service --type=microservice --complexity=advanced
  ```

#### Expert (Architect)
- **File Count**: 60+ files
- **Use Case**: Large organizations, complex domains
- **Features**: Full enterprise patterns, observability, governance
- **Example**: 
  ```bash
  go-starter new enterprise-platform --type=workspace --complexity=expert
  ```

### Dry Run Mode

Preview project structure before generation:

```bash
go-starter new my-project --type=web-api --dry-run
```

**Output shows**:
- Complete file list
- Directory structure
- Configuration summary
- No files created

### Random Name Generation

Generate creative project names:

```bash
go-starter new --random-name --type=cli
# Creates project like "brave-gopher-cli" or "cosmic-falcon-tool"
```

### Configuration Files

#### Project-Level Configuration

Generated projects include configuration files:

```yaml
# config/config.yaml
server:
  port: 8080
  host: localhost
  
database:
  driver: postgres
  host: localhost
  port: 5432
  
logger:
  level: info
  format: json
```

#### User-Level Configuration

Create `~/.go-starter.yaml` for personal defaults:

```yaml
author: "Your Name"
email: "your.email@example.com"
license: "MIT"
defaults:
  goVersion: "1.21"
  framework: "gin"
  logger: "slog"
  complexity: "standard"
  outputDir: "./projects"
```

## Project Types Deep Dive

### CLI Applications

Command-line tools and utilities.

#### When to Choose CLI
- Building developer tools
- Creating automation scripts
- System administration utilities
- Command-line interfaces for APIs

#### CLI Blueprints

##### 1. Simple CLI
```bash
go-starter new my-tool --type=cli --complexity=simple
```

**Structure** (8 files):
```
my-tool/
├── main.go          # Entry point
├── cmd/
│   ├── root.go      # Root command
│   └── version.go   # Version command
├── config.go        # Configuration
├── Makefile         # Build automation
├── README.md        # Documentation
├── go.mod          # Dependencies
└── go.sum          # Dependency hashes
```

**Features**:
- Basic flag handling
- Simple configuration
- Minimal dependencies
- Quick setup

##### 2. Standard CLI
```bash
go-starter new my-tool --type=cli --complexity=standard
```

**Structure** (29 files):
```
my-tool/
├── main.go
├── cmd/
│   ├── root.go
│   ├── version.go
│   ├── completion.go
│   ├── create.go
│   ├── delete.go
│   ├── list.go
│   └── update.go
├── internal/
│   ├── config/
│   ├── logger/
│   ├── errors/
│   └── output/
├── configs/
│   └── config.yaml
├── Dockerfile
├── Makefile
└── README.md
```

**Features**:
- Multiple subcommands
- Configuration file support
- Structured logging
- Error handling
- Testing framework
- Docker support

#### CLI Configuration Options

```bash
# Logger selection
go-starter new my-tool --type=cli --logger=zap

# Framework selection (currently only Cobra)
go-starter new my-tool --type=cli --framework=cobra

# Go version
go-starter new my-tool --type=cli --go-version=1.21
```

### Web APIs

RESTful APIs with multiple architecture patterns.

#### When to Choose Web API
- Building REST APIs
- Creating microservices
- Developing backend services
- API-first applications

#### Architecture Patterns

##### 1. Standard Architecture
```bash
go-starter new my-api --type=web-api --architecture=standard
```

**Best for**: Most web APIs, rapid development, small to medium teams

**Structure**:
```
my-api/
├── cmd/server/main.go
├── internal/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   └── services/
├── configs/
├── migrations/
└── tests/
```

##### 2. Clean Architecture
```bash
go-starter new my-api --type=web-api --architecture=clean
```

**Best for**: Enterprise applications, complex business logic, testable systems

**Structure**:
```
my-api/
├── cmd/server/main.go
├── internal/
│   ├── adapters/     # External adapters
│   ├── domain/       # Business entities
│   ├── infrastructure/  # Technical details
│   └── usecases/     # Business logic
```

**Benefits**:
- Clear separation of concerns
- Highly testable
- Independent of frameworks
- Dependency inversion

##### 3. Domain-Driven Design (DDD)
```bash
go-starter new my-api --type=web-api --architecture=ddd
```

**Best for**: Complex domains, rich business models, event-driven systems

**Structure**:
```
my-api/
├── internal/
│   ├── domain/
│   │   ├── user/        # User aggregate
│   │   │   ├── entity.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── specifications.go
│   │   └── order/       # Order aggregate
│   ├── application/     # Application services
│   ├── infrastructure/  # Technical implementation
│   └── presentation/    # HTTP handlers
```

**Features**:
- Aggregate design
- Domain services
- Repository patterns
- Specifications
- Domain events

##### 4. Hexagonal Architecture
```bash
go-starter new my-api --type=web-api --architecture=hexagonal
```

**Best for**: Highly testable systems, multiple adapters, ports & adapters pattern

**Structure**:
```
my-api/
├── internal/
│   ├── adapters/
│   │   ├── primary/     # Driving adapters (HTTP, CLI)
│   │   └── secondary/   # Driven adapters (DB, Email)
│   ├── application/
│   │   ├── ports/       # Interface definitions
│   │   └── services/    # Business logic
│   └── domain/         # Core domain
```

**Benefits**:
- Testable architecture
- Multiple interface support
- Flexible adapter pattern
- Clean boundaries

#### Web API Configuration

##### Framework Selection
```bash
# Gin (default, fastest)
go-starter new my-api --type=web-api --framework=gin

# Echo (middleware-rich)
go-starter new my-api --type=web-api --framework=echo

# Fiber (Express-like)
go-starter new my-api --type=web-api --framework=fiber

# Chi (lightweight)
go-starter new my-api --type=web-api --framework=chi
```

##### Database Integration
```bash
# PostgreSQL with GORM
go-starter new my-api --type=web-api \
  --database-driver=postgres \
  --database-orm=gorm

# MySQL with SQLX
go-starter new my-api --type=web-api \
  --database-driver=mysql \
  --database-orm=sqlx

# MongoDB (no ORM)
go-starter new my-api --type=web-api \
  --database-driver=mongodb

# SQLite for development
go-starter new my-api --type=web-api \
  --database-driver=sqlite \
  --database-orm=gorm
```

##### Authentication
```bash
# JWT authentication
go-starter new my-api --type=web-api --auth-type=jwt

# OAuth2
go-starter new my-api --type=web-api --auth-type=oauth2

# Session-based
go-starter new my-api --type=web-api --auth-type=session

# API Key
go-starter new my-api --type=web-api --auth-type=api-key
```

### AWS Lambda

Serverless functions optimized for AWS Lambda.

#### When to Choose Lambda
- Event-driven processing
- Webhooks and API endpoints
- Scheduled tasks
- Serverless architectures

#### Lambda Types

##### 1. Standard Lambda
```bash
go-starter new my-lambda --type=lambda
```

**Best for**: Event handlers, background processing, simple functions

**Structure**:
```
my-lambda/
├── main.go          # Lambda entry point
├── handler.go       # Business logic
├── template.yaml    # SAM template
├── Makefile        # Build and deploy
└── internal/
    ├── logger/
    └── observability/
```

##### 2. Lambda API Proxy
```bash
go-starter new my-api --type=lambda-proxy
```

**Best for**: RESTful APIs on Lambda, API Gateway integration

**Structure**:
```
my-api/
├── main.go
├── handler.go
├── internal/
│   ├── handlers/    # HTTP handlers
│   ├── middleware/  # Lambda middleware
│   ├── models/      # Request/response models
│   └── observability/
├── template.yaml    # SAM template with API Gateway
└── terraform/       # Infrastructure as code
```

#### Lambda Configuration

```bash
# Runtime optimization
go-starter new my-lambda --type=lambda --complexity=advanced

# With observability
go-starter new my-lambda --type=lambda --logger=zap

# API Gateway integration
go-starter new my-api --type=lambda-proxy --framework=gin
```

### Libraries

Reusable Go packages for distribution.

#### When to Choose Library
- Creating reusable packages
- SDK development
- Open-source projects
- Shared utilities

```bash
go-starter new my-lib --type=library
```

**Structure**:
```
my-lib/
├── library.go       # Main package code
├── library_test.go  # Tests
├── examples/        # Usage examples
│   ├── basic/
│   └── advanced/
├── doc.go          # Package documentation
├── CHANGELOG.md    # Version history
└── LICENSE         # License file
```

**Features**:
- Public API design
- Comprehensive examples
- Documentation generation
- Versioning support

### Microservices

Production-ready microservices with gRPC and observability.

#### When to Choose Microservice
- Distributed systems
- Service mesh architectures
- Cloud-native applications
- Scalable backends

```bash
go-starter new my-service --type=microservice
```

**Structure**:
```
my-service/
├── cmd/server/main.go
├── internal/
│   ├── app/         # Application setup
│   ├── handler/     # gRPC handlers
│   ├── middleware/  # gRPC middleware
│   ├── metrics/     # Prometheus metrics
│   └── tracing/     # Distributed tracing
├── proto/           # Protocol buffers
├── deployments/
│   ├── k8s/         # Kubernetes manifests
│   └── istio/       # Service mesh config
└── configs/
    ├── local.yaml
    ├── staging.yaml
    └── production.yaml
```

**Features**:
- gRPC server and client
- Health checks
- Metrics collection
- Distributed tracing
- Circuit breakers
- Rate limiting
- Kubernetes deployment
- Service mesh integration

### Monoliths

Traditional web applications with all components in one deployable unit.

#### When to Choose Monolith
- Rapid prototyping
- Small teams
- Traditional web applications
- Simplified deployment

```bash
go-starter new my-app --type=monolith
```

**Structure**:
```
my-app/
├── main.go
├── controllers/     # HTTP controllers
├── models/         # Data models
├── services/       # Business logic
├── middleware/     # HTTP middleware
├── routes/         # Route definitions
├── views/          # HTML templates
├── static/         # CSS, JS, images
├── database/       # Migrations, seeders
└── config/         # Configuration
```

**Features**:
- MVC architecture
- HTML template engine
- Static asset serving
- Database migrations
- Session management
- Authentication
- Admin interface

### Workspaces

Multi-module projects for monorepos.

#### When to Choose Workspace
- Monorepos
- Multiple related services
- Shared libraries
- Complex project organization

```bash
go-starter new my-workspace --type=workspace
```

**Structure**:
```
my-workspace/
├── go.work         # Workspace definition
├── cmd/            # Multiple commands
│   ├── api/
│   ├── worker/
│   └── cli/
├── pkg/            # Shared packages
│   ├── auth/
│   ├── database/
│   └── utils/
├── services/       # Individual services
│   ├── user-service/
│   └── order-service/
└── tools/          # Development tools
    └── codegen/
```

**Features**:
- Go workspace configuration
- Shared module management
- Cross-service dependencies
- Unified build system
- Development tools

## Configuration Management

### Environment Configuration

Projects support multiple environments with configuration files:

```
configs/
├── config.dev.yaml      # Development
├── config.test.yaml     # Testing
├── config.staging.yaml  # Staging
└── config.prod.yaml     # Production
```

### Configuration Structure

```yaml
# config.yaml
app:
  name: "my-api"
  version: "1.0.0"
  environment: "development"

server:
  host: "localhost"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"

database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  name: "myapp"
  user: "myapp"
  password: "secret"
  ssl_mode: "disable"
  max_open_conns: 25
  max_idle_conns: 5

logger:
  level: "info"
  format: "json"
  output: "stdout"

auth:
  jwt_secret: "your-secret-key"
  token_ttl: "24h"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

### Environment Variables

All configuration values can be overridden with environment variables:

```bash
# Server configuration
export APP_SERVER_HOST=0.0.0.0
export APP_SERVER_PORT=3000

# Database configuration
export APP_DATABASE_HOST=db.example.com
export APP_DATABASE_PASSWORD=supersecret

# Logger configuration
export APP_LOGGER_LEVEL=debug
```

### Configuration Loading

Generated projects use a hierarchical configuration system:

1. **Default values** (in code)
2. **Configuration files** (YAML)
3. **Environment variables** (highest priority)

## Best Practices

### Project Organization

#### Naming Conventions
- Use kebab-case for project names: `my-awesome-api`
- Use valid Go module paths: `github.com/username/project-name`
- Follow Go package naming: lowercase, no underscores

#### Directory Structure
- Keep `internal/` for private packages
- Use `pkg/` for packages intended for external use
- Put binaries in `cmd/` directory
- Store configuration in `configs/`

#### Module Management
```bash
# Initialize with proper module path
go-starter new my-project --module=github.com/myorg/my-project

# Use semantic versioning
git tag v1.0.0
```

### Development Workflow

#### 1. Generation
```bash
# Start with dry run
go-starter new my-project --type=web-api --dry-run

# Generate with specific requirements
go-starter new my-project --type=web-api --database-driver=postgres

# Review generated structure
cd my-project && tree
```

#### 2. Customization
```bash
# Install dependencies
go mod download

# Run tests
make test

# Start development server
make dev
```

#### 3. Development
```bash
# Format code
make fmt

# Lint code
make lint

# Build for production
make build
```

### Performance Optimization

#### Logger Performance

Different loggers have different performance characteristics:

```bash
# Highest performance (zero allocation)
go-starter new my-api --logger=zerolog

# High performance (structured)
go-starter new my-api --logger=zap

# Balanced performance and features
go-starter new my-api --logger=slog

# Feature-rich (moderate performance)
go-starter new my-api --logger=logrus
```

#### Database Optimization

```bash
# High-performance database access
go-starter new my-api --database-orm=sqlx

# ORM with good performance
go-starter new my-api --database-orm=gorm

# Type-safe SQL generation
go-starter new my-api --database-orm=sqlc
```

### Security Best Practices

#### Generated Security Features
- Input validation middleware
- Security headers
- Rate limiting
- Authentication middleware
- SQL injection prevention
- XSS protection

#### Security Configuration
```bash
# JWT authentication with security headers
go-starter new my-api --type=web-api --auth-type=jwt

# HTTPS-ready configuration
go-starter new my-api --type=web-api --complexity=advanced
```

### Testing Strategy

#### Generated Test Structure
```
tests/
├── unit/           # Unit tests
├── integration/    # Integration tests
├── fixtures/       # Test data
└── helpers/        # Test utilities
```

#### Test Commands
```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run integration tests
make test-integration

# Run benchmarks
make benchmark
```

### Deployment

#### Container Deployment
All projects include Docker support:

```bash
# Build Docker image
make docker-build

# Run with Docker Compose
docker-compose up -d
```

#### Kubernetes Deployment
Microservices include Kubernetes manifests:

```bash
# Deploy to Kubernetes
kubectl apply -f deployments/k8s/
```

#### Cloud Deployment
Lambda projects include deployment scripts:

```bash
# Deploy to AWS
make deploy
```

## Troubleshooting

### Common Issues

#### 1. Build Failures

**Problem**: `go build` fails with dependency errors
```
go: module example.com/myproject: reading go.mod: open go.mod: no such file or directory
```

**Solution**: 
```bash
# Ensure you're in the project directory
cd my-project

# Verify go.mod exists
ls -la go.mod

# Download dependencies
go mod download
```

#### 2. Import Path Issues

**Problem**: Generated code has incorrect import paths

**Solution**: 
```bash
# Regenerate with correct module path
go-starter new my-project --module=github.com/myorg/my-project

# Or update go.mod manually
go mod edit -module=github.com/myorg/my-project
```

#### 3. Template Generation Errors

**Problem**: `template: parse error` during generation

**Solution**:
```bash
# Update to latest version
go install github.com/francknouama/go-starter@latest

# Verify installation
go-starter version

# Try generation again
go-starter new my-project --type=web-api
```

#### 4. Permission Issues

**Problem**: Permission denied when creating files

**Solution**:
```bash
# Check directory permissions
ls -la .

# Create in different directory
go-starter new my-project --output=/tmp/projects

# Fix permissions
chmod 755 .
```

#### 5. Database Connection Issues

**Problem**: Generated project can't connect to database

**Solution**:
```bash
# Check database configuration
cat configs/config.dev.yaml

# Update connection string
export APP_DATABASE_HOST=localhost
export APP_DATABASE_PORT=5432

# Verify database is running
docker-compose up -d db
```

### Debug Mode

Enable verbose output for troubleshooting:

```bash
# Verbose generation
go-starter new my-project --type=web-api --verbose

# Debug output
export GO_STARTER_DEBUG=1
go-starter new my-project --type=web-api
```

### Log Analysis

Check logs for generation issues:

```bash
# View generation logs
cat ~/.go-starter/logs/generation.log

# View error logs
cat ~/.go-starter/logs/errors.log
```

## Migration and Upgrades

### Upgrading go-starter

#### Go Install Method
```bash
go install github.com/francknouama/go-starter@latest
```

#### Homebrew Method
```bash
brew upgrade go-starter
```

#### Package Manager Method
```bash
# Chocolatey
choco upgrade go-starter

# APT
sudo apt update && sudo apt upgrade go-starter
```

### Migrating Existing Projects

#### From Simple to Standard CLI
1. Generate new standard CLI project
2. Compare structures
3. Gradually migrate code
4. Update dependencies

#### From Standard to Clean Architecture
1. Generate clean architecture project
2. Identify domain entities
3. Move business logic to domain layer
4. Implement repository interfaces
5. Update tests

#### Version Migration
When upgrading between major versions:

1. **Review changelog**
2. **Test in development**
3. **Update dependencies**
4. **Run test suite**
5. **Deploy gradually**

### Configuration Migration

#### v1.3 to v1.4 Migration

Configuration changes:
```yaml
# Old format (v1.3)
database:
  url: "postgres://user:pass@host:port/db"

# New format (v1.4)
database:
  driver: "postgres"
  host: "host"
  port: 5432
  name: "db"
  user: "user"
  password: "pass"
```

Migration script:
```bash
# Run migration helper
go-starter migrate --from=v1.3 --to=v1.4
```

## Advanced Topics

### Custom Templates

#### Template Customization
```bash
# Export templates for customization
go-starter export-templates --output=./custom-templates

# Use custom templates
go-starter new my-project --templates=./custom-templates
```

#### Template Development
Create custom project types:

1. **Define template structure**
2. **Create template.yaml**
3. **Add template files**
4. **Test generation**
5. **Share with community**

### Plugin System

#### Available Plugins
- **Database plugins**: Additional database drivers
- **Authentication plugins**: Extended auth methods
- **Deployment plugins**: Cloud-specific deployment
- **Testing plugins**: Additional testing frameworks

#### Plugin Installation
```bash
# Install plugin
go-starter plugin install github.com/example/go-starter-plugin

# List installed plugins
go-starter plugin list

# Use plugin
go-starter new my-project --plugin=example
```

### Enterprise Features

#### Organization Templates
```bash
# Set organization defaults
go-starter config set organization "My Company"
go-starter config set template-repo "github.com/myorg/go-starter-templates"

# Generate with organization templates
go-starter new my-project --org-template
```

#### Compliance and Governance
- **License compliance**: Automatic license file generation
- **Security scanning**: Integration with security tools
- **Audit logging**: Track project generation
- **Policy enforcement**: Organizational standards

---

This user guide provides comprehensive coverage of go-starter's features and capabilities. For specific questions or advanced use cases, refer to the [GitHub repository](https://github.com/francknouama/go-starter) or [community discussions](https://github.com/francknouama/go-starter/discussions).