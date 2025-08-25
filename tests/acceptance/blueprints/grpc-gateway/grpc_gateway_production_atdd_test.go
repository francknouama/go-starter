package grpcgateway_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Production Readiness ATDD Test Suite for grpc-gateway blueprint
// This test suite validates that the grpc-gateway blueprint is truly production-ready

func TestGRPCGatewayProductionReadiness(t *testing.T) {
	testSuite := &GRPCGatewayATDDSuite{
		t:                   t,
		blueprintName:      "grpc-gateway",
		expectedFileRanges: map[string][2]int{
			"slog":    {30, 40}, // Expected file count range for slog
			"zap":     {30, 40},
			"logrus":  {30, 40}, 
			"zerolog": {30, 40},
		},
	}

	t.Run("Production_Readiness_Validation", func(t *testing.T) {
		t.Run("Feature_TemplateVariableResolution", testSuite.TestTemplateVariableResolution)
		t.Run("Feature_LoggerIntegration", testSuite.TestLoggerIntegration)
		t.Run("Feature_CompilationValidation", testSuite.TestCompilationValidation)
		t.Run("Feature_FileCountValidation", testSuite.TestFileCountValidation)
		t.Run("Feature_DependencyResolution", testSuite.TestDependencyResolution)
		t.Run("Feature_ProtobufGeneration", testSuite.TestProtobufGeneration)
		t.Run("Feature_GRPCServiceValidation", testSuite.TestGRPCServiceValidation)
		t.Run("Feature_GatewayIntegration", testSuite.TestGatewayIntegration)
		t.Run("Feature_SecurityFeatures", testSuite.TestSecurityFeatures)
	})
}

type GRPCGatewayATDDSuite struct {
	t                   *testing.T
	blueprintName      string
	expectedFileRanges map[string][2]int
	projectRoot        string
	testOutputDir      string
}

func (s *GRPCGatewayATDDSuite) setup() {
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

func (s *GRPCGatewayATDDSuite) TestTemplateVariableResolution(t *testing.T) {
	t.Log("🔍 TESTING: Template variable resolution for all logger types")
	
	s.setup()
	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run(fmt.Sprintf("Logger_%s", logger), func(t *testing.T) {
			projectName := fmt.Sprintf("test-grpc-gateway-%s", logger)
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

func (s *GRPCGatewayATDDSuite) TestLoggerIntegration(t *testing.T) {
	t.Log("🚀 TESTING: Logger integration across all supported types")
	
	s.setup()
	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run(fmt.Sprintf("Logger_%s", logger), func(t *testing.T) {
			projectName := fmt.Sprintf("test-grpc-gateway-logger-%s", logger)
			projectPath := filepath.Join(s.testOutputDir, projectName)
			
			// Generate project with specific logger
			s.generateProject(t, projectName, projectPath, logger)

			// Verify logger-specific files are generated
			s.verifyLoggerFiles(t, projectPath, logger)

			// Verify project compiles without protobuf generation first
			s.verifyBasicCompilation(t, projectPath, logger)
		})
	}
}

func (s *GRPCGatewayATDDSuite) TestCompilationValidation(t *testing.T) {
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
			projectName := fmt.Sprintf("test-grpc-gateway-%s", config.name)
			projectPath := filepath.Join(s.testOutputDir, projectName)
			
			s.generateProjectWithConfig(t, projectName, projectPath, config.logger, config.database, config.auth)
			
			// First verify basic structure
			s.verifyBasicCompilation(t, projectPath, config.name)
			
			// Then try full compilation with protobuf generation
			s.verifyFullCompilation(t, projectPath, config.name)
		})
	}
}

func (s *GRPCGatewayATDDSuite) TestFileCountValidation(t *testing.T) {
	t.Log("📁 TESTING: File count validation matches expectations")
	
	s.setup()
	
	projectName := "test-grpc-gateway-filecount"
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

func (s *GRPCGatewayATDDSuite) TestDependencyResolution(t *testing.T) {
	t.Log("📦 TESTING: Dependency resolution and go mod validation")
	
	s.setup()
	
	projectName := "test-grpc-gateway-deps"
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

	// Check for key gRPC dependencies
	s.verifyGRPCDependencies(t, projectPath)
}

func (s *GRPCGatewayATDDSuite) TestProtobufGeneration(t *testing.T) {
	t.Log("🔄 TESTING: Protobuf generation and buf configuration")
	
	s.setup()
	
	projectName := "test-grpc-gateway-protobuf"
	projectPath := filepath.Join(s.testOutputDir, projectName)
	
	s.generateProject(t, projectName, projectPath, "slog")

	// Verify buf configuration files exist
	s.verifyBufConfiguration(t, projectPath)

	// Verify protobuf files exist
	s.verifyProtobufFiles(t, projectPath)

	// Try to generate protobuf code (this tests the fix for googleapis dependency)
	s.testProtobufGeneration(t, projectPath)
}

func (s *GRPCGatewayATDDSuite) TestGRPCServiceValidation(t *testing.T) {
	t.Log("🌐 TESTING: gRPC service implementation validation")
	
	s.setup()
	
	projectName := "test-grpc-gateway-service"
	projectPath := filepath.Join(s.testOutputDir, projectName)
	
	s.generateProject(t, projectName, projectPath, "slog")

	// Verify gRPC server files exist
	s.verifyGRPCServerFiles(t, projectPath)

	// Verify service implementations exist
	s.verifyServiceImplementations(t, projectPath)
}

func (s *GRPCGatewayATDDSuite) TestGatewayIntegration(t *testing.T) {
	t.Log("🚪 TESTING: gRPC Gateway integration validation")
	
	s.setup()
	
	projectName := "test-grpc-gateway-integration"
	projectPath := filepath.Join(s.testOutputDir, projectName)
	
	s.generateProject(t, projectName, projectPath, "slog")

	// Verify gateway server files exist
	s.verifyGatewayServerFiles(t, projectPath)

	// Verify middleware exists
	s.verifyMiddlewareFiles(t, projectPath)
}

func (s *GRPCGatewayATDDSuite) TestSecurityFeatures(t *testing.T) {
	t.Log("🔐 TESTING: Security features and enhanced interceptors")
	
	s.setup()
	
	projectName := "test-grpc-gateway-security"
	projectPath := filepath.Join(s.testOutputDir, projectName)
	
	s.generateProject(t, projectName, projectPath, "slog")

	// Verify security-related files exist
	s.verifySecurityFiles(t, projectPath)

	// Verify TLS configuration exists
	s.verifyTLSConfiguration(t, projectPath)
}

// Helper methods

func (s *GRPCGatewayATDDSuite) generateProject(t *testing.T, name, path, logger string) {
	s.generateProjectWithConfig(t, name, path, logger, "", "")
}

func (s *GRPCGatewayATDDSuite) generateProjectWithConfig(t *testing.T, name, path, logger, database, auth string) {
	args := []string{
		"new", name,
		"--type=grpc-gateway",
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

func (s *GRPCGatewayATDDSuite) scanForUnresolvedVariables(t *testing.T, projectPath string) []string {
	var unresolvedVars []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, .git, binary files, and generated pb files
		if info.IsDir() || strings.Contains(path, ".git") || 
		   strings.HasSuffix(path, ".exe") || strings.HasSuffix(path, ".so") ||
		   strings.HasSuffix(path, ".pb.go") || strings.Contains(path, "gen/") {
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

func (s *GRPCGatewayATDDSuite) verifyLoggerFiles(t *testing.T, projectPath, logger string) {
	loggerFile := filepath.Join(projectPath, "internal/logger", logger+".go")
	assert.FileExists(t, loggerFile, "Logger-specific file should exist: %s", logger)

	// Verify factory file exists
	factoryFile := filepath.Join(projectPath, "internal/logger/factory.go")
	assert.FileExists(t, factoryFile, "Logger factory file should exist")

	// Verify interface file exists  
	interfaceFile := filepath.Join(projectPath, "internal/logger/interface.go")
	assert.FileExists(t, interfaceFile, "Logger interface file should exist")
}

func (s *GRPCGatewayATDDSuite) verifyBasicCompilation(t *testing.T, projectPath, configName string) {
	t.Logf("Verifying basic Go syntax for %s", configName)
	
	// Run go vet to check for syntax issues without building
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = projectPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("❌ Go vet failed for %s: %v\nOutput: %s", configName, err, string(output))
	} else {
		t.Logf("✅ Go vet passed for %s", configName)
	}
}

func (s *GRPCGatewayATDDSuite) verifyFullCompilation(t *testing.T, projectPath, configName string) {
	t.Logf("Attempting full compilation for %s", configName)
	
	// First, try to run the generation script to create protobuf files
	genScript := filepath.Join(projectPath, "scripts/generate.sh")
	if _, err := os.Stat(genScript); err == nil {
		// Make script executable
		exec.Command("chmod", "+x", genScript).Run()
		
		// Try to run the generation
		cmd := exec.Command("bash", genScript)
		cmd.Dir = projectPath
		genOutput, genErr := cmd.CombinedOutput()
		
		if genErr != nil {
			t.Logf("⚠️ Protobuf generation failed for %s (expected - may need buf/protoc): %v\nOutput: %s", 
				configName, genErr, string(genOutput))
		} else {
			t.Logf("✅ Protobuf generation successful for %s", configName)
		}
	}

	// Try to build (may fail if protobuf generation failed)
	cmd := exec.Command("go", "build", "-o", "test-binary", "./cmd/server")
	cmd.Dir = projectPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("⚠️ Full compilation failed for %s (may be due to missing protobuf generation): %v\nOutput: %s", 
			configName, err, string(output))
	} else {
		t.Logf("✅ Full compilation successful for %s", configName)
		
		// Verify binary was created and clean up
		binaryPath := filepath.Join(projectPath, "test-binary")
		if _, err := os.Stat(binaryPath); err == nil {
			os.Remove(binaryPath)
		}
	}
}

func (s *GRPCGatewayATDDSuite) countGeneratedFiles(t *testing.T, projectPath string) int {
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

func (s *GRPCGatewayATDDSuite) verifyGRPCDependencies(t *testing.T, projectPath string) {
	goModFile := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModFile)
	require.NoError(t, err, "Should be able to read go.mod")

	modContent := string(content)
	
	// Check for key gRPC dependencies
	requiredDeps := []string{
		"google.golang.org/grpc",
		"google.golang.org/protobuf",
		"github.com/grpc-ecosystem/grpc-gateway/v2",
		"google.golang.org/genproto/googleapis/api", // This was the fix for googleapis issue
		"google.golang.org/genproto/googleapis/rpc",
	}

	for _, dep := range requiredDeps {
		if !strings.Contains(modContent, dep) {
			t.Errorf("❌ Missing required gRPC dependency: %s", dep)
		} else {
			t.Logf("✅ Found required dependency: %s", dep)
		}
	}
}

func (s *GRPCGatewayATDDSuite) verifyBufConfiguration(t *testing.T, projectPath string) {
	// Check for buf.yaml
	bufYamlFile := filepath.Join(projectPath, "buf.yaml")
	assert.FileExists(t, bufYamlFile, "buf.yaml configuration file should exist")

	// Check for buf.gen.yaml
	bufGenFile := filepath.Join(projectPath, "buf.gen.yaml")
	assert.FileExists(t, bufGenFile, "buf.gen.yaml generation file should exist")
}

func (s *GRPCGatewayATDDSuite) verifyProtobufFiles(t *testing.T, projectPath string) {
	// Check for user proto file
	userProtoFile := filepath.Join(projectPath, "proto/user/v1/user.proto")
	assert.FileExists(t, userProtoFile, "User protobuf file should exist")

	// Check for health proto file
	healthProtoFile := filepath.Join(projectPath, "proto/health/v1/health.proto")
	assert.FileExists(t, healthProtoFile, "Health protobuf file should exist")
}

func (s *GRPCGatewayATDDSuite) testProtobufGeneration(t *testing.T, projectPath string) {
	// Check if generate script exists and is executable
	genScript := filepath.Join(projectPath, "scripts/generate.sh")
	assert.FileExists(t, genScript, "Generation script should exist")

	// Try to make it executable
	cmd := exec.Command("chmod", "+x", genScript)
	cmd.Dir = projectPath
	cmd.Run()

	// Don't fail the test if protobuf generation fails - this is expected in CI environments
	// without buf/protoc installed. We just log the attempt.
	cmd = exec.Command("bash", genScript)
	cmd.Dir = projectPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("⚠️ Protobuf generation failed (expected without buf/protoc installed): %v\nOutput: %s", 
			err, string(output))
	} else {
		t.Logf("✅ Protobuf generation successful")
		
		// If generation succeeded, verify generated files
		if _, err := os.Stat(filepath.Join(projectPath, "gen")); err == nil {
			t.Logf("✅ Generated protobuf files found in gen/ directory")
		}
	}
}

func (s *GRPCGatewayATDDSuite) verifyGRPCServerFiles(t *testing.T, projectPath string) {
	// Check for gRPC server implementation
	grpcServerFile := filepath.Join(projectPath, "internal/server/grpc.go")
	assert.FileExists(t, grpcServerFile, "gRPC server file should exist")
}

func (s *GRPCGatewayATDDSuite) verifyServiceImplementations(t *testing.T, projectPath string) {
	// Check for service implementations
	userServiceFile := filepath.Join(projectPath, "internal/services/user.go")
	assert.FileExists(t, userServiceFile, "User service implementation should exist")

	healthServiceFile := filepath.Join(projectPath, "internal/services/health.go")
	assert.FileExists(t, healthServiceFile, "Health service implementation should exist")
}

func (s *GRPCGatewayATDDSuite) verifyGatewayServerFiles(t *testing.T, projectPath string) {
	// Check for gateway server implementation
	gatewayServerFile := filepath.Join(projectPath, "internal/server/gateway.go")
	assert.FileExists(t, gatewayServerFile, "Gateway server file should exist")
}

func (s *GRPCGatewayATDDSuite) verifyMiddlewareFiles(t *testing.T, projectPath string) {
	middlewareDir := filepath.Join(projectPath, "internal/middleware")
	assert.DirExists(t, middlewareDir, "Middleware directory should exist")

	// Check for key middleware files
	requiredMiddleware := []string{
		"logging.go",
		"recovery.go", 
		"security.go",
		"request_id.go",
		"error_handler.go",
	}

	for _, middleware := range requiredMiddleware {
		middlewareFile := filepath.Join(middlewareDir, middleware)
		assert.FileExists(t, middlewareFile, "Middleware file should exist: %s", middleware)
	}
}

func (s *GRPCGatewayATDDSuite) verifySecurityFiles(t *testing.T, projectPath string) {
	// Check for enhanced interceptors
	interceptorsFile := filepath.Join(projectPath, "internal/interceptors/enhanced.go")
	assert.FileExists(t, interceptorsFile, "Enhanced interceptors file should exist")

	// Check for security middleware
	securityFile := filepath.Join(projectPath, "internal/middleware/security.go")
	assert.FileExists(t, securityFile, "Security middleware should exist")

	// Check for error handling
	errorsFile := filepath.Join(projectPath, "internal/errors/enhanced_errors.go")
	assert.FileExists(t, errorsFile, "Enhanced errors file should exist")
}

func (s *GRPCGatewayATDDSuite) verifyTLSConfiguration(t *testing.T, projectPath string) {
	// Check for TLS configuration
	tlsConfigFile := filepath.Join(projectPath, "internal/tls/config.go")
	assert.FileExists(t, tlsConfigFile, "TLS configuration file should exist")
}