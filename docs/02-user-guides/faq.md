# Frequently Asked Questions (FAQ)

Quick answers to the most common questions about go-starter. Can't find what you're looking for? Check our [troubleshooting guide](troubleshooting.md) or [get help](#getting-help).

## Table of Contents

- [Getting Started](#getting-started)
- [Project Types & Blueprints](#project-types--blueprints)
- [Logger Selection](#logger-selection)
- [Configuration & Customization](#configuration--customization)
- [Technical Questions](#technical-questions)
- [Best Practices](#best-practices)
- [Common Issues](#common-issues)
- [Getting Help](#getting-help)

## Getting Started

### What is go-starter?

go-starter is a comprehensive Go project generator that creates production-ready projects with modern best practices. It features:

- **12 production-ready blueprints** (CLI, Web API, Lambda, etc.)
- **Progressive disclosure system** that adapts to your experience level
- **Simplified logger system** with unified interface across 4 popular loggers
- **Multiple architecture patterns** (Standard, Clean, DDD, Hexagonal)
- **Best practices built-in** including testing, Docker, CI/CD

### How is go-starter different from other generators?

**Unique Features**:
- **Progressive Disclosure**: Beginner mode (14 options) vs Advanced mode (18+ options)
- **Logger Selector**: Choose from slog, zap, logrus, zerolog with consistent interface
- **Complexity Levels**: Simple (8 files) to Expert (60+ files) blueprints
- **Architecture Patterns**: 4 different patterns for web APIs
- **Production Ready**: All blueprints are tested and include deployment configs

### What are the system requirements?

**Minimum Requirements**:
- **Go 1.21+** (recommended for slog support)
- **Git** (optional but recommended)
- **Make** (optional, for using Makefile commands)

**Platform Support**:
- ✅ Linux (all distributions)
- ✅ macOS (Intel and Apple Silicon)
- ✅ Windows (PowerShell and CMD)
- ✅ Docker (all platforms)

### How do I install go-starter?

**Recommended Method**:
```bash
go install github.com/francknouama/go-starter@latest
go-starter version
```

**Alternative Methods**:
- **Homebrew**: `brew install francknouama/go-starter/go-starter`
- **Binary Download**: From [GitHub Releases](https://github.com/francknouama/go-starter/releases)
- **Package Managers**: Available for Chocolatey, Scoop, Snap, AUR

See the [installation guide](../01-getting-started/installation.md) for detailed instructions.

## Project Types & Blueprints

### Which project type should I choose?

**Quick Guide**:

| Use Case | Recommended Type | Command |
|----------|------------------|---------|
| Learning Go | CLI Simple | `--type=cli --complexity=simple` |
| REST API | Web API Standard | `--type=web-api` |
| Enterprise System | Web API Clean | `--type=web-api --architecture=clean` |
| Developer Tool | CLI Standard | `--type=cli --complexity=standard` |
| Serverless Function | Lambda | `--type=lambda` |
| Shared Library | Library | `--type=library` |
| Complex Domain | Web API DDD | `--type=web-api --architecture=ddd` |
| Maximum Testability | Web API Hexagonal | `--type=web-api --architecture=hexagonal` |

See the [blueprint selection guide](blueprint-selection.md) for detailed decision trees.

### What's the difference between CLI Simple and Standard?

**CLI Simple (8 files)**:
- Perfect for learning Go
- Single command with basic flags
- Minimal structure and dependencies
- Great for scripts and utilities

**CLI Standard (29 files)**:
- Production-ready CLI applications
- Multiple subcommands support
- Configuration files and advanced logging
- Comprehensive testing framework
- Shell completion scripts

### What architecture patterns are available?

**Web API Architectures**:

1. **Standard** (35 files): Traditional layered architecture, fast development
2. **Clean** (45 files): Clean Architecture principles, highly testable
3. **DDD** (50 files): Domain-Driven Design, complex business logic
4. **Hexagonal** (55 files): Ports & Adapters, maximum testability

**When to Use Each**:
- **Standard**: MVPs, simple APIs, rapid prototyping
- **Clean**: Enterprise apps, complex business logic, long-term projects
- **DDD**: Domain-rich applications, event-driven systems
- **Hexagonal**: Multiple interfaces (HTTP, gRPC, CLI), maximum testability

### Can I migrate between architectures later?

**Scaling Up**: Yes, with effort
- Simple → Standard: Straightforward refactoring
- Standard → Clean: Requires restructuring business logic
- Clean → DDD: Need to identify domain boundaries
- Any → Hexagonal: Most complex, requires interface design

**Scaling Down**: Easier
- Complex architectures can be simplified by removing layers

**Recommendation**: Start simple and migrate when complexity justifies it.

## Logger Selection

### Which logger should I choose?

**Quick Comparison**:

| Logger | Performance | Use Case | Dependencies |
|--------|-------------|----------|-------------|
| **slog** | Good | Default choice, Go 1.21+ | None (stdlib) |
| **zap** | Excellent | High-performance APIs | go.uber.org/zap |
| **logrus** | Good | Feature-rich apps | github.com/sirupsen/logrus |
| **zerolog** | Excellent | JSON-heavy, cloud-native | github.com/rs/zerolog |

**Recommendations**:
- **New to Go**: Use `slog` (standard library)
- **High-performance APIs**: Use `zap` or `zerolog`
- **Rich formatting needs**: Use `logrus`
- **Cloud/container environments**: Use `zerolog`

### Can I switch loggers after generation?

**Simplified Logger System**: go-starter uses a unified interface across all loggers, but the implementation differs.

**Options**:
1. **Regenerate project** with different logger (recommended)
2. **Manual replacement** of logger implementation files
3. **Copy business logic** to new project with desired logger

**Note**: The application code using the logger remains the same due to the unified interface.

### What's the unified logger interface?

All loggers use the same simple API:

```go
// Same interface regardless of logger chosen
logger.Info("User logged in", "user_id", userID, "ip", clientIP)
logger.Error("Database error", "error", err)
logger.Debug("Processing request", "request_id", reqID)
logger.Warn("Rate limit approaching", "limit", rateLimit)
```

The implementation differs, but your application code stays consistent.

## Configuration & Customization

### How do I set default preferences?

Create a configuration file at `~/.go-starter.yaml`:

```yaml
profiles:
  default:
    author: "Your Name"
    email: "your.email@example.com"
    license: "MIT"
    defaults:
      goVersion: "1.22"
      framework: "gin"
      logger: "slog"
      complexity: "standard"
current_profile: "default"
```

See the [configuration guide](configuration.md) for advanced setup.

### Can I customize the generated projects?

**After Generation**: ✅ Yes, modify any generated files
**Blueprint Templates**: ❌ Templates are embedded in the binary
**Configuration**: ✅ Use config files to influence generation

**Customization Options**:
- Modify generated code after creation
- Use environment variables for settings
- Fork the repository to modify blueprints
- Use configuration profiles for team standards

### How do I set up team standards?

**Team Configuration File**:
```yaml
# team-standards.yaml
profiles:
  team:
    author: "{{DEVELOPER_NAME}}"
    email: "{{DEVELOPER_EMAIL}}"
    license: "Proprietary"
    defaults:
      goVersion: "1.22"
      framework: "gin"
      logger: "zap"
      architecture: "clean"
      database: "postgres"
```

**Team Setup Script**:
```bash
# Download and apply team standards
curl -o team-config.yaml https://company.com/go-starter-config.yaml
go-starter config import team-config.yaml
go-starter config set-profile team
```

### Can I use go-starter in corporate environments?

**✅ Yes!** go-starter is designed for professional use:

- **License**: MIT license allows commercial use
- **Security**: No external calls during generation
- **Proxy Support**: Works with corporate proxies
- **Air-gapped**: Can work offline after initial install
- **Compliance**: Generates secure code following best practices

## Technical Questions

### What Go version is required?

**For go-starter tool**:
- Minimum: Go 1.21+
- Recommended: Go 1.22+

**For generated projects**:
- Depends on chosen logger and features
- slog requires Go 1.21+
- Other loggers work with Go 1.19+

### Is go-starter production-ready?

**✅ Absolutely!** All blueprints include:

- Comprehensive test suites with high coverage
- Production configuration examples
- Docker and CI/CD configurations
- Proper error handling and recovery
- Security best practices
- Performance optimizations
- Monitoring and observability setup

### How do I handle secrets and environment variables?

**Generated projects include**:
- Environment-based configuration
- Config file structure for different environments
- Examples for secret management

```yaml
# config/config.prod.yaml
database:
  host: ${DB_HOST}
  password: ${DB_PASSWORD}
  
# Environment variables
export DB_HOST=prod-db.example.com
export DB_PASSWORD=secure-password
```

### Can I add databases to my project?

**Current Support**:
- Database configuration structure
- Connection examples for popular databases
- GORM integration examples
- Migration setup

**Future Versions** will include:
- Interactive database selection
- Multiple ORM options (GORM, SQLx, SQLC, Ent)
- Automatic migration generation
- Database-specific optimizations

### How do I deploy generated projects?

**Multiple Options Available**:

```bash
# Docker deployment
make docker
docker run -p 8080:8080 my-app:latest

# Binary deployment
make build
./bin/server --config=configs/config.prod.yaml

# Cloud platforms
# AWS Lambda, Google Cloud Run, Azure Container Instances
```

All blueprints include deployment configurations and examples.

## Best Practices

### What's the recommended development workflow?

**Standard Workflow**:

1. **Generate**: `go-starter new my-project --type=web-api`
2. **Setup**: `cd my-project && go mod tidy`
3. **Test**: `make test`
4. **Develop**: `make run` (with hot reload)
5. **Test Continuously**: `make test-watch`
6. **Build**: `make build`
7. **Deploy**: Use Docker or binary deployment

### How should I structure my code in generated projects?

**Generated Structure Follows Go Best Practices**:

```
my-project/
├── cmd/                # Application entry points
├── internal/           # Private application code
│   ├── handlers/       # HTTP handlers
│   ├── services/       # Business logic
│   ├── repository/     # Data access
│   └── logger/         # Logging implementation
├── pkg/                # Public library code (if applicable)
├── configs/            # Configuration files
├── tests/              # Test files and test data
├── Dockerfile          # Container configuration
├── Makefile           # Build automation
└── README.md          # Project documentation
```

### Should I start simple or advanced?

**Recommendation: Start Simple, Grow as Needed**

**Learning Path**:
1. **Week 1**: CLI Simple for Go basics
2. **Week 2**: CLI Standard for production patterns
3. **Week 3**: Web API Standard for web development
4. **Week 4**: Clean Architecture for advanced patterns

**Project Evolution**:
- Begin with Standard architecture
- Migrate to Clean when complexity increases
- Consider DDD for domain-rich applications
- Use Hexagonal for maximum testability needs

### How do I choose between frameworks?

**Web Framework Comparison**:

| Framework | Speed | Use Case | Learning Curve |
|-----------|-------|----------|---------------|
| **Gin** | Fastest | APIs, microservices | Easy |
| **Echo** | Fast | Middleware-rich apps | Moderate |
| **Fiber** | Very Fast | Express-like API | Easy |
| **Chi** | Fast | Stdlib-focused | Easy |

**Recommendation**: Start with Gin (default) unless you have specific requirements.

## Common Issues

### Command not found after installation

**Quick Fix**:
```bash
# Add Go bin to PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Make permanent
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

### Generated project won't compile

**Common Solutions**:
```bash
# 1. Update dependencies
go mod tidy

# 2. Download modules
go mod download

# 3. Clear cache if corrupted
go clean -modcache && go mod download

# 4. Check Go version
go version  # Should be 1.21+
```

### No log output visible

**Quick Fixes**:
```bash
# 1. Set log level
LOG_LEVEL=debug go run cmd/server/main.go

# 2. Check configuration
cat configs/config.yaml

# 3. Verify logger initialization
grep -r "logger.New" internal/
```

### Docker build fails

**Common Solutions**:
```bash
# 1. Check Dockerfile
cat Dockerfile

# 2. Build with verbose output
docker build -t my-app . --no-cache --progress=plain

# 3. Test dependencies
go mod tidy && go build ./...
```

### Port already in use

**Quick Fix**:
```bash
# Find and kill process
kill -9 $(lsof -ti:8080)

# Or use different port
PORT=8081 go run cmd/server/main.go
```

For more detailed troubleshooting, see our [comprehensive troubleshooting guide](troubleshooting.md).

## Getting Help

### Where can I get support?

**Official Channels**:
- **GitHub Issues**: [Bug reports and feature requests](https://github.com/francknouama/go-starter/issues)
- **GitHub Discussions**: [Community support and questions](https://github.com/francknouama/go-starter/discussions)
- **Documentation**: [Comprehensive guides](../01-getting-started/)

### How do I report a bug?

**Include This Information**:
1. **Version**: `go-starter version`
2. **Go Version**: `go version`
3. **OS**: `uname -a` (Linux/Mac) or system info (Windows)
4. **Command**: Exact command that failed
5. **Error**: Complete error message
6. **Expected**: What you expected to happen

### How do I request a feature?

**Feature Request Guidelines**:
1. Search existing issues first
2. Describe the use case clearly
3. Explain why it would benefit others
4. Provide examples if possible
5. Consider contributing the feature yourself

### Can I contribute to go-starter?

**✅ Yes! We welcome contributions**:

- **Bug Fixes**: Always appreciated
- **Documentation**: Help improve clarity
- **Features**: Discuss in issues first
- **Blueprints**: New project types and patterns
- **Testing**: Improve test coverage

See our [contributing guidelines](https://github.com/francknouama/go-starter/blob/main/CONTRIBUTING.md) for details.

### How do I stay updated?

**Stay Informed**:
- ⭐ **Star** the [GitHub repository](https://github.com/francknouama/go-starter)
- 👀 **Watch** for releases and updates
- 💬 **Follow** discussions for community insights
- 📚 **Check** documentation for new features

---

**Can't find your answer?** 
- Check the [troubleshooting guide](troubleshooting.md) for technical issues
- Browse [GitHub Discussions](https://github.com/francknouama/go-starter/discussions) for community help
- [Create an issue](https://github.com/francknouama/go-starter/issues/new) for bugs or feature requests

**Happy coding with go-starter!** 🚀