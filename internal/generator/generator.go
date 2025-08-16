package generator

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/francknouama/go-starter/internal/config"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
)

// GenerationTransaction tracks operations that can be rolled back
type GenerationTransaction struct {
	outputPath    string
	filesCreated  []string
	dirsCreated   []string
	hooksExecuted []string
}

// NewGenerationTransaction creates a new transaction for rollback support
func NewGenerationTransaction(outputPath string) *GenerationTransaction {
	return &GenerationTransaction{
		outputPath:    outputPath,
		filesCreated:  make([]string, 0),
		dirsCreated:   make([]string, 0),
		hooksExecuted: make([]string, 0),
	}
}

// AddFile tracks a created file for potential rollback
func (tx *GenerationTransaction) AddFile(path string) {
	tx.filesCreated = append(tx.filesCreated, path)
}

// AddDirectory tracks a created directory for potential rollback
func (tx *GenerationTransaction) AddDirectory(path string) {
	tx.dirsCreated = append(tx.dirsCreated, path)
}

// AddHook tracks an executed hook for logging
func (tx *GenerationTransaction) AddHook(hookName string) {
	tx.hooksExecuted = append(tx.hooksExecuted, hookName)
}

// Rollback removes all created files and directories
func (tx *GenerationTransaction) Rollback() error {
	if len(tx.filesCreated) == 0 && len(tx.dirsCreated) == 0 {
		return nil
	}

	// Note: Use a fallback logger since g.logger might not be available during rollback
	fmt.Printf("Rolling back failed generation at %s...\n", tx.outputPath)

	var rollbackErrors []string

	// Remove created files in reverse order
	for i := len(tx.filesCreated) - 1; i >= 0; i-- {
		file := tx.filesCreated[i]
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			rollbackErrors = append(rollbackErrors, fmt.Sprintf("file %s: %v", file, err))
		}
	}

	// Remove directories in reverse order (only if empty)
	for i := len(tx.dirsCreated) - 1; i >= 0; i-- {
		dir := tx.dirsCreated[i]
		if err := os.Remove(dir); err != nil && !os.IsNotExist(err) {
			// Don't report errors for non-empty directories - that's expected
			if !strings.Contains(err.Error(), "directory not empty") {
				rollbackErrors = append(rollbackErrors, fmt.Sprintf("dir %s: %v", dir, err))
			}
		}
	}

	if len(rollbackErrors) > 0 {
		fmt.Printf("Rollback completed with %d errors\n", len(rollbackErrors))
		return fmt.Errorf("rollback errors: %s", strings.Join(rollbackErrors, "; "))
	}

	fmt.Println("Rollback completed successfully")
	return nil
}

// Generator handles project generation
type Generator struct {
	registry           *templates.Registry
	loader             *templates.TemplateLoader
	currentTransaction *GenerationTransaction
	logger             *GeneratorLogger
	errorHandler       *ErrorHandler
}

// New creates a new Generator instance with default logging
func New() *Generator {
	logger := NewGeneratorLogger(LogLevelInfo)
	return &Generator{
		registry:     templates.NewRegistry(),
		loader:       templates.NewTemplateLoader(),
		logger:       logger,
		errorHandler: NewErrorHandler(logger),
	}
}

// NewWithLogger creates a new Generator instance with custom logging
func NewWithLogger(logLevel LogLevel) *Generator {
	logger := NewGeneratorLogger(logLevel)
	return &Generator{
		registry:     templates.NewRegistry(),
		loader:       templates.NewTemplateLoader(),
		logger:       logger,
		errorHandler: NewErrorHandler(logger),
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

	// Create transaction for rollback support
	tx := NewGenerationTransaction(options.OutputPath)

	// Set up recovery mechanism
	defer func() {
		if r := recover(); r != nil {
			g.logger.Error("Generation panic occurred: %v", r)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				fmt.Printf("Rollback failed: %v\n", rollbackErr)
			}
			panic(r) // Re-panic after cleanup
		}
	}()

	// Validate configuration
	if err := g.validateConfig(config); err != nil {
		result.Error = err
		return result, err
	}

	// Check if template exists
	templateID := g.getTemplateID(config)
	template, err := g.registry.Get(templateID)
	if err != nil {
		// Try fallback: look for templates by type if direct ID lookup fails
		templatesByType := g.registry.GetByType(config.Type)
		if len(templatesByType) > 0 {
			// Use the first template of this type as fallback
			template = templatesByType[0]
			g.logger.Debug("Using template '%s' for type '%s'", template.ID, config.Type)
		} else {
			// For Phase 0, we'll show a helpful message about upcoming templates
			return g.handleMissingTemplate(config, result)
		}
	}

	// Skip file system operations in dry run mode
	if options.DryRun {
		return g.performDryRun(config, &template, options, result, startTime)
	}

	// Check if output directory already exists and validate it
	if err := g.checkOutputDirectory(options.OutputPath); err != nil {
		result.Error = err
		return result, err
	}

	// Create output directory
	if err := os.MkdirAll(options.OutputPath, 0755); err != nil {
		result.Error = types.NewFileSystemError("failed to create output directory", err)
		return result, result.Error
	}
	tx.AddDirectory(options.OutputPath)

	// Generate project files with transaction tracking
	filesCreated, err := g.generateProjectFilesWithTransaction(template, config, options.OutputPath, tx)
	if err != nil {
		fmt.Printf("File generation error: %v\n", err)
		result.Error = err
		// Perform rollback on failure
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			fmt.Printf("Rollback failed: %v\n", rollbackErr)
		}
		return result, err
	}
	result.FilesCreated = filesCreated

	// Initialize git repository if requested
	if !options.NoGit {
		fmt.Printf("🔧 Initializing git repository...\n")
		if err := g.initGitRepository(options.OutputPath); err != nil {
			// Git init failure is not fatal, just log it
			fmt.Printf("⚠️  Warning: failed to initialize git repository: %v\n", err)
		} else {
			fmt.Printf("✓ Git repository initialized\n")
		}
	}

	result.Duration = time.Since(startTime)
	result.Success = true

	// Show completion summary
	g.logger.Success("Project '%s' generated successfully!", config.Name)
	g.logger.Info("Project location: %s", options.OutputPath)
	g.logger.Info("Generation completed in %v", result.Duration)

	return result, nil
}

// GenerateInMemory generates a project in memory and returns the file contents
func (g *Generator) GenerateInMemory(config *types.ProjectConfig, blueprintID string) (map[string][]byte, error) {
	// Validate configuration
	if err := g.validateConfig(*config); err != nil {
		return nil, err
	}

	// Get template
	tmpl, err := g.registry.Get(blueprintID)
	if err != nil {
		return nil, fmt.Errorf("template not found: %s", blueprintID)
	}

	// Generate files in memory
	files := make(map[string][]byte)
	context := g.createTemplateContext(*config, tmpl)

	for _, file := range tmpl.Files {
		// Skip files with failing conditions
		if file.Condition != "" {
			shouldInclude, err := g.evaluateCondition(file.Condition, context)
			if err != nil {
				fmt.Printf("Warning: Failed to evaluate condition %q: %v\n", file.Condition, err)
				continue
			}
			if !shouldInclude {
				continue
			}
		}

		// Process destination path
		destPath := g.processTemplatePath(file.Destination, *config, &tmpl)

		// Load and process template content
		content, err := g.loader.LoadTemplateFile(blueprintID, file.Source)
		if err != nil {
			return nil, fmt.Errorf("failed to load template %s: %w", file.Source, err)
		}

		// Execute template
		goTmpl, err := template.New(file.Source).Funcs(sprig.TxtFuncMap()).Parse(content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", file.Source, err)
		}

		var buf bytes.Buffer
		if err := goTmpl.Execute(&buf, context); err != nil {
			return nil, fmt.Errorf("failed to execute template %s: %w", file.Source, err)
		}

		files[destPath] = buf.Bytes()
	}

	return files, nil
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
		// Try fallback: look for templates by type if direct ID lookup fails
		templatesByType := g.registry.GetByType(config.Type)
		if len(templatesByType) > 0 {
			// Use the first template of this type as fallback
			template = templatesByType[0]
			g.logger.Debug("Using template '%s' for type '%s'", template.ID, config.Type)
		} else {
			fmt.Printf("\nTemplate '%s' not found.\n", templateID)
			fmt.Printf("Use 'go-starter list' to see available blueprints.\n")
			return nil
		}
	}

	fmt.Printf("\nFiles to be generated:\n")
	for _, file := range template.Files {
		destination := g.processTemplatePath(file.Destination, config, &template)
		fmt.Printf("  %s\n", destination)
	}

	return nil
}

// validateConfig validates the project configuration
func (g *Generator) validateConfig(cfg types.ProjectConfig) error {
	// Validate project name
	if err := config.ValidateProjectName(cfg.Name); err != nil {
		return types.NewValidationError(fmt.Sprintf("invalid project name: %v", err), nil)
	}

	// Validate module path
	if err := config.ValidateModulePath(cfg.Module); err != nil {
		return types.NewValidationError(fmt.Sprintf("invalid module path: %v", err), nil)
	}

	// Validate project type
	if err := config.ValidateTemplateType(cfg.Type); err != nil {
		return types.NewValidationError(fmt.Sprintf("invalid project type: %v", err), nil)
	}

	// Validate Go version
	if cfg.GoVersion != "" {
		if err := config.ValidateGoVersion(cfg.GoVersion); err != nil {
			return types.NewValidationError(fmt.Sprintf("invalid Go version: %v", err), nil)
		}
	}

	// Validate framework
	if err := config.ValidateFramework(cfg.Framework); err != nil {
		return types.NewValidationError(fmt.Sprintf("invalid framework: %v", err), nil)
	}

	// Validate architecture
	if err := config.ValidateArchitecture(cfg.Architecture); err != nil {
		return types.NewValidationError(fmt.Sprintf("invalid architecture: %v", err), nil)
	}

	// Validate logger
	if err := config.ValidateLogger(cfg.Logger); err != nil {
		return types.NewValidationError(fmt.Sprintf("invalid logger: %v", err), nil)
	}

	// Validate features (check for nil features first)
	if cfg.Features != nil {
		if cfg.Features.Database.Driver != "" {
			if err := config.ValidateDatabaseDriver(cfg.Features.Database.Driver); err != nil {
				return types.NewValidationError(fmt.Sprintf("invalid database driver: %v", err), nil)
			}
		}

		if cfg.Features.Database.ORM != "" {
			if err := config.ValidateORM(cfg.Features.Database.ORM); err != nil {
				return types.NewValidationError(fmt.Sprintf("invalid ORM: %v", err), nil)
			}
		}

		if cfg.Features.Authentication.Type != "" {
			if err := config.ValidateAuthType(cfg.Features.Authentication.Type); err != nil {
				return types.NewValidationError(fmt.Sprintf("invalid authentication type: %v", err), nil)
			}
		}
	}

	return nil
}

// getTemplateID maps project configuration to template ID
func (g *Generator) getTemplateID(config types.ProjectConfig) string {
	// First check if a specific blueprint_id is set by the interactive CLI
	if config.Variables != nil {
		if blueprintID, exists := config.Variables["blueprint_id"]; exists && blueprintID != "" {
			return blueprintID
		}
	}
	
	// Fall back to architecture-based selection
	if config.Architecture != "" && config.Architecture != "standard" {
		return fmt.Sprintf("%s-%s", config.Type, config.Architecture)
	}
	return config.Type
}

// handleMissingTemplate provides helpful feedback when templates aren't available yet
func (g *Generator) handleMissingTemplate(config types.ProjectConfig, result *types.GenerationResult) (*types.GenerationResult, error) {
	templateID := g.getTemplateID(config)

	fmt.Printf("Template '%s' not found.\n", templateID)
	fmt.Println("\nAvailable templates:")
	
	// Get list of available templates from registry
	templates := g.registry.List()
	if len(templates) > 0 {
		for _, tmpl := range templates {
			fmt.Printf("  • %s - %s\n", tmpl.ID, tmpl.Description)
		}
		fmt.Printf("\nUse 'go-starter list' to see all available blueprints with detailed descriptions.\n")
	} else {
		fmt.Println("  No templates currently available.")
	}

	err := types.NewTemplateNotFoundError(templateID)
	result.Error = err
	return result, err
}

// generateProjectFiles generates all files for the project
// generateProjectFilesWithTransaction generates project files with rollback support
func (g *Generator) generateProjectFilesWithTransaction(tmpl types.Template, config types.ProjectConfig, outputPath string, tx *GenerationTransaction) ([]string, error) {
	// Set the transaction in generator for file tracking
	g.currentTransaction = tx
	defer func() { g.currentTransaction = nil }()

	// Use the existing generateProjectFiles function
	return g.generateProjectFiles(tmpl, config, outputPath)
}

func (g *Generator) generateProjectFiles(tmpl types.Template, config types.ProjectConfig, outputPath string) ([]string, error) {
	// Create template context with all variables
	context := g.createTemplateContext(config, tmpl)

	// Get template directory from metadata
	templateDir, ok := tmpl.Metadata["path"].(string)
	if !ok {
		return nil, fmt.Errorf("template metadata missing path")
	}

	g.logger.Progress("Generating project files...")

	// Filter files that should be generated based on conditions
	var filesToProcess []types.TemplateFile
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
		filesToProcess = append(filesToProcess, templateFile)
	}

	// Use parallel processing for better performance
	numWorkers := runtime.NumCPU()
	if len(filesToProcess) < numWorkers {
		// For small projects, limit workers to avoid overhead
		numWorkers = len(filesToProcess)
	}
	if numWorkers < 1 {
		numWorkers = 1
	}

	// Create parallel processor with transaction support
	processor := NewParallelTemplateProcessor(numWorkers, g.currentTransaction)
	
	// Process all template files in parallel
	filesCreated, err := processor.ProcessTemplates(
		filesToProcess,
		templateDir,
		outputPath,
		context,
		config,
		&tmpl,
	)
	if err != nil {
		return nil, err
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

	// Track file creation for rollback if transaction is active
	if g.currentTransaction != nil {
		g.currentTransaction.AddFile(path)
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
func (g *Generator) processTemplatePath(path string, config types.ProjectConfig, tmplObj *types.Template) string {
	// Create a comprehensive context for path processing (same as file content)
	context := map[string]any{
		"ProjectName":  config.Name,
		"ModulePath":   config.Module,
		"GoVersion":    config.GoVersion,
		"Framework":    config.Framework,
		"Architecture": config.Architecture,
		"Logger":       config.Logger,
		"Author":       config.Author,
		"Email":        config.Email,
		"License":      config.License,
	}

	// Add variables from config.Variables map
	if config.Variables != nil {
		for key, value := range config.Variables {
			context[key] = value
		}
	}

	// Add template-specific variables with defaults
	if tmplObj != nil {
		for _, variable := range tmplObj.Variables {
			// Skip if already set from config
			if _, exists := context[variable.Name]; exists {
				continue
			}

			// Add default value if specified
			if variable.Default != nil {
				context[variable.Name] = variable.Default
			}
		}
	}

	// Use cached template parsing for paths
	cacheKey := fmt.Sprintf("path:%s", path)
	tmpl, err := templates.GetOrParseTemplate(cacheKey, path, sprig.FuncMap())
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
	// Resolve Go version (handle empty, "auto" cases)
	goVersion := config.GoVersion
	if goVersion == "" || goVersion == "auto" {
		goVersion = "1.21" // Default to latest stable supported version
	}
	
	context := map[string]any{
		"ProjectName":  config.Name,
		"ModulePath":   config.Module,
		"GoVersion":    goVersion,
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
	// Check Features struct first, then fall back to Variables map, then template defaults
	dbDriver := g.getFeatureValue(config, "database", "driver", "")
	if dbDriver == "" && config.Variables != nil {
		if val, exists := config.Variables["DatabaseDriver"]; exists {
			dbDriver = val
		}
	}
	
	// Apply template default for DatabaseDriver if still empty
	if dbDriver == "" {
		for _, variable := range tmpl.Variables {
			if variable.Name == "DatabaseDriver" && variable.Default != nil {
				if defaultVal, ok := variable.Default.(string); ok {
					dbDriver = defaultVal
				}
			}
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
	
	// Apply template default for AuthType if still empty
	if authType == "" {
		for _, variable := range tmpl.Variables {
			if variable.Name == "AuthType" && variable.Default != nil {
				if defaultVal, ok := variable.Default.(string); ok {
					authType = defaultVal
				}
			}
		}
	}
	context["AuthType"] = authType

	ormValue := g.getFeatureValue(config, "database", "orm", "")

	// Check Variables map if not found in Features (like we do for DatabaseDriver)
	if ormValue == "" && config.Variables != nil {
		if val, exists := config.Variables["DatabaseORM"]; exists {
			ormValue = val
		}
	}
	
	// Apply template default for DatabaseORM if still empty
	if ormValue == "" {
		for _, variable := range tmpl.Variables {
			if variable.Name == "DatabaseORM" && variable.Default != nil {
				if defaultVal, ok := variable.Default.(string); ok {
					ormValue = defaultVal
				}
			}
		}
	}
	
	// Validate ORM implementation (only if not empty - empty means raw SQL)
	if ormValue != "" {
		if err := g.validateORM(ormValue); err != nil {
			// Log warning but don't fail - fallback to raw SQL (empty string)
			ormValue = ""
		}
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
	// Use cached template parsing
	cacheKey := fmt.Sprintf("%s:%s", templateDir, sourceFile)
	tmpl, err := templates.GetOrParseTemplate(cacheKey, content, sprig.FuncMap())
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

	// Track file creation for rollback if transaction is active
	if g.currentTransaction != nil {
		g.currentTransaction.AddFile(destPath)
	}

	return nil
}

// evaluateCondition evaluates a template condition
func (g *Generator) evaluateCondition(condition string, context map[string]any) (bool, error) {
	// Use cached template parsing for conditions
	cacheKey := fmt.Sprintf("condition:%s", condition)
	tmpl, err := templates.GetOrParseTemplate(cacheKey, condition, sprig.FuncMap())
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

	fmt.Printf("📦 Processing dependencies...\n")
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
		if err := g.addDependencies(outputPath, dependencies); err != nil {
			return err
		}
		fmt.Printf("✓ Added %d dependencies\n", len(dependencies))
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
			// Handle special case of {{.OutputPath}} in hook work directory
			if hook.WorkDir == "{{.OutputPath}}" {
				workDir = outputPath
			} else {
				// Process other template variables in work directory
				workDir = g.processTemplatePath(hook.WorkDir, config, &tmpl)
				if !filepath.IsAbs(workDir) {
					workDir = filepath.Join(outputPath, workDir)
				}
			}
		}

		// Execute command
		var cmd *exec.Cmd
		if len(hook.Args) > 0 {
			cmd = exec.Command(hook.Command, hook.Args...)
		} else {
			// Check if command contains shell metacharacters that need expansion
			if strings.Contains(hook.Command, "*") || strings.Contains(hook.Command, "?") || strings.Contains(hook.Command, "[") {
				// Use shell for wildcard expansion
				cmd = exec.Command("sh", "-c", hook.Command)
			} else {
				// Split command string if no explicit args
				parts := strings.Fields(hook.Command)
				if len(parts) == 0 {
					continue
				}
				cmd = exec.Command(parts[0], parts[1:]...)
			}
		}

		cmd.Dir = workDir
		if output, err := cmd.CombinedOutput(); err != nil {
			// Don't fail the generation for hook errors, just warn
			if len(output) > 0 {
				fmt.Printf("Warning: Hook '%s' failed with error: %v\nOutput: %s\n", hook.Name, err, string(output))
			} else {
				fmt.Printf("Warning: Hook '%s' failed with error: %v\n", hook.Name, err)
			}
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
		return fmt.Errorf("failed to create .gitignore file: %w", err)
	}

	// Track file creation for rollback if transaction is active
	if g.currentTransaction != nil {
		g.currentTransaction.AddFile(gitignorePath)
	}

	return nil
}

// checkOutputDirectory validates the output directory before generation
func (g *Generator) checkOutputDirectory(outputPath string) error {
	// Check if the path exists
	info, err := os.Stat(outputPath)
	if os.IsNotExist(err) {
		// Directory doesn't exist, this is fine
		return nil
	}
	if err != nil {
		return types.NewFileSystemError("failed to check output directory", err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return types.NewValidationError(fmt.Sprintf("output path '%s' exists but is not a directory", outputPath), nil)
	}

	// Check if directory is empty
	entries, err := os.ReadDir(outputPath)
	if err != nil {
		return types.NewFileSystemError("failed to read output directory", err)
	}

	if len(entries) > 0 {
		return types.NewValidationError(fmt.Sprintf("directory '%s' already exists and is not empty. Please choose a different name or remove the existing directory", filepath.Base(outputPath)), nil)
	}

	return nil
}

// performDryRun simulates project generation and shows what would be created
func (g *Generator) performDryRun(config types.ProjectConfig, template *types.Template, options types.GenerationOptions, result *types.GenerationResult, startTime time.Time) (*types.GenerationResult, error) {
	fmt.Printf("🔍 DRY RUN MODE - Preview of project generation\n\n")
	fmt.Printf("📋 Project Configuration:\n")
	fmt.Printf("   Name: %s\n", config.Name)
	fmt.Printf("   Module: %s\n", config.Module)
	fmt.Printf("   Type: %s\n", config.Type)
	fmt.Printf("   Architecture: %s\n", config.Architecture)
	fmt.Printf("   Framework: %s\n", config.Framework)
	fmt.Printf("   Logger: %s\n", config.Logger)
	fmt.Printf("   Go Version: %s\n\n", config.GoVersion)

	fmt.Printf("📁 Files that would be generated:\n")

	// Create template context
	context := g.createTemplateContext(config, *template)

	var fileList []string
	for _, templateFile := range template.Files {
		// Evaluate condition if present
		if templateFile.Condition != "" {
			shouldGenerate, err := g.evaluateCondition(templateFile.Condition, context)
			if err != nil {
				continue // Skip on error
			}
			if !shouldGenerate {
				continue
			}
		}

		// Process template path with variables
		destPath := g.processTemplatePath(templateFile.Destination, config, template)
		relPath := filepath.Join(config.Name, destPath)
		fmt.Printf("   ✓ %s\n", relPath)
		fileList = append(fileList, relPath)
	}

	// Show dependency information
	if len(template.Dependencies) > 0 {
		fmt.Printf("\n📦 Dependencies that would be added:\n")
		for _, dep := range template.Dependencies {
			// Check if dependency condition is met
			if dep.Condition != "" {
				shouldInclude, err := g.evaluateCondition(dep.Condition, context)
				if err != nil || !shouldInclude {
					continue
				}
			}
			fmt.Printf("   + %s\n", dep.Module)
		}
	}

	// Show git initialization info
	if !options.NoGit {
		fmt.Printf("\n🔧 Additional actions that would be performed:\n")
		fmt.Printf("   ✓ Initialize git repository\n")
		fmt.Printf("   ✓ Create .gitignore file\n")
		fmt.Printf("   ✓ Run go mod tidy\n")
	}

	fmt.Printf("\n✨ Total files: %d\n", len(fileList))
	fmt.Printf("💡 To actually generate the project, run the same command without --dry-run\n\n")

	result.Success = true
	result.FilesCreated = fileList
	result.Duration = time.Since(startTime)
	return result, nil
}
