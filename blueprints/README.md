# blueprints Package

This directory contains all project blueprints (templates) for go-starter, organized by project type and architecture pattern.

## Overview

Blueprints are the heart of go-starter, providing ready-to-use project templates for various Go application types. Each blueprint includes a complete project structure, configuration files, and conditional logic for customization.

## Blueprint Organization

```
blueprints/
├── api/                    # Web API blueprints
│   ├── standard/          # Standard REST API
│   ├── clean/            # Clean Architecture API
│   ├── ddd/              # Domain-Driven Design API
│   └── hexagonal/        # Hexagonal Architecture API
├── cli/                   # Command-line applications
├── library/               # Go libraries/packages
├── lambda/                # AWS Lambda functions
├── lambda-proxy/          # API Gateway Lambda proxy
├── event-driven/          # Event-driven microservices
├── microservice/          # gRPC microservices
├── monolith/             # Traditional web applications
└── workspace/            # Go workspace (monorepo)
```

## Blueprint Types

### 1. API Blueprints (`api/`)
RESTful web APIs with multiple architecture options:

| Architecture | Use Case | Complexity |
|-------------|----------|------------|
| **Standard** | Simple APIs, prototypes | Low |
| **Clean** | Enterprise applications | Medium |
| **DDD** | Complex business domains | High |
| **Hexagonal** | Testable, adaptable systems | High |

### 2. CLI Blueprint (`cli/`)
Command-line applications using Cobra framework:
- Subcommands structure
- Configuration management
- Shell completion support

### 3. Library Blueprint (`library/`)
Reusable Go packages:
- Clean API design
- Comprehensive examples
- Documentation structure

### 4. Lambda Blueprints (`lambda/`, `lambda-proxy/`)
AWS Lambda functions:
- **lambda**: Direct Lambda invocation
- **lambda-proxy**: API Gateway integration

### 5. Event-Driven Blueprint (`event-driven/`)
CQRS and Event Sourcing patterns:
- Event bus implementation
- Command/Query separation
- Event store integration

### 6. Microservice Blueprint (`microservice/`)
gRPC-based microservices:
- Protocol Buffer definitions
- Service mesh ready
- Health checks and metrics

### 7. Monolith Blueprint (`monolith/`)
Traditional web applications:
- MVC structure
- Template rendering
- Session management

### 8. Workspace Blueprint (`workspace/`)
Multi-module Go workspace:
- Shared dependencies
- Cross-module development
- Monorepo structure

## Blueprint Structure

Each blueprint contains:

### template.yaml
Blueprint metadata and configuration:
```yaml
name: "Blueprint Name"
description: "Blueprint description"
type: "api|cli|library|..."
architecture: "standard|clean|ddd|hexagonal"
minGoVersion: "1.21"

variables:
  - name: "VariableName"
    type: "string|select|boolean"
    description: "Variable description"
    default: "default value"
    required: true

files:
  - source: "template.go.tmpl"
    destination: "path/to/file.go"
    condition: "{{.SomeCondition}}"

dependencies:
  - module: "github.com/example/package"
    version: "v1.0.0"
    condition: "{{.NeedsDependency}}"

hooks:
  post_generate:
    - command: "go mod tidy"
      description: "Clean up dependencies"
```

### Template Files (*.tmpl)
Go template files with placeholders:
```go
package main

import (
    "fmt"
    {{if eq .Framework "gin"}}
    "github.com/gin-gonic/gin"
    {{end}}
)

func main() {
    fmt.Println("Welcome to {{.ProjectName}}!")
    {{if eq .Framework "gin"}}
    r := gin.Default()
    r.Run(":{{.Port}}")
    {{end}}
}
```

## Creating New Blueprints

### Step 1: Create Directory
```bash
mkdir -p blueprints/mytype/myarch
```

### Step 2: Add template.yaml
Define blueprint metadata, variables, and files.

### Step 3: Create Templates
Add `.tmpl` files with Go template syntax.

### Step 4: Test Blueprint
```bash
go-starter new --type mytype --architecture myarch
```

## Template Variables

### Common Variables
Available in all blueprints:
- `{{.ProjectName}}` - Project name
- `{{.ModulePath}}` - Go module path
- `{{.GoVersion}}` - Go version
- `{{.Author}}` - Project author
- `{{.License}}` - License type

### Type-Specific Variables
API blueprints:
- `{{.Framework}}` - Web framework (gin, echo, fiber, chi)
- `{{.Port}}` - Server port
- `{{.Authentication}}` - Auth configuration

Database features:
- `{{.Database.Driver}}` - Database type
- `{{.Database.ORM}}` - ORM choice
- `{{.Database.Migrations}}` - Migration tool

## Conditional Logic

### File Generation
```yaml
files:
  - source: "auth.go.tmpl"
    destination: "internal/auth/auth.go"
    condition: "{{.HasAuthentication}}"
```

### Dependencies
```yaml
dependencies:
  - module: "gorm.io/gorm"
    version: "v1.25.0"
    condition: "{{eq .Database.ORM \"gorm\"}}"
```

### Template Content
```go
{{if .Database}}
// Database configuration
db, err := setupDatabase()
{{end}}
```

## Best Practices

1. **Modular Design** - Keep blueprints focused and composable
2. **Sensible Defaults** - Provide good defaults for all variables
3. **Documentation** - Include README.md in generated projects
4. **Testing** - Add test files and examples
5. **Validation** - Validate inputs and provide clear errors

## Testing Blueprints

### Manual Testing
```bash
# Test blueprint generation
go-starter new --type api --architecture clean --name testproject

# Verify generated project
cd testproject
go build ./...
go test ./...
```

### Automated Testing
See `internal/templates/blueprint_test.go` for blueprint validation tests.

## Contributing Blueprints

1. Follow existing patterns
2. Test thoroughly
3. Document variables
4. Add examples
5. Submit PR with tests

## Version Compatibility

- Blueprints specify minimum Go version
- Framework versions are pinned
- Regular updates for security