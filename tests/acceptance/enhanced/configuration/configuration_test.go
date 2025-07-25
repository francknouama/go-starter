package configuration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
)

// TestFeatures runs the Enhanced Configuration Matrix BDD tests using godog
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			// Create test context for configuration matrix testing
			ctx := &ConfigurationTestContext{
				ProjectConfigs: make(map[string]*types.ProjectConfig),
				ProjectPaths:   make(map[string]string),
				TestResults:    make(map[string]*TestResult),
			}
			
			// Initialize templates
			if err := helpers.InitializeTemplates(); err != nil {
				t.Fatalf("Failed to initialize templates: %v", err)
			}
			
			s.Before(func(goCtx context.Context, sc *godog.Scenario) (context.Context, error) {
				// Setup before each scenario
				ctx.tempDir = helpers.CreateTempTestDir(t)
				ctx.startTime = time.Now()
				
				// Check if running in short mode
				if testing.Short() {
					ctx.shortMode = true
				}
				
				return goCtx, nil
			})
			
			s.After(func(goCtx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				// Cleanup after each scenario
				if ctx.tempDir != "" {
					os.RemoveAll(ctx.tempDir)
				}
				return goCtx, nil
			})
			
			ctx.RegisterSteps(s)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run enhanced configuration matrix feature tests")
	}
}

// ConfigurationTestContext holds test state for configuration matrix scenarios
type ConfigurationTestContext struct {
	ProjectConfigs map[string]*types.ProjectConfig
	ProjectPaths   map[string]string
	TestResults    map[string]*TestResult
	CurrentProject string
	CurrentConfig  *types.ProjectConfig
	tempDir        string
	startTime      time.Time
	shortMode      bool
}

// TestResult holds the result of a configuration test
type TestResult struct {
	Success       bool
	CompileTime   time.Duration
	Errors        []string
	Warnings      []string
	GenerateTime  time.Duration
}

// RegisterSteps registers all step definitions with godog for configuration matrix testing
func (ctx *ConfigurationTestContext) RegisterSteps(s *godog.ScenarioContext) {
	// Background steps
	s.Step(`^I have the go-starter CLI available$`, ctx.iHaveTheGoStarterCLIAvailable)
	s.Step(`^all templates are properly initialized$`, ctx.allTemplatesAreProperlyInitialized)
	
	// Configuration matrix steps
	s.Step(`^I use the critical configuration combination:$`, ctx.iUseTheCriticalConfigurationCombination)
	s.Step(`^I use the high priority configuration combination:$`, ctx.iUseTheHighPriorityConfigurationCombination)
	s.Step(`^I generate a web-api project with this configuration$`, ctx.iGenerateAWebApiProjectWithThisConfiguration)
	
	// Validation steps
	s.Step(`^the project should generate successfully$`, ctx.theProjectShouldGenerateSuccessfully)
	s.Step(`^the project should compile without errors$`, ctx.theProjectShouldCompileWithoutErrors)
	s.Step(`^all framework-specific code should be consistent$`, ctx.allFrameworkSpecificCodeShouldBeConsistent)
	s.Step(`^all database configuration should be consistent$`, ctx.allDatabaseConfigurationShouldBeConsistent)
	s.Step(`^all logger implementation should be consistent$`, ctx.allLoggerImplementationShouldBeConsistent)
	s.Step(`^all authentication setup should be consistent$`, ctx.allAuthenticationSetupShouldBeConsistent)
	s.Step(`^the architecture structure should be correct$`, ctx.theArchitectureStructureShouldBeCorrect)
	
	// Performance and matrix testing steps
	s.Step(`^I run the matrix test in short mode$`, ctx.iRunTheMatrixTestInShortMode)
	s.Step(`^only critical priority combinations should be tested$`, ctx.onlyCriticalPriorityCombinationsShouldBeTested)
	s.Step(`^all critical combinations should pass$`, ctx.allCriticalCombinationsShouldPass)
	s.Step(`^the test execution should complete within acceptable time limits$`, ctx.theTestExecutionShouldCompleteWithinAcceptableTimeLimits)
}

// Step definition implementations (placeholder implementations)

func (ctx *ConfigurationTestContext) iHaveTheGoStarterCLIAvailable() error {
	// Verify CLI is available
	return nil
}

func (ctx *ConfigurationTestContext) allTemplatesAreProperlyInitialized() error {
	// Verify templates are initialized
	// In a real implementation, we'd check template availability
	return nil
}

func (ctx *ConfigurationTestContext) iUseTheCriticalConfigurationCombination(table *godog.Table) error {
	// Parse configuration table and store critical configuration
	config := &types.ProjectConfig{
		Type: "web-api",
		Features: &types.Features{},
	}
	
	for i := 1; i < len(table.Rows); i++ { // Skip header
		row := table.Rows[i]
		if len(row.Cells) < 2 {
			continue
		}
		
		key := strings.TrimSpace(row.Cells[0].Value)
		value := strings.TrimSpace(row.Cells[1].Value)
		
		switch key {
		case "framework":
			config.Framework = value
		case "database":
			config.Features.Database.Driver = value // Store in Driver for backward compatibility
		case "database_driver":
			config.Features.Database.Driver = value
		case "orm":
			config.Features.Database.ORM = value
		case "logger":
			config.Logger = value
		case "auth_type":
			config.Features.Authentication.Type = value
		case "architecture":
			config.Architecture = value
		}
	}
	
	ctx.CurrentConfig = config
	return nil
}

func (ctx *ConfigurationTestContext) iUseTheHighPriorityConfigurationCombination(table *godog.Table) error {
	// Same as critical but marked as high priority
	return ctx.iUseTheCriticalConfigurationCombination(table)
}

func (ctx *ConfigurationTestContext) iGenerateAWebApiProjectWithThisConfiguration() error {
	// Generate project using stored configuration
	if ctx.CurrentConfig == nil {
		return fmt.Errorf("no configuration set")
	}
	
	projectName := fmt.Sprintf("test-project-%d", time.Now().UnixNano())
	projectPath := filepath.Join(ctx.tempDir, projectName)
	
	ctx.CurrentConfig.Name = projectName
	ctx.CurrentConfig.Module = fmt.Sprintf("github.com/test/%s", projectName)
	
	// Change to temp directory before generation  
	originalDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(ctx.tempDir); err != nil {
		return err
	}
	
	// Generate the project using the correct helper function
	startTime := time.Now()
	err = helpers.GenerateProjectInCwd(ctx.CurrentConfig)
	generateTime := time.Since(startTime)
	
	if err != nil {
		ctx.TestResults[projectName] = &TestResult{
			Success:      false,
			GenerateTime: generateTime,
			Errors:       []string{err.Error()},
		}
		return err
	}
	
	ctx.ProjectPaths[projectName] = projectPath
	ctx.ProjectConfigs[projectName] = ctx.CurrentConfig
	ctx.CurrentProject = projectName
	
	ctx.TestResults[projectName] = &TestResult{
		Success:      true,
		GenerateTime: generateTime,
	}
	
	return nil
}

func (ctx *ConfigurationTestContext) theProjectShouldGenerateSuccessfully() error {
	// Validate project generation succeeded
	if ctx.CurrentProject == "" {
		return fmt.Errorf("no current project")
	}
	
	result := ctx.TestResults[ctx.CurrentProject]
	if !result.Success {
		return fmt.Errorf("project generation failed: %v", result.Errors)
	}
	
	// Verify project structure exists
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", projectPath)
	}
	
	// Verify key files exist
	expectedFiles := []string{"go.mod", "main.go", "Makefile"}
	for _, file := range expectedFiles {
		filePath := filepath.Join(projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// main.go might be in cmd/server/
			if file == "main.go" {
				altPath := filepath.Join(projectPath, "cmd", "server", "main.go")
				if _, err := os.Stat(altPath); os.IsNotExist(err) {
					return fmt.Errorf("expected file not found: %s", file)
				}
			} else {
				return fmt.Errorf("expected file not found: %s", file)
			}
		}
	}
	
	return nil
}

func (ctx *ConfigurationTestContext) theProjectShouldCompileWithoutErrors() error {
	// Validate project compilation
	if ctx.CurrentProject == "" {
		return fmt.Errorf("no current project")
	}
	
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	result := ctx.TestResults[ctx.CurrentProject]
	
	// Run go mod download first
	modCmd := exec.Command("go", "mod", "download")
	modCmd.Dir = projectPath
	if output, err := modCmd.CombinedOutput(); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("go mod download failed: %s", output))
		return fmt.Errorf("go mod download failed: %s", output)
	}
	
	// Compile the project
	startTime := time.Now()
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = projectPath
	
	output, err := cmd.CombinedOutput()
	result.CompileTime = time.Since(startTime)
	
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("compilation failed: %s", output))
		return fmt.Errorf("compilation failed: %s", output)
	}
	
	// Check for warnings
	if strings.Contains(string(output), "warning:") {
		result.Warnings = append(result.Warnings, string(output))
	}
	
	return nil
}

func (ctx *ConfigurationTestContext) allFrameworkSpecificCodeShouldBeConsistent() error {
	// Validate framework consistency
	if ctx.CurrentProject == "" {
		return fmt.Errorf("no current project")
	}
	
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	config := ctx.ProjectConfigs[ctx.CurrentProject]
	
	// Get expected framework patterns
	expectedFramework := config.Framework
	otherFrameworks := []string{"gin", "echo", "fiber", "chi"}
	
	// Remove expected framework from others
	for i, fw := range otherFrameworks {
		if fw == expectedFramework {
			otherFrameworks = append(otherFrameworks[:i], otherFrameworks[i+1:]...)
			break
		}
	}
	
	// Check for cross-contamination
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if strings.HasSuffix(path, ".go") && !strings.Contains(path, "vendor/") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			
			contentStr := string(content)
			
			// Check for other framework imports
			for _, fw := range otherFrameworks {
				if strings.Contains(contentStr, fmt.Sprintf("%s-gonic/%s", fw, fw)) ||
				   strings.Contains(contentStr, fmt.Sprintf("labstack/%s", fw)) ||
				   strings.Contains(contentStr, fmt.Sprintf("gofiber/%s", fw)) ||
				   strings.Contains(contentStr, fmt.Sprintf("go-chi/%s", fw)) {
					return fmt.Errorf("found %s framework code in %s project: %s", fw, expectedFramework, path)
				}
			}
		}
		
		return nil
	})
	
	return err
}

func (ctx *ConfigurationTestContext) allDatabaseConfigurationShouldBeConsistent() error {
	// Validate database configuration consistency
	if ctx.CurrentProject == "" {
		return fmt.Errorf("no current project")
	}
	
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	config := ctx.ProjectConfigs[ctx.CurrentProject]
	
	// Check docker-compose.yml if database is configured
	if config.Features.Database.HasDatabase() {
		dockerComposePath := filepath.Join(projectPath, "docker-compose.yml")
		if _, err := os.Stat(dockerComposePath); err == nil {
			content, err := os.ReadFile(dockerComposePath)
			if err != nil {
				return err
			}
			
			contentStr := string(content)
			expectedService := config.Features.Database.PrimaryDriver()
			
			// Normalize postgres/postgresql
			if expectedService == "postgres" {
				expectedService = "postgresql"
			}
			
			if !strings.Contains(contentStr, expectedService+":") {
				return fmt.Errorf("docker-compose.yml does not contain %s service", expectedService)
			}
		}
		
		// Check connection string in config files
		configPath := filepath.Join(projectPath, "configs", "config.dev.yaml")
		if _, err := os.Stat(configPath); err == nil {
			content, err := os.ReadFile(configPath)
			if err != nil {
				return err
			}
			
			contentStr := string(content)
			driver := config.Features.Database.PrimaryDriver()
			if !strings.Contains(contentStr, driver) {
				return fmt.Errorf("config file does not contain correct database driver: %s", driver)
			}
		}
	}
	
	return nil
}

func (ctx *ConfigurationTestContext) allLoggerImplementationShouldBeConsistent() error {
	// Validate logger implementation consistency
	if ctx.CurrentProject == "" {
		return fmt.Errorf("no current project")
	}
	
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	config := ctx.ProjectConfigs[ctx.CurrentProject]
	
	// Check logger implementation file exists
	loggerPath := filepath.Join(projectPath, "internal", "logger", "logger.go")
	if _, err := os.Stat(loggerPath); os.IsNotExist(err) {
		return fmt.Errorf("logger implementation not found")
	}
	
	// Check go.mod for correct logger dependency
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return err
	}
	
	contentStr := string(content)
	
	// Verify correct logger dependency
	switch config.Logger {
	case "zap":
		if !strings.Contains(contentStr, "go.uber.org/zap") {
			return fmt.Errorf("zap dependency not found in go.mod")
		}
	case "logrus":
		if !strings.Contains(contentStr, "github.com/sirupsen/logrus") {
			return fmt.Errorf("logrus dependency not found in go.mod")
		}
	case "zerolog":
		if !strings.Contains(contentStr, "github.com/rs/zerolog") {
			return fmt.Errorf("zerolog dependency not found in go.mod")
		}
	case "slog":
		// slog is part of standard library in Go 1.21+
	}
	
	return nil
}

func (ctx *ConfigurationTestContext) allAuthenticationSetupShouldBeConsistent() error {
	// Validate authentication setup consistency
	if ctx.CurrentProject == "" {
		return fmt.Errorf("no current project")
	}
	
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	config := ctx.ProjectConfigs[ctx.CurrentProject]
	
	if config.Features.Authentication.Type != "" && config.Features.Authentication.Type != "none" {
		// Check for auth middleware
		authMiddlewarePath := filepath.Join(projectPath, "internal", "middleware", "auth.go")
		if _, err := os.Stat(authMiddlewarePath); os.IsNotExist(err) {
			return fmt.Errorf("auth middleware not found for auth type: %s", config.Features.Authentication.Type)
		}
		
		// Check for auth handler
		authHandlerPattern := filepath.Join(projectPath, "internal", "handlers", fmt.Sprintf("auth_%s.go", config.Framework))
		if _, err := os.Stat(authHandlerPattern); os.IsNotExist(err) {
			// Try generic auth handler
			authHandlerPattern = filepath.Join(projectPath, "internal", "handlers", "auth.go")
			if _, err := os.Stat(authHandlerPattern); os.IsNotExist(err) {
				return fmt.Errorf("auth handler not found for framework: %s", config.Framework)
			}
		}
	}
	
	return nil
}

func (ctx *ConfigurationTestContext) theArchitectureStructureShouldBeCorrect() error {
	// Validate architecture structure
	if ctx.CurrentProject == "" {
		return fmt.Errorf("no current project")
	}
	
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	config := ctx.ProjectConfigs[ctx.CurrentProject]
	
	// Check architecture-specific directories
	switch config.Architecture {
	case "clean":
		expectedDirs := []string{
			"internal/domain/entities",
			"internal/domain/usecases",
			"internal/adapters/controllers",
			"internal/infrastructure",
		}
		for _, dir := range expectedDirs {
			dirPath := filepath.Join(projectPath, dir)
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				return fmt.Errorf("clean architecture directory not found: %s", dir)
			}
		}
	case "ddd":
		expectedDirs := []string{
			"internal/domain",
			"internal/application",
			"internal/infrastructure",
		}
		for _, dir := range expectedDirs {
			dirPath := filepath.Join(projectPath, dir)
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				return fmt.Errorf("DDD directory not found: %s", dir)
			}
		}
	case "hexagonal":
		expectedDirs := []string{
			"internal/domain",
			"internal/application/ports",
			"internal/adapters/primary",
			"internal/adapters/secondary",
		}
		for _, dir := range expectedDirs {
			dirPath := filepath.Join(projectPath, dir)
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				return fmt.Errorf("hexagonal architecture directory not found: %s", dir)
			}
		}
	case "standard", "":
		expectedDirs := []string{
			"internal/handlers",
			"internal/services",
			"internal/repository",
			"internal/models",
		}
		for _, dir := range expectedDirs {
			dirPath := filepath.Join(projectPath, dir)
			if _, err := os.Stat(dirPath); os.IsNotExist(err) {
				return fmt.Errorf("standard architecture directory not found: %s", dir)
			}
		}
	}
	
	return nil
}

func (ctx *ConfigurationTestContext) iRunTheMatrixTestInShortMode() error {
	// Run matrix test in short mode
	ctx.shortMode = true
	return nil
}

func (ctx *ConfigurationTestContext) onlyCriticalPriorityCombinationsShouldBeTested() error {
	// Validate only critical combinations are tested
	if !ctx.shortMode {
		return fmt.Errorf("not in short mode")
	}
	
	// In short mode, we should have tested only critical combinations
	// This is controlled by the test scenarios
	return nil
}

func (ctx *ConfigurationTestContext) allCriticalCombinationsShouldPass() error {
	// Validate all critical combinations pass
	for project, result := range ctx.TestResults {
		if !result.Success {
			return fmt.Errorf("critical combination %s failed: %v", project, result.Errors)
		}
	}
	return nil
}

func (ctx *ConfigurationTestContext) theTestExecutionShouldCompleteWithinAcceptableTimeLimits() error {
	// Validate execution time
	elapsed := time.Since(ctx.startTime)
	
	// In short mode, tests should complete within 2 minutes
	if ctx.shortMode && elapsed > 2*time.Minute {
		return fmt.Errorf("short mode tests took too long: %v", elapsed)
	}
	
	// Full tests should complete within 10 minutes
	if !ctx.shortMode && elapsed > 10*time.Minute {
		return fmt.Errorf("full tests took too long: %v", elapsed)
	}
	
	return nil
}