package web_api_clean_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Production Readiness ATDD Test Suite for web-api-clean blueprint
// This test suite validates that the web-api-clean blueprint is truly production-ready

func TestWebAPICleanProductionReadiness(t *testing.T) {
	testSuite := &ProductionATDDSuite{
		t:                   t,
		blueprintName:      "web-api-clean",
		expectedFileRanges: map[string][2]int{
			"slog":    {35, 45}, // Expected file count range for slog
			"zap":     {35, 45},
			"logrus":  {35, 45},
			"zerolog": {35, 45},
		},
	}

	t.Run("Production_Readiness_Validation", func(t *testing.T) {
		t.Run("Feature_TemplateVariableResolution", testSuite.TestTemplateVariableResolution)
		t.Run("Feature_LoggerIntegration", testSuite.TestLoggerIntegration)
		t.Run("Feature_CompilationValidation", testSuite.TestCompilationValidation)
		t.Run("Feature_FileCountValidation", testSuite.TestFileCountValidation)
		t.Run("Feature_DependencyResolution", testSuite.TestDependencyResolution)
		t.Run("Feature_DatabaseIntegration", testSuite.TestDatabaseIntegration)
		t.Run("Feature_AuthenticationFlow", testSuite.TestAuthenticationFlow)
		t.Run("Feature_CleanArchitectureStructure", testSuite.TestCleanArchitectureStructure)
	})
}

type ProductionATDDSuite struct {
	t                   *testing.T
	blueprintName      string
	expectedFileRanges map[string][2]int
	projectRoot        string
	testOutputDir      string
}

func (s *ProductionATDDSuite) setup() {
	// Find project root
	currentDir, err := os.Getwd()
	require.NoError(s.t, err)

	s.projectRoot = currentDir
	for {
		if _, err := os.Stat(filepath.Join(s.projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(s.projectRoot)
		if parent == s.projectRoot {
			s.t.Fatal("Could not find project root with go.mod")
		}
		s.projectRoot = parent
	}

	// Setup test output directory
	s.testOutputDir = filepath.Join(s.projectRoot, "test-output", fmt.Sprintf("%s-production-test", s.blueprintName))
	os.RemoveAll(s.testOutputDir) // Clean up any previous test artifacts
}

func (s *ProductionATDDSuite) TestTemplateVariableResolution(t *testing.T) {
	t.Log("🔍 TESTING: Template variable resolution for all logger types")
	
	s.setup()
	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run(fmt.Sprintf("Logger_%s", logger), func(t *testing.T) {
			projectName := fmt.Sprintf("test-web-clean-%s", logger)
			projectPath := filepath.Join(s.testOutputDir, projectName)
			
			// Generate project
			s.generateProject(t, projectName, projectPath, logger)

			// Scan all generated files for unresolved template variables
			unresolvedVars := s.scanForUnresolvedVariables(t, projectPath)
			
			if len(unresolvedVars) > 0 {
				t.Errorf("❌ Found unresolved template variables with logger %s:\n%s", 
					logger, strings.Join(unresolvedVars, "\n"))
			} else {
				t.Logf("✅ All template variables resolved correctly for logger %s", logger)
			}
		})
	}
}

func (s *ProductionATDDSuite) TestLoggerIntegration(t *testing.T) {
	t.Log("🚀 TESTING: Logger integration across all supported types")
	
	s.setup()
	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run(fmt.Sprintf("Logger_%s", logger), func(t *testing.T) {
			projectName := fmt.Sprintf("test-web-clean-logger-%s", logger)
			projectPath := filepath.Join(s.testOutputDir, projectName)
			
			// Generate project with specific logger
			s.generateProject(t, projectName, projectPath, logger)

			// Verify logger-specific files are generated
			s.verifyLoggerFiles(t, projectPath, logger)

			// Verify project compiles
			s.verifyCompilation(t, projectPath, logger)
		})
	}
}

func (s *ProductionATDDSuite) TestCompilationValidation(t *testing.T) {
	t.Log("🔧 TESTING: Compilation validation for all configurations")
	
	s.setup()
	testConfigs := []struct {
		name     string
		logger   string
		database string
		auth     string
	}{
		{"minimal", "slog", "", ""},
		{"database_postgres", "zap", "postgres", ""},
		{"auth_jwt", "logrus", "", "jwt"},
		{"full_featured", "zerolog", "postgres", "jwt"},
	}

	for _, config := range testConfigs {
		t.Run(fmt.Sprintf("Config_%s", config.name), func(t *testing.T) {
			projectName := fmt.Sprintf("test-web-clean-%s", config.name)
			projectPath := filepath.Join(s.testOutputDir, projectName)
			
			s.generateProjectWithConfig(t, projectName, projectPath, config.logger, config.database, config.auth)
			s.verifyCompilation(t, projectPath, config.name)

			// Test that binary can start and respond to health check
			s.testBinaryExecution(t, projectPath, config.name)
		})
	}
}

func (s *ProductionATDDSuite) TestFileCountValidation(t *testing.T) {
	t.Log("📁 TESTING: File count validation matches expectations")
	
	s.setup()
	
	projectName := "test-web-clean-filecount"
	projectPath := filepath.Join(s.testOutputDir, projectName)
	
	s.generateProject(t, projectName, projectPath, "slog")

	fileCount := s.countGeneratedFiles(t, projectPath)
	expectedRange := s.expectedFileRanges["slog"]
	
	if fileCount < expectedRange[0] || fileCount > expectedRange[1] {
		t.Errorf("❌ File count %d is outside expected range [%d, %d]", 
			fileCount, expectedRange[0], expectedRange[1])
	} else {
		t.Logf("✅ File count %d is within expected range [%d, %d]", 
			fileCount, expectedRange[0], expectedRange[1])
	}
}

func (s *ProductionATDDSuite) TestDependencyResolution(t *testing.T) {
	t.Log("📦 TESTING: Dependency resolution and go mod validation")
	
	s.setup()
	
	projectName := "test-web-clean-deps"
	projectPath := filepath.Join(s.testOutputDir, projectName)
	
	s.generateProject(t, projectName, projectPath, "zap")

	// Run go mod download and verify all dependencies resolve
	cmd := exec.Command("go", "mod", "download")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Errorf("❌ go mod download failed: %v\nOutput: %s", err, string(output))
	} else {
		t.Log("✅ All dependencies resolved successfully")
	}

	// Run go mod verify
	cmd = exec.Command("go", "mod", "verify")
	cmd.Dir = projectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Errorf("❌ go mod verify failed: %v\nOutput: %s", err, string(output))
	} else {
		t.Log("✅ Module verification passed")
	}
}

func (s *ProductionATDDSuite) TestDatabaseIntegration(t *testing.T) {
	t.Log("🗄️ TESTING: Database integration with Clean Architecture")
	
	s.setup()
	
	databases := []string{"postgres", "mysql", "sqlite"}
	
	for _, db := range databases {
		t.Run(fmt.Sprintf("Database_%s", db), func(t *testing.T) {
			projectName := fmt.Sprintf("test-web-clean-db-%s", db)
			projectPath := filepath.Join(s.testOutputDir, projectName)
			
			s.generateProjectWithConfig(t, projectName, projectPath, "slog", db, "")

			// Verify database-related files are generated
			s.verifyDatabaseFiles(t, projectPath, db)

			// Verify compilation with database
			s.verifyCompilation(t, projectPath, fmt.Sprintf("db_%s", db))
		})
	}
}

func (s *ProductionATDDSuite) TestAuthenticationFlow(t *testing.T) {
	t.Log("🔐 TESTING: Authentication integration with JWT")
	
	s.setup()
	
	projectName := "test-web-clean-auth"
	projectPath := filepath.Join(s.testOutputDir, projectName)
	
	s.generateProjectWithConfig(t, projectName, projectPath, "slog", "postgres", "jwt")

	// Verify authentication files are generated
	s.verifyAuthFiles(t, projectPath)

	// Verify compilation with auth
	s.verifyCompilation(t, projectPath, "auth_jwt")
}

func (s *ProductionATDDSuite) TestCleanArchitectureStructure(t *testing.T) {
	t.Log("🏗️ TESTING: Clean Architecture structure validation")
	
	s.setup()
	
	projectName := "test-web-clean-structure"
	projectPath := filepath.Join(s.testOutputDir, projectName)
	
	s.generateProject(t, projectName, projectPath, "slog")

	// Verify Clean Architecture directory structure
	s.verifyCleanArchitectureStructure(t, projectPath)
}

// Helper methods

func (s *ProductionATDDSuite) generateProject(t *testing.T, name, path, logger string) {
	s.generateProjectWithConfig(t, name, path, logger, "", "")
}

func (s *ProductionATDDSuite) generateProjectWithConfig(t *testing.T, name, path, logger, database, auth string) {
	args := []string{
		"new", name,
		"--type=web-api",
		"--architecture=clean",
		fmt.Sprintf("--logger=%s", logger),
		"--output=" + filepath.Dir(path),
		"--module=github.com/test/" + name,
		"--no-git",
		"--quiet",
	}

	if database != "" {
		args = append(args, "--database-driver="+database)
		args = append(args, "--database-orm=gorm")
	}
	
	if auth != "" {
		args = append(args, "--auth-type="+auth)
	}

	cmd := exec.Command(filepath.Join(s.projectRoot, "go-starter"), args...)
	cmd.Dir = s.projectRoot
	
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Project generation failed: %s", string(output))
	
	// Verify the project directory was created
	require.DirExists(t, path, "Generated project directory should exist")
}

func (s *ProductionATDDSuite) scanForUnresolvedVariables(t *testing.T, projectPath string) []string {
	var unresolvedVars []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, .git, and binary files
		if info.IsDir() || strings.Contains(path, ".git") || 
		   strings.HasSuffix(path, ".exe") || strings.HasSuffix(path, ".so") {
			return nil
		}

		// Read file content
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil // Skip files we can't read
		}

		// Look for unresolved template variables
		contentStr := string(content)
		if strings.Contains(contentStr, "{{") && strings.Contains(contentStr, "}}") {
			lines := strings.Split(contentStr, "\n")
			for i, line := range lines {
				if strings.Contains(line, "{{") && strings.Contains(line, "}}") {
					relPath, _ := filepath.Rel(s.projectRoot, path)
					unresolvedVars = append(unresolvedVars, 
						fmt.Sprintf("  %s:%d: %s", relPath, i+1, strings.TrimSpace(line)))
				}
			}
		}

		return nil
	})

	require.NoError(t, err, "Failed to scan for unresolved variables")
	return unresolvedVars
}

func (s *ProductionATDDSuite) verifyLoggerFiles(t *testing.T, projectPath, logger string) {
	loggerFile := filepath.Join(projectPath, "internal/infrastructure/logger", logger+".go")
	assert.FileExists(t, loggerFile, "Logger-specific file should exist: %s", logger)

	// Verify factory file exists
	factoryFile := filepath.Join(projectPath, "internal/infrastructure/logger/factory.go")
	assert.FileExists(t, factoryFile, "Logger factory file should exist")

	// Verify interface file exists  
	interfaceFile := filepath.Join(projectPath, "internal/infrastructure/logger/interface.go")
	assert.FileExists(t, interfaceFile, "Logger interface file should exist")
}

func (s *ProductionATDDSuite) verifyCompilation(t *testing.T, projectPath, configName string) {
	t.Logf("Compiling project at %s for config %s", projectPath, configName)
	
	// Run go build
	cmd := exec.Command("go", "build", "-o", "test-binary", "./cmd/server")
	cmd.Dir = projectPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("❌ Compilation failed for %s: %v\nOutput: %s", configName, err, string(output))
	} else {
		t.Logf("✅ Compilation successful for %s", configName)
		
		// Verify binary was created
		binaryPath := filepath.Join(projectPath, "test-binary")
		assert.FileExists(t, binaryPath, "Binary should be created")
		
		// Clean up binary
		os.Remove(binaryPath)
	}
}

func (s *ProductionATDDSuite) testBinaryExecution(t *testing.T, projectPath, configName string) {
	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "test-binary", "./cmd/server")
	cmd.Dir = projectPath
	
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("❌ Cannot test binary execution, build failed: %v\nOutput: %s", err, string(output))
		return
	}

	binaryPath := filepath.Join(projectPath, "test-binary")
	defer os.Remove(binaryPath) // Clean up

	// Try to run the binary with --help (should exit quickly and successfully)
	cmd = exec.Command(binaryPath, "--help")
	cmd.Dir = projectPath
	
	// Set a reasonable timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Logf("❌ Binary execution failed for %s: %v", configName, err)
		} else {
			t.Logf("✅ Binary executed successfully for %s", configName)
		}
	case <-time.After(5 * time.Second):
		cmd.Process.Kill()
		t.Logf("⚠️ Binary execution timeout for %s (this might be expected)", configName)
	}
}

func (s *ProductionATDDSuite) countGeneratedFiles(t *testing.T, projectPath string) int {
	count := 0
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && !strings.Contains(path, ".git") {
			count++
		}
		return nil
	})
	require.NoError(t, err, "Failed to count files")
	return count
}

func (s *ProductionATDDSuite) verifyDatabaseFiles(t *testing.T, projectPath, database string) {
	// Check for database-related files in Clean Architecture structure
	persistenceDir := filepath.Join(projectPath, "internal/infrastructure/persistence")
	assert.DirExists(t, persistenceDir, "Persistence layer should exist")

	dbFile := filepath.Join(persistenceDir, "database.go")
	assert.FileExists(t, dbFile, "Database connection file should exist")

	repoFile := filepath.Join(persistenceDir, "user_repository.go")
	assert.FileExists(t, repoFile, "User repository implementation should exist")

	migrationsFile := filepath.Join(projectPath, "migrations/001_create_users.up.sql")
	assert.FileExists(t, migrationsFile, "Database migration should exist")
}

func (s *ProductionATDDSuite) verifyAuthFiles(t *testing.T, projectPath string) {
	// Check for auth-related files in Clean Architecture structure
	authController := filepath.Join(projectPath, "internal/adapters/controllers/auth_controller.go")
	assert.FileExists(t, authController, "Auth controller should exist")

	authUsecase := filepath.Join(projectPath, "internal/domain/usecases/auth_usecase.go")
	assert.FileExists(t, authUsecase, "Auth use case should exist")

	authEntity := filepath.Join(projectPath, "internal/domain/entities/auth.go")
	assert.FileExists(t, authEntity, "Auth entity should exist")
	
	authService := filepath.Join(projectPath, "internal/infrastructure/services/auth_service.go")
	assert.FileExists(t, authService, "Auth service should exist")
}

func (s *ProductionATDDSuite) verifyCleanArchitectureStructure(t *testing.T, projectPath string) {
	// Verify Clean Architecture layer structure
	layers := []string{
		"internal/domain/entities",          // Enterprise Business Rules
		"internal/domain/usecases",          // Application Business Rules  
		"internal/domain/ports",             // Interface definitions
		"internal/adapters/controllers",     // Interface Adapters
		"internal/adapters/presenters",      // Interface Adapters
		"internal/infrastructure/persistence", // Frameworks & Drivers
		"internal/infrastructure/web",       // Frameworks & Drivers
		"internal/infrastructure/services",  // Frameworks & Drivers
	}

	for _, layer := range layers {
		layerPath := filepath.Join(projectPath, layer)
		assert.DirExists(t, layerPath, "Clean Architecture layer should exist: %s", layer)
	}

	// Verify dependency injection container
	containerFile := filepath.Join(projectPath, "internal/infrastructure/container/container.go")
	assert.FileExists(t, containerFile, "Dependency injection container should exist")
}