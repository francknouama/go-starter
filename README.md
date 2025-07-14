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

## 🎯 What You Get

✅ **Compiles immediately** - Zero setup, zero errors  
✅ **Production-ready** - Industry best practices built-in  
✅ **Complete tests** - Unit, integration, benchmarks  
✅ **Docker ready** - Dockerfile and docker-compose  
✅ **CI/CD included** - GitHub Actions configured  
✅ **Full documentation** - README, API docs, examples  

## 🚀 Four Project Types

| Type | Use Case | Framework | 
|------|----------|-----------|
| **🌐 Web API** | REST APIs, microservices | Gin + database |
| **🖥️ CLI Tool** | DevOps, automation | Cobra + subcommands |
| **📦 Library** | SDKs, packages | Clean API + examples |
| **⚡ Lambda** | Serverless, events | AWS runtime + SAM |

## 🎛️ Unique Logger Selector

**Choose your logging strategy:**

```bash
go-starter new api --logger zap        # Zero allocations ⚡
go-starter new app --logger slog       # Standard library 📚  
go-starter new service --logger zerolog # JSON optimized ☁️
go-starter new tool --logger logrus    # Feature-rich 🔧
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

```bash
# 1. Install
go install github.com/francknouama/go-starter@latest

# 2. Generate (interactive mode)
go-starter new my-project

# 3. Ship
cd my-project && make run
```

**Alternative installation:** [Download binaries](docs/INSTALLATION.md) • [All methods](docs/INSTALLATION.md)

## 📚 Documentation

| Guide | Description |
|-------|-------------|
| 🚀 **[Quick Start](docs/GETTING_STARTED.md)** | Your first project in 5 minutes |
| ⚙️ **[Installation](docs/INSTALLATION.md)** | All installation methods |
| 📖 **[Project Types](docs/PROJECT_TYPES.md)** | Choose the right template |
| 📊 **[Logger Guide](docs/LOGGER_GUIDE.md)** | Master the logger selector |
| 🔧 **[Configuration](docs/CONFIGURATION.md)** | Customize your setup |
| 📋 **[Complete Docs](docs/README.md)** | Full documentation index |

## 🛣️ What's Next

**Current (v1.3.1):** 4 project types, 4 logger options, production-ready code

**Coming Soon:**
- 🏗️ **Advanced Architectures** - Clean, DDD, Hexagonal patterns  
- 🌐 **More Frameworks** - Echo, Fiber, Chi, Bun Router
- 🗃️ **Database Options** - GORM, sqlx, sqlc, ent, Bun ORM
- 📊 **Analytics Databases** - ClickHouse, TimescaleDB
- 🔍 **Monitoring & APM** - Prometheus, OpenTelemetry, Uptrace
- 📱 **Web Interface** - Browser-based project generator

## ❤️ Community

- 🐛 **[Report Issues](https://github.com/francknouama/go-starter/issues)** - Found a bug?
- 💬 **[Discussions](https://github.com/francknouama/go-starter/discussions)** - Questions and ideas
- 🤝 **[Contributing](CONTRIBUTING.md)** - Make go-starter better

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

**Ready to 10x your Go development?**

```bash
go install github.com/francknouama/go-starter@latest
go-starter new my-project
```

⭐ **[Star us on GitHub](https://github.com/francknouama/go-starter)** • 🐛 **[Report Issues](https://github.com/francknouama/go-starter/issues)** • 💬 **[Join Discussions](https://github.com/francknouama/go-starter/discussions)**