# Schema-First Development Approach

## Overview

The GraphQL API blueprint adopts a schema-first development approach, where the GraphQL schema is designed before implementation. This approach provides several benefits:

- **Contract-driven development**: API contract is defined upfront
- **Type safety**: Generated code ensures type correctness
- **Documentation**: Schema serves as living documentation
- **Code generation**: Reduces boilerplate and maintains consistency
- **Team collaboration**: Frontend and backend teams can work in parallel

## Schema Design Principles

### 1. Business Domain Focus
Design the schema around business concepts, not database tables:

```graphql
# Good: Business-focused design
type User {
    id: ID!
    profile: UserProfile!
    orders: [Order!]!
    preferences: UserPreferences!
}

type UserProfile {
    name: String!
    email: String!
    avatar: String
}

# Avoid: Database-centric design
type User {
    id: ID!
    name: String!
    email: String!
    avatar_url: String
    created_at: DateTime!
    updated_at: DateTime!
}
```

### 2. Nullable vs Non-Nullable Fields
Be explicit about field nullability:

```graphql
type Product {
    id: ID!              # Always present
    name: String!        # Required field
    description: String  # Optional field
    price: Money!        # Required business data
    reviews: [Review!]!  # Array is never null, but can be empty
    tags: [String!]      # Array can be null, but elements are never null
}
```

### 3. Pagination Patterns
Implement consistent pagination using Relay-style connections:

```graphql
type Query {
    products(
        first: Int
        after: String
        last: Int
        before: String
        filter: ProductFilter
    ): ProductConnection!
}

type ProductConnection {
    edges: [ProductEdge!]!
    pageInfo: PageInfo!
    totalCount: Int!
}

type ProductEdge {
    node: Product!
    cursor: String!
}

type PageInfo {
    hasNextPage: Boolean!
    hasPreviousPage: Boolean!
    startCursor: String
    endCursor: String
}
```

### 4. Input Types and Validation
Use input types for mutations and add validation directives:

```graphql
input CreateUserInput {
    name: String! @validate(min: 2, max: 50)
    email: String! @validate(format: "email")
    age: Int @validate(min: 18, max: 120)
}

input UpdateUserInput {
    name: String @validate(min: 2, max: 50)
    email: String @validate(format: "email")
    age: Int @validate(min: 18, max: 120)
}

type Mutation {
    createUser(input: CreateUserInput!): User!
    updateUser(id: ID!, input: UpdateUserInput!): User!
}
```

## Code Generation Workflow

### 1. Schema Definition
Create the GraphQL schema file:

```graphql
# schema/schema.graphql
directive @auth(requires: Role = USER) on FIELD_DEFINITION
directive @validate(format: String, min: Int, max: Int) on INPUT_FIELD_DEFINITION

scalar DateTime
scalar Money

enum Role {
    USER
    ADMIN
    MODERATOR
}

type Query {
    users: [User!]! @auth
    user(id: ID!): User @auth
    products(filter: ProductFilter): [Product!]!
}

type Mutation {
    createUser(input: CreateUserInput!): User! @auth(requires: ADMIN)
    updateUser(id: ID!, input: UpdateUserInput!): User! @auth
}

type Subscription {
    userUpdated(id: ID!): User! @auth
}

type User {
    id: ID!
    name: String!
    email: String!
    role: Role!
    createdAt: DateTime!
    orders: [Order!]!
}

type Product {
    id: ID!
    name: String!
    description: String
    price: Money!
    inStock: Boolean!
}

type Order {
    id: ID!
    user: User!
    products: [Product!]!
    total: Money!
    status: OrderStatus!
}

enum OrderStatus {
    PENDING
    CONFIRMED
    SHIPPED
    DELIVERED
    CANCELLED
}

input CreateUserInput {
    name: String! @validate(min: 2, max: 50)
    email: String! @validate(format: "email")
    role: Role = USER
}

input UpdateUserInput {
    name: String @validate(min: 2, max: 50)
    email: String @validate(format: "email")
}

input ProductFilter {
    name: String
    minPrice: Money
    maxPrice: Money
    inStock: Boolean
}
```

### 2. gqlgen Configuration
Configure code generation with `gqlgen.yml`:

```yaml
# gqlgen.yml.tmpl
schema:
  - schema/*.graphql

exec:
  filename: internal/generated/generated.go
  package: generated

model:
  filename: internal/model/models_gen.go
  package: model

resolver:
  layout: follow-schema
  dir: internal/resolvers
  package: resolvers
  filename_template: "{{ .PackageName }}.resolvers.go"

# Custom scalar mappings
models:
  DateTime:
    model: time.Time
  Money:
    model: "{{.ModulePath}}/internal/types.Money"

# Skip generating models for external types
  User:
    model: "{{.ModulePath}}/internal/domain.User"
  Product:
    model: "{{.ModulePath}}/internal/domain.Product"

# Directive implementations
directives:
  auth:
    skip_runtime: false
  validate:
    skip_runtime: false
```

### 3. Custom Types and Models
Define custom scalar types:

```go
// internal/types/money.go.tmpl
package types

import (
    "fmt"
    "io"
    "strconv"
    
    "github.com/99designs/gqlgen/graphql"
)

type Money struct {
    Amount   int64  `json:"amount"`   // Amount in cents
    Currency string `json:"currency"` // ISO 4217 currency code
}

func (m Money) String() string {
    return fmt.Sprintf("%.2f %s", float64(m.Amount)/100, m.Currency)
}

// MarshalGQL implements the graphql.Marshaler interface
func (m Money) MarshalGQL(w io.Writer) {
    fmt.Fprintf(w, `"%s"`, m.String())
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (m *Money) UnmarshalGQL(v interface{}) error {
    s, ok := v.(string)
    if !ok {
        return fmt.Errorf("money must be a string")
    }
    
    // Parse money string (e.g., "19.99 USD")
    // Implementation details...
    
    return nil
}
```

### 4. Domain Models
Define domain entities separate from GraphQL models:

```go
// internal/domain/user.go.tmpl
package domain

import (
    "time"
)

type User struct {
    ID        string    `db:"id" json:"id"`
    Name      string    `db:"name" json:"name"`
    Email     string    `db:"email" json:"email"`
    Role      Role      `db:"role" json:"role"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Role string

const (
    RoleUser      Role = "USER"
    RoleAdmin     Role = "ADMIN"
    RoleModerator Role = "MODERATOR"
)

func (r Role) String() string {
    return string(r)
}

func (r Role) IsValid() bool {
    switch r {
    case RoleUser, RoleAdmin, RoleModerator:
        return true
    default:
        return false
    }
}
```

## Resolver Implementation Patterns

### 1. Standard Resolver Pattern
Simple, direct implementation:

```go
// internal/resolvers/user.resolvers.go.tmpl
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
    // Validate authentication
    user := auth.UserFromContext(ctx)
    if user == nil {
        return nil, errors.New("unauthorized")
    }
    
    // Fetch data
    {{if .Database.Enabled}}
    {{if eq .Database.ORM "gorm"}}
    var users []*domain.User
    err := r.DB.Find(&users).Error
    if err != nil {
        return nil, err
    }
    {{else if eq .Database.ORM "sqlx"}}
    var users []*domain.User
    err := r.DB.Select(&users, "SELECT * FROM users")
    if err != nil {
        return nil, err
    }
    {{end}}
    
    // Convert domain models to GraphQL models
    result := make([]*model.User, len(users))
    for i, user := range users {
        result[i] = &model.User{
            ID:        user.ID,
            Name:      user.Name,
            Email:     user.Email,
            Role:      model.Role(user.Role),
            CreatedAt: user.CreatedAt,
        }
    }
    
    return result, nil
    {{else}}
    // Mock data for non-database setup
    return []*model.User{
        {
            ID:        "1",
            Name:      "John Doe",
            Email:     "john@example.com",
            Role:      model.RoleUser,
            CreatedAt: time.Now(),
        },
    }, nil
    {{end}}
}
```

### 2. Clean Architecture Resolver Pattern
Use cases and repositories:

```go
// internal/resolvers/user.resolvers.go.tmpl
type UserResolver struct {
    userUseCase *usecase.UserUseCase
    logger      {{.LoggerType}}.Logger
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
    // Validate authentication
    user := auth.UserFromContext(ctx)
    if user == nil {
        return nil, errors.New("unauthorized")
    }
    
    // Use case execution
    users, err := r.userUseCase.GetUsers(ctx)
    if err != nil {
        r.logger.Error("failed to get users", "error", err)
        return nil, errors.New("internal server error")
    }
    
    // Convert domain models to GraphQL models
    result := make([]*model.User, len(users))
    for i, user := range users {
        result[i] = r.userToGraphQL(user)
    }
    
    return result, nil
}

func (r *UserResolver) userToGraphQL(user *domain.User) *model.User {
    return &model.User{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        Role:      model.Role(user.Role),
        CreatedAt: user.CreatedAt,
    }
}
```

### 3. Field Resolvers for Complex Relationships
Handle N+1 queries with field resolvers:

```go
// internal/resolvers/user.resolvers.go.tmpl
func (r *userResolver) Orders(ctx context.Context, user *model.User) ([]*model.Order, error) {
    // Use DataLoader to prevent N+1 queries
    loader := dataloader.GetOrderLoader(ctx)
    orders, err := loader.Load(user.ID)
    if err != nil {
        return nil, err
    }
    
    result := make([]*model.Order, len(orders))
    for i, order := range orders {
        result[i] = r.orderToGraphQL(order)
    }
    
    return result, nil
}
```

## Subscription Implementation

### WebSocket Server Setup
```go
// internal/server/subscriptions.go.tmpl
{{if .GraphQL.Subscriptions}}
func (s *Server) setupSubscriptions() {
    s.graphqlServer.AddTransport(&transport.Websocket{
        Upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool {
                // Configure CORS for WebSocket connections
                origin := r.Header.Get("Origin")
                return s.isAllowedOrigin(origin)
            },
        },
        InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
            // Authenticate WebSocket connection
            token := initPayload["Authorization"]
            if token == nil {
                return nil, errors.New("unauthorized")
            }
            
            user, err := s.auth.ValidateToken(token.(string))
            if err != nil {
                return nil, err
            }
            
            return auth.WithUser(ctx, user), nil
        },
    })
}
{{end}}
```

### Subscription Resolvers
```go
// internal/resolvers/subscription.resolvers.go.tmpl
{{if .GraphQL.Subscriptions}}
func (r *subscriptionResolver) UserUpdated(ctx context.Context, id string) (<-chan *model.User, error) {
    // Validate authorization
    user := auth.UserFromContext(ctx)
    if user == nil {
        return nil, errors.New("unauthorized")
    }
    
    // Create channel for subscription
    ch := make(chan *model.User, 1)
    
    // Subscribe to user updates
    r.pubsub.Subscribe(ctx, "user_updated_"+id, func(data interface{}) {
        if user, ok := data.(*domain.User); ok {
            select {
            case ch <- r.userToGraphQL(user):
            case <-ctx.Done():
                return
            }
        }
    })
    
    return ch, nil
}
{{end}}
```

## Testing Schema-First Development

### Schema Validation Tests
```go
// internal/schema/validation_test.go.tmpl
func TestSchemaValidation(t *testing.T) {
    schema, err := ioutil.ReadFile("../../schema/schema.graphql")
    require.NoError(t, err)
    
    _, err = gqlparser.LoadSchema(&ast.Source{Input: string(schema)})
    require.NoError(t, err, "Schema should be valid GraphQL")
}

func TestSchemaComplexity(t *testing.T) {
    schema := generated.NewExecutableSchema(generated.Config{})
    
    query := `
        query {
            users {
                id
                name
                orders {
                    id
                    products {
                        id
                        name
                    }
                }
            }
        }
    `
    
    complexity := analysis.ComplexityCalculator(query, schema.Schema())
    assert.LessOrEqual(t, complexity, 1000, "Query complexity should be within limits")
}
```

### Resolver Tests
```go
// internal/resolvers/user_test.go.tmpl
func TestUserResolver(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    resolver := &Resolver{DB: db}
    
    // Test data
    testUser := &domain.User{
        ID:    "1",
        Name:  "John Doe",
        Email: "john@example.com",
        Role:  domain.RoleUser,
    }
    db.Create(testUser)
    
    // Test query
    ctx := auth.WithUser(context.Background(), testUser)
    users, err := resolver.Query().Users(ctx)
    
    require.NoError(t, err)
    assert.Len(t, users, 1)
    assert.Equal(t, "John Doe", users[0].Name)
}
```

## Benefits of Schema-First Approach

### 1. Type Safety
- Generated types ensure compile-time safety
- Automatic validation of resolver signatures
- Prevents runtime type errors

### 2. Documentation
- Schema serves as API contract
- Self-documenting through GraphQL introspection
- Tools can generate documentation automatically

### 3. Code Generation
- Reduces boilerplate code
- Maintains consistency across resolvers
- Automatic updates when schema changes

### 4. Team Collaboration
- Frontend teams can work with schema before backend implementation
- Clear contract between teams
- Parallel development possible

### 5. Tooling Support
- Better IDE support with generated types
- Automatic validation and linting
- Integration with GraphQL ecosystem tools

This schema-first approach provides a solid foundation for building maintainable, type-safe GraphQL APIs with excellent developer experience.