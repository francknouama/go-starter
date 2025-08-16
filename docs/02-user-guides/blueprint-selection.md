# Blueprint Selection Guide

Choose the perfect blueprint for your Go project. This guide provides decision frameworks, comparisons, and real-world scenarios to help you make the right choice.

## 🎯 Quick Decision Tree

```
What are you building?
│
├── 🖥️  Command-line tool?
│   ├── Simple utility/script? → CLI Simple (8 files)
│   └── Production CLI tool? → CLI Standard (29 files)
│
├── 🌐 Web API/Backend service?
│   ├── Simple REST API? → Web API Standard (35 files)
│   ├── Enterprise application? → Web API Clean (45 files)
│   ├── Complex business domain? → Web API DDD (50 files)
│   └── Maximum testability needed? → Web API Hexagonal (55 files)
│
├── ⚡ Serverless function?
│   ├── Event processing/background jobs? → Lambda Standard (12 files)
│   └── API Gateway integration? → Lambda Proxy (25 files)
│
├── 📦 Shared/reusable package?
│   └── Library/SDK for others to use? → Library (15 files)
│
├── 🏗️ Distributed system?
│   ├── Single microservice? → Microservice (60 files)
│   └── Multiple related services? → Workspace (varies)
│
└── 🏢 Traditional web application?
    └── All-in-one deployment? → Monolith (70+ files)
```

## 📊 Blueprint Comparison Matrix

| Blueprint | Complexity | Files | Learning Curve | Team Size | Use Cases |
|-----------|------------|-------|----------------|-----------|-----------|
| **CLI Simple** | ⭐ | 8 | Easy | 1 | Scripts, utilities, prototypes |
| **CLI Standard** | ⭐⭐ | 29 | Moderate | 1-3 | Production tools, automation |
| **Web API Standard** | ⭐⭐ | 35 | Moderate | 2-5 | REST APIs, microservices |
| **Web API Clean** | ⭐⭐⭐ | 45 | Hard | 3-8 | Enterprise APIs, complex logic |
| **Web API DDD** | ⭐⭐⭐⭐ | 50 | Very Hard | 4-10 | Domain-rich applications |
| **Web API Hexagonal** | ⭐⭐⭐⭐⭐ | 55 | Expert | 5-12 | Maximum testability |
| **Lambda Standard** | ⭐ | 12 | Easy | 1-2 | Event processing, webhooks |
| **Lambda Proxy** | ⭐⭐ | 25 | Moderate | 2-4 | Serverless APIs |
| **Library** | ⭐ | 15 | Easy | 1-3 | SDKs, shared utilities |
| **Microservice** | ⭐⭐⭐⭐ | 60 | Very Hard | 4-8 | Service mesh, gRPC |
| **Monolith** | ⭐⭐⭐ | 70+ | Hard | 3-10 | Traditional web apps |
| **Workspace** | ⭐⭐⭐⭐ | Varies | Very Hard | 5-15 | Monorepos, multi-service |

## 🎮 Interactive Selection Guide

### Step 1: What's Your Primary Goal?

#### 🚀 **Getting Started / Learning**
**Recommended**: CLI Simple or Web API Standard
- Start with simple, understandable structures
- Focus on Go fundamentals, not architecture complexity
- Easy to modify and experiment with

#### 🏢 **Production Application**
**Consider**: Team size, domain complexity, scalability needs
- Small team (1-3): Web API Standard, CLI Standard
- Medium team (3-8): Web API Clean, Microservice
- Large team (8+): Web API DDD, Workspace

#### 🔬 **Experimentation / Prototyping**
**Recommended**: CLI Simple, Lambda Standard
- Minimal setup overhead
- Quick iteration and testing
- Easy to throw away and restart

### Step 2: Architecture Complexity Assessment

#### ⭐ **Simple (Beginner-Friendly)**
```bash
# Perfect for learning and simple projects
go-starter new my-project --type=cli --complexity=simple
go-starter new my-api --type=web-api --framework=gin
go-starter new my-function --type=lambda
```

#### ⭐⭐⭐ **Intermediate (Production-Ready)**
```bash
# Balanced complexity and maintainability
go-starter new my-tool --type=cli --complexity=standard
go-starter new my-service --type=web-api --architecture=clean
go-starter new my-microservice --type=microservice
```

#### ⭐⭐⭐⭐⭐ **Advanced (Enterprise-Grade)**
```bash
# Maximum patterns and testability
go-starter new enterprise-api --type=web-api --architecture=hexagonal
go-starter new domain-service --type=web-api --architecture=ddd
go-starter new distributed-system --type=workspace
```

## 🏗️ Architecture Pattern Guide

### Standard Architecture
**Best for**: Simple APIs, quick prototypes, learning
```
├── main.go
├── internal/
│   ├── handlers/
│   ├── services/
│   └── models/
└── pkg/
```
**Pros**: Simple, fast development, easy to understand  
**Cons**: Can become messy with complex logic

### Clean Architecture
**Best for**: Enterprise applications, team development
```
├── cmd/
├── internal/
│   ├── delivery/     # Controllers/Handlers
│   ├── usecase/      # Business logic
│   ├── repository/   # Data access
│   └── domain/       # Business entities
└── pkg/
```
**Pros**: Clear separation, testable, scalable  
**Cons**: More files, steeper learning curve

### Domain-Driven Design (DDD)
**Best for**: Complex business domains, large teams
```
├── internal/
│   ├── domain/       # Core business logic
│   ├── application/  # Use cases
│   ├── infrastructure/ # External concerns
│   └── interfaces/   # Adapters
└── pkg/
```
**Pros**: Domain-focused, handles complexity well  
**Cons**: High learning curve, can be over-engineered

### Hexagonal Architecture
**Best for**: Maximum testability, complex integrations
```
├── internal/
│   ├── core/         # Business logic (hexagon center)
│   ├── ports/        # Interfaces
│   └── adapters/     # Implementation details
└── pkg/
```
**Pros**: Extremely testable, decoupled  
**Cons**: Complexity overhead, expert-level pattern

## 🔍 Real-World Scenarios

### Scenario 1: Startup API
**Context**: Small team, MVP, rapid iteration needed
**Recommendation**: Web API Standard with Gin
```bash
go-starter new startup-api --type=web-api --framework=gin --logger=slog
```
**Why**: Fast development, proven patterns, easy to scale later

### Scenario 2: Enterprise System
**Context**: Large team, complex domain, long-term maintenance
**Recommendation**: Web API Clean or DDD
```bash
go-starter new enterprise-system --type=web-api --architecture=clean --database-driver=postgres
```
**Why**: Maintainable architecture, team scalability

### Scenario 3: DevOps Tool
**Context**: Internal automation, CLI interface needed
**Recommendation**: CLI Standard
```bash
go-starter new deploy-tool --type=cli --complexity=standard --logger=zap
```
**Why**: Rich CLI features, good structure, production-ready

### Scenario 4: Event Processing
**Context**: Serverless, event-driven, cost-sensitive
**Recommendation**: Lambda Standard
```bash
go-starter new event-processor --type=lambda --logger=slog
```
**Why**: Minimal cold start, cost-effective, event-focused

### Scenario 5: SDK/Library
**Context**: Public API client, reusable code
**Recommendation**: Library
```bash
go-starter new api-client --type=library --logger=slog
```
**Why**: Clean public API, good documentation, testing focus

## 🛠️ Framework Selection Guide

### Web Frameworks

#### Gin (Most Popular)
**Best for**: General-purpose APIs, learning, rapid development
```bash
go-starter new api --type=web-api --framework=gin
```
**Pros**: Mature, fast, large community  
**Cons**: Less opinionated

#### Echo (Performance-Focused)
**Best for**: High-performance APIs, microservices
```bash
go-starter new api --type=web-api --framework=echo
```
**Pros**: High performance, good middleware  
**Cons**: Smaller community

#### Fiber (Express.js-like)
**Best for**: Node.js developers, ultra-fast APIs
```bash
go-starter new api --type=web-api --framework=fiber
```
**Pros**: Familiar to JS developers, very fast  
**Cons**: Different from standard Go patterns

#### Chi (Lightweight)
**Best for**: Minimal APIs, standard library feel
```bash
go-starter new api --type=web-api --framework=chi
```
**Pros**: Lightweight, idiomatic Go  
**Cons**: Less built-in features

### Logger Selection

#### slog (Recommended)
**Standard library, structured logging**
```bash
--logger=slog
```

#### Zap (High Performance)
**Zero-allocation, production-grade**
```bash
--logger=zap
```

#### Logrus (Popular)
**Mature, feature-rich**
```bash
--logger=logrus
```

#### Zerolog (Fast)
**Zero-allocation, chainable**
```bash
--logger=zerolog
```

## 🚦 Progressive Complexity Strategy

### Phase 1: Start Simple
Begin with simpler blueprints to learn patterns:
```bash
# Week 1-2: Learn basics
go-starter new learning-cli --type=cli --complexity=simple
go-starter new learning-api --type=web-api --framework=gin
```

### Phase 2: Add Structure
Move to production-ready patterns:
```bash
# Week 3-4: Production patterns
go-starter new production-cli --type=cli --complexity=standard
go-starter new production-api --type=web-api --architecture=clean
```

### Phase 3: Master Complexity
Tackle advanced patterns:
```bash
# Month 2+: Advanced patterns
go-starter new enterprise-api --type=web-api --architecture=ddd
go-starter new distributed-system --type=workspace
```

## 🎯 Decision Checklist

Before choosing a blueprint, answer these questions:

### Technical Requirements
- [ ] What's the primary interface? (CLI, REST API, events, library)
- [ ] What's the expected scale? (requests/day, concurrent users)
- [ ] What integrations are needed? (databases, external APIs, message queues)
- [ ] What's the performance requirement? (latency, throughput)

### Team Considerations
- [ ] How large is the development team?
- [ ] What's the team's Go experience level?
- [ ] How long will this project be maintained?
- [ ] Are there existing architecture standards?

### Project Context
- [ ] Is this a prototype or production system?
- [ ] How quickly do you need to deliver?
- [ ] What's the deployment environment? (cloud, on-premise, serverless)
- [ ] Are there compliance or security requirements?

## 💡 Best Practices

### When to Choose Simple Blueprints
- ✅ Learning Go or new patterns
- ✅ Prototyping and experimentation
- ✅ Small, focused utilities
- ✅ Time-sensitive projects

### When to Choose Complex Blueprints
- ✅ Long-term production systems
- ✅ Large development teams
- ✅ Complex business domains
- ✅ High testability requirements

### Common Mistakes to Avoid
- ❌ Over-engineering simple projects
- ❌ Under-architecting complex systems
- ❌ Choosing patterns the team doesn't understand
- ❌ Ignoring future scalability needs

## 🆘 Still Not Sure?

### Quick Recommendations by Use Case

**Building a REST API?** → Start with Web API Standard + Gin  
**Creating a CLI tool?** → Start with CLI Simple, upgrade to Standard if needed  
**Processing events?** → Start with Lambda Standard  
**Building a library?** → Use Library blueprint  
**Complex enterprise system?** → Web API Clean or DDD  
**Maximum testability?** → Web API Hexagonal  

### Get Help
- **Community**: [GitHub Discussions](https://github.com/francknouama/go-starter/discussions)
- **Examples**: [Community Showcases](../05-community/showcases.md)
- **Migration**: [Upgrading Blueprints Guide](upgrading-blueprints.md)

---

**🎯 Ready to choose?** Use this guide to make confident blueprint decisions and build better Go projects!