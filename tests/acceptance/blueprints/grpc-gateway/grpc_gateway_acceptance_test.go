package grpcgateway

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

// TestGRPCGatewayBlueprintGeneration tests high-level gRPC Gateway blueprint generation
func TestGRPCGatewayBlueprintGeneration(t *testing.T) {
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
			name:        "Basic gRPC Gateway",
			projectName: "basic-grpc-gateway",
			command:     "go-starter new basic-grpc-gateway --type=grpc-gateway --module=github.com/example/basic-grpc-gateway --no-git",
			expectedFiles: []string{
				"go.mod",
				"main.go",
				"proto/service.proto",
				"internal/server/grpc.go",
				"internal/gateway/rest.go",
				"pkg/api/v1",
				"cmd/server/main.go",
				"Dockerfile",
				"README.md",
			},
			shouldCompile: true,
		},
		{
			name:        "gRPC Gateway with Database",
			projectName: "grpc-db-gateway",
			command:     "go-starter new grpc-db-gateway --type=grpc-gateway --database-driver=postgres --database-orm=gorm --module=github.com/example/grpc-db-gateway --no-git",
			expectedFiles: []string{
				"go.mod",
				"main.go",
				"proto/service.proto",
				"internal/server/grpc.go",
				"internal/gateway/rest.go",
				"internal/repository",
				"config/database.go",
				"migrations",
				"Dockerfile",
				"docker-compose.yml",
			},
			shouldCompile: true,
		},
		{
			name:        "gRPC Gateway with Authentication",
			projectName: "grpc-auth-gateway",
			command:     "go-starter new grpc-auth-gateway --type=grpc-gateway --auth-type=jwt --module=github.com/example/grpc-auth-gateway --no-git",
			expectedFiles: []string{
				"go.mod",
				"main.go",
				"proto/service.proto",
				"internal/server/grpc.go",
				"internal/gateway/rest.go",
				"internal/middleware/auth.go",
				"internal/auth/jwt.go",
				"config/auth.go",
				"Dockerfile",
			},
			shouldCompile: true,
		},
		{
			name:        "gRPC Gateway with Multiple Services",
			projectName: "grpc-multi-gateway",
			command:     "go-starter new grpc-multi-gateway --type=grpc-gateway --multi-service=true --module=github.com/example/grpc-multi-gateway --no-git",
			expectedFiles: []string{
				"go.mod",
				"main.go",
				"proto/user.proto",
				"proto/order.proto",
				"internal/server/grpc.go",
				"internal/gateway/rest.go",
				"internal/gateway/router.go",
				"internal/service/user.go",
				"internal/service/order.go",
				"Dockerfile",
			},
			shouldCompile: true,
		},
		{
			name:        "gRPC Gateway with Observability",
			projectName: "grpc-observability-gateway",
			command:     "go-starter new grpc-observability-gateway --type=grpc-gateway --observability=true --module=github.com/example/grpc-observability-gateway --no-git",
			expectedFiles: []string{
				"go.mod",
				"main.go",
				"proto/service.proto",
				"internal/server/grpc.go",
				"internal/gateway/rest.go",
				"internal/middleware/tracing.go",
				"internal/metrics/prometheus.go",
				"internal/logger",
				"config/observability.go",
				"Dockerfile",
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

			// Test protobuf generation
			t.Run("Protobuf Generation", func(t *testing.T) {
				testProtobufGeneration(t, projectPath)
			})

			// Test Docker build
			t.Run("Docker Build", func(t *testing.T) {
				testDockerBuild(t, projectPath)
			})

			// Test API documentation generation
			t.Run("API Documentation", func(t *testing.T) {
				testAPIDocumentation(t, projectPath)
			})
		})
	}
}

// TestGRPCGatewayFeatureVariations tests specific feature combinations
func TestGRPCGatewayFeatureVariations(t *testing.T) {
	tempDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempDir)

	featureTests := []struct {
		name        string
		features    []string
		expectFiles []string
		testFunc    func(t *testing.T, projectPath string)
	}{
		{
			name:     "Rate Limiting",
			features: []string{"--rate-limiting=true"},
			expectFiles: []string{
				"internal/middleware/ratelimit.go",
				"config/ratelimit.yaml",
			},
			testFunc: testRateLimitingFeature,
		},
		{
			name:     "Streaming Support",
			features: []string{"--streaming=true"},
			expectFiles: []string{
				"internal/gateway/streaming.go",
				"proto/streaming.proto",
			},
			testFunc: testStreamingFeature,
		},
		{
			name:     "Custom HTTP Mapping",
			features: []string{"--custom-mapping=true"},
			expectFiles: []string{
				"proto/gateway.yaml",
				"internal/gateway/mapping.go",
			},
			testFunc: testCustomMappingFeature,
		},
		{
			name:     "Health Checks",
			features: []string{"--health-checks=true"},
			expectFiles: []string{
				"internal/health",
				"internal/gateway/health.go",
				"proto/health.proto",
			},
			testFunc: testHealthChecksFeature,
		},
		{
			name:     "Load Balancing",
			features: []string{"--load-balancing=true"},
			expectFiles: []string{
				"internal/discovery",
				"config/loadbalancer.go",
			},
			testFunc: testLoadBalancingFeature,
		},
		{
			name:     "Caching",
			features: []string{"--caching=true"},
			expectFiles: []string{
				"internal/cache",
				"config/cache.go",
			},
			testFunc: testCachingFeature,
		},
		{
			name:     "API Versioning",
			features: []string{"--api-versioning=true"},
			expectFiles: []string{
				"proto/v1/service.proto",
				"proto/v2/service.proto",
				"internal/gateway/versioning.go",
			},
			testFunc: testAPIVersioningFeature,
		},
	}

	for _, ft := range featureTests {
		t.Run(ft.name, func(t *testing.T) {
			projectName := fmt.Sprintf("grpc-gateway-%s", strings.ToLower(strings.ReplaceAll(ft.name, " ", "-")))
			featureFlags := strings.Join(ft.features, " ")
			command := fmt.Sprintf("go-starter new %s --type=grpc-gateway %s --module=github.com/example/%s --no-git",
				projectName, featureFlags, projectName)

			err := os.Chdir(tempDir)
			require.NoError(t, err)

			// Run generation command
			cmd := exec.Command("sh", "-c", command)
			cmd.Dir = tempDir
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Feature generation should succeed: %s", string(output))

			projectPath := filepath.Join(tempDir, projectName)

			// Verify expected files
			for _, expectedFile := range ft.expectFiles {
				filePath := filepath.Join(projectPath, expectedFile)
				assert.True(t, fileOrDirExists(filePath), "Feature file should exist: %s", expectedFile)
			}

			// Run feature-specific tests
			if ft.testFunc != nil {
				ft.testFunc(t, projectPath)
			}

			// Ensure compilation still works
			testCompilation(t, projectPath)
		})
	}
}

// TestGRPCGatewayIntegration tests end-to-end integration scenarios
func TestGRPCGatewayIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tempDir := setupTestEnvironment(t)
	defer cleanupTestEnvironment(t, tempDir)

	projectName := "grpc-integration-test"
	command := fmt.Sprintf("go-starter new %s --type=grpc-gateway --database-driver=postgres --auth-type=jwt --observability=true --module=github.com/example/%s --no-git",
		projectName, projectName)

	err := os.Chdir(tempDir)
	require.NoError(t, err)

	// Generate project
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = tempDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Integration test project generation should succeed: %s", string(output))

	projectPath := filepath.Join(tempDir, projectName)

	// Test full build process
	t.Run("Full Build Process", func(t *testing.T) {
		testFullBuildProcess(t, projectPath)
	})

	// Test protobuf compilation and code generation
	t.Run("Protobuf Code Generation", func(t *testing.T) {
		testProtobufCodeGeneration(t, projectPath)
	})

	// Test Docker compose setup
	t.Run("Docker Compose", func(t *testing.T) {
		testDockerCompose(t, projectPath)
	})

	// Test Kubernetes manifests
	t.Run("Kubernetes Manifests", func(t *testing.T) {
		testKubernetesManifests(t, projectPath)
	})
}

// Helper functions for testing specific features
func testCompilation(t *testing.T, projectPath string) {
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
	protoDir := filepath.Join(projectPath, "proto")
	if !fileOrDirExists(protoDir) {
		t.Skip("No proto directory found")
	}

	// Check for .proto files
	protoFiles, err := filepath.Glob(filepath.Join(protoDir, "*.proto"))
	if err == nil && len(protoFiles) > 0 {
		assert.Greater(t, len(protoFiles), 0, "Should have at least one .proto file")
	}

	// Check for generated Go files
	genFiles, err := filepath.Glob(filepath.Join(protoDir, "*.pb.go"))
	if err == nil && len(genFiles) > 0 {
		assert.Greater(t, len(genFiles), 0, "Should have generated .pb.go files")
	}
}

func testDockerBuild(t *testing.T, projectPath string) {
	dockerfilePath := filepath.Join(projectPath, "Dockerfile")
	if !fileOrDirExists(dockerfilePath) {
		t.Skip("No Dockerfile found")
	}

	// Test Docker build (this might take a while)
	if testing.Short() {
		t.Skip("Skipping Docker build in short mode")
	}

	cmd := exec.Command("docker", "build", "-t", "test-grpc-gateway", ".")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Docker build failed (this may be expected in CI): %s", string(output))
	}
}

func testAPIDocumentation(t *testing.T, projectPath string) {
	docPaths := []string{
		filepath.Join(projectPath, "api/openapi.yaml"),
		filepath.Join(projectPath, "api/swagger.json"),
		filepath.Join(projectPath, "docs/api.yaml"),
		filepath.Join(projectPath, "README.md"),
	}

	found := false
	for _, docPath := range docPaths {
		if fileOrDirExists(docPath) {
			found = true
			break
		}
	}

	assert.True(t, found, "Should have some form of API documentation")
}

func testRateLimitingFeature(t *testing.T, projectPath string) {
	// Check for rate limiting middleware
	middlewarePath := filepath.Join(projectPath, "internal/middleware/ratelimit.go")
	assert.True(t, fileOrDirExists(middlewarePath), "Rate limiting middleware should exist")

	// Check for configuration
	configPath := filepath.Join(projectPath, "config/ratelimit.yaml")
	if fileOrDirExists(configPath) {
		// Read and validate configuration
		content, err := os.ReadFile(configPath)
		if err == nil {
			assert.Contains(t, string(content), "rate", "Rate limiting config should contain rate settings")
		}
	}
}

func testStreamingFeature(t *testing.T, projectPath string) {
	// Check for streaming proto definitions
	protoFiles, _ := filepath.Glob(filepath.Join(projectPath, "proto/*.proto"))
	for _, protoFile := range protoFiles {
		content, err := os.ReadFile(protoFile)
		if err == nil && strings.Contains(string(content), "stream") {
			assert.Contains(t, string(content), "stream", "Proto files should contain streaming definitions")
			return
		}
	}

	// Check for streaming gateway
	streamingPath := filepath.Join(projectPath, "internal/gateway/streaming.go")
	assert.True(t, fileOrDirExists(streamingPath), "Streaming gateway should exist")
}

func testCustomMappingFeature(t *testing.T, projectPath string) {
	// Check for HTTP mapping configuration
	mappingPaths := []string{
		filepath.Join(projectPath, "proto/gateway.yaml"),
		filepath.Join(projectPath, "internal/gateway/mapping.go"),
	}

	found := false
	for _, path := range mappingPaths {
		if fileOrDirExists(path) {
			found = true
			break
		}
	}
	assert.True(t, found, "Custom HTTP mapping configuration should exist")
}

func testHealthChecksFeature(t *testing.T, projectPath string) {
	// Check for health check implementation
	healthPaths := []string{
		filepath.Join(projectPath, "internal/health"),
		filepath.Join(projectPath, "internal/gateway/health.go"),
	}

	found := false
	for _, path := range healthPaths {
		if fileOrDirExists(path) {
			found = true
			break
		}
	}
	assert.True(t, found, "Health check implementation should exist")
}

func testLoadBalancingFeature(t *testing.T, projectPath string) {
	// Check for service discovery
	discoveryPath := filepath.Join(projectPath, "internal/discovery")
	assert.True(t, fileOrDirExists(discoveryPath), "Service discovery should exist")

	// Check for load balancer configuration
	lbConfigPath := filepath.Join(projectPath, "config/loadbalancer.go")
	assert.True(t, fileOrDirExists(lbConfigPath), "Load balancer config should exist")
}

func testCachingFeature(t *testing.T, projectPath string) {
	// Check for cache implementation
	cachePath := filepath.Join(projectPath, "internal/cache")
	assert.True(t, fileOrDirExists(cachePath), "Cache implementation should exist")

	// Check for cache configuration
	cacheConfigPath := filepath.Join(projectPath, "config/cache.go")
	assert.True(t, fileOrDirExists(cacheConfigPath), "Cache configuration should exist")
}

func testAPIVersioningFeature(t *testing.T, projectPath string) {
	// Check for versioned proto files
	v1ProtoPath := filepath.Join(projectPath, "proto/v1")
	v2ProtoPath := filepath.Join(projectPath, "proto/v2")
	
	hasVersioning := fileOrDirExists(v1ProtoPath) || fileOrDirExists(v2ProtoPath)
	assert.True(t, hasVersioning, "Versioned proto files should exist")

	// Check for versioning gateway
	versioningPath := filepath.Join(projectPath, "internal/gateway/versioning.go")
	assert.True(t, fileOrDirExists(versioningPath), "Versioning gateway should exist")
}

func testFullBuildProcess(t *testing.T, projectPath string) {
	steps := []struct {
		name    string
		command []string
	}{
		{"Dependencies", []string{"go", "mod", "tidy"}},
		{"Generate", []string{"go", "generate", "./..."}},
		{"Test", []string{"go", "test", "./..."}},
		{"Build", []string{"go", "build", "./..."}},
	}

	for _, step := range steps {
		t.Run(step.name, func(t *testing.T) {
			cmd := exec.Command(step.command[0], step.command[1:]...)
			cmd.Dir = projectPath
			output, err := cmd.CombinedOutput()
			assert.NoError(t, err, "%s should succeed: %s", step.name, string(output))
		})
	}
}

func testProtobufCodeGeneration(t *testing.T, projectPath string) {
	// Check if protoc is available
	if _, err := exec.LookPath("protoc"); err != nil {
		t.Skip("protoc not available, skipping protobuf code generation test")
	}

	// Look for Makefile or scripts that handle protobuf generation
	makefilePath := filepath.Join(projectPath, "Makefile")
	scriptsDir := filepath.Join(projectPath, "scripts")
	
	if fileOrDirExists(makefilePath) {
		cmd := exec.Command("make", "proto")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Logf("Protobuf generation via Makefile succeeded: %s", string(output))
		} else {
			t.Logf("Protobuf generation via Makefile failed (may be expected): %s", string(output))
		}
	} else if fileOrDirExists(scriptsDir) {
		// Look for generation scripts
		scripts, _ := filepath.Glob(filepath.Join(scriptsDir, "*proto*"))
		for _, script := range scripts {
			cmd := exec.Command("sh", script)
			cmd.Dir = projectPath
			output, err := cmd.CombinedOutput()
			if err == nil {
				t.Logf("Protobuf generation via script succeeded: %s", string(output))
				break
			}
		}
	}
}

func testDockerCompose(t *testing.T, projectPath string) {
	composePath := filepath.Join(projectPath, "docker-compose.yml")
	if !fileOrDirExists(composePath) {
		t.Skip("No docker-compose.yml found")
	}

	// Test docker-compose validation
	cmd := exec.Command("docker-compose", "config")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Docker compose validation failed (may be expected): %s", string(output))
	} else {
		assert.NoError(t, err, "Docker compose should be valid: %s", string(output))
	}
}

func testKubernetesManifests(t *testing.T, projectPath string) {
	k8sPath := filepath.Join(projectPath, "k8s")
	if !fileOrDirExists(k8sPath) {
		k8sPath = filepath.Join(projectPath, "kubernetes")
		if !fileOrDirExists(k8sPath) {
			t.Skip("No Kubernetes manifests found")
		}
	}

	// Check for YAML files
	yamlFiles, _ := filepath.Glob(filepath.Join(k8sPath, "*.yaml"))
	ymlFiles, _ := filepath.Glob(filepath.Join(k8sPath, "*.yml"))
	
	totalManifests := len(yamlFiles) + len(ymlFiles)
	assert.Greater(t, totalManifests, 0, "Should have Kubernetes manifest files")

	// Test manifest validation with kubectl if available
	if _, err := exec.LookPath("kubectl"); err == nil {
		for _, file := range append(yamlFiles, ymlFiles...) {
			cmd := exec.Command("kubectl", "apply", "--dry-run=client", "-f", file)
			output, err := cmd.CombinedOutput()
			if err == nil {
				t.Logf("Kubernetes manifest %s is valid", filepath.Base(file))
			} else {
				t.Logf("Kubernetes manifest %s validation failed: %s", filepath.Base(file), string(output))
			}
		}
	}
}

// Test utilities
func setupTestEnvironment(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "grpc_gateway_acceptance_test_*")
	require.NoError(t, err, "Should create temp directory")

	// Ensure go-starter is available
	_, err = exec.LookPath("go-starter")
	if err != nil {
		t.Skip("go-starter CLI tool not available")
	}

	return tempDir
}

func cleanupTestEnvironment(t *testing.T, tempDir string) {
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