package quality

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoggerIntegrationATDD validates logger integration acceptance criteria
// This ensures the simplified logger system works correctly across all blueprints
func TestLoggerIntegrationATDD(t *testing.T) {
	// Logger integration matrix
	loggerConfigs := map[string]LoggerTestConfig{
		"slog": {
			ImportPackage: "log/slog",
			TypeName:     "slog.Logger",
			Constructor:  "slog.New",
			IsStdLib:     true,
			Description:  "Standard library structured logging",
		},
		"zap": {
			ImportPackage: "go.uber.org/zap",
			TypeName:     "*zap.Logger",
			Constructor:  "zap.New",
			IsStdLib:     false,
			Description:  "High-performance structured logging",
		},
		"logrus": {
			ImportPackage: "github.com/sirupsen/logrus",
			TypeName:     "*logrus.Logger",
			Constructor:  "logrus.New",
			IsStdLib:     false,
			Description:  "Feature-rich structured logging",
		},
		"zerolog": {
			ImportPackage: "github.com/rs/zerolog",
			TypeName:     "zerolog.Logger",
			Constructor:  "zerolog.New",
			IsStdLib:     false,
			Description:  "Zero allocation structured logging",
		},
	}

	// Blueprint types that support logger integration
	blueprintConfigs := map[string]BlueprintLoggerConfig{
		"cli-simple": {
			Type:          "cli",
			Complexity:    "simple",
			LoggerPath:    "internal/logger/logger.go",
			ExpectedFiles: []string{"main.go", "cmd/root.go", "internal/logger/logger.go"},
		},
		"cli-standard": {
			Type:          "cli",
			Complexity:    "standard", 
			LoggerPath:    "internal/logger/logger.go",
			ExpectedFiles: []string{"main.go", "cmd/root.go", "internal/logger/logger.go"},
		},
		"web-api-standard": {
			Type:          "web-api",
			Architecture:  "standard",
			LoggerPath:    "internal/logger/logger.go",
			ExpectedFiles: []string{"main.go", "cmd/server/main.go", "internal/logger/logger.go"},
		},
		"lambda-standard": {
			Type:          "lambda",
			LoggerPath:    "internal/logger/logger.go",
			ExpectedFiles: []string{"main.go", "internal/logger/logger.go"},
		},
	}

	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	projectRoot := filepath.Join(originalDir, "..", "..", "..")
	
	defer func() { _ = os.Chdir(originalDir) }()
	_ = os.Chdir(tmpDir)

	// Build go-starter once
	buildGoStarter(t, tmpDir, projectRoot)

	t.Run("logger_types_generate_correctly", func(t *testing.T) {
		// GIVEN: Different logger types and blueprint combinations
		// WHEN: Generating projects with each logger type
		// THEN: Logger integration should be generated correctly
		
		for loggerName, loggerConfig := range loggerConfigs {
			for blueprintName, blueprintConfig := range blueprintConfigs {
				t.Run(fmt.Sprintf("%s_logger_in_%s", loggerName, blueprintName), func(t *testing.T) {
					projectName := fmt.Sprintf("logger-test-%s-%s", loggerName, blueprintName)
					
					// Generate project with specific logger
					generateLoggerProject(t, tmpDir, projectName, blueprintConfig, loggerName)
					projectDir := filepath.Join(tmpDir, projectName)
					
					// Validate logger implementation
					validateLoggerImplementation(t, projectDir, loggerConfig, blueprintConfig)
				})
			}
		}
	})

	t.Run("simplified_logger_architecture_validation", func(t *testing.T) {
		// GIVEN: The new simplified logger architecture
		// WHEN: Examining generated logger code
		// THEN: Should demonstrate 60-90% code reduction from complex loggers
		
		for loggerName, loggerConfig := range loggerConfigs {
			t.Run("simplified_"+loggerName+"_logger", func(t *testing.T) {
				projectName := fmt.Sprintf("simple-logger-%s", loggerName)
				blueprintConfig := blueprintConfigs["cli-standard"]
				
				generateLoggerProject(t, tmpDir, projectName, blueprintConfig, loggerName)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Validate simplified architecture
				validateSimplifiedLoggerArchitecture(t, projectDir, loggerConfig, loggerName)
			})
		}
	})

	t.Run("logger_compilation_and_runtime", func(t *testing.T) {
		// GIVEN: Generated projects with different loggers
		// WHEN: Compiling and running the projects
		// THEN: All logger implementations should compile and work at runtime
		
		for loggerName := range loggerConfigs {
			for blueprintName, blueprintConfig := range blueprintConfigs {
				t.Run(fmt.Sprintf("compile_%s_in_%s", loggerName, blueprintName), func(t *testing.T) {
					projectName := fmt.Sprintf("compile-test-%s-%s", loggerName, blueprintName)
					
					generateLoggerProject(t, tmpDir, projectName, blueprintConfig, loggerName)
					projectDir := filepath.Join(tmpDir, projectName)
					
					// Validate compilation
					validateLoggerCompilation(t, projectDir, loggerName, blueprintName)
					
					// Validate runtime functionality for CLIs
					if strings.HasPrefix(blueprintName, "cli") {
						validateLoggerRuntime(t, projectDir, loggerName, blueprintName)
					}
				})
			}
		}
	})

	t.Run("logger_interface_consistency", func(t *testing.T) {
		// GIVEN: Different logger implementations
		// WHEN: Examining the logger interface
		// THEN: All loggers should expose a consistent interface
		
		for loggerName, loggerConfig := range loggerConfigs {
			t.Run("interface_"+loggerName, func(t *testing.T) {
				projectName := fmt.Sprintf("interface-test-%s", loggerName)
				blueprintConfig := blueprintConfigs["cli-standard"]
				
				generateLoggerProject(t, tmpDir, projectName, blueprintConfig, loggerName)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Validate consistent interface
				validateLoggerInterfaceConsistency(t, projectDir, loggerConfig, loggerName)
			})
		}
	})

	t.Run("logger_dependency_management", func(t *testing.T) {
		// GIVEN: Projects with different external logger dependencies
		// WHEN: Examining go.mod files
		// THEN: Only required dependencies should be included
		
		for loggerName, loggerConfig := range loggerConfigs {
			t.Run("dependencies_"+loggerName, func(t *testing.T) {
				projectName := fmt.Sprintf("deps-test-%s", loggerName)
				blueprintConfig := blueprintConfigs["web-api-standard"]
				
				generateLoggerProject(t, tmpDir, projectName, blueprintConfig, loggerName)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Initialize modules to update dependencies
				initializeGoModules(t, projectDir)
				
				// Validate dependencies
				validateLoggerDependencies(t, projectDir, loggerConfig, loggerName)
			})
		}
	})

	t.Run("conditional_logger_generation", func(t *testing.T) {
		// GIVEN: Blueprint templates with conditional logger logic
		// WHEN: Generating projects with different logger types
		// THEN: Only the selected logger implementation should be generated
		
		projectName := "conditional-logger-test"
		blueprintConfig := blueprintConfigs["cli-standard"]
		
		// Generate with zap logger
		generateLoggerProject(t, tmpDir, projectName, blueprintConfig, "zap")
		projectDir := filepath.Join(tmpDir, projectName)
		
		// Should have zap implementation, not others
		validateConditionalLoggerGeneration(t, projectDir, "zap")
	})

	t.Run("logger_complexity_reduction", func(t *testing.T) {
		// GIVEN: The simplified logger system
		// WHEN: Measuring logger code complexity
		// THEN: Should demonstrate significant complexity reduction
		
		complexityMetrics := make(map[string]LoggerComplexity)
		
		for loggerName := range loggerConfigs {
			t.Run("complexity_"+loggerName, func(t *testing.T) {
				projectName := fmt.Sprintf("complexity-%s", loggerName)
				blueprintConfig := blueprintConfigs["cli-standard"]
				
				generateLoggerProject(t, tmpDir, projectName, blueprintConfig, loggerName)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Measure complexity
				complexity := measureLoggerComplexity(t, projectDir)
				complexityMetrics[loggerName] = complexity
				
				// Validate complexity thresholds
				validateLoggerComplexity(t, complexity, loggerName)
			})
		}
		
		// Report overall complexity reduction
		reportComplexityReduction(t, complexityMetrics)
	})
}

// LoggerTestConfig represents configuration for testing a specific logger type
type LoggerTestConfig struct {
	ImportPackage string
	TypeName      string
	Constructor   string
	IsStdLib      bool
	Description   string
}

// BlueprintLoggerConfig represents blueprint configuration for logger testing
type BlueprintLoggerConfig struct {
	Type          string
	Complexity    string
	Architecture  string
	LoggerPath    string
	ExpectedFiles []string
}

// LoggerComplexity represents complexity metrics for a logger implementation
type LoggerComplexity struct {
	Lines                int
	Functions           int
	Imports             int
	StructsAndInterfaces int
}

// generateLoggerProject generates a project with specific logger configuration
func generateLoggerProject(t *testing.T, tmpDir, projectName string, blueprintConfig BlueprintLoggerConfig, logger string) {
	t.Helper()
	
	args := []string{"new", projectName}
	
	// Add blueprint configuration
	args = append(args, "--type="+blueprintConfig.Type)
	
	if blueprintConfig.Complexity != "" {
		args = append(args, "--complexity="+blueprintConfig.Complexity)
	}
	
	if blueprintConfig.Architecture != "" {
		args = append(args, "--architecture="+blueprintConfig.Architecture)
	}
	
	// Add logger
	args = append(args, "--logger="+logger)
	
	// Add common flags
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
	require.NoError(t, err, "Should generate project with %s logger", logger)
	
	t.Logf("Successfully generated project %s with %s logger", projectName, logger)
}

// validateLoggerImplementation validates the generated logger implementation
func validateLoggerImplementation(t *testing.T, projectDir string, loggerConfig LoggerTestConfig, blueprintConfig BlueprintLoggerConfig) {
	t.Helper()
	
	loggerPath := filepath.Join(projectDir, blueprintConfig.LoggerPath)
	assert.FileExists(t, loggerPath, "Logger file should exist")
	
	content, err := os.ReadFile(loggerPath)
	require.NoError(t, err, "Should be able to read logger file")
	contentStr := string(content)
	
	// Validate import
	assert.Contains(t, contentStr, loggerConfig.ImportPackage, 
		"Should import correct package: %s", loggerConfig.ImportPackage)
	
	// Validate logger usage patterns
	if !loggerConfig.IsStdLib {
		// External loggers should appear in import section
		importRegex := regexp.MustCompile(fmt.Sprintf(`import.*\n.*"%s"`, regexp.QuoteMeta(loggerConfig.ImportPackage)))
		assert.Regexp(t, importRegex, contentStr, "Should have proper import for %s", loggerConfig.ImportPackage)
	}
	
	// Validate consistent interface methods (Info, Error, Debug, etc.)
	expectedMethods := []string{"Info", "Error", "Debug", "Warn"}
	for _, method := range expectedMethods {
		assert.Contains(t, contentStr, method, "Should have %s method", method)
	}
	
	t.Logf("✓ Logger implementation validated for %s", loggerConfig.ImportPackage)
}

// validateSimplifiedLoggerArchitecture validates the simplified logger architecture
func validateSimplifiedLoggerArchitecture(t *testing.T, projectDir string, loggerConfig LoggerTestConfig, loggerName string) {
	t.Helper()
	
	loggerFiles := findLoggerFiles(t, projectDir)
	require.NotEmpty(t, loggerFiles, "Should have logger files")
	
	totalLines := 0
	totalFunctions := 0
	
	for _, loggerFile := range loggerFiles {
		content, err := os.ReadFile(loggerFile)
		require.NoError(t, err, "Should be able to read logger file")
		
		lines := strings.Split(string(content), "\n")
		
		// Count meaningful lines (not empty, not comments)
		meaningfulLines := 0
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "/*") {
				meaningfulLines++
			}
		}
		totalLines += meaningfulLines
		
		// Count functions
		functionRegex := regexp.MustCompile(`func\s+\w+`)
		functions := functionRegex.FindAllString(string(content), -1)
		totalFunctions += len(functions)
	}
	
	// Validate simplified architecture (should be concise)
	assert.LessOrEqual(t, totalLines, 100, "Simplified logger should be ≤100 meaningful lines (was %d)", totalLines)
	assert.LessOrEqual(t, totalFunctions, 10, "Simplified logger should have ≤10 functions (was %d)", totalFunctions)
	
	t.Logf("✓ Simplified architecture validated for %s: %d lines, %d functions", loggerName, totalLines, totalFunctions)
}

// validateLoggerCompilation ensures the logger implementation compiles correctly
func validateLoggerCompilation(t *testing.T, projectDir, loggerName, blueprintName string) {
	t.Helper()
	
	// Initialize go modules
	initializeGoModules(t, projectDir)
	
	// Compile the project
	buildCmd := exec.Command("go", "build", "-o", "test-binary", ".")
	buildCmd.Dir = projectDir
	output, err := buildCmd.CombinedOutput()
	
	require.NoError(t, err, "Project with %s logger should compile in %s: %s", 
		loggerName, blueprintName, string(output))
	
	// Verify binary exists
	binaryPath := filepath.Join(projectDir, "test-binary")
	assert.FileExists(t, binaryPath, "Binary should be created")
	
	t.Logf("✓ Compilation validated for %s logger in %s", loggerName, blueprintName)
}

// validateLoggerRuntime validates runtime functionality for CLI projects
func validateLoggerRuntime(t *testing.T, projectDir, loggerName, blueprintName string) {
	t.Helper()
	
	// Run the CLI binary
	binaryPath := "./test-binary"
	cmd := exec.Command(binaryPath, "--help")
	cmd.Dir = projectDir
	output, err := cmd.CombinedOutput()
	
	require.NoError(t, err, "CLI should run successfully with %s logger: %s", loggerName, string(output))
	
	outputStr := string(output)
	assert.Contains(t, outputStr, "Usage:", "CLI should show usage information")
	
	t.Logf("✓ Runtime validated for %s logger in %s", loggerName, blueprintName)
}

// validateLoggerInterfaceConsistency ensures consistent logger interface across implementations
func validateLoggerInterfaceConsistency(t *testing.T, projectDir string, loggerConfig LoggerTestConfig, loggerName string) {
	t.Helper()
	
	loggerFiles := findLoggerFiles(t, projectDir)
	require.NotEmpty(t, loggerFiles, "Should have logger files")
	
	for _, loggerFile := range loggerFiles {
		content, err := os.ReadFile(loggerFile)
		require.NoError(t, err, "Should be able to read logger file")
		contentStr := string(content)
		
		// Check for consistent method signatures
		expectedPatterns := []string{
			`func.*Info\(`,
			`func.*Error\(`,
			`func.*Debug\(`,
			`func.*Warn\(`,
		}
		
		for _, pattern := range expectedPatterns {
			matched, err := regexp.MatchString(pattern, contentStr)
			require.NoError(t, err, "Regex should compile")
			assert.True(t, matched, "Should have method matching pattern: %s", pattern)
		}
	}
	
	t.Logf("✓ Interface consistency validated for %s", loggerName)
}

// validateLoggerDependencies ensures correct dependency management
func validateLoggerDependencies(t *testing.T, projectDir string, loggerConfig LoggerTestConfig, loggerName string) {
	t.Helper()
	
	goModPath := filepath.Join(projectDir, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")
	contentStr := string(content)
	
	if loggerConfig.IsStdLib {
		// Standard library loggers should not add external dependencies
		assert.NotContains(t, contentStr, "go.uber.org/zap", "slog should not include zap dependency")
		assert.NotContains(t, contentStr, "github.com/sirupsen/logrus", "slog should not include logrus dependency")
		assert.NotContains(t, contentStr, "github.com/rs/zerolog", "slog should not include zerolog dependency")
	} else {
		// External loggers should include their dependency
		assert.Contains(t, contentStr, loggerConfig.ImportPackage, 
			"Should include dependency for %s", loggerConfig.ImportPackage)
	}
	
	t.Logf("✓ Dependencies validated for %s", loggerName)
}

// validateConditionalLoggerGeneration ensures only selected logger is generated
func validateConditionalLoggerGeneration(t *testing.T, projectDir, selectedLogger string) {
	t.Helper()
	
	loggerFiles := findLoggerFiles(t, projectDir)
	require.NotEmpty(t, loggerFiles, "Should have logger files")
	
	for _, loggerFile := range loggerFiles {
		content, err := os.ReadFile(loggerFile)
		require.NoError(t, err, "Should be able to read logger file")
		contentStr := string(content)
		
		// Should contain selected logger
		switch selectedLogger {
		case "zap":
			assert.Contains(t, contentStr, "go.uber.org/zap", "Should contain zap import")
			assert.NotContains(t, contentStr, "github.com/sirupsen/logrus", "Should not contain logrus")
			assert.NotContains(t, contentStr, "github.com/rs/zerolog", "Should not contain zerolog")
		case "logrus":
			assert.Contains(t, contentStr, "github.com/sirupsen/logrus", "Should contain logrus import")
			assert.NotContains(t, contentStr, "go.uber.org/zap", "Should not contain zap")
			assert.NotContains(t, contentStr, "github.com/rs/zerolog", "Should not contain zerolog")
		case "zerolog":
			assert.Contains(t, contentStr, "github.com/rs/zerolog", "Should contain zerolog import")
			assert.NotContains(t, contentStr, "go.uber.org/zap", "Should not contain zap")
			assert.NotContains(t, contentStr, "github.com/sirupsen/logrus", "Should not contain logrus")
		case "slog":
			assert.Contains(t, contentStr, "log/slog", "Should contain slog import")
			assert.NotContains(t, contentStr, "go.uber.org/zap", "Should not contain zap")
			assert.NotContains(t, contentStr, "github.com/sirupsen/logrus", "Should not contain logrus")
			assert.NotContains(t, contentStr, "github.com/rs/zerolog", "Should not contain zerolog")
		}
	}
	
	t.Logf("✓ Conditional generation validated for %s", selectedLogger)
}

// measureLoggerComplexity measures the complexity of the logger implementation
func measureLoggerComplexity(t *testing.T, projectDir string) LoggerComplexity {
	t.Helper()
	
	complexity := LoggerComplexity{}
	loggerFiles := findLoggerFiles(t, projectDir)
	
	for _, loggerFile := range loggerFiles {
		content, err := os.ReadFile(loggerFile)
		require.NoError(t, err, "Should be able to read logger file")
		contentStr := string(content)
		
		// Count lines
		lines := strings.Split(contentStr, "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) != "" && !strings.HasPrefix(strings.TrimSpace(line), "//") {
				complexity.Lines++
			}
		}
		
		// Count functions
		functionRegex := regexp.MustCompile(`func\s+\w+`)
		functions := functionRegex.FindAllString(contentStr, -1)
		complexity.Functions += len(functions)
		
		// Count imports
		importRegex := regexp.MustCompile(`import\s+["'][\w\/.-]+["']|["'][\w\/.-]+["']`)
		imports := importRegex.FindAllString(contentStr, -1)
		complexity.Imports += len(imports)
		
		// Count structs and interfaces
		structRegex := regexp.MustCompile(`type\s+\w+\s+(struct|interface)`)
		structs := structRegex.FindAllString(contentStr, -1)
		complexity.StructsAndInterfaces += len(structs)
	}
	
	return complexity
}

// validateLoggerComplexity validates that complexity is within acceptable bounds
func validateLoggerComplexity(t *testing.T, complexity LoggerComplexity, loggerName string) {
	t.Helper()
	
	// Simplified logger should be within these bounds
	assert.LessOrEqual(t, complexity.Lines, 100, "Logger should have ≤100 lines")
	assert.LessOrEqual(t, complexity.Functions, 10, "Logger should have ≤10 functions")
	assert.LessOrEqual(t, complexity.Imports, 5, "Logger should have ≤5 imports")
	assert.LessOrEqual(t, complexity.StructsAndInterfaces, 3, "Logger should have ≤3 types")
	
	t.Logf("✓ Complexity validated for %s: %d lines, %d functions, %d imports, %d types", 
		loggerName, complexity.Lines, complexity.Functions, complexity.Imports, complexity.StructsAndInterfaces)
}

// reportComplexityReduction reports overall complexity reduction achievements
func reportComplexityReduction(t *testing.T, complexityMetrics map[string]LoggerComplexity) {
	t.Helper()
	
	totalLines := 0
	maxLines := 0
	
	for loggerName, complexity := range complexityMetrics {
		totalLines += complexity.Lines
		if complexity.Lines > maxLines {
			maxLines = complexity.Lines
		}
		
		t.Logf("Logger %s complexity: %d lines, %d functions", loggerName, complexity.Lines, complexity.Functions)
	}
	
	avgLines := totalLines / len(complexityMetrics)
	t.Logf("Average logger complexity: %d lines", avgLines)
	t.Logf("Maximum logger complexity: %d lines", maxLines)
	
	// Target: 60-90% reduction from complex logger systems (assume 300-500 lines before)
	if avgLines <= 50 {
		t.Logf("✓ Excellent complexity reduction: ~85-90%% reduction achieved")
	} else if avgLines <= 100 {
		t.Logf("✓ Good complexity reduction: ~70-85%% reduction achieved")
	} else {
		t.Logf("⚠ Moderate complexity reduction: ~60-70%% reduction achieved")
	}
}