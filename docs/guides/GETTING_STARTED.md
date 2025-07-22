# Getting Started with go-starter

Welcome to **go-starter**! This guide will help you get up and running with the most comprehensive Go project generator available. With **12 blueprints**, **progressive disclosure**, and **enterprise architecture patterns**, you'll be creating production-ready Go applications that scale from simple scripts to complex enterprise systems.

## Table of Contents

- [Installation](#installation)
- [Progressive Disclosure System](#progressive-disclosure-system)
- [Your First Project](#your-first-project)
- [Understanding Project Complexity](#understanding-project-complexity)
- [All 12 Blueprint Types](#all-12-blueprint-types)
- [Understanding the Logger Selector](#understanding-the-logger-selector)
- [Working with Generated Projects](#working-with-generated-projects)
- [Advanced Usage](#advanced-usage)
- [Tips and Best Practices](#tips-and-best-practices)

## Installation

### Prerequisites

- **Go 1.21 or later** (optional but recommended - go-starter automatically detects your Go version)
- **Git** (optional but recommended)
- **Make** (optional, for using Makefile commands)

### Quick Install

#### Option 1: Homebrew (Recommended for macOS/Linux)

```bash
brew tap francknouama/tap
brew install go-starter
```

#### Option 2: Go Install

```bash
go install github.com/francknouama/go-starter@latest
```

#### Option 3: Download Binary

Visit the [releases page](https://github.com/francknouama/go-starter/releases/latest) and download the appropriate binary for your system.

### Verify Installation

```bash
go-starter version
```

You should see output like:
```
Version: 2.0.0
Commit: 1f9d312
Built: 2025-07-22
```

## Progressive Disclosure System

go-starter features a **progressive disclosure system** that adapts to your experience level and project complexity needs. No more overwhelming beginners with advanced options or limiting experts with simplified interfaces.

### 🎯 Two Disclosure Modes

#### Basic Mode (Default - Beginner Friendly)
Shows only **14 essential options** with clear descriptions:

```bash
go-starter new --help
# Shows: --type, --name, --module, --logger, --framework, etc.
# Hides: --database-driver, --auth-type, --architecture, etc.
```

#### Advanced Mode (Power Users)
Shows all **18+ options** including enterprise features:

```bash
go-starter new --advanced --help
# Shows: All database options, authentication types, deployment configs
```

### 🎓 Smart Learning Path

The system automatically suggests the right complexity level:

```bash
# Beginner: Simple CLI (8 files)
go-starter new my-script --type cli --complexity simple

# Intermediate: Standard CLI (29 files) 
go-starter new my-tool --type cli --complexity standard

# Advanced: Enterprise architectures
go-starter new my-api --type web-api --architecture hexagonal --advanced
```

### 💡 Key Benefits

- **No Overwhelm**: Beginners see only what they need
- **No Limits**: Experts access all advanced features  
- **Smart Defaults**: Reasonable choices for quick generation
- **Clear Progression**: Obvious path from simple to complex

## Your First Project

### Progressive Disclosure Workflow

The system adapts to your experience level automatically:

#### Beginner Workflow (Basic Mode)
```bash
go-starter new my-first-project
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
    [Show more options...]

? Select logger:
  ❯ slog - Go standard library (recommended)
    zap - High performance, zero allocation  
    logrus - Feature-rich, popular choice
    zerolog - Zero allocation JSON

? Module path: github.com/yourusername/my-first-project
```

#### Advanced Workflow (Power Users)
```bash
go-starter new enterprise-system --advanced
```

Access **all 18+ options** for complex projects:

```
? Choose architecture pattern:
  ❯ Standard - Simple layered architecture
    Clean - Clean Architecture principles
    DDD - Domain-Driven Design  
    Hexagonal - Ports & Adapters
    Event-Driven - CQRS + Event Sourcing

? Select databases (multiple allowed):
  ◯ PostgreSQL - Production RDBMS
  ◯ MySQL - Popular RDBMS  
  ◯ MongoDB - Document database
  ◯ Redis - In-memory cache
  ◯ SQLite - File-based DB

? Authentication type:
  ❯ None
    JWT - JSON Web Tokens
    OAuth2 - OAuth2 providers
    Session - Server-side sessions

? Deployment targets:
  ◯ Docker - Container deployment
  ◯ Kubernetes - Container orchestration
  ◯ AWS Lambda - Serverless
  ◯ Terraform - Infrastructure as Code
```

### Direct Mode (For Power Users)

Skip prompts when you know exactly what you want:

```bash
# Progressive complexity examples
go-starter new my-script --type cli --complexity simple         # 8 files
go-starter new my-tool --type cli --complexity standard         # 29 files  
go-starter new fast-api --type web-api --architecture clean --logger zap
go-starter new my-service --type microservice --logger zerolog
go-starter new my-workspace --type workspace                    # Monorepo

# Enterprise patterns with advanced mode
go-starter new enterprise-api --type web-api --architecture hexagonal --advanced
go-starter new event-system --type event-driven --logger zap --advanced
```

## Understanding Project Complexity

go-starter uses a **4-tier complexity system** to help you choose the right blueprint:

### 🟢 Beginner (Simple)
**Perfect for**: Learning, scripts, prototypes  
**Files**: 8-15 files  
**Features**: Basic structure, minimal dependencies  

```bash
go-starter new my-script --type cli --complexity simple
# Generates: 8 files, single main.go, basic config
```

### 🟡 Intermediate (Standard)  
**Perfect for**: Production tools, standard applications  
**Files**: 20-30 files  
**Features**: Full structure, testing, CI/CD  

```bash
go-starter new my-api --type web-api  
# Generates: 29 files, layered architecture, Docker, tests
```

### 🟠 Advanced (Enterprise Patterns)
**Perfect for**: Complex systems, team projects  
**Files**: 40-60 files  
**Features**: Advanced architectures, multiple patterns  

```bash
go-starter new my-system --type web-api --architecture clean
# Generates: Clean Architecture, dependency injection, mocks
```

### 🔴 Expert (Full-Featured)
**Perfect for**: Large-scale systems, microservices  
**Files**: 60+ files  
**Features**: All patterns, monitoring, deployment  

```bash
go-starter new my-platform --type microservice --advanced
# Generates: Service mesh, K8s, monitoring, distributed tracing
```

## All 12 Blueprint Types

go-starter provides **12 production-ready blueprints** covering every Go project type:

### 📊 Core Web APIs (4 Architecture Patterns)

#### 1. 🌐 Standard Web API
**Use Case**: REST APIs, CRUD services, microservices  
**Architecture**: Standard layered (handlers → services → repository)  
**Files**: ~29 files

```bash
go-starter new my-api --type web-api
```

#### 2. 🏗️ Clean Architecture Web API  
**Use Case**: Enterprise applications, complex business logic  
**Architecture**: Clean Architecture (entities → use cases → adapters)  
**Files**: ~45 files

```bash
go-starter new enterprise-api --type web-api --architecture clean
```

#### 3. ⚙️ DDD Web API
**Use Case**: Domain-rich applications, business rule heavy systems  
**Architecture**: Domain-Driven Design (domain → application → infrastructure)  
**Files**: ~42 files

```bash
go-starter new domain-api --type web-api --architecture ddd  
```

#### 4. 🔩 Hexagonal Architecture Web API
**Use Case**: Highly testable systems, multiple adapters  
**Architecture**: Ports & Adapters (core → ports → adapters)  
**Files**: ~48 files  

```bash
go-starter new testable-api --type web-api --architecture hexagonal
```

### 🖥️ CLI Applications (2 Complexity Levels)

#### 5. 📱 Simple CLI
**Use Case**: Scripts, utilities, quick tools  
**Complexity**: Beginner (8 files)  
**Features**: Single command, basic flags, minimal structure  

```bash
go-starter new my-script --type cli --complexity simple
```

#### 6. ⚙️ Standard CLI  
**Use Case**: Production CLI tools, DevOps utilities  
**Complexity**: Professional (29 files)  
**Features**: Subcommands, config files, advanced patterns  

```bash
go-starter new my-tool --type cli --complexity standard
```

### 🏢 Enterprise & Cloud-Native (4 Blueprints)

#### 7. 🌐 gRPC Gateway
**Use Case**: API Gateway with dual HTTP/gRPC support  
**Features**: Protocol buffers, TLS, dual endpoints  

```bash
go-starter new api-gateway --type grpc-gateway
```

#### 8. 🔄 Event-Driven Architecture
**Use Case**: CQRS, Event Sourcing, distributed systems  
**Features**: Event streams, command/query separation, projections  

```bash
go-starter new event-system --type event-driven
```

#### 9. 🏗️ Microservice  
**Use Case**: Service mesh, Kubernetes, distributed systems  
**Features**: Service discovery, circuit breakers, health checks  

```bash
go-starter new user-service --type microservice
```

#### 10. 🏢 Monolith
**Use Case**: Traditional web applications, full-stack systems  
**Features**: HTML templating, sessions, asset pipeline  

```bash
go-starter new webapp --type monolith
```

### ☁️ Serverless & Tools (2 Blueprints)

#### 11. ⚡ AWS Lambda
**Use Case**: Event-driven functions, webhooks  
**Runtime**: AWS Lambda Go  

```bash
go-starter new my-function --type lambda
```

#### 12. 🌉 Lambda Proxy  
**Use Case**: API Gateway integration, RESTful serverless  
**Features**: HTTP proxy patterns, API Gateway optimized  

```bash
go-starter new serverless-api --type lambda-proxy
```

#### 13. 📦 Library
**Use Case**: SDKs, reusable packages, open source  
**Features**: Clean public API, examples, documentation  

```bash
go-starter new awesome-lib --type library
```

#### 14. 🔧 Go Workspace
**Use Case**: Monorepos, multi-module projects  
**Features**: Multiple services, shared packages  

```bash
go-starter new my-platform --type workspace
```

### 🎯 Choosing the Right Blueprint

**Starting a new project?** Use this decision tree:

1. **Web API needed?** → Choose architecture complexity:
   - Simple CRUD → `web-api` (standard)
   - Enterprise system → `web-api --architecture clean`
   - Domain-heavy → `web-api --architecture ddd`  
   - Maximum testability → `web-api --architecture hexagonal`

2. **CLI tool needed?** → Choose complexity:
   - Quick script → `cli --complexity simple`
   - Production tool → `cli --complexity standard`

3. **Distributed system?** → Choose pattern:
   - Event-driven → `event-driven`
   - Microservices → `microservice`
   - API Gateway → `grpc-gateway`

4. **Serverless?** → Choose runtime:
   - Simple functions → `lambda`
   - REST API → `lambda-proxy`

5. **Package/Library?** → `library`

6. **Multiple services?** → `workspace`

## Understanding the Logger Selector

One of go-starter's most powerful features is the **Logger Selector System**. Here's how to choose the right logger for your project:

### Logger Comparison

| Logger | When to Use | Performance | Key Features |
|--------|-------------|-------------|--------------|
| **slog** | Default choice, standard Go projects | Good | Built-in to Go, structured logging, no dependencies |
| **zap** | High-traffic APIs, microservices | Excellent | Zero allocation, fastest option, JSON output |
| **logrus** | Feature-rich applications, legacy systems | Moderate | Hooks, formatters, extensive ecosystem |
| **zerolog** | JSON-heavy services, cloud-native apps | Excellent | Zero allocation, chainable API, clean JSON |

### Real-World Examples

#### Example 1: High-Performance API Service

```bash
go-starter new payment-api --type web-api --logger zap
```

Why zap? For a payment API handling thousands of requests per second, zap's zero-allocation design ensures logging doesn't impact performance.

#### Example 2: Developer-Friendly CLI Tool

```bash
go-starter new devtool --type cli --logger logrus
```

Why logrus? CLI tools benefit from logrus's readable output formats and extensive formatting options for terminal display.

#### Example 3: Cloud-Native Microservice

```bash
go-starter new user-service --type web-api --logger zerolog
```

Why zerolog? Cloud platforms often require JSON logs, and zerolog provides the cleanest JSON output with excellent performance.

## Project Types

### 1. Web API

Perfect for REST APIs, web services, and backend applications.

**Features:**
- Gin web framework (default)
- Health check endpoints
- Structured project layout
- Database support (optional)
- Docker & Docker Compose
- Graceful shutdown

**Generated Structure:**
```
my-api/
├── cmd/api/main.go        # Entry point
├── internal/
│   ├── config/           # Configuration
│   ├── handlers/         # HTTP handlers
│   ├── middleware/       # Middleware
│   └── logger/          # Your chosen logger
├── Makefile             # Development tasks
└── docker-compose.yml   # Local development
```

### 2. CLI Application

For building command-line tools and utilities.

**Features:**
- Cobra framework
- Subcommands support
- Configuration management
- Auto-completion scripts

**Generated Structure:**
```
my-cli/
├── cmd/
│   ├── root.go         # Root command
│   └── version.go      # Version command
├── internal/
│   ├── config/        # Configuration
│   └── logger/        # Your chosen logger
└── main.go
```

### 3. Library

For creating reusable Go packages.

**Features:**
- Clean public API
- Example implementations
- Comprehensive documentation structure
- No main function

**Generated Structure:**
```
my-lib/
├── mylib.go           # Public API
├── mylib_test.go      # Tests
├── examples/
│   └── basic/
│       └── main.go    # Usage example
└── internal/          # Private implementation
```

### 4. AWS Lambda

For serverless functions on AWS.

**Features:**
- AWS Lambda handler
- CloudWatch-optimized logging
- SAM template
- Local testing setup

**Generated Structure:**
```
my-lambda/
├── cmd/lambda/main.go  # Lambda handler
├── internal/
│   ├── handler/       # Business logic
│   └── logger/        # Your chosen logger
├── template.yaml      # SAM template
└── Makefile          # Build & deploy
```

## Working with Generated Projects

### First Steps

After generating your project:

```bash
cd my-first-api
make help  # See available commands
```

### Common Commands

All project types include a Makefile with helpful commands:

```bash
make run        # Start the application
make test       # Run tests with coverage
make lint       # Run golangci-lint
make build      # Build production binary
make docker     # Build Docker image
make clean      # Clean build artifacts
```

### Project Configuration

Each project includes a configuration system. For example, in a web API:

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

### Working with the Logger

Your chosen logger is available throughout the application:

```go
// Using slog (default)
logger.Info("Server started", 
    slog.String("port", "8080"),
    slog.String("version", "1.0.0"))

// Using zap
logger.Info("Server started",
    zap.String("port", "8080"),
    zap.String("version", "1.0.0"))

// Using logrus
logger.WithFields(logrus.Fields{
    "port": "8080",
    "version": "1.0.0",
}).Info("Server started")

// Using zerolog
logger.Info().
    Str("port", "8080").
    Str("version", "1.0.0").
    Msg("Server started")
```

## Advanced Usage

### Multiple Databases

go-starter supports multiple databases in a single project:

**Available Databases:**
- **PostgreSQL** - Relational database, ACID compliant
- **MySQL** - Popular relational database
- **MongoDB** - Document database for flexible schemas
- **SQLite** - File-based database, great for development
- **Redis** - In-memory data store for caching

Select multiple databases during project creation:

```bash
# Interactive mode
go-starter new multi-db-api --type web-api
# Then select multiple databases with spacebar

# Direct mode
go-starter new multi-db-api --type web-api --database postgres,redis,mongodb
```

This generates:
- Multi-service Docker Compose with all selected databases
- Connection code for each database with proper error handling
- Health checks for all services
- Environment-based configuration for each database
- GORM ORM integration for SQL databases

### Force Overwrite

Regenerate a project with different settings:

```bash
go-starter new my-api --type web-api --logger zap --force
```

### Custom Module Paths

Specify your module path directly:

```bash
go-starter new my-api --type web-api --module github.com/company/my-api
```

### Configuration Profiles

Create a configuration file at `~/.go-starter.yaml`:

```yaml
profiles:
  work:
    author: "Your Name"
    email: "you@company.com"
    defaults:
      logger: zap
      goVersion: "1.22"
  personal:
    author: "Your Name"
    email: "personal@email.com"
    defaults:
      logger: slog
current_profile: work
```

Switch profiles:
```bash
go-starter config set current_profile personal
```

## Tips and Best Practices

### 1. Choose the Right Logger

- **Starting out?** Use `slog` - it's the standard library choice
- **Need performance?** Use `zap` or `zerolog`
- **Need features?** Use `logrus` for its extensive ecosystem

### 2. Start Simple

Begin with the standard architecture and add complexity as needed:

```bash
# Start simple
go-starter new my-api --type web-api

# Later, regenerate with more features
go-starter new my-api --type web-api --database postgres,redis --force
```

### 3. Use Make Commands

The generated Makefile contains many helpful commands:

```bash
make help        # Always start here
make run         # Development mode with hot reload
make test-watch  # Run tests on file changes
```

### 4. Leverage Docker Compose

For local development, use the included Docker Compose:

```bash
docker-compose up -d  # Start dependencies
make run             # Run your application
```

### 5. Follow the Structure

The generated projects follow Go best practices:
- `cmd/` for executables
- `internal/` for private code
- `pkg/` for public packages (if needed)
- Configuration in `config/`

### 6. Switching Loggers

Need to switch loggers later? Just regenerate:

```bash
# Started with slog but need better performance?
go-starter new my-api --type web-api --logger zap --force
```

Your application code remains the same - only the logger implementation changes!

## Next Steps

Now that you've mastered the basics:

1. **Explore Blueprints**: Try different project types to see what works best
2. **Experiment with Loggers**: Test different loggers to understand their strengths
3. **Read Blueprint Docs**: Check the blueprint-specific documentation for deeper insights
4. **Join the Community**: Star the [GitHub repo](https://github.com/francknouama/go-starter) and report issues

## Getting Help

- **Documentation**: Full docs at [GitHub](https://github.com/francknouama/go-starter)
- **Issues**: Report bugs or request features on [GitHub Issues](https://github.com/francknouama/go-starter/issues)
- **Examples**: Check the `examples/` directory in the repository

---

**Happy coding!** 🚀 You're now ready to build amazing Go applications with go-starter.