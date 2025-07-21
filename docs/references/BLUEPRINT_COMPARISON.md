# Blueprint Comparison Guide

A comprehensive comparison of all go-starter blueprints to help you choose the right starting point for your project.

## Quick Decision Matrix

| Blueprint | Best For | Not Ideal For | Key Features | Logger Recommendation |
|----------|----------|---------------|--------------|----------------------|
| **Web API** | REST services, microservices, backend APIs | CLI tools, libraries, static sites | HTTP routing, middleware, database integration | `zap` for high-traffic, `slog` for standard |
| **CLI** | Command-line tools, scripts, automation | Web services, libraries | Cobra commands, config management, completions | `logrus` for rich output, `slog` for simplicity |
| **Library** | Reusable packages, SDKs, shared code | Standalone applications | Clean API, examples, minimal dependencies | `slog` for compatibility, minimal interface |
| **Lambda** | Serverless functions, event processing | Long-running services, stateful apps | AWS integration, CloudWatch logging, SAM deploy | `zerolog` for JSON logs, `zap` for performance |

## Detailed Blueprint Comparison

### 1. Web API Blueprint

#### Overview
The Web API blueprint creates a production-ready REST API server with comprehensive features for building scalable web services.

#### Architecture
```
web-api/
├── cmd/server/          # Application entry point
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware chain
│   ├── models/          # Data models
│   ├── services/        # Business logic layer
│   ├── repository/      # Data access layer
│   └── database/        # Database connections
├── configs/             # Environment configs
├── migrations/          # Database migrations
└── docker/              # Container configs
```

#### Key Features
- **Framework**: Gin (lightweight, high-performance)
- **Database**: GORM ORM with migration support
- **Middleware**: CORS, logging, recovery, authentication
- **API Documentation**: OpenAPI/Swagger support
- **Health Checks**: Liveness and readiness probes
- **Graceful Shutdown**: Proper connection cleanup
- **Docker Support**: Multi-stage builds, compose files

#### When to Use
✅ Building REST APIs or microservices  
✅ Need database integration  
✅ Require authentication and authorization  
✅ Building backend for web/mobile apps  
✅ Creating internal services  

#### When NOT to Use
❌ Building command-line tools  
❌ Creating reusable libraries  
❌ Serverless/FaaS environments  
❌ Simple scripts or utilities  

#### Configuration Example
```yaml
server:
  port: 8080
  timeout: 30s
  
database:
  driver: postgres
  host: localhost
  port: 5432
  
logger:
  level: info
  format: json
```

#### Sample Code
```go
// Handler example
func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    user, err := h.service.CreateUser(c.Request.Context(), req)
    if err != nil {
        h.logger.Error("Failed to create user", "error", err)
        c.JSON(500, gin.H{"error": "internal server error"})
        return
    }
    
    c.JSON(201, user)
}
```

---

### 2. CLI Application Blueprint

#### Overview
The CLI blueprint creates a feature-rich command-line application using the Cobra framework with professional CLI patterns.

#### Architecture
```
cli-app/
├── cmd/                 # Command definitions
│   ├── root.go          # Root command setup
│   ├── version.go       # Version subcommand
│   └── completion.go    # Shell completions
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   └── logger/          # Logger implementation
└── configs/             # Default configurations
```

#### Key Features
- **Framework**: Cobra (industry standard for CLIs)
- **Subcommands**: Hierarchical command structure
- **Configuration**: Viper integration with multiple sources
- **Completions**: Bash, Zsh, Fish, PowerShell support
- **Global Flags**: Consistent flag handling
- **Help System**: Auto-generated help text
- **Version Info**: Built-in version command

#### When to Use
✅ Building developer tools  
✅ Creating system utilities  
✅ Automation scripts  
✅ DevOps tooling  
✅ Interactive terminal applications  

#### When NOT to Use
❌ Web services or APIs  
❌ Libraries or packages  
❌ GUI applications  
❌ Mobile/web applications  

#### Configuration Example
```yaml
# CLI config supports multiple formats
output:
  format: table  # table, json, yaml
  color: auto    # auto, always, never
  
defaults:
  timeout: 30s
  retries: 3
```

#### Sample Code
```go
// Command example
var processCmd = &cobra.Command{
    Use:   "process [files...]",
    Short: "Process one or more files",
    Long: `Process files with various transformations.
Supports multiple file formats and concurrent processing.`,
    Args: cobra.MinimumNArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        format, _ := cmd.Flags().GetString("format")
        concurrent, _ := cmd.Flags().GetBool("concurrent")
        
        logger.Info("Processing files", 
            "count", len(args),
            "format", format,
            "concurrent", concurrent,
        )
        
        return processFiles(args, format, concurrent)
    },
}
```

---

### 3. Go Library Blueprint

#### Overview
The Library blueprint creates a well-structured Go package designed for reuse with clean APIs and comprehensive documentation.

#### Architecture
```
library/
├── mylib.go             # Main library interface
├── types.go             # Public types/constants
├── errors.go            # Error definitions
├── options.go           # Configuration options
├── internal/            # Private implementation
│   └── core/            # Core logic
├── examples/            # Usage examples
│   ├── basic/           # Simple examples
│   └── advanced/        # Complex examples
└── docs/                # Additional documentation
```

#### Key Features
- **Clean API**: Well-defined public interface
- **Zero Dependencies**: Minimal external dependencies
- **Functional Options**: Flexible configuration pattern
- **Examples**: Runnable example programs
- **Documentation**: GoDoc-friendly comments
- **Testing**: Comprehensive test coverage
- **Benchmarks**: Performance benchmarks included

#### When to Use
✅ Creating reusable packages  
✅ Building SDKs or clients  
✅ Sharing code between projects  
✅ Open source libraries  
✅ Internal company packages  

#### When NOT to Use
❌ Standalone applications  
❌ Web services  
❌ CLI tools  
❌ One-off scripts  

#### API Design Example
```go
// Clean public API
package mylib

// Client represents the main library interface
type Client interface {
    // Connect establishes a connection with the given options
    Connect(ctx context.Context, opts ...Option) error
    
    // Process handles the input data and returns results
    Process(ctx context.Context, data []byte) (*Result, error)
    
    // Close cleanly shuts down the client
    Close() error
}

// Option configures the client
type Option func(*clientOptions)

// WithTimeout sets the operation timeout
func WithTimeout(d time.Duration) Option {
    return func(o *clientOptions) {
        o.timeout = d
    }
}

// WithLogger sets a custom logger
func WithLogger(logger Logger) Option {
    return func(o *clientOptions) {
        o.logger = logger
    }
}
```

---

### 4. AWS Lambda Blueprint

#### Overview
The Lambda blueprint creates a serverless function optimized for AWS Lambda with API Gateway integration and CloudWatch logging.

#### Architecture
```
lambda/
├── main.go              # Lambda entry point
├── internal/            # Private application code
│   ├── handler/         # Request handling logic
│   ├── logger/          # CloudWatch logger
│   └── response/        # API Gateway responses
├── template.yaml        # SAM deployment config
└── events/              # Test event samples
```

#### Key Features
- **Runtime**: AWS Lambda Go runtime
- **API Gateway**: Request/response handling
- **CloudWatch**: Optimized structured logging
- **SAM Support**: Infrastructure as code
- **Local Testing**: SAM CLI integration
- **Cross-compilation**: Build scripts for Linux
- **Environment Config**: Lambda environment variables

#### When to Use
✅ Event-driven processing  
✅ API Gateway backends  
✅ Scheduled tasks  
✅ S3/SQS/SNS triggers  
✅ Serverless architectures  

#### When NOT to Use
❌ Long-running processes (>15 min)  
❌ Stateful applications  
❌ WebSocket servers  
❌ High-frequency/low-latency requirements  

#### Handler Example
```go
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // Extract request ID for tracing
    requestID := request.RequestContext.RequestID
    logger := logger.With("requestId", requestID)
    
    logger.Info("Processing request",
        "method", request.HTTPMethod,
        "path", request.Path,
        "sourceIP", request.RequestContext.Identity.SourceIP,
    )
    
    // Process request
    result, err := processRequest(request.Body)
    if err != nil {
        logger.Error("Processing failed", "error", err)
        return response.Error(500, "Internal server error"), nil
    }
    
    return response.Success(result), nil
}
```

## Logger Selection by Blueprint

### Web API Logger Recommendations

| Logger | When to Use | Configuration |
|--------|------------|---------------|
| **zap** | High-traffic APIs, microservices | Zero allocation, fastest performance |
| **slog** | Standard APIs, moderate traffic | Built-in, good performance |
| **zerolog** | JSON-heavy APIs, cloud-native | Clean JSON output |
| **logrus** | Legacy systems, rich features | Extensive ecosystem |

### CLI Logger Recommendations

| Logger | When to Use | Configuration |
|--------|------------|---------------|
| **logrus** | Rich terminal output, colors | Best for human-readable logs |
| **slog** | Simple CLIs, standard logging | No dependencies |
| **zap** | Performance-critical tools | Fast structured logging |
| **zerolog** | JSON output CLIs | Machine-readable output |

### Library Logger Recommendations

| Logger | When to Use | Configuration |
|--------|------------|---------------|
| **slog** | Maximum compatibility | Standard library |
| **Interface only** | Let users choose | Minimal logger interface |
| **zap** | Performance-critical libs | Document the dependency |
| **zerolog** | JSON-focused libraries | Clean API |

### Lambda Logger Recommendations

| Logger | When to Use | Configuration |
|--------|------------|---------------|
| **zerolog** | CloudWatch integration | Best JSON format |
| **zap** | High-performance lambdas | Fast with good JSON |
| **slog** | Simple functions | Standard library |
| **logrus** | Feature-rich logging | More overhead |

## Feature Comparison Matrix

| Feature | Web API | CLI | Library | Lambda |
|---------|---------|-----|---------|---------|
| **HTTP Server** | ✅ Gin Framework | ❌ | ❌ | ✅ API Gateway |
| **Database** | ✅ GORM + Migrations | ⚠️ Optional | ❌ | ⚠️ Optional |
| **Configuration** | ✅ YAML/Env | ✅ Viper | ⚠️ Options pattern | ✅ Environment |
| **Docker** | ✅ Multi-stage | ✅ Simple | ❌ | ❌ SAM instead |
| **Testing** | ✅ Unit + Integration | ✅ Command tests | ✅ Unit + Benchmarks | ✅ Handler tests |
| **CI/CD** | ✅ GitHub Actions | ✅ GitHub Actions | ✅ GitHub Actions | ✅ SAM Deploy |
| **Middleware** | ✅ Full chain | ❌ | ❌ | ⚠️ Limited |
| **Examples** | ✅ In tests | ✅ In help text | ✅ Dedicated folder | ✅ Event samples |
| **Documentation** | ✅ OpenAPI | ✅ Cobra help | ✅ GoDoc | ✅ README |
| **Deployment** | ✅ Container/Binary | ✅ Binary | ❌ Library | ✅ SAM/Serverless |

## Decision Flow Chart

```
Start: What are you building?
│
├─> Is it a web service or API?
│   └─> Yes: Use Web API Template
│       └─> High traffic? → Use zap logger
│       └─> Standard use? → Use slog logger
│
├─> Is it a command-line tool?
│   └─> Yes: Use CLI Template
│       └─> Rich output needed? → Use logrus logger
│       └─> Simple tool? → Use slog logger
│
├─> Is it a reusable package?
│   └─> Yes: Use Library Template
│       └─> Always use slog or minimal interface
│
└─> Is it a serverless function?
    └─> Yes: Use Lambda Template
        └─> Use zerolog for CloudWatch
```

## Migration Paths

### From Web API to Microservice
1. Start with Web API blueprint
2. Add gRPC support when needed
3. Implement service mesh integration
4. Add distributed tracing

### From CLI to Web API
1. Extract business logic to services
2. Add HTTP handlers
3. Implement middleware
4. Add database layer

### From Library to Application
1. Create cmd/ directory
2. Add main.go entry point
3. Implement configuration
4. Add operational concerns

### From Lambda to Web API
1. Extract handler logic
2. Add HTTP routing
3. Implement persistent connections
4. Add stateful components

## Best Practices by Blueprint Type

### Web API Best Practices
- Use middleware for cross-cutting concerns
- Implement proper error handling
- Add request validation
- Use dependency injection
- Implement health checks
- Add metrics and monitoring

### CLI Best Practices
- Use clear command hierarchy
- Provide helpful error messages
- Add progress indicators
- Support configuration files
- Implement --dry-run flag
- Add shell completions

### Library Best Practices
- Keep public API minimal
- Use semantic versioning
- Provide comprehensive examples
- Document all public functions
- Avoid global state
- Minimize dependencies

### Lambda Best Practices
- Keep functions small and focused
- Handle cold starts efficiently
- Use environment variables
- Implement proper error handling
- Log structured data
- Monitor execution time

## Conclusion

Choose your blueprint based on:
1. **Application type**: What you're building
2. **Deployment target**: Where it will run
3. **Performance needs**: Request volume and latency
4. **Team experience**: Familiarity with patterns
5. **Future growth**: Scalability requirements

Remember: You can always generate multiple projects to experiment with different blueprints and find the best fit for your needs.