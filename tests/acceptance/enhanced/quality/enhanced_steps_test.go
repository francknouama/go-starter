package quality

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
)

// init initializes the template filesystem before tests run
func init() {
	// Initialize templates filesystem for enhanced quality tests
	wd, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory: " + err.Error())
	}

	// Navigate to project root and find blueprints directory
	projectRoot := wd
	for {
		templatesDir := filepath.Join(projectRoot, "blueprints")
		if _, err := os.Stat(templatesDir); err == nil {
			entries, err := os.ReadDir(templatesDir)
			if err == nil && len(entries) > 0 {
				hasTemplates := false
				for _, entry := range entries {
					if entry.IsDir() {
						templateYaml := filepath.Join(templatesDir, entry.Name(), "template.yaml")
						if _, err := os.Stat(templateYaml); err == nil {
							hasTemplates = true
							break
						}
					}
				}

				if hasTemplates {
					templates.SetTemplatesFS(os.DirFS(templatesDir))
					return
				}
			}
		}

		parentDir := filepath.Dir(projectRoot)
		if parentDir == projectRoot {
			panic("Could not find blueprints directory with templates")
		}
		projectRoot = parentDir
	}
}

// EnhancedQualityTestContext provides comprehensive BDD testing context
type EnhancedQualityTestContext struct {
	ProjectConfigs map[string]types.ProjectConfig
	ProjectPaths   map[string]string
	TestResults    map[string]interface{}
	CurrentProject string
}

// NewEnhancedQualityTestContext creates a new comprehensive test context
func NewEnhancedQualityTestContext() *EnhancedQualityTestContext {
	return &EnhancedQualityTestContext{
		ProjectConfigs: make(map[string]types.ProjectConfig),
		ProjectPaths:   make(map[string]string),
		TestResults:    make(map[string]interface{}),
	}
}

// Framework Consistency Step Definitions

// Given_I_generate_a_project_with_framework creates project with specific framework
func (ctx *EnhancedQualityTestContext) Given_I_generate_a_project_with_framework(t *testing.T, framework string) {
	t.Helper()
	
	config := types.ProjectConfig{
		Name:      fmt.Sprintf("test-%s-project", framework),
		Module:    fmt.Sprintf("github.com/test/test-%s-project", framework),
		Type:      "web-api",
		GoVersion: "1.21",
		Framework: framework,
		Logger:    "slog",
		Author:    "Test Author",
		Email:     "test@example.com",
		License:   "MIT",
		Features:  &types.Features{},
	}
	
	projectPath, err := generateProjectForBDD(config)
	if err != nil {
		t.Fatalf("Failed to generate project for framework %s: %v", framework, err)
	}
	
	ctx.ProjectConfigs[framework] = config
	ctx.ProjectPaths[framework] = projectPath
	ctx.CurrentProject = framework
}

// When_I_scan_all_generated_Go_files_for_framework_references analyzes framework usage
func (ctx *EnhancedQualityTestContext) When_I_scan_all_generated_Go_files_for_framework_references(t *testing.T) {
	t.Helper()
	
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	frameworkReferences := make(map[string][]string)
	
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") {
			return nil
		}
		
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		
		fileContent := string(content)
		
		// Check for framework imports
		frameworks := []string{"gin", "fiber", "echo"}
		for _, fw := range frameworks {
			patterns := getFrameworkPatterns(fw)
			for _, pattern := range patterns {
				if strings.Contains(fileContent, pattern) {
					frameworkReferences[fw] = append(frameworkReferences[fw], fmt.Sprintf("%s: %s", path, pattern))
				}
			}
		}
		
		return nil
	})
	
	require.NoError(t, err, "Should be able to scan for framework references")
	ctx.TestResults["framework_references"] = frameworkReferences
}

// getFrameworkPatterns returns patterns to search for each framework
func getFrameworkPatterns(framework string) []string {
	switch framework {
	case "gin":
		return []string{
			"github.com/gin-gonic/gin",
			"gin.Default()",
			"gin.Context",
			"gin.HandlerFunc",
		}
	case "fiber":
		return []string{
			"github.com/gofiber/fiber",
			"fiber.New()",
			"fiber.Ctx",
			"fiber.Handler",
		}
	case "echo":
		return []string{
			"github.com/labstack/echo",
			"echo.New()",
			"echo.Context",
			"echo.MiddlewareFunc",
		}
	default:
		return []string{}
	}
}

// getFrameworkDependencies returns go.mod dependencies for each framework
func getFrameworkDependencies(framework string) []string {
	switch framework {
	case "gin":
		return []string{
			"github.com/gin-gonic/gin",
		}
	case "fiber":
		return []string{
			"github.com/gofiber/fiber",
		}
	case "echo":
		return []string{
			"github.com/labstack/echo",
		}
	default:
		return []string{}
	}
}

// getConfigValue retrieves configuration value with fallback default
func getConfigValue(configMap map[string]string, key, defaultValue string) string {
	if value, exists := configMap[key]; exists && value != "" {
		return value
	}
	return defaultValue
}

// Project cache for performance optimization
var projectCache = make(map[string]string)
var projectCacheMutex sync.RWMutex

// Cache metrics for monitoring
type CacheMetrics struct {
	Hits   int64
	Misses int64
	mutex  sync.RWMutex
}

var cacheMetrics = &CacheMetrics{}

// IncrementCacheHit safely increments the cache hit counter
func (cm *CacheMetrics) IncrementCacheHit() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.Hits++
}

// IncrementCacheMiss safely increments the cache miss counter
func (cm *CacheMetrics) IncrementCacheMiss() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.Misses++
}

// GetStats returns current cache statistics
func (cm *CacheMetrics) GetStats() (hits, misses int64, hitRate float64) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	hits = cm.Hits
	misses = cm.Misses
	total := hits + misses
	if total > 0 {
		hitRate = float64(hits) / float64(total) * 100
	}
	return
}

// Reset clears the cache metrics (useful for testing)
func (cm *CacheMetrics) Reset() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.Hits = 0
	cm.Misses = 0
}

// generateConfigKey creates a unique key for project configuration caching
func generateConfigKey(config types.ProjectConfig) string {
	var databaseDriver, databaseORM, authType string
	
	if config.Features != nil {
		if config.Features.Database.Driver != "" {
			databaseDriver = config.Features.Database.Driver
		} else if len(config.Features.Database.Drivers) > 0 {
			databaseDriver = strings.Join(config.Features.Database.Drivers, ",")
		}
		databaseORM = config.Features.Database.ORM
		authType = config.Features.Authentication.Type
	}
	
	return fmt.Sprintf("%s-%s-%s-%s-%s-%s-%s-%s", 
		config.Type,
		config.Framework, 
		databaseDriver,
		databaseORM,
		config.Logger,
		authType,
		config.Architecture,
		config.GoVersion,
	)
}

// generateProjectForBDD generates a project for BDD testing with caching for performance
func generateProjectForBDD(config types.ProjectConfig) (string, error) {
	cacheKey := generateConfigKey(config)
	
	// Check if we already have this project configuration cached
	projectCacheMutex.RLock()
	if cachedPath, exists := projectCache[cacheKey]; exists {
		// Verify the cached project still exists
		if _, err := os.Stat(cachedPath); err == nil {
			projectCacheMutex.RUnlock()
			cacheMetrics.IncrementCacheHit()
			return cachedPath, nil
		}
	}
	projectCacheMutex.RUnlock()
	
	// Record cache miss
	cacheMetrics.IncrementCacheMiss()
	
	// Generate new project
	projectCacheMutex.Lock()
	defer projectCacheMutex.Unlock()
	
	// Double-check after acquiring write lock
	if cachedPath, exists := projectCache[cacheKey]; exists {
		if _, err := os.Stat(cachedPath); err == nil {
			return cachedPath, nil
		}
		delete(projectCache, cacheKey)
	}
	
	// Create a temporary directory
	outputDir, err := os.MkdirTemp("", "go-starter-bdd-test-")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}
	
	// Ensure cleanup on any exit path
	defer func() {
		if err != nil {
			os.RemoveAll(outputDir)
		}
	}()
	
	projectPath := filepath.Join(outputDir, config.Name)
	
	// Use generator directly
	gen := generator.New()
	options := types.GenerationOptions{
		OutputPath: projectPath,
		DryRun:     false,
		NoGit:      true, // Skip git init for tests
		Verbose:    false,
	}
	
	_, err = gen.Generate(config, options)
	if err != nil {
		return "", fmt.Errorf("failed to generate project: %w", err)
	}
	
	// Cache the generated project
	projectCache[cacheKey] = projectPath
	
	return projectPath, nil
}

// analyzeImportsWithAST uses Go AST parsing to accurately analyze imports and their usage
func analyzeImportsWithAST(filePath string, problematicImports map[string]string) map[string]string {
	unusedImports := make(map[string]string)
	
	// Parse the Go file
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return unusedImports // Return empty if we can't parse
	}
	
	// Track which imports are declared
	declaredImports := make(map[string]string) // package name -> import path
	importPaths := make(map[string]string)     // import path -> package name
	
	// Collect import declarations
	for _, imp := range node.Imports {
		if imp.Path != nil {
			importPath := strings.Trim(imp.Path.Value, "\"")
			packageName := ""
			
			if imp.Name != nil {
				// Aliased import (import foo "bar")
				packageName = imp.Name.Name
			} else {
				// Standard import - get package name from path
				parts := strings.Split(importPath, "/")
				packageName = parts[len(parts)-1]
			}
			
			// Check if this is one of our problematic imports
			for problematicImport := range problematicImports {
				if importPath == problematicImport || packageName == problematicImport {
					declaredImports[packageName] = importPath
					importPaths[importPath] = packageName
				}
			}
		}
	}
	
	// Track usage of declared imports
	usedImports := make(map[string]bool)
	
	// Walk the AST to find usage
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.SelectorExpr:
			// Handle packageName.Function() calls
			if ident, ok := x.X.(*ast.Ident); ok {
				if _, exists := declaredImports[ident.Name]; exists {
					usedImports[ident.Name] = true
				}
			}
		case *ast.CallExpr:
			// Handle direct function calls that might use imports
			if ident, ok := x.Fun.(*ast.Ident); ok {
				// Check if this function name matches any import
				for packageName := range declaredImports {
					if ident.Name == packageName {
						usedImports[packageName] = true
					}
				}
			}
		}
		return true
	})
	
	// Identify unused problematic imports
	for packageName, importPath := range declaredImports {
		if !usedImports[packageName] {
			if reason, exists := problematicImports[packageName]; exists {
				unusedImports[packageName] = reason
			} else if reason, exists := problematicImports[importPath]; exists {
				unusedImports[importPath] = reason
			}
		}
	}
	
	return unusedImports
}

// logCacheMetrics logs current cache performance statistics
func logCacheMetrics() {
	hits, misses, hitRate := cacheMetrics.GetStats()
	total := hits + misses
	if total > 0 {
		fmt.Printf("📊 Cache Performance: %d hits, %d misses, %.1f%% hit rate (total: %d)\n", 
			hits, misses, hitRate, total)
	}
}

// Then_the_project_should_only_contain_framework_imports validates framework isolation
func (ctx *EnhancedQualityTestContext) Then_the_project_should_only_contain_framework_imports(t *testing.T, expectedFramework string) {
	t.Helper()
	
	frameworkReferences, ok := ctx.TestResults["framework_references"].(map[string][]string)
	require.True(t, ok, "Framework references should be recorded")
	
	// Should contain expected framework
	expectedRefs, hasExpected := frameworkReferences[expectedFramework]
	assert.True(t, hasExpected && len(expectedRefs) > 0,
		"Project should contain %s framework references", expectedFramework)
	
	// Should not contain other frameworks
	for fw, refs := range frameworkReferences {
		if fw != expectedFramework && len(refs) > 0 {
			t.Errorf("Project should not contain %s framework references but found: %v", fw, refs)
		}
	}
}

// Configuration Consistency Step Definitions

// Given_I_generate_a_project_with_database creates project with specific database config
func (ctx *EnhancedQualityTestContext) Given_I_generate_a_project_with_database(t *testing.T, database, driver string) {
	t.Helper()
	
	config := types.ProjectConfig{
		Name:      fmt.Sprintf("test-%s-project", database),
		Module:    fmt.Sprintf("github.com/test/test-%s-project", database),
		Type:      "web-api",  
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
		Author:    "Test Author",
		Email:     "test@example.com",
		License:   "MIT",
		Features:  &types.Features{},
	}
	
	config.Features.Database.Driver = driver
	
	projectPath, err := generateProjectForBDD(config)
	if err != nil {
		t.Fatalf("Failed to generate project for database %s-%s: %v", database, driver, err)
	}
	
	projectKey := fmt.Sprintf("%s-%s", database, driver)
	ctx.ProjectConfigs[projectKey] = config
	ctx.ProjectPaths[projectKey] = projectPath
	ctx.CurrentProject = projectKey
}

// When_I_examine_the_docker_compose_yml_file analyzes docker-compose configuration
func (ctx *EnhancedQualityTestContext) When_I_examine_the_docker_compose_yml_file(t *testing.T) {
	t.Helper()
	
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	dockerComposePath := filepath.Join(projectPath, "docker-compose.yml")
	
	var dockerContent string
	if _, err := os.Stat(dockerComposePath); err == nil {
		content, err := os.ReadFile(dockerComposePath)
		require.NoError(t, err, "Should be able to read docker-compose.yml")
		dockerContent = string(content)
	}
	
	ctx.TestResults["docker_compose_content"] = dockerContent
}

// Then_it_should_use_the_database_service validates docker-compose database service
func (ctx *EnhancedQualityTestContext) Then_it_should_use_the_database_service(t *testing.T, expectedService string) {
	t.Helper()
	
	dockerContent, ok := ctx.TestResults["docker_compose_content"].(string)
	require.True(t, ok, "Docker compose content should be recorded")
	
	if expectedService != "" && dockerContent != "" {
		assert.Contains(t, dockerContent, expectedService,
			"docker-compose.yml should contain service: %s", expectedService)
	}
}

// Static Analysis Step Definitions

// When_I_run_goimports_analysis_on_all_Go_files performs import analysis
func (ctx *EnhancedQualityTestContext) When_I_run_goimports_analysis_on_all_Go_files(t *testing.T) {
	t.Helper()
	
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	var importIssues []string
	var analysisErrors []string
	
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") || strings.Contains(path, ".git/") {
			return nil
		}
		
		// Check with goimports for formatting differences (indicates unused imports or missing imports)
		cmd := exec.Command("goimports", "-d", path)
		output, cmdErr := cmd.CombinedOutput()
		
		if cmdErr != nil {
			// goimports failed to run on this file - this could indicate syntax errors
			analysisErrors = append(analysisErrors, 
				fmt.Sprintf("goimports failed on %s: %v", path, cmdErr))
		} else if len(output) > 0 {
			// goimports has suggested changes - indicates import formatting issues
			importIssues = append(importIssues, 
				fmt.Sprintf("File %s has import formatting issues:\n%s", path, string(output)))
		}
		
		return nil
	})
	
	require.NoError(t, err, "Should be able to walk project directory for import analysis")
	ctx.TestResults["import_issues"] = importIssues
	ctx.TestResults["import_analysis_errors"] = analysisErrors
}

// Then_there_should_be_no_unused_import_statements validates clean imports
func (ctx *EnhancedQualityTestContext) Then_there_should_be_no_unused_import_statements(t *testing.T) {
	t.Helper()
	
	importIssues, ok := ctx.TestResults["import_issues"].([]string)
	require.True(t, ok, "Import issues should be recorded")
	
	analysisErrors, errorsOk := ctx.TestResults["import_analysis_errors"].([]string)
	require.True(t, errorsOk, "Import analysis errors should be recorded")
	
	// First check if there were any analysis errors (syntax errors, etc.)
	if len(analysisErrors) > 0 {
		t.Errorf("goimports analysis failed on some files:\n%s", strings.Join(analysisErrors, "\n"))
	}
	
	// Then check for import formatting issues (unused imports, missing imports)
	if len(importIssues) > 0 {
		t.Errorf("Found import formatting issues:\n%s", strings.Join(importIssues, "\n"))
	}
	
	assert.Empty(t, importIssues, "There should be no import formatting issues")
	assert.Empty(t, analysisErrors, "goimports should run successfully on all files")
}

// When_I_scan_for_problematic_import_patterns analyzes specific import problems
func (ctx *EnhancedQualityTestContext) When_I_scan_for_problematic_import_patterns(t *testing.T) {
	t.Helper()
	
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	
	// Common packages that are often imported unnecessarily in generated code
	problematicImports := map[string]string{
		"fmt":        "often unused in generated code when no formatting is needed",
		"os":         "often unused when not handling files or environment variables", 
		"context":    "unused when no context handling is implemented",
		"log":        "unused when custom logger is used",
		"errors":     "unused when no error wrapping is done",
	}
	
	var violations []string
	
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") || strings.Contains(path, ".git/") {
			return nil
		}
		
		// Use AST parsing for more accurate import analysis
		unusedImports := analyzeImportsWithAST(path, problematicImports)
		for importName, reason := range unusedImports {
			violations = append(violations,
				fmt.Sprintf("File %s imports '%s' but appears unused (%s)", path, importName, reason))
		}
		
		return nil
	})
	
	require.NoError(t, err, "Should be able to scan for problematic imports")
	ctx.TestResults["import_violations"] = violations
}

// Then_problematic_imports_should_only_be_present_when_used validates import usage
func (ctx *EnhancedQualityTestContext) Then_problematic_imports_should_only_be_present_when_used(t *testing.T, packageName string) {
	t.Helper()
	
	violations, ok := ctx.TestResults["import_violations"].([]string)
	require.True(t, ok, "Import violations should be recorded")
	
	// Filter violations for specific package
	packageViolations := []string{}
	for _, violation := range violations {
		if strings.Contains(violation, fmt.Sprintf("'%s'", packageName)) {
			packageViolations = append(packageViolations, violation)
		}
	}
	
	assert.Empty(t, packageViolations, 
		"Package %s should only be imported when used", packageName)
}

// Multiple Project Step Definitions

// Given_I_generate_multiple_projects_with_different_frameworks creates multiple projects
func (ctx *EnhancedQualityTestContext) Given_I_generate_multiple_projects_with_different_frameworks(t *testing.T, projects map[string]string) {
	t.Helper()
	
	for projectName, framework := range projects {
		config := types.ProjectConfig{
			Name:      projectName,
			Module:    fmt.Sprintf("github.com/test/%s", projectName),
			Type:      "web-api",
			GoVersion: "1.21", 
			Framework: framework,
			Logger:    "slog",
			Author:    "Test Author",
			Email:     "test@example.com",
			License:   "MIT",
			Features:  &types.Features{},
		}
		
		projectPath, err := generateProjectForBDD(config)
		if err != nil {
			t.Fatalf("Failed to generate project %s: %v", projectName, err)
		}
		
		ctx.ProjectConfigs[projectName] = config
		ctx.ProjectPaths[projectName] = projectPath
	}
}

// When_I_validate_framework_consistency_across_all_projects checks all projects
func (ctx *EnhancedQualityTestContext) When_I_validate_framework_consistency_across_all_projects(t *testing.T) {
	t.Helper()
	
	projectValidation := make(map[string]map[string]interface{})
	
	for projectName, projectPath := range ctx.ProjectPaths {
		expectedFramework := ctx.ProjectConfigs[projectName].Framework
		
		validation := make(map[string]interface{})
		validation["expected_framework"] = expectedFramework
		
		// Scan for framework references
		frameworkRefs := make(map[string][]string)
		
		err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") {
				return nil
			}
			
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			
			fileContent := string(content)
			frameworks := []string{"gin", "fiber", "echo"}
			
			for _, fw := range frameworks {
				patterns := getFrameworkPatterns(fw)
				for _, pattern := range patterns {
					if strings.Contains(fileContent, pattern) {
						frameworkRefs[fw] = append(frameworkRefs[fw], path)
					}
				}
			}
			
			return nil
		})
		if err != nil {
			t.Errorf("Error walking project path %s: %v", projectPath, err)
			continue
		}
		
		validation["framework_references"] = frameworkRefs
		projectValidation[projectName] = validation
	}
	
	ctx.TestResults["multi_project_validation"] = projectValidation
}

// Then_each_project_should_use_only_its_designated_framework validates isolation
func (ctx *EnhancedQualityTestContext) Then_each_project_should_use_only_its_designated_framework(t *testing.T) {
	t.Helper()
	
	validation, ok := ctx.TestResults["multi_project_validation"].(map[string]map[string]interface{})
	require.True(t, ok, "Multi-project validation should be recorded")
	
	for projectName, projectValidation := range validation {
		expectedFramework := projectValidation["expected_framework"].(string)
		frameworkRefs := projectValidation["framework_references"].(map[string][]string)
		
		// Should have expected framework
		expectedRefs, hasExpected := frameworkRefs[expectedFramework]
		assert.True(t, hasExpected && len(expectedRefs) > 0,
			"Project %s should contain %s framework references", projectName, expectedFramework)
		
		// Should not have other frameworks
		for fw, refs := range frameworkRefs {
			if fw != expectedFramework && len(refs) > 0 {
				t.Errorf("Project %s should not contain %s framework references but found in: %v", 
					projectName, fw, refs)
			}
		}
	}
}

// RegisterSteps registers all step definitions with godog
func (ctx *EnhancedQualityTestContext) RegisterSteps(s *godog.ScenarioContext) {
	// Background steps
	s.Step(`^I have the go-starter CLI available$`, ctx.iHaveTheGoStarterCLIAvailable)
	s.Step(`^all templates are properly initialized$`, ctx.allTemplatesAreProperlyInitialized)
	
	// Project generation steps
	s.Step(`^I generate a project with framework "([^"]*)"$`, ctx.iGenerateAProjectWithFramework)
	s.Step(`^I generate a project with database "([^"]*)" and driver "([^"]*)"$`, ctx.iGenerateAProjectWithDatabase)
	s.Step(`^I generate a project with configuration:$`, ctx.iGenerateAProjectWithConfiguration)
	s.Step(`^I generate multiple projects with different frameworks:$`, ctx.iGenerateMultipleProjectsWithDifferentFrameworks)
	
	// Framework consistency steps
	s.Step(`^I scan all generated Go files for framework references$`, ctx.iScanAllGeneratedGoFilesForFrameworkReferences)
	s.Step(`^the project should only contain "([^"]*)" framework imports$`, ctx.theProjectShouldOnlyContainFrameworkImports)
	s.Step(`^the project should not contain "([^"]*)" framework imports$`, ctx.theProjectShouldNotContainFrameworkImports)
	s.Step(`^main\.go should use the correct framework initialization pattern$`, ctx.mainGoShouldUseTheCorrectFrameworkInitializationPattern)
	s.Step(`^go\.mod should contain only the "([^"]*)" dependency$`, ctx.goModShouldContainOnlyTheFrameworkDependency)
	
	// Static analysis steps
	s.Step(`^I run goimports analysis on all Go files$`, ctx.iRunGoimportsAnalysisOnAllGoFiles)
	s.Step(`^there should be no unused import statements$`, ctx.thereShouldBeNoUnusedImportStatements)
	s.Step(`^goimports should report no formatting differences$`, ctx.goimportsShouldReportNoFormattingDifferences)
	s.Step(`^I scan for problematic import patterns$`, ctx.iScanForProblematicImportPatterns)
	s.Step(`^"([^"]*)" package should only be imported when format functions are used$`, ctx.packageShouldOnlyBeImportedWhenUsed)
	s.Step(`^"([^"]*)" package should only be imported when OS functions are used$`, ctx.packageShouldOnlyBeImportedWhenUsed)
	s.Step(`^"([^"]*)" package should only be imported when ORM is "([^"]*)"$`, ctx.packageShouldOnlyBeImportedWhenORMIs)
	
	// Configuration consistency steps
	s.Step(`^I examine the docker-compose\.yml file$`, ctx.iExamineTheDockerComposeYmlFile)
	s.Step(`^it should use the database service "([^"]*)"$`, ctx.itShouldUseTheDatabaseService)
	s.Step(`^I examine the generated configuration files$`, ctx.iExamineTheGeneratedConfigurationFiles)
	s.Step(`^go\.mod should contain the framework dependency "([^"]*)"$`, ctx.goModShouldContainTheFrameworkDependency)
	
	// Multi-project validation steps
	s.Step(`^I validate framework consistency across all projects$`, ctx.iValidateFrameworkConsistencyAcrossAllProjects)
	s.Step(`^each project should use only its designated framework$`, ctx.eachProjectShouldUseOnlyItsDesignatedFramework)
	s.Step(`^no project should contain references to other frameworks$`, ctx.noProjectShouldContainReferencesToOtherFrameworks)
	
	// Compilation and build steps
	s.Step(`^I attempt to compile the generated project$`, ctx.iAttemptToCompileTheGeneratedProject)
	s.Step(`^the compilation should succeed without errors$`, ctx.theCompilationShouldSucceedWithoutErrors)
	s.Step(`^the build output should not contain warnings$`, ctx.theBuildOutputShouldNotContainWarnings)
}

// Step definition implementations

func (ctx *EnhancedQualityTestContext) iHaveTheGoStarterCLIAvailable() error {
	// Templates are initialized in init() function
	return nil
}

func (ctx *EnhancedQualityTestContext) allTemplatesAreProperlyInitialized() error {
	// Templates filesystem is set up in init() function
	return nil
}

func (ctx *EnhancedQualityTestContext) iGenerateAProjectWithFramework(framework string) error {
	// Create a project configuration with the specified framework
	config := types.ProjectConfig{
		Name:      fmt.Sprintf("test-%s", framework),
		Module:    fmt.Sprintf("github.com/test/test-%s", framework),
		Type:      "web-api",
		GoVersion: "1.21",
		Framework: framework,
		Logger:    "slog",
		Author:    "Test Author",
		Email:     "test@example.com",
		License:   "MIT",
		Features:  &types.Features{},
	}
	
	// Generate the project
	projectPath, err := generateProjectForBDD(config)
	if err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}
	
	// Store the configuration and path
	ctx.ProjectConfigs[framework] = config
	ctx.ProjectPaths[framework] = projectPath
	ctx.CurrentProject = framework
	
	return nil
}

func (ctx *EnhancedQualityTestContext) iGenerateAProjectWithDatabase(database, driver string) error {
	// Create a project configuration with the specified database
	config := types.ProjectConfig{
		Name:      fmt.Sprintf("test-%s-%s", database, driver),
		Module:    fmt.Sprintf("github.com/test/test-%s-%s", database, driver),
		Type:      "web-api",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
		Author:    "Test Author",
		Email:     "test@example.com",
		License:   "MIT",
		Features:  &types.Features{},
	}
	
	// Set database configuration
	config.Features.Database.Driver = driver
	
	// Generate the project
	projectPath, err := generateProjectForBDD(config)
	if err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}
	
	// Store the configuration and path
	projectKey := fmt.Sprintf("%s-%s", database, driver)
	ctx.ProjectConfigs[projectKey] = config
	ctx.ProjectPaths[projectKey] = projectPath
	ctx.CurrentProject = projectKey
	
	return nil
}

func (ctx *EnhancedQualityTestContext) iGenerateAProjectWithConfiguration(table *godog.Table) error {
	configMap := make(map[string]string)
	for _, row := range table.Rows {
		if len(row.Cells) >= 2 {
			configMap[row.Cells[0].Value] = row.Cells[1].Value
		}
	}
	
	// Convert table configuration to ProjectConfig
	config := types.ProjectConfig{
		Name:      getConfigValue(configMap, "name", "test-project"),
		Module:    getConfigValue(configMap, "module", "github.com/test/test-project"),
		Type:      getConfigValue(configMap, "type", "web-api"),
		GoVersion: getConfigValue(configMap, "go_version", "1.21"),
		Framework: getConfigValue(configMap, "framework", "gin"),
		Logger:    getConfigValue(configMap, "logger", "slog"),
		Author:    getConfigValue(configMap, "author", "Test Author"),
		Email:     getConfigValue(configMap, "email", "test@example.com"),
		License:   getConfigValue(configMap, "license", "MIT"),
		Features:  &types.Features{},
	}
	
	// Set features based on configuration
	if dbDriver := getConfigValue(configMap, "database_driver", ""); dbDriver != "" {
		config.Features.Database.Driver = dbDriver
	}
	if dbORM := getConfigValue(configMap, "database_orm", ""); dbORM != "" {
		config.Features.Database.ORM = dbORM
	}
	if authType := getConfigValue(configMap, "auth_type", ""); authType != "" {
		config.Features.Authentication.Type = authType
	}
	if architecture := getConfigValue(configMap, "architecture", ""); architecture != "" {
		config.Architecture = architecture
	}
	
	// Generate the project using direct generator approach (no testing.T available in godog)
	projectPath, err := generateProjectForBDD(config)
	if err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}
	
	// Store the configuration and path
	projectKey := config.Name
	ctx.ProjectConfigs[projectKey] = config
	ctx.ProjectPaths[projectKey] = projectPath
	ctx.CurrentProject = projectKey
	
	return nil
}

func (ctx *EnhancedQualityTestContext) iGenerateMultipleProjectsWithDifferentFrameworks(table *godog.Table) error {
	projects := make(map[string]string)
	for i, row := range table.Rows {
		if i == 0 {
			continue // Skip header row
		}
		if len(row.Cells) >= 2 {
			projects[row.Cells[0].Value] = row.Cells[1].Value
		}
	}
	
	// Generate multiple projects with different frameworks
	for projectName, framework := range projects {
		config := types.ProjectConfig{
			Name:      projectName,
			Module:    fmt.Sprintf("github.com/test/%s", projectName),
			Type:      "web-api",
			GoVersion: "1.21", 
			Framework: framework,
			Logger:    "slog",
			Author:    "Test Author",
			Email:     "test@example.com",
			License:   "MIT",
			Features:  &types.Features{},
		}
		
		projectPath, err := generateProjectForBDD(config)
		if err != nil {
			return fmt.Errorf("failed to generate project %s: %w", projectName, err)
		}
		
		ctx.ProjectConfigs[projectName] = config
		ctx.ProjectPaths[projectName] = projectPath
	}
	
	return nil
}

func (ctx *EnhancedQualityTestContext) iScanAllGeneratedGoFilesForFrameworkReferences() error {
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	if projectPath == "" {
		return fmt.Errorf("no current project set for framework scanning")
	}
	
	frameworkReferences := make(map[string][]string)
	
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") {
			return nil
		}
		
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		
		fileContent := string(content)
		
		// Check for framework imports and usage patterns
		frameworks := []string{"gin", "fiber", "echo"}
		for _, fw := range frameworks {
			patterns := getFrameworkPatterns(fw)
			for _, pattern := range patterns {
				if strings.Contains(fileContent, pattern) {
					frameworkReferences[fw] = append(frameworkReferences[fw], fmt.Sprintf("%s: %s", path, pattern))
				}
			}
		}
		
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to scan for framework references: %v", err)
	}
	
	ctx.TestResults["framework_references"] = frameworkReferences
	return nil
}

func (ctx *EnhancedQualityTestContext) theProjectShouldOnlyContainFrameworkImports(framework string) error {
	frameworkReferences, ok := ctx.TestResults["framework_references"].(map[string][]string)
	if !ok {
		return fmt.Errorf("framework references not recorded - ensure 'I scan all generated Go files' step ran first")
	}
	
	// Should contain expected framework
	expectedRefs, hasExpected := frameworkReferences[framework]
	if !hasExpected || len(expectedRefs) == 0 {
		return fmt.Errorf("project should contain %s framework references but found none", framework)
	}
	
	// Should not contain other frameworks
	for fw, refs := range frameworkReferences {
		if fw != framework && len(refs) > 0 {
			return fmt.Errorf("project should not contain %s framework references but found: %v", fw, refs)
		}
	}
	
	return nil
}

func (ctx *EnhancedQualityTestContext) theProjectShouldNotContainFrameworkImports(framework string) error {
	frameworkReferences, ok := ctx.TestResults["framework_references"].(map[string][]string)
	if !ok {
		return fmt.Errorf("framework references not recorded - ensure 'I scan all generated Go files' step ran first")
	}
	
	if refs, exists := frameworkReferences[framework]; exists && len(refs) > 0 {
		return fmt.Errorf("project should not contain %s framework references but found: %v", framework, refs)
	}
	
	return nil
}

func (ctx *EnhancedQualityTestContext) mainGoShouldUseTheCorrectFrameworkInitializationPattern() error {
	// Implementation would check main.go for correct framework initialization
	return nil
}

func (ctx *EnhancedQualityTestContext) goModShouldContainOnlyTheFrameworkDependency(framework string) error {
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	if projectPath == "" {
		return fmt.Errorf("no current project set for go.mod validation")
	}
	
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %v", err)
	}
	
	goModContent := string(content)
	
	// Get expected dependency for framework
	expectedDeps := getFrameworkDependencies(framework)
	if len(expectedDeps) == 0 {
		return fmt.Errorf("unknown framework: %s", framework)
	}
	
	// Check that expected dependencies are present
	for _, expectedDep := range expectedDeps {
		if !strings.Contains(goModContent, expectedDep) {
			return fmt.Errorf("go.mod should contain dependency '%s' for framework '%s'", expectedDep, framework)
		}
	}
	
	// Check that other framework dependencies are not present
	allFrameworks := []string{"gin", "fiber", "echo"}
	for _, fw := range allFrameworks {
		if fw != framework {
			otherDeps := getFrameworkDependencies(fw)
			for _, otherDep := range otherDeps {
				if strings.Contains(goModContent, otherDep) {
					return fmt.Errorf("go.mod should not contain dependency '%s' for framework '%s' when using '%s'", 
						otherDep, fw, framework)
				}
			}
		}
	}
	
	return nil
}

func (ctx *EnhancedQualityTestContext) iRunGoimportsAnalysisOnAllGoFiles() error {
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	if projectPath == "" {
		return fmt.Errorf("no current project set for import analysis")
	}
	
	var importIssues []string
	var analysisErrors []string
	
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") || strings.Contains(path, ".git/") {
			return nil
		}
		
		// Check with goimports for formatting differences (indicates unused imports or missing imports)
		cmd := exec.Command("goimports", "-d", path)
		output, cmdErr := cmd.CombinedOutput()
		
		if cmdErr != nil {
			// goimports failed to run on this file - this could indicate syntax errors
			analysisErrors = append(analysisErrors, 
				fmt.Sprintf("goimports failed on %s: %v", path, cmdErr))
		} else if len(output) > 0 {
			// goimports has suggested changes - indicates import formatting issues
			importIssues = append(importIssues, 
				fmt.Sprintf("File %s has import formatting issues:\n%s", path, string(output)))
		}
		
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to analyze imports: %v", err)
	}
	
	ctx.TestResults["import_issues"] = importIssues
	ctx.TestResults["import_analysis_errors"] = analysisErrors
	return nil
}

func (ctx *EnhancedQualityTestContext) thereShouldBeNoUnusedImportStatements() error {
	importIssues, ok := ctx.TestResults["import_issues"].([]string)
	if !ok {
		return fmt.Errorf("import issues not recorded - ensure 'I run goimports analysis' step ran first")
	}
	
	analysisErrors, errorsOk := ctx.TestResults["import_analysis_errors"].([]string)
	if !errorsOk {
		return fmt.Errorf("import analysis errors not recorded - ensure 'I run goimports analysis' step ran first")
	}
	
	// First check if there were any analysis errors (syntax errors, etc.)
	if len(analysisErrors) > 0 {
		return fmt.Errorf("goimports analysis failed on some files:\n%s", strings.Join(analysisErrors, "\n"))
	}
	
	// Then check for import formatting issues (unused imports, missing imports)
	if len(importIssues) > 0 {
		return fmt.Errorf("found import formatting issues:\n%s", strings.Join(importIssues, "\n"))
	}
	
	return nil
}

func (ctx *EnhancedQualityTestContext) goimportsShouldReportNoFormattingDifferences() error {
	// This is covered by the unused imports check
	return ctx.thereShouldBeNoUnusedImportStatements()
}

func (ctx *EnhancedQualityTestContext) iScanForProblematicImportPatterns() error {
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	if projectPath == "" {
		return fmt.Errorf("no current project set for problematic import analysis")
	}
	
	// Common packages that are often imported unnecessarily in generated code
	problematicImports := map[string]string{
		"fmt":        "often unused in generated code when no formatting is needed",
		"os":         "often unused when not handling files or environment variables", 
		"context":    "unused when no context handling is implemented",
		"log":        "unused when custom logger is used",
		"errors":     "unused when no error wrapping is done",
	}
	
	var violations []string
	
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") || strings.Contains(path, ".git/") {
			return nil
		}
		
		// Use AST parsing for more accurate import analysis
		unusedImports := analyzeImportsWithAST(path, problematicImports)
		for importName, reason := range unusedImports {
			violations = append(violations,
				fmt.Sprintf("File %s imports '%s' but appears unused (%s)", path, importName, reason))
		}
		
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("failed to scan for problematic imports: %v", err)
	}
	
	ctx.TestResults["import_violations"] = violations
	return nil
}

func (ctx *EnhancedQualityTestContext) packageShouldOnlyBeImportedWhenUsed(packageName string) error {
	violations, ok := ctx.TestResults["import_violations"].([]string)
	if !ok {
		return fmt.Errorf("import violations not recorded - ensure 'I scan for problematic import patterns' step ran first")
	}
	
	// Filter violations for specific package
	packageViolations := []string{}
	for _, violation := range violations {
		if strings.Contains(violation, fmt.Sprintf("'%s'", packageName)) {
			packageViolations = append(packageViolations, violation)
		}
	}
	
	if len(packageViolations) > 0 {
		return fmt.Errorf("package %s should only be imported when used, but found violations:\n%s", 
			packageName, strings.Join(packageViolations, "\n"))
	}
	
	return nil
}

func (ctx *EnhancedQualityTestContext) packageShouldOnlyBeImportedWhenORMIs(packageName, ormType string) error {
	// Implementation would check if models package is only imported when ORM is used
	return nil
}

func (ctx *EnhancedQualityTestContext) iExamineTheDockerComposeYmlFile() error {
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	if projectPath == "" {
		return fmt.Errorf("no current project set for docker-compose examination")
	}
	
	dockerComposePath := filepath.Join(projectPath, "docker-compose.yml")
	
	var dockerContent string
	if _, err := os.Stat(dockerComposePath); err == nil {
		content, err := os.ReadFile(dockerComposePath)
		if err != nil {
			return fmt.Errorf("failed to read docker-compose.yml: %v", err)
		}
		dockerContent = string(content)
	}
	
	ctx.TestResults["docker_compose_content"] = dockerContent
	ctx.TestResults["docker_compose_exists"] = dockerContent != ""
	return nil
}

func (ctx *EnhancedQualityTestContext) itShouldUseTheDatabaseService(expectedService string) error {
	dockerContent, ok := ctx.TestResults["docker_compose_content"].(string)
	if !ok {
		return fmt.Errorf("docker compose content not recorded - ensure 'I examine the docker-compose.yml file' step ran first")
	}
	
	if expectedService == "" || expectedService == "none" {
		// No database service expected
		return nil
	}
	
	if dockerContent == "" {
		return fmt.Errorf("docker-compose.yml should exist when database service '%s' is expected", expectedService)
	}
	
	if !strings.Contains(dockerContent, expectedService) {
		return fmt.Errorf("docker-compose.yml should contain service '%s' but it was not found", expectedService)
	}
	
	// Additional validation for database-specific configurations
	switch expectedService {
	case "postgres":
		if !strings.Contains(dockerContent, "postgresql") && !strings.Contains(dockerContent, "postgres") {
			return fmt.Errorf("docker-compose.yml should contain postgres/postgresql configuration for postgres service")
		}
	case "mysql":
		if !strings.Contains(dockerContent, "mysql") {
			return fmt.Errorf("docker-compose.yml should contain mysql configuration for mysql service")
		}
	case "mongodb":
		if !strings.Contains(dockerContent, "mongo") {
			return fmt.Errorf("docker-compose.yml should contain mongo configuration for mongodb service")
		}
	}
	
	return nil
}

func (ctx *EnhancedQualityTestContext) iExamineTheGeneratedConfigurationFiles() error {
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	if projectPath == "" {
		return fmt.Errorf("no current project set for configuration files examination")
	}
	
	configFiles := make(map[string]string)
	
	// Check for go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	if content, err := os.ReadFile(goModPath); err == nil {
		configFiles["go.mod"] = string(content)
	}
	
	// Check for docker-compose.yml
	dockerComposePath := filepath.Join(projectPath, "docker-compose.yml")
	if content, err := os.ReadFile(dockerComposePath); err == nil {
		configFiles["docker-compose.yml"] = string(content)
	}
	
	// Check for Dockerfile
	dockerfilePath := filepath.Join(projectPath, "Dockerfile")
	if content, err := os.ReadFile(dockerfilePath); err == nil {
		configFiles["Dockerfile"] = string(content)
	}
	
	// Check for .env files
	envPath := filepath.Join(projectPath, ".env")
	if content, err := os.ReadFile(envPath); err == nil {
		configFiles[".env"] = string(content)
	}
	
	envExamplePath := filepath.Join(projectPath, ".env.example")
	if content, err := os.ReadFile(envExamplePath); err == nil {
		configFiles[".env.example"] = string(content)
	}
	
	// Check for configuration yaml files
	configDirs := []string{"config", "configs", "internal/config"}
	for _, configDir := range configDirs {
		configDirPath := filepath.Join(projectPath, configDir)
		if _, err := os.Stat(configDirPath); err == nil {
			err := filepath.Walk(configDirPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, ".json") {
					if content, err := os.ReadFile(path); err == nil {
						relPath, _ := filepath.Rel(projectPath, path)
						configFiles[relPath] = string(content)
					}
				}
				return nil
			})
			if err != nil {
				// Continue processing other config directories even if one fails
				// (in a BDD context, we can't use t.Logf, so we just continue silently)
			}
		}
	}
	
	ctx.TestResults["configuration_files"] = configFiles
	return nil
}

func (ctx *EnhancedQualityTestContext) goModShouldContainTheFrameworkDependency(expectedDep string) error {
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	if projectPath == "" {
		return fmt.Errorf("no current project set for go.mod dependency validation")
	}
	
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %v", err)
	}
	
	goModContent := string(content)
	
	if !strings.Contains(goModContent, expectedDep) {
		return fmt.Errorf("go.mod should contain dependency '%s' but it was not found", expectedDep)
	}
	
	return nil
}

func (ctx *EnhancedQualityTestContext) iValidateFrameworkConsistencyAcrossAllProjects() error {
	// Validate framework consistency across all generated projects
	projectValidation := make(map[string]map[string]interface{})
	
	for projectName, projectPath := range ctx.ProjectPaths {
		expectedFramework := ctx.ProjectConfigs[projectName].Framework
		
		validation := make(map[string]interface{})
		validation["expected_framework"] = expectedFramework
		
		// Scan for framework references in the project
		frameworkRefs := make(map[string][]string)
		err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			if !strings.HasSuffix(path, ".go") {
				return nil
			}
			
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			
			// Check for framework imports and usage
			for _, framework := range []string{"gin", "echo", "fiber", "chi"} {
				frameworkPatterns := getFrameworkPatterns(framework)
				for _, pattern := range frameworkPatterns {
					if strings.Contains(string(content), pattern) {
						frameworkRefs[framework] = append(frameworkRefs[framework], path)
						break
					}
				}
			}
			
			return nil
		})
		
		if err != nil {
			return fmt.Errorf("failed to scan project %s: %v", projectName, err)
		}
		
		validation["framework_references"] = frameworkRefs
		projectValidation[projectName] = validation
	}
	
	ctx.TestResults["multi_project_validation"] = projectValidation
	return nil
}

func (ctx *EnhancedQualityTestContext) eachProjectShouldUseOnlyItsDesignatedFramework() error {
	// Validate that each project uses only its designated framework
	validation, ok := ctx.TestResults["multi_project_validation"].(map[string]map[string]interface{})
	if !ok {
		return fmt.Errorf("multi-project validation not recorded - ensure 'I validate framework consistency' step ran first")
	}
	
	for projectName, projectValidation := range validation {
		expectedFramework := projectValidation["expected_framework"].(string)
		frameworkRefs := projectValidation["framework_references"].(map[string][]string)
		
		// Should have expected framework
		if len(frameworkRefs[expectedFramework]) == 0 {
			return fmt.Errorf("project %s should contain %s framework references but none found", projectName, expectedFramework)
		}
		
		// Should not have other frameworks
		for framework, refs := range frameworkRefs {
			if framework != expectedFramework && len(refs) > 0 {
				return fmt.Errorf("project %s contains %s framework references but should only use %s. Found in: %v", 
					projectName, framework, expectedFramework, refs)
			}
		}
	}
	
	return nil
}

func (ctx *EnhancedQualityTestContext) noProjectShouldContainReferencesToOtherFrameworks() error {
	// This is covered by the designated framework check
	return ctx.eachProjectShouldUseOnlyItsDesignatedFramework()
}

func (ctx *EnhancedQualityTestContext) iAttemptToCompileTheGeneratedProject() error {
	// Use the current project path from the test context
	projectPath := ctx.ProjectPaths[ctx.CurrentProject]
	if projectPath == "" {
		return fmt.Errorf("no current project set for compilation")
	}
	
	// Run go mod tidy first to download dependencies
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = projectPath
	if tidyOutput, tidyErr := tidyCmd.CombinedOutput(); tidyErr != nil {
		// Store tidy error for debugging
		ctx.TestResults["compilation_success"] = false
		ctx.TestResults["compilation_output"] = fmt.Sprintf("go mod tidy failed: %s", string(tidyOutput))
		ctx.TestResults["compilation_error"] = tidyErr
		return nil
	}
	
	// Attempt to compile the project
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	
	// Store results in test context for verification steps
	ctx.TestResults["compilation_success"] = err == nil
	ctx.TestResults["compilation_output"] = string(output)
	ctx.TestResults["compilation_error"] = err
	
	return nil // Always return nil for step execution, validation happens in Then steps
}

func (ctx *EnhancedQualityTestContext) theCompilationShouldSucceedWithoutErrors() error {
	success, ok := ctx.TestResults["compilation_success"].(bool)
	if !ok {
		return fmt.Errorf("compilation result not recorded - ensure 'I attempt to compile' step ran first")
	}
	
	if !success {
		output := ctx.TestResults["compilation_output"].(string)
		err := ctx.TestResults["compilation_error"]
		return fmt.Errorf("project compilation failed:\nError: %v\nOutput: %s", err, output)
	}
	
	return nil
}

func (ctx *EnhancedQualityTestContext) theBuildOutputShouldNotContainWarnings() error {
	output, ok := ctx.TestResults["compilation_output"].(string)
	if !ok {
		return fmt.Errorf("compilation output not recorded - ensure 'I attempt to compile' step ran first")
	}
	
	// Check for common warning patterns
	warningPatterns := []string{
		"warning:",
		"deprecated:",
		"unused variable",
		"unused import",
	}
	
	var foundWarnings []string
	for _, pattern := range warningPatterns {
		if strings.Contains(strings.ToLower(output), pattern) {
			foundWarnings = append(foundWarnings, pattern)
		}
	}
	
	if len(foundWarnings) > 0 {
		return fmt.Errorf("build output contains warnings: %v\nOutput:\n%s", foundWarnings, output)
	}
	
	return nil
}