package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ValidateProtocolBuffers uses buf to validate proto files
func ValidateProtocolBuffers(projectPath string) error {
	protoDir := filepath.Join(projectPath, "proto")
	if !fileOrDirExists(protoDir) {
		return fmt.Errorf("proto directory not found: %s", protoDir)
	}

	// Check if buf is available
	if _, err := exec.LookPath("buf"); err != nil {
		return fmt.Errorf("buf tool not available: %w", err)
	}

	// Run buf lint
	cmd := exec.Command("buf", "lint", protoDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("proto validation failed: %s", string(output))
	}
	return nil
}

// GenerateProtocolBuffers runs proto generation using make or buf
func GenerateProtocolBuffers(projectPath string) error {
	// Try using Makefile first
	makefilePath := filepath.Join(projectPath, "Makefile")
	if fileOrDirExists(makefilePath) {
		cmd := exec.Command("make", "generate")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("proto generation via make failed: %s", string(output))
		}
		return nil
	}

	// Fall back to buf generate if available
	if _, err := exec.LookPath("buf"); err == nil {
		cmd := exec.Command("buf", "generate")
		cmd.Dir = projectPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("proto generation via buf failed: %s", string(output))
		}
		return nil
	}

	return fmt.Errorf("no proto generation tool available (make or buf)")
}

// TestGRPCServiceHealth tests gRPC health check endpoint
func TestGRPCServiceHealth(projectPath string, port int) error {
	// Check if grpcurl is available
	if _, err := exec.LookPath("grpcurl"); err != nil {
		return fmt.Errorf("grpcurl not available for health check testing: %w", err)
	}

	// Test health check endpoint
	cmd := exec.Command("grpcurl", "-plaintext", 
		fmt.Sprintf("localhost:%d", port), 
		"grpc.health.v1.Health/Check")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("health check failed: %s", string(output))
	}
	
	return nil
}

// StartGRPCService starts a gRPC service for testing and returns a cleanup function
func StartGRPCService(projectPath string, port int) (func(), error) {
	// Build the service first
	buildCmd := exec.Command("go", "build", "-o", "test-server", "./cmd/server")
	buildCmd.Dir = projectPath
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to build gRPC service: %s", string(output))
	}

	// Start the service
	serverCmd := exec.Command("./test-server")
	serverCmd.Dir = projectPath
	serverCmd.Env = append(os.Environ(), fmt.Sprintf("GRPC_PORT=%d", port))
	
	err = serverCmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start gRPC service: %w", err)
	}

	// Wait a bit for the service to start
	time.Sleep(2 * time.Second)

	// Return cleanup function
	cleanup := func() {
		if serverCmd.Process != nil {
			_ = serverCmd.Process.Kill()
		}
		// Clean up binary
		_ = os.Remove(filepath.Join(projectPath, "test-server"))
	}

	return cleanup, nil
}

// LoadTestGRPCService performs load testing on a gRPC service
func LoadTestGRPCService(projectPath string, port int, requests int, concurrency int) error {
	// Check if ghz is available
	if _, err := exec.LookPath("ghz"); err != nil {
		return fmt.Errorf("ghz not available for load testing: %w", err)
	}

	// Find proto files
	protoFiles, err := findProtoFiles(filepath.Join(projectPath, "proto"))
	if err != nil || len(protoFiles) == 0 {
		return fmt.Errorf("no proto files found for load testing")
	}

	// Use the first proto file for load testing
	protoFile := protoFiles[0]
	
	// Extract service and method names from proto file
	serviceName, methodName, err := extractServiceMethod(protoFile)
	if err != nil {
		return fmt.Errorf("failed to extract service/method from proto: %w", err)
	}

	// Run load test
	cmd := exec.Command("ghz",
		"--insecure",
		"--proto", protoFile,
		"--call", fmt.Sprintf("%s.%s", serviceName, methodName),
		"--total", fmt.Sprintf("%d", requests),
		"--concurrency", fmt.Sprintf("%d", concurrency),
		fmt.Sprintf("localhost:%d", port))
	
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("load test failed: %s", string(output))
	}

	return nil
}

// ValidateGRPCInterceptors checks if interceptors are properly configured
func ValidateGRPCInterceptors(projectPath string, expectedInterceptors []string) error {
	serverConfigPath := filepath.Join(projectPath, "internal/server/grpc.go")
	content, err := os.ReadFile(serverConfigPath)
	if err != nil {
		return fmt.Errorf("failed to read server config: %w", err)
	}

	contentStr := string(content)
	for _, interceptor := range expectedInterceptors {
		if !strings.Contains(contentStr, interceptor) {
			return fmt.Errorf("interceptor %s not found in server configuration", interceptor)
		}
	}

	return nil
}

// CheckGRPCDependencies validates that required gRPC dependencies are present
func CheckGRPCDependencies(projectPath string, expectedDeps []string) error {
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	contentStr := string(content)
	for _, dep := range expectedDeps {
		if !strings.Contains(contentStr, dep) {
			return fmt.Errorf("dependency %s not found in go.mod", dep)
		}
	}

	return nil
}

// ValidateGRPCServiceImplementation checks if service interfaces are properly implemented
func ValidateGRPCServiceImplementation(projectPath string, serviceName string) error {
	servicePath := filepath.Join(projectPath, "internal/services", strings.ToLower(serviceName)+".go")
	if !fileOrDirExists(servicePath) {
		return fmt.Errorf("service implementation not found: %s", servicePath)
	}

	content, err := os.ReadFile(servicePath)
	if err != nil {
		return fmt.Errorf("failed to read service implementation: %w", err)
	}

	contentStr := string(content)
	
	// Check for basic service structure
	if !strings.Contains(contentStr, "type") {
		return fmt.Errorf("service type definition not found")
	}

	if !strings.Contains(contentStr, "func") {
		return fmt.Errorf("service methods not found")
	}

	return nil
}

// ValidateGRPCHealthCheck ensures health check service is properly implemented
func ValidateGRPCHealthCheck(projectPath string) error {
	healthServicePath := filepath.Join(projectPath, "internal/services/health.go")
	if !fileOrDirExists(healthServicePath) {
		return fmt.Errorf("health service implementation not found")
	}

	healthServerPath := filepath.Join(projectPath, "internal/server/health.go")
	if !fileOrDirExists(healthServerPath) {
		return fmt.Errorf("health server implementation not found")
	}

	// Check for health proto
	healthProtoPath := filepath.Join(projectPath, "proto/health/v1/health.proto")
	if !fileOrDirExists(healthProtoPath) {
		return fmt.Errorf("health proto definition not found")
	}

	return nil
}

// CompileAndTestGRPCProject compiles project and runs basic tests
func CompileAndTestGRPCProject(projectPath string) error {
	// Clean dependencies
	modTidyCmd := exec.Command("go", "mod", "tidy")
	modTidyCmd.Dir = projectPath
	output, err := modTidyCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %s", string(output))
	}

	// Generate protobuf code if needed
	if err := GenerateProtocolBuffers(projectPath); err != nil {
		// Not a fatal error, continue with compilation
		fmt.Printf("Warning: Proto generation failed: %v\n", err)
	}

	// Compile project
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Dir = projectPath
	output, err = buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output))
	}

	// Run tests
	testCmd := exec.Command("go", "test", "./...")
	testCmd.Dir = projectPath
	output, err = testCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("tests failed: %s", string(output))
	}

	return nil
}

// FindProtoFiles finds all .proto files in a directory and its subdirectories
func FindProtoFiles(dir string) ([]string, error) {
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

// Helper functions

func fileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func findProtoFiles(dir string) ([]string, error) {
	return FindProtoFiles(dir)
}

func extractServiceMethod(protoFile string) (string, string, error) {
	content, err := os.ReadFile(protoFile)
	if err != nil {
		return "", "", err
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	var serviceName, methodName string
	inService := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if strings.HasPrefix(line, "service ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				serviceName = parts[1]
				inService = true
			}
		} else if inService && strings.HasPrefix(line, "rpc ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				methodName = parts[1]
				break
			}
		} else if inService && line == "}" {
			break
		}
	}

	if serviceName == "" || methodName == "" {
		return "", "", fmt.Errorf("could not extract service and method from proto file")
	}

	return serviceName, methodName, nil
}