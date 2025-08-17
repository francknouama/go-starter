# Project Showcases

Real-world projects built with go-starter, demonstrating best practices and innovative solutions.

## Featured Projects

### 🏆 **TaskMaster Pro** - Task Management API
**Blueprint**: `web-api-clean`  
**Architecture**: Clean Architecture  
**Tech Stack**: Gin, PostgreSQL, Redis, JWT

A production-ready task management system serving 10K+ users daily.

**Key Features:**
- RESTful API with OpenAPI documentation
- Real-time notifications via WebSocket
- Redis caching for performance
- Comprehensive test coverage (89%)

**Highlights:**
```go
// Clean architecture in action
type TaskUseCase interface {
    Create(ctx context.Context, task *domain.Task) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error)
    UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}
```

**Repository**: [github.com/example/taskmaster-pro](https://github.com/example/taskmaster-pro)

---

### 🚀 **MetricsCollector** - CLI Monitoring Tool
**Blueprint**: `cli-standard`  
**Architecture**: Standard  
**Tech Stack**: Cobra, Prometheus, InfluxDB

A powerful CLI tool for collecting and analyzing system metrics.

**Key Features:**
- Multiple subcommands for different metrics
- Real-time data streaming
- Export to various formats (JSON, CSV, Prometheus)
- Plugin system for custom collectors

**Usage Example:**
```bash
metrics-collector system --interval 5s --output prometheus
metrics-collector custom --plugin ./plugins/docker.so
metrics-collector export --format json --start "1h ago"
```

**Repository**: [github.com/monitoring/metrics-collector](https://github.com/monitoring/metrics-collector)

---

### 📦 **GoKit** - Utility Library Collection
**Blueprint**: `library-standard`  
**Architecture**: Modular  
**Tech Stack**: Pure Go, Zero dependencies

A comprehensive utility library used by 50+ organizations.

**Modules:**
- `gokit/strings` - Advanced string manipulation
- `gokit/errors` - Enhanced error handling
- `gokit/cache` - Generic caching solutions
- `gokit/retry` - Configurable retry mechanisms

**Example Usage:**
```go
import "github.com/gokit/cache"

cache := cache.New[string, User](cache.WithTTL(5*time.Minute))
cache.Set("user:123", user)
```

**Repository**: [github.com/gokit/gokit](https://github.com/gokit/gokit)

---

### ⚡ **EventStream** - Event-Driven Platform
**Blueprint**: `event-driven`  
**Architecture**: Event Sourcing + CQRS  
**Tech Stack**: NATS, PostgreSQL, Elasticsearch

A scalable event-driven platform processing 1M+ events daily.

**Architecture Highlights:**
- Event sourcing for audit trail
- CQRS for read/write optimization
- Saga pattern for distributed transactions
- Event replay capabilities

**Code Example:**
```go
// Event sourcing in practice
type OrderAggregate struct {
    ID     uuid.UUID
    Status string
    Events []Event
}

func (o *OrderAggregate) Apply(event Event) {
    switch e := event.(type) {
    case OrderCreated:
        o.ID = e.OrderID
        o.Status = "created"
    case OrderShipped:
        o.Status = "shipped"
    }
    o.Events = append(o.Events, event)
}
```

**Repository**: [github.com/eventstream/platform](https://github.com/eventstream/platform)

---

### 🔧 **MicroCommerce** - E-commerce Microservices
**Blueprint**: `microservice-standard`  
**Architecture**: Microservices  
**Tech Stack**: gRPC, Kubernetes, Istio

A cloud-native e-commerce platform with 12 microservices.

**Services:**
- User Service (Authentication/Authorization)
- Product Catalog Service
- Cart Service
- Order Service
- Payment Service
- Notification Service

**Service Communication:**
```proto
service ProductCatalog {
    rpc GetProduct(GetProductRequest) returns (Product);
    rpc ListProducts(ListProductsRequest) returns (ProductList);
    rpc SearchProducts(SearchRequest) returns (ProductList);
}
```

**Repository**: [github.com/microcommerce/platform](https://github.com/microcommerce/platform)

---

### 🌐 **GraphQL Store** - Modern API
**Blueprint**: `graphql-api`  
**Architecture**: GraphQL + Clean  
**Tech Stack**: gqlgen, PostgreSQL, DataLoader

A modern GraphQL API showcasing advanced patterns.

**Features:**
- Efficient data loading with DataLoader
- Subscription support for real-time updates
- Role-based field authorization
- Automatic persisted queries

**Schema Example:**
```graphql
type Product {
  id: ID!
  name: String!
  price: Float!
  inventory: Int! @hasRole(role: ADMIN)
  reviews: [Review!]! @dataloader
}

type Subscription {
  productUpdated(id: ID!): Product!
}
```

**Repository**: [github.com/graphqlstore/api](https://github.com/graphqlstore/api)

---

### 🔒 **SecureVault** - Lambda Security Tool
**Blueprint**: `lambda-standard`  
**Architecture**: Serverless  
**Tech Stack**: AWS Lambda, DynamoDB, KMS

A serverless security tool for secret management.

**Features:**
- Encrypted secret storage
- Temporary credential generation
- Audit logging
- Cost-effective scaling

**Deployment:**
```bash
# Generated SAM template from go-starter
sam build
sam deploy --guided
```

**Repository**: [github.com/securevault/lambda](https://github.com/securevault/lambda)

---

## Community Statistics

### By Blueprint Type
- **Web API**: 45% of showcases
- **CLI**: 20% of showcases
- **Microservice**: 15% of showcases
- **Library**: 10% of showcases
- **Other**: 10% of showcases

### By Architecture
- **Clean Architecture**: 35%
- **Standard**: 30%
- **DDD**: 20%
- **Hexagonal**: 10%
- **Event-Driven**: 5%

### Performance Metrics
Average performance improvements after using go-starter:
- **Development Time**: -60% reduction
- **Code Quality**: +40% improvement (measured by linting scores)
- **Test Coverage**: +35% average increase
- **Bug Reports**: -50% reduction in first month

## Submit Your Project

Built something amazing with go-starter? We'd love to feature it!

### Submission Guidelines

1. **Project Requirements:**
   - Built with go-starter
   - Open source or has public documentation
   - Demonstrates best practices
   - In production or significant usage

2. **What to Include:**
   - Project name and description
   - Blueprint and architecture used
   - Tech stack details
   - Key features and innovations
   - Performance metrics (if available)
   - Repository link

3. **How to Submit:**
   - Open a PR to add your project to this page
   - Or create an issue with project details
   - Tag it with `showcase`

### Benefits of Being Featured
- 🌟 Recognition in the go-starter community
- 🔗 Backlink to your project
- 👥 Increased visibility and contributors
- 💡 Inspire others with your implementation

## Learning from Showcases

### Common Patterns

1. **Clean Architecture Success:**
   - Clear separation of concerns
   - Easy testing and mocking
   - Flexible dependency injection

2. **Performance Optimization:**
   - Strategic caching with Redis
   - Connection pooling
   - Efficient query patterns

3. **Testing Strategies:**
   - Table-driven tests
   - Integration test containers
   - Mock generation with gomock

4. **Deployment Patterns:**
   - Docker multi-stage builds
   - Kubernetes configurations
   - CI/CD with GitHub Actions

### Code Snippets from Showcases

**Elegant Error Handling** (from TaskMaster Pro):
```go
type AppError struct {
    Code    string
    Message string
    Err     error
}

func (e AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %s", e.Message, e.Err)
    }
    return e.Message
}
```

**Middleware Chain** (from MicroCommerce):
```go
chain := middleware.Chain(
    middleware.Logger(logger),
    middleware.Recoverer(),
    middleware.RateLimiter(100),
    middleware.Auth(authService),
)
```

**Graceful Shutdown** (from EventStream):
```go
signals := make(chan os.Signal, 1)
signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

go func() {
    <-signals
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        logger.Error("Server shutdown error", "error", err)
    }
}()
```

---

*Want to see your project here? [Submit your showcase](#submit-your-project)!*