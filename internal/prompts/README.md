# prompts Package

This package handles all interactive user prompts for the go-starter CLI, providing an intuitive interface for project configuration.

## Overview

The prompts package implements interactive terminal prompts using the AlecAivazis/survey library, supporting both basic and advanced modes with progressive disclosure of options.

## Key Components

### Types

- **`ProjectPrompts`** - Main prompts handler
- **`PromptConfig`** - Configuration for prompt behavior
- **`ValidationRules`** - Input validation rules
- **`PromptMode`** - Basic or Advanced mode

### Core Functions

- **`New(config *PromptConfig) *ProjectPrompts`** - Create prompts instance
- **`RunInteractiveMode(mode PromptMode) (*ProjectConfig, error)`** - Run prompts
- **`GetProjectName() (string, error)`** - Prompt for project name
- **`GetBlueprintType() (string, error)`** - Select blueprint type
- **`GetArchitecture(blueprintType string) (string, error)`** - Select architecture
- **`GetFramework(blueprintType string) (string, error)`** - Select framework
- **`GetFeatures(advanced bool) (*Features, error)`** - Configure features

## Prompt Flow

### Basic Mode
1. Project name
2. Blueprint type (with descriptions)
3. Framework selection (context-aware)
4. Logger selection
5. Basic features (database, authentication)

### Advanced Mode
Includes all basic prompts plus:
- Architecture pattern selection
- Advanced database options (ORM, migrations)
- Authentication providers
- Deployment targets
- Testing configuration
- Performance optimizations

## Features

### Smart Defaults
```go
// Suggests sensible defaults based on selections
if blueprintType == "api" {
    defaultFramework = "gin"
    defaultArchitecture = "standard"
}
```

### Input Validation
```go
// Project name validation
- Must start with letter
- Alphanumeric, dash, underscore only
- No spaces or special characters
- Length between 2-50 characters

// Module path validation
- Valid Go module syntax
- No trailing slashes
- Proper domain format
```

### Context-Aware Options
- Shows only relevant frameworks for blueprint type
- Filters architecture options by blueprint
- Adjusts feature options based on selections

## Usage Example

```go
import "github.com/yourusername/go-starter/internal/prompts"

// Create prompts handler
p := prompts.New(&prompts.PromptConfig{
    UseDefaults: false,
    SkipConfirm: false,
})

// Run interactive mode
config, err := p.RunInteractiveMode(prompts.BasicMode)
if err != nil {
    log.Fatal(err)
}

// Use configuration
fmt.Printf("Creating %s project: %s\n", config.Type, config.Name)
```

## Prompt Types

### Select Prompts
- Single choice from list
- Descriptions for each option
- Search/filter support
- Default selection

### Input Prompts
- Text input with validation
- Default values
- Help text
- Transform functions

### Confirm Prompts
- Yes/No questions
- Default to safe option
- Clear messaging

### Multi-Select Prompts
- Multiple choices
- Toggle with space
- Select all/none options

## Validation Functions

```go
// Common validators
ValidateProjectName(val interface{}) error
ValidateModulePath(val interface{}) error
ValidateGoVersion(val interface{}) error
ValidatePort(val interface{}) error
```

## Error Handling

- User cancellation (Ctrl+C)
- Invalid input recovery
- Validation error messages
- Graceful exit options

## Dependencies

- **github.com/AlecAivazis/survey/v2** - Interactive prompts
- **github.com/fatih/color** - Colored output
- **github.com/muesli/termenv** - Terminal detection