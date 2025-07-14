# templates Package

This package manages the blueprint template system for go-starter, including loading, parsing, and registering project blueprints.

## Overview

The templates package provides the core template engine that powers project generation, supporting embedded blueprints, template processing with Go templates, and a flexible registry system.

## Key Components

### Types

- **`Blueprint`** - Blueprint definition with metadata
- **`TemplateRegistry`** - Central registry of all blueprints
- **`TemplateEngine`** - Template processing engine
- **`BlueprintMetadata`** - Blueprint information and variables
- **`TemplateFile`** - Individual template file definition

### Core Functions

- **`NewRegistry() *TemplateRegistry`** - Create blueprint registry
- **`LoadBlueprints() error`** - Load all embedded blueprints
- **`GetBlueprint(name string) (*Blueprint, error)`** - Retrieve blueprint
- **`ListBlueprints() []BlueprintInfo`** - List available blueprints
- **`ProcessTemplate(tmpl string, vars interface{}) (string, error)`** - Render template

## Blueprint Structure

```
blueprints/
├── api/
│   ├── standard/
│   │   ├── template.yaml
│   │   ├── main.go.tmpl
│   │   ├── go.mod.tmpl
│   │   └── ...
│   ├── clean/
│   ├── ddd/
│   └── hexagonal/
├── cli/
├── library/
└── ...
```

## Template Configuration (template.yaml)

```yaml
name: "Web API - Standard"
description: "RESTful API with standard project structure"
type: "api"
architecture: "standard"
minGoVersion: "1.21"

variables:
  - name: "ProjectName"
    description: "Name of the project"
    type: "string"
    required: true
  - name: "Framework"
    description: "Web framework to use"
    type: "select"
    options: ["gin", "echo", "fiber", "chi"]
    default: "gin"

files:
  - source: "main.go.tmpl"
    destination: "main.go"
  - source: "internal/server.go.tmpl"
    destination: "internal/server/server.go"
    condition: "{{ne .Framework \"\"}}"
  - source: "config/{{.Environment}}.yaml.tmpl"
    destination: "config/{{.Environment}}.yaml"

dependencies:
  - module: "github.com/gin-gonic/gin"
    version: "v1.9.1"
    condition: "{{eq .Framework \"gin\"}}"

hooks:
  post_generate:
    - command: "go mod download"
      description: "Download dependencies"
    - command: "go fmt ./..."
      description: "Format code"
```

## Template Engine Features

### Built-in Functions
All Go template functions plus:
- Sprig template functions
- Custom helper functions
- String manipulation
- Path handling
- Conditional logic

### Custom Functions
```go
// Available in templates
{{.ProjectName | toUpper}}
{{.ModulePath | base}}
{{hasFeature "database"}}
{{sanitizePath .FilePath}}
```

### Conditional Generation
```go
{{if eq .Framework "gin"}}
    router := gin.Default()
{{else if eq .Framework "echo"}}
    e := echo.New()
{{end}}
```

## Registry Management

### Blueprint Loading
```go
// Embedded blueprints
//go:embed blueprints/*
var blueprintFS embed.FS

// Load at startup
registry := NewRegistry()
err := registry.LoadFromFS(blueprintFS)
```

### Blueprint Validation
- Schema validation for template.yaml
- Template syntax checking
- Variable reference validation
- Dependency verification

## Usage Example

```go
import "github.com/yourusername/go-starter/internal/templates"

// Get registry
registry := templates.GetRegistry()

// List blueprints
blueprints := registry.ListBlueprints()
for _, bp := range blueprints {
    fmt.Printf("%s: %s\n", bp.Name, bp.Description)
}

// Get specific blueprint
blueprint, err := registry.GetBlueprint("api-standard")
if err != nil {
    log.Fatal(err)
}

// Process template
vars := map[string]interface{}{
    "ProjectName": "myapp",
    "Framework": "gin",
}
content, err := templates.ProcessTemplate(templateContent, vars)
```

## Template Best Practices

1. **Use Clear Variable Names**
   ```go
   {{.ProjectName}} not {{.Name}}
   {{.DatabaseDriver}} not {{.DB}}
   ```

2. **Provide Defaults**
   ```yaml
   variables:
     - name: "Port"
       default: "8080"
   ```

3. **Document Templates**
   ```go
   {{/* This section configures the database connection */}}
   ```

4. **Handle All Cases**
   ```go
   {{if .Database}}
       // Database config
   {{else}}
       // No database
   {{end}}
   ```

## Performance Optimization

- Template caching
- Parallel processing
- Lazy loading
- Memory efficiency

## Error Handling

- Template parse errors
- Missing variable errors
- Invalid syntax detection
- Circular dependency prevention

## Dependencies

- **text/template** - Go template engine
- **github.com/Masterminds/sprig/v3** - Template functions
- **gopkg.in/yaml.v3** - YAML parsing
- **embed** - Embedded file system