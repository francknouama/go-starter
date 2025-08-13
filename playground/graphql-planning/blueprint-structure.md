# GraphQL API Blueprint Structure Planning

## Overview
This document outlines the planned structure for the GraphQL API blueprint, focusing on schema-first development and multi-framework support.

## Blueprint Directory Structure

```
blueprints/graphql-api/
├── template.yaml                 # Blueprint metadata and configuration
├── schema/
│   ├── schema.graphql.tmpl       # Main GraphQL schema definition
│   └── fragments/
│       ├── types.graphql.tmpl    # Common types and interfaces
│       └── scalars.graphql.tmpl  # Custom scalar definitions
├── resolvers/
│   ├── resolver.go.tmpl          # Main resolver struct
│   ├── query.resolvers.go.tmpl   # Query resolvers
│   ├── mutation.resolvers.go.tmpl # Mutation resolvers
│   └── subscription.resolvers.go.tmpl # Subscription resolvers (optional)
├── models/
│   ├── models_gen.go.tmpl        # Generated models (gqlgen)
│   └── custom_models.go.tmpl     # Custom model implementations
├── directives/
│   └── auth.go.tmpl              # Custom directives (auth, validation)
├── middleware/
│   ├── cors.go.tmpl              # CORS middleware
│   ├── auth.go.tmpl              # Authentication middleware
│   └── logging.go.tmpl           # Logging middleware
├── server/
│   ├── server.go.tmpl            # HTTP server setup
│   └── playground.go.tmpl        # GraphQL playground (conditional)
├── config/
│   ├── config.go.tmpl            # Configuration management
│   └── database.go.tmpl          # Database connection (conditional)
├── internal/
│   ├── services/
│   │   └── user.go.tmpl          # Business logic services
│   └── repository/
│       └── user.go.tmpl          # Data access layer (conditional)
├── cmd/
│   └── server/
│       └── main.go.tmpl          # Application entry point
├── deployments/
│   ├── docker/
│   │   └── Dockerfile.tmpl       # Docker containerization
│   └── k8s/
│       ├── deployment.yaml.tmpl  # Kubernetes deployment
│       └── service.yaml.tmpl     # Kubernetes service
├── scripts/
│   ├── generate.sh.tmpl          # Code generation script
│   └── dev.sh.tmpl               # Development script
├── docs/
│   ├── README.md.tmpl            # Project documentation
│   └── API.md.tmpl               # API documentation
├── tests/
│   ├── integration/
│   │   └── graphql_test.go.tmpl  # Integration tests
│   └── unit/
│       └── resolver_test.go.tmpl # Unit tests for resolvers
├── go.mod.tmpl                   # Go module definition
├── go.sum.tmpl                   # Go module checksums (empty initially)
├── gqlgen.yml.tmpl               # gqlgen configuration
├── tools.go.tmpl                 # Go tools dependencies
└── .gitignore.tmpl               # Git ignore rules
```

## File Count Estimation
- **Simple**: ~15 files (basic schema, resolvers, server)
- **Standard**: ~25 files (full structure, middleware, tests)
- **Advanced**: ~35 files (all features, deployment, documentation)

## Conditional Generation Logic

### Database Integration
- Files generated only when database driver is specified
- Repository pattern implementation
- Database migrations (optional)

### Authentication
- Auth middleware and directives
- JWT token validation
- User context injection

### Subscriptions
- WebSocket server setup
- Subscription resolvers
- Real-time event handling

### Deployment
- Docker configuration
- Kubernetes manifests
- CI/CD pipeline files

## Framework Support Priority

### Primary: 99designs/gqlgen
- Most popular Go GraphQL library
- Schema-first approach
- Excellent tooling and code generation
- Strong community support

### Secondary: graphql-go/graphql (Future)
- Programmatic schema definition
- More flexible but verbose
- Good for complex custom logic

### Tertiary: vektah/gqlparser (Future)
- Lower-level parser
- For advanced custom implementations