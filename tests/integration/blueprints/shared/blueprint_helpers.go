package shared

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// BlueprintTestConfig holds configuration for blueprint integration tests
type BlueprintTestConfig struct {
	BlueprintType string
	ProjectName   string
	ModulePath    string
	OutputDir     string
	Framework     string
	Logger        string
	NoGit         bool
	Timeout       time.Duration
}

// DefaultBlueprintConfig returns a default configuration for blueprint testing
func DefaultBlueprintConfig(blueprintType, projectName string) BlueprintTestConfig {
	return BlueprintTestConfig{
		BlueprintType: blueprintType,
		ProjectName:   projectName,
		ModulePath:    "github.com/test/" + projectName,
		Framework:     "gin",
		Logger:        "slog",
		NoGit:         true,
		Timeout:       5 * time.Minute,
	}
}

// GenerateProject generates a project using go-starter with the given configuration
func GenerateProject(t *testing.T, config BlueprintTestConfig) string {
	tmpDir := t.TempDir()
	projectPath := filepath.Join(tmpDir, config.ProjectName)

	// Build go-starter CLI
	originalDir, err := os.Getwd()
	require.NoError(t, err)

	// Navigate to project root (4 levels up from tests/integration/blueprints/shared)
	projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")
	goStarterPath := filepath.Join(tmpDir, "go-starter")

	buildCmd := exec.Command("go", "build", "-o", goStarterPath, ".")
	buildCmd.Dir = projectRoot
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Failed to build go-starter: %s", string(output))

	// Change to temporary directory for generation
	err = os.Chdir(tmpDir)
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	// Build generation command
	args := []string{
		"new", config.ProjectName,
		"--type=" + config.BlueprintType,
		"--module=" + config.ModulePath,
		"--framework=" + config.Framework,
		"--logger=" + config.Logger,
	}

	if config.NoGit {
		args = append(args, "--no-git")
	}

	// Execute generation command
	generateCmd := exec.Command(goStarterPath, args...)
	generateCmd.Dir = tmpDir
	output, err = generateCmd.CombinedOutput()
	require.NoError(t, err, "Project generation failed: %s", string(output))

	// Verify project was created
	require.DirExists(t, projectPath, "Generated project directory should exist")

	return projectPath
}

// CompileProject attempts to compile the generated project to verify it builds successfully
func CompileProject(t *testing.T, projectPath string) {
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Dir = projectPath
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Generated project should compile successfully: %s", string(output))
}

// ValidateFileStructure validates that expected files exist in the generated project
func ValidateFileStructure(t *testing.T, projectPath string, expectedFiles []string) {
	for _, expectedFile := range expectedFiles {
		fullPath := filepath.Join(projectPath, expectedFile)
		require.FileExists(t, fullPath, "Expected file should exist: %s", expectedFile)
	}
}

// ValidateDirectoryStructure validates that expected directories exist
func ValidateDirectoryStructure(t *testing.T, projectPath string, expectedDirs []string) {
	for _, expectedDir := range expectedDirs {
		fullPath := filepath.Join(projectPath, expectedDir)
		require.DirExists(t, fullPath, "Expected directory should exist: %s", expectedDir)
	}
}

// CountGeneratedFiles counts the total number of files in the generated project
func CountGeneratedFiles(t *testing.T, projectPath string) int {
	fileCount := 0
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileCount++
		}
		return nil
	})
	require.NoError(t, err, "Should be able to walk project directory")
	return fileCount
}

// ValidateGoMod checks that go.mod file has correct module name and Go version
func ValidateGoMod(t *testing.T, projectPath string, expectedModule string) {
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod file")

	contentStr := string(content)
	require.Contains(t, contentStr, "module "+expectedModule, "go.mod should contain correct module name")
	require.Contains(t, contentStr, "go 1.", "go.mod should specify Go version")
}

// ValidateProjectCompiles is a comprehensive validation that ensures the project compiles and basic commands work
func ValidateProjectCompiles(t *testing.T, projectPath string) {
	// First verify project compiles
	CompileProject(t, projectPath)

	// Check if it's a CLI project and try to run help command
	if hasCLIStructure(projectPath) {
		validateCLIProject(t, projectPath)
	}

	// Verify go mod tidy works
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = projectPath
	output, err := tidyCmd.CombinedOutput()
	require.NoError(t, err, "go mod tidy should succeed: %s", string(output))
}

// hasCLIStructure checks if the project appears to be a CLI application
func hasCLIStructure(projectPath string) bool {
	// Check for cmd directory or main.go with CLI patterns
	cmdDir := filepath.Join(projectPath, "cmd")
	if _, err := os.Stat(cmdDir); err == nil {
		return true
	}

	// Check main.go for CLI indicators
	mainGoPath := filepath.Join(projectPath, "main.go")
	if content, err := os.ReadFile(mainGoPath); err == nil {
		contentStr := string(content)
		return strings.Contains(contentStr, "cobra") || strings.Contains(contentStr, "flag")
	}

	return false
}

// validateCLIProject performs CLI-specific validation
func validateCLIProject(t *testing.T, projectPath string) {
	// Build the CLI binary
	binaryName := filepath.Base(projectPath)
	buildCmd := exec.Command("go", "build", "-o", binaryName, ".")
	buildCmd.Dir = projectPath
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "CLI project should build successfully: %s", string(output))

	binaryPath := filepath.Join(projectPath, binaryName)
	require.FileExists(t, binaryPath, "CLI binary should be created")

	// Try to run help command
	helpCmd := exec.Command("./"+binaryName, "--help")
	helpCmd.Dir = projectPath
	output, _ = helpCmd.CombinedOutput()

	// It's okay if help command fails, but it shouldn't panic
	outputStr := string(output)
	require.NotContains(t, outputStr, "panic:", "CLI help should not panic")
}

// ArchitectureValidator provides validation for specific architecture patterns
type ArchitectureValidator struct {
	ProjectPath string
}

// NewArchitectureValidator creates a new architecture validator for the given project
func NewArchitectureValidator(projectPath string) *ArchitectureValidator {
	return &ArchitectureValidator{ProjectPath: projectPath}
}

// ValidateCleanArchitecture validates Clean Architecture pattern compliance
func (av *ArchitectureValidator) ValidateCleanArchitecture(t *testing.T) {
	expectedDirs := []string{
		"internal/domain/entities",
		"internal/domain/usecases",
		"internal/infrastructure/repository",
		"internal/interfaces/controllers",
	}
	ValidateDirectoryStructure(t, av.ProjectPath, expectedDirs)

	// Validate that domain layer doesn't import infrastructure
	av.validateDependencyDirection(t, "internal/domain", []string{"internal/infrastructure", "internal/interfaces"})
}

// ValidateDDDArchitecture validates Domain-Driven Design pattern compliance
func (av *ArchitectureValidator) ValidateDDDArchitecture(t *testing.T) {
	expectedDirs := []string{
		"internal/domain/aggregates",
		"internal/domain/services",
		"internal/domain/valueobjects",
		"internal/infrastructure",
		"internal/application",
	}
	ValidateDirectoryStructure(t, av.ProjectPath, expectedDirs)
}

// ValidateHexagonalArchitecture validates Hexagonal Architecture pattern compliance
func (av *ArchitectureValidator) ValidateHexagonalArchitecture(t *testing.T) {
	expectedDirs := []string{
		"internal/core/domain",
		"internal/core/ports",
		"internal/adapters/primary",
		"internal/adapters/secondary",
	}
	ValidateDirectoryStructure(t, av.ProjectPath, expectedDirs)

	// Validate core doesn't depend on adapters
	av.validateDependencyDirection(t, "internal/core", []string{"internal/adapters"})
}

// ValidateStandardArchitecture validates standard layered architecture
func (av *ArchitectureValidator) ValidateStandardArchitecture(t *testing.T) {
	expectedDirs := []string{
		"internal/handlers",
		"internal/services",
		"internal/repository",
		"internal/models",
	}
	ValidateDirectoryStructure(t, av.ProjectPath, expectedDirs)
}

// validateDependencyDirection ensures certain directories don't import others
func (av *ArchitectureValidator) validateDependencyDirection(t *testing.T, sourceDir string, forbiddenImports []string) {
	sourcePath := filepath.Join(av.ProjectPath, sourceDir)

	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		contentStr := string(content)
		for _, forbidden := range forbiddenImports {
			require.NotContains(t, contentStr, `"`+forbidden,
				"File %s should not import %s (violates architecture)",
				strings.TrimPrefix(path, av.ProjectPath+"/"), forbidden)
		}

		return nil
	})

	require.NoError(t, err, "Should be able to validate dependency direction")
}
