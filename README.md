# go-starter

[![Go Report Card](https://goreportcard.com/badge/github.com/francknouama/go-starter)](https://goreportcard.com/report/github.com/francknouama/go-starter)
[![Go Reference](https://pkg.go.dev/badge/github.com/francknouama/go-starter.svg)](https://pkg.go.dev/github.com/francknouama/go-starter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/v/release/francknouama/go-starter)](https://github.com/francknouama/go-starter/releases)

Stop fighting boilerplate. Generate complete, production-ready Go projects that compile immediately and follow best practices.

## ⚡ See It In Action

```bash
# Install once
go install github.com/francknouama/go-starter@latest

# Create a production-ready API in 30 seconds
go-starter new my-api --type web-api --logger zap

cd my-api
make run    # Server running on :8080
make test   # All tests pass ✅
make build  # Production binary ready 🚀
```

**That's it.** No configuration files. No dependency hunting. No project structure decisions. Just working, production-ready code.

## 🎯 What You Get Instantly

Every generated project includes:

- ✅ **Compiles immediately** - Zero setup, zero errors
- ✅ **Production-ready structure** - Industry best practices built-in
- ✅ **Complete test suite** - Unit tests, integration tests, benchmarks
- ✅ **Docker ready** - Dockerfile and docker-compose included
- ✅ **CI/CD pipeline** - GitHub Actions workflow configured
- ✅ **Documentation** - README, API docs, examples included

## 🚀 Choose Your Path

### 🌐 REST APIs & Web Services
```bash
go-starter new user-service --type web-api
```
**Perfect for:** Microservices, REST APIs, web backends  
**Includes:** Gin framework, database integration, middleware, health checks

### 🖥️ CLI Tools & Automation
```bash
go-starter new deploy-tool --type cli  
```
**Perfect for:** DevOps tools, automation scripts, utilities  
**Includes:** Cobra framework, subcommands, configuration, shell completion

### 📦 Libraries & SDKs
```bash
go-starter new awesome-sdk --type library
```
**Perfect for:** Reusable packages, SDKs, shared components  
**Includes:** Clean API design, examples, benchmarks, documentation

### ⚡ Serverless Functions
```bash
go-starter new processor --type lambda
```
**Perfect for:** AWS Lambda, event processing, serverless APIs  
**Includes:** Lambda runtime, SAM templates, deployment automation

## 🎛️ Unique Logger Selector

**The only Go generator with pluggable logging.** Choose your logging strategy, not your vendor:

```bash
# High-performance APIs
go-starter new api --logger zap        # Zero allocations ⚡

# Cloud-native services  
go-starter new service --logger zerolog # JSON optimized ☁️

# Standard library approach
go-starter new app --logger slog       # Go 1.21+ built-in 📚

# Feature-rich applications
go-starter new platform --logger logrus # Hooks & formatters 🔧
```

**Switch anytime** without changing a single line of application code. Same interface, different performance characteristics.

## 📈 Real Results

**Before go-starter:**
- 🕐 2-4 hours setting up project structure
- 🐛 Configuration bugs and dependency conflicts  
- 📚 Reading docs for project layout best practices
- ⚠️ Missing tests, Docker configs, or CI/CD

**After go-starter:**
- ⚡ 30 seconds to working project
- ✅ Everything works out of the box
- 🏆 Industry best practices by default
- 🚀 Focus on business logic, not boilerplate

## 🏃‍♂️ Quick Start

### 1. Install
```bash
go install github.com/francknouama/go-starter@latest
```

### 2. Generate
```bash
go-starter new my-project
# Follow the interactive prompts or use direct mode
```

### 3. Ship
```bash
cd my-project
make run    # Development server
make test   # Run tests  
make build  # Production binary
make docker # Container image
```

## 📚 Learn More

- 🚀 **[Quick Start Guide](docs/GETTING_STARTED.md)** - Your first project in 5 minutes
- 📖 **[Project Types](docs/PROJECT_TYPES.md)** - Choose the right template  
- 📊 **[Logger Guide](docs/LOGGER_GUIDE.md)** - Master the logger selector
- ⚙️ **[Installation](docs/INSTALLATION.md)** - All installation methods
- 🔧 **[Configuration](docs/CONFIGURATION.md)** - Customize your setup

## 🛣️ What's Next

**Current (v1.3.1):** 4 project types, 4 logger options, bulletproof basics

**Coming Soon:**
- 🏗️ **Advanced Architectures** - Clean Architecture, DDD, Hexagonal patterns
- 🌐 **More Frameworks** - Echo, Fiber, Chi web frameworks  
- 🗃️ **Database Options** - GORM, sqlx, sqlc, ent ORMs
- 📱 **Web Interface** - Browser-based project generator

See [PROJECT_ROADMAP.md](PROJECT_ROADMAP.md) for detailed plans.

## ❤️ Community

- 🐛 **[Report Issues](https://github.com/francknouama/go-starter/issues)** - Found a bug?
- 💬 **[Discussions](https://github.com/francknouama/go-starter/discussions)** - Questions and ideas
- 🤝 **[Contributing](CONTRIBUTING.md)** - Make go-starter better

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

**Ready to stop fighting boilerplate?**

```bash
go install github.com/francknouama/go-starter@latest
go-starter new my-project
```