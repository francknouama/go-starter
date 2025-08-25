# Getting Started with go-starter

Welcome to go-starter! This comprehensive guide will take you from installation to your first production-ready Go project using our progressive disclosure system.

## Table of Contents

- [Installation](#installation)
- [Understanding Progressive Disclosure](#understanding-progressive-disclosure)
- [Your First Project](#your-first-project)
- [Project Types Deep Dive](#project-types-deep-dive)
- [Configuration Options](#configuration-options)
- [Development Workflow](#development-workflow)
- [Next Steps](#next-steps)

## Installation

### Prerequisites

- **Go 1.21+** (recommended, but optional - go-starter auto-detects your version)
- **Git** (optional but recommended)
- **Make** (optional, for using Makefile commands)

### Install go-starter

#### Option 1: Go Install (Recommended)
```bash
go install github.com/francknouama/go-starter@latest
```

#### Option 2: Homebrew (macOS/Linux)
```bash
brew tap francknouama/go-starter
brew install go-starter
```

#### Option 3: Direct Download
Download the latest binary from [GitHub Releases](https://github.com/francknouama/go-starter/releases).

### Verify Installation
```bash
go-starter version
# Expected output:
# Version: 2.0.0
# Commit: 1f9d312
# Built: 2025-07-22
```

## Understanding Progressive Disclosure

go-starter features a unique **progressive disclosure system** that adapts to your experience level, preventing beginners from being overwhelmed while giving experts full control.

### 🎯 Two Interface Modes

#### Basic Mode (Default - Beginner Friendly)
Shows only **14 essential options**:
```bash
go-starter new --help
# Shows: --type, --name, --module, --framework, --logger, etc.
# Hides: --database-driver, --auth-type, --architecture, etc.
```

#### Advanced Mode (Power Users)
Shows all **18+ options** including enterprise features:
```bash
go-starter new --advanced --help
# Shows: All database options, authentication types, deployment configs
```

### 📊 Complexity Levels

| Level | Description | Files | Best For |
|-------|-------------|-------|----------|
| **Simple** | Minimal structure | 8-15 | Learning, prototypes, scripts |
| **Standard** | Production-ready | 25-35 | Most projects, teams |
| **Advanced** | Enterprise patterns | 40-60 | Complex business logic |
| **Expert** | Full-featured | 60+ | Large organizations |

## Your First Project

### Beginner Workflow: Interactive Mode

Start with the guided experience:
```bash
go-starter new
```

You'll see **simplified prompts** with only essential options:
```
? Choose your project type:
  ❯ 🌐 Web API (Standard) - REST APIs, CRUD services  
    🏗️ Clean Architecture API - Enterprise applications
    🖥️ Simple CLI - Scripts, utilities (8 files)
    ⚙️ Standard CLI - Production tools (29 files)
    📦 Library - SDKs, packages
    ⚡ AWS Lambda - Event functions

? Select logger:
  ❯ slog - Go standard library (recommended)
    zap - High performance, zero allocation  
    logrus - Feature-rich, popular choice
    zerolog - Zero allocation JSON

? Module path: github.com/yourusername/my-project
```

### Expert Workflow: Direct Generation

Skip prompts when you know exactly what you want:
```bash
# Enterprise API with all the bells and whistles
go-starter new enterprise-api \
  --type=web-api \
  --architecture=clean \
  --database-driver=postgres \
  --database-orm=gorm \
  --auth-type=jwt \
  --logger=zap \
  --framework=gin \
  --advanced

# High-performance microservice  
go-starter new user-service \
  --type=microservice \
  --logger=zerolog \
  --advanced

# Simple automation script
go-starter new my-script \
  --type=cli \
  --complexity=simple \
  --logger=slog
```

## 🎉 7 Production-Ready Blueprints - Milestone Achievement!

**Major Update**: We've reached our **7 production-ready blueprints milestone**! This includes the newly production-ready **gRPC Gateway** blueprint for dual HTTP/gRPC APIs.

### ✅ Production-Ready Blueprints (Ready for Immediate Use)
- **CLI-Simple** (10 files) - Learning Go, quick utilities
- **CLI-Standard** (28 files) - Production CLI tools  
- **Web-API-Standard** (44 files) - REST APIs, CRUD services
- **Lambda-Standard** (17 files) - AWS serverless functions
- **Library-Standard** (19 files) - Go packages, SDKs
- **Web-API-Clean** (69 files) - Enterprise Clean Architecture
- **gRPC-Gateway** (45 files) - Dual HTTP/gRPC APIs ✅ **NEWLY PRODUCTION READY**

## Project Types Deep Dive

Each production-ready blueprint has been thoroughly tested and validated for immediate production use:

### 🖥️ CLI Applications

#### Simple CLI (8 files)
**Perfect for**: Scripts, utilities, learning Go
```bash
go-starter new my-script --type=cli --complexity=simple
```

**Generated structure**:
```
my-script/
├── main.go          # Entry point with basic logic
├── cmd/
│   ├── root.go      # Root command definition
│   └── version.go   # Version command
├── config.go        # Simple configuration
├── Makefile         # Build automation
├── README.md        # Usage documentation
└── go.mod          # Go module definition
```

#### Standard CLI (29 files)
**Perfect for**: Production tools, DevOps utilities, complex CLIs
```bash
go-starter new my-tool --type=cli --complexity=standard
```

**Features**: Subcommands, config files, interactive prompts, shell completion, structured logging, comprehensive testing.

### 🌐 Web APIs (4 Architecture Patterns)

#### Standard Web API
**Use Case**: REST APIs, microservices, standard backend services
```bash
go-starter new my-api --type=web-api
```

#### Clean Architecture Web API
**Use Case**: Enterprise applications, complex business logic
```bash
go-starter new enterprise-api --type=web-api --architecture=clean
```

#### Domain-Driven Design (DDD) Web API
**Use Case**: Domain-rich applications, complex business rules
```bash
go-starter new domain-api --type=web-api --architecture=ddd
```

#### Hexagonal Architecture Web API
**Use Case**: Maximum testability, multiple adapters
```bash
go-starter new testable-api --type=web-api --architecture=hexagonal
```

### ☁️ Serverless Functions

#### AWS Lambda
**Use Case**: Event processing, webhooks, background tasks
```bash
go-starter new my-function --type=lambda
```

#### Lambda API Proxy
**Use Case**: REST APIs on Lambda, API Gateway integration
```bash
go-starter new serverless-api --type=lambda-proxy
```

### 🏢 Enterprise & Specialized

#### Library
**Use Case**: SDKs, reusable packages, open-source projects
```bash
go-starter new awesome-lib --type=library
```

#### gRPC Gateway ✅ **NEWLY PRODUCTION READY**
**Use Case**: API gateway patterns, dual-protocol APIs, protocol translation
```bash
go-starter new my-gateway --type=grpc-gateway
```
**Production Features**: Dual HTTP/gRPC support, enhanced interceptors, unified middleware, rate limiting, metrics collection

#### Microservice
**Use Case**: Distributed systems, service mesh, cloud-native
```bash
go-starter new user-service --type=microservice
```

#### Monolith
**Use Case**: Traditional web apps, full-stack systems
```bash
go-starter new webapp --type=monolith
```

#### Go Workspace
**Use Case**: Monorepos, multi-module projects
```bash
go-starter new platform --type=workspace
```

## Configuration Options

### Database Integration

**Supported Databases**:
- **PostgreSQL** - Production RDBMS (recommended)
- **MySQL** - Popular alternative
- **MongoDB** - Document database
- **SQLite** - Development/testing
- **Redis** - Caching/sessions

```bash
# Single database
go-starter new my-api --type=web-api --database-driver=postgres

# Multiple databases
go-starter new my-api --type=web-api --database-driver=postgres,redis
```

### Logger Selection

go-starter features a **simplified logger system** with unified interface:

| Logger | Performance | Use Case | Package |
|--------|-------------|----------|---------|
| **slog** | Good | Standard library choice | `log/slog` |
| **zap** | Excellent | High-performance apps | `go.uber.org/zap` |
| **logrus** | Good | Feature-rich requirements | `github.com/sirupsen/logrus` |
| **zerolog** | Excellent | JSON-heavy, cloud-native | `github.com/rs/zerolog` |

```bash
# High-performance API
go-starter new fast-api --type=web-api --logger=zap

# Cloud-native service
go-starter new cloud-service --type=web-api --logger=zerolog

# Standard development
go-starter new my-app --logger=slog
```

### Authentication

```bash
# JWT authentication (most common)
go-starter new my-api --type=web-api --auth-type=jwt

# OAuth2 providers
go-starter new my-api --type=web-api --auth-type=oauth2

# Session-based (traditional)
go-starter new my-api --type=web-api --auth-type=session
```

### Framework Selection

#### Web Frameworks
```bash
--framework=gin     # Fastest, most popular (default)
--framework=echo    # Middleware-rich
--framework=fiber   # Express-like API
--framework=chi     # Lightweight router
```

## Development Workflow

### After Generation

```bash
cd my-project
make help           # See all available commands
```

### Essential Make Commands

Every project includes a comprehensive Makefile:
```bash
make run           # Start development server
make test          # Run tests with coverage
make lint          # Run golangci-lint
make build         # Build production binary
make docker        # Build Docker image
make clean         # Clean build artifacts
make dev           # Development mode with hot reload
```

### Configuration Management

Each project includes environment-based configuration:

```yaml
# config/config.yaml
server:
  port: 8080
  timeout: 30s

database:
  host: localhost
  port: 5432
  name: myapp

logging:
  level: info
  format: json
```

### Working with Generated Code

Your chosen logger is available throughout the application with a **unified interface**:

```go
// All loggers use the same simple interface
logger.Info("Server started", "port", 8080, "version", "1.0.0")
logger.Error("Database connection failed", "error", err)
logger.Debug("Processing request", "user_id", userID)
```

### Local Development with Docker

```bash
# Start dependencies (databases, cache, etc.)
docker-compose up -d

# Run your application
make run

# View logs
docker-compose logs -f
```

## Real-World Examples

### Startup MVP
```bash
go-starter new startup-api \
  --type=web-api \
  --framework=gin \
  --database-driver=postgres \
  --auth-type=jwt \
  --logger=slog
```

### Enterprise System
```bash
go-starter new enterprise-system \
  --type=web-api \
  --architecture=clean \
  --database-driver=postgres \
  --database-orm=gorm \
  --auth-type=jwt \
  --logger=zap \
  --advanced
```

### Developer Tool
```bash
go-starter new dev-tool \
  --type=cli \
  --complexity=standard \
  --logger=logrus
```

### Event Processing
```bash
go-starter new event-processor \
  --type=lambda \
  --logger=zerolog
```

### Microservices Platform
```bash
# Create workspace for multiple services
go-starter new platform --type=workspace

# Add individual services
cd platform/services
go-starter new user-service --type=microservice --logger=zap
go-starter new order-service --type=microservice --logger=zap
go-starter new notification-service --type=microservice --logger=zerolog
```

## Configuration Profiles

Create persistent defaults with `~/.go-starter.yaml`:

```yaml
profiles:
  work:
    author: "Your Name"
    email: "you@company.com"
    defaults:
      logger: zap
      goVersion: "1.22"
      framework: gin
  personal:
    author: "Your Name"  
    email: "personal@email.com"
    defaults:
      logger: slog
      complexity: simple
      
current_profile: work
```

Switch profiles:
```bash
go-starter config set current_profile personal
```

## Pro Tips

### 1. Use Dry Run for Planning
```bash
go-starter new my-project --type=web-api --architecture=clean --dry-run
```
Shows exactly what files will be generated without creating them.

### 2. Start Simple, Upgrade Later
```bash
# Start simple
go-starter new my-tool --type=cli --complexity=simple

# Later, generate standard version for comparison
go-starter new my-tool-v2 --type=cli --complexity=standard
```

### 3. Preview with Different Loggers
```bash
# Compare logger implementations
go-starter new api-slog --type=web-api --logger=slog --dry-run
go-starter new api-zap --type=web-api --logger=zap --dry-run
```

### 4. Leverage Progressive Disclosure
- **New to Go?** Use basic mode and interactive prompts
- **Experienced?** Use advanced mode with direct flags  
- **Team lead?** Use dry-run to plan project structure

## Troubleshooting

### Common Issues

#### Permission Denied
```bash
chmod +x $(which go-starter)
```

#### Module Path Issues
```bash
# Use valid module path format
go-starter new my-project --module=github.com/username/my-project
```

#### Build Failures
```bash
# Check Go version compatibility
go version
# Should be 1.21 or higher
```

#### Template Generation Errors
```bash
# Verify installation
go-starter version

# Reinstall if needed
go install github.com/francknouama/go-starter@latest
```

## Next Steps

Now that you've mastered the basics:

1. **Experiment with Different Types**: Try CLI, web API, and Lambda projects
2. **Explore Architecture Patterns**: Compare standard vs clean vs DDD architectures
3. **Test Logger Performance**: Benchmark different loggers for your use case
4. **Read Project-Specific Docs**: Check the README.md in your generated projects
5. **Join the Community**: Star the [GitHub repo](https://github.com/francknouama/go-starter) and contribute

## Getting Help

- **CLI Help**: `go-starter --help` or `go-starter new --advanced --help`
- **Documentation**: Full guides in [user guides](../02-user-guides/)
- **Issues**: [Report bugs](https://github.com/francknouama/go-starter/issues)
- **Discussions**: [Community forum](https://github.com/francknouama/go-starter/discussions)
- **Troubleshooting**: [Common issues guide](../02-user-guides/troubleshooting.md)

---

**Happy coding!** 🚀 You're now ready to build amazing Go applications with go-starter.

**Next**: Explore the [blueprint selection guide](../02-user-guides/blueprint-selection.md) to choose the perfect project type for your needs.