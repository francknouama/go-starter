# Quick Start Guide

Get your first go-starter project running in 5 minutes! ⚡

> **🎉 HISTORIC 100% COVERAGE ACHIEVED**: go-starter now has **12 production-ready blueprints**, including enterprise-enhanced **Monolith** and **Microservice** with advanced observability and resilience patterns!

## 🎯 Goal

By the end of this guide, you'll have:
- ✅ go-starter installed
- ✅ Your first project generated  
- ✅ A working Go application running

## 📦 Step 1: Install go-starter

### macOS (Homebrew)
```bash
brew install francknouama/tap/go-starter
```

### Linux/Windows (Binary)
```bash
# Download latest release
curl -L https://github.com/francknouama/go-starter/releases/latest/download/go-starter-linux-amd64 -o go-starter
chmod +x go-starter
sudo mv go-starter /usr/local/bin/
```

### Go Install
```bash
go install github.com/francknouama/go-starter@latest
```

### Verify Installation
```bash
go-starter version
# Should output: go-starter v2.x.x
```

## 🚀 Step 2: Generate Your First Project

### Simple CLI Tool (Recommended for beginners)
```bash
go-starter new my-first-cli --type=cli --complexity=simple
```

This creates a minimal CLI tool with:
- ✅ 8 files (not overwhelming!)
- ✅ Basic command structure
- ✅ Built-in help system
- ✅ Ready to run

### Alternative: Web API (For backend developers)
```bash
go-starter new my-api --type=web-api --framework=gin
```

### New: gRPC Gateway ✅ **NEWLY PRODUCTION READY** (For dual-protocol APIs)
```bash
go-starter new my-gateway --type=grpc-gateway
```

This creates a production-ready API gateway with:
- ✅ Both HTTP REST and gRPC endpoints
- ✅ Enhanced interceptors and middleware
- ✅ Protocol buffer integration

## 🎮 Step 3: Run Your Project

```bash
cd my-first-cli  # or my-api

# Build and run
go run main.go

# For CLI: Try the help command
go run main.go --help

# For API: Server starts on http://localhost:8080
```

## 🎉 Success!

You now have a working Go project! Here's what you can do next:

### 📚 Learn More
- **[Complete Tutorial](getting-started.md)** - Detailed walkthrough with examples
- **[Blueprint Guide](../02-user-guides/blueprint-selection.md)** - Choose the right project type
- **[Configuration](../02-user-guides/configuration.md)** - Customize for your needs

### 🛠️ Customize Your Project
```bash
# Add features to your CLI
go-starter new enhanced-cli --type=cli --complexity=standard --logger=zap

# Try different web frameworks
go-starter new fiber-api --type=web-api --framework=fiber

# Explore architectures
go-starter new clean-api --type=web-api --architecture=clean

# Try the new gRPC Gateway
go-starter new gateway-api --type=grpc-gateway
```

### 🔍 Explore Available Options
```bash
# See all project types
go-starter list

# Get help with advanced options
go-starter new --advanced --help

# Preview without creating files
go-starter new test-project --type=cli --dry-run
```

## 🆘 Need Help?

- **Common Issues**: [Troubleshooting Guide](../02-user-guides/troubleshooting.md)
- **Quick Answers**: [FAQ](../02-user-guides/faq.md)
- **Full Reference**: [CLI Commands](../03-reference/cli-commands.md)

## 🎯 What's Next?

### For Beginners
Continue with the [Complete Getting Started Guide](getting-started.md) for a deeper understanding.

### For Experienced Developers
Jump to [Blueprint Selection](../02-user-guides/blueprint-selection.md) to explore all project types and architectures.

### For Teams
Check out [Configuration Guide](../02-user-guides/configuration.md) for team standards and shared configurations.

---

**🚀 Ready to build something amazing?** You've got all the tools you need to start your Go project journey!