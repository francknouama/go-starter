# Blueprint Externalization Plan

This document outlines the comprehensive plan to externalize the blueprints into their own GitHub repository and transform the current engine into a plugin-based system.

## Overview

The current blueprint system embeds 9 sophisticated blueprints directly into the binary. This plan details the migration to an external, pluggable architecture that enables:
- Community blueprint contributions
- Faster blueprint updates without tool releases
- Custom organizational blueprint repositories
- Blueprint versioning and dependency management

## Current System Analysis

**Embedded Blueprints (9 total):**
- `cli-standard` - CLI applications with Cobra framework
- `grpc-gateway-standard` - gRPC services with REST gateway
- `lambda-standard` - AWS Lambda functions
- `library-standard` - Go libraries with examples
- `microservice-standard` - Microservice architecture
- `web-api-clean` - Clean Architecture REST APIs (4 layers)
- `web-api-ddd` - Domain-Driven Design APIs
- `web-api-standard` - Standard REST APIs

**Current Registry System:**
- Location: `internal/templates/registry.go`
- Embedded via: `embed.go` with `//go:embed blueprints/*`
- Templates loaded at startup via `loadEmbeddedTemplates()`

## 1. Create New GitHub Repository for Blueprints

### 1.1 Repository Structure
```
go-starter-blueprints/
├── README.md
├── LICENSE
├── .github/
│   ├── workflows/
│   │   ├── validate-blueprints.yml
│   │   ├── release.yml
│   │   └── test-blueprints.yml
│   └── ISSUE_TEMPLATE/
│       └── blueprint-request.md
├── blueprints/
│   ├── cli-standard/
│   ├── grpc-gateway-standard/
│   ├── lambda-standard/
│   ├── library-standard/
│   ├── microservice-standard/
│   ├── web-api-clean/
│   ├── web-api-ddd/
│   └── web-api-standard/
├── docs/
│   ├── BLUEPRINT_SPEC.md
│   ├── CONTRIBUTING.md
│   └── TESTING.md
├── scripts/
│   ├── validate.sh
│   └── test-generation.sh
└── metadata.yaml
```

### 1.2 Repository Setup Commands
```bash
# Create and initialize repository
gh repo create go-starter-blueprints --public --description "Official blueprints for go-starter project generator"
cd go-starter-blueprints
git init
git remote add origin https://github.com/your-org/go-starter-blueprints.git

# Add repository metadata
cat > metadata.yaml << 'EOF'
name: go-starter-blueprints
version: "1.0.0"
description: "Official blueprints for go-starter project generator"
maintained_by: "go-starter team"
compatibility: ">=1.0.0"
blueprints_count: 9
last_updated: "2024-01-01T00:00:00Z"
EOF
```

## 2. Blueprint Migration Strategy

### 2.1 Copy Current Blueprints
```bash
# Copy blueprints with git history preservation
rsync -av --progress go-starter/blueprints/ go-starter-blueprints/blueprints/

# Validate blueprint structure
for blueprint in blueprints/*/; do
  echo "Validating $blueprint"
  if [[ ! -f "$blueprint/template.yaml" ]]; then
    echo "ERROR: Missing template.yaml in $blueprint"
  fi
done
```

### 2.2 Remove Embedded Blueprints
```go
// Remove from go-starter/embed.go
// Delete line: //go:embed blueprints/*

// Update internal/templates/registry.go
func (r *Registry) loadEmbeddedTemplates() {
  // Remove embedded blueprint loading
  // Will be replaced with remote blueprint loading
}
```

## 3. Update the Generator Engine

### 3.1 Blueprint Discovery System

**New Configuration File: `~/.go-starter/config.yaml`**
```yaml
blueprints:
  repositories:
    - name: "official"
      url: "https://github.com/your-org/go-starter-blueprints"
      branch: "main"
      enabled: true
      priority: 1
    - name: "community"
      url: "https://github.com/community/custom-blueprints"
      branch: "main"
      enabled: false
      priority: 2
    - name: "enterprise"
      url: "https://github.com/company/internal-blueprints"
      branch: "main"
      enabled: true
      priority: 0
      auth:
        type: "token"
        token_env: "GITHUB_TOKEN"
  
  cache:
    directory: "~/.go-starter/cache"
    ttl: "24h"
    auto_update: true
  
  defaults:
    repository: "official"
    blueprint: "web-api-standard"
```

**New Remote Registry Implementation:**
```go
// internal/remote/registry.go
package remote

import (
    "context"
    "fmt"
    "net/http"
    "time"
    
    "github.com/go-git/go-git/v5"
    "github.com/francknouama/go-starter/pkg/types"
)

type RemoteRegistry struct {
    config     *Config
    cache      Cache
    httpClient *http.Client
    repos      map[string]*Repository
}

type Repository struct {
    Name     string    `yaml:"name"`
    URL      string    `yaml:"url"`
    Branch   string    `yaml:"branch"`
    Enabled  bool      `yaml:"enabled"`
    Priority int       `yaml:"priority"`
    Auth     *AuthConfig `yaml:"auth,omitempty"`
    LastSync time.Time
}

type AuthConfig struct {
    Type     string `yaml:"type"`     // "token", "basic", "ssh"
    Token    string `yaml:"token"`
    TokenEnv string `yaml:"token_env"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
}

func NewRemoteRegistry(config *Config) *RemoteRegistry {
    return &RemoteRegistry{
        config:     config,
        cache:      NewFileCache(config.Cache.Directory),
        httpClient: &http.Client{Timeout: 30 * time.Second},
        repos:      make(map[string]*Repository),
    }
}

func (r *RemoteRegistry) LoadRepositories(ctx context.Context) error {
    for _, repo := range r.config.Repositories {
        if !repo.Enabled {
            continue
        }
        
        if err := r.syncRepository(ctx, repo); err != nil {
            return fmt.Errorf("failed to sync repository %s: %w", repo.Name, err)
        }
        
        r.repos[repo.Name] = repo
    }
    return nil
}

func (r *RemoteRegistry) syncRepository(ctx context.Context, repo *Repository) error {
    localPath := r.cache.RepositoryPath(repo.Name)
    
    // Check if repository exists locally
    if r.cache.Exists(repo.Name) {
        return r.updateRepository(ctx, repo, localPath)
    }
    
    return r.cloneRepository(ctx, repo, localPath)
}
```

### 3.2 Blueprint Caching System

**Cache Implementation:**
```go
// internal/cache/cache.go
package cache

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "time"
    
    "github.com/francknouama/go-starter/pkg/types"
)

type FileCache struct {
    baseDir string
}

type CacheEntry struct {
    Blueprint  types.Template `json:"blueprint"`
    Repository string         `json:"repository"`
    Version    string         `json:"version"`
    CachedAt   time.Time      `json:"cached_at"`
    ExpiresAt  time.Time      `json:"expires_at"`
}

func NewFileCache(baseDir string) Cache {
    return &FileCache{baseDir: baseDir}
}

func (c *FileCache) Get(repo, blueprintID string) (*CacheEntry, error) {
    cachePath := c.blueprintPath(repo, blueprintID)
    
    data, err := os.ReadFile(cachePath)
    if err != nil {
        return nil, err
    }
    
    var entry CacheEntry
    if err := json.Unmarshal(data, &entry); err != nil {
        return nil, err
    }
    
    // Check if cache entry is expired
    if time.Now().After(entry.ExpiresAt) {
        return nil, fmt.Errorf("cache entry expired")
    }
    
    return &entry, nil
}

func (c *FileCache) Set(repo, blueprintID string, blueprint types.Template, ttl time.Duration) error {
    entry := CacheEntry{
        Blueprint:  blueprint,
        Repository: repo,
        Version:    blueprint.Version,
        CachedAt:   time.Now(),
        ExpiresAt:  time.Now().Add(ttl),
    }
    
    data, err := json.Marshal(entry)
    if err != nil {
        return err
    }
    
    cachePath := c.blueprintPath(repo, blueprintID)
    if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
        return err
    }
    
    return os.WriteFile(cachePath, data, 0644)
}

func (c *FileCache) blueprintPath(repo, blueprintID string) string {
    return filepath.Join(c.baseDir, "repositories", repo, "blueprints", blueprintID+".json")
}
```

### 3.3 Enhanced Commands

**Updated `list` Command:**
```go
// cmd/list.go additions
func listRemoteBlueprints(cmd *cobra.Command, args []string) error {
    registry, err := remote.NewRegistryFromConfig()
    if err != nil {
        return err
    }
    
    if err := registry.LoadRepositories(cmd.Context()); err != nil {
        return err
    }
    
    blueprints, err := registry.ListAllBlueprints()
    if err != nil {
        return err
    }
    
    // Group by repository for display
    byRepo := make(map[string][]types.Template)
    for _, bp := range blueprints {
        byRepo[bp.Repository] = append(byRepo[bp.Repository], bp)
    }
    
    for repo, bps := range byRepo {
        fmt.Printf("\n%s Repository (%s):\n", 
            color.New(color.FgCyan, color.Bold).Sprint(repo),
            registry.GetRepositoryURL(repo))
        
        for _, bp := range bps {
            fmt.Printf("  %-20s %s\n", bp.ID, bp.Description)
        }
    }
    
    return nil
}
```

**Blueprint Selection Enhancement:**
```go
// internal/prompts/blueprint.go
func (p *Prompter) SelectRemoteBlueprint(registry *remote.RemoteRegistry) (types.Template, error) {
    repositories, err := registry.ListRepositories()
    if err != nil {
        return types.Template{}, err
    }
    
    // First, select repository
    var repoOptions []string
    for _, repo := range repositories {
        repoOptions = append(repoOptions, 
            fmt.Sprintf("%s (%s)", repo.Name, repo.URL))
    }
    
    repoIndex, err := p.Select("Select blueprint repository:", repoOptions)
    if err != nil {
        return types.Template{}, err
    }
    
    selectedRepo := repositories[repoIndex]
    
    // Then, select blueprint from repository
    blueprints, err := registry.ListBlueprints(selectedRepo.Name)
    if err != nil {
        return types.Template{}, err
    }
    
    var bpOptions []string
    for _, bp := range blueprints {
        bpOptions = append(bpOptions, 
            fmt.Sprintf("%s - %s", bp.ID, bp.Description))
    }
    
    bpIndex, err := p.Select("Select blueprint:", bpOptions)
    if err != nil {
        return types.Template{}, err
    }
    
    return blueprints[bpIndex], nil
}
```

## 4. Plugin System Implementation

### 4.1 Blueprint Management Commands

**New CLI Commands:**
```bash
# Repository management
go-starter blueprint repo add official https://github.com/your-org/go-starter-blueprints
go-starter blueprint repo add enterprise https://github.com/company/internal-blueprints --auth token
go-starter blueprint repo list
go-starter blueprint repo remove enterprise
go-starter blueprint repo sync official

# Blueprint management
go-starter blueprint list
go-starter blueprint list --repo official
go-starter blueprint sync
go-starter blueprint validate my-custom-blueprint
go-starter blueprint init my-blueprint-repo

# Cache management
go-starter blueprint cache clear
go-starter blueprint cache status
go-starter blueprint cache rebuild
```

**Implementation:**
```go
// cmd/blueprint.go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/francknouama/go-starter/internal/remote"
)

var blueprintCmd = &cobra.Command{
    Use:   "blueprint",
    Short: "Manage blueprint repositories and cache",
    Long:  `Commands for managing remote blueprint repositories, caching, and validation.`,
}

var blueprintRepoCmd = &cobra.Command{
    Use:   "repo",
    Short: "Manage blueprint repositories",
}

var blueprintRepoAddCmd = &cobra.Command{
    Use:   "add [name] [url]",
    Short: "Add a new blueprint repository",
    Args:  cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
        name, url := args[0], args[1]
        
        config, err := remote.LoadConfig()
        if err != nil {
            return err
        }
        
        authType, _ := cmd.Flags().GetString("auth")
        branch, _ := cmd.Flags().GetString("branch")
        
        repo := &remote.Repository{
            Name:    name,
            URL:     url,
            Branch:  branch,
            Enabled: true,
            Auth:    remote.ParseAuthConfig(authType),
        }
        
        if err := config.AddRepository(repo); err != nil {
            return err
        }
        
        if err := config.Save(); err != nil {
            return err
        }
        
        fmt.Printf("Successfully added repository '%s'\n", name)
        return nil
    },
}

func init() {
    blueprintRepoAddCmd.Flags().String("auth", "", "Authentication type (token, basic, ssh)")
    blueprintRepoAddCmd.Flags().String("branch", "main", "Repository branch")
    
    blueprintRepoCmd.AddCommand(blueprintRepoAddCmd)
    blueprintCmd.AddCommand(blueprintRepoCmd)
    rootCmd.AddCommand(blueprintCmd)
}
```

### 4.2 Custom Blueprint Repository Structure

**Blueprint Repository Template:**
```
my-custom-blueprints/
├── metadata.yaml
├── README.md
├── LICENSE
├── .github/
│   └── workflows/
│       └── validate.yml
├── blueprints/
│   ├── nextjs-api/
│   │   ├── template.yaml
│   │   ├── package.json.tmpl
│   │   ├── pages/
│   │   │   └── api/
│   │   │       └── hello.ts.tmpl
│   │   └── README.md.tmpl
│   └── fastapi-service/
│       ├── template.yaml
│       ├── main.py.tmpl
│       ├── requirements.txt.tmpl
│       └── Dockerfile.tmpl
└── docs/
    ├── BLUEPRINT_GUIDE.md
    └── TESTING.md
```

**Custom Blueprint Validation:**
```yaml
# .github/workflows/validate.yml
name: Validate Blueprints
on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install go-starter
        run: |
          go install github.com/your-org/go-starter@latest
      
      - name: Validate blueprint metadata
        run: |
          for blueprint in blueprints/*/; do
            echo "Validating $blueprint"
            go-starter blueprint validate "$blueprint"
          done
      
      - name: Test blueprint generation
        run: |
          for blueprint in blueprints/*/; do
            blueprint_name=$(basename "$blueprint")
            echo "Testing generation for $blueprint_name"
            
            # Generate project from blueprint
            go-starter new test-$blueprint_name \
              --blueprint local:$blueprint_name \
              --no-interactive \
              --output /tmp/test-$blueprint_name
            
            # Verify generated project compiles
            cd /tmp/test-$blueprint_name
            if [[ -f "go.mod" ]]; then
              go build ./...
            elif [[ -f "package.json" ]]; then
              npm install && npm run build
            elif [[ -f "requirements.txt" ]]; then
              pip install -r requirements.txt
            fi
          done
```

### 4.3 Blueprint Discovery API

**Repository Metadata API:**
```go
// internal/remote/api.go
package remote

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type RepositoryMetadata struct {
    Name            string    `json:"name"`
    Version         string    `json:"version"`
    Description     string    `json:"description"`
    MaintainedBy    string    `json:"maintained_by"`
    Compatibility   string    `json:"compatibility"`
    BlueprintsCount int       `json:"blueprints_count"`
    LastUpdated     time.Time `json:"last_updated"`
    Blueprints      []BlueprintInfo `json:"blueprints"`
}

type BlueprintInfo struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Version     string   `json:"version"`
    Tags        []string `json:"tags"`
    Language    string   `json:"language"`
    Framework   string   `json:"framework"`
    Complexity  string   `json:"complexity"` // "beginner", "intermediate", "advanced"
}

func (r *RemoteRegistry) FetchRepositoryMetadata(repoURL string) (*RepositoryMetadata, error) {
    metadataURL := fmt.Sprintf("%s/raw/main/metadata.yaml", repoURL)
    
    resp, err := r.httpClient.Get(metadataURL)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch metadata: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("metadata not found: %s", resp.Status)
    }
    
    var metadata RepositoryMetadata
    if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
        return nil, fmt.Errorf("failed to parse metadata: %w", err)
    }
    
    return &metadata, nil
}

func (r *RemoteRegistry) DiscoverBlueprints(repoPath string) ([]BlueprintInfo, error) {
    blueprintsPath := filepath.Join(repoPath, "blueprints")
    
    entries, err := os.ReadDir(blueprintsPath)
    if err != nil {
        return nil, err
    }
    
    var blueprints []BlueprintInfo
    for _, entry := range entries {
        if !entry.IsDir() {
            continue
        }
        
        templatePath := filepath.Join(blueprintsPath, entry.Name(), "template.yaml")
        template, err := r.loadTemplate(templatePath)
        if err != nil {
            continue // Skip invalid blueprints
        }
        
        blueprint := BlueprintInfo{
            ID:          template.ID,
            Name:        template.Name,
            Description: template.Description,
            Version:     template.Version,
            Tags:        template.Tags,
            Language:    template.Language,
            Framework:   template.Framework,
            Complexity:  template.Complexity,
        }
        
        blueprints = append(blueprints, blueprint)
    }
    
    return blueprints, nil
}
```

## 5. Configuration Schema & Examples

### 5.1 Complete Configuration Schema

```yaml
# ~/.go-starter/config.yaml
version: "1.0"

blueprints:
  repositories:
    - name: "official"                    # Repository identifier
      url: "https://github.com/your-org/go-starter-blueprints"
      branch: "main"                      # Git branch to use
      enabled: true                       # Enable/disable repository
      priority: 1                         # Lower = higher priority
      description: "Official go-starter blueprints"
      
    - name: "community"
      url: "https://github.com/community/awesome-go-blueprints"
      branch: "main"
      enabled: false
      priority: 2
      description: "Community contributed blueprints"
      
    - name: "enterprise"
      url: "https://github.com/company/internal-blueprints"
      branch: "main"
      enabled: true
      priority: 0                         # Highest priority
      description: "Internal company blueprints"
      auth:
        type: "token"                     # "token", "basic", "ssh"
        token_env: "GITHUB_TOKEN"         # Environment variable
        # Alternative auth methods:
        # token: "ghp_xxxxxxxxxxxx"       # Direct token (not recommended)
        # username: "user"                # For basic auth
        # password_env: "GITHUB_PASSWORD" # For basic auth
        # ssh_key_path: "~/.ssh/id_rsa"   # For SSH auth
  
  cache:
    directory: "~/.go-starter/cache"      # Cache location
    ttl: "24h"                           # Time to live for cache entries
    max_size: "500MB"                    # Maximum cache size
    auto_update: true                    # Auto-update on list/new commands
    compress: true                       # Compress cached blueprints
    
  defaults:
    repository: "official"               # Default repository for new projects
    blueprint: "web-api-standard"        # Default blueprint
    
  filters:
    language: []                         # Filter by language: ["go", "typescript"]
    complexity: []                       # Filter by complexity: ["beginner", "intermediate"]
    tags: []                            # Filter by tags: ["api", "microservice"]

# CLI behavior configuration
cli:
  interactive_mode: true                 # Enable interactive prompts
  color_output: true                    # Enable colored output
  progress_bars: true                   # Show progress bars
  auto_completion: true                 # Enable shell completion
  
# Logging configuration
logging:
  level: "info"                         # debug, info, warn, error
  format: "console"                     # console, json
  file: ""                             # Log to file (empty = stdout only)

# Project generation defaults
generation:
  output_directory: "./{{.ProjectName}}" # Default output pattern
  overwrite_protection: true            # Prompt before overwriting
  git_init: true                       # Initialize git repository
  git_initial_commit: true             # Create initial commit
  validate_generated: true             # Validate generated projects
  
  hooks:
    pre_generation: []                  # Commands to run before generation
    post_generation:                    # Commands to run after generation
      - "go mod download"
      - "go build ./..."
      
# Update checking
updates:
  check_frequency: "daily"             # never, daily, weekly, monthly
  auto_update: false                   # Auto-update go-starter
  include_prereleases: false           # Include beta/rc versions
```

### 5.2 Environment Variable Overrides

```bash
# Override configuration via environment variables
export GO_STARTER_CONFIG_DIR="$HOME/.config/go-starter"
export GO_STARTER_CACHE_DIR="$HOME/.cache/go-starter"
export GO_STARTER_DEFAULT_REPO="enterprise"
export GO_STARTER_LOG_LEVEL="debug"
export GO_STARTER_INTERACTIVE="false"

# Authentication
export GITHUB_TOKEN="ghp_xxxxxxxxxxxx"
export GITLAB_TOKEN="glpat-xxxxxxxxxxxx"
export CUSTOM_REPO_TOKEN="custom_token_here"
```

## 6. Testing Strategy

### 6.1 Test Categories

**Unit Tests:**
```go
// internal/remote/registry_test.go
func TestRemoteRegistry_LoadRepositories(t *testing.T) {
    tests := []struct {
        name        string
        config      *Config
        expectError bool
        expectRepos int
    }{
        {
            name: "load single enabled repository",
            config: &Config{
                Repositories: []*Repository{
                    {Name: "test", URL: "https://github.com/test/repo", Enabled: true},
                },
            },
            expectError: false,
            expectRepos: 1,
        },
        {
            name: "skip disabled repositories",
            config: &Config{
                Repositories: []*Repository{
                    {Name: "enabled", URL: "https://github.com/test/enabled", Enabled: true},
                    {Name: "disabled", URL: "https://github.com/test/disabled", Enabled: false},
                },
            },
            expectError: false,
            expectRepos: 1,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            registry := NewRemoteRegistry(tt.config)
            err := registry.LoadRepositories(context.Background())
            
            if tt.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Len(t, registry.repos, tt.expectRepos)
            }
        })
    }
}
```

**Integration Tests:**
```go
// tests/integration/remote_blueprints_test.go
func TestRemoteBlueprintGeneration(t *testing.T) {
    // Setup test repository
    testRepo := setupTestRepository(t)
    defer cleanupTestRepository(t, testRepo)
    
    // Configure registry with test repository
    config := &remote.Config{
        Repositories: []*remote.Repository{
            {Name: "test", URL: testRepo.URL, Enabled: true},
        },
    }
    
    registry := remote.NewRemoteRegistry(config)
    err := registry.LoadRepositories(context.Background())
    require.NoError(t, err)
    
    // Test blueprint listing
    blueprints, err := registry.ListBlueprints("test")
    require.NoError(t, err)
    assert.Greater(t, len(blueprints), 0)
    
    // Test blueprint generation
    for _, blueprint := range blueprints {
        t.Run(fmt.Sprintf("generate_%s", blueprint.ID), func(t *testing.T) {
            tempDir := t.TempDir()
            
            err := generateProject(blueprint, tempDir)
            assert.NoError(t, err)
            
            // Verify generated project structure
            assertProjectStructure(t, tempDir, blueprint.ID)
            
            // Verify project compiles
            assertProjectCompiles(t, tempDir)
        })
    }
}
```

**Cache Tests:**
```go
// internal/cache/cache_test.go
func TestFileCache_SetAndGet(t *testing.T) {
    tempDir := t.TempDir()
    cache := NewFileCache(tempDir)
    
    blueprint := types.Template{
        ID:          "test-blueprint",
        Name:        "Test Blueprint",
        Description: "A test blueprint",
        Version:     "1.0.0",
    }
    
    // Test cache set
    err := cache.Set("test-repo", "test-blueprint", blueprint, time.Hour)
    assert.NoError(t, err)
    
    // Test cache get
    entry, err := cache.Get("test-repo", "test-blueprint")
    assert.NoError(t, err)
    assert.Equal(t, blueprint.ID, entry.Blueprint.ID)
    assert.Equal(t, "test-repo", entry.Repository)
}

func TestFileCache_Expiration(t *testing.T) {
    tempDir := t.TempDir()
    cache := NewFileCache(tempDir)
    
    blueprint := types.Template{ID: "test"}
    
    // Set with very short TTL
    err := cache.Set("repo", "test", blueprint, time.Millisecond)
    assert.NoError(t, err)
    
    // Wait for expiration
    time.Sleep(time.Millisecond * 10)
    
    // Should return error for expired cache
    _, err = cache.Get("repo", "test")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "expired")
}
```

### 6.2 Validation Framework

```go
// internal/validation/blueprint.go
package validation

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    
    "gopkg.in/yaml.v3"
    "github.com/francknouama/go-starter/pkg/types"
)

type BlueprintValidator struct {
    rules []ValidationRule
}

type ValidationRule interface {
    Validate(blueprint *types.Template, blueprintPath string) []ValidationError
}

type ValidationError struct {
    Rule     string `json:"rule"`
    Severity string `json:"severity"` // "error", "warning", "info"
    Message  string `json:"message"`
    File     string `json:"file,omitempty"`
    Line     int    `json:"line,omitempty"`
}

func NewBlueprintValidator() *BlueprintValidator {
    return &BlueprintValidator{
        rules: []ValidationRule{
            &TemplateYamlRule{},
            &RequiredFilesRule{},
            &TemplateConventionRule{},
            &SecurityRule{},
            &PerformanceRule{},
        },
    }
}

func (v *BlueprintValidator) ValidateBlueprint(blueprintPath string) ([]ValidationError, error) {
    // Load template.yaml
    templatePath := filepath.Join(blueprintPath, "template.yaml")
    template, err := v.loadTemplate(templatePath)
    if err != nil {
        return nil, fmt.Errorf("failed to load template: %w", err)
    }
    
    var allErrors []ValidationError
    
    // Run all validation rules
    for _, rule := range v.rules {
        errors := rule.Validate(template, blueprintPath)
        allErrors = append(allErrors, errors...)
    }
    
    return allErrors, nil
}

// Template YAML validation rule
type TemplateYamlRule struct{}

func (r *TemplateYamlRule) Validate(template *types.Template, blueprintPath string) []ValidationError {
    var errors []ValidationError
    
    if template.ID == "" {
        errors = append(errors, ValidationError{
            Rule:     "template-yaml",
            Severity: "error",
            Message:  "template.yaml must have a non-empty 'id' field",
            File:     "template.yaml",
        })
    }
    
    if template.Name == "" {
        errors = append(errors, ValidationError{
            Rule:     "template-yaml",
            Severity: "error",
            Message:  "template.yaml must have a non-empty 'name' field",
            File:     "template.yaml",
        })
    }
    
    if template.Description == "" {
        errors = append(errors, ValidationError{
            Rule:     "template-yaml",
            Severity: "warning",
            Message:  "template.yaml should have a 'description' field",
            File:     "template.yaml",
        })
    }
    
    return errors
}
```

## 7. Error Handling & Migration Strategy

### 7.1 Graceful Degradation

```go
// internal/fallback/registry.go
package fallback

import (
    "fmt"
    "log"
    
    "github.com/francknouama/go-starter/internal/remote"
    "github.com/francknouama/go-starter/internal/templates"
    "github.com/francknouama/go-starter/pkg/types"
)

type HybridRegistry struct {
    remote   *remote.RemoteRegistry
    embedded *templates.Registry
    fallback bool
}

func NewHybridRegistry(config *remote.Config) *HybridRegistry {
    return &HybridRegistry{
        remote:   remote.NewRemoteRegistry(config),
        embedded: templates.NewRegistry(), // Contains embedded blueprints as fallback
        fallback: config.EnableFallback,
    }
}

func (h *HybridRegistry) GetTemplate(templateID string) (types.Template, error) {
    // Try remote registry first
    template, err := h.remote.Get(templateID)
    if err == nil {
        return template, nil
    }
    
    // Log remote failure
    log.Printf("Remote blueprint fetch failed for %s: %v", templateID, err)
    
    // Fallback to embedded blueprints if enabled
    if h.fallback {
        log.Printf("Falling back to embedded blueprint for %s", templateID)
        return h.embedded.Get(templateID)
    }
    
    return types.Template{}, fmt.Errorf("blueprint %s not found in any repository", templateID)
}

func (h *HybridRegistry) ListTemplates() ([]types.Template, error) {
    var allTemplates []types.Template
    
    // Get remote templates
    remoteTemplates, err := h.remote.ListAll()
    if err != nil {
        log.Printf("Failed to list remote templates: %v", err)
    } else {
        allTemplates = append(allTemplates, remoteTemplates...)
    }
    
    // Add embedded templates if fallback enabled or remote failed
    if h.fallback || err != nil {
        embeddedTemplates, embeddedErr := h.embedded.ListAll()
        if embeddedErr != nil {
            log.Printf("Failed to list embedded templates: %v", embeddedErr)
        } else {
            // Deduplicate templates (remote takes precedence)
            for _, embedded := range embeddedTemplates {
                if !h.containsTemplate(allTemplates, embedded.ID) {
                    allTemplates = append(allTemplates, embedded)
                }
            }
        }
    }
    
    if len(allTemplates) == 0 {
        return nil, fmt.Errorf("no templates available from any source")
    }
    
    return allTemplates, nil
}
```

### 7.2 Migration Script

```bash
#!/bin/bash
# scripts/migrate-to-remote-blueprints.sh

set -e

echo "🚀 Migrating go-starter to remote blueprint system..."

# Backup current configuration
if [[ -f "$HOME/.go-starter.yaml" ]]; then
    echo "📋 Backing up existing configuration..."
    cp "$HOME/.go-starter.yaml" "$HOME/.go-starter.yaml.backup"
fi

# Create new configuration directory
echo "📁 Setting up configuration directory..."
mkdir -p "$HOME/.go-starter"

# Generate default configuration
echo "⚙️  Generating default configuration..."
cat > "$HOME/.go-starter/config.yaml" << 'EOF'
version: "1.0"
blueprints:
  repositories:
    - name: "official"
      url: "https://github.com/your-org/go-starter-blueprints"
      branch: "main"
      enabled: true
      priority: 1
      description: "Official go-starter blueprints"
  
  cache:
    directory: "~/.go-starter/cache"
    ttl: "24h"
    auto_update: true
  
  defaults:
    repository: "official"
    blueprint: "web-api-standard"

cli:
  interactive_mode: true
  color_output: true

generation:
  git_init: true
  validate_generated: true
EOF

# Initialize cache directory
echo "🗂️  Initializing cache directory..."
mkdir -p "$HOME/.go-starter/cache"

# Test configuration
echo "🧪 Testing new configuration..."
if command -v go-starter >/dev/null 2>&1; then
    echo "✅ Testing blueprint listing..."
    go-starter blueprint list --timeout 10s || {
        echo "⚠️  Remote blueprint fetch failed, enabling fallback mode..."
        
        # Enable fallback in configuration
        sed -i.bak 's/auto_update: true/auto_update: true\n    enable_fallback: true/' "$HOME/.go-starter/config.yaml"
    }
else
    echo "ℹ️  go-starter not found in PATH, skipping test"
fi

echo "✅ Migration completed successfully!"
echo ""
echo "📖 Next steps:"
echo "   1. Run 'go-starter blueprint list' to see available blueprints"
echo "   2. Run 'go-starter blueprint repo list' to manage repositories"
echo "   3. Check '$HOME/.go-starter/config.yaml' to customize settings"
echo ""
echo "🔄 To revert: restore from '$HOME/.go-starter.yaml.backup'"
```

### 7.3 Rollback Strategy

```go
// internal/migration/rollback.go
package migration

import (
    "fmt"
    "os"
    "path/filepath"
)

type RollbackManager struct {
    configDir string
    backupDir string
}

func NewRollbackManager() *RollbackManager {
    homeDir, _ := os.UserHomeDir()
    return &RollbackManager{
        configDir: filepath.Join(homeDir, ".go-starter"),
        backupDir: filepath.Join(homeDir, ".go-starter", "backups"),
    }
}

func (r *RollbackManager) CreateBackup(version string) error {
    backupPath := filepath.Join(r.backupDir, fmt.Sprintf("backup-%s", version))
    
    if err := os.MkdirAll(backupPath, 0755); err != nil {
        return err
    }
    
    // Backup configuration
    configPath := filepath.Join(r.configDir, "config.yaml")
    if _, err := os.Stat(configPath); err == nil {
        return r.copyFile(configPath, filepath.Join(backupPath, "config.yaml"))
    }
    
    return nil
}

func (r *RollbackManager) Rollback(version string) error {
    backupPath := filepath.Join(r.backupDir, fmt.Sprintf("backup-%s", version))
    
    if _, err := os.Stat(backupPath); os.IsNotExist(err) {
        return fmt.Errorf("backup version %s not found", version)
    }
    
    // Restore configuration
    backupConfigPath := filepath.Join(backupPath, "config.yaml")
    configPath := filepath.Join(r.configDir, "config.yaml")
    
    return r.copyFile(backupConfigPath, configPath)
}
```

## 8. Implementation Phases

### Phase 1: Core Infrastructure (Weeks 1-2)

**Deliverables:**
- [ ] Create `go-starter-blueprints` repository with all existing blueprints
- [ ] Implement `internal/remote/` package with basic repository cloning
- [ ] Implement `internal/cache/` package with file-based caching
- [ ] Add new configuration schema support (`~/.go-starter/config.yaml`)
- [ ] Update `list` command to support remote repositories

**Implementation Steps:**
```bash
# Week 1: Repository setup and basic remote functionality
1. Create go-starter-blueprints repository
2. Migrate existing blueprints with validation
3. Implement basic git cloning and caching
4. Add configuration file support

# Week 2: Integration and commands
5. Update list command for remote blueprints
6. Add basic repository management commands
7. Implement hybrid registry (remote + embedded fallback)
8. Add comprehensive unit tests
```

**Success Criteria:**
- Users can list blueprints from remote repository
- Caching works correctly with TTL expiration
- Fallback to embedded blueprints when remote fails
- All existing functionality continues to work

### Phase 2: Full Plugin System (Weeks 3-4)

**Deliverables:**
- [ ] Complete blueprint management CLI commands
- [ ] Authentication support for private repositories
- [ ] Blueprint validation framework
- [ ] Enhanced error handling and recovery
- [ ] Repository priority and filtering system

**Implementation Steps:**
```bash
# Week 3: Plugin system core
1. Implement all blueprint management commands
2. Add authentication (token, SSH, basic auth)
3. Build blueprint validation framework
4. Add repository priority handling

# Week 4: Advanced features and polish
5. Implement blueprint filtering and search
6. Add repository metadata API
7. Build comprehensive error handling
8. Create migration and rollback tools
```

**Success Criteria:**
- Users can add/remove custom blueprint repositories
- Private repository authentication works
- Blueprint validation catches common issues
- Error handling provides clear guidance

### Phase 3: Documentation & Testing (Week 5)

**Deliverables:**
- [ ] Comprehensive testing suite (unit, integration, e2e)
- [ ] Updated documentation and guides
- [ ] Performance optimization and benchmarking
- [ ] Security audit and hardening

**Implementation Steps:**
```bash
# Week 5: Quality assurance
1. Build comprehensive test suite (80%+ coverage)
2. Update all documentation with new features
3. Performance testing and optimization
4. Security review and fixes
5. Beta testing with community
```

**Success Criteria:**
- Test coverage >80% for all new code
- Documentation complete and tested
- Performance benchmarks meet targets
- Security review passes

## 9. Success Metrics & Monitoring

### 9.1 Key Performance Indicators (KPIs)

**User Adoption:**
- Blueprint repository additions per week
- Community blueprint contributions
- Custom blueprint usage vs official blueprints
- User retention after remote blueprint introduction

**System Performance:**
- Repository sync time (target: <30s for official repo)
- Cache hit rate (target: >90% for repeated operations)
- Blueprint generation time (should not increase >10%)
- Network failure recovery rate (target: >95% fallback success)

**Developer Experience:**
- Command completion time (target: <5s for list commands)
- Error recovery success rate
- Documentation page views and engagement
- Support ticket reduction related to blueprint issues

### 9.2 Monitoring Implementation

```go
// internal/telemetry/metrics.go
package telemetry

import (
    "context"
    "time"
    
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    blueprintSyncDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "go_starter_blueprint_sync_duration_seconds",
            Help: "Time taken to sync blueprint repositories",
        },
        []string{"repository", "status"},
    )
    
    cacheHitRate = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "go_starter_cache_hits_total",
            Help: "Number of cache hits",
        },
        []string{"repository", "blueprint"},
    )
    
    repositoryErrors = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "go_starter_repository_errors_total",
            Help: "Number of repository access errors",
        },
        []string{"repository", "error_type"},
    )
)

func RecordSyncDuration(repo string, duration time.Duration, success bool) {
    status := "success"
    if !success {
        status = "failure"
    }
    blueprintSyncDuration.WithLabelValues(repo, status).Observe(duration.Seconds())
}

func RecordCacheHit(repo, blueprint string) {
    cacheHitRate.WithLabelValues(repo, blueprint).Inc()
}

func RecordRepositoryError(repo, errorType string) {
    repositoryErrors.WithLabelValues(repo, errorType).Inc()
}
```

### 9.3 Health Checks

```go
// cmd/health.go
package cmd

import (
    "context"
    "fmt"
    "time"
    
    "github.com/spf13/cobra"
    "github.com/francknouama/go-starter/internal/remote"
)

var healthCmd = &cobra.Command{
    Use:   "health",
    Short: "Check system health and repository status",
    RunE: func(cmd *cobra.Command, args []string) error {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        return runHealthCheck(ctx)
    },
}

func runHealthCheck(ctx context.Context) error {
    config, err := remote.LoadConfig()
    if err != nil {
        return fmt.Errorf("❌ Configuration load failed: %w", err)
    }
    
    fmt.Printf("🏥 go-starter Health Check\n\n")
    
    registry := remote.NewRemoteRegistry(config)
    
    // Check each repository
    for _, repo := range config.Repositories {
        if !repo.Enabled {
            fmt.Printf("⏭️  %s: SKIPPED (disabled)\n", repo.Name)
            continue
        }
        
        start := time.Now()
        err := registry.TestConnection(ctx, repo)
        duration := time.Since(start)
        
        if err != nil {
            fmt.Printf("❌ %s: FAILED (%s) - %v\n", repo.Name, duration, err)
        } else {
            fmt.Printf("✅ %s: OK (%s)\n", repo.Name, duration)
        }
    }
    
    // Check cache
    cacheStats, err := registry.GetCacheStats()
    if err != nil {
        fmt.Printf("❌ Cache: FAILED - %v\n", err)
    } else {
        fmt.Printf("📦 Cache: %d entries, %s used\n", 
            cacheStats.EntryCount, 
            formatBytes(cacheStats.SizeBytes))
    }
    
    return nil
}

func init() {
    rootCmd.AddCommand(healthCmd)
}
```

## 10. Security Considerations

### 10.1 Security Model

**Authentication Security:**
- Never store tokens in configuration files (use environment variables)
- Support credential managers (macOS Keychain, Windows Credential Store)
- Implement token expiry detection and refresh
- Add audit logging for authentication attempts

**Repository Security:**
```go
// internal/security/repository.go
package security

import (
    "crypto/sha256"
    "fmt"
    "net/url"
    "strings"
)

type RepositoryValidator struct {
    allowedHosts []string
    blockedHosts []string
}

func NewRepositoryValidator() *RepositoryValidator {
    return &RepositoryValidator{
        allowedHosts: []string{
            "github.com",
            "gitlab.com",
            "bitbucket.org",
        },
        blockedHosts: []string{
            "localhost",
            "127.0.0.1",
            "10.*",
            "192.168.*",
            "172.16.*",
        },
    }
}

func (v *RepositoryValidator) ValidateRepositoryURL(repoURL string) error {
    u, err := url.Parse(repoURL)
    if err != nil {
        return fmt.Errorf("invalid repository URL: %w", err)
    }
    
    // Check protocol
    if u.Scheme != "https" && u.Scheme != "ssh" {
        return fmt.Errorf("repository URL must use https or ssh protocol")
    }
    
    // Check against allowed hosts
    host := strings.ToLower(u.Host)
    for _, allowedHost := range v.allowedHosts {
        if strings.Contains(host, allowedHost) {
            return nil
        }
    }
    
    // Check against blocked hosts
    for _, blockedHost := range v.blockedHosts {
        if strings.Contains(host, blockedHost) {
            return fmt.Errorf("repository host %s is not allowed", host)
        }
    }
    
    return fmt.Errorf("repository host %s is not in the allowed list", host)
}

func (v *RepositoryValidator) GenerateRepositoryChecksum(repoURL string) string {
    hash := sha256.Sum256([]byte(repoURL))
    return fmt.Sprintf("%x", hash[:16]) // First 16 bytes as hex
}
```

### 10.2 Blueprint Security Scanning

```go
// internal/security/blueprint.go
package security

import (
    "bufio"
    "os"
    "path/filepath"
    "regexp"
    "strings"
)

type BlueprintSecurityScanner struct {
    dangerousPatterns []*regexp.Regexp
    suspiciousFiles   map[string]bool
}

func NewBlueprintSecurityScanner() *BlueprintSecurityScanner {
    return &BlueprintSecurityScanner{
        dangerousPatterns: []*regexp.Regexp{
            regexp.MustCompile(`rm\s+-rf\s+/`),                    // Dangerous file operations
            regexp.MustCompile(`\bcurl\s+.*\|\s*sh`),              // Pipe to shell
            regexp.MustCompile(`eval\s*\(\s*.*\)`),                // Eval statements
            regexp.MustCompile(`\$\(.*\)`),                        // Command substitution
            regexp.MustCompile(`\b(password|secret|token)\s*=`),   // Hardcoded secrets
        },
        suspiciousFiles: map[string]bool{
            ".env":         true,
            "id_rsa":       true,
            "private.key":  true,
            "secrets.yaml": true,
        },
    }
}

func (s *BlueprintSecurityScanner) ScanBlueprint(blueprintPath string) ([]SecurityIssue, error) {
    var issues []SecurityIssue
    
    err := filepath.Walk(blueprintPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if info.IsDir() {
            return nil
        }
        
        // Check for suspicious file names
        fileName := info.Name()
        if s.suspiciousFiles[fileName] {
            issues = append(issues, SecurityIssue{
                Type:        "suspicious_file",
                Severity:    "high",
                File:        path,
                Description: fmt.Sprintf("Found suspicious file: %s", fileName),
            })
        }
        
        // Scan file contents for dangerous patterns
        if strings.HasSuffix(fileName, ".tmpl") || strings.HasSuffix(fileName, ".sh") {
            fileIssues, err := s.scanFileContent(path)
            if err != nil {
                return err
            }
            issues = append(issues, fileIssues...)
        }
        
        return nil
    })
    
    return issues, err
}

type SecurityIssue struct {
    Type        string `json:"type"`
    Severity    string `json:"severity"` // "low", "medium", "high", "critical"
    File        string `json:"file"`
    Line        int    `json:"line,omitempty"`
    Description string `json:"description"`
    Pattern     string `json:"pattern,omitempty"`
}
```

## 11. Future Enhancements

### 11.1 Marketplace Integration (Phase 4)

**Community Blueprint Marketplace:**
- Web interface for blueprint discovery and sharing
- Rating and review system for community blueprints
- Blueprint analytics and usage statistics
- Automated security scanning for community contributions

**Enterprise Features:**
- Private blueprint marketplaces for organizations
- Blueprint approval workflows
- Usage analytics and compliance reporting
- Integration with corporate identity providers

### 11.2 Advanced Features

**AI-Powered Blueprint Generation:**
```go
// Future: AI blueprint assistant
type AIBlueprintAssistant struct {
    client OpenAIClient
}

func (a *AIBlueprintAssistant) GenerateBlueprint(requirements string) (*types.Template, error) {
    prompt := fmt.Sprintf(`
    Generate a Go project blueprint based on these requirements:
    %s
    
    Return a complete template.yaml and necessary template files.
    `, requirements)
    
    response, err := a.client.Complete(prompt)
    if err != nil {
        return nil, err
    }
    
    return parseAIResponse(response)
}
```

**Blueprint Versioning and Dependencies:**
```yaml
# Future: Blueprint dependencies
dependencies:
  blueprints:
    - id: "base-go-service"
      version: ">=1.2.0"
      repository: "official"
  
  tools:
    - name: "protoc"
      version: ">=3.15.0"
    - name: "docker"
      version: ">=20.0.0"

upgrade_path:
  from: "1.0.0"
  to: "2.0.0"
  breaking_changes:
    - "Configuration format changed"
    - "Removed deprecated middleware"
```

This comprehensive blueprint externalization plan provides concrete, actionable steps to transform go-starter into a powerful, extensible project generator with a robust plugin ecosystem. The implementation is designed to be secure, performant, and user-friendly while maintaining backward compatibility.
```
