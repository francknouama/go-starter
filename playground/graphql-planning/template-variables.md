# GraphQL API Template Variables

## Core Template Variables

### Project Configuration
```go
type ProjectConfig struct {
    ProjectName    string `yaml:"project_name"`    // e.g., "ecommerce-api"
    ModulePath     string `yaml:"module_path"`     // e.g., "github.com/user/ecommerce-api"
    GoVersion      string `yaml:"go_version"`      // e.g., "1.21"
    Framework      string `yaml:"framework"`       // e.g., "gin", "echo", "chi"
    LoggerType     string `yaml:"logger_type"`     // e.g., "slog", "zap", "logrus"
    Architecture   string `yaml:"architecture"`    // e.g., "standard", "clean", "ddd"
}
```

### GraphQL-Specific Configuration
```go
type GraphQLConfig struct {
    Library        string   `yaml:"library"`         // "gqlgen", "graphql-go" (future)
    Playground     bool     `yaml:"playground"`      // Enable GraphQL playground
    Subscriptions  bool     `yaml:"subscriptions"`   // Enable WebSocket subscriptions
    SchemaPath     string   `yaml:"schema_path"`     // "schema/schema.graphql"
    GeneratedPath  string   `yaml:"generated_path"`  // "internal/generated"
    ComplexityLimit int     `yaml:"complexity_limit"` // Query complexity limit
    CustomScalars  []string `yaml:"custom_scalars"`  // ["DateTime", "Upload"]
    Directives     []string `yaml:"directives"`      // ["auth", "validate"]
}
```

### Database Configuration
```go
type DatabaseConfig struct {
    Enabled  bool   `yaml:"enabled"`   // Enable database integration
    Driver   string `yaml:"driver"`    // "postgres", "mysql", "sqlite"
    ORM      string `yaml:"orm"`       // "gorm", "sqlx", "ent"
    Host     string `yaml:"host"`      // Database host
    Port     int    `yaml:"port"`      // Database port
    Name     string `yaml:"name"`      // Database name
    Migrations bool `yaml:"migrations"` // Include migration support
}
```

### Authentication Configuration
```go
type AuthConfig struct {
    Enabled   bool     `yaml:"enabled"`    // Enable authentication
    Type      string   `yaml:"type"`       // "jwt", "oauth2", "apikey"
    Providers []string `yaml:"providers"`  // ["google", "github"] for OAuth2
    JWTSecret string   `yaml:"jwt_secret"` // JWT signing secret
    Roles     []string `yaml:"roles"`      // ["USER", "ADMIN", "MODERATOR"]
}
```

### Features Configuration
```go
type FeaturesConfig struct {
    CORS        CORSConfig        `yaml:"cors"`
    RateLimit   RateLimitConfig   `yaml:"rate_limit"`
    Monitoring  MonitoringConfig  `yaml:"monitoring"`
    Testing     TestingConfig     `yaml:"testing"`
    Docker      DockerConfig      `yaml:"docker"`
    Kubernetes  KubernetesConfig  `yaml:"kubernetes"`
}

type CORSConfig struct {
    Enabled      bool     `yaml:"enabled"`
    AllowOrigins []string `yaml:"allow_origins"`
    AllowMethods []string `yaml:"allow_methods"`
    AllowHeaders []string `yaml:"allow_headers"`
}

type RateLimitConfig struct {
    Enabled  bool `yaml:"enabled"`
    Requests int  `yaml:"requests"` // Requests per window
    Window   int  `yaml:"window"`   // Window in seconds
}

type MonitoringConfig struct {
    Metrics bool `yaml:"metrics"`   // Prometheus metrics
    Tracing bool `yaml:"tracing"`   // OpenTelemetry tracing
    Health  bool `yaml:"health"`    // Health check endpoint
}

type TestingConfig struct {
    Unit        bool `yaml:"unit"`        // Unit tests
    Integration bool `yaml:"integration"` // Integration tests
    E2E         bool `yaml:"e2e"`         // End-to-end tests
    Benchmarks  bool `yaml:"benchmarks"`  // Benchmark tests
}

type DockerConfig struct {
    Enabled    bool   `yaml:"enabled"`
    BaseImage  string `yaml:"base_image"`  // "alpine", "scratch", "ubuntu"
    Multistage bool   `yaml:"multistage"`  // Use multi-stage build
}

type KubernetesConfig struct {
    Enabled    bool `yaml:"enabled"`
    Replicas   int  `yaml:"replicas"`
    AutoScale  bool `yaml:"auto_scale"`
    MinReplicas int `yaml:"min_replicas"`
    MaxReplicas int `yaml:"max_replicas"`
}
```

## Template Variable Usage Examples

### Schema Generation
```graphql
# schema/schema.graphql.tmpl
directive @auth(requires: Role = USER) on FIELD_DEFINITION
{{if .Features.Auth.Enabled}}
directive @validate(format: String) on INPUT_FIELD_DEFINITION
{{end}}

{{range .GraphQL.CustomScalars}}
scalar {{.}}
{{end}}

type Query {
    users: [User!]! {{if .Features.Auth.Enabled}}@auth{{end}}
    user(id: ID!): User {{if .Features.Auth.Enabled}}@auth{{end}}
}

{{if .GraphQL.Subscriptions}}
type Subscription {
    userUpdated(id: ID!): User! @auth
}
{{end}}

type User {
    id: ID!
    name: String!
    email: String!
    {{if .Features.Auth.Enabled}}
    role: Role!
    {{end}}
}

{{if .Features.Auth.Enabled}}
enum Role {
    {{range .Features.Auth.Roles}}
    {{.}}
    {{end}}
}
{{end}}
```

### Resolver Generation
```go
// resolvers/resolver.go.tmpl
package resolvers

import (
    "context"
    {{if .Database.Enabled}}
    {{if eq .Database.ORM "gorm"}}
    "gorm.io/gorm"
    {{else if eq .Database.ORM "sqlx"}}
    "github.com/jmoiron/sqlx"
    {{end}}
    {{end}}
    "{{.ModulePath}}/internal/generated"
    "{{.ModulePath}}/internal/model"
)

type Resolver struct {
    {{if .Database.Enabled}}
    {{if eq .Database.ORM "gorm"}}
    DB *gorm.DB
    {{else if eq .Database.ORM "sqlx"}}
    DB *sqlx.DB
    {{end}}
    {{end}}
}

func (r *Resolver) Query() generated.QueryResolver {
    return &queryResolver{r}
}

{{if .GraphQL.Subscriptions}}
func (r *Resolver) Subscription() generated.SubscriptionResolver {
    return &subscriptionResolver{r}
}
{{end}}

type queryResolver struct{ *Resolver }
{{if .GraphQL.Subscriptions}}
type subscriptionResolver struct{ *Resolver }
{{end}}
```

### Configuration File Generation
```go
// config/config.go.tmpl
package config

import (
    "os"
    "strconv"
)

type Config struct {
    Server ServerConfig `yaml:"server"`
    {{if .Database.Enabled}}
    Database DatabaseConfig `yaml:"database"`
    {{end}}
    {{if .Features.Auth.Enabled}}
    Auth AuthConfig `yaml:"auth"`
    {{end}}
    GraphQL GraphQLConfig `yaml:"graphql"`
}

type ServerConfig struct {
    Port string `yaml:"port"`
    Host string `yaml:"host"`
}

{{if .Database.Enabled}}
type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    Database string `yaml:"database"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    SSLMode  string `yaml:"ssl_mode"`
}
{{end}}

{{if .Features.Auth.Enabled}}
type AuthConfig struct {
    JWTSecret string `yaml:"jwt_secret"`
    TokenTTL  int    `yaml:"token_ttl"`
}
{{end}}

type GraphQLConfig struct {
    ComplexityLimit int  `yaml:"complexity_limit"`
    Playground     bool `yaml:"playground"`
}

func Load() (*Config, error) {
    cfg := &Config{
        Server: ServerConfig{
            Port: getEnv("PORT", "8080"),
            Host: getEnv("HOST", "0.0.0.0"),
        },
        {{if .Database.Enabled}}
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnvAsInt("DB_PORT", {{.Database.Port}}),
            Database: getEnv("DB_NAME", "{{.Database.Name}}"),
            Username: getEnv("DB_USER", "{{.ProjectName}}"),
            Password: getEnv("DB_PASSWORD", ""),
            SSLMode:  getEnv("DB_SSL_MODE", "disable"),
        },
        {{end}}
        {{if .Features.Auth.Enabled}}
        Auth: AuthConfig{
            JWTSecret: getEnv("JWT_SECRET", "{{.Features.Auth.JWTSecret}}"),
            TokenTTL:  getEnvAsInt("TOKEN_TTL", 3600),
        },
        {{end}}
        GraphQL: GraphQLConfig{
            ComplexityLimit: {{.GraphQL.ComplexityLimit}},
            Playground:     {{.GraphQL.Playground}},
        },
    }
    
    return cfg, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}
```

### Docker Configuration
```dockerfile
# deployments/docker/Dockerfile.tmpl
{{if .Features.Docker.Multistage}}
# Build stage
FROM golang:{{.GoVersion}}-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
{{if eq .GraphQL.Library "gqlgen"}}
RUN go generate ./...
{{end}}
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Production stage
FROM {{.Features.Docker.BaseImage}}
{{if eq .Features.Docker.BaseImage "alpine"}}
RUN apk --no-cache add ca-certificates
{{end}}
WORKDIR /root/

COPY --from=builder /app/server .
{{else}}
FROM golang:{{.GoVersion}}-alpine

WORKDIR /app
COPY . .

{{if eq .GraphQL.Library "gqlgen"}}
RUN go generate ./...
{{end}}
RUN go build -o server cmd/server/main.go
{{end}}

EXPOSE 8080

CMD ["./server"]
```

## Variable Validation Rules

### Required Variables
- `ProjectName`: Must be valid directory name
- `ModulePath`: Must be valid Go module path
- `GoVersion`: Must be supported Go version
- `Framework`: Must be supported HTTP framework

### Conditional Requirements
- If `Database.Enabled = true`: `Database.Driver` and `Database.ORM` required
- If `Features.Auth.Enabled = true`: `Features.Auth.Type` required
- If `GraphQL.Subscriptions = true`: WebSocket support required

### Default Values
```yaml
# Default configuration
framework: "gin"
logger_type: "slog"
architecture: "standard"
go_version: "1.21"

graphql:
  library: "gqlgen"
  playground: true
  subscriptions: false
  complexity_limit: 1000
  schema_path: "schema/schema.graphql"
  generated_path: "internal/generated"

features:
  cors:
    enabled: true
    allow_origins: ["*"]
    allow_methods: ["GET", "POST", "OPTIONS"]
    allow_headers: ["Content-Type", "Authorization"]
  
  rate_limit:
    enabled: false
    requests: 100
    window: 60
  
  monitoring:
    metrics: false
    tracing: false
    health: true
  
  testing:
    unit: true
    integration: true
    e2e: false
    benchmarks: false
  
  docker:
    enabled: true
    base_image: "alpine"
    multistage: true
  
  kubernetes:
    enabled: false
    replicas: 3
    auto_scale: false
```