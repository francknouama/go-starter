# go-starter

[![Go Report Card](https://goreportcard.com/badge/github.com/francknouama/go-starter)](https://goreportcard.com/report/github.com/francknouama/go-starter)
[![Go Reference](https://pkg.go.dev/badge/github.com/francknouama/go-starter.svg)](https://pkg.go.dev/github.com/francknouama/go-starter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/francknouama/go-starter)](https://github.com/francknouama/go-starter/releases)
[![GitHub Actions](https://github.com/francknouama/go-starter/workflows/Release/badge.svg)](https://github.com/francknouama/go-starter/actions)

A comprehensive Go project generator that combines the simplicity of create-react-app with the flexibility of Spring Initializr. Features a unique **Logger Selector System** that lets you choose your logging strategy without vendor lock-in.

## 🌟 Key Features

### ⚡ Unique Logger Selector System

**Choose your logging strategy, not your vendor.** go-starter's Logger Selector System is the first of its kind:

| Logger | Performance | Best For | Zero Allocation |
|--------|-------------|----------|-----------------|
| **slog** | Good | Standard library choice | ✅ |
| **zap** | Excellent | High-performance APIs | ✅ |
| **logrus** | Good | Feature-rich applications | ❌ |
| **zerolog** | Excellent | JSON-heavy workloads | ✅ |

- **Switch Anytime**: Change loggers without touching application code
- **Consistent Interface**: Same logging calls across all implementations  
- **Clean Dependencies**: Only selected logger included in your project
- **Performance Optimized**: Each logger tuned for its specific strengths

### 🚀 v1.1.0 - Latest Features

- 📦 **Multi-Database Selection**: PostgreSQL, MySQL, MongoDB, SQLite, Redis support
- 🔄 **Dynamic Go Version Detection**: Uses your current Go version automatically
- 🛠️ **Enhanced CLI**: Better error handling when Go is not installed
- 🐳 **Multi-Service Docker**: Docker Compose for multiple databases
- 🎯 **16 Tested Combinations**: All template+logger combinations validated
- ⚡ **Instant Setup**: Generate complete, compilable projects in under 10 seconds

### 🛣️ Roadmap - Future Development
See [PROJECT_ROADMAP.md](PROJECT_ROADMAP.md) for detailed future plans:
- 🏗️ **Advanced Architecture Templates**: Clean Architecture, DDD, Hexagonal patterns
- 📦 **Additional Frameworks**: Echo, Fiber, Chi web framework support  
- 🔧 **Enhanced Features**: Database drivers, authentication methods, deployment options
- 🌐 **Web UI Interface**: Browser-based project generator with live preview
- 🏢 **Enterprise Templates**: Microservices, event-driven, workspace patterns

## 💻 Installation

### Using Homebrew (Recommended)

```bash
brew tap francknouama/tap
brew install go-starter
```

### Using Go Install

```bash
go install github.com/francknouama/go-starter@latest
```

### Download Binary

Download the latest release for your platform from [GitHub Releases](https://github.com/francknouama/go-starter/releases/latest).

```bash
# Example for Linux AMD64
curl -L https://github.com/francknouama/go-starter/releases/download/v1.1.0/go-starter_1.1.0_Linux_x86_64.tar.gz | tar -xz
sudo mv go-starter /usr/local/bin/

# Example for macOS Apple Silicon
curl -L https://github.com/francknouama/go-starter/releases/download/v1.1.0/go-starter_1.1.0_Darwin_arm64.tar.gz | tar -xz
sudo mv go-starter /usr/local/bin/

# Example for Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/francknouama/go-starter/releases/download/v1.1.0/go-starter_1.1.0_Windows_x86_64.zip" -OutFile "go-starter.zip"
Expand-Archive go-starter.zip -DestinationPath .
# Add to PATH or move to desired location
```

### Using Package Managers

**Linux packages** (deb, rpm, apk) are available from [GitHub Releases](https://github.com/francknouama/go-starter/releases/latest):

```bash
# Debian/Ubuntu
wget https://github.com/francknouama/go-starter/releases/download/v1.1.0/go-starter_1.1.0_linux_amd64.deb
sudo dpkg -i go-starter_1.1.0_linux_amd64.deb

# RHEL/CentOS/Fedora  
wget https://github.com/francknouama/go-starter/releases/download/v1.1.0/go-starter_1.1.0_linux_amd64.rpm
sudo rpm -i go-starter_1.1.0_linux_amd64.rpm

# Alpine Linux
wget https://github.com/francknouama/go-starter/releases/download/v1.1.0/go-starter_1.1.0_linux_amd64.apk
sudo apk add --allow-untrusted go-starter_1.1.0_linux_amd64.apk
```

### Verify Installation

```bash
go-starter version
go-starter list
```

### From Source

```bash
git clone https://github.com/francknouama/go-starter.git
cd go-starter
make install
```

## 📚 Documentation

- 📖 **[Getting Started Guide](docs/GETTING_STARTED.md)** - Comprehensive guide for beginners
- 🚀 **[Quick Reference](docs/QUICK_REFERENCE.md)** - Command cheatsheet and examples
- 📋 **[Release Notes](RELEASE_NOTES.md)** - Latest features and changes

## 🚀 Quick Start

### ⚡ Logger Selector in Action

See how easy it is to switch between different logging strategies:

```bash
# High-performance API with zap (zero allocation)
go-starter new fast-api --type web-api --logger zap

# Feature-rich app with structured logging
go-starter new rich-app --type web-api --logger logrus  

# Standard library approach
go-starter new simple-api --type web-api --logger slog

# JSON-optimized service
go-starter new json-service --type web-api --logger zerolog
```

**Same application code, different performance characteristics!**

### 🎯 Interactive Mode (Recommended for beginners)

```bash
go-starter new my-awesome-project

# Follow the interactive prompts:
? Project type: › Web API
? Framework: › gin  
? Logger: › zap (High-performance, zero allocation)
? Database: › PostgreSQL, Redis
? Module path: › github.com/yourusername/my-awesome-project

✅ Project generated successfully!
🚀 Run 'cd my-awesome-project && make run' to start development
```

### ⚡ Direct Mode (For experienced developers)

```bash
# Web API with high-performance logging
go-starter new my-api --type=web-api --framework=gin --logger=zap

# CLI tool with structured logging  
go-starter new my-cli --type=cli --framework=cobra --logger=logrus

# Go library with standard logging
go-starter new my-lib --type=library --logger=slog

# AWS Lambda with zero-allocation JSON logging
go-starter new my-lambda --type=lambda --logger=zerolog
```

### 🔥 What You Get Instantly

Every generated project includes:
- ✅ **Compiles immediately** - no setup required
- ✅ **Production-ready structure** with best practices
- ✅ **Complete testing setup** with examples
- ✅ **Docker configuration** for containerization  
- ✅ **Makefile** with common development tasks
- ✅ **GitHub Actions** CI/CD pipeline
- ✅ **Comprehensive documentation** and examples

## 🏗️ Project Types (v1.0.0)

### 🌐 Web API
Production-ready REST API with Gin framework:
- **Framework**: Gin (Echo, Fiber, Chi planned for future)
- **Architecture**: Standard structure (Clean/DDD/Hexagonal patterns planned)
- **Features**: Middleware, routing, health checks, Docker support
- **Generated**: Complete API with database integration, tests, CI/CD

### 🖥️ CLI Application  
Professional command-line tools with Cobra:
- **Framework**: Cobra with subcommands and configuration
- **Features**: Interactive prompts, completion, version management
- **Use Cases**: DevOps tools, utilities, automation scripts
- **Generated**: Complete CLI with config management, tests, Docker support

### 📦 Go Library
Well-structured reusable packages:
- **Features**: Clean public API, comprehensive documentation, examples
- **Testing**: Unit tests, benchmarks, CI/CD integration
- **Use Cases**: SDKs, shared functionality, open source packages
- **Generated**: Complete library with examples, docs, and publishing setup

### ⚡ AWS Lambda
Optimized serverless functions:
- **Runtime**: AWS Lambda Go runtime with API Gateway integration
- **Logging**: CloudWatch-optimized structured logging
- **Deployment**: SAM templates with automated deployment scripts
- **Generated**: Complete Lambda with infrastructure-as-code and CI/CD

## 📝 Logger Options

Choose the perfect logging solution for your needs:

| Logger | Performance | Use Case | Key Features |
|--------|-------------|----------|-------------|
| **slog** ⭐ | Good | General purpose, stdlib | Standard library, structured logging, Go 1.21+ |
| **zap** | Excellent | High performance apps | Zero allocation, blazing fast, Uber's choice |
| **logrus** | Good | Feature-rich apps | JSON/Text, hooks, popular ecosystem |
| **zerolog** | Excellent | Cloud-native, APIs | Zero allocation, chainable, minimal memory |

### 🔄 Consistent Interface

All loggers implement the same interface, so you can switch between them without changing your code:

```go
// Works with any logger choice
logger.Info("Server starting", "port", 8080, "env", "production")
logger.Error("Database connection failed", "error", err)
logger.Debug("Processing request", "method", "GET", "path", "/api/users")
```

**💡 Recommendation:**
- **slog** for most projects (stdlib, no dependencies)
- **zap** for high-throughput applications
- **zerolog** for cloud/container deployments
- **logrus** for feature-rich logging needs

## 📊 Implementation Status

### ✅ Currently Available (v1.0.0)
| Feature | Templates | Loggers | Status |
|---------|-----------|---------|--------|
| **Project Types** | 4 (web-api, cli, library, lambda) | 4 (slog, zap, logrus, zerolog) | ✅ Production Ready |
| **Template Combinations** | 16 total combinations | All tested | ✅ Fully Validated |
| **Frameworks** | Gin (web), Cobra (cli) | - | ✅ Complete |
| **Architecture Patterns** | Standard | - | ✅ Complete |
| **Docker Support** | All templates | - | ✅ Complete |
| **CI/CD Integration** | GitHub Actions | - | ✅ Complete |

### 🔮 Planned for Future Releases
| Feature | Target | Status |
|---------|--------|--------|
| Clean Architecture Templates | Phase 8 | ❌ Not Started |
| Additional Web Frameworks | Phase 7 | ❌ Not Started |
| Database Driver Selection | Phase 7 | ❌ Not Started |
| Web UI Interface | Phase 9 | ❌ Not Started |
| Microservice Templates | Phase 8 | ❌ Not Started |

*See [PROJECT_ROADMAP.md](PROJECT_ROADMAP.md) for detailed development timeline*

## Configuration

### Global Configuration

```yaml
# ~/.go-starter.yaml
profiles:
  default:
    author: "John Doe"
    email: "john@example.com"
    license: "MIT"
    defaults:
      goVersion: "1.21"
      framework: "gin"
      logger: "slog"
```

### Advanced Mode

Enable advanced configuration options:

```bash
go-starter new my-project --advanced
```

Advanced mode includes:
- Database selection (PostgreSQL, MySQL, MongoDB, SQLite)
- Authentication methods (JWT, OAuth2, API Key)
- Message queue integration (RabbitMQ, Kafka, Redis)
- Observability tools (Prometheus, Jaeger, OpenTelemetry)

## 🚀 Real-World Examples

### Building a REST API
```bash
# Generate a high-performance API
go-starter new user-service --type=web-api --framework=gin --logger=zap

cd user-service
make run    # Starts server on :8080
make test   # Runs all tests
make build  # Creates production binary
make docker # Builds Docker image
```

### Creating a CLI Tool
```bash
# Generate a CLI application
go-starter new deployment-tool --type=cli --framework=cobra --logger=logrus

cd deployment-tool
go run main.go --help           # See available commands
go run main.go version          # Check version
make build && ./bin/deployment-tool deploy --env=prod
```

### Publishing a Go Library
```bash
# Generate a reusable library
go-starter new awesome-sdk --type=library --logger=slog

cd awesome-sdk
make test      # Run tests and benchmarks
make lint      # Check code quality
make docs      # Generate documentation
```

### Deploying to AWS Lambda
```bash
# Generate a Lambda function
go-starter new data-processor --type=lambda --logger=zerolog

cd data-processor
make build-lambda  # Cross-compile for Linux
make deploy        # Deploy with SAM
make logs          # View CloudWatch logs
```

## Development

### Prerequisites

- Go 1.21+
- Make
- Git

### Building

```bash
# Build the CLI
make build

# Run tests
make test

# Run linting
make lint
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Roadmap

- [x] **Phase 1-4: Core Templates + Logger Selector** ✅
  - [x] Web API template with Gin framework
  - [x] CLI application template with Cobra
  - [x] Go library template for reusable packages
  - [x] AWS Lambda template for serverless functions
  - [x] Logger selector (slog, zap, logrus, zerolog)
  - [x] Conditional dependencies and consistent interfaces
- [ ] **Phase 5: Enhancements** (Optional)
  - [ ] Additional frameworks and database drivers
  - [ ] Web UI with live preview
  - [ ] Template marketplace and GitHub integration

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by create-react-app and Spring Initializr
- Built with [Cobra](https://github.com/spf13/cobra) and [Viper](https://github.com/spf13/viper)
- Templates use Go's text/template with [Sprig](https://github.com/Masterminds/sprig) functions

## 📚 Documentation

- 📋 **[Template Usage Guide](docs/TEMPLATES.md)** - Comprehensive guide for all project types
- 🪵 **[Logger Selector Guide](docs/LOGGER_GUIDE.md)** - Complete logging documentation
- ❓ **[Troubleshooting & FAQ](docs/FAQ.md)** - Common issues and solutions
- 🛣️ **[Project Roadmap](PROJECT_ROADMAP.md)** - Future development plans

## Support

- 📖 [Documentation](docs/)
- 🐛 [Issue Tracker](https://github.com/francknouama/go-starter/issues)
- 💬 [Discussions](https://github.com/francknouama/go-starter/discussions)