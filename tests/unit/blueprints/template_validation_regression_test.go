package blueprints_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDDDValueObjectStringConversion_RegressionTest tests that the DDD template 
// properly calls String() methods on value objects (Issue #149 fix validation)
func TestDDDValueObjectStringConversion_RegressionTest(t *testing.T) {
	blueprintPath := "../../../blueprints/web-api-ddd"
	
	// Test the DTO file has correct string conversions
	dtoFile := filepath.Join(blueprintPath, "internal/application/user/dto.go.tmpl")
	content, err := os.ReadFile(dtoFile)
	require.NoError(t, err, "Should be able to read DTO template file")
	
	dtoContent := string(content)
	
	// Check that value objects are converted to strings properly
	assert.Contains(t, dtoContent, "{{.DomainName}}Entity.Name().String()",
		"DTO should call .String() on UserName value object")
	assert.Contains(t, dtoContent, "{{.DomainName}}Entity.Email().String()",
		"DTO should call .String() on EmailAddress value object")
	assert.Contains(t, dtoContent, "{{.DomainName}}Entity.Description().String()",
		"DTO should call .String() on UserDescription value object")
	
	// Ensure we're not trying to use value objects directly as strings
	assert.NotContains(t, dtoContent, "Name:        {{.DomainName}}Entity.Name(),",
		"Should not use value object directly without .String() conversion")
	assert.NotContains(t, dtoContent, "Email:       {{.DomainName}}Entity.Email(),",
		"Should not use value object directly without .String() conversion")
	assert.NotContains(t, dtoContent, "Description: {{.DomainName}}Entity.Description(),",
		"Should not use value object directly without .String() conversion")
}

func TestDDDRepositoryValueObjectStringConversion_RegressionTest(t *testing.T) {
	blueprintPath := "../../../blueprints/web-api-ddd"
	
	// Test the repository file has correct string conversions
	repoFile := filepath.Join(blueprintPath, "internal/infrastructure/persistence/user_repository.go.tmpl")
	content, err := os.ReadFile(repoFile)
	require.NoError(t, err, "Should be able to read repository template file")
	
	repoContent := string(content)
	
	// Check that entityToModel method properly converts value objects to strings
	assert.Contains(t, repoContent, "Name:        entity.Name().String()",
		"Repository should call .String() on UserName value object")
	assert.Contains(t, repoContent, "Email:       entity.Email().String()",
		"Repository should call .String() on EmailAddress value object")
	assert.Contains(t, repoContent, "Description: entity.Description().String()",
		"Repository should call .String() on UserDescription value object")
	
	// Ensure we're not trying to use value objects directly as strings
	assert.NotContains(t, repoContent, "Name:        entity.Name(),",
		"Should not use value object directly without .String() conversion")
	assert.NotContains(t, repoContent, "Email:       entity.Email(),",
		"Should not use value object directly without .String() conversion")
	assert.NotContains(t, repoContent, "Description: entity.Description(),",
		"Should not use value object directly without .String() conversion")
}

func TestEventDrivenBasicFilesExist_RegressionTest(t *testing.T) {
	blueprintPath := "../../../blueprints/event-driven"
	
	// Test that basic template files exist (Issue #149 fix validation)
	requiredFiles := []string{
		"go.mod.tmpl",
		"go.sum.tmpl",
		"main.go.tmpl",
		"README.md.tmpl",
		"Dockerfile.tmpl",
		".dockerignore.tmpl",
		".gitignore.tmpl",
		"docker-compose.yml.tmpl",
		"Makefile.tmpl",
	}
	
	for _, file := range requiredFiles {
		filePath := filepath.Join(blueprintPath, file)
		_, err := os.Stat(filePath)
		assert.NoError(t, err, "Required template file should exist: %s", file)
	}
}

func TestMicroserviceScriptFilesExist_RegressionTest(t *testing.T) {
	blueprintPath := "../../../blueprints/microservice-standard"
	
	// Test that script files exist (Issue #149 fix validation)
	requiredScripts := []string{
		"scripts/generate.sh.tmpl",
		"scripts/test.sh.tmpl",
	}
	
	for _, file := range requiredScripts {
		filePath := filepath.Join(blueprintPath, file)
		_, err := os.Stat(filePath)
		assert.NoError(t, err, "Required script file should exist: %s", file)
	}
	
	// Check that the scripts have proper content
	generateScript := filepath.Join(blueprintPath, "scripts/generate.sh.tmpl")
	content, err := os.ReadFile(generateScript)
	require.NoError(t, err, "Should be able to read generate script")
	
	scriptContent := string(content)
	assert.Contains(t, scriptContent, "#!/bin/bash", "Script should have bash shebang")
	assert.Contains(t, scriptContent, "echo", "Script should have echo statements")
	assert.Contains(t, scriptContent, "go build", "Script should build the project")
}

func TestMicroserviceDockerComposeExists_RegressionTest(t *testing.T) {
	blueprintPath := "../../../blueprints/microservice-standard"
	
	dockerComposeFile := filepath.Join(blueprintPath, "docker-compose.yml.tmpl")
	_, err := os.Stat(dockerComposeFile)
	assert.NoError(t, err, "docker-compose.yml.tmpl should exist")
	
	content, err := os.ReadFile(dockerComposeFile)
	require.NoError(t, err, "Should be able to read docker-compose template")
	
	composeContent := string(content)
	assert.Contains(t, composeContent, "version:", "Docker Compose should have version")
	assert.Contains(t, composeContent, "services:", "Docker Compose should have services")
	assert.Contains(t, composeContent, "{{.ProjectName}}:", "Docker Compose should reference project name")
}

// TestValueObjectsHaveStringMethods_RegressionTest tests that DDD value objects
// implement the String() method (validation for the compilation fix)
func TestValueObjectsHaveStringMethods_RegressionTest(t *testing.T) {
	blueprintPath := "../../../blueprints/web-api-ddd"
	
	valueObjectsFile := filepath.Join(blueprintPath, "internal/domain/user/value_objects.go.tmpl")
	content, err := os.ReadFile(valueObjectsFile)
	require.NoError(t, err, "Should be able to read value objects template file")
	
	voContent := string(content)
	
	// Check that value objects have String() methods
	assert.Contains(t, voContent, "func (n UserName) String() string",
		"UserName should have String() method")
	assert.Contains(t, voContent, "func (e EmailAddress) String() string",
		"EmailAddress should have String() method") 
	assert.Contains(t, voContent, "func (d UserDescription) String() string",
		"UserDescription should have String() method")
	assert.Contains(t, voContent, "func (s Status) String() string",
		"Status should have String() method")
}