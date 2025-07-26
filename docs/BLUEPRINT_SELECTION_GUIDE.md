# Blueprint Selection Guide

Comprehensive guide to choosing the right blueprint for your Go project. This guide helps you make informed decisions based on project requirements, team expertise, and architectural needs.

## 🎯 Quick Decision Tree

```
Start Here: What are you building?
│
├── Command-line tool?
│   ├── Simple utility? → CLI Simple
│   └── Production tool? → CLI Standard
│
├── Web API/Backend?
│   ├── Simple REST API? → Web API Standard
│   ├── Enterprise system? → Web API Clean
│   ├── Complex domain? → Web API DDD
│   └── Maximum testability? → Web API Hexagonal
│
├── Serverless function?
│   ├── Event processing? → Lambda Standard
│   └── API Gateway? → Lambda Proxy
│
├── Shared code?
│   └── Reusable package? → Library
│
├── Distributed system?
│   └── Service architecture? → Microservice
│
├── Traditional web app?
│   └── All-in-one deployment? → Monolith
│
└── Multiple related services?
    └── Monorepo structure? → Workspace
```

## 📊 Blueprint Comparison Matrix

| Blueprint | Complexity | Files | Learning Curve | Team Size | Maintenance | Performance |
|-----------|------------|-------|----------------|-----------|-------------|-------------|
| **CLI Simple** | ⭐ | 8 | Easy | 1 | Low | Good |
| **CLI Standard** | ⭐⭐ | 29 | Moderate | 1-3 | Medium | Good |
| **Web API Standard** | ⭐⭐ | 35 | Moderate | 2-5 | Medium | Excellent |
| **Web API Clean** | ⭐⭐⭐ | 45 | Hard | 3-8 | High | Excellent |
| **Web API DDD** | ⭐⭐⭐⭐ | 50 | Very Hard | 4-10 | High | Good |
| **Web API Hexagonal** | ⭐⭐⭐⭐⭐ | 55 | Expert | 5-12 | Very High | Good |
| **Lambda Standard** | ⭐ | 12 | Easy | 1-2 | Low | Excellent |
| **Lambda Proxy** | ⭐⭐ | 25 | Moderate | 2-4 | Medium | Excellent |
| **Library** | ⭐ | 15 | Easy | 1-3 | Low | N/A |
| **Microservice** | ⭐⭐⭐⭐ | 60 | Very Hard | 4-8 | High | Excellent |
| **Monolith** | ⭐⭐⭐ | 65 | Hard | 3-8 | High | Good |
| **Workspace** | ⭐⭐⭐ | 40 | Hard | 3-10 | High | Good |

## 🏗️ Architecture Pattern Comparison

### Standard Architecture
> **Best for**: Most projects, rapid development, small to medium teams

**Structure**: Traditional layered architecture
```
handlers/ → services/ → repository/ → database
```

**Pros**:
- ✅ Simple and familiar
- ✅ Fast development
- ✅ Easy to understand
- ✅ Good for most use cases

**Cons**:
- ❌ Can become monolithic
- ❌ Less testable
- ❌ Tight coupling possible

**When to Choose**:
- Building MVP or prototype
- Small to medium team
- Straightforward business logic
- Time-to-market is priority

### Clean Architecture
> **Best for**: Enterprise applications, complex business logic, long-term projects

**Structure**: Dependency inversion with clear boundaries
```
adapters/ ← usecases/ ← domain/
```

**Pros**:
- ✅ Highly testable
- ✅ Independent of frameworks
- ✅ Clear separation of concerns
- ✅ Maintainable long-term

**Cons**:
- ❌ Higher initial complexity
- ❌ More files and interfaces
- ❌ Steeper learning curve

**When to Choose**:
- Enterprise applications
- Complex business requirements
- Long-term maintainability needed
- Team familiar with clean architecture

### Domain-Driven Design (DDD)
> **Best for**: Complex domains, rich business models, event-driven systems

**Structure**: Domain-centric with aggregates and bounded contexts
```
domain/ (aggregates, entities, services) → application/ → infrastructure/
```

**Pros**:
- ✅ Rich domain models
- ✅ Business logic centralization
- ✅ Clear business boundaries
- ✅ Event-driven capabilities

**Cons**:
- ❌ Very complex initially
- ❌ Requires domain expertise
- ❌ Can be over-engineered

**When to Choose**:
- Complex business domains
- Rich business rules
- Event-driven architecture
- Domain experts available

### Hexagonal Architecture
> **Best for**: Maximum testability, multiple interfaces, ports & adapters

**Structure**: Ports and adapters with business logic at center
```
primary adapters → application (ports) → secondary adapters
```

**Pros**:
- ✅ Maximum testability
- ✅ Multiple interface support
- ✅ Technology independence
- ✅ Easy adapter swapping

**Cons**:
- ❌ Highest complexity
- ❌ Many interfaces
- ❌ Can be over-abstracted

**When to Choose**:
- Multiple interfaces needed (HTTP, gRPC, CLI)
- Maximum testability required
- Frequent technology changes
- Expert team available

## 📋 Detailed Blueprint Analysis

### CLI Applications

#### CLI Simple
**Perfect for**:
- Learning Go fundamentals
- Quick automation scripts
- Personal utilities
- Prototyping CLI tools

**Generated Structure**:
```
my-tool/
├── main.go          # Entry point
├── cmd/             # Commands
├── config.go        # Simple config
├── Makefile         # Build automation
└── README.md        # Documentation
```

**Command**: `go-starter new my-tool --type=cli --complexity=simple`

#### CLI Standard
**Perfect for**:
- Production CLI tools
- Developer tools
- System administration utilities
- CI/CD tools

**Advanced Features**:
- Multiple subcommands
- Configuration files
- Interactive prompts
- Shell completion
- Structured logging
- Comprehensive testing

**Command**: `go-starter new my-tool --type=cli --complexity=standard`

### Web APIs

#### Web API Standard
**Perfect for**:
- REST APIs
- Microservices
- Backend services
- API-first applications

**Key Features**:
- HTTP framework integration (Gin, Echo, Fiber)
- Database support (PostgreSQL, MySQL, MongoDB)
- Authentication (JWT, OAuth2)
- Middleware stack
- OpenAPI documentation

**Command**: 
```bash
go-starter new my-api --type=web-api \
  --framework=gin \
  --database-driver=postgres \
  --auth-type=jwt
```

#### Web API Clean
**Perfect for**:
- Enterprise applications
- Complex business logic
- Long-term projects
- High testability requirements

**Architecture Benefits**:
- Clear separation of concerns
- Framework independence
- Testable business logic
- Dependency inversion

**Command**:
```bash
go-starter new my-api --type=web-api --architecture=clean \
  --database-driver=postgres \
  --auth-type=jwt
```

#### Web API DDD
**Perfect for**:
- Complex business domains
- Rich domain models
- Event-driven systems
- Domain expert collaboration

**Domain Features**:
- Aggregates and entities
- Domain services
- Specifications (business rules)
- Domain events
- Bounded contexts

**Command**:
```bash
go-starter new my-api --type=web-api --architecture=ddd \
  --database-driver=postgres \
  --auth-type=jwt
```

#### Web API Hexagonal
**Perfect for**:
- Maximum testability
- Multiple adapters
- Technology independence
- Complex integration requirements

**Hexagonal Benefits**:
- Ports and adapters pattern
- Multiple primary adapters (HTTP, gRPC, CLI)
- Multiple secondary adapters (DB, cache, email)
- Complete isolation of business logic

**Command**:
```bash
go-starter new my-api --type=web-api --architecture=hexagonal \
  --database-driver=postgres \
  --auth-type=jwt
```

### Serverless

#### Lambda Standard
**Perfect for**:
- Event processing
- Background tasks
- Webhooks
- Simple serverless functions

**AWS Integration**:
- CloudWatch Events
- X-Ray tracing
- CloudWatch metrics
- SAM deployment

**Command**: `go-starter new my-lambda --type=lambda`

#### Lambda Proxy
**Perfect for**:
- REST APIs on Lambda
- API Gateway integration
- Serverless web backends
- Cost-optimized APIs

**API Features**:
- API Gateway integration
- HTTP routing
- Request/response handling
- Serverless deployment

**Command**: `go-starter new my-api --type=lambda-proxy`

### Specialized Blueprints

#### Library
**Perfect for**:
- Reusable packages
- SDK development
- Open-source projects
- Shared utilities

**Library Features**:
- Public API design
- Comprehensive examples
- Documentation generation
- Versioning support

**Command**: `go-starter new my-lib --type=library`

#### Microservice
**Perfect for**:
- Distributed systems
- Service mesh architectures
- Cloud-native applications
- Scalable backends

**Microservice Features**:
- gRPC server/client
- Health checks
- Metrics collection
- Distributed tracing
- Kubernetes deployment

**Command**: `go-starter new my-service --type=microservice`

#### Monolith
**Perfect for**:
- Traditional web applications
- Rapid prototyping
- Small teams
- Simplified deployment

**Monolith Features**:
- MVC architecture
- HTML templates
- Static assets
- Database migrations
- Admin interface

**Command**: `go-starter new my-app --type=monolith`

#### Workspace
**Perfect for**:
- Monorepos
- Multiple related services
- Shared libraries
- Complex project organization

**Workspace Features**:
- Go workspace configuration
- Multiple modules
- Shared dependencies
- Unified build system

**Command**: `go-starter new my-workspace --type=workspace`

## 🎭 Use Case Scenarios

### Startup MVP
**Scenario**: Building a minimum viable product for a startup
**Recommendation**: Web API Standard
**Reasoning**: Fast development, good performance, scalable foundation

```bash
go-starter new startup-api --type=web-api \
  --framework=gin \
  --database-driver=postgres \
  --auth-type=jwt \
  --logger=slog
```

### Enterprise System
**Scenario**: Large enterprise with complex business rules
**Recommendation**: Web API Clean or DDD
**Reasoning**: Maintainable, testable, handles complex business logic

```bash
go-starter new enterprise-system --type=web-api --architecture=clean \
  --database-driver=postgres \
  --database-orm=gorm \
  --auth-type=jwt \
  --logger=zap \
  --advanced
```

### Developer Tool
**Scenario**: Building a CLI tool for developers
**Recommendation**: CLI Standard
**Reasoning**: Rich feature set, professional polish, good UX

```bash
go-starter new dev-tool --type=cli --complexity=standard \
  --logger=slog
```

### Event Processing
**Scenario**: Processing events from various AWS services
**Recommendation**: Lambda Standard
**Reasoning**: Serverless, event-driven, cost-effective

```bash
go-starter new event-processor --type=lambda \
  --logger=zerolog
```

### Microservices Platform
**Scenario**: Building a platform with multiple microservices
**Recommendation**: Workspace + Microservice
**Reasoning**: Shared libraries, consistent patterns, service organization

```bash
# Create workspace
go-starter new platform --type=workspace

# Add individual services
cd platform/services
go-starter new user-service --type=microservice
go-starter new order-service --type=microservice
go-starter new notification-service --type=microservice
```

### API with Multiple Interfaces
**Scenario**: Need HTTP REST, gRPC, and CLI interfaces
**Recommendation**: Web API Hexagonal
**Reasoning**: Multiple adapters, clean separation, testable

```bash
go-starter new multi-interface-api --type=web-api --architecture=hexagonal \
  --database-driver=postgres \
  --auth-type=jwt \
  --logger=zap \
  --advanced
```

## 🚀 Migration Paths

### Scaling Up Complexity

#### Simple → Standard CLI
When your CLI tool needs more features:
1. Generate new standard CLI project
2. Copy business logic from simple version
3. Add new features (subcommands, config files)
4. Migrate gradually

#### Standard → Clean Architecture
When your API needs better structure:
1. Generate clean architecture version
2. Identify domain entities and use cases
3. Move business logic to domain layer
4. Implement repository interfaces
5. Migrate endpoint by endpoint

#### Clean → DDD
When you need richer domain modeling:
1. Identify aggregates and bounded contexts
2. Create domain services and specifications
3. Add domain events
4. Implement event handlers

### Scaling Down Complexity

Sometimes you might need to simplify:

#### DDD → Clean
If DDD proves too complex:
1. Flatten aggregates to simple entities
2. Move domain services to use cases
3. Remove complex specifications
4. Simplify event handling

#### Clean → Standard
If clean architecture is overkill:
1. Merge use cases into service layer
2. Remove adapter interfaces
3. Simplify to traditional layers
4. Keep good practices (testing, structure)

## 📊 Performance Characteristics

### Runtime Performance

| Blueprint | Startup Time | Memory Usage | Throughput | Latency |
|-----------|--------------|--------------|------------|---------|
| **CLI Simple** | Fast | Low | N/A | N/A |
| **CLI Standard** | Moderate | Medium | N/A | N/A |
| **Web API Standard** | Fast | Medium | High | Low |
| **Web API Clean** | Moderate | Medium | High | Low |
| **Web API DDD** | Slow | High | Medium | Medium |
| **Web API Hexagonal** | Slow | High | Medium | Medium |
| **Lambda Standard** | Very Fast | Very Low | High | Very Low |
| **Lambda Proxy** | Fast | Low | High | Low |
| **Microservice** | Moderate | Medium | High | Medium |

### Development Performance

| Blueprint | Initial Setup | Feature Addition | Bug Fixing | Testing |
|-----------|---------------|------------------|------------|---------|
| **CLI Simple** | Very Fast | Fast | Fast | Easy |
| **CLI Standard** | Fast | Moderate | Moderate | Moderate |
| **Web API Standard** | Fast | Fast | Moderate | Moderate |
| **Web API Clean** | Moderate | Moderate | Easy | Easy |
| **Web API DDD** | Slow | Slow | Easy | Easy |
| **Web API Hexagonal** | Very Slow | Moderate | Very Easy | Very Easy |
| **Lambda Standard** | Very Fast | Fast | Fast | Easy |
| **Lambda Proxy** | Fast | Moderate | Moderate | Moderate |

## 🎯 Team Expertise Requirements

### Beginner-Friendly
- **CLI Simple**: Perfect for Go beginners
- **Lambda Standard**: Good serverless introduction
- **Library**: Focuses on Go fundamentals

### Intermediate
- **CLI Standard**: Good for CLI development learning
- **Web API Standard**: Solid web development foundation
- **Lambda Proxy**: Serverless with more complexity

### Advanced
- **Web API Clean**: Requires architecture knowledge
- **Microservice**: Needs distributed systems understanding
- **Monolith**: Complex but familiar patterns

### Expert
- **Web API DDD**: Requires domain modeling expertise
- **Web API Hexagonal**: Advanced architecture patterns
- **Workspace**: Complex project organization

## 🛠️ Customization Options

### Framework Selection

#### Web Frameworks
```bash
# Gin (fastest, most popular)
--framework=gin

# Echo (middleware-rich)
--framework=echo

# Fiber (Express-like)
--framework=fiber

# Chi (lightweight)
--framework=chi
```

#### CLI Frameworks
```bash
# Cobra (most popular)
--framework=cobra
```

### Database Options

```bash
# Relational databases
--database-driver=postgres  # Recommended for production
--database-driver=mysql     # Alternative relational DB
--database-driver=sqlite    # Development/testing

# NoSQL databases
--database-driver=mongodb   # Document database

# ORM/Query Builders
--database-orm=gorm        # Feature-rich ORM
--database-orm=sqlx        # Lightweight extensions
--database-orm=sqlc        # Type-safe SQL generation
--database-orm=ent         # Facebook's entity framework
```

### Authentication

```bash
# JSON Web Tokens (most popular)
--auth-type=jwt

# OAuth2 providers
--auth-type=oauth2

# Session-based (traditional)
--auth-type=session

# API key authentication
--auth-type=api-key
```

### Logger Selection

```bash
# Standard library (Go 1.21+)
--logger=slog

# High performance
--logger=zap

# Zero allocation
--logger=zerolog

# Feature-rich
--logger=logrus
```

## 📋 Quick Reference Commands

### Most Common Combinations

```bash
# Simple learning project
go-starter new my-project --type=cli --complexity=simple

# Production CLI tool
go-starter new my-tool --type=cli --complexity=standard --logger=slog

# Standard REST API
go-starter new my-api --type=web-api --framework=gin --database-driver=postgres --auth-type=jwt

# Enterprise API
go-starter new enterprise-api --type=web-api --architecture=clean --database-driver=postgres --auth-type=jwt --logger=zap

# Serverless function
go-starter new my-lambda --type=lambda --logger=zerolog

# Shared library
go-starter new my-lib --type=library

# Microservice
go-starter new my-service --type=microservice --logger=zap

# Complex domain API
go-starter new domain-api --type=web-api --architecture=ddd --database-driver=postgres --auth-type=jwt --advanced
```

### Preview Before Generation

```bash
# See what will be generated
go-starter new my-project --type=web-api --dry-run

# Advanced preview
go-starter new my-project --type=web-api --architecture=clean --dry-run --advanced
```

## 🎯 Final Recommendations

### Choose Based On:

1. **Team Size & Expertise**
   - 1 person: CLI Simple, Lambda Standard
   - 2-5 people: Web API Standard, CLI Standard
   - 5+ people: Web API Clean, Microservice

2. **Project Longevity**
   - Short-term: Standard architectures
   - Long-term: Clean/DDD architectures

3. **Business Complexity**
   - Simple: Standard patterns
   - Complex: DDD or Clean Architecture

4. **Performance Requirements**
   - High throughput: Web API Standard, Lambda
   - Low latency: Web API Standard with optimized logger

5. **Testability Needs**
   - Standard: Web API Standard
   - High: Web API Clean
   - Maximum: Web API Hexagonal

Remember: You can always start simple and migrate to more complex architectures as your project grows. go-starter's consistent patterns make this migration easier.