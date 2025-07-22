# go-starter v1.1.0 Release Notes

## 🎉 Introducing go-starter: The Comprehensive Go Project Generator

We're excited to announce the release of **go-starter v1.1.0**, a powerful CLI tool that revolutionizes how Go developers bootstrap new projects. With its unique **Logger Selector System** and modern best practices baked in, go-starter helps you create production-ready Go applications in seconds.

## 🌟 Key Features

### 🔧 Logger Selector System - Choose Your Logging Strategy

One of go-starter's standout features is the **Logger Selector System**, allowing you to choose from four high-performance logging libraries while maintaining a consistent interface across your entire codebase:

| Logger | Package | Performance | Best For |
|--------|---------|-------------|----------|
| **slog** | `log/slog` | Good | Standard library choice, structured logging |
| **zap** | `go.uber.org/zap` | Excellent | High-performance, zero allocation |
| **logrus** | `github.com/sirupsen/logrus` | Good | Feature-rich, popular choice |
| **zerolog** | `github.com/rs/zerolog` | Excellent | Zero allocation, chainable API |

#### Why This Matters:
- **No Lock-in**: Switch between loggers without changing your application code
- **Performance Optimization**: Choose the logger that best fits your performance requirements
- **Clean Dependencies**: Only the selected logger's dependencies are included in your project
- **Consistent Interface**: All loggers implement the same interface for seamless integration

### 📦 Project Templates

go-starter v1.1.0 ships with 4 production-ready templates:

1. **Web API** - REST API with Gin framework
   - Database integration (PostgreSQL, MySQL, MongoDB, SQLite, Redis)
   - Docker support with multi-stage builds
   - Health checks and graceful shutdown
   - Structured project layout

2. **CLI Application** - Command-line tools with Cobra
   - Subcommands and flag management
   - Configuration file support
   - Auto-completion scripts

3. **Go Library** - Reusable packages
   - Clean public API design
   - Comprehensive examples
   - Proper documentation structure

4. **AWS Lambda** - Serverless functions
   - CloudWatch-optimized logging
   - Event handling patterns
   - SAM/Serverless Framework ready

### 🚀 What's New in v1.1.0

- **Multi-Database Selection**: Select multiple databases during project generation
- **Dynamic Go Version Detection**: Automatically uses your current Go version
- **Improved CLI Behavior**: Graceful handling when Go is not installed
- **Enhanced Docker Support**: Multi-service Docker Compose for multiple databases

## 💻 Installation

### Using Homebrew (Recommended)
```bash
brew tap francknouama/tap
brew install go-starter
```

### Using Go Install
```bash
go install github.com/francknouama/go-starter@v1.1.0
```

### Download Binary
Download the appropriate binary for your platform from the [releases page](https://github.com/francknouama/go-starter/releases/tag/v1.1.0).

## 🎯 Quick Start

### Interactive Mode (Recommended for Beginners)
```bash
go-starter new my-awesome-api
```

The interactive mode guides you through:
- Project type selection
- Framework choice (for applicable templates)
- Logger selection
- Database configuration
- Advanced options

### Direct Mode (For Power Users)
```bash
# Create a web API with zap logger and PostgreSQL
go-starter new my-api --type web-api --logger zap --database postgres

# Create a CLI tool with logrus logger
go-starter new my-cli --type cli --logger logrus

# Create a library with the default slog logger
go-starter new my-lib --type library
```

## 📊 Logger Selection in Action

### Example: Creating a High-Performance API
```bash
$ go-starter new high-perf-api --type web-api --logger zap
✓ Project 'high-perf-api' created successfully!
✓ Type: web-api
✓ Framework: gin
✓ Logger: zap (high-performance, zero allocation)
✓ Files created: 23
✓ Git repository initialized

Get started:
  cd high-perf-api
  make run
```

### Generated Logger Configuration
```go
// internal/logger/logger.go
package logger

import "go.uber.org/zap"

var log *zap.Logger

func init() {
    config := zap.NewProductionConfig()
    config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
    log, _ = config.Build()
}

func Get() *zap.Logger {
    return log
}
```

## 🏗️ Generated Project Structure

```
my-awesome-api/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── handlers/            # HTTP handlers
│   ├── logger/              # Logger implementation (your choice!)
│   ├── middleware/          # HTTP middleware
│   └── models/              # Data models
├── pkg/                     # Public packages
├── docker/
│   └── Dockerfile           # Multi-stage production build
├── docker-compose.yml       # Local development setup
├── Makefile                 # Common tasks
├── go.mod                   # Go modules
└── README.md               # Project documentation
```

## 📈 Performance Comparison

Based on our benchmarks with the Logger Selector System:

| Logger | Allocations/Op | ns/Op | Best Use Case |
|--------|---------------|-------|---------------|
| slog | 0 | 45 | Standard applications |
| zap | 0 | 28 | High-throughput APIs |
| logrus | 23 | 584 | Feature-rich applications |
| zerolog | 0 | 31 | JSON-heavy workloads |

## 🔄 Switching Loggers

One of the most powerful features is the ability to switch loggers without changing application code:

```bash
# Start with slog (default)
go-starter new my-api --type web-api

# Later, need better performance? Regenerate with zap:
go-starter new my-api --type web-api --logger zap --force
```

Your application code remains unchanged - only the logger implementation swaps out!

## 🛠️ Development Experience

### Instant Productivity
```bash
$ go-starter new my-service --type web-api
$ cd my-service
$ make run

INFO[0000] Starting server on :8080
INFO[0000] Database connected successfully
INFO[0000] Server is ready to handle requests
```

### Built-in Development Commands
- `make run` - Start the application
- `make test` - Run tests with coverage
- `make lint` - Run golangci-lint
- `make build` - Build production binary
- `make docker` - Build Docker image

## 🤝 Community and Support

- **GitHub**: [github.com/francknouama/go-starter](https://github.com/francknouama/go-starter)
- **Issues**: [Report bugs or request features](https://github.com/francknouama/go-starter/issues)
- **Discussions**: [Join the conversation](https://github.com/francknouama/go-starter/discussions)

## 🎯 Who Should Use go-starter?

- **Beginners**: Get started with Go using production-ready patterns
- **Teams**: Standardize project structure across your organization
- **Freelancers**: Bootstrap client projects quickly
- **Open Source**: Create consistent, well-structured libraries

## 🚀 What's Next?

We're actively working on:
- Additional templates (Clean Architecture, DDD, Hexagonal, Microservices)
- More web frameworks (Echo, Fiber, Chi)
- Authentication templates (JWT, OAuth2)
- Web UI for visual project configuration
- Template marketplace for community contributions

## 📝 Migration from Other Tools

Coming from other project generators? go-starter offers:
- More comprehensive templates than `go mod init`
- Logger flexibility not found in other generators
- Better defaults than manual project setup
- Active maintenance and community support

## 🙏 Acknowledgments

Special thanks to the Go community for inspiration and feedback. This project stands on the shoulders of giants:
- The Go team for the excellent standard library
- Uber for zap
- Sirupsen for logrus  
- RS for zerolog
- All our early adopters and contributors

---

**Ready to supercharge your Go development?** Install go-starter today and experience the difference a well-crafted project generator can make!

```bash
brew tap francknouama/tap
brew install go-starter
go-starter new my-next-project
```

Happy coding! 🚀