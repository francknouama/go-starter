package generator

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
)

// Generator handles project generation
type Generator struct {
	registry *templates.Registry
	loader   *templates.TemplateLoader
}

// New creates a new Generator instance
func New() *Generator {
	return &Generator{
		registry: templates.NewRegistry(),
		loader:   templates.NewTemplateLoader(),
	}
}

// Generate generates a new project based on the configuration
func (g *Generator) Generate(config types.ProjectConfig, options types.GenerationOptions) (*types.GenerationResult, error) {
	startTime := time.Now()

	result := &types.GenerationResult{
		ProjectPath:  options.OutputPath,
		FilesCreated: []string{},
		Success:      false,
	}

	// Validate configuration
	if err := g.validateConfig(config); err != nil {
		result.Error = err
		return result, err
	}

	// Check if template exists
	template, err := g.registry.Get(g.getTemplateID(config))
	if err != nil {
		// For Phase 0, we'll show a helpful message about upcoming templates
		return g.handleMissingTemplate(config, result)
	}

	// Create output directory
	if err := os.MkdirAll(options.OutputPath, 0755); err != nil {
		result.Error = types.NewFileSystemError("failed to create output directory", err)
		return result, result.Error
	}

	// Generate project files
	filesCreated, err := g.generateProjectFiles(template, config, options.OutputPath)
	if err != nil {
		result.Error = err
		return result, err
	}
	result.FilesCreated = filesCreated

	// Initialize git repository if requested
	if !options.NoGit {
		if err := g.initGitRepository(options.OutputPath); err != nil {
			// Git init failure is not fatal, just log it
			fmt.Printf("Warning: failed to initialize git repository: %v\n", err)
		}
	}

	result.Duration = time.Since(startTime)
	result.Success = true
	return result, nil
}

// Preview shows what would be generated without creating files
func (g *Generator) Preview(config types.ProjectConfig, outputDir string) error {
	fmt.Printf("Preview for project '%s':\n", config.Name)
	fmt.Printf("  Type: %s\n", config.Type)
	fmt.Printf("  Module: %s\n", config.Module)

	if config.Framework != "" {
		fmt.Printf("  Framework: %s\n", config.Framework)
	}
	if config.Architecture != "" {
		fmt.Printf("  Architecture: %s\n", config.Architecture)
	}

	outputPath := filepath.Join(outputDir, config.Name)
	fmt.Printf("  Output path: %s\n", outputPath)

	templateID := g.getTemplateID(config)
	template, err := g.registry.Get(templateID)
	if err != nil {
		fmt.Printf("\nTemplate '%s' is not available yet.\n", templateID)
		fmt.Printf("This template will be implemented in upcoming phases.\n")
		return nil
	}

	fmt.Printf("\nFiles to be generated:\n")
	for _, file := range template.Files {
		destination := g.processTemplatePath(file.Destination, config)
		fmt.Printf("  %s\n", destination)
	}

	return nil
}

// validateConfig validates the project configuration
func (g *Generator) validateConfig(config types.ProjectConfig) error {
	if config.Name == "" {
		return types.NewValidationError("project name is required", nil)
	}
	if config.Module == "" {
		return types.NewValidationError("module path is required", nil)
	}
	if config.Type == "" {
		return types.NewValidationError("project type is required", nil)
	}
	return nil
}

// getTemplateID maps project configuration to template ID
func (g *Generator) getTemplateID(config types.ProjectConfig) string {
	if config.Architecture != "" && config.Architecture != "standard" {
		return fmt.Sprintf("%s-%s", config.Type, config.Architecture)
	}
	return config.Type
}

// handleMissingTemplate provides helpful feedback when templates aren't available yet
func (g *Generator) handleMissingTemplate(config types.ProjectConfig, result *types.GenerationResult) (*types.GenerationResult, error) {
	templateID := g.getTemplateID(config)

	fmt.Printf("Template '%s' is not available yet.\n", templateID)
	fmt.Println("\ngo-starter is currently in Phase 0 (Foundation).")
	fmt.Println("Project templates will be added in the following phases:")
	fmt.Println("  • Phase 1: Web API template (Gin framework)")
	fmt.Println("  • Phase 2: CLI Application template (Cobra framework)")
	fmt.Println("  • Phase 3: Go Library template")
	fmt.Println("  • Phase 4: AWS Lambda template")
	fmt.Println("  • Phase 5+: Additional templates and architectures")

	err := types.NewTemplateNotFoundError(templateID)
	result.Error = err
	return result, err
}

// generateProjectFiles generates all files for the project
func (g *Generator) generateProjectFiles(tmpl types.Template, config types.ProjectConfig, outputPath string) ([]string, error) {
	var filesCreated []string

	// Create template context with all variables
	context := g.createTemplateContext(config, tmpl)

	// Get template directory from metadata
	templateDir, ok := tmpl.Metadata["path"].(string)
	if !ok {
		return nil, fmt.Errorf("template metadata missing path")
	}

	// Process each file in the template
	for _, templateFile := range tmpl.Files {
		// Evaluate condition if present
		if templateFile.Condition != "" {
			shouldGenerate, err := g.evaluateCondition(templateFile.Condition, context)
			if err != nil {
				return nil, fmt.Errorf("failed to evaluate condition for %s: %w", templateFile.Source, err)
			}
			if !shouldGenerate {
				continue
			}
		}

		// Process template path with variables
		destPath := g.processTemplatePath(templateFile.Destination, config)
		fullDestPath := filepath.Join(outputPath, destPath)

		// Create directory if it doesn't exist
		if err := os.MkdirAll(filepath.Dir(fullDestPath), 0755); err != nil {
			return nil, types.NewFileSystemError("failed to create directory", err)
		}

		// Generate file from template
		if err := g.processTemplateFile(templateDir, templateFile.Source, fullDestPath, context); err != nil {
			return nil, fmt.Errorf("failed to process template file %s: %w", templateFile.Source, err)
		}

		// Set executable permission if needed
		if templateFile.Executable {
			if err := os.Chmod(fullDestPath, 0755); err != nil {
				return nil, types.NewFileSystemError("failed to set executable permission", err)
			}
		}

		filesCreated = append(filesCreated, fullDestPath)
	}

	// Process dependencies
	if err := g.processDependencies(tmpl, config, outputPath, context); err != nil {
		return nil, fmt.Errorf("failed to process dependencies: %w", err)
	}

	// Execute post-generation hooks
	g.executeHooks(tmpl, config, outputPath, context)

	return filesCreated, nil
}

// createGoMod creates a basic go.mod file
func (g *Generator) createGoMod(config types.ProjectConfig, path string) error {
	content := fmt.Sprintf("module %s\n\ngo %s\n", config.Module, config.GoVersion)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return types.NewFileSystemError("failed to create go.mod", err)
	}

	return nil
}

// initGitRepository initializes a git repository
func (g *Generator) initGitRepository(projectPath string) error {
	return g.initGit(projectPath)
}

// initGit handles the actual git initialization
func (g *Generator) initGit(projectPath string) error {
	// Check if git is available
	if !g.isGitAvailable() {
		return types.NewFileSystemError("git is not available in PATH", nil)
	}

	// Check if directory already has git
	if g.hasGitRepository(projectPath) {
		return nil // Already a git repository
	}

	// Initialize git repository
	if err := g.runGitCommand(projectPath, "init"); err != nil {
		return types.NewFileSystemError("failed to initialize git repository", err)
	}

	// Create .gitignore file
	if err := g.createGitignore(projectPath); err != nil {
		return types.NewFileSystemError("failed to create .gitignore", err)
	}

	return nil
}

// processTemplatePath processes template variables in file paths
func (g *Generator) processTemplatePath(path string, config types.ProjectConfig) string {
	// Create a simple context for path processing
	context := map[string]any{
		"ProjectName":  config.Name,
		"ModulePath":   config.Module,
		"GoVersion":    config.GoVersion,
		"Framework":    config.Framework,
		"Architecture": config.Architecture,
	}

	// Use text/template to process the path
	tmpl, err := template.New("path").Funcs(sprig.FuncMap()).Parse(path)
	if err != nil {
		// Log the error but continue with original path for backwards compatibility
		fmt.Printf("Warning: Failed to parse template path %q: %v\n", path, err)
		return path
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		// Log the error but continue with original path for backwards compatibility
		fmt.Printf("Warning: Failed to execute template path %q: %v\n", path, err)
		return path
	}

	return buf.String()
}

// createTemplateContext creates a template context with all variables
func (g *Generator) createTemplateContext(config types.ProjectConfig, tmpl types.Template) map[string]any {
	context := map[string]any{
		"ProjectName":  config.Name,
		"ModulePath":   config.Module,
		"GoVersion":    config.GoVersion,
		"Framework":    config.Framework,
		"Architecture": config.Architecture,
		"Author":       config.Author,
		"Email":        config.Email,
		"License":      config.License,
		"Type":         config.Type,
		"Logger":       config.Logger,
	}

	// Add features from the config
	if config.Features != nil {
		context["Features"] = config.Features
	}

	// Add logger configuration to context
	if config.Logger != "" {
		context["LoggerType"] = config.Logger

		// Create logger configuration for templates
		if config.Features != nil && config.Features.Logging.Type != "" {
			context["LoggerConfig"] = map[string]interface{}{
				"Type":       config.Features.Logging.Type,
				"Level":      config.Features.Logging.Level,
				"Format":     config.Features.Logging.Format,
				"Structured": config.Features.Logging.Structured,
			}
		} else {
			// Default logger configuration
			context["LoggerConfig"] = map[string]interface{}{
				"Type":       config.Logger,
				"Level":      "info",
				"Format":     "json",
				"Structured": true,
			}
		}

		// Add logger-specific template variables
		switch config.Logger {
		case "slog":
			context["UseSlog"] = true
		case "zap":
			context["UseZap"] = true
		case "logrus":
			context["UseLogrus"] = true
		case "zerolog":
			context["UseZerolog"] = true
		}
	}

	// Add variables from config.Variables map
	if config.Variables != nil {
		for key, value := range config.Variables {
			context[key] = value
		}
	}

	// Add template-specific variables with defaults
	for _, variable := range tmpl.Variables {
		// Skip if already set from config
		if _, exists := context[variable.Name]; exists {
			continue
		}

		// Add default value if specified
		if variable.Default != nil {
			context[variable.Name] = variable.Default
		}
	}

	// Add convenience variables for common patterns
	// Check Features struct first, then fall back to Variables map
	dbDriver := g.getFeatureValue(config, "database", "driver", "")
	if dbDriver == "" && config.Variables != nil {
		if val, exists := config.Variables["DatabaseDriver"]; exists {
			dbDriver = val
		}
	}
	context["DatabaseDriver"] = dbDriver

	// Add multi-database support
	dbDrivers := g.getDatabaseDrivers(config)
	context["DatabaseDrivers"] = dbDrivers
	context["HasDatabase"] = len(dbDrivers) > 0 || dbDriver != ""

	// For backward compatibility, ensure DatabaseDriver is set to primary database
	if dbDriver == "" && len(dbDrivers) > 0 {
		context["DatabaseDriver"] = dbDrivers[0]
	}

	// Add convenience flags for each database type
	allDrivers := dbDrivers
	if dbDriver != "" {
		// Add legacy single driver to the list if not already present
		found := false
		for _, d := range dbDrivers {
			if d == dbDriver {
				found = true
				break
			}
		}
		if !found {
			allDrivers = append([]string{dbDriver}, dbDrivers...)
		}
	}

	for _, driver := range allDrivers {
		switch driver {
		case "postgresql", "postgres":
			context["HasPostgreSQL"] = true
		case "mysql":
			context["HasMySQL"] = true
		case "mongodb", "mongo":
			context["HasMongoDB"] = true
		case "sqlite":
			context["HasSQLite"] = true
		case "redis":
			context["HasRedis"] = true
		}
	}

	// Add secondary database flags (for multi-database setups)
	if len(dbDrivers) > 1 {
		context["HasMultipleDatabases"] = true
		for _, driver := range dbDrivers[1:] { // Skip primary database
			switch driver {
			case "redis":
				context["HasRedisCache"] = true
			case "mongodb", "mongo":
				context["HasMongoAnalytics"] = true
			}
		}
	}

	authType := g.getFeatureValue(config, "authentication", "type", "")
	if authType == "" && config.Variables != nil {
		if val, exists := config.Variables["AuthType"]; exists {
			authType = val
		}
	}
	context["AuthType"] = authType

	ormValue := g.getFeatureValue(config, "database", "orm", "gorm")

	// Validate ORM implementation
	if err := g.validateORM(ormValue); err != nil {
		// Log warning but don't fail - fallback to default
		ormValue = "gorm"
	}

	context["ORM"] = ormValue
	context["DatabaseORM"] = ormValue

	return context
}

// getFeatureValue safely gets a nested feature value
func (g *Generator) getFeatureValue(config types.ProjectConfig, feature, key, defaultValue string) string {
	if config.Features == nil {
		return defaultValue
	}

	switch feature {
	case "database":
		switch key {
		case "driver":
			// For backward compatibility, return primary driver
			if len(config.Features.Database.Drivers) > 0 {
				return config.Features.Database.Drivers[0]
			}
			return config.Features.Database.Driver //nolint:staticcheck // kept for backward compatibility
		case "orm":
			if config.Features.Database.ORM != "" {
				return config.Features.Database.ORM
			}
			return defaultValue
		}
	case "authentication":
		switch key {
		case "type":
			if config.Features.Authentication.Type != "" {
				return config.Features.Authentication.Type
			}
			return defaultValue
		}
	}

	return defaultValue
}

// validateORM checks if the specified ORM is currently implemented
func (g *Generator) validateORM(orm string) error {
	supportedORMs := map[string]bool{
		"gorm": true,
		"raw":  true,
		"":     true, // empty/default is valid
	}

	if !supportedORMs[orm] {
		return fmt.Errorf("ORM '%s' is not yet implemented. Currently supported: gorm, raw. See PROJECT_ROADMAP.md for implementation timeline", orm)
	}

	return nil
}

// getDatabaseDrivers returns all configured database drivers
func (g *Generator) getDatabaseDrivers(config types.ProjectConfig) []string {
	if config.Features == nil {
		return []string{}
	}

	return config.Features.Database.GetDrivers()
}

// processTemplateFile processes a single template file
func (g *Generator) processTemplateFile(templateDir, sourceFile, destPath string, context map[string]any) error {
	// Load template content
	content, err := g.loader.LoadTemplateFile(templateDir, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to load template file: %w", err)
	}

	// Parse template with Sprig functions
	tmpl, err := template.New(sourceFile).Funcs(sprig.FuncMap()).Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Write to destination
	if err := os.WriteFile(destPath, buf.Bytes(), 0644); err != nil {
		return types.NewFileSystemError("failed to write file", err)
	}

	return nil
}

// evaluateCondition evaluates a template condition
func (g *Generator) evaluateCondition(condition string, context map[string]any) (bool, error) {
	// Parse condition as a template
	tmpl, err := template.New("condition").Funcs(sprig.FuncMap()).Parse(condition)
	if err != nil {
		return false, fmt.Errorf("failed to parse condition: %w", err)
	}

	// Execute condition
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return false, fmt.Errorf("failed to execute condition: %w", err)
	}

	// Parse result
	result := strings.TrimSpace(buf.String())

	// Handle boolean values
	if result == "true" {
		return true, nil
	}
	if result == "false" {
		return false, nil
	}

	// Handle numeric values (0 = false, non-zero = true)
	if num, err := strconv.Atoi(result); err == nil {
		return num != 0, nil
	}

	// Handle string values (empty = false, non-empty = true)
	return result != "", nil
}

// processDependencies processes template dependencies
func (g *Generator) processDependencies(tmpl types.Template, _ types.ProjectConfig, outputPath string, context map[string]any) error {
	if len(tmpl.Dependencies) == 0 {
		return nil
	}

	var dependencies []string

	// Process each dependency
	for _, dep := range tmpl.Dependencies {
		// Check condition if present
		if dep.Condition != "" {
			shouldInclude, err := g.evaluateCondition(dep.Condition, context)
			if err != nil {
				return fmt.Errorf("failed to evaluate dependency condition: %w", err)
			}
			if !shouldInclude {
				continue
			}
		}

		// Add dependency
		if dep.Version != "" {
			dependencies = append(dependencies, fmt.Sprintf("%s@%s", dep.Module, dep.Version))
		} else {
			dependencies = append(dependencies, dep.Module)
		}
	}

	// If we have dependencies, add them to go.mod
	if len(dependencies) > 0 {
		return g.addDependencies(outputPath, dependencies)
	}

	return nil
}

// addDependencies adds dependencies to go.mod
func (g *Generator) addDependencies(projectPath string, dependencies []string) error {
	// Check if Go is available before trying to add dependencies
	if !g.isGoAvailable() {
		// Go is not available - generate a warning and instructions instead of failing
		g.logGoUnavailableWarning(dependencies)
		return nil
	}

	for _, dep := range dependencies {
		cmd := exec.Command("go", "get", dep)
		cmd.Dir = projectPath

		if output, err := cmd.CombinedOutput(); err != nil {
			// Clean up the output to make it more user-friendly
			outputStr := strings.TrimSpace(string(output))
			if outputStr == "" {
				return fmt.Errorf("failed to add dependency %q: %w", dep, err)
			}
			return fmt.Errorf("failed to add dependency %q: %s", dep, outputStr)
		}
	}
	return nil
}

// isGoAvailable checks if Go is installed and available
func (g *Generator) isGoAvailable() bool {
	cmd := exec.Command("go", "version")
	return cmd.Run() == nil
}

// logGoUnavailableWarning logs a warning when Go is not available
func (g *Generator) logGoUnavailableWarning(dependencies []string) {
	fmt.Println("⚠️  Warning: Go is not installed or not available in PATH")
	fmt.Println("   Project structure has been generated successfully, but dependencies were not installed.")
	fmt.Println()
	fmt.Println("   To complete the setup:")
	fmt.Println("   1. Install Go from https://golang.org/dl/")
	fmt.Println("   2. Navigate to your project directory")
	fmt.Println("   3. Run the following commands:")
	fmt.Println()
	for _, dep := range dependencies {
		fmt.Printf("      go get %s\n", dep)
	}
	fmt.Println("      go mod tidy")
	fmt.Println()
}

// executeHooks executes post-generation hooks
func (g *Generator) executeHooks(tmpl types.Template, config types.ProjectConfig, outputPath string, _ map[string]any) {
	for _, hook := range tmpl.PostHooks {
		// Check condition if present - for now, we'll execute all hooks
		// In the future, we can add conditional hook execution based on hook.Name

		// Determine working directory
		workDir := outputPath
		if hook.WorkDir != "" {
			// Process template variables in work directory
			workDir = g.processTemplatePath(hook.WorkDir, config)
			if !filepath.IsAbs(workDir) {
				workDir = filepath.Join(outputPath, workDir)
			}
		}

		// Execute command
		var cmd *exec.Cmd
		if len(hook.Args) > 0 {
			cmd = exec.Command(hook.Command, hook.Args...)
		} else {
			// Split command string if no explicit args
			parts := strings.Fields(hook.Command)
			if len(parts) == 0 {
				continue
			}
			cmd = exec.Command(parts[0], parts[1:]...)
		}

		cmd.Dir = workDir
		if output, err := cmd.CombinedOutput(); err != nil {
			// Don't fail the generation for hook errors, just warn
			fmt.Printf("Warning: Hook '%s' failed: %s\n", hook.Name, string(output))
		}
	}
}

// isGitAvailable checks if git is available in the system PATH
func (g *Generator) isGitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// hasGitRepository checks if a directory already contains a git repository
func (g *Generator) hasGitRepository(projectPath string) bool {
	gitDir := filepath.Join(projectPath, ".git")
	_, err := os.Stat(gitDir)
	return err == nil
}

// runGitCommand executes a git command in the specified directory
func (g *Generator) runGitCommand(projectPath string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git command failed: %s, output: %s", err, string(output))
	}

	return nil
}

// createGitignore creates a basic .gitignore file for Go projects
func (g *Generator) createGitignore(projectPath string) error {
	gitignoreContent := `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with "go test -c"
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# Environment variables
.env
.env.local

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Build artifacts
dist/
build/
bin/

# Logs
*.log
logs/

# Temporary files
tmp/
temp/
`

	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(strings.TrimSpace(gitignoreContent)), 0644); err != nil {
		return err
	}

	return nil
}
