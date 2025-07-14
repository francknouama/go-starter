# generator Package

This package contains the core project generation logic for go-starter, responsible for creating new Go projects from blueprints.

## Overview

The generator package orchestrates the entire project generation process, from loading blueprints to rendering templates and setting up the project structure.

## Key Components

### Types

- **`Generator`** - Main generator instance
- **`GeneratorConfig`** - Configuration for generation process
- **`GenerationResult`** - Result of project generation
- **`FileGenerator`** - Individual file generation logic

### Core Functions

- **`New(config *GeneratorConfig) *Generator`** - Create new generator
- **`Generate(blueprint *Blueprint, vars Variables) (*GenerationResult, error)`** - Generate project
- **`ValidateBlueprint(blueprint *Blueprint) error`** - Validate blueprint before generation
- **`RollbackOnError(path string) error`** - Clean up on generation failure

## Generation Process

1. **Blueprint Loading** - Load and parse selected blueprint
2. **Variable Resolution** - Merge user inputs with defaults
3. **Directory Creation** - Set up project directory structure
4. **Template Rendering** - Process templates with variables
5. **File Generation** - Write rendered files with proper permissions
6. **Post-Generation** - Run hooks and initialization commands
7. **Validation** - Ensure generated project is valid

## Features

### Conditional Generation
```go
// Files are generated based on conditions
if evaluateCondition(file.Condition, variables) {
    generateFile(file, variables)
}
```

### Template Functions
- Standard Go template functions
- Sprig template functions
- Custom helper functions

### Recovery Mechanism
- Automatic rollback on errors
- Partial generation recovery
- Transaction-like behavior

### Memory Mode
- Generate projects in memory for preview
- Used by web interface
- No disk writes until confirmed

## Usage Example

```go
import "github.com/yourusername/go-starter/internal/generator"

// Create generator
gen := generator.New(&generator.GeneratorConfig{
    OutputDir: "./output",
    DryRun:    false,
})

// Generate project
result, err := gen.Generate(blueprint, variables)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Generated project at: %s\n", result.ProjectPath)
```

## Error Handling

- Comprehensive error messages
- Rollback on failure
- Validation before generation
- Permission checks

## Dependencies

- **text/template** - Go template engine
- **github.com/Masterminds/sprig/v3** - Template functions
- **github.com/go-git/go-git/v5** - Git initialization