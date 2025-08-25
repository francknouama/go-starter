# go-starter

[![Go Report Card](https://goreportcard.com/badge/github.com/francknouama/go-starter)](https://goreportcard.com/report/github.com/francknouama/go-starter)
[![Go Reference](https://pkg.go.dev/badge/github.com/francknouama/go-starter.svg)](https://pkg.go.dev/github.com/francknouama/go-starter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/francknouama/go-starter)](https://github.com/francknouama/go-starter/releases)

The most comprehensive Go project generator with **7 production-ready blueprints** ✅ **MILESTONE ACHIEVED!** (3 enhancement-ready, 10 more in development), **progressive disclosure system**, and **simplified logger architecture**. Generate production-ready Go projects in 30 seconds that compile immediately and follow enterprise best practices.

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

> **ATDD Validation Complete**: Comprehensive testing framework implemented and validated

✅ **Production-Ready Blueprints (7)** 🎉 **MILESTONE ACHIEVED**
- **CLI-Simple** (10 files) • **CLI-Standard** (28 files)
- **Web-API-Standard** (44 files) • **Lambda-Standard** (17 files)  
- **Library-Standard** (19 files) • **Web-API-Clean** (69 files)
- **gRPC-Gateway** (45 files) ✅ **NEWLY PRODUCTION READY**

🔧 **Enhancement-Ready Blueprints (3)**
- **Lambda-Event-Processing** (22 files, template fixes needed)
- **Monolith** (66 files, module resolution fixes needed)
- **Microservice-Standard** (template processing fixes needed)

🔄 **Development Phase Blueprints (10)**
- **Event-Driven** (58 missing files identified)

**ATDD Framework Achievements:**
- **Sprig Template Functions**: Fixed 296+ template errors
- **Cross-Platform Testing**: Windows, macOS, Linux validated
- **Progressive Disclosure**: 66.7% file reduction (CLI-Simple vs CLI-Standard)
- **Logger Integration**: All 4 logger types tested and working

## 🚀 Blueprint Production Status

> **Current Status**: 7 production-ready, 3 enhancement-ready, 10 in development pipeline

### ✅ Production-Ready Blueprints (Immediate Use)

**Comprehensive ATDD validation confirms these blueprints are ready for production:**

| Blueprint | Type | Files | Status | Key Features |
|-----------|------|-------|--------|--------------|
| **CLI-Simple** | CLI Tool | **10** | ✅ **VALIDATED** | Minimal structure, perfect for learning |
| **CLI-Standard** | CLI Tool | **28** | ✅ **VALIDATED** | Cobra framework, production-ready |
| **Web-API-Standard** | REST API | **44** | ✅ **VALIDATED** | Complete CRUD, middleware, testing |
| **Lambda-Standard** | Serverless | **17** | ✅ **VALIDATED** | AWS SDK v2, CloudWatch, X-Ray |
| **Library-Standard** | Go Package | **19** | ✅ **VALIDATED** | Clean API, examples, documentation |
| **Web-API-Clean** | Clean Architecture | **69** | ✅ **VALIDATED** | Enterprise patterns, Clean Architecture |
| **gRPC-Gateway** | Microservice | **45** | ✅ **VALIDATED** | Dual HTTP/gRPC APIs, gateway pattern |

### 🔧 Enhancement-Ready Blueprints (Phase 2B)

**Files generate successfully, template fixes needed for production readiness:**

| Priority | Blueprint | Files | Status | Fix Required |
|----------|-----------|-------|--------|--------------|
| **P1** | **Lambda-Event-Processing** | 22 | Template syntax issues | Fix string escaping in tracing.go |
| **P1** | **Monolith** | 66 | Module resolution | Fix import paths and package structure |
| **P2** | **Microservice-Standard** | 47 | Template processing hangs | Fix conditional template logic |

### 🔄 Development Pipeline (Phase 3 Priorities)

**Major development effort required - comprehensive implementation needed:**

| Priority | Blueprint | Complexity | Status | Development Effort |
|----------|-----------|------------|--------|--------------------|
| **P1** | **Event-Driven** | Expert | 58 missing files | Major architecture implementation |
| **P2** | **Go Workspace** | Advanced | Multi-module structure | Monorepo tooling |

**Future roadmap includes 8 additional specialized blueprints for advanced use cases.**

## 🎛️ Simplified Logger Architecture ✨

> **Major Enhancement**: 60-90% code reduction while maintaining all logger functionality

**Choose your logging strategy with simplified implementation:**

```bash
go-starter new api --logger zap        # Zero allocations ⚡ (91% less code)
go-starter new app --logger slog       # Standard library 📚 (default)
go-starter new service --logger zerolog # JSON optimized ☁️ (72% less code)
go-starter new tool --logger logrus    # Feature-rich 🔧 (simplified interface)
```

**Logger Complexity Reduction:**
- **CLI Standard**: 1,051 → 98 lines (91% reduction)
- **Web API Standard**: 398 → 110 lines (72% reduction)  
- **Single Interface**: Consistent API across all loggers
- **Conditional Dependencies**: Only selected logger included

**Progressive Complexity Examples:**

```bash
# Start simple, grow as needed - VALIDATED file counts
go-starter new my-tool --type cli --complexity simple   # 10 files ✅
go-starter new my-tool --type cli --complexity standard # 28 files ✅
go-starter new my-api --type web-api                    # 44 files ✅
go-starter new my-lambda --type lambda                  # 17 files ✅
go-starter new my-lib --type library                    # 19 files ✅
go-starter new my-clean --type web-api --architecture clean # 69 files ✅
go-starter new my-gateway --type grpc-gateway           # 45 files ✅
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

**Current (v2.1+):** ✅ **7 production-ready blueprints (MILESTONE ACHIEVED!)**, progressive disclosure system, simplified logger architecture, comprehensive ATDD testing framework, proven hook-enhanced agent workflows

### ✅ Phase 2.1 Core Infrastructure - 🚀 COMPLETED
- ✅ **Progressive Disclosure System** - Smart help with basic/advanced modes
- ✅ **ATDD Testing Framework** - Comprehensive blueprint validation  
- ✅ **Template Function Fixes** - Resolved 296+ Sprig template errors
- ✅ **CLI Two-Tier Approach** - Simple (10 files) vs Standard (28 files)
- ✅ **Simplified Logger Architecture** - 60-90% code reduction
- ✅ **Production Blueprint Validation** - 7 blueprints fully verified
- ✅ **Cross-Platform Testing** - Windows, macOS, Linux compatibility

### ✅ Phase 2.2 - Enhancement Sprint (COMPLETED)
- ✅ **Web-API-Clean** - Now production-ready with Clean Architecture patterns
- ✅ **gRPC-Gateway** - Now production-ready with dual HTTP/gRPC support ✅ **MILESTONE ACHIEVED**
- ✅ **Blueprint Metrics** - File count accuracy validated and updated
- ✅ **Hook-Enhanced Agent Workflow** - Multi-specialist coordination proved highly effective

### 🔄 Phase 3 - Major Blueprint Development
- 🔄 **Event-Driven Architecture** - 58 missing files implementation
- 🏗️ **Microservice-Standard** - Service mesh and observability patterns
- 🏢 **Monolith** - Modular architecture with background job processing
- 🔧 **Go Workspace** - Multi-module monorepo tooling

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