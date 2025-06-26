# go-starter Documentation

Welcome to the comprehensive documentation for go-starter - the powerful Go project generator with a unique logger selector system.

## 📚 Documentation Overview

### Getting Started
- **[Getting Started Guide](GETTING_STARTED.md)** - Installation, first project, and basic workflows
- **[Quick Reference Card](QUICK_REFERENCE_CARD.md)** - All commands at a glance

### Choosing the Right Template
- **[Template Comparison Guide](TEMPLATE_COMPARISON.md)** - Detailed comparison of all templates
- **[Template Usage Guide](TEMPLATES.md)** - In-depth guide for each template type

### Technical Guides
- **[Logger Guide](LOGGER_GUIDE.md)** - Understanding the logger selector system
- **[ORM Selection Guide](ORM_GUIDE.md)** - Choosing between GORM and raw SQL

### Help & Support
- **[FAQ](FAQ.md)** - Frequently asked questions
- **[Troubleshooting Guide](TROUBLESHOOTING.md)** - Comprehensive problem-solving guide

## 🚀 Quick Start

```bash
# Install go-starter
brew install go-starter

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
├── TEMPLATE_COMPARISON.md      # Choosing the right template
├── TEMPLATES.md                # Detailed template documentation
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