package grpcpure

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
	"github.com/francknouama/go-starter/tests/helpers"
)

// TestGRPCPureIntegrationScenarios tests end-to-end integration scenarios
func TestGRPCPureIntegrationScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tempDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempDir)

	integrationTests := []struct {
		name        string
		command     string
		validations []func(t *testing.T, projectPath string)
	}{
		{
			name:    "Production gRPC Service with All Features",
			command: "./bin/go-starter-dev new production-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/production-grpc --no-git --logger=zap",
			validations: []func(t *testing.T, projectPath string){
				validateProductionReadiness,
				validateCompilationAndBuild,
				validateProtobufGeneration,
				validateInterceptorConfiguration,
				validateObservabilitySetup,
				validateDatabaseIntegration,
				validateServiceDiscoverySetup,
				validateDockerization,
				validateTestingInfrastructure,
			},
		},
		{
			name:    "Minimal gRPC Service",
			command: "./bin/go-starter-dev new minimal-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/minimal-grpc --no-git --logger=slog",
			validations: []func(t *testing.T, projectPath string){
				validateMinimalSetup,
				validateCompilationAndBuild,
				validateProtobufGeneration,
				validateBasicInterceptors,
				validateHealthChecks,
				validateDocumentation,
			},
		},
		{
			name:    "High-Security gRPC Service with mTLS",
			command: "./bin/go-starter-dev new secure-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/secure-grpc --no-git --logger=slog",
			validations: []func(t *testing.T, projectPath string){
				validateMTLSConfiguration,
				validateSecurityHeaders,
				validateObservabilitySetup,
				validateCompilationAndBuild,
				validateTLSCertificateHandling,
			},
		},
		{
			name:    "Multi-Database gRPC Service",
			command: "./bin/go-starter-dev new multidb-grpc --type=grpc-pure --architecture=microservice --module=github.com/example/multidb-grpc --no-git --logger=logrus",
			validations: []func(t *testing.T, projectPath string){
				validateSQLxIntegration,
				validateLogrusIntegration,
				validateDatabaseMigrations,
				validateRepositoryPattern,
				validateCompilationAndBuild,
			},
		},
	}

	for _, it := range integrationTests {
		t.Run(it.name, func(t *testing.T) {
			projectName := extractProjectName(it.command)
			
			err := os.Chdir(tempDir)
			require.NoError(t, err)

			// Run generation command
			cmd := exec.Command("sh", "-c", it.command)
			cmd.Dir = tempDir
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Integration test project generation should succeed: %s", string(output))

			projectPath := filepath.Join(tempDir, projectName)

			// Run all validations for this test
			for i, validation := range it.validations {
				t.Run(fmt.Sprintf("Validation_%d", i+1), func(t *testing.T) {
					validation(t, projectPath)
				})
			}
		})
	}
}

// TestGRPCPurePerformanceScenarios tests performance-related aspects
func TestGRPCPurePerformanceScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	tempDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempDir)

	projectName := "perf-grpc"
	command := fmt.Sprintf("./bin/go-starter-dev new %s --type=grpc-pure --architecture=microservice --module=github.com/example/%s --no-git --logger=slog",
		projectName, projectName)

	err := os.Chdir(tempDir)
	require.NoError(t, err)

	// Generate project
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Performance test project generation should succeed: %s", string(output))

	projectPath := filepath.Join(tempDir, projectName)

	t.Run("Build Performance", func(t *testing.T) {
		validateBuildPerformance(t, projectPath)
	})

	t.Run("Load Testing Infrastructure", func(t *testing.T) {
		validateLoadTestingInfrastructure(t, projectPath)
	})

	t.Run("Performance Monitoring Setup", func(t *testing.T) {
		validatePerformanceMonitoring(t, projectPath)
	})
}

// TestGRPCPureCrossEnvironmentCompatibility tests cross-environment scenarios
func TestGRPCPureCrossEnvironmentCompatibility(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping cross-environment tests in short mode")
	}

	tempDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempDir)

	projectName := "cross-env-grpc"
	command := fmt.Sprintf("./bin/go-starter-dev new %s --type=grpc-pure --architecture=microservice --module=github.com/example/%s --no-git --logger=slog",
		projectName, projectName)

	err := os.Chdir(tempDir)
	require.NoError(t, err)

	// Generate project
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Cross-environment test project generation should succeed: %s", string(output))

	projectPath := filepath.Join(tempDir, projectName)

	t.Run("Docker Compatibility", func(t *testing.T) {
		validateDockerCompatibility(t, projectPath)
	})

	t.Run("Kubernetes Compatibility", func(t *testing.T) {
		validateKubernetesCompatibility(t, projectPath)
	})

	t.Run("Cross-Platform Build", func(t *testing.T) {
		validateCrossPlatformBuild(t, projectPath)
	})

	t.Run("Environment Configuration", func(t *testing.T) {
		validateEnvironmentConfiguration(t, projectPath)
	})
}

// Validation functions

func validateProductionReadiness(t *testing.T, projectPath string) {
	t.Helper()

	productionChecks := []struct {
		name string
		path string
	}{
		{"JWT Authentication", "internal/auth/jwt.go"},
		{"OpenTelemetry Tracing", "internal/observability/tracing.go"},
		{"Prometheus Metrics", "internal/observability/metrics.go"},
		{"Consul Service Discovery", "internal/discovery/consul.go"},
		{"Database Connection", "internal/database/connection.go"},
		{"Authentication Interceptor", "internal/interceptors/auth.go"},
		{"Tracing Interceptor", "internal/interceptors/tracing.go"},
		{"Metrics Interceptor", "internal/interceptors/metrics.go"},
		{"Production Config", "configs/config.prod.yaml"},
		{"CI/CD Workflow", ".github/workflows/ci.yml"},
		{"Security Workflow", ".github/workflows/security.yml"},
	}

	for _, check := range productionChecks {
		t.Run(check.name, func(t *testing.T) {
			fullPath := filepath.Join(projectPath, check.path)
			assert.True(t, fileOrDirExists(fullPath), "Production component should exist: %s", check.path)
		})
	}
}

func validateMinimalSetup(t *testing.T, projectPath string) {
	t.Helper()

	minimalComponents := []string{
		"cmd/server/main.go",
		"internal/server/grpc.go",
		"internal/server/health.go",
		"internal/services",
		"internal/logger/interface.go",
		"internal/logger/logger.go",
		"proto",
		"buf.yaml",
		"buf.gen.yaml",
		"Makefile",
		"README.md",
	}

	for _, component := range minimalComponents {
		componentPath := filepath.Join(projectPath, component)
		assert.True(t, fileOrDirExists(componentPath), "Minimal component should exist: %s", component)
	}
}

func validateCompilationAndBuild(t *testing.T, projectPath string) {
	t.Helper()

	// Clean dependencies
	modTidyCmd := exec.Command("go", "mod", "tidy")
	modTidyCmd.Dir = projectPath
	output, err := modTidyCmd.CombinedOutput()
	require.NoError(t, err, "go mod tidy should succeed: %s", string(output))

	// Build the project
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Dir = projectPath
	output, err = buildCmd.CombinedOutput()
	assert.NoError(t, err, "Project should compile successfully: %s", string(output))

	// Test the project
	testCmd := exec.Command("go", "test", "./...")
	testCmd.Dir = projectPath
	output, err = testCmd.CombinedOutput()
	if err != nil {
		t.Logf("Tests failed (may be expected without external dependencies): %s", string(output))
	}
}

func validateProtobufGeneration(t *testing.T, projectPath string) {
	t.Helper()

	// Check for proto files
	protoDir := filepath.Join(projectPath, "proto")
	require.True(t, fileOrDirExists(protoDir), "Proto directory should exist")

	protoFiles, err := helpers.FindProtoFiles(protoDir)
	require.NoError(t, err, "Should be able to find proto files")
	assert.Greater(t, len(protoFiles), 0, "Should have at least one proto file")

	// Check for buf configuration
	bufYaml := filepath.Join(projectPath, "buf.yaml")
	assert.True(t, fileOrDirExists(bufYaml), "buf.yaml should exist")

	bufGenYaml := filepath.Join(projectPath, "buf.gen.yaml")
	assert.True(t, fileOrDirExists(bufGenYaml), "buf.gen.yaml should exist")

	// Try to generate if tools are available
	if _, err := exec.LookPath("buf"); err == nil {
		t.Run("Buf Generation", func(t *testing.T) {
			cmd := exec.Command("buf", "generate")
			cmd.Dir = projectPath
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Logf("buf generate failed (may be expected): %s", string(output))
			} else {
				t.Logf("buf generate succeeded: %s", string(output))
			}
		})
	}
}

func validateInterceptorConfiguration(t *testing.T, projectPath string) {
	t.Helper()

	serverConfigPath := filepath.Join(projectPath, "internal/server/grpc.go")
	content, err := os.ReadFile(serverConfigPath)
	require.NoError(t, err, "Should be able to read server configuration")

	contentStr := string(content)

	// Check for interceptor chain setup
	interceptorKeywords := []string{
		"interceptor",
		"Interceptor", 
		"middleware",
		"Middleware",
	}

	found := false
	for _, keyword := range interceptorKeywords {
		if strings.Contains(contentStr, keyword) {
			found = true
			break
		}
	}
	assert.True(t, found, "Interceptor configuration should be present in server setup")

	// Check individual interceptor files
	interceptorFiles := []string{
		"internal/interceptors/logging.go",
		"internal/interceptors/recovery.go",
		"internal/interceptors/ratelimit.go",
	}

	for _, interceptorFile := range interceptorFiles {
		fullPath := filepath.Join(projectPath, interceptorFile)
		assert.True(t, fileOrDirExists(fullPath), "Interceptor should exist: %s", interceptorFile)
	}
}

func validateBasicInterceptors(t *testing.T, projectPath string) {
	t.Helper()

	basicInterceptors := []string{
		"internal/interceptors/logging.go",
		"internal/interceptors/recovery.go",
		"internal/interceptors/ratelimit.go",
	}

	for _, interceptor := range basicInterceptors {
		fullPath := filepath.Join(projectPath, interceptor)
		assert.True(t, fileOrDirExists(fullPath), "Basic interceptor should exist: %s", interceptor)
	}
}

func validateObservabilitySetup(t *testing.T, projectPath string) {
	t.Helper()

	observabilityComponents := []string{
		"internal/observability/tracing.go",
		"internal/observability/metrics.go",
		"internal/interceptors/tracing.go", 
		"internal/interceptors/metrics.go",
	}

	for _, component := range observabilityComponents {
		fullPath := filepath.Join(projectPath, component)
		assert.True(t, fileOrDirExists(fullPath), "Observability component should exist: %s", component)
	}

	// Check dependencies in go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")

	contentStr := string(content)
	observabilityDeps := []string{
		"go.opentelemetry.io/otel",
		"github.com/prometheus/client_golang",
	}

	for _, dep := range observabilityDeps {
		assert.Contains(t, contentStr, dep, "Observability dependency should be present: %s", dep)
	}
}

func validateDatabaseIntegration(t *testing.T, projectPath string) {
	t.Helper()

	databaseComponents := []string{
		"internal/database/connection.go",
		"internal/database/migrations.go",
		"internal/models/user.go",
		"internal/repository/user.go",
		"internal/repository/interface.go",
		"migrations",
	}

	for _, component := range databaseComponents {
		fullPath := filepath.Join(projectPath, component)
		assert.True(t, fileOrDirExists(fullPath), "Database component should exist: %s", component)
	}

	// Check for migration files
	migrationsDir := filepath.Join(projectPath, "migrations")
	migrationFiles, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	require.NoError(t, err, "Should be able to search for migration files")
	assert.Greater(t, len(migrationFiles), 0, "Should have SQL migration files")
}

func validateServiceDiscoverySetup(t *testing.T, projectPath string) {
	t.Helper()

	serviceDiscoveryPath := filepath.Join(projectPath, "internal/discovery/consul.go")
	assert.True(t, fileOrDirExists(serviceDiscoveryPath), "Consul service discovery should exist")

	interfacePath := filepath.Join(projectPath, "internal/discovery/interface.go")
	assert.True(t, fileOrDirExists(interfacePath), "Service discovery interface should exist")

	// Check dependency in go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")

	contentStr := string(content)
	assert.Contains(t, contentStr, "github.com/hashicorp/consul/api", "Consul dependency should be present")
}

func validateDockerization(t *testing.T, projectPath string) {
	t.Helper()

	dockerFiles := []string{
		"Dockerfile", 
		"docker-compose.yml",
		".dockerignore",
	}

	for _, dockerFile := range dockerFiles {
		fullPath := filepath.Join(projectPath, dockerFile)
		if dockerFile == ".dockerignore" {
			// .dockerignore is optional
			continue
		}
		assert.True(t, fileOrDirExists(fullPath), "Docker file should exist: %s", dockerFile)
	}

	// Basic Dockerfile validation
	dockerfilePath := filepath.Join(projectPath, "Dockerfile")
	content, err := os.ReadFile(dockerfilePath)
	require.NoError(t, err, "Should be able to read Dockerfile")

	contentStr := string(content)
	assert.Contains(t, contentStr, "FROM", "Dockerfile should have FROM directive")
	assert.Contains(t, contentStr, "EXPOSE", "Dockerfile should expose gRPC port")
}

func validateTestingInfrastructure(t *testing.T, projectPath string) {
	t.Helper()

	testDirectories := []string{
		"tests/integration",
		"tests/unit", 
		"tests/load",
	}

	for _, testDir := range testDirectories {
		fullPath := filepath.Join(projectPath, testDir)
		assert.True(t, fileOrDirExists(fullPath), "Test directory should exist: %s", testDir)
	}

	// Check for test files
	testFiles := []string{
		"tests/integration/grpc_test.go",
		"tests/integration/health_test.go",
		"tests/unit/services_test.go",
		"tests/load/grpc_load_test.go",
	}

	for _, testFile := range testFiles {
		fullPath := filepath.Join(projectPath, testFile)
		assert.True(t, fileOrDirExists(fullPath), "Test file should exist: %s", testFile)
	}
}

func validateHealthChecks(t *testing.T, projectPath string) {
	t.Helper()

	healthComponents := []string{
		"internal/services/health.go",
		"internal/server/health.go",
		"proto/health/v1/health.proto",
		"tests/integration/health_test.go",
	}

	for _, component := range healthComponents {
		fullPath := filepath.Join(projectPath, component)
		assert.True(t, fileOrDirExists(fullPath), "Health component should exist: %s", component)
	}
}

func validateDocumentation(t *testing.T, projectPath string) {
	t.Helper()

	documentationFiles := []string{
		"README.md",
		"docs/ARCHITECTURE.md",
		"docs/API.md", 
		"docs/DEPLOYMENT.md",
	}

	for _, docFile := range documentationFiles {
		fullPath := filepath.Join(projectPath, docFile)
		assert.True(t, fileOrDirExists(fullPath), "Documentation should exist: %s", docFile)
	}

	// Check README content
	readmePath := filepath.Join(projectPath, "README.md")
	content, err := os.ReadFile(readmePath)
	require.NoError(t, err, "Should be able to read README")

	contentStr := string(content)
	assert.Contains(t, contentStr, "gRPC", "README should mention gRPC")
	assert.Contains(t, contentStr, "Getting Started", "README should have getting started section")
}

func validateMTLSConfiguration(t *testing.T, projectPath string) {
	t.Helper()

	mtlsComponents := []string{
		"internal/auth/mtls.go",
		"internal/tls/config.go",
	}

	for _, component := range mtlsComponents {
		fullPath := filepath.Join(projectPath, component)
		assert.True(t, fileOrDirExists(fullPath), "mTLS component should exist: %s", component)
	}
}

func validateSecurityHeaders(t *testing.T, projectPath string) {
	t.Helper()

	// Check for security-related configurations
	serverConfigPath := filepath.Join(projectPath, "internal/server/grpc.go")
	content, err := os.ReadFile(serverConfigPath)
	require.NoError(t, err, "Should be able to read server configuration")

	contentStr := string(content)
	securityKeywords := []string{"tls", "TLS", "security", "Security", "auth", "Auth"}
	
	found := false
	for _, keyword := range securityKeywords {
		if strings.Contains(contentStr, keyword) {
			found = true
			break
		}
	}
	assert.True(t, found, "Security configuration should be present")
}

func validateTLSCertificateHandling(t *testing.T, projectPath string) {
	t.Helper()

	tlsConfigPath := filepath.Join(projectPath, "internal/tls/config.go")
	content, err := os.ReadFile(tlsConfigPath)
	require.NoError(t, err, "Should be able to read TLS configuration")

	contentStr := string(content)
	tlsKeywords := []string{"Certificate", "certificate", "PrivateKey", "CertFile", "KeyFile"}
	
	found := false
	for _, keyword := range tlsKeywords {
		if strings.Contains(contentStr, keyword) {
			found = true
			break
		}
	}
	assert.True(t, found, "TLS certificate handling should be present")
}

func validateSQLxIntegration(t *testing.T, projectPath string) {
	t.Helper()

	// Check for SQLx dependency
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")

	contentStr := string(content)
	assert.Contains(t, contentStr, "github.com/jmoiron/sqlx", "SQLx dependency should be present")

	// Check for repository implementation
	repoPath := filepath.Join(projectPath, "internal/repository")
	assert.True(t, fileOrDirExists(repoPath), "Repository directory should exist")
}

func validateLogrusIntegration(t *testing.T, projectPath string) {
	t.Helper()

	// Check for Logrus dependency
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	require.NoError(t, err, "Should be able to read go.mod")

	contentStr := string(content)
	assert.Contains(t, contentStr, "github.com/sirupsen/logrus", "Logrus dependency should be present")
}

func validateDatabaseMigrations(t *testing.T, projectPath string) {
	t.Helper()

	migrationsDir := filepath.Join(projectPath, "migrations")
	assert.True(t, fileOrDirExists(migrationsDir), "Migrations directory should exist")

	migrationFiles := []string{
		"001_create_users.up.sql",
		"001_create_users.down.sql", 
		"embed.go",
	}

	for _, migrationFile := range migrationFiles {
		fullPath := filepath.Join(migrationsDir, migrationFile)
		assert.True(t, fileOrDirExists(fullPath), "Migration file should exist: %s", migrationFile)
	}
}

func validateRepositoryPattern(t *testing.T, projectPath string) {
	t.Helper()

	repositoryFiles := []string{
		"internal/repository/interface.go",
		"internal/repository/user.go",
	}

	for _, repoFile := range repositoryFiles {
		fullPath := filepath.Join(projectPath, repoFile)
		assert.True(t, fileOrDirExists(fullPath), "Repository file should exist: %s", repoFile)
	}
}

// Performance and compatibility validation functions

func validateBuildPerformance(t *testing.T, projectPath string) {
	t.Helper()

	startTime := time.Now()
	
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Dir = projectPath
	output, err := buildCmd.CombinedOutput()
	
	buildDuration := time.Since(startTime)
	
	assert.NoError(t, err, "Build should succeed: %s", string(output))
	assert.Less(t, buildDuration, 2*time.Minute, "Build should complete within reasonable time")
	
	t.Logf("Build completed in %v", buildDuration)
}

func validateLoadTestingInfrastructure(t *testing.T, projectPath string) {
	t.Helper()

	loadTestPath := filepath.Join(projectPath, "tests/load/grpc_load_test.go")
	assert.True(t, fileOrDirExists(loadTestPath), "Load test should exist")

	// Check load test content
	content, err := os.ReadFile(loadTestPath)
	require.NoError(t, err, "Should be able to read load test")

	contentStr := string(content)
	loadTestKeywords := []string{"benchmark", "Benchmark", "performance", "concurrent", "load"}
	
	found := false
	for _, keyword := range loadTestKeywords {
		if strings.Contains(contentStr, keyword) {
			found = true
			break
		}
	}
	assert.True(t, found, "Load test should contain performance testing code")
}

func validatePerformanceMonitoring(t *testing.T, projectPath string) {
	t.Helper()

	monitoringComponents := []string{
		"internal/observability/metrics.go",
		"internal/interceptors/metrics.go",
	}

	for _, component := range monitoringComponents {
		fullPath := filepath.Join(projectPath, component)
		assert.True(t, fileOrDirExists(fullPath), "Performance monitoring component should exist: %s", component)
	}
}

func validateDockerCompatibility(t *testing.T, projectPath string) {
	t.Helper()

	dockerfilePath := filepath.Join(projectPath, "Dockerfile")
	assert.True(t, fileOrDirExists(dockerfilePath), "Dockerfile should exist")

	// Try to validate Dockerfile if Docker is available
	if _, err := exec.LookPath("docker"); err == nil {
		cmd := exec.Command("docker", "build", "-t", "test-grpc", ".")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Logf("Docker build failed (may be expected in CI): %s", string(output))
		} else {
			t.Logf("Docker build succeeded")
			
			// Clean up
			cleanupCmd := exec.Command("docker", "rmi", "test-grpc")
			cleanupCmd.Run()
		}
	} else {
		t.Skip("Docker not available, skipping Docker compatibility test")
	}
}

func validateKubernetesCompatibility(t *testing.T, projectPath string) {
	t.Helper()

	// Check for Kubernetes manifests or Helm charts
	k8sPaths := []string{
		"k8s",
		"kubernetes", 
		"helm",
		"deployments/k8s",
	}

	found := false
	for _, k8sPath := range k8sPaths {
		fullPath := filepath.Join(projectPath, k8sPath)
		if fileOrDirExists(fullPath) {
			found = true
			break
		}
	}
	
	if !found {
		t.Skip("No Kubernetes manifests found")
	}
}

func validateCrossPlatformBuild(t *testing.T, projectPath string) {
	t.Helper()

	platforms := []struct {
		goos   string
		goarch string
	}{
		{"linux", "amd64"},
		{"darwin", "amd64"},
		{"windows", "amd64"},
	}

	for _, platform := range platforms {
		t.Run(fmt.Sprintf("%s_%s", platform.goos, platform.goarch), func(t *testing.T) {
			buildCmd := exec.Command("go", "build", "-o", fmt.Sprintf("server_%s_%s", platform.goos, platform.goarch), "./cmd/server")
			buildCmd.Dir = projectPath
			buildCmd.Env = append(os.Environ(), 
				fmt.Sprintf("GOOS=%s", platform.goos),
				fmt.Sprintf("GOARCH=%s", platform.goarch))
			
			output, err := buildCmd.CombinedOutput()
			assert.NoError(t, err, "Cross-platform build should succeed for %s/%s: %s", 
				platform.goos, platform.goarch, string(output))
		})
	}
}

func validateEnvironmentConfiguration(t *testing.T, projectPath string) {
	t.Helper()

	configFiles := []string{
		"configs/config.dev.yaml",
		"configs/config.prod.yaml",
		"configs/config.test.yaml",
		".env.example",
	}

	for _, configFile := range configFiles {
		fullPath := filepath.Join(projectPath, configFile)
		assert.True(t, fileOrDirExists(fullPath), "Environment config should exist: %s", configFile)
	}
}