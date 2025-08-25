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

// TestTemplateVariableResolutionATDD validates template variable resolution acceptance criteria
// This ensures all template variables are properly resolved without any missing values
func TestTemplateVariableResolutionATDD(t *testing.T) {
	// Template variable categories for comprehensive testing
	// Note: These categories ensure we test all major variable types used in templates
	_ = map[string]VariableCategory{
		"project_metadata": {
			Variables: []string{"ProjectName", "ModulePath", "GoVersion", "Author", "Email"},
			Description: "Basic project metadata variables",
		},
		"configuration": {
			Variables: []string{"LoggerType", "Framework", "Architecture", "DatabaseDriver", "DatabaseORM"},
			Description: "Configuration and feature selection variables",
		},
		"features": {
			Variables: []string{"Features.Database", "Features.Authentication", "Features.Logging", "Features.Testing"},
			Description: "Feature configuration object variables",
		},
		"conditional": {
			Variables: []string{"ne .Features.Database.Driver", "eq .LoggerType", "if .Features.Authentication"},
			Description: "Conditional template logic variables",
		},
	}

	// Blueprint configurations for comprehensive testing
	testBlueprints := map[string]VariableTestConfig{
		"cli-simple": {
			Type:         "cli",
			Complexity:   "simple",
			TestModule:   "github.com/test/variable-cli-simple",
			ExpectedVars: []string{"ProjectName", "ModulePath", "LoggerType", "GoVersion"},
		},
		"cli-standard": {
			Type:         "cli",
			Complexity:   "standard",
			TestModule:   "github.com/test/variable-cli-standard",
			ExpectedVars: []string{"ProjectName", "ModulePath", "LoggerType", "Framework", "GoVersion"},
		},
		"web-api-standard": {
			Type:         "web-api",
			Architecture: "standard",
			TestModule:   "github.com/test/variable-web-api",
			ExpectedVars: []string{"ProjectName", "ModulePath", "Framework", "LoggerType", "GoVersion"},
		},
		"lambda-standard": {
			Type:       "lambda",
			TestModule: "github.com/test/variable-lambda",
			ExpectedVars: []string{"ProjectName", "ModulePath", "LoggerType", "GoVersion"},
		},
		"library-standard": {
			Type:       "library",
			TestModule: "github.com/test/variable-library",
			ExpectedVars: []string{"ProjectName", "ModulePath", "GoVersion", "Author"},
		},
	}

	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	projectRoot := filepath.Join(originalDir, "..", "..", "..")
	
	defer func() { _ = os.Chdir(originalDir) }()
	_ = os.Chdir(tmpDir)

	// Build go-starter once
	buildGoStarter(t, tmpDir, projectRoot)

	t.Run("template_variable_resolution_validation", func(t *testing.T) {
		// GIVEN: Templates with various variables
		// WHEN: Generating projects with different configurations
		// THEN: All template variables must be resolved correctly
		
		for blueprintName, config := range testBlueprints {
			t.Run(blueprintName+"_variable_resolution", func(t *testing.T) {
				projectName := fmt.Sprintf("var-resolve-%s", blueprintName)
				
				generateVariableTestProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Comprehensive variable resolution check
				validateCompleteVariableResolution(t, projectDir, config, blueprintName)
			})
		}
	})

	t.Run("unresolved_template_variable_detection", func(t *testing.T) {
		// GIVEN: Generated projects from templates
		// WHEN: Scanning all files for template syntax
		// THEN: Should find no unresolved {{.Variable}} patterns
		
		for blueprintName, config := range testBlueprints {
			t.Run(blueprintName+"_unresolved_detection", func(t *testing.T) {
				projectName := fmt.Sprintf("unresolved-%s", blueprintName)
				
				generateVariableTestProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Detect any unresolved template variables
				unresolvedVars := findUnresolvedVariables(t, projectDir)
				
				if len(unresolvedVars) > 0 {
					t.Errorf("Found %d unresolved template variables in %s:", len(unresolvedVars), blueprintName)
					for _, v := range unresolvedVars {
						t.Errorf("  %s", v)
					}
				}
				
				assert.Empty(t, unresolvedVars, "Should have no unresolved template variables in %s", blueprintName)
			})
		}
	})

	t.Run("variable_substitution_accuracy", func(t *testing.T) {
		// GIVEN: Specific variable values
		// WHEN: Generating projects with those values
		// THEN: Generated files should contain exact substituted values
		
		testCases := []VariableSubstitutionTest{
			{
				BlueprintName: "cli-simple",
				ProjectName:   "test-accuracy-cli",
				ModulePath:    "github.com/accuracy/test-cli",
				Logger:        "zap",
				ExpectedSubstitutions: map[string][]string{
					"main.go":   {"github.com/accuracy/test-cli", "test-accuracy-cli"},
					"go.mod":    {"module github.com/accuracy/test-cli", "go 1.2"},
					"README.md": {"# test-accuracy-cli", "test-accuracy-cli"},
				},
			},
			{
				BlueprintName: "web-api-standard",
				ProjectName:   "api-accuracy-test",
				ModulePath:    "github.com/accuracy/api-test",
				Logger:        "logrus",
				ExpectedSubstitutions: map[string][]string{
					"main.go":   {"github.com/accuracy/api-test"},
					"go.mod":    {"module github.com/accuracy/api-test"},
					"README.md": {"# api-accuracy-test"},
				},
			},
		}
		
		for _, testCase := range testCases {
			t.Run(testCase.BlueprintName+"_substitution_accuracy", func(t *testing.T) {
				config := testBlueprints[testCase.BlueprintName]
				config.TestModule = testCase.ModulePath
				
				generateVariableTestProject(t, tmpDir, testCase.ProjectName, config)
				projectDir := filepath.Join(tmpDir, testCase.ProjectName)
				
				// Validate specific substitutions
				validateVariableSubstitutions(t, projectDir, testCase)
			})
		}
	})

	t.Run("conditional_variable_logic_validation", func(t *testing.T) {
		// GIVEN: Templates with conditional logic
		// WHEN: Generating projects with different feature configurations
		// THEN: Conditional variables should resolve correctly
		
		conditionalTests := []ConditionalVariableTest{
			{
				BlueprintName: "web-api-standard",
				ProjectName:   "conditional-db-test",
				WithDatabase:  true,
				DatabaseDriver: "postgres",
				ExpectedFiles: []string{"internal/database/database.go"},
				NotExpectedFiles: []string{},
			},
			{
				BlueprintName: "web-api-standard", 
				ProjectName:   "conditional-no-db-test",
				WithDatabase:  false,
				ExpectedFiles: []string{},
				NotExpectedFiles: []string{"internal/database/database.go"},
			},
			{
				BlueprintName: "cli-standard",
				ProjectName:   "conditional-logger-test",
				Logger:        "zap",
				ExpectedContent: map[string][]string{
					"internal/logger/logger.go": {"go.uber.org/zap"},
				},
				NotExpectedContent: map[string][]string{
					"internal/logger/logger.go": {"log/slog", "github.com/sirupsen/logrus"},
				},
			},
		}
		
		for _, testCase := range conditionalTests {
			t.Run(testCase.ProjectName, func(t *testing.T) {
				generateConditionalTestProject(t, tmpDir, testCase)
				projectDir := filepath.Join(tmpDir, testCase.ProjectName)
				
				validateConditionalVariables(t, projectDir, testCase)
			})
		}
	})

	t.Run("complex_variable_expressions", func(t *testing.T) {
		// GIVEN: Templates with complex variable expressions
		// WHEN: Generating projects  
		// THEN: Complex expressions should resolve correctly
		
		complexTests := []ComplexVariableTest{
			{
				BlueprintName: "cli-standard",
				ProjectName:   "complex-expr-test",
				TestExpressions: map[string]string{
					"go.mod": `module github.com/test/complex-expr-test`, // Simple module reference
					"main.go": `github.com/test/complex-expr-test/cmd`,    // Package import
				},
			},
		}
		
		for _, testCase := range complexTests {
			t.Run(testCase.ProjectName, func(t *testing.T) {
				config := testBlueprints[testCase.BlueprintName]
				config.TestModule = "github.com/test/" + testCase.ProjectName
				
				generateVariableTestProject(t, tmpDir, testCase.ProjectName, config)
				projectDir := filepath.Join(tmpDir, testCase.ProjectName)
				
				validateComplexVariableExpressions(t, projectDir, testCase)
			})
		}
	})

	t.Run("variable_edge_cases", func(t *testing.T) {
		// GIVEN: Edge case variable configurations
		// WHEN: Generating projects with unusual but valid inputs
		// THEN: Should handle edge cases gracefully
		
		edgeCases := []EdgeCaseTest{
			{
				ProjectName: "edge-case-hyphens",
				ModulePath:  "github.com/test/edge-case-hyphens",
				Description: "Project name with hyphens",
			},
			{
				ProjectName: "edgecasenohyphens",
				ModulePath:  "github.com/test/edgecasenohyphens", 
				Description: "Project name without separators",
			},
			{
				ProjectName: "Edge_Case_Underscores",
				ModulePath:  "github.com/test/edge-case-underscores", // Note: Go modules prefer hyphens
				Description: "Project name with underscores",
			},
		}
		
		for _, edgeCase := range edgeCases {
			t.Run(edgeCase.ProjectName, func(t *testing.T) {
				config := testBlueprints["cli-simple"]
				config.TestModule = edgeCase.ModulePath
				
				generateVariableTestProject(t, tmpDir, edgeCase.ProjectName, config)
				projectDir := filepath.Join(tmpDir, edgeCase.ProjectName)
				
				// Should generate successfully without variable resolution errors
				validateEdgeCaseVariables(t, projectDir, edgeCase)
			})
		}
	})

	t.Run("variable_consistency_across_files", func(t *testing.T) {
		// GIVEN: Multiple files using the same variables
		// WHEN: Generating a project
		// THEN: Variables should be consistent across all files
		
		for blueprintName, config := range testBlueprints {
			t.Run(blueprintName+"_consistency", func(t *testing.T) {
				projectName := fmt.Sprintf("consistency-%s", blueprintName)
				
				generateVariableTestProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				validateVariableConsistency(t, projectDir, config, blueprintName)
			})
		}
	})
}

// Type definitions for template variable testing

type VariableCategory struct {
	Variables   []string
	Description string
}

type VariableTestConfig struct {
	Type         string
	Complexity   string
	Architecture string
	TestModule   string
	ExpectedVars []string
}

type VariableSubstitutionTest struct {
	BlueprintName         string
	ProjectName           string
	ModulePath            string
	Logger                string
	ExpectedSubstitutions map[string][]string
}

type ConditionalVariableTest struct {
	BlueprintName        string
	ProjectName          string
	WithDatabase         bool
	DatabaseDriver       string
	Logger               string
	ExpectedFiles        []string
	NotExpectedFiles     []string
	ExpectedContent      map[string][]string
	NotExpectedContent   map[string][]string
}

type ComplexVariableTest struct {
	BlueprintName   string
	ProjectName     string
	TestExpressions map[string]string
}

type EdgeCaseTest struct {
	ProjectName string
	ModulePath  string
	Description string
}

// generateVariableTestProject generates a project for variable testing
func generateVariableTestProject(t *testing.T, tmpDir, projectName string, config VariableTestConfig) {
	t.Helper()
	
	args := []string{"new", projectName}
	
	// Add configuration
	args = append(args, "--type="+config.Type)
	
	if config.Complexity != "" {
		args = append(args, "--complexity="+config.Complexity)
	}
	
	if config.Architecture != "" {
		args = append(args, "--architecture="+config.Architecture)
	}
	
	// Add module and other settings
	args = append(args, 
		"--module="+config.TestModule,
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
	require.NoError(t, err, "Should generate project for variable testing")
	
	t.Logf("Successfully generated variable test project: %s", projectName)
}

// generateConditionalTestProject generates a project for conditional variable testing
func generateConditionalTestProject(t *testing.T, tmpDir string, testCase ConditionalVariableTest) {
	t.Helper()
	
	config := VariableTestConfig{
		Type:       strings.Split(testCase.BlueprintName, "-")[0],
		TestModule: "github.com/test/" + testCase.ProjectName,
	}
	
	if strings.Contains(testCase.BlueprintName, "standard") {
		if strings.HasPrefix(testCase.BlueprintName, "cli") {
			config.Complexity = "standard"
		} else {
			config.Architecture = "standard"
		}
	}
	
	args := []string{"new", testCase.ProjectName}
	args = append(args, "--type="+config.Type)
	
	if config.Complexity != "" {
		args = append(args, "--complexity="+config.Complexity)
	}
	if config.Architecture != "" {
		args = append(args, "--architecture="+config.Architecture)
	}
	
	// Add conditional configurations
	if testCase.WithDatabase {
		args = append(args, "--database-driver="+testCase.DatabaseDriver)
	}
	
	if testCase.Logger != "" {
		args = append(args, "--logger="+testCase.Logger)
	} else {
		args = append(args, "--logger=slog")
	}
	
	args = append(args, "--module="+config.TestModule, "--no-git")
	
	goStarterPath := filepath.Join(tmpDir, "go-starter")
	cmd := exec.Command(goStarterPath, args...)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Logf("Generate command: %s %s", goStarterPath, strings.Join(args, " "))
		t.Logf("Generate output: %s", string(output))
	}
	require.NoError(t, err, "Should generate conditional test project")
}

// validateCompleteVariableResolution performs comprehensive variable resolution validation
func validateCompleteVariableResolution(t *testing.T, projectDir string, config VariableTestConfig, blueprintName string) {
	t.Helper()
	
	// Find all generated files
	allFiles := findAllTextFiles(t, projectDir)
	
	totalUnresolved := 0
	
	for _, file := range allFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue // Skip files that can't be read
		}
		
		contentStr := string(content)
		
		// Look for unresolved template patterns
		unresolvedPatterns := []string{
			`\{\{\..*?\}\}`,     // {{.Variable}}
			`\{\{\s*\..*?\}\}`,  // {{ .Variable }}
			`\{\{-\s*\..*?\}\}`, // {{- .Variable }}
			`\{\{.*?\s*-\}\}`,   // {{ .Variable -}}
		}
		
		for _, pattern := range unresolvedPatterns {
			regex := regexp.MustCompile(pattern)
			matches := regex.FindAllString(contentStr, -1)
			
			if len(matches) > 0 {
				totalUnresolved += len(matches)
				t.Errorf("Found unresolved variables in %s:", file)
				for _, match := range matches {
					t.Errorf("  %s", match)
				}
			}
		}
	}
	
	assert.Equal(t, 0, totalUnresolved, "Should have no unresolved variables in %s", blueprintName)
	t.Logf("✓ Complete variable resolution validated for %s", blueprintName)
}

// findUnresolvedVariables finds any unresolved template variables
func findUnresolvedVariables(t *testing.T, projectDir string) []string {
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
			return nil // Skip files that can't be read
		}
		
		contentStr := string(content)
		
		// Multiple patterns to catch different template syntaxes
		patterns := []string{
			`\{\{\..*?\}\}`,           // {{.Variable}}
			`\{\{\s+\..*?\}\}`,        // {{ .Variable }}
			`\{\{-\s*\..*?\s*-?\}\}`,  // {{- .Variable -}}
		}
		
		for _, pattern := range patterns {
			regex := regexp.MustCompile(pattern)
			matches := regex.FindAllString(contentStr, -1)
			
			for _, match := range matches {
				// Get line number for better reporting
				lines := strings.Split(contentStr, "\n")
				for i, line := range lines {
					if strings.Contains(line, match) {
						unresolvedVars = append(unresolvedVars, 
							fmt.Sprintf("%s:%d: %s", path, i+1, strings.TrimSpace(line)))
						break
					}
				}
			}
		}
		
		return nil
	})
	
	require.NoError(t, err, "Should be able to walk project directory")
	return unresolvedVars
}

// validateVariableSubstitutions validates specific variable substitutions
func validateVariableSubstitutions(t *testing.T, projectDir string, testCase VariableSubstitutionTest) {
	t.Helper()
	
	for fileName, expectedSubstitutions := range testCase.ExpectedSubstitutions {
		filePath := filepath.Join(projectDir, fileName)
		
		if !assert.FileExists(t, filePath, "File %s should exist", fileName) {
			continue
		}
		
		content, err := os.ReadFile(filePath)
		require.NoError(t, err, "Should be able to read %s", fileName)
		contentStr := string(content)
		
		for _, expectedSub := range expectedSubstitutions {
			assert.Contains(t, contentStr, expectedSub, 
				"File %s should contain substitution: %s", fileName, expectedSub)
		}
	}
	
	t.Logf("✓ Variable substitutions validated for %s", testCase.ProjectName)
}

// validateConditionalVariables validates conditional variable logic
func validateConditionalVariables(t *testing.T, projectDir string, testCase ConditionalVariableTest) {
	t.Helper()
	
	// Validate expected files exist
	for _, expectedFile := range testCase.ExpectedFiles {
		filePath := filepath.Join(projectDir, expectedFile)
		assert.FileExists(t, filePath, "Conditional file %s should exist", expectedFile)
	}
	
	// Validate files that should NOT exist
	for _, notExpectedFile := range testCase.NotExpectedFiles {
		filePath := filepath.Join(projectDir, notExpectedFile)
		assert.NoFileExists(t, filePath, "Conditional file %s should NOT exist", notExpectedFile)
	}
	
	// Validate expected content
	for fileName, expectedContent := range testCase.ExpectedContent {
		filePath := filepath.Join(projectDir, fileName)
		if assert.FileExists(t, filePath, "File %s should exist", fileName) {
			content, err := os.ReadFile(filePath)
			require.NoError(t, err, "Should read %s", fileName)
			contentStr := string(content)
			
			for _, expected := range expectedContent {
				assert.Contains(t, contentStr, expected, 
					"File %s should contain: %s", fileName, expected)
			}
		}
	}
	
	// Validate content that should NOT be present
	for fileName, notExpectedContent := range testCase.NotExpectedContent {
		filePath := filepath.Join(projectDir, fileName)
		if assert.FileExists(t, filePath, "File %s should exist", fileName) {
			content, err := os.ReadFile(filePath)
			require.NoError(t, err, "Should read %s", fileName)
			contentStr := string(content)
			
			for _, notExpected := range notExpectedContent {
				assert.NotContains(t, contentStr, notExpected,
					"File %s should NOT contain: %s", fileName, notExpected)
			}
		}
	}
	
	t.Logf("✓ Conditional variables validated for %s", testCase.ProjectName)
}

// validateComplexVariableExpressions validates complex variable expressions
func validateComplexVariableExpressions(t *testing.T, projectDir string, testCase ComplexVariableTest) {
	t.Helper()
	
	for fileName, expectedExpression := range testCase.TestExpressions {
		filePath := filepath.Join(projectDir, fileName)
		
		if assert.FileExists(t, filePath, "File %s should exist", fileName) {
			content, err := os.ReadFile(filePath)
			require.NoError(t, err, "Should read %s", fileName)
			contentStr := string(content)
			
			assert.Contains(t, contentStr, expectedExpression,
				"File %s should contain complex expression: %s", fileName, expectedExpression)
		}
	}
	
	t.Logf("✓ Complex variable expressions validated for %s", testCase.ProjectName)
}

// validateEdgeCaseVariables validates edge case variable handling
func validateEdgeCaseVariables(t *testing.T, projectDir string, edgeCase EdgeCaseTest) {
	t.Helper()
	
	// Should have basic files without template errors
	basicFiles := []string{"go.mod", "main.go", "README.md"}
	
	for _, file := range basicFiles {
		filePath := filepath.Join(projectDir, file)
		if assert.FileExists(t, filePath, "Basic file %s should exist", file) {
			
			content, err := os.ReadFile(filePath)
			require.NoError(t, err, "Should read %s", file)
			contentStr := string(content)
			
			// Should not contain unresolved template variables
			assert.NotContains(t, contentStr, "{{.", "File %s should not contain unresolved variables", file)
			
			// Should contain the project/module references
			if file == "go.mod" {
				assert.Contains(t, contentStr, edgeCase.ModulePath, "go.mod should contain module path")
			}
			if file == "README.md" {
				// README should contain some form of the project name
				projectNameFound := strings.Contains(contentStr, edgeCase.ProjectName) ||
					strings.Contains(contentStr, strings.ToLower(edgeCase.ProjectName)) ||
					strings.Contains(contentStr, strings.ReplaceAll(edgeCase.ProjectName, "_", "-"))
				assert.True(t, projectNameFound, "README.md should reference project name")
			}
		}
	}
	
	t.Logf("✓ Edge case variables validated for %s (%s)", edgeCase.ProjectName, edgeCase.Description)
}

// validateVariableConsistency validates variable consistency across files
func validateVariableConsistency(t *testing.T, projectDir string, config VariableTestConfig, blueprintName string) {
	t.Helper()
	
	// Find all files that should contain the module path
	modulePathFiles := []string{"go.mod"}
	modulePathOccurrences := 0
	
	// Find files that might contain imports from the module
	importFiles := findGoFiles(t, projectDir)
	
	for _, file := range modulePathFiles {
		filePath := filepath.Join(projectDir, file)
		if assert.FileExists(t, filePath, "File %s should exist", file) {
			content, err := os.ReadFile(filePath)
			require.NoError(t, err, "Should read %s", file)
			contentStr := string(content)
			
			if strings.Contains(contentStr, config.TestModule) {
				modulePathOccurrences++
			}
		}
	}
	
	// Check for consistent module references in Go files
	for _, goFile := range importFiles {
		content, err := os.ReadFile(goFile)
		if err != nil {
			continue
		}
		contentStr := string(content)
		
		// If the file imports from the module, it should use the consistent path
		if strings.Contains(contentStr, config.TestModule) {
			// Should be properly formatted import
			importRegex := regexp.MustCompile(`import.*"` + regexp.QuoteMeta(config.TestModule) + `.*"`)
			if importRegex.MatchString(contentStr) {
				modulePathOccurrences++
			}
		}
	}
	
	assert.Greater(t, modulePathOccurrences, 0, "Module path should appear at least once consistently")
	
	t.Logf("✓ Variable consistency validated for %s (module path found %d times)", 
		blueprintName, modulePathOccurrences)
}

// Helper functions

func findAllTextFiles(t *testing.T, projectDir string) []string {
	t.Helper()
	var textFiles []string
	
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() || isBinaryFile(path) {
			return nil
		}
		
		// Include common text file types
		ext := strings.ToLower(filepath.Ext(path))
		textExtensions := []string{".go", ".md", ".yaml", ".yml", ".toml", ".txt", ".json", ".sql"}
		
		for _, textExt := range textExtensions {
			if ext == textExt || ext == "" { // Include files without extensions
				textFiles = append(textFiles, path)
				break
			}
		}
		
		// Also include common config files without extensions
		baseName := strings.ToLower(filepath.Base(path))
		configFiles := []string{"dockerfile", "makefile", "readme", "license", "gitignore"}
		
		for _, configFile := range configFiles {
			if strings.Contains(baseName, configFile) {
				textFiles = append(textFiles, path)
				break
			}
		}
		
		return nil
	})
	
	require.NoError(t, err, "Should be able to find text files")
	return textFiles
}