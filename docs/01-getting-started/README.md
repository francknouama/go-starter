# 🌟 Getting Started with go-starter

Welcome to **go-starter** - the comprehensive Go project generator with progressive disclosure for all skill levels!

## 🚀 Quick Start

### Installation
```bash
# Install go-starter
go install github.com/francknouama/go-starter@latest

# Verify installation
go-starter --version
```

### Generate Your First Project

#### Interactive Mode (Recommended)
```bash
# Start interactive guided project creation
go-starter new

# The wizard will guide you through:
# ✅ Project name (with random generator option)
# ✅ Blueprint selection (12 templates)
# ✅ Complexity level (Simple → Expert)
# ✅ Framework choice (context-aware)
# ✅ Logger selection (slog, zap, logrus, zerolog)
```

#### Direct Mode
```bash
# Skip interactive prompts with flags
go-starter new my-api --type web-api --logger zap

# Generate with specific complexity
go-starter new my-tool --type cli --complexity simple
```

## 📖 Documentation Index

### 💻 Core Documentation
| Document | Purpose | Time | Audience |
|----------|---------|------|----------|
| [Installation](installation.md) | Install go-starter (all platforms) | 5 min | First-time users |
| [Quick Start](quick-start.md) | Generate your first project | 5 min | All users |
| [Getting Started](getting-started.md) | Complete tutorial with examples | 30 min | Developers |

### 🎯 Blueprint & Configuration
| Document | Purpose | Time | Audience |
|----------|---------|------|----------|
| [Blueprint Selection](../02-user-guides/blueprint-selection.md) | Choosing the right template | 10 min | All users |
| [Configuration Guide](../02-user-guides/configuration.md) | Advanced configuration options | 15 min | Power users |
| [Progressive Disclosure](../02-user-guides/progressive-disclosure.md) | Understanding complexity levels | 10 min | All users |

### 📊 Reference & Support
| Document | Purpose | Time | Audience |
|----------|---------|------|----------|
| [CLI Reference](../03-reference/README.md) | All commands and flags | Variable | Reference |
| [FAQ](../02-user-guides/faq.md) | Frequently asked questions | 5 min | Troubleshooting |
| [Troubleshooting](../02-user-guides/troubleshooting.md) | Problem-solving guide | Variable | Support |

## 🎯 Progressive Disclosure System

go-starter adapts to your experience level:

### Basic Mode (Default)
```bash
go-starter new              # Shows 14 essential flags
go-starter new --help       # Basic help with common options
```

### Advanced Mode
```bash
go-starter new --advanced            # All 18+ flags available
go-starter new --advanced --help    # Complete help documentation
```

## 🏗️ Available Blueprints

### Quick Selection by Complexity
```bash
# Simple projects (8-15 files) - Learning & Prototypes
go-starter new my-tool --type cli --complexity simple
go-starter new my-func --type lambda --complexity simple

# Standard projects (25-30 files) - Production Ready
go-starter new my-api --type web-api --complexity standard
go-starter new my-cli --type cli --complexity standard

# Advanced projects (40+ files) - Enterprise Scale
go-starter new my-service --type microservice --complexity advanced
go-starter new my-app --type monolith --complexity advanced
```

### Complete Blueprint List
- **CLI-Simple** (8 files) - Quick utilities and learning
- **CLI-Standard** (29 files) - Production CLI applications
- **Library-Standard** (10 files) - Reusable Go packages
- **Web-API-Standard** (25 files) - REST APIs with middleware
- **Web-API-Clean** (40 files) - Clean architecture APIs
- **Web-API-Echo** (25 files) - Echo framework APIs
- **Web-API-Fiber** (25 files) - Fiber framework APIs
- **Lambda-Standard** (15 files) - AWS serverless functions
- **Lambda-Proxy** (20 files) - API Gateway integration
- **gRPC-Gateway** (45 files) - Dual HTTP/gRPC APIs
- **Monolith** (72 files) - Full-stack applications
- **Microservice** (47 files) - Enterprise gRPC services

## 💡 Success Tips

### For Beginners
1. Start with interactive mode - let the wizard guide you
2. Choose `simple` complexity for learning projects
3. Use the default logger (slog) and framework options
4. Read the generated README.md in your project

### For Power Users
1. Use `--advanced` flag to access all options
2. Create aliases for common project types
3. Explore different complexity levels for the same blueprint
4. Customize with configuration files

### Best Practices
- Always run `make test` after generation to verify setup
- Use `--dry-run` to preview what will be generated
- Check `go.mod` is initialized correctly
- Review the generated Makefile for available commands

## 🚀 Next Steps

After generating your project:

1. **Navigate to your project**
   ```bash
   cd my-project
   ```

2. **Run tests to verify setup**
   ```bash
   make test
   ```

3. **Start development**
   ```bash
   make run    # For applications
   make build  # For libraries
   ```

4. **Explore generated structure**
   - Review README.md for project-specific documentation
   - Check Makefile for available commands
   - Examine the architecture in your chosen blueprint

## 🤝 Getting Help

- **[FAQ](../02-user-guides/faq.md)** - Quick answers to common questions
- **[Troubleshooting](../02-user-guides/troubleshooting.md)** - Solve common problems
- **[GitHub Issues](https://github.com/francknouama/go-starter/issues)** - Report bugs or request features
- **[GitHub Discussions](https://github.com/francknouama/go-starter/discussions)** - Community support

---

Ready to start? → [Installation](installation.md) → [Quick Start](quick-start.md)