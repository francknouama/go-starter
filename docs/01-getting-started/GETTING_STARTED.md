# Getting Started with go-starter

Welcome to go-starter! This guide will help you create your first Go project using our comprehensive project generator. Whether you're a beginner or an experienced developer, go-starter adapts to your needs with progressive disclosure and smart defaults.

## 🚀 Quick Start (2 Minutes)

### 1. Install go-starter

Choose your preferred installation method:

#### Go Install (Recommended)
```bash
go install github.com/francknouama/go-starter@latest
```

#### Homebrew (macOS/Linux)
```bash
brew tap francknouama/go-starter
brew install go-starter
```

#### Quick Install Script
```bash
curl -sSL https://raw.githubusercontent.com/francknouama/go-starter/main/scripts/install.sh | bash
```

#### Manual Download
Download the latest binary from [releases](https://github.com/francknouama/go-starter/releases).

### 2. Create Your First Project

#### For Beginners - Interactive Mode
```bash
go-starter new
```
This launches an interactive wizard that guides you through all options.

#### For Quick Setup - Direct Generation
```bash
# Simple CLI tool
go-starter new my-tool --type=cli --complexity=simple

# Standard web API
go-starter new my-api --type=web-api

# AWS Lambda function
go-starter new my-lambda --type=lambda
```

### 3. Explore Your Generated Project

```bash
cd my-tool
ls -la
```

Your project is ready to use! Check the README.md for project-specific instructions.

## 📚 Understanding go-starter

### Progressive Disclosure System

go-starter adapts its interface based on your experience level:

- **Basic Mode** (Default): Shows 14 essential options for beginners
- **Advanced Mode**: Shows 18+ options including database, authentication, and deployment settings

```bash
# Basic help (beginner-friendly)
go-starter new --help

# Advanced help (all options)
go-starter new --advanced --help
```

### Project Complexity Levels

| Complexity | Description | File Count | Best For |
|------------|-------------|------------|----------|
| **Simple** | Minimal structure | 8-15 files | Learning, prototypes, simple utilities |
| **Standard** | Production-ready | 25-35 files | Most projects, team development |
| **Advanced** | Enterprise patterns | 40-60 files | Complex business logic, multiple teams |
| **Expert** | Full-featured | 60+ files | Large organizations, complex domains |

## 🏗️ Project Types Guide

### 1. CLI Applications

Perfect for command-line tools, scripts, and utilities.

#### Simple CLI (8 files)
```bash
go-starter new my-tool --type=cli --complexity=simple
```
**When to use**: Quick utilities, learning Go, simple automation scripts

**Generated structure**:
```
my-tool/
├── main.go          # Entry point
├── cmd/
│   ├── root.go      # Root command
│   └── version.go   # Version command
├── config.go        # Configuration
├── Makefile         # Build automation
├── README.md        # Documentation
└── go.mod          # Dependencies
```

#### Standard CLI (29 files)
```bash
go-starter new my-tool --type=cli --complexity=standard
```
**When to use**: Production tools, team projects, complex CLI applications

**Features**: Subcommands, configuration files, advanced logging, testing framework

### 2. Web APIs

RESTful APIs with multiple architecture patterns.

#### Standard Web API
```bash
go-starter new my-api --type=web-api
```
**Best for**: Most web APIs, microservices, standard REST applications

#### Clean Architecture
```bash
go-starter new my-api --type=web-api --architecture=clean
```
**Best for**: Enterprise applications, complex business logic, testable systems

#### Domain-Driven Design (DDD)
```bash
go-starter new my-api --type=web-api --architecture=ddd
```
**Best for**: Complex domains, rich business models, event-driven systems

#### Hexagonal Architecture
```bash
go-starter new my-api --type=web-api --architecture=hexagonal
```
**Best for**: Highly testable systems, multiple adapters, ports & adapters pattern

### 3. AWS Lambda Functions

Serverless functions optimized for AWS Lambda.

#### Simple Lambda
```bash
go-starter new my-lambda --type=lambda
```
**Best for**: Event handlers, webhooks, simple processing

#### Lambda API Proxy
```bash
go-starter new my-api --type=lambda-proxy
```
**Best for**: RESTful APIs on Lambda, API Gateway integration

### 4. Libraries

Reusable Go packages for distribution.

```bash
go-starter new my-lib --type=library
```
**Best for**: Shared utilities, SDK development, open-source packages

### 5. Microservices

Production-ready microservices with gRPC and observability.

```bash
go-starter new my-service --type=microservice
```
**Best for**: Distributed systems, service mesh architectures, cloud-native applications

### 6. Monoliths

Traditional web applications with all components in one deployable unit.

```bash
go-starter new my-app --type=monolith
```
**Best for**: Rapid prototyping, small teams, traditional web applications

### 7. Workspaces

Multi-module projects for monorepos.

```bash
go-starter new my-workspace --type=workspace
```
**Best for**: Monorepos, shared libraries, multiple related services

## ⚙️ Configuration Options

### Database Integration

```bash
# PostgreSQL with GORM
go-starter new my-api --type=web-api --database-driver=postgres --database-orm=gorm

# MongoDB
go-starter new my-api --type=web-api --database-driver=mongodb

# SQLite for development
go-starter new my-api --type=web-api --database-driver=sqlite
```

### Logger Selection

go-starter supports multiple logging libraries with a simplified, unified interface:

```bash
# Standard library (default)
go-starter new my-app --logger=slog

# High performance
go-starter new my-app --logger=zap

# Feature-rich
go-starter new my-app --logger=logrus

# Zero allocation
go-starter new my-app --logger=zerolog
```

### Authentication

```bash
# JWT authentication
go-starter new my-api --type=web-api --auth-type=jwt

# OAuth2
go-starter new my-api --type=web-api --auth-type=oauth2

# Session-based
go-starter new my-api --type=web-api --auth-type=session
```

### Frameworks

#### Web Frameworks
```bash
# Gin (default, fastest)
go-starter new my-api --type=web-api --framework=gin

# Echo (middleware-rich)
go-starter new my-api --type=web-api --framework=echo

# Fiber (Express-like)
go-starter new my-api --type=web-api --framework=fiber
```

#### CLI Frameworks
```bash
# Cobra (default, most popular)
go-starter new my-tool --type=cli --framework=cobra
```

## 🎯 Workflow Examples

### Beginner Developer Workflow

1. **Start with interactive mode**:
   ```bash
   go-starter new
   ```

2. **Choose simple complexity**:
   - Select "CLI Application"
   - Choose "Simple" complexity
   - Accept default options

3. **Explore the generated code**:
   ```bash
   cd my-project
   code .  # Open in VS Code
   ```

4. **Build and run**:
   ```bash
   go run .
   make build
   ./my-project --help
   ```

### Experienced Developer Workflow

1. **Use direct generation with specific options**:
   ```bash
   go-starter new enterprise-api \
     --type=web-api \
     --architecture=hexagonal \
     --database-driver=postgres \
     --database-orm=gorm \
     --auth-type=jwt \
     --logger=zap \
     --framework=gin \
     --advanced
   ```

2. **Review and customize**:
   ```bash
   cd enterprise-api
   cat README.md
   ```

3. **Start development**:
   ```bash
   make dev
   ```

### Team Lead Workflow

1. **Preview project structure**:
   ```bash
   go-starter new team-project \
     --type=microservice \
     --dry-run \
     --architecture=clean
   ```

2. **Generate with team standards**:
   ```bash
   go-starter new team-project \
     --type=microservice \
     --architecture=clean \
     --database-driver=postgres \
     --auth-type=jwt \
     --logger=zap
   ```

3. **Set up team repository**:
   ```bash
   cd team-project
   git init
   git add .
   git commit -m "Initial project setup with go-starter"
   ```

## 🧰 Essential Commands

### Project Generation

```bash
# Interactive mode
go-starter new

# Quick generation
go-starter new <name> --type=<type>

# Advanced configuration
go-starter new <name> --type=<type> --advanced

# Preview without creating files
go-starter new <name> --type=<type> --dry-run
```

### Information Commands

```bash
# List available project types
go-starter list

# Show version
go-starter version

# Show help
go-starter --help

# Advanced help with all options
go-starter new --advanced --help
```

### Utility Commands

```bash
# Generate shell completion
go-starter completion bash > /etc/bash_completion.d/go-starter

# Generate random project name
go-starter new --random-name --type=cli
```

## ✨ Pro Tips

### 1. Use Dry Run for Planning
```bash
go-starter new my-project --type=web-api --dry-run
```
This shows you exactly what files will be generated without creating them.

### 2. Start Simple, Upgrade Later
Begin with simple complexity and refactor to standard when needed:
```bash
# Start simple
go-starter new my-tool --type=cli --complexity=simple

# Later, generate standard version for comparison
go-starter new my-tool-v2 --type=cli --complexity=standard
```

### 3. Leverage Progressive Disclosure
- New to Go? Use basic mode and interactive prompts
- Experienced? Use advanced mode with direct flags
- Team lead? Use dry-run to plan project structure

### 4. Customize After Generation
All generated projects are fully customizable:
- Modify templates in your local copy
- Add additional dependencies
- Adjust configuration files
- Extend with custom middleware

### 5. Use Project Type Selection Matrix

| Use Case | Recommended Type | Complexity | Architecture |
|----------|------------------|------------|--------------|
| Learning Go | CLI | Simple | Standard |
| Utility script | CLI | Simple | Standard |
| Production CLI | CLI | Standard | Standard |
| Simple API | Web API | Standard | Standard |
| Enterprise API | Web API | Advanced | Clean/Hexagonal |
| Event processing | Lambda | Standard | Standard |
| API Gateway | Lambda Proxy | Standard | Standard |
| Shared library | Library | Standard | Standard |
| Cloud service | Microservice | Advanced | Clean |
| Web application | Monolith | Standard | Standard |
| Multiple services | Workspace | Advanced | Standard |

## 🔧 Customization

### Environment Variables

```bash
# Default output directory
export GO_STARTER_OUTPUT_DIR="./projects"

# Default Go version
export GO_STARTER_GO_VERSION="1.21"

# Default author information
export GO_STARTER_AUTHOR="Your Name"
export GO_STARTER_EMAIL="your.email@example.com"
```

### Configuration File

Create `~/.go-starter.yaml` for persistent defaults:

```yaml
author: "Your Name"
email: "your.email@example.com"
license: "MIT"
defaults:
  goVersion: "1.21"
  framework: "gin"
  logger: "slog"
  complexity: "standard"
```

## 🆘 Troubleshooting

### Common Issues

#### 1. Permission Denied
```bash
# Fix: Ensure go-starter is executable
chmod +x $(which go-starter)
```

#### 2. Module Path Issues
```bash
# Fix: Use valid module path format
go-starter new my-project --module=github.com/username/my-project
```

#### 3. Build Failures
```bash
# Check Go version compatibility
go version
# Should be 1.21 or higher
```

#### 4. Template Generation Errors
```bash
# Verify installation
go-starter version

# Reinstall if needed
go install github.com/francknouama/go-starter@latest
```

### Getting Help

- **Documentation**: [github.com/francknouama/go-starter](https://github.com/francknouama/go-starter)
- **Issues**: [Report bugs](https://github.com/francknouama/go-starter/issues)
- **Discussions**: [Community forum](https://github.com/francknouama/go-starter/discussions)
- **CLI Help**: `go-starter --help`

## 🎉 What's Next?

1. **Generate your first project** using this guide
2. **Explore the generated code** and documentation
3. **Build and run** your project
4. **Customize** according to your needs
5. **Share feedback** with the community

Welcome to the go-starter community! We're excited to see what you'll build. 🚀

---

**Quick Reference Card**: For a condensed version of this guide, see [QUICK_REFERENCE.md](QUICK_REFERENCE.md).

**Advanced Topics**: For advanced configuration and customization, see [CONFIGURATION.md](CONFIGURATION.md).

**Template Documentation**: For details about specific project types, see [BLUEPRINTS.md](BLUEPRINTS.md).