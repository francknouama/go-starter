# Project Types Guide

## Available Templates

### 🌐 Web API
**Perfect for**: REST APIs, microservices, web backends

**What you get**:
- Production-ready Gin framework setup
- Database integration (PostgreSQL, MySQL, SQLite, etc.)
- Middleware for CORS, logging, authentication
- Health check endpoints
- Request/response validation
- Comprehensive test suite
- Docker configuration
- CI/CD pipeline

**Generated structure**:
```
my-api/
├── cmd/server/           # Application entry point
├── internal/            # Private application code
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Custom middleware
│   ├── models/          # Data models
│   └── services/        # Business logic
├── pkg/                 # Public packages
├── tests/               # Integration tests
├── docker-compose.yml   # Local development
├── Dockerfile          # Production container
└── Makefile            # Development commands
```

**Example usage**:
```bash
go-starter new user-api --type web-api --logger zap --database postgres
cd user-api && make run
```

### 🖥️ CLI Application
**Perfect for**: DevOps tools, utilities, automation scripts

**What you get**:
- Cobra framework with subcommands
- Configuration management with Viper
- Interactive prompts and validation
- Shell completion (bash, zsh, fish)
- Version management and updates
- Comprehensive help system
- Cross-platform builds
- Release automation

**Generated structure**:
```
my-cli/
├── cmd/                 # Commands and subcommands
│   ├── root.go         # Root command
│   ├── version.go      # Version command
│   └── deploy.go       # Example subcommand
├── internal/           # Internal packages
│   ├── config/         # Configuration
│   └── utils/          # Utilities
├── pkg/                # Public packages
├── completion/         # Shell completion scripts
└── scripts/           # Build and release scripts
```

**Example usage**:
```bash
go-starter new deploy-tool --type cli --logger logrus
cd deploy-tool
go run main.go --help
./deploy-tool completion bash > /etc/bash_completion.d/deploy-tool
```

### 📦 Go Library
**Perfect for**: SDKs, shared packages, reusable components

**What you get**:
- Clean public API design
- Comprehensive documentation with examples
- Benchmark tests and performance tracking
- Semantic versioning setup
- GitHub releases automation
- Go modules best practices
- Badge generation
- Example applications

**Generated structure**:
```
awesome-sdk/
├── pkg/                # Public API
│   └── awesome/        # Main package
├── internal/           # Private implementation
├── examples/           # Usage examples
├── docs/               # Documentation
├── benchmarks/         # Performance tests
└── .goreleaser.yml     # Release configuration
```

**Example usage**:
```bash
go-starter new awesome-sdk --type library --logger slog
cd awesome-sdk
make test && make benchmark
```

### ⚡ AWS Lambda
**Perfect for**: Serverless functions, event processing, APIs

**What you get**:
- AWS Lambda Go runtime setup
- API Gateway integration
- CloudWatch logging optimization
- SAM template for infrastructure
- Local development with SAM CLI
- Automated deployment scripts
- Environment-specific configurations
- Cold start optimization

**Generated structure**:
```
data-processor/
├── cmd/lambda/         # Lambda function code
├── internal/           # Business logic
├── template.yaml       # SAM infrastructure template
├── events/             # Test events
├── scripts/            # Deployment scripts
└── Makefile           # Lambda-specific commands
```

**Example usage**:
```bash
go-starter new data-processor --type lambda --logger zerolog
cd data-processor
make build-lambda && make deploy
make logs  # View CloudWatch logs
```

## Choosing the Right Type

### Decision Matrix

| Need | Web API | CLI | Library | Lambda |
|------|---------|-----|---------|--------|
| **REST API** | ✅ Perfect | ❌ No | ❌ No | ⚠️ Limited |
| **Command line tool** | ❌ No | ✅ Perfect | ❌ No | ❌ No |
| **Reusable package** | ❌ No | ❌ No | ✅ Perfect | ❌ No |
| **Event processing** | ⚠️ Possible | ⚠️ Possible | ❌ No | ✅ Perfect |
| **Microservice** | ✅ Great | ❌ No | ❌ No | ✅ Great |
| **Background jobs** | ⚠️ Possible | ✅ Great | ❌ No | ✅ Great |
| **Web scraping** | ⚠️ Possible | ✅ Great | ⚠️ Possible | ✅ Great |

### Common Use Cases

**Web API** - Choose when you need:
- REST/GraphQL APIs
- Web application backends
- Microservices
- Real-time APIs with WebSockets

**CLI** - Choose when you need:
- DevOps automation tools
- Data processing scripts
- System administration utilities
- Developer productivity tools

**Library** - Choose when you need:
- SDKs for external services
- Shared business logic
- Open source packages
- Internal company libraries

**Lambda** - Choose when you need:
- Event-driven processing
- Serverless APIs
- Scheduled tasks
- Image/data processing

## Performance Characteristics

### Startup Time
- **Web API**: ~100-500ms (persistent)
- **CLI**: ~50-200ms (per invocation)
- **Library**: N/A (embedded)
- **Lambda**: ~100ms warm, ~1-3s cold

### Resource Usage
- **Web API**: Moderate (persistent memory)
- **CLI**: Low (short-lived)
- **Library**: Minimal (caller-dependent)
- **Lambda**: Very low (pay-per-use)

### Scaling
- **Web API**: Horizontal scaling with load balancer
- **CLI**: Process-level parallelism
- **Library**: Scales with host application
- **Lambda**: Automatic scaling (0-15k concurrent)

## Next Steps

- 🚀 **[Quick Start](GETTING_STARTED.md)** - Create your first project
- ⚙️ **[Configuration](CONFIGURATION.md)** - Customize your setup
- 📊 **[Logger Guide](LOGGER_GUIDE.md)** - Choose logging strategy