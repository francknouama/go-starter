# Getting Started with go-starter

Welcome to **go-starter**! This guide will help you get up and running with the most powerful Go project generator available. In just a few minutes, you'll be creating production-ready Go applications with the perfect logging strategy for your needs.

## Table of Contents

- [Installation](#installation)
- [Your First Project](#your-first-project)
- [Understanding the Logger Selector](#understanding-the-logger-selector)
- [Project Types](#project-types)
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
Version: 1.1.0
Commit: 760abac
Built: 2025-06-25
```

## Your First Project

### Interactive Mode (Recommended for Beginners)

The easiest way to get started is using interactive mode:

```bash
go-starter new my-first-api
```

You'll be guided through a series of prompts:

```
? Project type: 
  ❯ Web API - REST API or web service
    CLI Application - Command-line tool
    Library - Reusable Go package
    AWS Lambda - Serverless function

? Framework: 
  ❯ gin (recommended)

? Logger type: 
  ❯ slog - Standard library structured logging
    zap - High performance, zero allocation
    logrus - Feature-rich, popular logger
    zerolog - Zero allocation JSON logger

? Database support? 
  ❯ Yes
    No

? Select databases (space to select, enter to continue):
  ❯ ◉ PostgreSQL
    ◯ MySQL
    ◉ Redis
    ◯ MongoDB
    ◯ SQLite

? Module path: github.com/yourusername/my-first-api
```

### Direct Mode (For Power Users)

If you know exactly what you want, use direct mode:

```bash
# Create a high-performance API with zap logger
go-starter new fast-api --type web-api --logger zap --database postgres

# Create a CLI tool with logrus
go-starter new mycli --type cli --logger logrus

# Create a library with standard slog
go-starter new mylib --type library
```

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

1. **Explore Templates**: Try different project types to see what works best
2. **Experiment with Loggers**: Test different loggers to understand their strengths
3. **Read Template Docs**: Check the template-specific documentation for deeper insights
4. **Join the Community**: Star the [GitHub repo](https://github.com/francknouama/go-starter) and report issues

## Getting Help

- **Documentation**: Full docs at [GitHub](https://github.com/francknouama/go-starter)
- **Issues**: Report bugs or request features on [GitHub Issues](https://github.com/francknouama/go-starter/issues)
- **Examples**: Check the `examples/` directory in the repository

---

**Happy coding!** 🚀 You're now ready to build amazing Go applications with go-starter.