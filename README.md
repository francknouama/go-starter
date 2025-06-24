# go-starter

[![Go Report Card](https://goreportcard.com/badge/github.com/francknouama/go-starter)](https://goreportcard.com/report/github.com/francknouama/go-starter)
[![Go Reference](https://pkg.go.dev/badge/github.com/francknouama/go-starter.svg)](https://pkg.go.dev/github.com/francknouama/go-starter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/francknouama/go-starter)](https://github.com/francknouama/go-starter/releases)
[![GitHub Actions](https://github.com/francknouama/go-starter/workflows/Release/badge.svg)](https://github.com/francknouama/go-starter/actions)

A comprehensive Go project generator that combines the simplicity of create-react-app with the flexibility of Spring Initializr, offering both CLI and web interfaces with progressive disclosure for beginners and advanced developers.

## Features

### ✨ v1.0.0 - Production Ready!

- 🚀 **4 Project Types**: Web API (Gin), CLI application (Cobra), Go library, and AWS Lambda
- 📝 **Logger Selection**: Choose between slog, zap, logrus, or zerolog with consistent interface
- 🎯 **16 Combinations**: All template+logger combinations tested and validated
- 🔧 **Best Practices**: Pre-configured linting, testing, and development tools
- 🐳 **Docker Ready**: Multi-stage Dockerfiles for production deployment
- ⚡ **Fast Setup**: Generate production-ready projects in seconds

### 🛣️ Coming Soon (See [Roadmap](PROJECT_ROADMAP.md))
- 🏗️ **Advanced Architectures**: Clean Architecture, DDD, Hexagonal patterns
- 📦 **More Frameworks**: Echo, Fiber, Chi support
- 🌐 **Web UI**: Browser-based project generator
- 🔐 **Enterprise Features**: Authentication, databases, microservices

## Installation

### Using Go Install (Recommended)

```bash
go install github.com/francknouama/go-starter@latest
```

### Download Binary

Download the latest release for your platform from [GitHub Releases](https://github.com/francknouama/go-starter/releases/latest).

```bash
# Example for Linux
curl -L https://github.com/francknouama/go-starter/releases/download/v1.0.0/go-starter_1.0.0_Linux_x86_64.tar.gz | tar -xz
sudo mv go-starter /usr/local/bin/
```

### Using Homebrew (macOS/Linux)

```bash
brew tap francknouama/tap
brew install go-starter
```

### Using Docker

```bash
docker pull francknouama/go-starter:latest
docker run --rm -v $(pwd):/workspace francknouama/go-starter:latest new my-project
```

### From Source

```bash
git clone https://github.com/francknouama/go-starter.git
cd go-starter
make install
```

## Quick Start

### Interactive Mode

```bash
go-starter new my-awesome-api
```

### Direct Mode

```bash
# Create a web API with Gin and Zap logger
go-starter new my-api --type=web-api --framework=gin --logger=zap

# Create a CLI application with Cobra and Logrus
go-starter new my-cli --type=cli --framework=cobra --logger=logrus

# Create a Go library with slog
go-starter new my-lib --type=library --logger=slog

# Create an AWS Lambda function
go-starter new my-lambda --type=lambda --logger=zerolog
```

## Project Types

### Web API
REST API with your choice of web framework:
- **Frameworks**: Gin, Echo, Fiber, Chi
- **Features**: Middleware, routing, validation, OpenAPI docs
- **Architecture**: Standard, Clean, DDD, or Hexagonal

### CLI Application
Command-line tools with Cobra framework:
- **Features**: Subcommands, flags, configuration management
- **Use Cases**: DevOps tools, utilities, automation scripts

### Library
Reusable Go packages:
- **Features**: Well-documented API, examples, benchmarks
- **Use Cases**: Shared functionality, SDK development

### AWS Lambda
Serverless functions for AWS:
- **Features**: API Gateway integration, CloudWatch logging
- **Deployment**: SAM templates, deployment scripts

## Logger Options

Choose from four popular logging libraries:

- **slog**: Go's standard structured logging (default)
- **zap**: Uber's high-performance logger
- **logrus**: Feature-rich structured logger
- **zerolog**: Zero-allocation JSON logger

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

## Support

- 📖 [Documentation](https://github.com/francknouama/go-starter/wiki)
- 🐛 [Issue Tracker](https://github.com/francknouama/go-starter/issues)
- 💬 [Discussions](https://github.com/francknouama/go-starter/discussions)