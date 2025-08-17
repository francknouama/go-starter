# Best Practices

Community-driven best practices for building production-ready Go applications with go-starter.

## Table of Contents
- [Project Organization](#project-organization)
- [Code Style](#code-style)
- [Error Handling](#error-handling)
- [Testing Strategies](#testing-strategies)
- [Performance Optimization](#performance-optimization)
- [Security Best Practices](#security-best-practices)
- [Deployment Strategies](#deployment-strategies)
- [Monitoring & Observability](#monitoring--observability)

## Project Organization

### Standard Layout

Follow the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

```
myproject/
├── cmd/              # Main applications
├── internal/         # Private application code
├── pkg/              # Library code that can be used by external apps
├── api/              # API definitions (OpenAPI, Proto, etc.)
├── web/              # Web assets (if applicable)
├── configs/          # Configuration file templates
├── deployments/      # Deployment configurations
├── scripts/          # Scripts for build, install, analysis
├── tests/            # Additional test apps and test data
└── docs/             # Design and user documents
```

### Package Design

**✅ DO:**
```go
// Good: Clear package responsibility
package user

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}
```

**❌ DON'T:**
```go
// Bad: Mixed responsibilities
package utils

func ValidateUser() {}
func SendEmail() {}
func GenerateToken() {}
```

### Dependency Injection

Use constructor injection for testability:

```go
// internal/service/order.go
type OrderService struct {
    orderRepo     OrderRepository
    inventoryRepo InventoryRepository
    emailService  EmailService
    logger        *slog.Logger
}

func NewOrderService(
    orderRepo OrderRepository,
    inventoryRepo InventoryRepository,
    emailService EmailService,
    logger *slog.Logger,
) *OrderService {
    return &OrderService{
        orderRepo:     orderRepo,
        inventoryRepo: inventoryRepo,
        emailService:  emailService,
        logger:        logger,
    }
}
```

## Code Style

### Naming Conventions

**Interfaces:**
```go
// Good: -er suffix for single-method interfaces
type Reader interface {
    Read([]byte) (int, error)
}

// Good: Descriptive names for multi-method interfaces
type UserRepository interface {
    Create(context.Context, *User) error
    GetByID(context.Context, string) (*User, error)
    Update(context.Context, *User) error
    Delete(context.Context, string) error
}
```

**Constants:**
```go
// Good: CamelCase for exported, camelCase for unexported
const (
    MaxRetries = 3
    minTimeout = 100 * time.Millisecond
)

// Good: Typed constants
type Status string

const (
    StatusActive   Status = "active"
    StatusInactive Status = "inactive"
    StatusPending  Status = "pending"
)
```

### Function Design

**Keep functions small and focused:**
```go
// Good: Single responsibility
func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
    if err := s.validateInput(input); err != nil {
        return nil, err
    }
    
    user := s.buildUser(input)
    
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }
    
    s.sendWelcomeEmail(ctx, user)
    
    return user, nil
}
```

**Use functional options for complex configurations:**
```go
type ServerOption func(*Server)

func WithPort(port int) ServerOption {
    return func(s *Server) {
        s.port = port
    }
}

func WithTimeout(timeout time.Duration) ServerOption {
    return func(s *Server) {
        s.timeout = timeout
    }
}

func NewServer(opts ...ServerOption) *Server {
    s := &Server{
        port:    8080,
        timeout: 30 * time.Second,
    }
    
    for _, opt := range opts {
        opt(s)
    }
    
    return s
}
```

## Error Handling

### Error Types

Define domain-specific errors:

```go
// internal/errors/errors.go
type Error struct {
    Code    string
    Message string
    Err     error
}

func (e Error) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %s", e.Message, e.Err)
    }
    return e.Message
}

func (e Error) Unwrap() error {
    return e.Err
}

// Predefined errors
var (
    ErrNotFound     = Error{Code: "NOT_FOUND", Message: "resource not found"}
    ErrUnauthorized = Error{Code: "UNAUTHORIZED", Message: "unauthorized access"}
    ErrInvalidInput = Error{Code: "INVALID_INPUT", Message: "invalid input"}
)
```

### Error Wrapping

Always wrap errors with context:

```go
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("get user %s: %w", id, ErrNotFound)
        }
        return nil, fmt.Errorf("get user %s: %w", id, err)
    }
    
    return user, nil
}
```

### Error Handling in HTTP

```go
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    userID := chi.URLParam(r, "id")
    
    user, err := h.service.GetUser(r.Context(), userID)
    if err != nil {
        switch {
        case errors.Is(err, ErrNotFound):
            h.respondError(w, http.StatusNotFound, "User not found")
        case errors.Is(err, ErrUnauthorized):
            h.respondError(w, http.StatusUnauthorized, "Unauthorized")
        default:
            h.logger.Error("Get user failed", "error", err)
            h.respondError(w, http.StatusInternalServerError, "Internal server error")
        }
        return
    }
    
    h.respondJSON(w, http.StatusOK, user)
}
```

## Testing Strategies

### Table-Driven Tests

```go
func TestCalculateDiscount(t *testing.T) {
    tests := []struct {
        name     string
        amount   float64
        tier     string
        expected float64
        wantErr  bool
    }{
        {
            name:     "silver tier discount",
            amount:   100.00,
            tier:     "silver",
            expected: 95.00,
            wantErr:  false,
        },
        {
            name:     "gold tier discount",
            amount:   100.00,
            tier:     "gold",
            expected: 90.00,
            wantErr:  false,
        },
        {
            name:     "invalid tier",
            amount:   100.00,
            tier:     "platinum",
            expected: 0,
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := CalculateDiscount(tt.amount, tt.tier)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("CalculateDiscount() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if result != tt.expected {
                t.Errorf("CalculateDiscount() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Test Fixtures

```go
// tests/fixtures/user.go
func NewTestUser(t *testing.T) *domain.User {
    t.Helper()
    
    return &domain.User{
        ID:        uuid.New().String(),
        Email:     "test@example.com",
        Name:      "Test User",
        CreatedAt: time.Now(),
    }
}

func NewTestUserWithEmail(t *testing.T, email string) *domain.User {
    t.Helper()
    
    user := NewTestUser(t)
    user.Email = email
    return user
}
```

### Integration Testing

```go
// tests/integration/user_api_test.go
func TestUserAPI(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Setup
    container := setupTestContainer(t)
    defer container.Terminate(context.Background())
    
    db := setupTestDB(t, container)
    defer db.Close()
    
    app := setupTestApp(t, db)
    
    t.Run("Create User", func(t *testing.T) {
        payload := `{"email":"test@example.com","name":"Test User"}`
        
        req := httptest.NewRequest("POST", "/api/users", strings.NewReader(payload))
        req.Header.Set("Content-Type", "application/json")
        
        rec := httptest.NewRecorder()
        app.ServeHTTP(rec, req)
        
        assert.Equal(t, http.StatusCreated, rec.Code)
        
        var user map[string]interface{}
        err := json.Unmarshal(rec.Body.Bytes(), &user)
        require.NoError(t, err)
        
        assert.NotEmpty(t, user["id"])
        assert.Equal(t, "test@example.com", user["email"])
    })
}
```

## Performance Optimization

### Connection Pooling

```go
// internal/database/postgres.go
func NewPostgresDB(cfg Config) (*sql.DB, error) {
    db, err := sql.Open("postgres", cfg.DSN)
    if err != nil {
        return nil, fmt.Errorf("open database: %w", err)
    }
    
    // Configure connection pool
    db.SetMaxOpenConns(cfg.MaxOpenConns)       // Default: 25
    db.SetMaxIdleConns(cfg.MaxIdleConns)       // Default: 25
    db.SetConnMaxLifetime(cfg.ConnMaxLifetime) // Default: 5 minutes
    db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime) // Default: 5 minutes
    
    // Verify connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("ping database: %w", err)
    }
    
    return db, nil
}
```

### Caching Strategy

```go
// internal/cache/redis.go
type Cache interface {
    Get(ctx context.Context, key string, dest interface{}) error
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
}

type RedisCache struct {
    client *redis.Client
}

func (c *RedisCache) GetWithLoader(
    ctx context.Context,
    key string,
    dest interface{},
    loader func() (interface{}, error),
    ttl time.Duration,
) error {
    // Try cache first
    err := c.Get(ctx, key, dest)
    if err == nil {
        return nil
    }
    
    // Cache miss - load from source
    data, err := loader()
    if err != nil {
        return fmt.Errorf("loader failed: %w", err)
    }
    
    // Store in cache
    if err := c.Set(ctx, key, data, ttl); err != nil {
        // Log but don't fail
        log.Printf("cache set failed: %v", err)
    }
    
    // Copy to destination
    return copier.Copy(dest, data)
}
```

### Concurrent Processing

```go
// internal/service/batch.go
func (s *Service) ProcessBatch(ctx context.Context, items []Item) error {
    const workers = 10
    
    ch := make(chan Item, len(items))
    errCh := make(chan error, workers)
    
    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for item := range ch {
                if err := s.processItem(ctx, item); err != nil {
                    select {
                    case errCh <- err:
                    default:
                    }
                    return
                }
            }
        }()
    }
    
    // Send work
    for _, item := range items {
        select {
        case ch <- item:
        case <-ctx.Done():
            close(ch)
            return ctx.Err()
        }
    }
    close(ch)
    
    // Wait for completion
    wg.Wait()
    close(errCh)
    
    // Check for errors
    if err := <-errCh; err != nil {
        return err
    }
    
    return nil
}
```

## Security Best Practices

### Input Validation

```go
// internal/validation/user.go
func ValidateCreateUserInput(input CreateUserInput) error {
    if err := ValidateEmail(input.Email); err != nil {
        return fmt.Errorf("email: %w", err)
    }
    
    if err := ValidatePassword(input.Password); err != nil {
        return fmt.Errorf("password: %w", err)
    }
    
    if len(input.Name) < 2 || len(input.Name) > 100 {
        return errors.New("name must be between 2 and 100 characters")
    }
    
    return nil
}

func ValidatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("must be at least 8 characters")
    }
    
    var hasUpper, hasLower, hasNumber bool
    for _, ch := range password {
        switch {
        case unicode.IsUpper(ch):
            hasUpper = true
        case unicode.IsLower(ch):
            hasLower = true
        case unicode.IsNumber(ch):
            hasNumber = true
        }
    }
    
    if !hasUpper || !hasLower || !hasNumber {
        return errors.New("must contain uppercase, lowercase, and number")
    }
    
    return nil
}
```

### SQL Injection Prevention

```go
// Always use parameterized queries
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
    query := `
        SELECT id, email, name, password_hash, created_at
        FROM users
        WHERE email = $1
    `
    
    var user User
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Email,
        &user.Name,
        &user.PasswordHash,
        &user.CreatedAt,
    )
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNotFound
        }
        return nil, fmt.Errorf("query user: %w", err)
    }
    
    return &user, nil
}
```

### Authentication & Authorization

```go
// internal/middleware/auth.go
func AuthMiddleware(authService AuthService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := extractToken(r)
            if token == "" {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            
            claims, err := authService.ValidateToken(token)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }
            
            // Add claims to context
            ctx := context.WithValue(r.Context(), "claims", claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func RequireRole(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims, ok := r.Context().Value("claims").(*Claims)
            if !ok || !claims.HasRole(role) {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

## Deployment Strategies

### Docker Multi-Stage Build

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build application
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server cmd/server/main.go

# Final stage
FROM scratch

# Copy certificates and timezone data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy binary
COPY --from=builder /app/server /server

# Copy config files
COPY --from=builder /app/configs /configs

EXPOSE 8080

ENTRYPOINT ["/server"]
```

### Kubernetes Deployment

```yaml
# deployments/k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-server
  labels:
    app: api-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-server
  template:
    metadata:
      labels:
        app: api-server
    spec:
      containers:
      - name: api-server
        image: myregistry/api-server:v1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: database-url
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### Health Checks

```go
// internal/health/health.go
type Checker interface {
    Check(ctx context.Context) error
}

type HealthService struct {
    checkers map[string]Checker
}

func (h *HealthService) LivenessHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "alive",
    })
}

func (h *HealthService) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()
    
    results := make(map[string]string)
    allHealthy := true
    
    for name, checker := range h.checkers {
        if err := checker.Check(ctx); err != nil {
            results[name] = err.Error()
            allHealthy = false
        } else {
            results[name] = "healthy"
        }
    }
    
    status := http.StatusOK
    if !allHealthy {
        status = http.StatusServiceUnavailable
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  results,
        "healthy": allHealthy,
    })
}
```

## Monitoring & Observability

### Structured Logging

```go
// internal/logger/logger.go
func NewLogger(cfg Config) *slog.Logger {
    var handler slog.Handler
    
    opts := &slog.HandlerOptions{
        Level: cfg.Level,
        AddSource: cfg.AddSource,
    }
    
    switch cfg.Format {
    case "json":
        handler = slog.NewJSONHandler(os.Stdout, opts)
    default:
        handler = slog.NewTextHandler(os.Stdout, opts)
    }
    
    logger := slog.New(handler)
    
    // Add default fields
    logger = logger.With(
        slog.String("service", cfg.ServiceName),
        slog.String("version", cfg.Version),
        slog.String("environment", cfg.Environment),
    )
    
    return logger
}

// Usage
logger.Info("Processing order",
    slog.String("order_id", orderID),
    slog.String("customer_id", customerID),
    slog.Float64("amount", amount),
    slog.Duration("processing_time", duration),
)
```

### Metrics Collection

```go
// internal/metrics/prometheus.go
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
}

func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
        
        next.ServeHTTP(ww, r)
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(ww.Status())
        
        httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
        httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
    })
}
```

### Distributed Tracing

```go
// internal/tracing/tracing.go
func InitTracer(cfg Config) (func(), error) {
    exporter, err := jaeger.New(
        jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.JaegerEndpoint)),
    )
    if err != nil {
        return nil, fmt.Errorf("create jaeger exporter: %w", err)
    }
    
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceName(cfg.ServiceName),
            semconv.ServiceVersion(cfg.Version),
        )),
    )
    
    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.TraceContext{})
    
    return func() {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        tp.Shutdown(ctx)
    }, nil
}

// Usage in handler
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    ctx, span := otel.Tracer("api").Start(r.Context(), "GetUser")
    defer span.End()
    
    userID := chi.URLParam(r, "id")
    span.SetAttributes(attribute.String("user.id", userID))
    
    user, err := h.service.GetUser(ctx, userID)
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        // ... error handling
        return
    }
    
    span.SetStatus(codes.Ok, "")
    // ... success response
}
```

## Summary

These best practices come from real-world experience building production Go applications. Key takeaways:

1. **Keep it simple**: Don't over-engineer solutions
2. **Test thoroughly**: Aim for >80% coverage
3. **Handle errors gracefully**: Always wrap with context
4. **Monitor everything**: You can't fix what you can't see
5. **Secure by default**: Validate inputs, use parameterized queries
6. **Optimize when needed**: Profile first, optimize second

Remember: These are guidelines, not rules. Adapt them to your specific needs and constraints.

---
*Have a best practice to share? [Submit a PR](https://github.com/francknouama/go-starter/blob/main/docs/05-community/best-practices.md)!*