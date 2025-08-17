# Template Documentation

Comprehensive documentation for all project templates available in go-starter, including architecture patterns, features, and best practices for each project type.

## 📋 Template Overview

| Template | Complexity | Files | Use Case | Architecture |
|----------|------------|-------|----------|--------------|
| [CLI Simple](#cli-simple) | Beginner | 8 | Quick utilities | None |
| [CLI Standard](#cli-standard) | Intermediate | 29 | Production CLIs | MVC-lite |
| [Web API Standard](#web-api-standard) | Intermediate | 35 | REST APIs | Layered |
| [Web API Clean](#web-api-clean) | Advanced | 45 | Enterprise APIs | Clean Architecture |
| [Web API DDD](#web-api-ddd) | Advanced | 50 | Domain-rich APIs | Domain-Driven Design |
| [Web API Hexagonal](#web-api-hexagonal) | Expert | 55 | Testable APIs | Ports & Adapters |
| [Lambda Standard](#lambda-standard) | Beginner | 12 | Event handlers | Functional |
| [Lambda Proxy](#lambda-proxy) | Intermediate | 25 | API Gateway | Proxy pattern |
| [Library](#library) | Beginner | 15 | Shared packages | Public API |
| [Microservice](#microservice) | Advanced | 60 | Distributed systems | Service-oriented |
| [Monolith](#monolith) | Intermediate | 65 | Web applications | Modular monolith |
| [Workspace](#workspace) | Advanced | 40 | Multi-module | Workspace pattern |

---

## CLI Simple

> **Best for**: Learning Go, quick utilities, simple automation scripts

### 🎯 When to Use
- Building your first Go CLI tool
- Creating simple automation scripts
- Learning Go fundamentals
- Prototyping command-line utilities

### 📁 Project Structure (8 files)

```
my-tool/
├── main.go              # Entry point with basic CLI setup
├── cmd/
│   ├── root.go          # Cobra root command
│   └── version.go       # Version command
├── config.go            # Simple configuration
├── Makefile             # Build automation
├── README.md            # Getting started documentation
├── go.mod              # Module definition
└── go.sum              # Dependency checksums
```

### 🔧 Generated Features

#### Core Functionality
- **Basic CLI framework** using Cobra
- **Version command** with build info
- **Simple configuration** management
- **Help system** with usage examples
- **Build automation** via Makefile

#### Configuration
```go
// config.go
type Config struct {
    Verbose bool   `json:"verbose"`
    Output  string `json:"output"`
    Format  string `json:"format"`
}

func LoadConfig() (*Config, error) {
    return &Config{
        Verbose: false,
        Output:  "stdout",
        Format:  "text",
    }, nil
}
```

#### Main CLI Structure
```go
// cmd/root.go
var rootCmd = &cobra.Command{
    Use:   "my-tool",
    Short: "A brief description of your application",
    Long:  `A longer description...`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Hello from my-tool!")
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}
```

### 🚀 Quick Start
```bash
# Generate the project
go-starter new my-tool --type=cli --complexity=simple

# Build and run
cd my-tool
make build
./my-tool --help
```

### 🎓 Learning Path
1. **Understand the structure**: Review generated files
2. **Add basic commands**: Implement core functionality
3. **Add flags**: Learn Cobra flag handling
4. **Add configuration**: Implement config file support
5. **Graduate to standard**: When you need more features

---

## CLI Standard

> **Best for**: Production CLI tools, team projects, complex command-line applications

### 🎯 When to Use
- Building production CLI tools
- Multi-command applications
- Team development projects
- CLI tools with complex logic

### 📁 Project Structure (29 files)

```
my-tool/
├── main.go
├── cmd/
│   ├── root.go          # Root command with global flags
│   ├── version.go       # Version information
│   ├── completion.go    # Shell completion
│   ├── create.go        # Create command
│   ├── delete.go        # Delete command
│   ├── list.go         # List command
│   ├── update.go       # Update command
│   └── root_test.go    # Command tests
├── internal/
│   ├── config/
│   │   ├── config.go   # Configuration management
│   │   └── config_test.go
│   ├── logger/
│   │   ├── interface.go # Logger interface
│   │   └── logger.go   # Logger implementation
│   ├── errors/
│   │   └── errors.go   # Custom error types
│   ├── interactive/
│   │   └── prompt.go   # Interactive prompts
│   ├── output/
│   │   └── output.go   # Output formatting
│   └── version/
│       └── version.go  # Version info
├── configs/
│   └── config.yaml     # Default configuration
├── Dockerfile          # Container support
├── Makefile           # Build automation
└── README.md          # Comprehensive documentation
```

### 🔧 Generated Features

#### Advanced CLI Framework
- **Multiple subcommands** with full Cobra integration
- **Global and command-specific flags**
- **Configuration file support** (YAML, JSON, TOML)
- **Interactive prompts** for user input
- **Shell completion** (bash, zsh, fish, PowerShell)
- **Structured logging** with configurable levels
- **Output formatting** (table, JSON, YAML)
- **Comprehensive testing** framework

#### Configuration Management
```go
// internal/config/config.go
type Config struct {
    // Global settings
    LogLevel   string            `mapstructure:"log_level"`
    OutputDir  string            `mapstructure:"output_dir"`
    ConfigFile string            `mapstructure:"config_file"`
    
    // Command-specific settings
    Create CreateConfig `mapstructure:"create"`
    List   ListConfig   `mapstructure:"list"`
    Update UpdateConfig `mapstructure:"update"`
}

type CreateConfig struct {
    Template    string            `mapstructure:"template"`
    Variables   map[string]string `mapstructure:"variables"`
    Overwrite   bool             `mapstructure:"overwrite"`
}
```

#### Command Structure
```go
// cmd/create.go
var createCmd = &cobra.Command{
    Use:   "create [name]",
    Short: "Create a new resource",
    Long:  `Create a new resource with the specified configuration.`,
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        cfg := config.Get()
        logger := logger.Get()
        
        name := args[0]
        logger.Info("Creating resource", "name", name)
        
        // Implementation here
        return nil
    },
}

func init() {
    rootCmd.AddCommand(createCmd)
    createCmd.Flags().String("template", "", "Template to use")
    createCmd.Flags().Bool("overwrite", false, "Overwrite existing files")
}
```

### 🚀 Generation Command
```bash
go-starter new my-tool --type=cli --complexity=standard --logger=slog
```

### 📊 Use Cases
- **Developer tools**: Code generators, build tools, deployment scripts
- **System administration**: Server management, monitoring tools
- **Data processing**: ETL tools, data migration utilities
- **CI/CD tools**: Custom deployment and testing tools

---

## Web API Standard

> **Best for**: Most REST APIs, microservices, standard web backends

### 🎯 When to Use
- Building REST APIs
- Creating microservices
- Standard web backends
- API-first applications

### 📁 Project Structure (35 files)

```
my-api/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── database/
│   │   ├── connection.go    # Database connection
│   │   └── migrations.go    # Database migrations
│   ├── handlers/
│   │   ├── handlers.go      # HTTP handlers
│   │   ├── auth.go         # Authentication handlers
│   │   ├── health.go       # Health check handlers
│   │   └── users.go        # User handlers
│   ├── middleware/
│   │   ├── auth.go         # Authentication middleware
│   │   ├── cors.go         # CORS middleware
│   │   ├── logger.go       # Request logging
│   │   ├── recovery.go     # Panic recovery
│   │   ├── request_id.go   # Request ID generation
│   │   └── security_headers.go # Security headers
│   ├── models/
│   │   ├── base.go         # Base model
│   │   └── user.go         # User model
│   ├── repository/
│   │   ├── interfaces.go   # Repository interfaces
│   │   └── user.go         # User repository
│   ├── services/
│   │   ├── auth.go         # Authentication service
│   │   └── user.go         # User service
│   └── logger/
│       ├── interface.go    # Logger interface
│       └── logger.go       # Logger implementation
├── configs/
│   ├── config.dev.yaml     # Development config
│   ├── config.prod.yaml    # Production config
│   └── config.test.yaml    # Testing config
├── migrations/
│   ├── 001_create_users.up.sql
│   ├── 001_create_users.down.sql
│   └── embed.go            # Embedded migrations
├── api/
│   └── openapi.yaml        # OpenAPI specification
├── tests/
│   ├── integration/
│   │   └── api_test.go     # Integration tests
│   └── unit/
│       └── services_test.go # Unit tests
├── scripts/
│   ├── dev.sh             # Development setup
│   └── migrate.sh         # Migration script
├── docker-compose.yml     # Development environment
├── Dockerfile            # Container definition
├── Makefile             # Build automation
└── README.md            # Documentation
```

### 🔧 Generated Features

#### Web Framework Integration
- **HTTP server** with graceful shutdown
- **Routing** with middleware support
- **Request/response handling**
- **Content negotiation** (JSON, XML)
- **Error handling** with proper HTTP status codes

#### Database Integration
```go
// internal/repository/user.go
type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id string) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, limit, offset int) ([]*models.User, error)
}

type userRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
    query := `
        INSERT INTO users (id, email, name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
    `
    _, err := r.db.ExecContext(ctx, query, 
        user.ID, user.Email, user.Name, user.CreatedAt, user.UpdatedAt)
    return err
}
```

#### Service Layer
```go
// internal/services/user.go
type UserService interface {
    CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, error)
    GetUser(ctx context.Context, id string) (*models.User, error)
    UpdateUser(ctx context.Context, id string, req UpdateUserRequest) (*models.User, error)
    DeleteUser(ctx context.Context, id string) error
    ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error)
}

type userService struct {
    repo   repository.UserRepository
    logger logger.Logger
}

func NewUserService(repo repository.UserRepository, logger logger.Logger) UserService {
    return &userService{
        repo:   repo,
        logger: logger,
    }
}

func (s *userService) CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, error) {
    s.logger.Info("Creating user", "email", req.Email)
    
    user := &models.User{
        ID:        uuid.New().String(),
        Email:     req.Email,
        Name:      req.Name,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    if err := s.repo.Create(ctx, user); err != nil {
        s.logger.Error("Failed to create user", "error", err)
        return nil, err
    }
    
    return user, nil
}
```

#### HTTP Handlers
```go
// internal/handlers/users.go
type UserHandler struct {
    service services.UserService
    logger  logger.Logger
}

func NewUserHandler(service services.UserService, logger logger.Logger) *UserHandler {
    return &UserHandler{
        service: service,
        logger:  logger,
    }
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    user, err := h.service.CreateUser(c.Request.Context(), req)
    if err != nil {
        h.logger.Error("Failed to create user", "error", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    
    c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
    users := r.Group("/users")
    {
        users.POST("", h.CreateUser)
        users.GET("/:id", h.GetUser)
        users.PUT("/:id", h.UpdateUser)
        users.DELETE("/:id", h.DeleteUser)
        users.GET("", h.ListUsers)
    }
}
```

### 🚀 Generation Commands
```bash
# Basic web API
go-starter new my-api --type=web-api

# With PostgreSQL and JWT auth
go-starter new my-api --type=web-api \
  --database-driver=postgres \
  --database-orm=gorm \
  --auth-type=jwt \
  --logger=zap

# With all advanced features
go-starter new my-api --type=web-api \
  --framework=gin \
  --database-driver=postgres \
  --auth-type=jwt \
  --logger=zap \
  --advanced
```

### 📊 Supported Options

#### Frameworks
- **gin** (default): Fast HTTP web framework
- **echo**: High performance, extensible web framework
- **fiber**: Express-inspired web framework
- **chi**: Lightweight, idiomatic HTTP router

#### Databases
- **postgres**: PostgreSQL with connection pooling
- **mysql**: MySQL with optimized settings
- **sqlite**: SQLite for development/testing
- **mongodb**: MongoDB with official driver

#### ORMs/Database Libraries
- **gorm**: Feature-rich ORM
- **sqlx**: Extensions on database/sql
- **sqlc**: Generate type-safe code from SQL
- **ent**: Entity framework for Go

#### Authentication
- **jwt**: JSON Web Token authentication
- **oauth2**: OAuth2 provider integration
- **session**: Session-based authentication
- **api-key**: API key authentication

---

## Web API Clean

> **Best for**: Enterprise applications, complex business logic, testable systems

### 🎯 When to Use
- Enterprise applications
- Complex business logic
- Systems requiring high testability
- Long-term maintainable projects

### 📁 Project Structure (45 files)

```
my-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── adapters/
│   │   ├── primary/             # Driving adapters
│   │   │   ├── http/           # HTTP adapter
│   │   │   │   ├── handlers/   # HTTP handlers
│   │   │   │   ├── middleware/ # HTTP middleware
│   │   │   │   └── server.go   # HTTP server setup
│   │   │   └── cli/            # CLI adapter (if needed)
│   │   └── secondary/          # Driven adapters
│   │       ├── database/       # Database adapter
│   │       │   ├── postgres/   # PostgreSQL implementation
│   │       │   └── memory/     # In-memory implementation
│   │       ├── email/          # Email adapter
│   │       └── cache/          # Cache adapter
│   ├── domain/
│   │   ├── entities/           # Domain entities
│   │   │   ├── user.go        # User entity
│   │   │   └── errors.go      # Domain errors
│   │   ├── repositories/       # Repository interfaces
│   │   │   └── user.go        # User repository interface
│   │   └── services/          # Domain services
│   │       └── user.go        # User domain service
│   ├── usecases/              # Application layer
│   │   ├── interfaces/        # Use case interfaces
│   │   ├── user/             # User use cases
│   │   │   ├── create.go     # Create user use case
│   │   │   ├── get.go        # Get user use case
│   │   │   └── list.go       # List users use case
│   │   └── common/           # Common use case logic
│   └── infrastructure/        # Infrastructure layer
│       ├── config/           # Configuration
│       ├── database/         # Database setup
│       ├── logger/           # Logging setup
│       └── container/        # Dependency injection
├── tests/
│   ├── unit/
│   │   ├── entities_test.go  # Entity tests
│   │   ├── usecases_test.go  # Use case tests
│   │   └── services_test.go  # Service tests
│   ├── integration/
│   │   └── api_test.go       # Integration tests
│   └── mocks/               # Generated mocks
│       ├── mock_user_repository.go
│       └── mock_email_service.go
└── [standard files...]
```

### 🔧 Clean Architecture Principles

#### Domain Layer (Core)
```go
// internal/domain/entities/user.go
type User struct {
    id       UserID
    email    Email
    name     Name
    status   UserStatus
    createdAt time.Time
    updatedAt time.Time
}

func NewUser(email Email, name Name) (*User, error) {
    if err := email.Validate(); err != nil {
        return nil, err
    }
    
    if err := name.Validate(); err != nil {
        return nil, err
    }
    
    return &User{
        id:        NewUserID(),
        email:     email,
        name:      name,
        status:    UserStatusActive,
        createdAt: time.Now(),
        updatedAt: time.Now(),
    }, nil
}

func (u *User) ChangeName(name Name) error {
    if err := name.Validate(); err != nil {
        return err
    }
    
    u.name = name
    u.updatedAt = time.Now()
    return nil
}

// Domain repository interface
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id UserID) (*User, error)
    FindByEmail(ctx context.Context, email Email) (*User, error)
}
```

#### Use Cases (Application Layer)
```go
// internal/usecases/user/create.go
type CreateUserUseCase interface {
    Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error)
}

type CreateUserInput struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}

type CreateUserOutput struct {
    ID        string    `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}

type createUserUseCase struct {
    userRepo      domain.UserRepository
    emailService  EmailService
    logger        logger.Logger
}

func NewCreateUserUseCase(
    userRepo domain.UserRepository,
    emailService EmailService,
    logger logger.Logger,
) CreateUserUseCase {
    return &createUserUseCase{
        userRepo:     userRepo,
        emailService: emailService,
        logger:       logger,
    }
}

func (uc *createUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
    // Validate input
    email, err := domain.NewEmail(input.Email)
    if err != nil {
        return nil, fmt.Errorf("invalid email: %w", err)
    }
    
    name, err := domain.NewName(input.Name)
    if err != nil {
        return nil, fmt.Errorf("invalid name: %w", err)
    }
    
    // Check if user already exists
    existingUser, err := uc.userRepo.FindByEmail(ctx, email)
    if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
        return nil, fmt.Errorf("failed to check existing user: %w", err)
    }
    
    if existingUser != nil {
        return nil, domain.ErrUserAlreadyExists
    }
    
    // Create new user
    user, err := domain.NewUser(email, name)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    // Save user
    if err := uc.userRepo.Save(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to save user: %w", err)
    }
    
    // Send welcome email (async)
    go func() {
        if err := uc.emailService.SendWelcomeEmail(user.Email(), user.Name()); err != nil {
            uc.logger.Error("Failed to send welcome email", "error", err, "user_id", user.ID())
        }
    }()
    
    return &CreateUserOutput{
        ID:        user.ID().String(),
        Email:     user.Email().String(),
        Name:      user.Name().String(),
        Status:    user.Status().String(),
        CreatedAt: user.CreatedAt(),
    }, nil
}
```

#### Adapters (Infrastructure Layer)
```go
// internal/adapters/secondary/database/postgres/user.go
type userRepository struct {
    db     *sql.DB
    logger logger.Logger
}

func NewUserRepository(db *sql.DB, logger logger.Logger) domain.UserRepository {
    return &userRepository{
        db:     db,
        logger: logger,
    }
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (id, email, name, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (id) DO UPDATE SET
            email = EXCLUDED.email,
            name = EXCLUDED.name,
            status = EXCLUDED.status,
            updated_at = EXCLUDED.updated_at
    `
    
    _, err := r.db.ExecContext(ctx, query,
        user.ID().String(),
        user.Email().String(),
        user.Name().String(),
        user.Status().String(),
        user.CreatedAt(),
        user.UpdatedAt(),
    )
    
    if err != nil {
        r.logger.Error("Failed to save user", "error", err, "user_id", user.ID())
        return fmt.Errorf("failed to save user: %w", err)
    }
    
    return nil
}
```

### 🚀 Generation Command
```bash
go-starter new my-api --type=web-api --architecture=clean \
  --database-driver=postgres \
  --auth-type=jwt \
  --logger=zap
```

### ✅ Benefits
- **Clear separation of concerns**
- **Highly testable architecture**
- **Independent of frameworks and databases**
- **Business logic isolation**
- **Dependency inversion**

---

## Web API DDD

> **Best for**: Complex domains, rich business models, event-driven systems

### 🎯 When to Use
- Complex business domains
- Rich domain models
- Event-driven architectures
- Domain expertise is critical

### 📁 Project Structure (50 files)

```
my-api/
├── internal/
│   ├── domain/
│   │   ├── user/                    # User aggregate
│   │   │   ├── entity.go           # User entity
│   │   │   ├── value_objects.go    # Email, Name, etc.
│   │   │   ├── repository.go       # Repository interface
│   │   │   ├── service.go          # Domain service
│   │   │   ├── specifications.go   # Business rules
│   │   │   ├── events.go          # Domain events
│   │   │   └── errors.go          # Domain errors
│   │   ├── order/                  # Order aggregate
│   │   │   ├── entity.go
│   │   │   ├── order_item.go      # Child entity
│   │   │   ├── repository.go
│   │   │   ├── service.go
│   │   │   └── specifications.go
│   │   ├── shared/                 # Shared domain concepts
│   │   │   ├── value_objects/     # Common value objects
│   │   │   ├── events/            # Event infrastructure
│   │   │   └── specifications/    # Base specifications
│   │   └── common/                # Domain primitives
│   ├── application/               # Application services
│   │   ├── user/
│   │   │   ├── commands/          # Command handlers
│   │   │   ├── queries/           # Query handlers
│   │   │   ├── dto/              # Data transfer objects
│   │   │   └── handlers/         # Application event handlers
│   │   ├── order/
│   │   └── shared/               # Shared application logic
│   ├── infrastructure/           # Technical implementations
│   │   ├── persistence/
│   │   │   ├── user/            # User repository implementation
│   │   │   └── order/           # Order repository implementation
│   │   ├── messaging/           # Event publishing
│   │   ├── external/            # External service adapters
│   │   └── config/              # Configuration
│   └── presentation/            # HTTP layer
│       ├── http/
│       │   ├── handlers/        # HTTP handlers
│       │   ├── middleware/      # HTTP middleware
│       │   └── dto/            # HTTP DTOs
│       └── grpc/               # gRPC handlers (optional)
└── [standard files...]
```

### 🔧 DDD Implementation

#### Domain Entity with Business Logic
```go
// internal/domain/user/entity.go
type User struct {
    id           UserID
    email        Email
    profile      Profile
    subscription Subscription
    status       UserStatus
    domainEvents []shared.DomainEvent
    version      int
    createdAt    time.Time
    updatedAt    time.Time
}

func NewUser(email Email, profile Profile) (*User, error) {
    // Business rule: Email must be unique
    spec := NewUniqueEmailSpecification()
    if !spec.IsSatisfiedBy(email) {
        return nil, ErrEmailAlreadyExists
    }
    
    user := &User{
        id:           NewUserID(),
        email:        email,
        profile:      profile,
        subscription: NewBasicSubscription(),
        status:       UserStatusActive,
        domainEvents: []shared.DomainEvent{},
        version:      1,
        createdAt:    time.Now(),
        updatedAt:    time.Now(),
    }
    
    // Raise domain event
    user.RaiseDomainEvent(NewUserRegisteredEvent(user.id, user.email))
    
    return user, nil
}

func (u *User) UpgradeSubscription(newPlan SubscriptionPlan) error {
    // Business rule: Can only upgrade, not downgrade
    if newPlan.Level() <= u.subscription.Plan().Level() {
        return ErrCannotDowngradeSubscription
    }
    
    // Business rule: Active users only
    if u.status != UserStatusActive {
        return ErrInactiveUserCannotUpgrade
    }
    
    oldPlan := u.subscription.Plan()
    u.subscription = NewSubscription(newPlan)
    u.updatedAt = time.Now()
    u.version++
    
    // Raise domain event
    u.RaiseDomainEvent(NewSubscriptionUpgradedEvent(
        u.id, oldPlan, newPlan, time.Now(),
    ))
    
    return nil
}

func (u *User) RaiseDomainEvent(event shared.DomainEvent) {
    u.domainEvents = append(u.domainEvents, event)
}

func (u *User) DomainEvents() []shared.DomainEvent {
    return u.domainEvents
}

func (u *User) ClearDomainEvents() {
    u.domainEvents = []shared.DomainEvent{}
}
```

#### Value Objects
```go
// internal/domain/user/value_objects.go
type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    if value == "" {
        return Email{}, ErrInvalidEmail
    }
    
    if !isValidEmail(value) {
        return Email{}, ErrInvalidEmail
    }
    
    return Email{value: strings.ToLower(value)}, nil
}

func (e Email) String() string {
    return e.value
}

func (e Email) Equals(other Email) bool {
    return e.value == other.value
}

type Profile struct {
    firstName FirstName
    lastName  LastName
    avatar    Avatar
}

func NewProfile(firstName, lastName string, avatar string) (Profile, error) {
    fn, err := NewFirstName(firstName)
    if err != nil {
        return Profile{}, err
    }
    
    ln, err := NewLastName(lastName)
    if err != nil {
        return Profile{}, err
    }
    
    av, err := NewAvatar(avatar)
    if err != nil {
        return Profile{}, err
    }
    
    return Profile{
        firstName: fn,
        lastName:  ln,
        avatar:    av,
    }, nil
}

func (p Profile) FullName() string {
    return fmt.Sprintf("%s %s", p.firstName.String(), p.lastName.String())
}
```

#### Specifications (Business Rules)
```go
// internal/domain/user/specifications.go
type Specification interface {
    IsSatisfiedBy(user *User) bool
    Reason() string
}

type CanUpgradeSubscriptionSpecification struct {
    userRepo UserRepository
}

func NewCanUpgradeSubscriptionSpecification(userRepo UserRepository) Specification {
    return &CanUpgradeSubscriptionSpecification{
        userRepo: userRepo,
    }
}

func (s *CanUpgradeSubscriptionSpecification) IsSatisfiedBy(user *User) bool {
    // Business rule: User must be active
    if user.Status() != UserStatusActive {
        return false
    }
    
    // Business rule: No outstanding payments
    if user.HasOutstandingPayments() {
        return false
    }
    
    // Business rule: Not already on highest tier
    if user.Subscription().Plan().IsHighestTier() {
        return false
    }
    
    return true
}

func (s *CanUpgradeSubscriptionSpecification) Reason() string {
    return "User must be active with no outstanding payments and not on highest tier"
}

// Composite specifications
type AndSpecification struct {
    left, right Specification
}

func (s *AndSpecification) IsSatisfiedBy(user *User) bool {
    return s.left.IsSatisfiedBy(user) && s.right.IsSatisfiedBy(user)
}

func And(left, right Specification) Specification {
    return &AndSpecification{left: left, right: right}
}
```

#### Domain Events
```go
// internal/domain/user/events.go
type UserRegisteredEvent struct {
    UserID    UserID    `json:"user_id"`
    Email     string    `json:"email"`
    Timestamp time.Time `json:"timestamp"`
}

func NewUserRegisteredEvent(userID UserID, email Email) shared.DomainEvent {
    return &UserRegisteredEvent{
        UserID:    userID,
        Email:     email.String(),
        Timestamp: time.Now(),
    }
}

func (e *UserRegisteredEvent) EventType() string {
    return "user.registered"
}

func (e *UserRegisteredEvent) AggregateID() string {
    return e.UserID.String()
}

func (e *UserRegisteredEvent) OccurredAt() time.Time {
    return e.Timestamp
}

type SubscriptionUpgradedEvent struct {
    UserID      UserID           `json:"user_id"`
    OldPlan     SubscriptionPlan `json:"old_plan"`
    NewPlan     SubscriptionPlan `json:"new_plan"`
    UpgradedAt  time.Time        `json:"upgraded_at"`
}

func NewSubscriptionUpgradedEvent(userID UserID, oldPlan, newPlan SubscriptionPlan, upgradedAt time.Time) shared.DomainEvent {
    return &SubscriptionUpgradedEvent{
        UserID:     userID,
        OldPlan:    oldPlan,
        NewPlan:    newPlan,
        UpgradedAt: upgradedAt,
    }
}
```

#### Domain Service
```go
// internal/domain/user/service.go
type UserDomainService interface {
    CanUserUpgradeSubscription(ctx context.Context, user *User, newPlan SubscriptionPlan) error
    ProcessSubscriptionUpgrade(ctx context.Context, user *User, newPlan SubscriptionPlan) error
}

type userDomainService struct {
    userRepo        UserRepository
    paymentService  PaymentService
    planRepo        SubscriptionPlanRepository
}

func NewUserDomainService(
    userRepo UserRepository,
    paymentService PaymentService,
    planRepo SubscriptionPlanRepository,
) UserDomainService {
    return &userDomainService{
        userRepo:       userRepo,
        paymentService: paymentService,
        planRepo:       planRepo,
    }
}

func (s *userDomainService) CanUserUpgradeSubscription(ctx context.Context, user *User, newPlan SubscriptionPlan) error {
    // Check business rules using specifications
    canUpgradeSpec := NewCanUpgradeSubscriptionSpecification(s.userRepo)
    if !canUpgradeSpec.IsSatisfiedBy(user) {
        return fmt.Errorf("cannot upgrade subscription: %s", canUpgradeSpec.Reason())
    }
    
    // Check payment capability
    if !s.paymentService.CanProcessPayment(ctx, user.ID(), newPlan.Price()) {
        return ErrInsufficientFunds
    }
    
    return nil
}

func (s *userDomainService) ProcessSubscriptionUpgrade(ctx context.Context, user *User, newPlan SubscriptionPlan) error {
    // Domain logic for complex subscription upgrade
    if err := s.CanUserUpgradeSubscription(ctx, user, newPlan); err != nil {
        return err
    }
    
    // Process payment
    payment, err := s.paymentService.ProcessPayment(ctx, user.ID(), newPlan.Price())
    if err != nil {
        return fmt.Errorf("payment failed: %w", err)
    }
    
    // Upgrade subscription
    if err := user.UpgradeSubscription(newPlan); err != nil {
        // Rollback payment
        s.paymentService.RefundPayment(ctx, payment.ID())
        return err
    }
    
    return nil
}
```

### 🚀 Generation Command
```bash
go-starter new my-api --type=web-api --architecture=ddd \
  --database-driver=postgres \
  --auth-type=jwt \
  --logger=zap \
  --advanced
```

### ✅ DDD Benefits
- **Rich domain models** with business logic
- **Clear business rules** through specifications
- **Event-driven architecture** with domain events
- **Ubiquitous language** shared with domain experts
- **Complex business logic** properly encapsulated

---

## Web API Hexagonal

> **Best for**: Highly testable systems, multiple adapters, ports & adapters pattern

### 🎯 When to Use
- Maximum testability required
- Multiple interface types (HTTP, CLI, gRPC)
- Frequent adapter changes
- Complex integration requirements

### 📁 Project Structure (55 files)

```
my-api/
├── internal/
│   ├── adapters/
│   │   ├── primary/                # Driving adapters (inbound)
│   │   │   ├── http/              # HTTP adapter
│   │   │   │   ├── handlers/      # HTTP handlers
│   │   │   │   ├── middleware/    # HTTP middleware
│   │   │   │   ├── mappers/       # DTO mappers
│   │   │   │   └── server.go      # HTTP server
│   │   │   ├── grpc/              # gRPC adapter
│   │   │   │   ├── handlers/      # gRPC handlers
│   │   │   │   └── server.go      # gRPC server
│   │   │   └── cli/               # CLI adapter
│   │   │       └── commands/      # CLI commands
│   │   └── secondary/             # Driven adapters (outbound)
│   │       ├── database/          # Database adapters
│   │       │   ├── postgres/      # PostgreSQL adapter
│   │       │   ├── mongodb/       # MongoDB adapter
│   │       │   └── memory/        # In-memory adapter
│   │       ├── email/             # Email adapters
│   │       │   ├── smtp/          # SMTP adapter
│   │       │   ├── sendgrid/      # SendGrid adapter
│   │       │   └── mock/          # Mock adapter
│   │       ├── cache/             # Cache adapters
│   │       │   ├── redis/         # Redis adapter
│   │       │   └── memory/        # In-memory cache
│   │       ├── messaging/         # Message queue adapters
│   │       │   ├── rabbitmq/      # RabbitMQ adapter
│   │       │   └── sqs/           # AWS SQS adapter
│   │       └── external/          # External service adapters
│   │           ├── payment/       # Payment service adapter
│   │           └── notification/  # Notification service adapter
│   ├── application/               # Application layer (use cases)
│   │   ├── ports/                 # Port interfaces
│   │   │   ├── primary/           # Primary ports (inbound)
│   │   │   │   ├── user_service.go
│   │   │   │   └── order_service.go
│   │   │   └── secondary/         # Secondary ports (outbound)
│   │   │       ├── user_repository.go
│   │   │       ├── email_service.go
│   │   │       ├── cache_service.go
│   │   │       └── payment_service.go
│   │   ├── services/              # Application services
│   │   │   ├── user_service.go    # User application service
│   │   │   └── order_service.go   # Order application service
│   │   ├── commands/              # Command objects
│   │   │   ├── create_user.go
│   │   │   └── place_order.go
│   │   ├── queries/               # Query objects
│   │   │   ├── get_user.go
│   │   │   └── list_orders.go
│   │   └── events/                # Application events
│   │       ├── user_created.go
│   │       └── order_placed.go
│   └── domain/                    # Domain layer (core business logic)
│       ├── entities/              # Domain entities
│       │   ├── user.go
│       │   └── order.go
│       ├── value_objects/         # Value objects
│       │   ├── email.go
│       │   ├── money.go
│       │   └── address.go
│       ├── aggregates/            # Aggregate roots
│       │   ├── user_aggregate.go
│       │   └── order_aggregate.go
│       ├── specifications/        # Business rules
│       │   ├── user_specs.go
│       │   └── order_specs.go
│       └── events/                # Domain events
│           ├── user_events.go
│           └── order_events.go
└── [test structure mirrors implementation...]
```

### 🔧 Hexagonal Architecture Implementation

#### Primary Ports (Inbound)
```go
// internal/application/ports/primary/user_service.go
type UserService interface {
    CreateUser(ctx context.Context, cmd CreateUserCommand) (*UserResponse, error)
    GetUser(ctx context.Context, query GetUserQuery) (*UserResponse, error)
    UpdateUser(ctx context.Context, cmd UpdateUserCommand) (*UserResponse, error)
    DeleteUser(ctx context.Context, cmd DeleteUserCommand) error
    ListUsers(ctx context.Context, query ListUsersQuery) (*UsersResponse, error)
}

type CreateUserCommand struct {
    Email     string            `json:"email"`
    Name      string            `json:"name"`
    Metadata  map[string]string `json:"metadata,omitempty"`
}

type GetUserQuery struct {
    ID string `json:"id"`
}

type UserResponse struct {
    ID        string            `json:"id"`
    Email     string            `json:"email"`
    Name      string            `json:"name"`
    Status    string            `json:"status"`
    Metadata  map[string]string `json:"metadata"`
    CreatedAt time.Time         `json:"created_at"`
    UpdatedAt time.Time         `json:"updated_at"`
}
```

#### Secondary Ports (Outbound)
```go
// internal/application/ports/secondary/user_repository.go
type UserRepository interface {
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id domain.UserID) (*domain.User, error)
    FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error)
    FindAll(ctx context.Context, limit, offset int) ([]*domain.User, error)
    Delete(ctx context.Context, id domain.UserID) error
    Count(ctx context.Context) (int, error)
}

// internal/application/ports/secondary/email_service.go
type EmailService interface {
    SendWelcomeEmail(ctx context.Context, email string, name string) error
    SendPasswordResetEmail(ctx context.Context, email string, token string) error
    SendNotificationEmail(ctx context.Context, email string, subject string, body string) error
}

// internal/application/ports/secondary/cache_service.go
type CacheService interface {
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Get(ctx context.Context, key string, dest interface{}) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
}
```

#### Application Service (Use Cases)
```go
// internal/application/services/user_service.go
type userService struct {
    userRepo     ports.UserRepository
    emailService ports.EmailService
    cacheService ports.CacheService
    logger       logger.Logger
    eventBus     ports.EventBus
}

func NewUserService(
    userRepo ports.UserRepository,
    emailService ports.EmailService,
    cacheService ports.CacheService,
    logger logger.Logger,
    eventBus ports.EventBus,
) ports.UserService {
    return &userService{
        userRepo:     userRepo,
        emailService: emailService,
        cacheService: cacheService,
        logger:       logger,
        eventBus:     eventBus,
    }
}

func (s *userService) CreateUser(ctx context.Context, cmd ports.CreateUserCommand) (*ports.UserResponse, error) {
    s.logger.Info("Creating user", "email", cmd.Email)
    
    // Create domain objects
    email, err := domain.NewEmail(cmd.Email)
    if err != nil {
        return nil, fmt.Errorf("invalid email: %w", err)
    }
    
    name, err := domain.NewName(cmd.Name)
    if err != nil {
        return nil, fmt.Errorf("invalid name: %w", err)
    }
    
    // Check if user already exists
    existingUser, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
        return nil, fmt.Errorf("failed to check existing user: %w", err)
    }
    
    if existingUser != nil {
        return nil, domain.ErrUserAlreadyExists
    }
    
    // Create new user
    user := domain.NewUser(email, name)
    
    // Save user
    if err := s.userRepo.Save(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to save user: %w", err)
    }
    
    // Cache user
    cacheKey := fmt.Sprintf("user:%s", user.ID().String())
    if err := s.cacheService.Set(ctx, cacheKey, user, time.Hour); err != nil {
        s.logger.Warn("Failed to cache user", "error", err, "user_id", user.ID())
    }
    
    // Send welcome email (async)
    go func() {
        if err := s.emailService.SendWelcomeEmail(context.Background(), user.Email().String(), user.Name().String()); err != nil {
            s.logger.Error("Failed to send welcome email", "error", err, "user_id", user.ID())
        }
    }()
    
    // Publish domain events
    for _, event := range user.DomainEvents() {
        if err := s.eventBus.Publish(ctx, event); err != nil {
            s.logger.Error("Failed to publish event", "error", err, "event_type", event.Type())
        }
    }
    user.ClearDomainEvents()
    
    return s.mapUserToResponse(user), nil
}
```

#### Primary Adapter (HTTP)
```go
// internal/adapters/primary/http/handlers/user_handler.go
type UserHandler struct {
    userService ports.UserService
    logger      logger.Logger
}

func NewUserHandler(userService ports.UserService, logger logger.Logger) *UserHandler {
    return &UserHandler{
        userService: userService,
        logger:      logger,
    }
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
        return
    }
    
    // Map HTTP request to command
    cmd := ports.CreateUserCommand{
        Email:    req.Email,
        Name:     req.Name,
        Metadata: req.Metadata,
    }
    
    user, err := h.userService.CreateUser(c.Request.Context(), cmd)
    if err != nil {
        h.handleError(c, err)
        return
    }
    
    c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) handleError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, domain.ErrUserAlreadyExists):
        c.JSON(http.StatusConflict, ErrorResponse{Error: "User already exists"})
    case errors.Is(err, domain.ErrInvalidEmail):
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid email format"})
    case errors.Is(err, domain.ErrInvalidName):
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid name format"})
    default:
        h.logger.Error("Internal server error", "error", err)
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
    }
}
```

#### Secondary Adapter (Database)
```go
// internal/adapters/secondary/database/postgres/user_repository.go
type userRepository struct {
    db     *sql.DB
    logger logger.Logger
}

func NewUserRepository(db *sql.DB, logger logger.Logger) ports.UserRepository {
    return &userRepository{
        db:     db,
        logger: logger,
    }
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (id, email, name, status, metadata, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO UPDATE SET
            email = EXCLUDED.email,
            name = EXCLUDED.name,
            status = EXCLUDED.status,
            metadata = EXCLUDED.metadata,
            updated_at = EXCLUDED.updated_at
    `
    
    metadataJSON, err := json.Marshal(user.Metadata())
    if err != nil {
        return fmt.Errorf("failed to marshal metadata: %w", err)
    }
    
    _, err = r.db.ExecContext(ctx, query,
        user.ID().String(),
        user.Email().String(),
        user.Name().String(),
        user.Status().String(),
        metadataJSON,
        user.CreatedAt(),
        user.UpdatedAt(),
    )
    
    if err != nil {
        r.logger.Error("Failed to save user", "error", err, "user_id", user.ID())
        return fmt.Errorf("failed to save user: %w", err)
    }
    
    return nil
}

func (r *userRepository) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
    query := `
        SELECT id, email, name, status, metadata, created_at, updated_at
        FROM users
        WHERE id = $1
    `
    
    var (
        userID       string
        email        string
        name         string
        status       string
        metadataJSON []byte
        createdAt    time.Time
        updatedAt    time.Time
    )
    
    err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
        &userID, &email, &name, &status, &metadataJSON, &createdAt, &updatedAt,
    )
    
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domain.ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    
    var metadata map[string]string
    if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
        return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
    }
    
    // Reconstruct domain object
    domainEmail, _ := domain.NewEmail(email)
    domainName, _ := domain.NewName(name)
    domainStatus := domain.UserStatus(status)
    
    user := domain.ReconstructUser(
        domain.UserID(userID),
        domainEmail,
        domainName,
        domainStatus,
        metadata,
        createdAt,
        updatedAt,
    )
    
    return user, nil
}
```

### 🚀 Generation Command
```bash
go-starter new my-api --type=web-api --architecture=hexagonal \
  --database-driver=postgres \
  --auth-type=jwt \
  --logger=zap \
  --advanced
```

### ✅ Hexagonal Benefits
- **Multiple adapters** support (HTTP, gRPC, CLI)
- **Easy testing** with port interfaces
- **Adapter swapping** without core changes
- **Clear boundaries** between layers
- **Technology independence**

---

## Lambda Standard

> **Best for**: Event handlers, background processing, simple serverless functions

### 🎯 When to Use
- AWS event processing
- Webhooks and callbacks
- Scheduled tasks
- Simple serverless functions

### 📁 Project Structure (12 files)

```
my-lambda/
├── main.go                      # Lambda entry point
├── handler.go                   # Business logic
├── template.yaml                # SAM template
├── Makefile                    # Build and deploy automation
├── README.md                   # Deployment instructions
├── deploy.sh                   # Deployment script
├── go.mod                      # Module definition
├── internal/
│   ├── logger/
│   │   └── logger.go           # Lambda-optimized logger
│   └── observability/
│       ├── cloudwatch.go       # CloudWatch metrics
│       ├── metrics.go          # Custom metrics
│       └── tracing.go          # X-Ray tracing
└── template.yaml.tmpl          # SAM template (if configurable)
```

### 🔧 Lambda Implementation

#### Main Handler
```go
// main.go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "github.com/aws/aws-xray-sdk-go/xray"
    
    "my-lambda/internal/logger"
    "my-lambda/internal/observability"
)

var (
    log     *logger.Logger
    metrics *observability.Metrics
)

func init() {
    log = logger.New(logger.Config{
        Level:  "info",
        Format: "json",
    })
    
    metrics = observability.NewMetrics()
}

func main() {
    lambda.Start(xray.LambdaHandler(handleRequest))
}

func handleRequest(ctx context.Context, event events.CloudWatchEvent) error {
    log.Info("Processing CloudWatch event", 
        "source", event.Source,
        "detail_type", event.DetailType,
        "region", event.Region,
    )
    
    // Start tracing segment
    _, seg := xray.BeginSubsegment(ctx, "process-event")
    defer seg.Close(nil)
    
    // Record custom metric
    metrics.IncrementCounter("events_processed", map[string]string{
        "source": event.Source,
        "region": event.Region,
    })
    
    // Process the event
    if err := processEvent(ctx, event); err != nil {
        log.Error("Failed to process event", "error", err)
        metrics.IncrementCounter("events_failed", map[string]string{
            "source": event.Source,
            "error":  err.Error(),
        })
        return err
    }
    
    metrics.IncrementCounter("events_succeeded", map[string]string{
        "source": event.Source,
    })
    
    log.Info("Event processed successfully")
    return nil
}
```

#### Business Logic
```go
// handler.go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-xray-sdk-go/xray"
)

type EventProcessor interface {
    ProcessEvent(ctx context.Context, event events.CloudWatchEvent) error
}

type eventProcessor struct {
    log     *logger.Logger
    metrics *observability.Metrics
}

func NewEventProcessor(log *logger.Logger, metrics *observability.Metrics) EventProcessor {
    return &eventProcessor{
        log:     log,
        metrics: metrics,
    }
}

func processEvent(ctx context.Context, event events.CloudWatchEvent) error {
    processor := NewEventProcessor(log, metrics)
    return processor.ProcessEvent(ctx, event)
}

func (p *eventProcessor) ProcessEvent(ctx context.Context, event events.CloudWatchEvent) error {
    // Start X-Ray subsegment
    _, seg := xray.BeginSubsegment(ctx, "event-processing")
    defer seg.Close(nil)
    
    startTime := time.Now()
    defer func() {
        duration := time.Since(startTime)
        p.metrics.RecordDuration("event_processing_duration", duration, map[string]string{
            "source": event.Source,
        })
    }()
    
    switch event.Source {
    case "aws.s3":
        return p.processS3Event(ctx, event)
    case "aws.dynamodb":
        return p.processDynamoDBEvent(ctx, event)
    case "custom.application":
        return p.processCustomEvent(ctx, event)
    default:
        return fmt.Errorf("unsupported event source: %s", event.Source)
    }
}

func (p *eventProcessor) processS3Event(ctx context.Context, event events.CloudWatchEvent) error {
    p.log.Info("Processing S3 event", "detail_type", event.DetailType)
    
    // Parse S3 event details
    var s3Detail map[string]interface{}
    if err := json.Unmarshal(event.Detail, &s3Detail); err != nil {
        return fmt.Errorf("failed to parse S3 event detail: %w", err)
    }
    
    bucket, ok := s3Detail["bucket"].(map[string]interface{})
    if !ok {
        return fmt.Errorf("invalid S3 event format: missing bucket")
    }
    
    bucketName, ok := bucket["name"].(string)
    if !ok {
        return fmt.Errorf("invalid S3 event format: missing bucket name")
    }
    
    p.log.Info("Processing S3 bucket event", "bucket", bucketName)
    
    // Add your S3-specific business logic here
    // Example: Process uploaded files, trigger workflows, etc.
    
    return nil
}

func (p *eventProcessor) processDynamoDBEvent(ctx context.Context, event events.CloudWatchEvent) error {
    p.log.Info("Processing DynamoDB event", "detail_type", event.DetailType)
    
    // Parse DynamoDB event details
    var dynamoDetail map[string]interface{}
    if err := json.Unmarshal(event.Detail, &dynamoDetail); err != nil {
        return fmt.Errorf("failed to parse DynamoDB event detail: %w", err)
    }
    
    // Add your DynamoDB-specific business logic here
    // Example: Process stream records, update indexes, etc.
    
    return nil
}

func (p *eventProcessor) processCustomEvent(ctx context.Context, event events.CloudWatchEvent) error {
    p.log.Info("Processing custom event", "detail_type", event.DetailType)
    
    // Parse custom event details
    var customDetail map[string]interface{}
    if err := json.Unmarshal(event.Detail, &customDetail); err != nil {
        return fmt.Errorf("failed to parse custom event detail: %w", err)
    }
    
    // Add your custom business logic here
    
    return nil
}
```

#### SAM Template
```yaml
# template.yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31
Description: >
  my-lambda
  
  Lambda function for processing CloudWatch events

Globals:
  Function:
    Timeout: 30
    Runtime: go1.x
    Environment:
      Variables:
        LOG_LEVEL: info
        _X_AMZN_TRACE_ID: !Ref AWS::NoValue

Resources:
  EventProcessorFunction:
    Type: AWS::Serverless::Function
    Properties:
      CodeUri: bin/
      Handler: main
      Architectures:
        - x86_64
      Events:
        CloudWatchRule:
          Type: CloudWatchEvent
          Properties:
            Pattern:
              source:
                - "aws.s3"
                - "aws.dynamodb"
                - "custom.application"
      Environment:
        Variables:
          LOG_LEVEL: !Ref LogLevel
      Tracing: Active
      Policies:
        - CloudWatchPutMetricPolicy: {}
        - AWSXRayDaemonWriteAccess
        - Version: "2012-10-17"
          Statement:
            - Effect: Allow
              Action:
                - logs:CreateLogGroup
                - logs:CreateLogStream
                - logs:PutLogEvents
              Resource: !Sub "arn:aws:logs:${AWS::Region}:${AWS::AccountId}:*"

Parameters:
  LogLevel:
    Type: String
    Default: info
    AllowedValues:
      - debug
      - info
      - warn
      - error
    Description: Log level for the Lambda function

Outputs:
  EventProcessorFunction:
    Description: "Event Processor Lambda Function ARN"
    Value: !GetAtt EventProcessorFunction.Arn
  EventProcessorFunctionIamRole:
    Description: "Implicit IAM Role created for Event Processor function"
    Value: !GetAtt EventProcessorFunctionRole.Arn
```

### 🚀 Generation Command
```bash
go-starter new my-lambda --type=lambda --logger=slog
```

### 📊 Deployment
```bash
# Build and deploy
make build
make deploy

# Or use SAM CLI directly
sam build
sam deploy --guided
```

---

## Summary

Each template is designed for specific use cases and complexity levels. The unified logger interface and consistent project structure make it easy to:

1. **Start with the right template** for your use case
2. **Scale complexity** as your project grows
3. **Switch between similar templates** when requirements change
4. **Maintain consistency** across different project types

### Quick Selection Guide

| Need | Template | Command |
|------|----------|---------|
| **Learning CLI** | CLI Simple | `--type=cli --complexity=simple` |
| **Production CLI** | CLI Standard | `--type=cli --complexity=standard` |
| **Simple API** | Web API Standard | `--type=web-api` |
| **Enterprise API** | Web API Clean | `--type=web-api --architecture=clean` |
| **Complex Domain** | Web API DDD | `--type=web-api --architecture=ddd` |
| **Maximum Testing** | Web API Hexagonal | `--type=web-api --architecture=hexagonal` |
| **Event Processing** | Lambda Standard | `--type=lambda` |
| **API on Lambda** | Lambda Proxy | `--type=lambda-proxy` |

For detailed implementation examples and best practices, refer to the generated README.md in each project.