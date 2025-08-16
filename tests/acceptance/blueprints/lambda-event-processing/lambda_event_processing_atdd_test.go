package lambda_event_processing_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
)

func init() {
	// Initialize templates filesystem for testing
	// Get the path to the blueprints directory from the project root
	blueprintsPath := filepath.Join("..", "..", "..", "..", "blueprints")
	templates.SetTemplatesFS(os.DirFS(blueprintsPath))
}

func TestLambdaEventProcessingBlueprint(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "test-lambda-event-processing")

	// Create generator
	gen := generator.New()

	// Define project configuration
	config := types.ProjectConfig{
		Name:        "test-lambda-event-processing",
		Module:      "github.com/test/lambda-event-processing",
		Type:        "lambda-event-processing",
		GoVersion:   "1.21",
		Logger:      "slog",
		Author:      "Test Author",
		Email:       "test@example.com",
		License:     "MIT",
		Features:    &types.Features{}, // Initialize empty features to prevent nil pointer
	}

	// Generate project
	options := types.GenerationOptions{
		OutputPath: outputDir,
		DryRun:     false,
		NoGit:      true,
	}

	result, err := gen.Generate(config, options)
	require.NoError(t, err, "Generation should succeed")
	assert.True(t, result.Success, "Generation should be successful")

	// Verify expected files are created
	expectedFiles := []string{
		"main.go",
		"go.mod",
		"README.md",
		"internal/handler/event_router.go",
		"internal/handler/sqs_handler.go",
		"internal/domain/event.go",
		"internal/domain/message.go",
		"internal/config/config.go",
		"internal/observability/logger.go",
		"internal/observability/metrics.go",
		"internal/observability/tracing.go",
		"internal/performance/cold_start.go",
		"internal/security/secrets_manager.go",
	}

	t.Logf("Expected %d files, generated %d files", len(expectedFiles), len(result.FilesCreated))

	for _, expectedFile := range expectedFiles {
		filePath := filepath.Join(outputDir, expectedFile)
		assert.FileExists(t, filePath, "Expected file %s should exist", expectedFile)
		
		// Verify file is not empty
		stat, err := os.Stat(filePath)
		require.NoError(t, err)
		assert.Greater(t, stat.Size(), int64(0), "File %s should not be empty", expectedFile)
	}

	// Test that multi-select defaults are properly applied
	t.Run("MultiSelectDefaults", func(t *testing.T) {
		// Should include SQS handler (EventSources default includes "sqs")
		sqsHandlerPath := filepath.Join(outputDir, "internal/handler/sqs_handler.go")
		assert.FileExists(t, sqsHandlerPath, "SQS handler should be generated based on EventSources default")

		// Should include secrets manager (SecurityFeatures default includes "secrets-manager")
		secretsPath := filepath.Join(outputDir, "internal/security/secrets_manager.go")
		assert.FileExists(t, secretsPath, "Secrets manager should be generated based on SecurityFeatures default")

		// Should include metrics and tracing (ObservabilityLevel default is "advanced")
		metricsPath := filepath.Join(outputDir, "internal/observability/metrics.go")
		assert.FileExists(t, metricsPath, "Metrics should be generated based on ObservabilityLevel default")

		tracingPath := filepath.Join(outputDir, "internal/observability/tracing.go")
		assert.FileExists(t, tracingPath, "Tracing should be generated based on ObservabilityLevel default")
	})

	// Test compilation
	t.Run("Compilation", func(t *testing.T) {
		// Check if go.mod exists and has proper content
		goModPath := filepath.Join(outputDir, "go.mod")
		assert.FileExists(t, goModPath, "go.mod should exist")

		// Verify main.go contains expected imports and structure
		mainGoPath := filepath.Join(outputDir, "main.go")
		content, err := os.ReadFile(mainGoPath)
		require.NoError(t, err)
		
		mainGoContent := string(content)
		assert.Contains(t, mainGoContent, "package main", "main.go should have package main")
		assert.Contains(t, mainGoContent, "github.com/aws/aws-lambda-go/lambda", "Should import AWS Lambda package")
		assert.Contains(t, mainGoContent, config.Module, "Should import from the generated module path")
	})

	// Verify file count matches expectation (should be at least 12 files)
	assert.GreaterOrEqual(t, len(result.FilesCreated), 12, "Should generate at least 12 files for lambda-event-processing")
}

func TestLambdaEventProcessingConditionalGeneration(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "test-lambda-conditional")

	// Create generator
	gen := generator.New()

	// Test with different variable combinations to ensure conditional logic works
	config := types.ProjectConfig{
		Name:      "test-lambda-conditional",
		Module:    "github.com/test/lambda-conditional",
		Type:      "lambda-event-processing",
		GoVersion: "1.21",
		Logger:    "slog",
		Features:  &types.Features{}, // Initialize empty features to prevent nil pointer
		Variables: map[string]string{
			"EventSources":        `["sqs"]`,              // Only SQS
			"ObservabilityLevel":  "basic",                // Basic observability
			"SecurityFeatures":    `["input-validation"]`, // Only input validation
		},
	}

	// Generate project
	options := types.GenerationOptions{
		OutputPath: outputDir,
		DryRun:     false,
		NoGit:      true,
	}

	result, err := gen.Generate(config, options)
	require.NoError(t, err, "Generation should succeed")
	assert.True(t, result.Success, "Generation should be successful")

	// With basic observability, metrics and tracing should not be generated
	metricsPath := filepath.Join(outputDir, "internal/observability/metrics.go")
	tracingPath := filepath.Join(outputDir, "internal/observability/tracing.go")
	coldStartPath := filepath.Join(outputDir, "internal/performance/cold_start.go")

	// These files should NOT exist with basic observability
	assert.NoFileExists(t, metricsPath, "Metrics should not be generated with basic observability")
	assert.NoFileExists(t, tracingPath, "Tracing should not be generated with basic observability")
	assert.NoFileExists(t, coldStartPath, "Cold start optimization should not be generated with basic observability")

	// SQS handler should still exist
	sqsHandlerPath := filepath.Join(outputDir, "internal/handler/sqs_handler.go")
	assert.FileExists(t, sqsHandlerPath, "SQS handler should be generated")

	// Secrets manager should NOT exist (not in SecurityFeatures)
	secretsPath := filepath.Join(outputDir, "internal/security/secrets_manager.go")
	assert.NoFileExists(t, secretsPath, "Secrets manager should not be generated when not in SecurityFeatures")
}