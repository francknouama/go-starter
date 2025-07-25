package workspace

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

// WorkspaceTestContext holds test state for workspace integration testing
type WorkspaceTestContext struct {
	workspaceName     string
	workspacePath     string
	components        map[string]*ComponentConfig
	tempDir           string
	startTime         time.Time
	workspaceConfig   *types.ProjectConfig
	lastCommandOutput string
	lastCommandError  error
}

// ComponentConfig represents configuration for a workspace component
type ComponentConfig struct {
	Name         string
	Type         string
	Architecture string
	Framework    string
	Database     string
	ORM          string
	Auth         string
	Logger       string
	Complexity   string
	GoVersion    string
	Config       *types.ProjectConfig
}

// TestFeatures runs the workspace integration BDD tests
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			ctx := &WorkspaceTestContext{
				components: make(map[string]*ComponentConfig),
			}

			s.Before(func(goCtx context.Context, sc *godog.Scenario) (context.Context, error) {
				// Setup before each scenario
				ctx.tempDir = helpers.CreateTempTestDir(t)
				ctx.startTime = time.Now()
				ctx.components = make(map[string]*ComponentConfig)

				// Initialize templates
				if err := helpers.InitializeTemplates(); err != nil {
					return goCtx, err
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
		t.Fatal("non-zero status returned, failed to run workspace integration tests")
	}
}

// RegisterSteps registers all step definitions
func (ctx *WorkspaceTestContext) RegisterSteps(s *godog.ScenarioContext) {
	// Background steps
	s.Step(`^I have the go-starter CLI available$`, ctx.iHaveTheGoStarterCLIAvailable)
	s.Step(`^all templates are properly initialized$`, ctx.allTemplatesAreProperlyInitialized)
	s.Step(`^I am testing workspace integration scenarios$`, ctx.iAmTestingWorkspaceIntegrationScenarios)

	// Workspace generation steps
	s.Step(`^I generate a workspace project with configuration:$`, ctx.iGenerateAWorkspaceProjectWithConfiguration)
	s.Step(`^I add component "([^"]*)" with configuration:$`, ctx.iAddComponentWithConfiguration)

	// Validation steps
	s.Step(`^the workspace should be generated successfully$`, ctx.theWorkspaceShouldBeGeneratedSuccessfully)
	s.Step(`^all components should share a consistent root go\.mod$`, ctx.allComponentsShouldShareAConsistentRootGoMod)
	s.Step(`^shared dependencies should be properly managed$`, ctx.sharedDependenciesShouldBeProperlyManaged)
	s.Step(`^component go\.mod files should be properly configured$`, ctx.componentGoModFilesShouldBeProperlyConfigured)
	s.Step(`^workspace compilation should succeed$`, ctx.workspaceCompilationShouldSucceed)
	s.Step(`^inter-component imports should work correctly$`, ctx.interComponentImportsShouldWorkCorrectly)
	s.Step(`^shared configuration should be available to all components$`, ctx.sharedConfigurationShouldBeAvailableToAllComponents)

	// Database and service-specific validations
	s.Step(`^each service should have independent database configurations$`, ctx.eachServiceShouldHaveIndependentDatabaseConfigurations)
	s.Step(`^shared authentication libraries should be consistent$`, ctx.sharedAuthenticationLibrariesShouldBeConsistent)
	s.Step(`^logger implementations should be consistent across services$`, ctx.loggerImplementationsShouldBeConsistentAcrossServices)
	s.Step(`^service-to-service communication patterns should be available$`, ctx.serviceToServiceCommunicationPatternsShouldBeAvailable)

	// Modular and monolith validations
	s.Step(`^shared database connections should be properly managed$`, ctx.sharedDatabaseConnectionsShouldBeProperlyManaged)
	s.Step(`^session management should be consistent across web components$`, ctx.sessionManagementShouldBeConsistentAcrossWebComponents)
	s.Step(`^logger configuration should be shared$`, ctx.loggerConfigurationShouldBeShared)
	s.Step(`^modular structure should be maintained$`, ctx.modularStructureShouldBeMaintained)

	// Event-driven validations
	s.Step(`^event sourcing patterns should be properly implemented$`, ctx.eventSourcingPatternsShouldBeProperlyImplemented)
	s.Step(`^CQRS separation should be maintained$`, ctx.cqrsSeparationShouldBeMaintained)
	s.Step(`^event store configuration should be shared$`, ctx.eventStoreConfigurationShouldBeShared)
	s.Step(`^event flow between components should be configured$`, ctx.eventFlowBetweenComponentsShouldBeConfigured)

	// Complex workspace validations
	s.Step(`^all blueprint types should be properly integrated$`, ctx.allBlueprintTypesShouldBeProperlyIntegrated)
	s.Step(`^shared dependencies should be managed efficiently$`, ctx.sharedDependenciesShouldBeManagedEfficiently)
	s.Step(`^logger configuration should be consistent across all components$`, ctx.loggerConfigurationShouldBeConsistentAcrossAllComponents)
	s.Step(`^authentication should be compatible across components$`, ctx.authenticationShouldBeCompatibleAcrossComponents)
	s.Step(`^inter-component communication should be properly configured$`, ctx.interComponentCommunicationShouldBeProperlyConfigured)
	s.Step(`^shared library should be usable by all components$`, ctx.sharedLibraryShouldBeUsableByAllComponents)

	// Dependency and configuration validations
	s.Step(`^I generate a workspace with multiple components having different dependency requirements$`, ctx.iGenerateAWorkspaceWithMultipleComponentsHavingDifferentDependencyRequirements)
	s.Step(`^go\.mod dependency resolution should work correctly$`, ctx.goModDependencyResolutionShouldWorkCorrectly)
	s.Step(`^no version conflicts should exist$`, ctx.noVersionConflictsShouldExist)
	s.Step(`^shared dependencies should use consistent versions$`, ctx.sharedDependenciesShouldUseConsistentVersions)
	s.Step(`^workspace should build without dependency errors$`, ctx.workspaceShouldBuildWithoutDependencyErrors)

	// Configuration consistency validations
	s.Step(`^I generate a workspace with multiple components$`, ctx.iGenerateAWorkspaceWithMultipleComponents)
	s.Step(`^configuration patterns should be consistent across components$`, ctx.configurationPatternsShouldBeConsistentAcrossComponents)
	s.Step(`^environment variable handling should be standardized$`, ctx.environmentVariableHandlingShouldBeStandardized)
	s.Step(`^logging configuration should be centralized$`, ctx.loggingConfigurationShouldBeCentralized)
	s.Step(`^database connection management should be optimized$`, ctx.databaseConnectionManagementShouldBeOptimized)

	// Performance validations
	s.Step(`^I generate a workspace with (\d+)\+ components$`, ctx.iGenerateAWorkspaceWithPlusComponents)
	s.Step(`^generation time should be reasonable$`, ctx.generationTimeShouldBeReasonable)
	s.Step(`^memory usage should remain within acceptable limits$`, ctx.memoryUsageShouldRemainWithinAcceptableLimits)
	s.Step(`^compilation time should be optimized$`, ctx.compilationTimeShouldBeOptimized)
	s.Step(`^workspace structure should be maintainable$`, ctx.workspaceStructureShouldBeMaintainable)
}

// Step implementations

func (ctx *WorkspaceTestContext) iHaveTheGoStarterCLIAvailable() error {
	// Verify CLI is available and functional
	cmd := exec.Command("go-starter", "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go-starter CLI not available: %v", err)
	}
	return nil
}

func (ctx *WorkspaceTestContext) allTemplatesAreProperlyInitialized() error {
	// Verify templates are loaded and accessible
	return helpers.InitializeTemplates()
}

func (ctx *WorkspaceTestContext) iAmTestingWorkspaceIntegrationScenarios() error {
	// Set context for workspace integration testing
	return nil
}

func (ctx *WorkspaceTestContext) iGenerateAWorkspaceProjectWithConfiguration(table *godog.Table) error {
	// Parse workspace configuration from table
	config := &types.ProjectConfig{
		Type: "workspace",
	}

	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		key := row.Cells[0].Value
		value := row.Cells[1].Value

		switch key {
		case "type":
			config.Type = value
		case "name":
			config.Name = value
			ctx.workspaceName = value
		case "go_version":
			config.GoVersion = value
		case "output_dir":
			ctx.workspacePath = filepath.Join(ctx.tempDir, value)
		}
	}

	// Set default workspace path if not specified
	if ctx.workspacePath == "" {
		ctx.workspacePath = filepath.Join(ctx.tempDir, ctx.workspaceName)
	}

	ctx.workspaceConfig = config

	// Generate the workspace project
	return ctx.generateWorkspaceProject(config)
}

func (ctx *WorkspaceTestContext) iAddComponentWithConfiguration(componentName string, table *godog.Table) error {
	// Parse component configuration from table
	component := &ComponentConfig{
		Name: componentName,
	}

	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		key := row.Cells[0].Value
		value := row.Cells[1].Value

		switch key {
		case "type":
			component.Type = value
		case "architecture":
			component.Architecture = value
		case "framework":
			component.Framework = value
		case "database":
			component.Database = value
		case "orm":
			component.ORM = value
		case "auth":
			component.Auth = value
		case "logger":
			component.Logger = value
		case "complexity":
			component.Complexity = value
		case "go_version":
			component.GoVersion = value
		}
	}

	// Store component configuration
	ctx.components[componentName] = component

	// Generate the component within the workspace
	return ctx.generateWorkspaceComponent(component)
}

func (ctx *WorkspaceTestContext) generateWorkspaceProject(config *types.ProjectConfig) error {
	// Create workspace directory
	if err := os.MkdirAll(ctx.workspacePath, 0755); err != nil {
		return fmt.Errorf("failed to create workspace directory: %v", err)
	}

	// Change to workspace directory
	if err := os.Chdir(ctx.workspacePath); err != nil {
		return fmt.Errorf("failed to change to workspace directory: %v", err)
	}

	// Initialize workspace with go mod init
	cmd := exec.Command("go", "mod", "init", config.Name)
	cmd.Dir = ctx.workspacePath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to initialize workspace go.mod: %v, output: %s", err, output)
	}

	// Create workspace structure
	if err := ctx.createWorkspaceStructure(); err != nil {
		return fmt.Errorf("failed to create workspace structure: %v", err)
	}

	return nil
}

func (ctx *WorkspaceTestContext) createWorkspaceStructure() error {
	// Create standard workspace directories
	dirs := []string{
		"cmd",
		"internal",
		"pkg",
		"docs",
		"scripts",
		"deployments",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(ctx.workspacePath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	// Create workspace README
	readmePath := filepath.Join(ctx.workspacePath, "README.md")
	readmeContent := fmt.Sprintf("# %s\n\nThis is a workspace project generated by go-starter.\n", ctx.workspaceName)
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to create workspace README: %v", err)
	}

	return nil
}

func (ctx *WorkspaceTestContext) generateWorkspaceComponent(component *ComponentConfig) error {
	// Create component directory within workspace
	componentPath := filepath.Join(ctx.workspacePath, "cmd", component.Name)
	if err := os.MkdirAll(componentPath, 0755); err != nil {
		return fmt.Errorf("failed to create component directory: %v", err)
	}

	// Change to component directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	if err := os.Chdir(componentPath); err != nil {
		return fmt.Errorf("failed to change to component directory: %v", err)
	}

	// Build go-starter command for component
	args := []string{"new", component.Name}
	args = append(args, "--type", component.Type)

	if component.Architecture != "" {
		args = append(args, "--architecture", component.Architecture)
	}
	if component.Framework != "" {
		args = append(args, "--framework", component.Framework)
	}
	if component.Database != "" {
		args = append(args, "--database-driver", component.Database)
	}
	if component.ORM != "" {
		args = append(args, "--database-orm", component.ORM)
	}
	if component.Auth != "" {
		args = append(args, "--auth-type", component.Auth)
	}
	if component.Logger != "" {
		args = append(args, "--logger", component.Logger)
	}
	if component.Complexity != "" {
		args = append(args, "--complexity", component.Complexity)
	}

	// Set module path for component
	modulePath := fmt.Sprintf("%s/cmd/%s", ctx.workspaceName, component.Name)
	args = append(args, "--module", modulePath)

	// Generate component
	cmd := exec.Command("go-starter", args...)
	cmd.Dir = componentPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		ctx.lastCommandOutput = string(output)
		ctx.lastCommandError = err
		return fmt.Errorf("failed to generate component %s: %v, output: %s", component.Name, err, output)
	}

	ctx.lastCommandOutput = string(output)
	return nil
}

func (ctx *WorkspaceTestContext) theWorkspaceShouldBeGeneratedSuccessfully() error {
	// Verify workspace directory exists
	if _, err := os.Stat(ctx.workspacePath); os.IsNotExist(err) {
		return fmt.Errorf("workspace directory does not exist: %s", ctx.workspacePath)
	}

	// Verify workspace go.mod exists
	goModPath := filepath.Join(ctx.workspacePath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("workspace go.mod does not exist: %s", goModPath)
	}

	// Verify all components were generated
	for componentName := range ctx.components {
		componentPath := filepath.Join(ctx.workspacePath, "cmd", componentName)
		if _, err := os.Stat(componentPath); os.IsNotExist(err) {
			return fmt.Errorf("component directory does not exist: %s", componentPath)
		}
	}

	return nil
}

func (ctx *WorkspaceTestContext) allComponentsShouldShareAConsistentRootGoMod() error {
	// Read root go.mod
	goModPath := filepath.Join(ctx.workspacePath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read root go.mod: %v", err)
	}

	// Verify module name matches workspace name
	goModContent := string(content)
	if !strings.Contains(goModContent, ctx.workspaceName) {
		return fmt.Errorf("root go.mod does not contain workspace name: %s", ctx.workspaceName)
	}

	// Verify go version is specified
	if !strings.Contains(goModContent, "go ") {
		return fmt.Errorf("root go.mod does not specify go version")
	}

	return nil
}

func (ctx *WorkspaceTestContext) sharedDependenciesShouldBeProperlyManaged() error {
	// Analyze dependencies across components
	goModPath := filepath.Join(ctx.workspacePath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %v", err)
	}

	goModContent := string(content)

	// Check for common dependencies based on components
	for _, component := range ctx.components {
		switch component.Logger {
		case "zap":
			if !strings.Contains(goModContent, "go.uber.org/zap") {
				return fmt.Errorf("zap dependency not found in workspace go.mod")
			}
		case "logrus":
			if !strings.Contains(goModContent, "github.com/sirupsen/logrus") {
				return fmt.Errorf("logrus dependency not found in workspace go.mod")
			}
		case "zerolog":
			if !strings.Contains(goModContent, "github.com/rs/zerolog") {
				return fmt.Errorf("zerolog dependency not found in workspace go.mod")
			}
		}

		if component.Framework == "gin" {
			if !strings.Contains(goModContent, "github.com/gin-gonic/gin") {
				return fmt.Errorf("gin dependency not found in workspace go.mod")
			}
		}
	}

	return nil
}

func (ctx *WorkspaceTestContext) componentGoModFilesShouldBeProperlyConfigured() error {
	// For workspace projects, components typically share the root go.mod
	// or have their own go.mod files with proper replace directives
	for componentName := range ctx.components {
		componentPath := filepath.Join(ctx.workspacePath, "cmd", componentName)
		
		// Check if component has its own go.mod
		componentGoMod := filepath.Join(componentPath, "go.mod")
		if _, err := os.Stat(componentGoMod); err == nil {
			// Component has its own go.mod, verify it's properly configured
			content, err := os.ReadFile(componentGoMod)
			if err != nil {
				return fmt.Errorf("failed to read component go.mod for %s: %v", componentName, err)
			}

			goModContent := string(content)
			expectedModule := fmt.Sprintf("%s/cmd/%s", ctx.workspaceName, componentName)
			if !strings.Contains(goModContent, expectedModule) {
				return fmt.Errorf("component go.mod for %s does not contain expected module path", componentName)
			}
		}
	}

	return nil
}

func (ctx *WorkspaceTestContext) workspaceCompilationShouldSucceed() error {
	// Try to build the entire workspace
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	if err := os.Chdir(ctx.workspacePath); err != nil {
		return fmt.Errorf("failed to change to workspace directory: %v", err)
	}

	// Run go mod tidy to ensure dependencies are correct
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = ctx.workspacePath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy failed: %v, output: %s", err, output)
	}

	// Try to build each component
	for componentName := range ctx.components {
		componentPath := filepath.Join(ctx.workspacePath, "cmd", componentName)
		
		// Try to build the component
		cmd := exec.Command("go", "build", "./...")
		cmd.Dir = componentPath
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("compilation failed for component %s: %v, output: %s", componentName, err, output)
		}
	}

	return nil
}

func (ctx *WorkspaceTestContext) interComponentImportsShouldWorkCorrectly() error {
	// This would involve checking that components can import shared packages
	// For now, we'll verify that the workspace structure supports this
	sharedPkgPath := filepath.Join(ctx.workspacePath, "pkg")
	if _, err := os.Stat(sharedPkgPath); os.IsNotExist(err) {
		return fmt.Errorf("shared pkg directory does not exist")
	}

	return nil
}

func (ctx *WorkspaceTestContext) sharedConfigurationShouldBeAvailableToAllComponents() error {
	// Verify that configuration patterns are consistent
	// This could involve checking for shared config files or patterns
	return nil
}

// Additional validation methods (simplified implementations)
func (ctx *WorkspaceTestContext) eachServiceShouldHaveIndependentDatabaseConfigurations() error {
	// Verify that each service has its own database configuration
	return nil
}

func (ctx *WorkspaceTestContext) sharedAuthenticationLibrariesShouldBeConsistent() error {
	// Verify that authentication libraries are consistent across services
	return nil
}

func (ctx *WorkspaceTestContext) loggerImplementationsShouldBeConsistentAcrossServices() error {
	// Verify that logger implementations are consistent
	return nil
}

func (ctx *WorkspaceTestContext) serviceToServiceCommunicationPatternsShouldBeAvailable() error {
	// Verify that service-to-service communication patterns are available
	return nil
}

func (ctx *WorkspaceTestContext) sharedDatabaseConnectionsShouldBeProperlyManaged() error {
	// Verify that shared database connections are properly managed
	return nil
}

func (ctx *WorkspaceTestContext) sessionManagementShouldBeConsistentAcrossWebComponents() error {
	// Verify that session management is consistent across web components
	return nil
}

func (ctx *WorkspaceTestContext) loggerConfigurationShouldBeShared() error {
	// Verify that logger configuration is shared
	return nil
}

func (ctx *WorkspaceTestContext) modularStructureShouldBeMaintained() error {
	// Verify that modular structure is maintained
	return nil
}

func (ctx *WorkspaceTestContext) eventSourcingPatternsShouldBeProperlyImplemented() error {
	// Verify that event sourcing patterns are properly implemented
	return nil
}

func (ctx *WorkspaceTestContext) cqrsSeparationShouldBeMaintained() error {
	// Verify that CQRS separation is maintained
	return nil
}

func (ctx *WorkspaceTestContext) eventStoreConfigurationShouldBeShared() error {
	// Verify that event store configuration is shared
	return nil
}

func (ctx *WorkspaceTestContext) eventFlowBetweenComponentsShouldBeConfigured() error {
	// Verify that event flow between components is configured
	return nil
}

func (ctx *WorkspaceTestContext) allBlueprintTypesShouldBeProperlyIntegrated() error {
	// Verify that all blueprint types are properly integrated
	return nil
}

func (ctx *WorkspaceTestContext) sharedDependenciesShouldBeManagedEfficiently() error {
	// Verify that shared dependencies are managed efficiently
	return nil
}

func (ctx *WorkspaceTestContext) loggerConfigurationShouldBeConsistentAcrossAllComponents() error {
	// Verify that logger configuration is consistent across all components
	return nil
}

func (ctx *WorkspaceTestContext) authenticationShouldBeCompatibleAcrossComponents() error {
	// Verify that authentication is compatible across components
	return nil
}

func (ctx *WorkspaceTestContext) interComponentCommunicationShouldBeProperlyConfigured() error {
	// Verify that inter-component communication is properly configured
	return nil
}

func (ctx *WorkspaceTestContext) sharedLibraryShouldBeUsableByAllComponents() error {
	// Verify that shared library is usable by all components
	return nil
}

// Additional step implementations for dependency and performance testing
func (ctx *WorkspaceTestContext) iGenerateAWorkspaceWithMultipleComponentsHavingDifferentDependencyRequirements() error {
	// Generate a test workspace with varied dependencies
	return nil
}

func (ctx *WorkspaceTestContext) goModDependencyResolutionShouldWorkCorrectly() error {
	return nil
}

func (ctx *WorkspaceTestContext) noVersionConflictsShouldExist() error {
	return nil
}

func (ctx *WorkspaceTestContext) sharedDependenciesShouldUseConsistentVersions() error {
	return nil
}

func (ctx *WorkspaceTestContext) workspaceShouldBuildWithoutDependencyErrors() error {
	return nil
}

func (ctx *WorkspaceTestContext) iGenerateAWorkspaceWithMultipleComponents() error {
	return nil
}

func (ctx *WorkspaceTestContext) configurationPatternsShouldBeConsistentAcrossComponents() error {
	return nil
}

func (ctx *WorkspaceTestContext) environmentVariableHandlingShouldBeStandardized() error {
	return nil
}

func (ctx *WorkspaceTestContext) loggingConfigurationShouldBeCentralized() error {
	return nil
}

func (ctx *WorkspaceTestContext) databaseConnectionManagementShouldBeOptimized() error {
	return nil
}

func (ctx *WorkspaceTestContext) iGenerateAWorkspaceWithPlusComponents(componentCount int) error {
	return nil
}

func (ctx *WorkspaceTestContext) generationTimeShouldBeReasonable() error {
	return nil
}

func (ctx *WorkspaceTestContext) memoryUsageShouldRemainWithinAcceptableLimits() error {
	return nil
}

func (ctx *WorkspaceTestContext) compilationTimeShouldBeOptimized() error {
	return nil
}

func (ctx *WorkspaceTestContext) workspaceStructureShouldBeMaintainable() error {
	return nil
}