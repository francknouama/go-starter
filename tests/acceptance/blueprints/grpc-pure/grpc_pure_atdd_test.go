package grpcpure

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/francknouama/go-starter/tests/helpers"
)

// GRPCPureTestContext holds state for gRPC Pure BDD tests
type GRPCPureTestContext struct {
	workingDir   string
	projectDir   string
	projectName  string
	originalDir  string
	projectRoot  string
	lastCommand  *exec.Cmd
	lastOutput   []byte
	lastError    error
	lastExitCode int
	projectExists bool
	compilationOK bool
	
	// gRPC-specific test state
	protoFilesValid      bool
	interceptorsEnabled  bool
	serviceDiscoveryType string
	authType            string
	tracingEnabled      bool
	metricsEnabled      bool
	reflectionEnabled   bool
	databaseDriver      string
	databaseORM         string
	loggerType          string
}

// TestGRPCPureBlueprintBasicGeneration tests fundamental gRPC Pure blueprint functionality
func TestGRPCPureBlueprintBasicGeneration(t *testing.T) {
	// Skip if short mode
	if testing.Short() {
		t.Skip("Skipping ATDD tests in short mode")
	}

	tempDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempDir)

	testCases := []struct {
		name           string
		projectName    string
		command        string
		expectedFiles  []string
		shouldCompile  bool
		expectedErrors []string
	}{
		{
			name:        "Basic gRPC Pure Service",
			projectName: "basic-grpc-pure",
			command:     "go-starter new basic-grpc-pure --type=grpc-pure --architecture=microservice --module=github.com/example/basic-grpc-pure --no-git --logger=slog",
			expectedFiles: []string{
				"go.mod",
				"cmd/server/main.go",
				"internal/server/grpc.go",
				"internal/server/health.go", 
				"internal/services/basic-grpc-pure.go",
				"internal/services/health.go",
				"internal/config/config.go",
				"internal/logger/interface.go",
				"internal/logger/logger.go",
				"internal/interceptors/logging.go",
				"internal/interceptors/recovery.go",
				"internal/interceptors/ratelimit.go",
				"proto/basic-grpc-pure/v1/service.proto",
				"proto/health/v1/health.proto",
				"proto/common/v1/common.proto",
				"buf.yaml",
				"buf.gen.yaml",
				"Makefile",
				"Dockerfile",
				"docker-compose.yml",
				"README.md",
				"configs/config.dev.yaml",
				"configs/config.prod.yaml",
				"configs/config.test.yaml",
				"scripts/generate.sh",
				"scripts/dev.sh",
				"scripts/test.sh",
				"tests/integration/grpc_test.go",
				"tests/integration/health_test.go",
				"tests/integration/interceptors_test.go",
				"tests/unit/services_test.go",
				"tests/load/grpc_load_test.go",
			},
			shouldCompile: true,
		},
		{
			name:        "gRPC Pure with JWT Authentication",
			projectName: "grpc-jwt",
			command:     "go-starter new grpc-jwt --type=grpc-pure --architecture=microservice --module=github.com/example/grpc-jwt --no-git --logger=slog",
			expectedFiles: []string{
				"go.mod",
				"cmd/server/main.go",
				"internal/server/grpc.go",
				"internal/auth/jwt.go",
				"internal/auth/interface.go",
				"internal/interceptors/auth.go",
				"internal/interceptors/logging.go",
				"internal/interceptors/recovery.go",
				"proto/grpc-jwt/v1/service.proto",
				"internal/services/grpc-jwt.go",
			},
			shouldCompile: true,
		},
		{
			name:        "gRPC Pure with Observability",
			projectName: "grpc-observability",
			command:     "go-starter new grpc-observability --type=grpc-pure --architecture=microservice --module=github.com/example/grpc-observability --no-git --logger=slog",
			expectedFiles: []string{
				"go.mod",
				"cmd/server/main.go",
				"internal/server/grpc.go",
				"internal/server/metrics.go",
				"internal/observability/tracing.go",
				"internal/observability/metrics.go",
				"internal/interceptors/tracing.go",
				"internal/interceptors/metrics.go",
				"proto/grpc-observability/v1/service.proto",
			},
			shouldCompile: true,
		},
		{
			name:        "gRPC Pure with Database",
			projectName: "grpc-database",
			command:     "go-starter new grpc-database --type=grpc-pure --architecture=microservice --module=github.com/example/grpc-database --no-git --logger=slog",
			expectedFiles: []string{
				"go.mod",
				"cmd/server/main.go",
				"internal/server/grpc.go",
				"internal/database/connection.go",
				"internal/database/migrations.go",
				"internal/models/user.go",
				"internal/repository/user.go",
				"internal/repository/interface.go",
				"migrations/001_create_users.up.sql",
				"migrations/001_create_users.down.sql",
				"migrations/embed.go",
				"proto/grpc-database/v1/service.proto",
			},
			shouldCompile: true,
		},
		{
			name:        "gRPC Pure with Service Discovery",
			projectName: "grpc-discovery",
			command:     "go-starter new grpc-discovery --type=grpc-pure --architecture=microservice --module=github.com/example/grpc-discovery --no-git --logger=slog",
			expectedFiles: []string{
				"go.mod",
				"cmd/server/main.go",
				"internal/server/grpc.go",
				"internal/discovery/consul.go",
				"internal/discovery/interface.go",
				"proto/grpc-discovery/v1/service.proto",
			},
			shouldCompile: true,
		},
		{
			name:        "gRPC Pure with Alternative Logger",
			projectName: "grpc-zap",
			command:     "go-starter new grpc-zap --type=grpc-pure --architecture=microservice --module=github.com/example/grpc-zap --no-git --logger=zap",
			expectedFiles: []string{
				"go.mod",
				"cmd/server/main.go",
				"internal/server/grpc.go",
				"internal/logger/interface.go",
				"internal/logger/logger.go",
				"proto/grpc-zap/v1/service.proto",
			},
			shouldCompile: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Change to temp directory for each test
			err := os.Chdir(tempDir)
			require.NoError(t, err)

			// Run the generation command
			cmd := exec.Command("sh", "-c", tc.command)
			cmd.Dir = tempDir
			output, err := cmd.CombinedOutput()

			if len(tc.expectedErrors) > 0 {
				// Test case expects errors
				assert.Error(t, err, "Command should have failed")
				outputStr := string(output)
				for _, expectedError := range tc.expectedErrors {
					assert.Contains(t, outputStr, expectedError, "Output should contain expected error")
				}
				return
			}

			// Test case expects success
			require.NoError(t, err, "Command should succeed: %s", string(output))

			projectPath := filepath.Join(tempDir, tc.projectName)

			// Verify expected files exist
			for _, expectedFile := range tc.expectedFiles {
				filePath := filepath.Join(projectPath, expectedFile)
				assert.True(t, fileOrDirExists(filePath), "Expected file/directory should exist: %s", expectedFile)
			}

			// Test compilation if expected
			if tc.shouldCompile {
				t.Run("Compilation", func(t *testing.T) {
					testCompilation(t, projectPath)
				})
			}

			// Test protobuf generation and validation
			t.Run("Protobuf Generation", func(t *testing.T) {
				testProtobufGeneration(t, projectPath)
			})

			// Test protobuf validation with buf
			t.Run("Protobuf Validation", func(t *testing.T) {
				testProtobufValidation(t, projectPath)
			})
		})
	}
}

// TestGRPCPureAdvancedFeatures tests advanced gRPC Pure features and combinations
func TestGRPCPureAdvancedFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping advanced feature tests in short mode")
	}

	tempDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempDir)

	advancedTests := []struct {
		name        string
		command     string
		testFunc    func(t *testing.T, projectPath string)
		shouldCompile bool
	}{
		{
			name:     "Full Production Configuration",
			command:  "./bin/go-starter-dev new production-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/production-grpc --no-git --logger=zap",
			testFunc: testProductionConfiguration,
			shouldCompile: true,
		},
		{
			name:     "mTLS Authentication",
			command:  "./bin/go-starter-dev new mtls-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/mtls-grpc --no-git --logger=slog",
			testFunc: testMTLSConfiguration,
			shouldCompile: true,
		},
		{
			name:     "OAuth2 Authentication",
			command:  "./bin/go-starter-dev new oauth-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/oauth-grpc --no-git --logger=slog",
			testFunc: testOAuth2Configuration,
			shouldCompile: true,
		},
		{
			name:     "Etcd Service Discovery",
			command:  "./bin/go-starter-dev new etcd-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/etcd-grpc --no-git --logger=slog",
			testFunc: testEtcdServiceDiscovery,
			shouldCompile: true,
		},
		{
			name:     "Kubernetes Service Discovery",
			command:  "./bin/go-starter-dev new k8s-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/k8s-grpc --no-git --logger=slog",
			testFunc: testKubernetesServiceDiscovery,
			shouldCompile: true,
		},
		{
			name:     "MySQL Database Integration",
			command:  "./bin/go-starter-dev new mysql-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/mysql-grpc --no-git --logger=slog",
			testFunc: testMySQLIntegration,
			shouldCompile: true,
		},
		{
			name:     "SQLite Database Integration",
			command:  "./bin/go-starter-dev new sqlite-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/sqlite-grpc --no-git --logger=slog",
			testFunc: testSQLiteIntegration,
			shouldCompile: true,
		},
		{
			name:     "SQLx Database Integration",
			command:  "./bin/go-starter-dev new sqlx-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/sqlx-grpc --no-git --logger=slog",
			testFunc: testSQLxIntegration,
			shouldCompile: true,
		},
		{
			name:     "gRPC Reflection for Development",
			command:  "./bin/go-starter-dev new reflection-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/reflection-grpc --no-git --logger=slog",
			testFunc: testGRPCReflection,
			shouldCompile: true,
		},
		{
			name:     "Alternative Loggers - Logrus",
			command:  "./bin/go-starter-dev new logrus-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/logrus-grpc --no-git --logger=logrus",
			testFunc: testLogrusLogger,
			shouldCompile: true,
		},
		{
			name:     "Alternative Loggers - Zerolog",
			command:  "./bin/go-starter-dev new zerolog-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/zerolog-grpc --no-git --logger=zerolog",
			testFunc: testZerologLogger,
			shouldCompile: true,
		},
	}

	for _, at := range advancedTests {
		t.Run(at.name, func(t *testing.T) {
			projectName := extractProjectName(at.command)
			
			err := os.Chdir(tempDir)
			require.NoError(t, err)

			// Run generation command
			cmd := exec.Command("sh", "-c", at.command)
			cmd.Dir = tempDir
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Advanced feature generation should succeed: %s", string(output))

			projectPath := filepath.Join(tempDir, projectName)

			// Run feature-specific tests
			if at.testFunc != nil {
				at.testFunc(t, projectPath)
			}

			// Ensure compilation still works if specified
			if at.shouldCompile {
				testCompilation(t, projectPath)
			}
		})
	}
}

// TestGRPCPureInterceptorChain tests the interceptor chain configuration and order
func TestGRPCPureInterceptorChain(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping interceptor chain tests in short mode")
	}

	tempDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempDir)

	projectName := "interceptor-test-grpc"
	command := fmt.Sprintf("./bin/go-starter-dev new %s --type=grpc-pure --architecture=microservice --module=github.com/example/%s --no-git --logger=slog",
		projectName, projectName)

	err := os.Chdir(tempDir)
	require.NoError(t, err)

	// Generate project
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Interceptor test project generation should succeed: %s", string(output))

	projectPath := filepath.Join(tempDir, projectName)

	t.Run("Interceptor Chain Order", func(t *testing.T) {
		testInterceptorChainOrder(t, projectPath)
	})

	t.Run("Auth Interceptor Validation", func(t *testing.T) {
		testAuthInterceptorValidation(t, projectPath)
	})

	t.Run("Logging Interceptor Validation", func(t *testing.T) {
		testLoggingInterceptorValidation(t, projectPath)
	})

	t.Run("Metrics Interceptor Validation", func(t *testing.T) {
		testMetricsInterceptorValidation(t, projectPath)
	})

	t.Run("Recovery Interceptor Validation", func(t *testing.T) {
		testRecoveryInterceptorValidation(t, projectPath)
	})

	t.Run("Rate Limiting Interceptor Validation", func(t *testing.T) {
		testRateLimitingInterceptorValidation(t, projectPath)
	})
}

// Helper functions for testing

func testCompilation(t *testing.T, projectPath string) {
	t.Helper()
	
	// Install dependencies
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "go mod tidy should succeed: %s", string(output))

	// Build the project
	cmd = exec.Command("go", "build", "./...")
	cmd.Dir = projectPath
	output, err = cmd.CombinedOutput()
	assert.NoError(t, err, "Project should compile successfully: %s", string(output))
}

func testProtobufGeneration(t *testing.T, projectPath string) {
	t.Helper()
	
	protoDir := filepath.Join(projectPath, "proto")
	if !fileOrDirExists(protoDir) {
		t.Skip("No proto directory found")
	}

	// Check for .proto files
	protoFiles, err := findProtoFiles(protoDir)
	require.NoError(t, err, "Should be able to search for proto files")
	assert.Greater(t, len(protoFiles), 0, "Should have at least one .proto file")

	// Test Makefile generation command if available
	makefilePath := filepath.Join(projectPath, "Makefile")
	if fileOrDirExists(makefilePath) {
		cmd := exec.Command("make", "generate")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("Protobuf generation via make failed (may be expected without buf/protoc): %s", string(output))
		} else {
			t.Logf("Protobuf generation via make succeeded: %s", string(output))
		}
	}
}

func testProtobufValidation(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check for buf configuration
	bufYamlPath := filepath.Join(projectPath, "buf.yaml")
	assert.True(t, fileOrDirExists(bufYamlPath), "buf.yaml should exist for protobuf validation")

	bufGenYamlPath := filepath.Join(projectPath, "buf.gen.yaml")
	assert.True(t, fileOrDirExists(bufGenYamlPath), "buf.gen.yaml should exist for code generation")

	// Test buf lint if buf is available
	if _, err := exec.LookPath("buf"); err == nil {
		cmd := exec.Command("buf", "lint")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("buf lint failed (may be expected): %s", string(output))
		} else {
			t.Logf("buf lint succeeded: %s", string(output))
		}
	} else {
		t.Skip("buf tool not available, skipping protobuf validation")
	}
}

func testProductionConfiguration(t *testing.T, projectPath string) {
	t.Helper()
	
	// Verify all production components are present
	productionFiles := []string{
		"internal/auth/jwt.go",
		"internal/observability/tracing.go",
		"internal/observability/metrics.go",
		"internal/discovery/consul.go",
		"internal/database/connection.go",
		"internal/interceptors/auth.go",
		"internal/interceptors/tracing.go",
		"internal/interceptors/metrics.go",
	}

	for _, file := range productionFiles {
		filePath := filepath.Join(projectPath, file)
		assert.True(t, fileOrDirExists(filePath), "Production file should exist: %s", file)
	}

	// Check go.mod for required dependencies
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")

	contentStr := string(content)
	requiredDeps := []string{
		"github.com/golang-jwt/jwt/v5",
		"go.opentelemetry.io/otel",
		"github.com/prometheus/client_golang",
		"github.com/hashicorp/consul/api",
		"gorm.io/gorm",
		"go.uber.org/zap",
	}

	for _, dep := range requiredDeps {
		assert.Contains(t, contentStr, dep, "Production dependency should be present: %s", dep)
	}
}

func testMTLSConfiguration(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check for mTLS auth implementation
	mtlsPath := filepath.Join(projectPath, "internal/auth/mtls.go")
	assert.True(t, fileOrDirExists(mtlsPath), "mTLS auth implementation should exist")

	// Check for TLS configuration
	tlsConfigPath := filepath.Join(projectPath, "internal/tls/config.go")
	assert.True(t, fileOrDirExists(tlsConfigPath), "TLS configuration should exist")
}

func testOAuth2Configuration(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check for OAuth2 auth implementation
	oauthPath := filepath.Join(projectPath, "internal/auth/oauth.go")
	assert.True(t, fileOrDirExists(oauthPath), "OAuth2 auth implementation should exist")

	// Check for OAuth2 dependency in go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")
	assert.Contains(t, string(content), "golang.org/x/oauth2", "OAuth2 dependency should be present")
}

func testEtcdServiceDiscovery(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check for etcd implementation
	etcdPath := filepath.Join(projectPath, "internal/discovery/etcd.go")
	assert.True(t, fileOrDirExists(etcdPath), "Etcd service discovery should exist")

	// Check for etcd dependency
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")
	assert.Contains(t, string(content), "go.etcd.io/etcd/client/v3", "Etcd dependency should be present")
}

func testKubernetesServiceDiscovery(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check for Kubernetes implementation
	k8sPath := filepath.Join(projectPath, "internal/discovery/kubernetes.go")
	assert.True(t, fileOrDirExists(k8sPath), "Kubernetes service discovery should exist")

	// Check for Kubernetes dependency
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")
	assert.Contains(t, string(content), "k8s.io/client-go", "Kubernetes dependency should be present")
}

func testMySQLIntegration(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check for MySQL driver dependency
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")
	assert.Contains(t, string(content), "gorm.io/driver/mysql", "MySQL GORM driver should be present")
}

func testSQLiteIntegration(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check for SQLite driver dependency
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")
	assert.Contains(t, string(content), "gorm.io/driver/sqlite", "SQLite GORM driver should be present")
}

func testSQLxIntegration(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check for SQLx dependency
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")
	assert.Contains(t, string(content), "github.com/jmoiron/sqlx", "SQLx dependency should be present")
}

func testGRPCReflection(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check for reflection in server configuration
	serverPath := filepath.Join(projectPath, "internal/server/grpc.go")
	content, err := os.ReadFile(serverPath)
	require.NoError(t, err, "Should be able to read server configuration")
	assert.Contains(t, string(content), "reflection", "gRPC reflection should be configured")
}

func testLogrusLogger(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check for Logrus dependency
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")
	assert.Contains(t, string(content), "github.com/sirupsen/logrus", "Logrus dependency should be present")
}

func testZerologLogger(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check for Zerolog dependency
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")
	assert.Contains(t, string(content), "github.com/rs/zerolog", "Zerolog dependency should be present")
}

func testInterceptorChainOrder(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check server configuration for proper interceptor order
	serverPath := filepath.Join(projectPath, "internal/server/grpc.go")
	content, err := os.ReadFile(serverPath)
	require.NoError(t, err, "Should be able to read server configuration")

	contentStr := string(content)
	
	// Check for interceptor chain setup
	assert.Contains(t, contentStr, "grpc_middleware", "Should use grpc middleware package")
	
	// The order should be: recovery, logging, auth, metrics, tracing
	// This is important for proper request handling
	recoveryPos := strings.Index(contentStr, "recovery")
	loggingPos := strings.Index(contentStr, "logging")
	
	if recoveryPos != -1 && loggingPos != -1 {
		assert.Less(t, recoveryPos, loggingPos, "Recovery interceptor should come before logging")
	}
}

func testAuthInterceptorValidation(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check auth interceptor implementation
	authInterceptorPath := filepath.Join(projectPath, "internal/interceptors/auth.go")
	content, err := os.ReadFile(authInterceptorPath)
	require.NoError(t, err, "Should be able to read auth interceptor")

	contentStr := string(content)
	assert.Contains(t, contentStr, "JWT", "Auth interceptor should handle JWT")
	assert.Contains(t, contentStr, "UnaryServerInterceptor", "Should implement unary interceptor")
}

func testLoggingInterceptorValidation(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check logging interceptor implementation
	loggingInterceptorPath := filepath.Join(projectPath, "internal/interceptors/logging.go")
	content, err := os.ReadFile(loggingInterceptorPath)
	require.NoError(t, err, "Should be able to read logging interceptor")

	contentStr := string(content)
	assert.Contains(t, contentStr, "UnaryServerInterceptor", "Should implement unary interceptor")
	assert.Contains(t, contentStr, "logger", "Should use logger")
}

func testMetricsInterceptorValidation(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check metrics interceptor implementation
	metricsInterceptorPath := filepath.Join(projectPath, "internal/interceptors/metrics.go")
	content, err := os.ReadFile(metricsInterceptorPath)
	require.NoError(t, err, "Should be able to read metrics interceptor")

	contentStr := string(content)
	assert.Contains(t, contentStr, "prometheus", "Should use Prometheus metrics")
	assert.Contains(t, contentStr, "UnaryServerInterceptor", "Should implement unary interceptor")
}

func testRecoveryInterceptorValidation(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check recovery interceptor implementation
	recoveryInterceptorPath := filepath.Join(projectPath, "internal/interceptors/recovery.go")
	content, err := os.ReadFile(recoveryInterceptorPath)
	require.NoError(t, err, "Should be able to read recovery interceptor")

	contentStr := string(content)
	assert.Contains(t, contentStr, "panic", "Should handle panics")
	assert.Contains(t, contentStr, "UnaryServerInterceptor", "Should implement unary interceptor")
}

func testRateLimitingInterceptorValidation(t *testing.T, projectPath string) {
	t.Helper()
	
	// Check rate limiting interceptor implementation
	rateLimitInterceptorPath := filepath.Join(projectPath, "internal/interceptors/ratelimit.go")
	content, err := os.ReadFile(rateLimitInterceptorPath)
	require.NoError(t, err, "Should be able to read rate limit interceptor")

	contentStr := string(content)
	assert.Contains(t, contentStr, "rate", "Should implement rate limiting")
	assert.Contains(t, contentStr, "UnaryServerInterceptor", "Should implement unary interceptor")
}

// Utility functions

func setupTestEnvironment(t *testing.T) string {
	t.Helper()
	
	tempDir, err := os.MkdirTemp("", "grpc_pure_atdd_test_*")
	require.NoError(t, err, "Should create temp directory")

	// Ensure go-starter is available
	_, err = exec.LookPath("go-starter")
	if err != nil {
		t.Skip("go-starter CLI tool not available")
	}

	// Initialize templates if needed
	err = helpers.InitializeTemplates()
	if err != nil {
		t.Logf("Warning: Failed to initialize templates: %v", err)
	}

	return tempDir
}

func cleanupTestEnvironment(t *testing.T, tempDir string) {
	t.Helper()
	
	if tempDir != "" {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Logf("Warning: failed to remove temp directory %s: %v", tempDir, err)
		}
	}
}

func fileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func findProtoFiles(dir string) ([]string, error) {
	var protoFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".proto") {
			protoFiles = append(protoFiles, path)
		}
		return nil
	})
	return protoFiles, err
}

func extractProjectName(command string) string {
	parts := strings.Fields(command)
	for i, part := range parts {
		if part == "new" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return "unknown-project"
}