# go-starter

[![Go Report Card](https://goreportcard.com/badge/github.com/francknouama/go-starter)](https://goreportcard.com/report/github.com/francknouama/go-starter)
[![Go Reference](https://pkg.go.dev/badge/github.com/francknouama/go-starter.svg)](https://pkg.go.dev/github.com/francknouama/go-starter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/francknouama/go-starter)](https://github.com/francknouama/go-starter/releases)

The most comprehensive Go project generator with **12 blueprints**, **progressive disclosure**, and **4 logger options**. Generate production-ready Go projects in 30 seconds that compile immediately and follow enterprise best practices.

## ⚡ See It In Action

### 🎯 Interactive Mode (Recommended)
```bash
# Install once
go install github.com/francknouama/go-starter@latest

# Start interactive guided project creation
go-starter new

# The interactive wizard will guide you through:
# ✅ Project name (with random generator option)
# ✅ Blueprint selection (with descriptions)
# ✅ Complexity level (Simple → Expert)
# ✅ Framework choice (context-aware options)
# ✅ Logger selection (4 options)
# ✅ Additional features (database, auth, etc.)
```

### Direct Mode (Non-Interactive)
```bash
# Skip interactive prompts with flags
go-starter new my-api --type web-api --logger zap

cd my-api
make run    # Server running on :8080
make test   # All tests pass ✅
make build  # Production binary ready 🚀
```

### Progressive Disclosure - Start Simple, Scale Smart
```bash
# Interactive mode adapts to your experience level
go-starter new              # Basic options (beginner-friendly)
go-starter new --advanced   # All options (power users)

# Direct mode with complexity control
go-starter new my-tool --type cli --complexity simple    # 8 files
go-starter new my-tool --type cli --complexity standard  # 29 files
```

**That's it.** No configuration files. No dependency hunting. No project structure decisions. Just working, production-ready code that scales with your needs.

## 🎯 What You Get

✅ **Compiles immediately** - Zero setup, zero errors  
✅ **Production-ready** - Industry best practices built-in  
✅ **Complete tests** - Unit, integration, benchmarks  
✅ **Docker ready** - Dockerfile and docker-compose  
✅ **CI/CD included** - GitHub Actions configured  
✅ **Full documentation** - README, API docs, examples  

## 🧪 Enhanced Testing Infrastructure

- **114+ comprehensive test scenarios** validating all blueprint combinations
- **100% compilation guarantee** for all generated projects  
- **Architecture validation** with AST parsing for Clean/DDD/Hexagonal patterns
- **Cross-platform testing** ensuring Windows, macOS, and Linux compatibility
- **Performance monitoring** with resource usage tracking
- **Enhanced ATDD suite** covering enterprise scenarios

## 🚀 12 Enterprise-Grade Blueprints

### 📊 Core Web APIs (4 Architecture Patterns) - Production Hardened
| Blueprint | Use Case | Architecture | Production Features |
|-----------|----------|--------------|--------------------|

| **🌐 Standard Web API** | REST APIs, CRUD services | Standard layered | HTTP middleware, validation, CORS |
| **🏗️ Clean Architecture API** | Enterprise applications | Clean Architecture | Dependency injection, layered testing |
| **⚙️ DDD Web API** | Domain-rich applications | Domain-Driven Design | Rich domain models, event sourcing |
| **🔩 Hexagonal Architecture API** | Highly testable systems | Ports & Adapters | Multiple adapters, complete isolation |

### 🖥️ CLI Applications (2 Complexity Levels)
| Blueprint | Use Case | Files | Complexity |
|-----------|----------|--------|------------|
| **📱 Simple CLI** | Scripts, utilities | 8 files | Beginner |
| **⚙️ Standard CLI** | Production tools | 29 files | Professional |

### 🏢 Enterprise & Cloud-Native - 🚀 Phase 2 Enhanced
| Blueprint | Use Case | Production Features |
|-----------|----------|---------------------|
| **🌐 gRPC Gateway** | API Gateway + gRPC | **✨ Enhanced interceptors, rate limiting, unified metrics** |
| **🔄 Event-Driven** | CQRS, Event Sourcing | Event streams, projections |
| **🏗️ Microservice** | Service mesh, K8s | **✨ OpenTelemetry tracing, resilience patterns, performance optimization** |
| **🏢 Monolith** | Traditional web apps | **✨ Background job processing, multi-layer caching, performance monitoring** |

### ☁️ Serverless & Tools  
| Blueprint | Use Case | Runtime |
|-----------|----------|---------|
| **⚡ AWS Lambda** | Event functions | AWS Lambda Go |
| **🌉 Lambda Proxy** | API Gateway integration | HTTP proxy patterns |
| **📦 Library** | SDKs, packages | Clean API + examples |
| **🔧 Go Workspace** | Monorepo projects | Multi-module workspace |

## 🎛️ Unique Logger Selector

**Choose your logging strategy:**

```bash
go-starter new api --logger zap        # Zero allocations ⚡
go-starter new app --logger slog       # Standard library 📚  
go-starter new service --logger zerolog # JSON optimized ☁️
go-starter new tool --logger logrus    # Feature-rich 🔧
```

**Progressive Complexity Examples:**

```bash
# Start simple, grow as needed
go-starter new my-tool --type cli --complexity simple   # 8 files
go-starter new my-api --type web-api --architecture clean
go-starter new my-service --type microservice --logger zap
go-starter new my-workspace --type workspace   # Multi-module monorepo
```

**Switch anytime** without changing application code.

## 📈 Before vs After

| Before go-starter | After go-starter |
|-------------------|------------------|
| 🕐 2-4 hours setup | ⚡ 30 seconds |
| 🐛 Config bugs | ✅ Works out of the box |
| 📚 Research best practices | 🏆 Best practices by default |
| ⚠️ Missing tests/Docker/CI | 🚀 Everything included |

## 🏃‍♂️ Quick Start

### Basic Mode (Beginner-Friendly)
```bash
# 1. Install
go install github.com/francknouama/go-starter@latest

# 2. Generate with guided prompts (shows 14 essential options)
go-starter new my-project

# 3. Ship
cd my-project && make run
```

### Advanced Mode (Power Users)
```bash
# See all 18+ options for complex projects
go-starter new --advanced --help

# Generate enterprise patterns directly
go-starter new enterprise-api --type web-api --architecture hexagonal --advanced

# Create workspace for monorepos
go-starter new my-workspace --type workspace
```

**Alternative installation:** [Download binaries](docs/01-getting-started/installation.md) • [All methods](docs/01-getting-started/installation.md)

## 📚 Documentation

### 🚀 **Quick Start** (5 minutes to success)
- **[Installation Guide](docs/01-getting-started/installation.md)** - Install on any platform  
- **[Quick Start](docs/01-getting-started/quick-start.md)** - Generate your first project in 5 minutes
- **[Getting Started](docs/01-getting-started/getting-started.md)** - Complete tutorial with examples

### 📖 **User Guides** (Real-world usage)
- **[Blueprint Selection](docs/02-user-guides/blueprint-selection.md)** - Choose the perfect project type & architecture
- **[Configuration Guide](docs/02-user-guides/configuration.md)** - Team setup & shared settings  
- **[Troubleshooting](docs/02-user-guides/troubleshooting.md)** - Solve problems quickly
- **[FAQ](docs/02-user-guides/faq.md)** - Quick answers to common questions

### 📋 **Complete Documentation**
- **[📚 Full Documentation Hub](docs/README.md)** - All guides organized by user type
- **[🔧 Developer Guide](docs/04-developers/README.md)** - Contributing and development setup
- **[🌟 Community Resources](docs/05-community/README.md)** - Examples and showcases

## 🛣️ Current Status & Roadmap

**Current (v2.1+):** 12 production-ready blueprints, enterprise observability, advanced resilience patterns, comprehensive security features

### ✅ Phase 2 Enterprise Production Features - 🚀 Major Enhancements Complete
- 🏗️ **Advanced Architectures** - Clean, DDD, Hexagonal ✅
- 🔄 **Event-Driven Architecture** - CQRS, Event Sourcing ✅  
- 🏢 **Enterprise Microservices** - ✨ **OpenTelemetry tracing, rate limiting, resilience patterns** ✅
- 🌐 **Enhanced gRPC Gateway** - ✨ **Advanced interceptors, unified middleware, metrics collection** ✅
- 🏢 **Production Monoliths** - ✨ **Background jobs, multi-layer caching, performance monitoring** ✅
- 🔧 **Go Workspace** - Multi-module monorepos ✅
- 🧪 **Enhanced ATDD Testing** - 114+ scenarios, 100% blueprint validation ✅
- 🔍 **Enterprise Observability** - ✨ **OpenTelemetry, Prometheus metrics, distributed tracing** ✅
- 🔒 **Security & Resilience** - ✨ **Input validation, circuit breakers, graceful error handling** ✅

### 🚧 Phase 3 - Web Interface Development
- 📱 **Web Interface** - Browser-based project generator
- 🌐 **More Frameworks** - Echo, Fiber, Chi, Bun Router
- 🗃️ **Database Options** - GORM, sqlx, sqlc, ent, Bun ORM
- 📊 **Analytics Databases** - ClickHouse, TimescaleDB

### 🔮 Phase 4 - Future Vision
- 🔍 **Monitoring & APM** - Prometheus, OpenTelemetry, Uptrace
- ☁️ **Cloud Platforms** - AWS, GCP, Azure deployment
- 🏪 **Blueprint Marketplace** - Community templates

## ❤️ Community & Support

### 🆘 **Getting Help**
- **[FAQ](docs/02-user-guides/faq.md)** - Quick answers to common questions
- **[Troubleshooting](docs/02-user-guides/troubleshooting.md)** - Solve problems with proven solutions
- **[GitHub Discussions](https://github.com/francknouama/go-starter/discussions)** - Community Q&A and ideas

### 🤝 **Contributing**
- **[Development Guide](docs/04-developers/README.md)** - Set up your development environment
- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute code and documentation
- **[Report Issues](https://github.com/francknouama/go-starter/issues)** - Bug reports and feature requests

### 🌟 **Community Resources**
- **[Community Hub](docs/05-community/README.md)** - Examples, showcases, and best practices
- **[GitHub Discussions](https://github.com/francknouama/go-starter/discussions)** - Share your projects and get help
- **[Project Examples](docs/guides/README.md)** - Sample projects and proven patterns

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

**Ready to experience the most comprehensive Go generator?**

```bash
# Beginner? Start here
go install github.com/francknouama/go-starter@latest
go-starter new my-first-project

# Power user? Go advanced
go-starter new enterprise-system --type microservice --architecture hexagonal --advanced
```

🚀 **From simple scripts to enterprise architectures - go-starter scales with you.**

⭐ **[Star us on GitHub](https://github.com/francknouama/go-starter)** • 🐛 **[Report Issues](https://github.com/francknouama/go-starter/issues)** • 💬 **[Join Discussions](https://github.com/francknouama/go-starter/discussions)**