package helpers

import (
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/pkg/types"
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

// GenerateProject generates a project into a temporary directory for testing.
// It returns the path to the generated project.
func GenerateProject(t TB, config types.ProjectConfig) string {
	t.Helper()

	outputDir := CreateTempDir(t)
	projectPath := filepath.Join(outputDir, config.Name)

	gen := generator.New()
	options := types.GenerationOptions{
		OutputPath: projectPath,
		DryRun:     false,
		NoGit:      true, // Skip git init for tests
		Verbose:    false,
	}

	_, err := gen.Generate(config, options)
	assert.NoError(t, err, "Project generation failed")

	return projectPath
}
