package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/internal/templates"
)

// TB is an interface that testing.T and testing.B both implement
type TB interface {
	Helper()
	Cleanup(func())
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// CreateTempDir creates a temporary directory for testing and returns its path.
// It also registers a cleanup function with t.Cleanup().
func CreateTempDir(t TB) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "go-starter-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	return dir
}

// CleanupDir removes a directory and its contents.
func CleanupDir(t TB, dir string) {
	t.Helper()
	err := os.RemoveAll(dir)
	assert.NoError(t, err, "Failed to clean up directory")
}

// GenerateProjectInCwd generates a project using the current generator implementation in the current working directory.
func GenerateProjectInCwd(config *types.ProjectConfig) error {
	gen := generator.New()
	
	// Use current working directory to place the project
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	
	outputPath := filepath.Join(cwd, config.Name)
	
	// Generate the project
	_, err = gen.Generate(*config, types.GenerationOptions{
		OutputPath: outputPath,
		DryRun:     false,
		NoGit:      true,
		Verbose:    false,
	})
	
	return err
}

// CreateTempTestDir creates a temporary directory for testing that doesn't require cleanup interface
func CreateTempTestDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "go-starter-test-")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	return dir
}

// GenerateProject generates a project for testing with the given configuration and returns the project path
func GenerateProject(t *testing.T, config types.ProjectConfig) string {
	t.Helper()
	
	// Create temporary directory for the project
	tempDir := CreateTempTestDir(t)
	
	// Set output path to the temp directory
	projectPath := filepath.Join(tempDir, config.Name)
	
	// Generate the project
	gen := generator.New()
	_, err := gen.Generate(config, types.GenerationOptions{
		OutputPath: projectPath,
		DryRun:     false,
		NoGit:      true,
		Verbose:    false,
	})
	
	if err != nil {
		t.Fatalf("Failed to generate project: %v", err)
	}
	
	return projectPath
}

// InitializeTemplates initializes the template system for tests
func InitializeTemplates() error {
	// Find project root by looking for go.mod
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	
	projectRoot := currentDir
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			return os.ErrNotExist
		}
		projectRoot = parent
	}
	
	// Set templates FS
	blueprintsDir := filepath.Join(projectRoot, "blueprints")
	templates.SetTemplatesFS(os.DirFS(blueprintsDir))
	return nil
}
