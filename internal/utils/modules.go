package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// ModuleInfo represents Go module information
type ModuleInfo struct {
	Path    string
	Version string
	Dir     string
}

// InitGoModule initializes a new Go module in the specified directory
func InitGoModule(projectPath, modulePath string) error {
	cmd := exec.Command("go", "mod", "init", modulePath)
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize go module: %w", err)
	}

	return nil
}

// GoModTidy runs go mod tidy to clean up dependencies
func GoModTidy(projectPath string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go mod tidy: %w", err)
	}

	return nil
}

// GoModDownload downloads module dependencies
func GoModDownload(projectPath string) error {
	cmd := exec.Command("go", "mod", "download")
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to download dependencies: %w", err)
	}

	return nil
}

// GoGet adds a dependency to the module
func GoGet(projectPath string, packages ...string) error {
	if len(packages) == 0 {
		return fmt.Errorf("no packages specified")
	}

	args := append([]string{"get"}, packages...)
	cmd := exec.Command("go", args...)
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to go get packages: %w", err)
	}

	return nil
}

// GoVersion returns the installed Go version
func GoVersion() (string, error) {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get go version: %w", err)
	}

	// Parse version from output like "go version go1.21.0 darwin/amd64"
	versionRegex := regexp.MustCompile(`go(\d+\.\d+(?:\.\d+)?)`)
	matches := versionRegex.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return "", fmt.Errorf("failed to parse go version from: %s", string(output))
	}

	return matches[1], nil
}

// IsGoInstalled checks if Go is installed and available
func IsGoInstalled() bool {
	cmd := exec.Command("go", "version")
	return cmd.Run() == nil
}

// GetModulePath returns the module path from go.mod file
func GetModulePath(projectPath string) (string, error) {
	cmd := exec.Command("go", "list", "-m")
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get module path: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// HasGoMod checks if a directory contains a go.mod file
func HasGoMod(projectPath string) bool {
	goModPath := filepath.Join(projectPath, "go.mod")
	return FileExists(goModPath)
}

// GoModuleExists checks if a Go module exists at a path
func GoModuleExists(modulePath string) error {
	cmd := exec.Command("go", "list", "-m", modulePath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("module %s does not exist or is not accessible: %w", modulePath, err)
	}

	return nil
}

// GoList lists packages in the module
func GoList(projectPath string, pattern string) ([]string, error) {
	args := []string{"list"}
	if pattern != "" {
		args = append(args, pattern)
	} else {
		args = append(args, "./...")
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list packages: %w", err)
	}

	packages := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(packages) == 1 && packages[0] == "" {
		return []string{}, nil
	}

	return packages, nil
}

// GoBuild builds the Go project
func GoBuild(projectPath string, outputPath string, packages ...string) error {
	args := []string{"build"}

	if outputPath != "" {
		args = append(args, "-o", outputPath)
	}

	if len(packages) > 0 {
		args = append(args, packages...)
	} else {
		args = append(args, "./...")
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to build project: %w", err)
	}

	return nil
}

// GoTest runs tests for the Go project
func GoTest(projectPath string, packages ...string) error {
	args := []string{"test"}

	if len(packages) > 0 {
		args = append(args, packages...)
	} else {
		args = append(args, "./...")
	}

	cmd := exec.Command("go", args...)
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run tests: %w", err)
	}

	return nil
}

// GoFmt formats Go source code
func GoFmt(projectPath string) error {
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to format code: %w", err)
	}

	return nil
}

// GoVet runs go vet on the project
func GoVet(projectPath string) error {
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to vet code: %w", err)
	}

	return nil
}

// GoClean cleans build artifacts
func GoClean(projectPath string) error {
	cmd := exec.Command("go", "clean", "./...")
	cmd.Dir = projectPath

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clean project: %w", err)
	}

	return nil
}

// GoEnv gets Go environment variables
func GoEnv(vars ...string) (map[string]string, error) {
	args := append([]string{"env"}, vars...)
	cmd := exec.Command("go", args...)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get go env: %w", err)
	}

	env := make(map[string]string)
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := strings.Trim(parts[1], `"`)
			env[key] = value
		}
	}

	return env, nil
}

// GetGOPATH returns the GOPATH environment variable
func GetGOPATH() (string, error) {
	env, err := GoEnv("GOPATH")
	if err != nil {
		return "", err
	}

	return env["GOPATH"], nil
}

// GetGOROOT returns the GOROOT environment variable
func GetGOROOT() (string, error) {
	env, err := GoEnv("GOROOT")
	if err != nil {
		return "", err
	}

	return env["GOROOT"], nil
}

// ValidateGoInstallation checks if Go is properly installed and configured
func ValidateGoInstallation() error {
	if !IsGoInstalled() {
		return fmt.Errorf("go is not installed or not available in PATH")
	}

	version, err := GoVersion()
	if err != nil {
		return fmt.Errorf("failed to get Go version: %w", err)
	}

	// Check minimum Go version (1.18+)
	if !isValidGoVersion(version) {
		return fmt.Errorf("go version %s is not supported (minimum: 1.18)", version)
	}

	return nil
}

// isValidGoVersion checks if the Go version meets minimum requirements
func isValidGoVersion(version string) bool {
	// Simple version check - assumes format like "1.21.0" or "1.21"
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return false
	}

	if parts[0] != "1" {
		return false
	}

	// Check minor version - support Go 1.18+
	minorVersion := parts[1]

	// Convert to int for comparison
	var minor int
	if _, err := fmt.Sscanf(minorVersion, "%d", &minor); err != nil {
		return false
	}

	return minor >= 18
}

// GetOptimalGoVersion returns the best Go version to use for new projects
// It returns the current installed Go version in major.minor format,
// or falls back to a sensible default if detection fails
func GetOptimalGoVersion() string {
	// Supported versions in order of preference (latest first)
	supportedVersions := []string{"1.23", "1.22", "1.21"}
	
	// Try to get the current Go version
	currentVersion, err := GoVersion()
	if err != nil {
		// Fallback to a stable default if Go is not available or detection fails
		return "1.21"
	}

	// Convert to major.minor format (e.g., "1.21.5" -> "1.21")
	parts := strings.Split(currentVersion, ".")
	if len(parts) >= 2 {
		majorMinor := fmt.Sprintf("%s.%s", parts[0], parts[1])

		// Check if the current version is in our supported list
		for _, supportedVersion := range supportedVersions {
			if majorMinor == supportedVersion {
				return majorMinor
			}
		}
		
		// If current version is newer than supported (e.g., 1.24), use the latest supported
		if isValidGoVersion(majorMinor) {
			// Return the latest supported version for newer Go versions
			return "1.23"
		}
	}

	// If current version is too old or invalid, use minimum supported
	return "1.21"
}

// GenerateGoMod generates a go.mod file content
func GenerateGoMod(modulePath, goVersion string) string {
	if goVersion == "" {
		goVersion = "1.21"
	}

	return fmt.Sprintf(`module %s

go %s
`, modulePath, goVersion)
}

// CreateGoModule creates a new Go module with proper initialization
func CreateGoModule(projectPath, modulePath, goVersion string) error {
	// Initialize the module
	if err := InitGoModule(projectPath, modulePath); err != nil {
		return fmt.Errorf("failed to initialize module: %w", err)
	}

	// If a specific Go version is requested, update go.mod
	if goVersion != "" && goVersion != "1.21" {
		goModPath := filepath.Join(projectPath, "go.mod")
		content, err := ReadFile(goModPath)
		if err != nil {
			return fmt.Errorf("failed to read go.mod: %w", err)
		}

		// Replace the Go version
		goVersionRegex := regexp.MustCompile(`go \d+\.\d+(?:\.\d+)?`)
		updatedContent := goVersionRegex.ReplaceAllString(content, fmt.Sprintf("go %s", goVersion))

		if err := WriteFile(goModPath, updatedContent); err != nil {
			return fmt.Errorf("failed to update go.mod: %w", err)
		}
	}

	return nil
}

// AddDependencies adds multiple dependencies to a Go module
func AddDependencies(projectPath string, dependencies []string) error {
	if len(dependencies) == 0 {
		return nil
	}

	// Add dependencies
	if err := GoGet(projectPath, dependencies...); err != nil {
		return fmt.Errorf("failed to add dependencies: %w", err)
	}

	// Clean up
	if err := GoModTidy(projectPath); err != nil {
		return fmt.Errorf("failed to tidy module: %w", err)
	}

	return nil
}

// CheckModulePath validates if a module path is accessible
func CheckModulePath(modulePath string) error {
	// Basic format validation
	if modulePath == "" {
		return fmt.Errorf("module path cannot be empty")
	}

	// Check if it looks like a valid module path
	if !strings.Contains(modulePath, "/") {
		return fmt.Errorf("module path should contain at least one slash")
	}

	parts := strings.Split(modulePath, "/")
	if len(parts) < 2 {
		return fmt.Errorf("module path should have at least domain and path components")
	}

	// Check if domain part looks valid
	domain := parts[0]
	if !strings.Contains(domain, ".") {
		return fmt.Errorf("module path should start with a domain")
	}

	return nil
}
