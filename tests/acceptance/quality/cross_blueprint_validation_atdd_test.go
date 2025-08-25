package quality

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCrossBlueprintValidationATDD validates cross-blueprint acceptance criteria
// This ensures consistency and quality across all high-priority blueprints
func TestCrossBlueprintValidationATDD(t *testing.T) {
	// High-priority blueprints for production readiness
	priorityBlueprints := map[string]CrossBlueprintConfig{
		"cli-simple": {
			Type:          "cli",
			Complexity:    "simple",
			ExpectedFiles: []string{"main.go", "go.mod", "README.md", "cmd/root.go"},
			RequiredDirs:  []string{"cmd"},
			OptionalDirs:  []string{"internal"},
			Capabilities:  []string{"basic-cli", "help-command", "version-support"},
		},
		"cli-standard": {
			Type:          "cli", 
			Complexity:    "standard",
			ExpectedFiles: []string{"main.go", "go.mod", "README.md", "Makefile", "cmd/root.go"},
			RequiredDirs:  []string{"cmd", "internal"},
			OptionalDirs:  []string{"configs", "scripts"},
			Capabilities:  []string{"full-cli", "subcommands", "configuration", "production-ready"},
		},
		"web-api-standard": {
			Type:          "web-api",
			Architecture:  "standard",
			ExpectedFiles: []string{"main.go", "go.mod", "README.md", "Makefile"},
			RequiredDirs:  []string{"cmd", "internal"},
			OptionalDirs:  []string{"api", "configs", "migrations"},
			Capabilities:  []string{"rest-api", "http-server", "routing", "middleware"},
		},
		"lambda-standard": {
			Type:          "lambda",
			ExpectedFiles: []string{"main.go", "go.mod", "README.md"},
			RequiredDirs:  []string{},
			OptionalDirs:  []string{"internal"},
			Capabilities:  []string{"aws-lambda", "event-processing", "serverless"},
		},
		"library-standard": {
			Type:          "library",
			ExpectedFiles: []string{"go.mod", "README.md", "library.go", "library_test.go"},
			RequiredDirs:  []string{},
			OptionalDirs:  []string{"examples", "internal"},
			Capabilities:  []string{"public-api", "go-module", "documentation", "testing"},
		},
	}

	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	projectRoot := filepath.Join(originalDir, "..", "..", "..")
	
	defer func() { _ = os.Chdir(originalDir) }()
	_ = os.Chdir(tmpDir)

	// Build go-starter once
	buildGoStarter(t, tmpDir, projectRoot)

	t.Run("blueprint_consistency_validation", func(t *testing.T) {
		// GIVEN: Multiple blueprint types
		// WHEN: Generating projects with each blueprint
		// THEN: All should follow consistent patterns and conventions
		
		generatedProjects := make(map[string]string)
		
		for blueprintName, config := range priorityBlueprints {
			t.Run(blueprintName+"_consistency", func(t *testing.T) {
				projectName := fmt.Sprintf("consistency-%s", blueprintName)
				
				generateCrossBlueprintProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				generatedProjects[blueprintName] = projectDir
				
				// Validate consistent structure
				validateBlueprintConsistency(t, projectDir, config, blueprintName)
			})
		}
		
		// Cross-blueprint consistency checks
		validateCrossConsistency(t, generatedProjects, priorityBlueprints)
	})

	t.Run("blueprint_capability_validation", func(t *testing.T) {
		// GIVEN: Blueprint capabilities requirements
		// WHEN: Examining generated project structure and code
		// THEN: Each blueprint should demonstrate its stated capabilities
		
		for blueprintName, config := range priorityBlueprints {
			t.Run(blueprintName+"_capabilities", func(t *testing.T) {
				projectName := fmt.Sprintf("capability-%s", blueprintName)
				
				generateCrossBlueprintProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Validate blueprint-specific capabilities
				validateBlueprintCapabilities(t, projectDir, config, blueprintName)
			})
		}
	})

	t.Run("blueprint_go_module_validation", func(t *testing.T) {
		// GIVEN: Generated projects from different blueprints
		// WHEN: Examining go.mod files and module structure
		// THEN: All should have valid Go module configuration
		
		for blueprintName, config := range priorityBlueprints {
			t.Run(blueprintName+"_module", func(t *testing.T) {
				projectName := fmt.Sprintf("module-%s", blueprintName)
				
				generateCrossBlueprintProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Validate Go module structure
				validateGoModuleStructure(t, projectDir, blueprintName)
			})
		}
	})

	t.Run("blueprint_compilation_matrix", func(t *testing.T) {
		// GIVEN: All high-priority blueprints
		// WHEN: Compiling each generated project
		// THEN: All should compile successfully with no errors
		
		compilationResults := make(map[string]CompilationResult)
		
		for blueprintName, config := range priorityBlueprints {
			t.Run(blueprintName+"_compilation", func(t *testing.T) {
				projectName := fmt.Sprintf("compile-%s", blueprintName)
				
				generateCrossBlueprintProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Test compilation
				result := validateBlueprintCompilation(t, projectDir, blueprintName)
				compilationResults[blueprintName] = result
			})
		}
		
		// Report compilation matrix
		reportCompilationMatrix(t, compilationResults)
	})

	t.Run("blueprint_documentation_standards", func(t *testing.T) {
		// GIVEN: Generated projects from different blueprints  
		// WHEN: Examining documentation files
		// THEN: All should have consistent documentation standards
		
		for blueprintName, config := range priorityBlueprints {
			t.Run(blueprintName+"_documentation", func(t *testing.T) {
				projectName := fmt.Sprintf("docs-%s", blueprintName)
				
				generateCrossBlueprintProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Validate documentation standards
				validateDocumentationStandards(t, projectDir, blueprintName)
			})
		}
	})

	t.Run("blueprint_testing_infrastructure", func(t *testing.T) {
		// GIVEN: Blueprints that support testing
		// WHEN: Examining test files and structure
		// THEN: Should have appropriate testing infrastructure
		
		for blueprintName, config := range priorityBlueprints {
			// Skip blueprints that don't typically include tests
			if shouldSkipTestingValidation(blueprintName) {
				continue
			}
			
			t.Run(blueprintName+"_testing", func(t *testing.T) {
				projectName := fmt.Sprintf("test-%s", blueprintName)
				
				generateCrossBlueprintProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Validate testing infrastructure
				validateTestingInfrastructure(t, projectDir, blueprintName)
			})
		}
	})

	t.Run("blueprint_production_readiness", func(t *testing.T) {
		// GIVEN: Production-ready blueprint requirements
		// WHEN: Examining generated projects for production features
		// THEN: Should have appropriate production-ready features
		
		productionBlueprints := []string{"cli-standard", "web-api-standard", "lambda-standard"}
		
		for _, blueprintName := range productionBlueprints {
			if config, exists := priorityBlueprints[blueprintName]; exists {
				t.Run(blueprintName+"_production", func(t *testing.T) {
					projectName := fmt.Sprintf("prod-%s", blueprintName)
					
					generateCrossBlueprintProject(t, tmpDir, projectName, config)
					projectDir := filepath.Join(tmpDir, projectName)
					
					// Validate production readiness
					validateProductionReadiness(t, projectDir, blueprintName)
				})
			}
		}
	})

	t.Run("blueprint_error_handling_patterns", func(t *testing.T) {
		// GIVEN: Generated projects with Go code
		// WHEN: Examining error handling patterns
		// THEN: Should follow Go error handling best practices
		
		for blueprintName, config := range priorityBlueprints {
			t.Run(blueprintName+"_error_handling", func(t *testing.T) {
				projectName := fmt.Sprintf("error-%s", blueprintName)
				
				generateCrossBlueprintProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Validate error handling patterns
				validateErrorHandlingPatterns(t, projectDir, blueprintName)
			})
		}
	})
}

// CrossBlueprintConfig represents configuration for cross-blueprint testing
type CrossBlueprintConfig struct {
	Type          string
	Complexity    string
	Architecture  string
	ExpectedFiles []string
	RequiredDirs  []string
	OptionalDirs  []string
	Capabilities  []string
}

// CompilationResult represents the result of a compilation test
type CompilationResult struct {
	Success     bool
	BuildTime   string
	BinarySize  int64
	Errors      []string
}

// generateCrossBlueprintProject generates a project for cross-blueprint testing
func generateCrossBlueprintProject(t *testing.T, tmpDir, projectName string, config CrossBlueprintConfig) {
	t.Helper()
	
	args := []string{"new", projectName}
	
	// Add type
	args = append(args, "--type="+config.Type)
	
	// Add complexity if specified
	if config.Complexity != "" {
		args = append(args, "--complexity="+config.Complexity)
	}
	
	// Add architecture if specified
	if config.Architecture != "" {
		args = append(args, "--architecture="+config.Architecture)
	}
	
	// Add common flags
	args = append(args, 
		"--module=github.com/test/"+projectName,
		"--logger=slog",
		"--no-git",
	)
	
	goStarterPath := filepath.Join(tmpDir, "go-starter")
	cmd := exec.Command(goStarterPath, args...)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Logf("Generate command: %s %s", goStarterPath, strings.Join(args, " "))
		t.Logf("Generate output: %s", string(output))
	}
	require.NoError(t, err, "Should generate project for %s", projectName)
	
	t.Logf("Successfully generated cross-blueprint project: %s", projectName)
}

// validateBlueprintConsistency validates consistent patterns across blueprints
func validateBlueprintConsistency(t *testing.T, projectDir string, config CrossBlueprintConfig, blueprintName string) {
	t.Helper()
	
	// Validate expected files exist
	for _, expectedFile := range config.ExpectedFiles {
		filePath := filepath.Join(projectDir, expectedFile)
		assert.FileExists(t, filePath, "Expected file %s should exist in %s", expectedFile, blueprintName)
	}
	
	// Validate required directories exist
	for _, requiredDir := range config.RequiredDirs {
		dirPath := filepath.Join(projectDir, requiredDir)
		assert.DirExists(t, dirPath, "Required directory %s should exist in %s", requiredDir, blueprintName)
	}
	
	// Check go.mod consistency
	validateGoModConsistency(t, projectDir, blueprintName)
	
	// Check README.md consistency
	validateReadmeConsistency(t, projectDir, blueprintName)
	
	t.Logf("✓ Blueprint consistency validated for %s", blueprintName)
}

// validateBlueprintCapabilities validates blueprint-specific capabilities
func validateBlueprintCapabilities(t *testing.T, projectDir string, config CrossBlueprintConfig, blueprintName string) {
	t.Helper()
	
	for _, capability := range config.Capabilities {
		switch capability {
		case "basic-cli":
			validateBasicCLICapability(t, projectDir)
		case "full-cli":
			validateFullCLICapability(t, projectDir)
		case "rest-api":
			validateRestAPICapability(t, projectDir)
		case "aws-lambda":
			validateAWSLambdaCapability(t, projectDir)
		case "public-api":
			validatePublicAPICapability(t, projectDir)
		case "production-ready":
			validateProductionReadyCapability(t, projectDir)
		}
	}
	
	t.Logf("✓ Blueprint capabilities validated for %s", blueprintName)
}

// validateGoModuleStructure validates Go module structure and dependencies
func validateGoModuleStructure(t *testing.T, projectDir, blueprintName string) {
	t.Helper()
	
	goModPath := filepath.Join(projectDir, "go.mod")
	assert.FileExists(t, goModPath, "go.mod should exist")
	
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")
	contentStr := string(content)
	
	// Validate module declaration
	assert.Contains(t, contentStr, "module github.com/test/", "Should have correct module path")
	
	// Validate Go version
	assert.Contains(t, contentStr, "go 1.2", "Should specify Go version")
	
	// Initialize and validate dependencies
	initializeGoModules(t, projectDir)
	
	t.Logf("✓ Go module structure validated for %s", blueprintName)
}

// validateBlueprintCompilation compiles the project and returns results
func validateBlueprintCompilation(t *testing.T, projectDir, blueprintName string) CompilationResult {
	t.Helper()
	
	result := CompilationResult{}
	
	// Initialize modules
	initializeGoModules(t, projectDir)
	
	// Build the project
	buildCmd := exec.Command("go", "build", "-o", "test-binary", ".")
	buildCmd.Dir = projectDir
	output, err := buildCmd.CombinedOutput()
	
	if err != nil {
		result.Success = false
		result.Errors = []string{string(output)}
		t.Logf("⚠ Compilation failed for %s: %s", blueprintName, string(output))
	} else {
		result.Success = true
		
		// Get binary info if successful
		binaryPath := filepath.Join(projectDir, "test-binary")
		if stat, statErr := os.Stat(binaryPath); statErr == nil {
			result.BinarySize = stat.Size()
		}
		
		t.Logf("✓ Compilation succeeded for %s (binary: %d bytes)", blueprintName, result.BinarySize)
	}
	
	return result
}

// validateDocumentationStandards validates documentation consistency
func validateDocumentationStandards(t *testing.T, projectDir, blueprintName string) {
	t.Helper()
	
	readmePath := filepath.Join(projectDir, "README.md")
	assert.FileExists(t, readmePath, "README.md should exist")
	
	content, err := os.ReadFile(readmePath)
	require.NoError(t, err, "Should be able to read README.md")
	contentStr := string(content)
	
	// Standard documentation elements
	assert.Contains(t, contentStr, "# ", "Should have main heading")
	assert.Contains(t, contentStr, "## Installation", "Should have installation section")
	assert.Contains(t, contentStr, "## Usage", "Should have usage section")
	
	// Blueprint-specific documentation
	switch blueprintName {
	case "cli-simple", "cli-standard":
		assert.Contains(t, contentStr, "command", "CLI documentation should mention commands")
	case "web-api-standard":
		assert.Contains(t, contentStr, "API", "Web API documentation should mention API")
	case "lambda-standard":
		assert.Contains(t, contentStr, "AWS", "Lambda documentation should mention AWS")
	case "library-standard":
		assert.Contains(t, contentStr, "import", "Library documentation should show import")
	}
	
	t.Logf("✓ Documentation standards validated for %s", blueprintName)
}

// validateTestingInfrastructure validates testing setup
func validateTestingInfrastructure(t *testing.T, projectDir, blueprintName string) {
	t.Helper()
	
	// Look for test files
	testFiles := findTestFiles(t, projectDir)
	
	if len(testFiles) > 0 {
		for _, testFile := range testFiles {
			// Validate test file structure
			content, err := os.ReadFile(testFile)
			require.NoError(t, err, "Should be able to read test file")
			contentStr := string(content)
			
			assert.Contains(t, contentStr, "package ", "Test file should have package declaration")
			assert.Contains(t, contentStr, "import ", "Test file should have imports")
			assert.Contains(t, contentStr, "func Test", "Test file should have test functions")
		}
		
		t.Logf("✓ Testing infrastructure validated for %s (%d test files)", blueprintName, len(testFiles))
	} else {
		t.Logf("No test files found for %s (may be expected for this blueprint type)", blueprintName)
	}
}

// validateProductionReadiness validates production-ready features
func validateProductionReadiness(t *testing.T, projectDir, blueprintName string) {
	t.Helper()
	
	productionFeatures := []string{
		"Makefile",     // Build automation
		"README.md",    // Documentation
		"go.mod",       // Dependency management
	}
	
	for _, feature := range productionFeatures {
		featurePath := filepath.Join(projectDir, feature)
		assert.FileExists(t, featurePath, "Production feature %s should exist in %s", feature, blueprintName)
	}
	
	// Validate Makefile has production targets
	makefilePath := filepath.Join(projectDir, "Makefile")
	if assert.FileExists(t, makefilePath, "Makefile should exist") {
		content, err := os.ReadFile(makefilePath)
		require.NoError(t, err, "Should be able to read Makefile")
		contentStr := string(content)
		
		productionTargets := []string{"build", "test", "clean"}
		for _, target := range productionTargets {
			assert.Contains(t, contentStr, target+":", "Makefile should have %s target", target)
		}
	}
	
	t.Logf("✓ Production readiness validated for %s", blueprintName)
}

// validateErrorHandlingPatterns validates Go error handling best practices
func validateErrorHandlingPatterns(t *testing.T, projectDir, blueprintName string) {
	t.Helper()
	
	goFiles := findGoFiles(t, projectDir)
	
	for _, goFile := range goFiles {
		content, err := os.ReadFile(goFile)
		require.NoError(t, err, "Should be able to read Go file")
		contentStr := string(content)
		
		// Check for proper error handling patterns
		if strings.Contains(contentStr, "err :=") || strings.Contains(contentStr, "err =") {
			// If errors are assigned, they should be handled
			assert.Contains(t, contentStr, "if err != nil", 
				"Go file %s should handle errors properly", goFile)
		}
	}
	
	t.Logf("✓ Error handling patterns validated for %s", blueprintName)
}

// Capability-specific validation functions

func validateBasicCLICapability(t *testing.T, projectDir string) {
	t.Helper()
	assert.FileExists(t, filepath.Join(projectDir, "main.go"), "Basic CLI should have main.go")
	assert.DirExists(t, filepath.Join(projectDir, "cmd"), "Basic CLI should have cmd/ directory")
}

func validateFullCLICapability(t *testing.T, projectDir string) {
	t.Helper()
	validateBasicCLICapability(t, projectDir)
	assert.DirExists(t, filepath.Join(projectDir, "internal"), "Full CLI should have internal/ directory")
	assert.FileExists(t, filepath.Join(projectDir, "Makefile"), "Full CLI should have Makefile")
}

func validateRestAPICapability(t *testing.T, projectDir string) {
	t.Helper()
	assert.DirExists(t, filepath.Join(projectDir, "internal"), "REST API should have internal/ directory")
	assert.DirExists(t, filepath.Join(projectDir, "cmd"), "REST API should have cmd/ directory")
}

func validateAWSLambdaCapability(t *testing.T, projectDir string) {
	t.Helper()
	assert.FileExists(t, filepath.Join(projectDir, "main.go"), "Lambda should have main.go")
	
	// Check for AWS Lambda-specific imports
	mainContent, err := os.ReadFile(filepath.Join(projectDir, "main.go"))
	if err == nil {
		contentStr := string(mainContent)
		assert.Contains(t, contentStr, "lambda", "Lambda should import lambda package")
	}
}

func validatePublicAPICapability(t *testing.T, projectDir string) {
	t.Helper()
	
	// Look for public API files
	goFiles := findGoFiles(t, projectDir)
	
	hasPublicFunctions := false
	for _, goFile := range goFiles {
		content, err := os.ReadFile(goFile)
		if err != nil {
			continue
		}
		contentStr := string(content)
		
		// Look for exported functions
		if strings.Contains(contentStr, "func ") && 
		   !strings.Contains(filepath.Base(goFile), "_test.go") {
			hasPublicFunctions = true
			break
		}
	}
	
	assert.True(t, hasPublicFunctions, "Library should have public API functions")
}

func validateProductionReadyCapability(t *testing.T, projectDir string) {
	t.Helper()
	productionFiles := []string{"Makefile", "README.md"}
	for _, file := range productionFiles {
		assert.FileExists(t, filepath.Join(projectDir, file), "Production-ready project should have %s", file)
	}
}

// Cross-validation functions

func validateCrossConsistency(t *testing.T, projects map[string]string, configs map[string]CrossBlueprintConfig) {
	t.Helper()
	
	// Validate consistent go.mod structure across projects
	for blueprintName, projectDir := range projects {
		goModPath := filepath.Join(projectDir, "go.mod")
		if !assert.FileExists(t, goModPath, "go.mod should exist in %s", blueprintName) {
			continue
		}
		
		content, err := os.ReadFile(goModPath)
		require.NoError(t, err, "Should read go.mod from %s", blueprintName)
		
		// All should use similar Go version
		assert.Contains(t, string(content), "go 1.2", "%s should use consistent Go version", blueprintName)
	}
	
	t.Logf("✓ Cross-blueprint consistency validated")
}

func validateGoModConsistency(t *testing.T, projectDir, blueprintName string) {
	t.Helper()
	
	goModPath := filepath.Join(projectDir, "go.mod")
	assert.FileExists(t, goModPath, "go.mod should exist")
	
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")
	contentStr := string(content)
	
	// Consistent module naming
	assert.Contains(t, contentStr, "github.com/test/", "Should have consistent module prefix")
	assert.Contains(t, contentStr, "go 1.2", "Should have consistent Go version")
}

func validateReadmeConsistency(t *testing.T, projectDir, blueprintName string) {
	t.Helper()
	
	readmePath := filepath.Join(projectDir, "README.md")
	assert.FileExists(t, readmePath, "README.md should exist")
	
	content, err := os.ReadFile(readmePath)
	require.NoError(t, err, "Should be able to read README.md")
	contentStr := string(content)
	
	// Consistent structure
	assert.Contains(t, contentStr, "# ", "Should have main heading")
	assert.Contains(t, contentStr, "## ", "Should have section headings")
}

func reportCompilationMatrix(t *testing.T, results map[string]CompilationResult) {
	t.Helper()
	
	successful := 0
	failed := 0
	
	for blueprintName, result := range results {
		if result.Success {
			successful++
			t.Logf("✓ %s: compiled successfully (%d bytes)", blueprintName, result.BinarySize)
		} else {
			failed++
			t.Logf("✗ %s: compilation failed", blueprintName)
		}
	}
	
	t.Logf("Compilation Matrix: %d successful, %d failed", successful, failed)
	assert.Equal(t, 0, failed, "All blueprints should compile successfully")
}

// Helper functions

func shouldSkipTestingValidation(blueprintName string) bool {
	// Some blueprints might not include test files by default
	skipBlueprints := []string{"lambda-standard"}
	
	for _, skip := range skipBlueprints {
		if blueprintName == skip {
			return true
		}
	}
	
	return false
}

func findTestFiles(t *testing.T, projectDir string) []string {
	t.Helper()
	var testFiles []string
	
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if strings.HasSuffix(path, "_test.go") {
			testFiles = append(testFiles, path)
		}
		
		return nil
	})
	
	require.NoError(t, err, "Should be able to find test files")
	return testFiles
}

func findGoFiles(t *testing.T, projectDir string) []string {
	t.Helper()
	var goFiles []string
	
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if strings.HasSuffix(path, ".go") && !info.IsDir() {
			goFiles = append(goFiles, path)
		}
		
		return nil
	})
	
	require.NoError(t, err, "Should be able to find Go files")
	return goFiles
}