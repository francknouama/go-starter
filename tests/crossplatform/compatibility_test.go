package crossplatform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
)

// PlatformTestResult holds results for platform-specific testing
var (

	// Target platforms for cross-platform testing
	targetPlatforms = []struct {
		GOOS   string
		GOARCH string
		Name   string
	}{
		{"linux", "amd64", "Linux x86_64"},
		{"linux", "arm64", "Linux ARM64"},
		{"windows", "amd64", "Windows x86_64"},
		{"darwin", "amd64", "macOS x86_64"},
		{"darwin", "arm64", "macOS Apple Silicon"},
	}

	// Test blueprints for cross-platform validation
	testBlueprints = []struct {
		name   string
		config types.ProjectConfig
	}{
		{
			name: "cli-simple",
			config: types.ProjectConfig{
				Name:      "test-cli-cross",
				Module:    "github.com/test/cli-cross",
				Type:      "cli",
				GoVersion: "1.21",
				Framework: "cobra",
				Logger:    "slog",
				Variables: map[string]string{
					"blueprint_id": "cli-simple",
				},
			},
		},
		{
			name: "cli-standard",
			config: types.ProjectConfig{
				Name:      "test-cli-standard",
				Module:    "github.com/test/cli-standard",
				Type:      "cli",
				GoVersion: "1.21",
				Framework: "cobra",
				Logger:    "slog",
			},
		},
		{
			name: "web-api",
			config: types.ProjectConfig{
				Name:         "test-web-api",
				Module:       "github.com/test/web-api",
				Type:         "web-api",
				GoVersion:    "1.21",
				Framework:    "gin",
				Architecture: "standard",
				Logger:       "slog",
			},
		},
		{
			name: "lambda",
			config: types.ProjectConfig{
				Name:      "test-lambda",
				Module:    "github.com/test/lambda",
				Type:      "lambda",
				GoVersion: "1.21",
				Logger:    "slog",
			},
		},
	}
)

// Initialize test environment
func init() {
	// Set up embedded templates for testing
	blueprintsPath := findBlueprintsPath()
	if blueprintsPath != "" {
		templates.SetTemplatesFS(os.DirFS(blueprintsPath))
	}
}

// findBlueprintsPath locates the blueprints directory
func findBlueprintsPath() string {
	currentDir, _ := os.Getwd()
	for currentDir != "/" && currentDir != "" {
		blueprintsPath := filepath.Join(currentDir, "blueprints")
		if _, err := os.Stat(blueprintsPath); err == nil {
			return blueprintsPath
		}
		currentDir = filepath.Dir(currentDir)
	}
	return ""
}

// TestCrossPlatformCompatibility runs comprehensive cross-platform tests
func TestCrossPlatformCompatibility(t *testing.T) {
	if !isGoInstalled() {
		t.Skip("Go is not installed - skipping cross-platform compilation tests")
	}

	for _, blueprint := range testBlueprints {
		t.Run(fmt.Sprintf("CrossPlatform_%s", blueprint.name), func(t *testing.T) {
			testBlueprintCrossPlatform(t, blueprint.name, blueprint.config)
		})
	}

	// Generate compatibility report
	t.Run("GenerateCompatibilityReport", func(t *testing.T) {
		generateCompatibilityReport(t)
	})
}

// TestNativePlatformExecution tests execution on the current platform
func TestNativePlatformExecution(t *testing.T) {
	if !isGoInstalled() {
		t.Skip("Go is not installed - skipping native execution tests")
	}

	for _, blueprint := range testBlueprints {
		t.Run(fmt.Sprintf("Native_%s", blueprint.name), func(t *testing.T) {
			testNativeExecution(t, blueprint.name, blueprint.config)
		})
	}
}

// TestPathHandling tests cross-platform path handling
func TestPathHandling(t *testing.T) {
	testCases := []struct {
		name         string
		templatePath string
		expectedIssues []string
	}{
		{
			name:         "Unix paths",
			templatePath: "internal/config/config.go",
			expectedIssues: []string{},
		},
		{
			name:         "Nested paths",
			templatePath: "cmd/serve/handlers/health.go",
			expectedIssues: []string{},
		},
		{
			name:         "Special characters",
			templatePath: "docs/examples/hello-world.md",
			expectedIssues: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test path normalization across platforms
			normalizedPath := filepath.Clean(tc.templatePath)
			
			// Check for platform-specific issues
			var issues []string
			
			// Windows-specific checks
			if runtime.GOOS == "windows" {
				if strings.Contains(normalizedPath, ":") && !filepath.IsAbs(normalizedPath) {
					issues = append(issues, "Invalid colon in relative path on Windows")
				}
			}
			
			// Unix-specific checks
			if runtime.GOOS != "windows" {
				if strings.Contains(normalizedPath, "\\") {
					issues = append(issues, "Backslash in path on Unix system")
				}
			}
			
			// Report any unexpected issues
			for _, issue := range issues {
				if !contains(tc.expectedIssues, issue) {
					t.Errorf("Unexpected path issue: %s", issue)
				}
			}
		})
	}
}

// TestFilePermissions tests cross-platform file permissions
func TestFilePermissions(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "permission-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test-executable")
	
	// Create a test file
	if err := os.WriteFile(testFile, []byte("#!/bin/bash\necho hello"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test setting executable permissions
	if err := os.Chmod(testFile, 0755); err != nil {
		if runtime.GOOS == "windows" {
			t.Logf("chmod failed on Windows (expected): %v", err)
		} else {
			t.Errorf("Failed to set executable permissions: %v", err)
		}
	}

	// Verify permissions
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat test file: %v", err)
	}

	mode := info.Mode()
	if runtime.GOOS != "windows" {
		// On Unix systems, check for execute permission
		if mode&0111 == 0 {
			t.Error("Execute permission not set on Unix system")
		}
	}

	t.Logf("File permissions on %s: %s", runtime.GOOS, mode)
}

// testBlueprintCrossPlatform tests a blueprint across all target platforms
func testBlueprintCrossPlatform(t *testing.T, blueprintName string, config types.ProjectConfig) {
	for _, platform := range targetPlatforms {
		// Skip if cross-compilation is not supported
		if !supportsCrossCompilation(platform.GOOS, platform.GOARCH) {
			t.Logf("Skipping %s - cross-compilation not supported", platform.Name)
			continue
		}

		t.Run(platform.Name, func(t *testing.T) {
			result := testPlatformGeneration(blueprintName, config, platform.GOOS, platform.GOARCH)
			
			// Add result to compatibility report
			compatibilityReport.TestResults = append(compatibilityReport.TestResults, result)
			
			if !result.Success {
				t.Errorf("Generation failed for %s on %s: %s", 
					blueprintName, platform.Name, result.GenerationError)
			}
			
			// Test compilation for the target platform
			if result.Success {
				compileResult := testCrossCompilation(result, platform.GOOS, platform.GOARCH)
				
				// Update result with compilation info
				result.CompileTime = compileResult.CompileTime
				result.BinarySize = compileResult.BinarySize
				result.CompileError = compileResult.CompileError
				
				if compileResult.CompileError != "" {
					result.Success = false
					t.Errorf("Compilation failed for %s on %s: %s", 
						blueprintName, platform.Name, compileResult.CompileError)
				}
			}
		})
	}
}

// testNativeExecution tests execution on the current platform
func testNativeExecution(t *testing.T, blueprintName string, config types.ProjectConfig) {
	result := testPlatformGeneration(blueprintName, config, runtime.GOOS, runtime.GOARCH)
	
	if !result.Success {
		t.Fatalf("Generation failed: %s", result.GenerationError)
	}

	// Test compilation
	compileResult := testNativeCompilation(result)
	result.CompileTime = compileResult.CompileTime
	result.BinarySize = compileResult.BinarySize
	result.CompileError = compileResult.CompileError

	if compileResult.CompileError != "" {
		t.Fatalf("Compilation failed: %s", compileResult.CompileError)
	}

	// Test execution for CLI applications
	if strings.Contains(blueprintName, "cli") {
		execResult := testNativeExecution_CLI(result)
		result.ExecuteTime = execResult.ExecuteTime
		result.ExecuteError = execResult.ExecuteError

		if execResult.ExecuteError != "" {
			t.Errorf("Execution failed: %s", execResult.ExecuteError)
		} else {
			t.Logf("CLI execution successful in %v", execResult.ExecuteTime)
		}
	}

	// Add to compatibility report
	compatibilityReport.TestResults = append(compatibilityReport.TestResults, result)
}

// testPlatformGeneration generates a project for a specific platform
func testPlatformGeneration(blueprintName string, config types.ProjectConfig, goos, goarch string) PlatformTestResult {
	result := PlatformTestResult{
		Platform:     goos,
		Architecture: goarch,
		Blueprint:    blueprintName,
		Timestamp:    time.Now(),
	}

	// Create unique project config
	uniqueConfig := config
	uniqueConfig.Name = fmt.Sprintf("%s-%s-%s", config.Name, goos, goarch)
	uniqueConfig.Module = fmt.Sprintf("%s-%s-%s", config.Module, goos, goarch)

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("cross-platform-%s-*", blueprintName))
	if err != nil {
		result.GenerationError = fmt.Sprintf("Failed to create temp directory: %v", err)
		return result
	}
	defer os.RemoveAll(tempDir)

	projectPath := filepath.Join(tempDir, uniqueConfig.Name)

	// Generate project
	gen := generator.New()
	startTime := time.Now()

	genResult, err := gen.Generate(uniqueConfig, types.GenerationOptions{
		OutputPath: projectPath,
		DryRun:     false,
		NoGit:      true,
	})

	result.GenerationTime = time.Since(startTime)

	if err != nil {
		result.GenerationError = err.Error()
		return result
	}

	result.Success = true
	result.FilesGenerated = len(genResult.FilesCreated)

	// Check for platform-specific path issues
	result.PathIssues = checkPathIssues(genResult.FilesCreated, goos)
	
	// Check for permission issues
	result.PermissionIssues = checkPermissionIssues(genResult.FilesCreated, goos)

	return result
}

// CompilationResult holds compilation test results
type CompilationResult struct {
	CompileTime  time.Duration
	BinarySize   int64
	CompileError string
}

// ExecutionResult holds execution test results
type ExecutionResult struct {
	ExecuteTime  time.Duration
	ExecuteError string
	Output       string
}

// testCrossCompilation tests cross-compilation for a target platform
func testCrossCompilation(result PlatformTestResult, goos, goarch string) CompilationResult {
	compileResult := CompilationResult{}

	// Find the project directory from the result
	tempDir, err := os.MkdirTemp("", "cross-compile-*")
	if err != nil {
		compileResult.CompileError = fmt.Sprintf("Failed to create compile temp dir: %v", err)
		return compileResult
	}
	defer os.RemoveAll(tempDir)

	// Re-generate the project for compilation testing
	config := getConfigFromBlueprint(result.Blueprint)
	config.Name = fmt.Sprintf("%s-compile", config.Name)
	
	gen := generator.New()
	projectPath := filepath.Join(tempDir, config.Name)
	
	_, err = gen.Generate(config, types.GenerationOptions{
		OutputPath: projectPath,
		DryRun:     false,
		NoGit:      true,
	})
	
	if err != nil {
		compileResult.CompileError = fmt.Sprintf("Failed to regenerate project: %v", err)
		return compileResult
	}

	// Set cross-compilation environment
	env := append(os.Environ(),
		fmt.Sprintf("GOOS=%s", goos),
		fmt.Sprintf("GOARCH=%s", goarch),
		"CGO_ENABLED=0",
	)

	// Determine binary name
	binaryName := "main"
	if goos == "windows" {
		binaryName = "main.exe"
	}
	
	binaryPath := filepath.Join(projectPath, binaryName)

	// Compile the project
	startTime := time.Now()
	cmd := exec.Command("go", "build", "-o", binaryName, ".")
	cmd.Dir = projectPath
	cmd.Env = env

	output, err := cmd.CombinedOutput()
	compileResult.CompileTime = time.Since(startTime)

	if err != nil {
		compileResult.CompileError = fmt.Sprintf("Compilation failed: %v\nOutput: %s", err, string(output))
		return compileResult
	}

	// Get binary size
	if info, err := os.Stat(binaryPath); err == nil {
		compileResult.BinarySize = info.Size()
	}

	return compileResult
}

// testNativeCompilation tests compilation on the current platform
func testNativeCompilation(result PlatformTestResult) CompilationResult {
	return testCrossCompilation(result, runtime.GOOS, runtime.GOARCH)
}

// testNativeExecution_CLI tests CLI execution on the current platform
func testNativeExecution_CLI(result PlatformTestResult) ExecutionResult {
	execResult := ExecutionResult{}

	// Create a temporary project for execution testing
	tempDir, err := os.MkdirTemp("", "exec-test-*")
	if err != nil {
		execResult.ExecuteError = fmt.Sprintf("Failed to create exec temp dir: %v", err)
		return execResult
	}
	defer os.RemoveAll(tempDir)

	// Re-generate and compile the project
	config := getConfigFromBlueprint(result.Blueprint)
	config.Name = fmt.Sprintf("%s-exec", config.Name)
	
	gen := generator.New()
	projectPath := filepath.Join(tempDir, config.Name)
	
	_, err = gen.Generate(config, types.GenerationOptions{
		OutputPath: projectPath,
		DryRun:     false,
		NoGit:      true,
	})
	
	if err != nil {
		execResult.ExecuteError = fmt.Sprintf("Failed to regenerate project: %v", err)
		return execResult
	}

	// Compile the project
	binaryName := "main"
	if runtime.GOOS == "windows" {
		binaryName = "main.exe"
	}
	
	compileCmd := exec.Command("go", "build", "-o", binaryName, ".")
	compileCmd.Dir = projectPath
	
	if output, err := compileCmd.CombinedOutput(); err != nil {
		execResult.ExecuteError = fmt.Sprintf("Compilation for execution failed: %v\nOutput: %s", err, string(output))
		return execResult
	}

	// Execute the binary
	binaryPath := filepath.Join(projectPath, binaryName)
	startTime := time.Now()
	
	execCmd := exec.Command(binaryPath, "--version")
	execCmd.Dir = projectPath
	
	output, err := execCmd.CombinedOutput()
	execResult.ExecuteTime = time.Since(startTime)
	execResult.Output = string(output)

	if err != nil {
		// Try without --version flag
		execCmd = exec.Command(binaryPath, "--help")
		execCmd.Dir = projectPath
		
		output, err = execCmd.CombinedOutput()
		execResult.Output = string(output)
		
		if err != nil {
			execResult.ExecuteError = fmt.Sprintf("Execution failed: %v\nOutput: %s", err, string(output))
		}
	}

	return execResult
}

// checkPathIssues checks for platform-specific path issues
func checkPathIssues(filePaths []string, targetOS string) []string {
	var issues []string

	for _, filePath := range filePaths {
		// Check for Windows-specific issues
		if targetOS == "windows" {
			if strings.Contains(filePath, ":") && !filepath.IsAbs(filePath) {
				issues = append(issues, fmt.Sprintf("Invalid colon in relative path: %s", filePath))
			}
		}

		// Check for Unix-specific issues
		if targetOS != "windows" {
			if strings.Contains(filePath, "\\") {
				issues = append(issues, fmt.Sprintf("Backslash in path on Unix: %s", filePath))
			}
		}

		// Check for generally problematic characters
		if strings.ContainsAny(filePath, "<>|\"*?") {
			issues = append(issues, fmt.Sprintf("Problematic characters in path: %s", filePath))
		}

		// Check for excessively long paths
		if len(filePath) > 260 && targetOS == "windows" {
			issues = append(issues, fmt.Sprintf("Path too long for Windows: %s", filePath))
		}
	}

	return issues
}

// checkPermissionIssues checks for platform-specific permission issues
func checkPermissionIssues(filePaths []string, targetOS string) []string {
	var issues []string

	for _, filePath := range filePaths {
		// Check if executable files are properly marked
		if strings.HasSuffix(filePath, ".sh") || strings.HasSuffix(filePath, ".py") {
			if targetOS != "windows" {
				info, err := os.Stat(filePath)
				if err == nil {
					mode := info.Mode()
					if mode&0111 == 0 {
						issues = append(issues, fmt.Sprintf("Executable file without execute permission: %s", filePath))
					}
				}
			}
		}
	}

	return issues
}

// supportsCrossCompilation checks if cross-compilation is supported for the target
func supportsCrossCompilation(goos, goarch string) bool {
	// Go supports cross-compilation for most common platforms
	supportedPlatforms := map[string][]string{
		"linux":   {"amd64", "arm64", "386", "arm"},
		"windows": {"amd64", "386", "arm64"},
		"darwin":  {"amd64", "arm64"},
		"freebsd": {"amd64", "386", "arm64"},
	}

	if archs, exists := supportedPlatforms[goos]; exists {
		return contains(archs, goarch)
	}

	return false
}

// getConfigFromBlueprint returns a config for the given blueprint name
func getConfigFromBlueprint(blueprintName string) types.ProjectConfig {
	for _, blueprint := range testBlueprints {
		if blueprint.name == blueprintName {
			return blueprint.config
		}
	}
	
	// Return a default config if not found
	return types.ProjectConfig{
		Name:      "test-default",
		Module:    "github.com/test/default",
		Type:      "cli",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}
}

// isGoInstalled checks if Go is installed and available
func isGoInstalled() bool {
	cmd := exec.Command("go", "version")
	return cmd.Run() == nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// generateCompatibilityReport generates a comprehensive compatibility report
func generateCompatibilityReport(t *testing.T) {
	// Calculate platform summaries
	platformStats := make(map[string]*PlatformSummary)
	
	for _, result := range compatibilityReport.TestResults {
		platform := result.Platform
		
		if _, exists := platformStats[platform]; !exists {
			platformStats[platform] = &PlatformSummary{
				Platform: platform,
				Issues:   make([]string, 0),
			}
		}
		
		stats := platformStats[platform]
		stats.TotalTests++
		
		if result.Success {
			stats.SuccessfulTests++
			stats.AvgGenerationTime += result.GenerationTime
			stats.AvgCompileTime += result.CompileTime
			stats.AvgBinarySize += result.BinarySize
		} else {
			stats.FailedTests++
			if result.GenerationError != "" {
				stats.Issues = append(stats.Issues, fmt.Sprintf("Generation: %s", result.GenerationError))
			}
			if result.CompileError != "" {
				stats.Issues = append(stats.Issues, fmt.Sprintf("Compilation: %s", result.CompileError))
			}
		}
		
		// Add path and permission issues
		for _, issue := range result.PathIssues {
			stats.Issues = append(stats.Issues, fmt.Sprintf("Path: %s", issue))
		}
		for _, issue := range result.PermissionIssues {
			stats.Issues = append(stats.Issues, fmt.Sprintf("Permission: %s", issue))
		}
	}
	
	// Calculate averages
	for _, stats := range platformStats {
		if stats.SuccessfulTests > 0 {
			stats.AvgGenerationTime = time.Duration(int64(stats.AvgGenerationTime) / int64(stats.SuccessfulTests))
			stats.AvgCompileTime = time.Duration(int64(stats.AvgCompileTime) / int64(stats.SuccessfulTests))
			stats.AvgBinarySize = stats.AvgBinarySize / int64(stats.SuccessfulTests)
		}
		
		compatibilityReport.PlatformSummary[stats.Platform] = *stats
	}
	
	// Calculate compatibility score
	totalTests := len(compatibilityReport.TestResults)
	successfulTests := 0
	for _, result := range compatibilityReport.TestResults {
		if result.Success {
			successfulTests++
		}
	}
	
	if totalTests > 0 {
		compatibilityReport.CompatibilityScore = float64(successfulTests) / float64(totalTests) * 100
	}
	
	// Identify cross-platform issues
	identifyCrossPlatformIssues()
	
	// Generate report files
	if err := saveCrossPlatformReport(); err != nil {
		t.Errorf("Failed to save cross-platform report: %v", err)
	}
	
	t.Logf("Cross-platform compatibility score: %.1f%% (%d/%d tests passed)", 
		compatibilityReport.CompatibilityScore, successfulTests, totalTests)
}

// identifyCrossPlatformIssues analyzes results to identify platform-specific issues
func identifyCrossPlatformIssues() {
	issueMap := make(map[string][]string)
	
	for _, result := range compatibilityReport.TestResults {
		if !result.Success {
			key := result.GenerationError + result.CompileError + result.ExecuteError
			if key != "" {
				issueMap[key] = append(issueMap[key], result.Platform)
			}
		}
	}
	
	for issue, platforms := range issueMap {
		severity := "Low"
		if len(platforms) > 2 {
			severity = "High"
		} else if len(platforms) > 1 {
			severity = "Medium"
		}
		
		crossPlatformIssue := CrossPlatformIssue{
			Type:        "Compatibility",
			Description: issue,
			Platforms:   platforms,
			Severity:    severity,
			Suggestions: generateSuggestions(issue, platforms),
		}
		
		compatibilityReport.CrossPlatformIssues = append(compatibilityReport.CrossPlatformIssues, crossPlatformIssue)
	}
}

// generateSuggestions provides suggestions for fixing cross-platform issues
func generateSuggestions(issue string, platforms []string) []string {
	var suggestions []string
	
	if strings.Contains(issue, "path") {
		suggestions = append(suggestions, "Use filepath.Join() for cross-platform path handling")
		suggestions = append(suggestions, "Avoid hardcoded path separators")
	}
	
	if strings.Contains(issue, "permission") {
		suggestions = append(suggestions, "Set appropriate file permissions using os.Chmod()")
		suggestions = append(suggestions, "Check platform-specific permission requirements")
	}
	
	if strings.Contains(issue, "compilation") {
		suggestions = append(suggestions, "Check for platform-specific build constraints")
		suggestions = append(suggestions, "Verify CGO dependencies are compatible")
	}
	
	if len(platforms) == 1 && platforms[0] == "windows" {
		suggestions = append(suggestions, "Windows-specific issue - check file paths and line endings")
	}
	
	return suggestions
}

// saveCrossPlatformReport saves the compatibility report to files
func saveCrossPlatformReport() error {
	reportDir := "cross_platform_reports"
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return fmt.Errorf("failed to create report directory: %w", err)
	}
	
	timestamp := time.Now().Format("20060102_150405")
	
	// Save JSON report
	jsonPath := filepath.Join(reportDir, fmt.Sprintf("compatibility_report_%s.json", timestamp))
	if err := saveCrossPlatformJSONReport(jsonPath); err != nil {
		return fmt.Errorf("failed to save JSON report: %w", err)
	}
	
	// Save Markdown report
	mdPath := filepath.Join(reportDir, fmt.Sprintf("compatibility_report_%s.md", timestamp))
	if err := saveCrossPlatformMarkdownReport(mdPath); err != nil {
		return fmt.Errorf("failed to save Markdown report: %w", err)
	}
	
	fmt.Printf("Cross-platform compatibility reports generated:\n")
	fmt.Printf("  JSON: %s\n", jsonPath)
	fmt.Printf("  Markdown: %s\n", mdPath)
	
	return nil
}