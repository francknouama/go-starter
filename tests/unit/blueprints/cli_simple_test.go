package blueprints_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func init() {
	// Initialize templates filesystem for unit tests
	// We use DirFS pointing to the real blueprints directory
	wd, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory: " + err.Error())
	}

	// Navigate to project root and find blueprints directory
	projectRoot := wd
	for {
		templatesDir := filepath.Join(projectRoot, "blueprints")
		if _, err := os.Stat(templatesDir); err == nil {
			// Check if this directory actually contains template files
			// by looking for template.yaml files
			entries, err := os.ReadDir(templatesDir)
			if err == nil && len(entries) > 0 {
				// Check if any subdirectory contains template.yaml
				hasTemplates := false
				for _, entry := range entries {
					if entry.IsDir() {
						templateYaml := filepath.Join(templatesDir, entry.Name(), "template.yaml")
						if _, err := os.Stat(templateYaml); err == nil {
							hasTemplates = true
							break
						}
					}
				}

				if hasTemplates {
					templates.SetTemplatesFS(os.DirFS(templatesDir))
					return
				}
			}
		}

		// Move up one directory
		parentDir := filepath.Dir(projectRoot)
		if parentDir == projectRoot {
			panic("Could not find blueprints directory")
		}
		projectRoot = parentDir
	}
}

// TestCLISimpleBlueprint_TemplateMetadata tests the template.yaml metadata and structure
func TestCLISimpleBlueprint_TemplateMetadata(t *testing.T) {
	// Load the CLI-Simple blueprint template.yaml
	blueprintPath := findBlueprintPath(t, "cli-simple")
	templateYamlPath := filepath.Join(blueprintPath, "template.yaml")
	
	// Verify template.yaml exists
	helpers.AssertFileExists(t, templateYamlPath)
	
	// Read and parse template.yaml
	content := helpers.ReadFileContent(t, templateYamlPath)
	var templateData types.Template
	err := yaml.Unmarshal([]byte(content), &templateData)
	require.NoError(t, err, "template.yaml should be valid YAML")
	
	// Test metadata fields
	assert.Equal(t, "cli-simple", templateData.Name, "Template name should be cli-simple")
	assert.Equal(t, "cli", templateData.Type, "Template type should be cli")
	assert.Equal(t, "simple", templateData.Architecture, "Template architecture should be simple")
	assert.NotEmpty(t, templateData.Description, "Description should not be empty")
	assert.NotEmpty(t, templateData.Version, "Version should not be empty")
	
	// Test required variables
	expectedVariables := []string{"ProjectName", "ModulePath", "Author", "GoVersion"}
	actualVariables := make(map[string]bool)
	for _, variable := range templateData.Variables {
		actualVariables[variable.Name] = true
	}
	
	for _, expectedVar := range expectedVariables {
		assert.True(t, actualVariables[expectedVar], "Variable %s should be defined", expectedVar)
	}
	
	// Test files structure - CLI-Simple should have < 10 files
	assert.LessOrEqual(t, len(templateData.Files), 10, "CLI-Simple should have <= 10 files for simplicity")
	
	// Test essential files are included
	expectedFiles := map[string]string{
		"main.go.tmpl":        "main.go",
		"go.mod.tmpl":         "go.mod",
		"README.md.tmpl":      "README.md",
		"cmd/root.go.tmpl":    "cmd/root.go",
		"cmd/version.go.tmpl": "cmd/version.go",
		"config.go.tmpl":      "config.go",
	}
	
	actualFiles := make(map[string]string)
	for _, file := range templateData.Files {
		actualFiles[file.Source] = file.Destination
	}
	
	for source, expectedDest := range expectedFiles {
		actualDest, exists := actualFiles[source]
		assert.True(t, exists, "Essential file %s should be included", source)
		assert.Equal(t, expectedDest, actualDest, "File %s should map to %s", source, expectedDest)
	}
	
	// Test dependencies - should only have cobra (minimal dependencies)
	assert.LessOrEqual(t, len(templateData.Dependencies), 1, "CLI-Simple should have minimal dependencies")
	if len(templateData.Dependencies) > 0 {
		assert.Equal(t, "github.com/spf13/cobra", templateData.Dependencies[0].Module, "Should use Cobra framework")
	}
	
	// Test features - should have essential CLI features
	featureNames := make(map[string]bool)
	for _, feature := range templateData.Features {
		featureNames[feature.Name] = true
	}
	assert.True(t, featureNames["essential_cli"], "Should have essential_cli feature")
	assert.True(t, featureNames["shell_completion"], "Should have shell_completion feature")
}

// TestCLISimpleBlueprint_TemplateFileSyntax tests that all template files have valid Go template syntax
func TestCLISimpleBlueprint_TemplateFileSyntax(t *testing.T) {
	blueprintPath := findBlueprintPath(t, "cli-simple")
	
	// Define test template data
	testData := map[string]interface{}{
		"ProjectName": "test-cli",
		"ModulePath":  "github.com/test/test-cli",
		"Author":      "Test Author",
		"GoVersion":   "1.21",
	}
	
	// Define template functions that might be used
	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}
	
	// Find all .tmpl files
	templateFiles := helpers.FindFiles(t, blueprintPath, "*.tmpl")
	assert.NotEmpty(t, templateFiles, "Should find template files")
	
	for _, templateFile := range templateFiles {
		t.Run(filepath.Base(templateFile), func(t *testing.T) {
			// Read template content
			content := helpers.ReadFileContent(t, templateFile)
			
			// Parse template
			tmpl, err := template.New(filepath.Base(templateFile)).Funcs(funcMap).Parse(content)
			require.NoError(t, err, "Template %s should have valid syntax", templateFile)
			
			// Execute template to ensure it works
			var output strings.Builder
			err = tmpl.Execute(&output, testData)
			require.NoError(t, err, "Template %s should execute without errors", templateFile)
			
			// Verify output is not empty (unless it's meant to be)
			result := output.String()
			if !strings.Contains(templateFile, "gitignore") { // .gitignore might be empty
				assert.NotEmpty(t, result, "Template %s should produce non-empty output", templateFile)
			}
		})
	}
}

// TestCLISimpleBlueprint_MainGoTemplate tests the main.go template specifically
func TestCLISimpleBlueprint_MainGoTemplate(t *testing.T) {
	blueprintPath := findBlueprintPath(t, "cli-simple")
	mainGoPath := filepath.Join(blueprintPath, "main.go.tmpl")
	
	helpers.AssertFileExists(t, mainGoPath)
	content := helpers.ReadFileContent(t, mainGoPath)
	
	// Test essential characteristics of main.go
	assert.Contains(t, content, "package main", "Should be main package")
	assert.Contains(t, content, "log/slog", "Should use slog for logging")
	assert.Contains(t, content, "{{.ModulePath}}/cmd", "Should import cmd package")
	assert.Contains(t, content, "cmd.Execute()", "Should execute cmd")
	assert.NotContains(t, content, "factory", "Should not use factory patterns")
	assert.NotContains(t, content, "interface", "Should not use complex interfaces")
	
	// Test template execution
	testData := map[string]interface{}{
		"ModulePath": "github.com/test/simple-cli",
	}
	
	tmpl, err := template.New("main.go").Parse(content)
	require.NoError(t, err)
	
	var output strings.Builder
	err = tmpl.Execute(&output, testData)
	require.NoError(t, err)
	
	result := output.String()
	assert.Contains(t, result, "github.com/test/simple-cli/cmd", "Module path should be substituted")
}

// TestCLISimpleBlueprint_RootCommandTemplate tests the root.go command template
func TestCLISimpleBlueprint_RootCommandTemplate(t *testing.T) {
	blueprintPath := findBlueprintPath(t, "cli-simple")
	rootGoPath := filepath.Join(blueprintPath, "cmd", "root.go.tmpl")
	
	helpers.AssertFileExists(t, rootGoPath)
	content := helpers.ReadFileContent(t, rootGoPath)
	
	// Test essential CLI features
	assert.Contains(t, content, "cobra.Command", "Should use Cobra command")
	assert.Contains(t, content, "--quiet", "Should have quiet flag")
	assert.Contains(t, content, "--output", "Should have output format flag")
	assert.Contains(t, content, "json", "Should support JSON output")
	assert.Contains(t, content, "slog.Info", "Should use slog directly")
	assert.Contains(t, content, "CompletionOptions", "Should support shell completion")
	
	// Test simplicity - should not have complex patterns
	assert.NotContains(t, content, "factory", "Should not use factory patterns")
	// Note: map[string]interface{} is acceptable for CLI JSON output
	
	// Test template execution
	testData := map[string]interface{}{
		"ProjectName": "simple-cli",
	}
	
	tmpl, err := template.New("root.go").Parse(content)
	require.NoError(t, err)
	
	var output strings.Builder
	err = tmpl.Execute(&output, testData)
	require.NoError(t, err)
	
	result := output.String()
	assert.Contains(t, result, "simple-cli", "Project name should be substituted")
}

// TestCLISimpleBlueprint_ConfigTemplate tests the config.go template
func TestCLISimpleBlueprint_ConfigTemplate(t *testing.T) {
	blueprintPath := findBlueprintPath(t, "cli-simple")
	configGoPath := filepath.Join(blueprintPath, "config.go.tmpl")
	
	helpers.AssertFileExists(t, configGoPath)
	content := helpers.ReadFileContent(t, configGoPath)
	
	// Test configuration characteristics
	assert.Contains(t, content, "type Config struct", "Should define Config struct")
	assert.Contains(t, content, "LoadConfig", "Should have LoadConfig function")
	assert.Contains(t, content, "Validate", "Should have validation")
	assert.Contains(t, content, "os.Getenv", "Should use environment variables")
	assert.Contains(t, content, "slog.Level", "Should support slog levels")
	
	// Should not use complex config libraries
	assert.NotContains(t, content, "viper", "Should not use Viper")
	assert.NotContains(t, content, "cobra.Command", "Config should be separate from commands")
	
	// Test template functions (upper case conversion)
	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
	}
	
	testData := map[string]interface{}{
		"ProjectName": "test-cli",
	}
	
	tmpl, err := template.New("config.go").Funcs(funcMap).Parse(content)
	require.NoError(t, err)
	
	var output strings.Builder
	err = tmpl.Execute(&output, testData)
	require.NoError(t, err)
	
	result := output.String()
	assert.Contains(t, result, "TEST-CLI_", "Should use uppercase env var prefix")
}

// TestCLISimpleBlueprint_VersionCommandTemplate tests the version.go template
func TestCLISimpleBlueprint_VersionCommandTemplate(t *testing.T) {
	blueprintPath := findBlueprintPath(t, "cli-simple")
	versionGoPath := filepath.Join(blueprintPath, "cmd", "version.go.tmpl")
	
	helpers.AssertFileExists(t, versionGoPath)
	content := helpers.ReadFileContent(t, versionGoPath)
	
	// Test version command characteristics
	assert.Contains(t, content, "versionCmd", "Should define version command")
	assert.Contains(t, content, "cobra.Command", "Should use Cobra command")
	assert.Contains(t, content, "version", "Should handle version command")
	assert.Contains(t, content, "rootCmd.AddCommand", "Should add to root command")
	
	// Test template execution
	tmpl, err := template.New("version.go").Parse(content)
	require.NoError(t, err)
	
	var output strings.Builder
	err = tmpl.Execute(&output, map[string]interface{}{})
	require.NoError(t, err)
	
	result := output.String()
	assert.NotEmpty(t, result, "Should generate valid version command")
}

// TestCLISimpleBlueprint_GeneratedProjectStructure tests the overall generated project structure
func TestCLISimpleBlueprint_GeneratedProjectStructure(t *testing.T) {
	// Generate a test project
	config := types.ProjectConfig{
		Name:         "test-simple-cli",
		Module:       "github.com/test/test-simple-cli",
		Type:         "cli",
		Architecture: "simple",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}
	
	projectPath := helpers.GenerateProject(t, config)
	
	// Test file count - should be minimal (< 10 files)
	totalFiles := helpers.CountAllFiles(t, projectPath)
	assert.LessOrEqual(t, totalFiles, 10, "CLI-Simple should generate <= 10 files")
	
	// Test essential files exist
	essentialFiles := []string{
		"main.go",
		"go.mod",
		"README.md",
		"config.go",
		"cmd/root.go",
		"cmd/version.go",
	}
	
	for _, file := range essentialFiles {
		helpers.AssertFileExists(t, filepath.Join(projectPath, file))
	}
	
	// Test files that should NOT exist (complexity indicators)
	unwantedFiles := []string{
		"internal/logger/factory.go",
		"internal/logger/interface.go",
		"cmd/create.go",
		"cmd/update.go",
		"cmd/delete.go",
		"cmd/list.go",
	}
	
	for _, file := range unwantedFiles {
		helpers.AssertFileNotExists(t, filepath.Join(projectPath, file))
	}
	
	// Test directories that should NOT exist
	unwantedDirs := []string{
		"internal/logger",
		"configs",
	}
	
	for _, dir := range unwantedDirs {
		helpers.AssertDirectoryNotExists(t, filepath.Join(projectPath, dir))
	}
}

// TestCLISimpleBlueprint_GeneratedCodeCompilation tests that generated code compiles
func TestCLISimpleBlueprint_GeneratedCodeCompilation(t *testing.T) {
	// Generate a test project
	config := types.ProjectConfig{
		Name:         "test-compilation-cli",
		Module:       "github.com/test/test-compilation-cli",
		Type:         "cli",
		Architecture: "simple",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}
	
	projectPath := helpers.GenerateProject(t, config)
	
	// Test compilation
	helpers.AssertProjectCompiles(t, projectPath)
	
	// Test go.mod structure
	goModPath := filepath.Join(projectPath, "go.mod")
	helpers.AssertFileExists(t, goModPath)
	helpers.AssertFileContains(t, goModPath, "github.com/test/test-compilation-cli")
	helpers.AssertFileContains(t, goModPath, "go 1.21")
	helpers.AssertFileContains(t, goModPath, "github.com/spf13/cobra")
	
	// Should NOT contain complex dependencies
	helpers.AssertFileNotContains(t, goModPath, "github.com/spf13/viper")
	helpers.AssertFileNotContains(t, goModPath, "go.uber.org/zap")
	helpers.AssertFileNotContains(t, goModPath, "github.com/sirupsen/logrus")
}

// TestCLISimpleBlueprint_TemplateVariableSubstitution tests variable substitution in templates
func TestCLISimpleBlueprint_TemplateVariableSubstitution(t *testing.T) {
	// Generate project with specific values
	config := types.ProjectConfig{
		Name:         "my-awesome-tool",
		Module:       "github.com/myuser/my-awesome-tool",
		Type:         "cli",
		Architecture: "simple",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}
	
	projectPath := helpers.GenerateProject(t, config)
	
	// Test variable substitution in main.go
	mainGoPath := filepath.Join(projectPath, "main.go")
	mainGoContent := helpers.ReadFileContent(t, mainGoPath)
	assert.Contains(t, mainGoContent, "github.com/myuser/my-awesome-tool/cmd", "Module path should be substituted in imports")
	
	// Test variable substitution in root.go
	rootGoPath := filepath.Join(projectPath, "cmd", "root.go")
	rootGoContent := helpers.ReadFileContent(t, rootGoPath)
	assert.Contains(t, rootGoContent, "my-awesome-tool", "Project name should be substituted in command")
	
	// Test variable substitution in config.go
	configGoPath := filepath.Join(projectPath, "config.go")
	configGoContent := helpers.ReadFileContent(t, configGoPath)
	assert.Contains(t, configGoContent, "MY-AWESOME-TOOL_", "Uppercase project name should be used for env vars")
	
	// Test variable substitution in go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	goModContent := helpers.ReadFileContent(t, goModPath)
	assert.Contains(t, goModContent, "module github.com/myuser/my-awesome-tool", "Module should be substituted")
	assert.Contains(t, goModContent, "go 1.21", "Go version should be substituted")
}

// TestCLISimpleBlueprint_SlogIntegration tests slog integration and usage
func TestCLISimpleBlueprint_SlogIntegration(t *testing.T) {
	// Generate project
	config := types.ProjectConfig{
		Name:         "slog-test-cli",
		Module:       "github.com/test/slog-test-cli",
		Type:         "cli",
		Architecture: "simple",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}
	
	projectPath := helpers.GenerateProject(t, config)
	
	// Test slog usage in main.go
	mainGoPath := filepath.Join(projectPath, "main.go")
	mainGoContent := helpers.ReadFileContent(t, mainGoPath)
	assert.Contains(t, mainGoContent, "log/slog", "Should import slog")
	assert.Contains(t, mainGoContent, "slog.New", "Should create slog logger")
	assert.Contains(t, mainGoContent, "slog.SetDefault", "Should set default logger")
	assert.Contains(t, mainGoContent, "slog.Error", "Should use slog for errors")
	
	// Test slog usage in root.go
	rootGoPath := filepath.Join(projectPath, "cmd", "root.go")
	rootGoContent := helpers.ReadFileContent(t, rootGoPath)
	assert.Contains(t, rootGoContent, "log/slog", "Should import slog")
	assert.Contains(t, rootGoContent, "slog.Info", "Should use slog for info logging")
	
	// Test slog usage in config.go
	configGoPath := filepath.Join(projectPath, "config.go")
	configGoContent := helpers.ReadFileContent(t, configGoPath)
	assert.Contains(t, configGoContent, "log/slog", "Should import slog")
	assert.Contains(t, configGoContent, "slog.Level", "Should use slog.Level")
	assert.Contains(t, configGoContent, "slog.Error", "Should use slog for config errors")
	
	// Should NOT have logger factory or interfaces
	loggerFactoryPath := filepath.Join(projectPath, "internal", "logger", "factory.go")
	helpers.AssertFileNotExists(t, loggerFactoryPath)
	
	loggerInterfacePath := filepath.Join(projectPath, "internal", "logger", "interface.go")
	helpers.AssertFileNotExists(t, loggerInterfacePath)
}

// TestCLISimpleBlueprint_ComplexityValidation tests that the blueprint maintains simplicity
func TestCLISimpleBlueprint_ComplexityValidation(t *testing.T) {
	// Generate project
	config := types.ProjectConfig{
		Name:         "complexity-test-cli",
		Module:       "github.com/test/complexity-test-cli",
		Type:         "cli",
		Architecture: "simple",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}
	
	projectPath := helpers.GenerateProject(t, config)
	
	// Test file count (complexity metric)
	goFiles := helpers.CountFiles(t, projectPath, ".go")
	assert.LessOrEqual(t, goFiles, 6, "Should have <= 6 Go files for simplicity")
	
	totalFiles := helpers.CountAllFiles(t, projectPath)
	assert.LessOrEqual(t, totalFiles, 10, "Should have <= 10 total files for simplicity")
	
	// Test directory depth (complexity metric)
	maxDepth := helpers.GetMaxDirectoryDepth(t, projectPath)
	assert.LessOrEqual(t, maxDepth, 2, "Should have <= 2 directory levels for simplicity")
	
	// Test interface count (complexity metric)
	goFilesList := helpers.FindGoFiles(t, projectPath)
	totalInterfaces := 0
	for _, file := range goFilesList {
		content := helpers.ReadFileContent(t, file)
		interfaces := helpers.CountInterfaces(content)
		totalInterfaces += interfaces
	}
	assert.LessOrEqual(t, totalInterfaces, 2, "Should have <= 2 interfaces for simplicity (map[string]interface{} is acceptable)")
	
	// Test dependency count (complexity metric)
	goModContent := helpers.ReadFileContent(t, filepath.Join(projectPath, "go.mod"))
	depCount := helpers.CountDependencies(goModContent)
	assert.LessOrEqual(t, depCount, 1, "Should have <= 1 direct dependency for simplicity")
}

// Helper function to find blueprint path
func findBlueprintPath(t *testing.T, blueprintName string) string {
	t.Helper()
	
	// Get current working directory and navigate to project root
	wd, err := os.Getwd()
	require.NoError(t, err)
	
	// Navigate to project root by looking for go.mod
	projectRoot := wd
	for {
		goModPath := filepath.Join(projectRoot, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}
		parentDir := filepath.Dir(projectRoot)
		if parentDir == projectRoot {
			t.Fatal("Could not find project root (go.mod)")
		}
		projectRoot = parentDir
	}
	
	blueprintPath := filepath.Join(projectRoot, "blueprints", blueprintName)
	_, err = os.Stat(blueprintPath)
	require.NoError(t, err, "Blueprint %s should exist at %s", blueprintName, blueprintPath)
	
	return blueprintPath
}

// TestCLISimpleBlueprint_LoaderIntegration tests integration with the template loader
func TestCLISimpleBlueprint_LoaderIntegration(t *testing.T) {
	// Test that CLI-Simple blueprint can be loaded by the template system
	// Use the default loader which should use the embedded templates
	loader := templates.NewTemplateLoader()
	
	// Load the template
	template, err := loader.LoadTemplate("cli-simple")
	require.NoError(t, err, "Should load CLI-Simple template without errors")
	
	// Validate loaded template
	assert.Equal(t, "cli-simple", template.Name)
	assert.Equal(t, "cli", template.Type)
	assert.Equal(t, "simple", template.Architecture)
	assert.NotNil(t, template.Variables)
	assert.NotNil(t, template.Files)
	assert.NotNil(t, template.Dependencies)
}

// TestCLISimpleBlueprint_FileSystemStructure tests the blueprint's file system structure
func TestCLISimpleBlueprint_FileSystemStructure(t *testing.T) {
	blueprintPath := findBlueprintPath(t, "cli-simple")
	
	// Test required files exist
	requiredFiles := []string{
		"template.yaml",
		"main.go.tmpl",
		"go.mod.tmpl",
		"README.md.tmpl",
		"Makefile.tmpl",
		"config.go.tmpl",
		"cmd/root.go.tmpl",
		"cmd/version.go.tmpl",
	}
	
	for _, file := range requiredFiles {
		filePath := filepath.Join(blueprintPath, file)
		helpers.AssertFileExists(t, filePath)
	}
	
	// Test that unwanted files don't exist (complexity indicators)
	unwantedFiles := []string{
		"internal/logger/factory.go.tmpl",
		"internal/logger/interface.go.tmpl",
		"configs/config.yaml.tmpl",
		"cmd/create.go.tmpl",
		"cmd/update.go.tmpl",
		"cmd/delete.go.tmpl",
	}
	
	for _, file := range unwantedFiles {
		filePath := filepath.Join(blueprintPath, file)
		helpers.AssertFileNotExists(t, filePath)
	}
	
	// Test directory structure is simple
	cmdDir := filepath.Join(blueprintPath, "cmd")
	helpers.AssertDirectoryExists(t, cmdDir)
	
	// Should not have complex internal structure
	internalDir := filepath.Join(blueprintPath, "internal")
	helpers.AssertDirectoryNotExists(t, internalDir)
	
	configsDir := filepath.Join(blueprintPath, "configs")
	helpers.AssertDirectoryNotExists(t, configsDir)
}