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

// TestComprehensiveBlueprintQualityATDD validates comprehensive blueprint quality acceptance criteria
// This test ensures all blueprints meet production-ready standards for code generation quality
func TestComprehensiveBlueprintQualityATDD(t *testing.T) {
	// High-priority blueprints that must work perfectly (updated with realistic file counts)
	priorityBlueprints := map[string]BlueprintConfig{
		"cli-simple": {
			Type:           "cli",
			Complexity:     "simple",
			ExpectedFiles:  8,
			MaxFiles:       12,
			MinFiles:       7,
			Description:    "Simple CLI with minimal structure",
		},
		"cli-standard": {
			Type:           "cli",
			Complexity:     "standard", 
			ExpectedFiles:  25,
			MaxFiles:       35,
			MinFiles:       20,
			Description:    "Full-featured CLI application",
		},
		"web-api-standard": {
			Type:           "web-api",
			Architecture:   "standard",
			ExpectedFiles:  35,
			MaxFiles:       45,
			MinFiles:       30,
			Description:    "Standard web API with REST endpoints",
		},
		"lambda-standard": {
			Type:           "lambda",
			Architecture:   "standard",
			ExpectedFiles:  15,
			MaxFiles:       20,
			MinFiles:       10,
			Description:    "AWS Lambda serverless function",
		},
		"library-standard": {
			Type:           "library",
			Architecture:   "standard",
			ExpectedFiles:  15,
			MaxFiles:       20,
			MinFiles:       12,
			Description:    "Go library package",
		},
	}

	// Logger types to test across all blueprints
	loggerTypes := []string{"slog", "zap", "logrus", "zerolog"}

	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	projectRoot := filepath.Join(originalDir, "..", "..", "..")
	
	defer func() { _ = os.Chdir(originalDir) }()
	_ = os.Chdir(tmpDir)

	// Build go-starter once for all tests
	buildGoStarter(t, tmpDir, projectRoot)

	t.Run("blueprint_generation_compilation_quality", func(t *testing.T) {
		// GIVEN: High-priority blueprints with different configurations
		// WHEN: Generating projects with each blueprint
		// THEN: All generated projects must compile successfully
		
		for blueprintName, config := range priorityBlueprints {
			t.Run(blueprintName, func(t *testing.T) {
				projectName := fmt.Sprintf("test-%s", blueprintName)
				
				// Generate project
				generateProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Verify compilation
				validateProjectCompilation(t, projectDir, blueprintName)
				
				// Verify file count expectations
				validateFileCount(t, projectDir, config)
				
				// Verify basic structure
				validateBasicStructure(t, projectDir, config)
			})
		}
	})

	t.Run("logger_integration_across_blueprints", func(t *testing.T) {
		// GIVEN: Multiple logger types (slog, zap, logrus, zerolog)
		// WHEN: Generating projects with each logger type
		// THEN: All logger integrations must work correctly and compile
		
		for _, logger := range loggerTypes {
			for blueprintName, config := range priorityBlueprints {
				// Skip certain combinations that don't support loggers
				if shouldSkipLoggerTest(blueprintName, logger) {
					continue
				}
				
				t.Run(fmt.Sprintf("%s_with_%s_logger", blueprintName, logger), func(t *testing.T) {
					projectName := fmt.Sprintf("test-%s-%s", blueprintName, logger)
					configWithLogger := config
					configWithLogger.Logger = logger
					
					// Generate project with specific logger
					generateProject(t, tmpDir, projectName, configWithLogger)
					projectDir := filepath.Join(tmpDir, projectName)
					
					// Validate logger integration
					validateLoggerIntegration(t, projectDir, logger, blueprintName)
					
					// Ensure project compiles with logger
					validateProjectCompilation(t, projectDir, fmt.Sprintf("%s-%s", blueprintName, logger))
				})
			}
		}
	})

	t.Run("template_variable_resolution", func(t *testing.T) {
		// GIVEN: Templates with various variables
		// WHEN: Generating projects
		// THEN: All template variables must resolve correctly with no {{.Variable}} remnants
		
		for blueprintName, config := range priorityBlueprints {
			t.Run(blueprintName+"_variable_resolution", func(t *testing.T) {
				projectName := fmt.Sprintf("var-test-%s", blueprintName)
				
				generateProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Check for unresolved template variables
				validateTemplateVariableResolution(t, projectDir, blueprintName)
			})
		}
	})

	t.Run("go_module_dependency_correctness", func(t *testing.T) {
		// GIVEN: Generated projects with dependencies
		// WHEN: Running go mod tidy
		// THEN: All dependencies must be valid and accessible
		
		for blueprintName, config := range priorityBlueprints {
			t.Run(blueprintName+"_dependencies", func(t *testing.T) {
				projectName := fmt.Sprintf("deps-test-%s", blueprintName)
				
				generateProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Validate go.mod and dependencies
				validateGoDependencies(t, projectDir, blueprintName)
			})
		}
	})

	t.Run("progressive_complexity_validation", func(t *testing.T) {
		// GIVEN: CLI blueprints with different complexity levels
		// WHEN: Generating simple vs standard CLI
		// THEN: Simple should have 60-70% fewer files than standard
		
		// Generate simple CLI
		generateProject(t, tmpDir, "complexity-simple", priorityBlueprints["cli-simple"])
		simpleDir := filepath.Join(tmpDir, "complexity-simple")
		simpleCount := countProjectFiles(t, simpleDir)
		
		// Generate standard CLI  
		generateProject(t, tmpDir, "complexity-standard", priorityBlueprints["cli-standard"])
		standardDir := filepath.Join(tmpDir, "complexity-standard")
		standardCount := countProjectFiles(t, standardDir)
		
		// Validate file count reduction
		validateComplexityReduction(t, simpleCount, standardCount)
		
		// Validate both compile successfully
		validateProjectCompilation(t, simpleDir, "cli-simple-complexity")
		validateProjectCompilation(t, standardDir, "cli-standard-complexity")
	})

	t.Run("conditional_file_generation_logic", func(t *testing.T) {
		// GIVEN: Blueprints with conditional file generation
		// WHEN: Generating with different configurations
		// THEN: Conditional logic must work correctly
		
		// Test web-api with and without database
		t.Run("web_api_database_conditional", func(t *testing.T) {
			// Generate without database
			configNoDB := priorityBlueprints["web-api-standard"]
			generateProject(t, tmpDir, "webapi-no-db", configNoDB)
			noDatabaseFiles := countProjectFiles(t, filepath.Join(tmpDir, "webapi-no-db"))
			
			// Generate with database
			configWithDB := priorityBlueprints["web-api-standard"]
			configWithDB.DatabaseDriver = "postgres"
			generateProject(t, tmpDir, "webapi-with-db", configWithDB)
			withDatabaseFiles := countProjectFiles(t, filepath.Join(tmpDir, "webapi-with-db"))
			
			// With database should have more files
			assert.Greater(t, withDatabaseFiles, noDatabaseFiles, "Database configuration should generate additional files")
		})
	})

	t.Run("code_quality_and_go_best_practices", func(t *testing.T) {
		// GIVEN: Generated projects  
		// WHEN: Running go fmt, go vet, and code analysis
		// THEN: Generated code must follow Go best practices
		
		for blueprintName, config := range priorityBlueprints {
			t.Run(blueprintName+"_code_quality", func(t *testing.T) {
				projectName := fmt.Sprintf("quality-%s", blueprintName)
				
				generateProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Initialize go modules
				initializeGoModules(t, projectDir)
				
				// Validate code formatting
				validateGoFormatting(t, projectDir, blueprintName)
				
				// Validate code with go vet
				validateGoVet(t, projectDir, blueprintName)
				
				// Validate imports and package structure
				validatePackageStructure(t, projectDir, blueprintName)
			})
		}
	})

	t.Run("simplified_logger_architecture", func(t *testing.T) {
		// GIVEN: The new simplified logger system
		// WHEN: Examining generated logger code
		// THEN: Should have minimal complexity (60-90% reduction)
		
		for _, logger := range loggerTypes {
			t.Run("simplified_"+logger+"_logger", func(t *testing.T) {
				projectName := fmt.Sprintf("logger-simple-%s", logger)
				config := priorityBlueprints["cli-standard"]
				config.Logger = logger
				
				generateProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Validate simplified logger implementation
				validateSimplifiedLogger(t, projectDir, logger)
			})
		}
	})
}

// BlueprintConfig represents test configuration for a blueprint
type BlueprintConfig struct {
	Type           string
	Architecture   string
	Complexity     string
	Logger         string
	DatabaseDriver string
	ExpectedFiles  int
	MaxFiles       int
	MinFiles       int
	Description    string
}


// generateProject generates a project with the given configuration
func generateProject(t *testing.T, tmpDir, projectName string, config BlueprintConfig) {
	t.Helper()
	
	args := []string{"new", projectName}
	
	// Add type
	if config.Type != "" {
		args = append(args, "--type="+config.Type)
	}
	
	// Add complexity
	if config.Complexity != "" {
		args = append(args, "--complexity="+config.Complexity)
	}
	
	// Add architecture  
	if config.Architecture != "" {
		args = append(args, "--architecture="+config.Architecture)
	}
	
	// Add logger
	logger := config.Logger
	if logger == "" {
		logger = "slog" // Default logger
	}
	args = append(args, "--logger="+logger)
	
	// Add database driver if specified
	if config.DatabaseDriver != "" {
		args = append(args, "--database-driver="+config.DatabaseDriver)
	}
	
	// Add module path and other flags
	args = append(args, 
		"--module=github.com/test/"+projectName,
		"--no-git",
	)
	
	goStarterPath := filepath.Join(tmpDir, "go-starter")
	cmd := exec.Command(goStarterPath, args...)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Logf("Generate command: %s %s", goStarterPath, strings.Join(args, " "))
		t.Logf("Generate output: %s", string(output))
	}
	require.NoError(t, err, "Project generation should succeed for %s", projectName)
	t.Logf("Successfully generated project: %s", projectName)
}


// validateFileCount ensures the generated project has the expected number of files
func validateFileCount(t *testing.T, projectDir string, config BlueprintConfig) {
	t.Helper()
	
	fileCount := countProjectFiles(t, projectDir)
	
	assert.GreaterOrEqual(t, fileCount, config.MinFiles, "Project should have at least %d files", config.MinFiles)
	assert.LessOrEqual(t, fileCount, config.MaxFiles, "Project should have at most %d files", config.MaxFiles)
	
	t.Logf("✓ File count validation passed: %d files (expected %d±%d)", 
		fileCount, config.ExpectedFiles, config.MaxFiles-config.ExpectedFiles)
}

// validateBasicStructure ensures basic project structure exists
func validateBasicStructure(t *testing.T, projectDir string, config BlueprintConfig) {
	t.Helper()
	
	// Common files that should exist in most projects
	commonFiles := []string{
		"go.mod",
		"README.md",
	}
	
	for _, file := range commonFiles {
		filePath := filepath.Join(projectDir, file)
		assert.FileExists(t, filePath, "Basic file %s should exist", file)
	}
	
	// Type-specific validations
	switch config.Type {
	case "cli":
		assert.DirExists(t, filepath.Join(projectDir, "cmd"), "CLI should have cmd/ directory")
		assert.FileExists(t, filepath.Join(projectDir, "main.go"), "CLI should have main.go")
		if config.Complexity != "simple" {
			assert.DirExists(t, filepath.Join(projectDir, "internal"), "Standard CLI should have internal/ directory")
		}
	case "web-api":
		assert.DirExists(t, filepath.Join(projectDir, "internal"), "Web API should have internal/ directory")
		assert.DirExists(t, filepath.Join(projectDir, "cmd"), "Web API should have cmd/ directory")
		// Web API has main.go in cmd/server/main.go
		assert.FileExists(t, filepath.Join(projectDir, "cmd", "server", "main.go"), "Web API should have cmd/server/main.go")
	case "library":
		// Library has project-name.go instead of library.go
		libFiles := []string{}
		err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") && path != filepath.Join(projectDir, "doc.go") {
				libFiles = append(libFiles, path)
			}
			return nil
		})
		require.NoError(t, err, "Should be able to find library files")
		assert.NotEmpty(t, libFiles, "Library should have at least one main library file")
	case "lambda":
		assert.FileExists(t, filepath.Join(projectDir, "main.go"), "Lambda should have main.go")
	}
	
	t.Logf("✓ Basic structure validation passed for %s", config.Type)
}

// validateLoggerIntegration ensures logger integration works correctly  
func validateLoggerIntegration(t *testing.T, projectDir, logger, blueprintName string) {
	t.Helper()
	
	// Find logger files
	loggerFiles := findLoggerFiles(t, projectDir)
	require.NotEmpty(t, loggerFiles, "Should have logger files for %s", blueprintName)
	
	for _, loggerFile := range loggerFiles {
		content, err := os.ReadFile(loggerFile)
		require.NoError(t, err, "Should be able to read logger file")
		contentStr := string(content)
		
		// Validate logger-specific imports and usage
		switch logger {
		case "slog":
			assert.Contains(t, contentStr, "log/slog", "slog logger should import log/slog")
		case "zap":
			assert.Contains(t, contentStr, "go.uber.org/zap", "zap logger should import zap")
		case "logrus":
			assert.Contains(t, contentStr, "github.com/sirupsen/logrus", "logrus logger should import logrus")
		case "zerolog":
			assert.Contains(t, contentStr, "github.com/rs/zerolog", "zerolog logger should import zerolog")
		}
	}
	
	t.Logf("✓ Logger integration validated for %s with %s", blueprintName, logger)
}

// validateTemplateVariableResolution ensures no unresolved template variables remain
func validateTemplateVariableResolution(t *testing.T, projectDir, blueprintName string) {
	t.Helper()
	
	var unresolvedVars []string
	
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip binary files and directories
		if info.IsDir() || isBinaryFile(path) {
			return nil
		}
		
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		
		contentStr := string(content)
		
		// Look for unresolved template variables like {{.Variable}}
		if strings.Contains(contentStr, "{{.") || strings.Contains(contentStr, "}}") {
			lines := strings.Split(contentStr, "\n")
			for i, line := range lines {
				if strings.Contains(line, "{{.") || strings.Contains(line, "}}") {
					unresolvedVars = append(unresolvedVars, 
						fmt.Sprintf("%s:%d: %s", path, i+1, strings.TrimSpace(line)))
				}
			}
		}
		
		return nil
	})
	
	require.NoError(t, err, "Should be able to walk project directory")
	
	if len(unresolvedVars) > 0 {
		t.Logf("Unresolved template variables found in %s:", blueprintName)
		for _, v := range unresolvedVars {
			t.Logf("  %s", v)
		}
	}
	
	assert.Empty(t, unresolvedVars, "Should have no unresolved template variables in %s", blueprintName)
	t.Logf("✓ Template variable resolution validated for %s", blueprintName)
}

// validateGoDependencies ensures go.mod and dependencies are correct
func validateGoDependencies(t *testing.T, projectDir, blueprintName string) {
	t.Helper()
	
	// Check go.mod exists
	goModPath := filepath.Join(projectDir, "go.mod")
	assert.FileExists(t, goModPath, "go.mod should exist")
	
	// Initialize modules to verify dependencies
	initializeGoModules(t, projectDir)
	
	// Run go mod verify
	verifyCmd := exec.Command("go", "mod", "verify")
	verifyCmd.Dir = projectDir
	output, err := verifyCmd.CombinedOutput()
	
	require.NoError(t, err, "go mod verify should succeed for %s: %s", blueprintName, string(output))
	
	t.Logf("✓ Go dependencies validated for %s", blueprintName)
}

// validateComplexityReduction ensures simple CLI has significantly fewer files than standard
func validateComplexityReduction(t *testing.T, simpleCount, standardCount int) {
	t.Helper()
	
	reductionPercentage := float64(standardCount-simpleCount) / float64(standardCount) * 100
	
	t.Logf("File count: simple=%d, standard=%d, reduction=%.1f%%", 
		simpleCount, standardCount, reductionPercentage)
	
	assert.GreaterOrEqual(t, reductionPercentage, 60.0, "Simple CLI should have 60%+ fewer files than standard")
	assert.LessOrEqual(t, reductionPercentage, 80.0, "Reduction should be realistic (not over 80%)")
	
	t.Logf("✓ Progressive complexity validation passed: %.1f%% reduction", reductionPercentage)
}

// validateGoFormatting ensures generated code is properly formatted
func validateGoFormatting(t *testing.T, projectDir, blueprintName string) {
	t.Helper()
	
	fmtCmd := exec.Command("go", "fmt", "./...")
	fmtCmd.Dir = projectDir
	output, err := fmtCmd.CombinedOutput()
	
	require.NoError(t, err, "go fmt should succeed for %s", blueprintName)
	assert.Empty(t, string(output), "Generated code should be properly formatted for %s", blueprintName)
	
	t.Logf("✓ Go formatting validated for %s", blueprintName)
}

// validateGoVet ensures generated code passes go vet
func validateGoVet(t *testing.T, projectDir, blueprintName string) {
	t.Helper()
	
	vetCmd := exec.Command("go", "vet", "./...")
	vetCmd.Dir = projectDir
	output, err := vetCmd.CombinedOutput()
	
	require.NoError(t, err, "go vet should pass for %s: %s", blueprintName, string(output))
	
	t.Logf("✓ Go vet validation passed for %s", blueprintName)
}

// validatePackageStructure ensures proper Go package structure
func validatePackageStructure(t *testing.T, projectDir, blueprintName string) {
	t.Helper()
	
	// Find all .go files and validate package declarations
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		
		contentStr := string(content)
		lines := strings.Split(contentStr, "\n")
		
		// Find package declaration
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "package ") {
				packageName := strings.TrimPrefix(line, "package ")
				
				// Validate package name is valid Go identifier
				assert.Regexp(t, `^[a-z][a-z0-9_]*$`, packageName, 
					"Package name should be valid Go identifier in %s", path)
				break
			}
		}
		
		return nil
	})
	
	require.NoError(t, err, "Should be able to validate package structure")
	t.Logf("✓ Package structure validated for %s", blueprintName)
}

// validateSimplifiedLogger ensures simplified logger architecture
func validateSimplifiedLogger(t *testing.T, projectDir, logger string) {
	t.Helper()
	
	loggerFiles := findLoggerFiles(t, projectDir)
	
	for _, loggerFile := range loggerFiles {
		content, err := os.ReadFile(loggerFile)
		require.NoError(t, err)
		
		lines := strings.Split(string(content), "\n")
		nonEmptyLines := 0
		
		for _, line := range lines {
			if strings.TrimSpace(line) != "" && !strings.HasPrefix(strings.TrimSpace(line), "//") {
				nonEmptyLines++
			}
		}
		
		// Simplified logger should be concise (aim for <50 meaningful lines)
		assert.LessOrEqual(t, nonEmptyLines, 100, "Simplified logger should be concise (was %d lines)", nonEmptyLines)
		
		t.Logf("✓ Simplified logger validated: %s has %d lines", logger, nonEmptyLines)
	}
}

// Helper functions

func shouldSkipLoggerTest(blueprintName, logger string) bool {
	// Some blueprints might not support all logger types
	// Add logic here if needed
	return false
}

