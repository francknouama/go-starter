# go-starter Documentation

Welcome to the comprehensive go-starter documentation. Find everything you need to master Go project generation.

## 🚀 Getting Started

**New to go-starter?** Start here:

- 📖 **[Getting Started Guide](GETTING_STARTED.md)** - Your first project in 5 minutes
- ⚙️ **[Installation Guide](INSTALLATION.md)** - All installation methods
- 🏃‍♂️ **[Quick Reference](QUICK_REFERENCE_CARD.md)** - Common commands and patterns

## 📚 Core Guides

### Project Creation
- 📖 **[Project Types Guide](PROJECT_TYPES.md)** - Choose the right template (Web API, CLI, Library, Lambda)
- 🏗️ **[Blueprint Guide](BLUEPRINTS.md)** - Deep dive into project templates
- 📊 **[Blueprint Comparison](BLUEPRINT_COMPARISON.md)** - Side-by-side feature comparison

### Configuration & Customization
- ⚙️ **[Configuration Guide](CONFIGURATION.md)** - Global settings and profiles
- 📊 **[Logger Guide](LOGGER_GUIDE.md)** - Master the unique logger selector system
- 🗃️ **[ORM Guide](ORM_GUIDE.md)** - Database and ORM selection

## 🔧 Advanced Topics

### Development & Testing
- 🧪 **[Testing Guide](TESTING_GUIDE.md)** - Test your generated projects
- 📝 **[TDD Enforcement](TDD-ENFORCEMENT.md)** - Test-driven development practices

### Workflow & Productivity
- 📋 **[Task Master Guide](TASK_MASTER_GUIDE.md)** - Organize development with AI-powered task management
- 📝 **[Quick Reference](QUICK_REFERENCE.md)** - Command cheatsheet

## 🆘 Help & Support

- ❓ **[FAQ](FAQ.md)** - Frequently asked questions
- 🔧 **[Troubleshooting Guide](TROUBLESHOOTING.md)** - Solve common issues
- 🐛 **[Report Issues](https://github.com/francknouama/go-starter/issues)** - Found a bug?
- 💬 **[Discussions](https://github.com/francknouama/go-starter/discussions)** - Community support

## 🚀 Quick Start

```bash
# Install go-starter (using Go install - recommended)
go install github.com/francknouama/go-starter@latest

# Alternative: Download binary from GitHub releases
# Homebrew currently unavailable due to PAT issues

# Generate your first project
go-starter new my-awesome-api

# Navigate and run
cd my-awesome-api
make run
```

## 📊 Available Templates

| Template | Use Case | Key Features |
|----------|----------|--------------|
| **Web API** | REST services, microservices | Gin framework, middleware, database |
| **CLI** | Command-line tools | Cobra framework, subcommands |
| **Library** | Reusable packages | Clean API, examples, minimal deps |
| **Lambda** | Serverless functions | AWS integration, CloudWatch logging |

## 🪵 Logger Options

| Logger | Performance | Best For |
|--------|-------------|----------|
| **slog** | Good | Standard library choice |
| **zap** | Excellent | High-performance applications |
| **logrus** | Good | Feature-rich requirements |
| **zerolog** | Excellent | JSON-heavy, cloud-native |

## 🎯 Common Use Cases

### Building a High-Performance API
```bash
go-starter new api --type=web-api --logger=zap
```

### Creating a Developer Tool
```bash
go-starter new tool --type=cli --logger=logrus
```

### Developing a Reusable Library
```bash
go-starter new sdk --type=library --logger=slog
```

### Deploying Serverless Functions
```bash
go-starter new function --type=lambda --logger=zerolog
```

## 📖 Documentation Structure

```
docs/
├── README.md                    # This file
├── GETTING_STARTED.md          # Installation and first steps
├── BLUEPRINT_COMPARISON.md      # Choosing the right blueprint
├── BLUEPRINTS.md                # Detailed blueprint documentation
├── LOGGER_GUIDE.md             # Logger selector deep dive
├── ORM_GUIDE.md                # Database interaction patterns
├── QUICK_REFERENCE_CARD.md     # Commands and patterns reference
├── FAQ.md                      # Frequently asked questions
└── TROUBLESHOOTING.md          # Problem-solving guide
```

## 🤝 Contributing

We welcome contributions! Please see:
- [Contributing Guide](../CONTRIBUTING.md)
- [Project Roadmap](../PROJECT_ROADMAP.md)

## 🔗 Links

- **Repository**: [github.com/francknouama/go-starter](https://github.com/francknouama/go-starter)
- **Issues**: [GitHub Issues](https://github.com/francknouama/go-starter/issues)
- **Releases**: [Latest Releases](https://github.com/francknouama/go-starter/releases)

---

**Need help?** Check the [FAQ](FAQ.md) or [Troubleshooting Guide](TROUBLESHOOTING.md) first. If you can't find an answer, please [open an issue](https://github.com/francknouama/go-starter/issues/new).