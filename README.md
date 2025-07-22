# go-starter

[![Go Report Card](https://goreportcard.com/badge/github.com/francknouama/go-starter)](https://goreportcard.com/report/github.com/francknouama/go-starter)
[![Go Reference](https://pkg.go.dev/badge/github.com/francknouama/go-starter.svg)](https://pkg.go.dev/github.com/francknouama/go-starter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/francknouama/go-starter)](https://github.com/francknouama/go-starter/releases)

The most comprehensive Go project generator with **12 blueprints**, **progressive disclosure**, and **4 logger options**. Generate production-ready Go projects in 30 seconds that compile immediately and follow enterprise best practices.

## ⚡ See It In Action

### Beginner-Friendly Mode (Default)
```bash
# Install once
go install github.com/francknouama/go-starter@latest

# Generate with simple guided prompts
go-starter new my-api --type web-api --logger zap

cd my-api
make run    # Server running on :8080
make test   # All tests pass ✅
make build  # Production binary ready 🚀
```

### Progressive Disclosure - Start Simple, Scale Smart
```bash
# Simple CLI (8 files) - Perfect for learning
go-starter new my-tool --type cli --complexity simple

# Standard CLI (29 files) - Production-ready
go-starter new my-tool --type cli --complexity standard

# Advanced mode - See all 18+ options
go-starter new --advanced --help
```

**That's it.** No configuration files. No dependency hunting. No project structure decisions. Just working, production-ready code that scales with your needs.

## 🎯 What You Get

✅ **Compiles immediately** - Zero setup, zero errors  
✅ **Production-ready** - Industry best practices built-in  
✅ **Complete tests** - Unit, integration, benchmarks  
✅ **Docker ready** - Dockerfile and docker-compose  
✅ **CI/CD included** - GitHub Actions configured  
✅ **Full documentation** - README, API docs, examples  

## 🚀 12 Complete Blueprints

### 📊 Core Web APIs (4 Architecture Patterns)
| Blueprint | Use Case | Architecture | 
|-----------|----------|--------------|
| **🌐 Standard Web API** | REST APIs, CRUD services | Standard layered |
| **🏗️ Clean Architecture API** | Enterprise applications | Clean Architecture |
| **⚙️ DDD Web API** | Domain-rich applications | Domain-Driven Design |
| **🔩 Hexagonal Architecture API** | Highly testable systems | Ports & Adapters |

### 🖥️ CLI Applications (2 Complexity Levels)
| Blueprint | Use Case | Files | Complexity |
|-----------|----------|--------|------------|
| **📱 Simple CLI** | Scripts, utilities | 8 files | Beginner |
| **⚙️ Standard CLI** | Production tools | 29 files | Professional |

### 🏢 Enterprise & Cloud-Native
| Blueprint | Use Case | Key Features |
|-----------|----------|--------------|
| **🌐 gRPC Gateway** | API Gateway + gRPC | Dual HTTP/gRPC, TLS |
| **🔄 Event-Driven** | CQRS, Event Sourcing | Event streams, projections |
| **🏗️ Microservice** | Service mesh, K8s | Discovery, circuit breakers |
| **🏢 Monolith** | Traditional web apps | Full-stack, templating |

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

**Alternative installation:** [Download binaries](docs/guides/INSTALLATION.md) • [All methods](docs/guides/INSTALLATION.md)

## 📚 Documentation

| Guide | Description |
|-------|-------------|
| 🚀 **[Quick Start](docs/guides/GETTING_STARTED.md)** | Your first project in 5 minutes |
| ⚙️ **[Installation](docs/guides/INSTALLATION.md)** | All installation methods |
| 📖 **[Project Types](docs/references/PROJECT_TYPES.md)** | Choose the right template |
| 📊 **[Logger Guide](docs/references/LOGGER_GUIDE.md)** | Master the logger selector |
| 🔧 **[Configuration](docs/guides/CONFIGURATION.md)** | Customize your setup |
| 📋 **[Complete Docs](docs/README.md)** | Full documentation index |

## 🛣️ Current Status & Roadmap

**Current (v2.0+):** 12 complete blueprints, progressive disclosure, enterprise architecture patterns

### ✅ Phase 2 Complete - Advanced Architecture Patterns
- 🏗️ **Advanced Architectures** - Clean, DDD, Hexagonal ✅
- 🔄 **Event-Driven Architecture** - CQRS, Event Sourcing ✅  
- 🏢 **Enterprise Patterns** - Microservices, Monoliths ✅
- 🌐 **gRPC Gateway** - Dual HTTP/gRPC APIs ✅
- 🔧 **Go Workspace** - Multi-module monorepos ✅

### 🚧 Phase 3 - In Development
- 📱 **Web Interface** - Browser-based project generator
- 🌐 **More Frameworks** - Echo, Fiber, Chi, Bun Router
- 🗃️ **Database Options** - GORM, sqlx, sqlc, ent, Bun ORM
- 📊 **Analytics Databases** - ClickHouse, TimescaleDB

### 🔮 Phase 4 - Future Vision
- 🔍 **Monitoring & APM** - Prometheus, OpenTelemetry, Uptrace
- ☁️ **Cloud Platforms** - AWS, GCP, Azure deployment
- 🏪 **Blueprint Marketplace** - Community templates

## ❤️ Community

- 🐛 **[Report Issues](https://github.com/francknouama/go-starter/issues)** - Found a bug?
- 💬 **[Discussions](https://github.com/francknouama/go-starter/discussions)** - Questions and ideas
- 🤝 **[Contributing](CONTRIBUTING.md)** - Make go-starter better

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