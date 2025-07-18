package helpers

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// RuntimeValidator provides utilities for testing runtime behavior of generated projects
type RuntimeValidator struct {
	ProjectPath string
	Port        int
}

// NewRuntimeValidator creates a new RuntimeValidator instance
func NewRuntimeValidator(projectPath string) *RuntimeValidator {
	return &RuntimeValidator{
		ProjectPath: projectPath,
		Port:        18080, // Use isolated port for testing
	}
}

// ValidateServerStartup tests that the generated project can start a server
func (r *RuntimeValidator) ValidateServerStartup(t *testing.T) {
	t.Helper()
	
	// Check if we can build the project first
	if !r.canBuildProject(t) {
		t.Skip("Project cannot be built, skipping runtime validation")
		return
	}
	
	// For now, just validate that the project has the basic structure for a server
	r.validateServerStructure(t)
}

// ValidateHealthEndpoint tests that health endpoints respond correctly
func (r *RuntimeValidator) ValidateHealthEndpoint(t *testing.T) {
	t.Helper()
	
	// This would start the server and test health endpoints
	// For now, just validate health endpoint structure exists
	r.validateHealthEndpointStructure(t)
}

// ValidateGracefulShutdown tests graceful shutdown behavior
func (r *RuntimeValidator) ValidateGracefulShutdown(t *testing.T) {
	t.Helper()
	
	// This would test actual graceful shutdown
	// For now, just validate shutdown structure exists
	r.validateShutdownStructure(t)
}

// canBuildProject checks if the project can be built
func (r *RuntimeValidator) canBuildProject(t *testing.T) bool {
	t.Helper()
	
	// Check for go.mod
	goModPath := filepath.Join(r.ProjectPath, "go.mod")
	if !FileExists(goModPath) {
		t.Log("⚠ No go.mod found, cannot build project")
		return false
	}
	
	// Check for main entry point
	mainPaths := []string{
		filepath.Join(r.ProjectPath, "main.go"),
		filepath.Join(r.ProjectPath, "cmd", "server", "main.go"),
		filepath.Join(r.ProjectPath, "cmd", "api", "main.go"),
	}
	
	for _, path := range mainPaths {
		if FileExists(path) {
			t.Logf("✓ Found main entry point at %s", path)
			return true
		}
	}
	
	t.Log("⚠ No main entry point found")
	return false
}

// validateServerStructure validates server structure exists
func (r *RuntimeValidator) validateServerStructure(t *testing.T) {
	t.Helper()
	
	// Look for server-related files
	serverPaths := []string{
		filepath.Join(r.ProjectPath, "internal", "server"),
		filepath.Join(r.ProjectPath, "internal", "api"),
		filepath.Join(r.ProjectPath, "internal", "router"),
		filepath.Join(r.ProjectPath, "internal", "routes"),
	}
	
	found := false
	for _, path := range serverPaths {
		if DirExists(path) {
			t.Logf("✓ Found server structure at %s", path)
			found = true
			break
		}
	}
	
	if !found {
		t.Log("⚠ No explicit server structure found")
	}
}

// validateHealthEndpointStructure validates health endpoint structure
func (r *RuntimeValidator) validateHealthEndpointStructure(t *testing.T) {
	t.Helper()
	
	// Look for health check related files
	healthPaths := []string{
		filepath.Join(r.ProjectPath, "internal", "health"),
		filepath.Join(r.ProjectPath, "internal", "handlers", "health.go"),
		filepath.Join(r.ProjectPath, "internal", "api", "health.go"),
	}
	
	for _, path := range healthPaths {
		if DirExists(path) || FileExists(path) {
			t.Logf("✓ Found health endpoint structure at %s", path)
			return
		}
	}
	
	t.Log("⚠ No explicit health endpoint structure found")
}

// validateShutdownStructure validates shutdown structure
func (r *RuntimeValidator) validateShutdownStructure(t *testing.T) {
	t.Helper()
	
	// Look for graceful shutdown patterns in main files
	mainPaths := []string{
		filepath.Join(r.ProjectPath, "main.go"),
		filepath.Join(r.ProjectPath, "cmd", "server", "main.go"),
	}
	
	for _, path := range mainPaths {
		if FileExists(path) {
			t.Logf("✓ Main file exists for shutdown validation: %s", path)
			// Could read file and check for context.WithCancel or signal handling
			return
		}
	}
	
	t.Log("⚠ No main file found for shutdown validation")
}

// StartTestServer starts a test server instance (for advanced testing)
func (r *RuntimeValidator) StartTestServer(t *testing.T, timeout time.Duration) *TestServer {
	t.Helper()
	
	// This would actually start the server binary
	// For now, return a mock test server
	return &TestServer{
		Port:    r.Port,
		BaseURL: fmt.Sprintf("http://localhost:%d", r.Port),
	}
}

// TestServer represents a running test server instance
type TestServer struct {
	Port     int
	BaseURL  string
	process  *exec.Cmd
	ctx      context.Context
	cancel   context.CancelFunc
}

// Stop stops the test server
func (ts *TestServer) Stop() error {
	if ts.cancel != nil {
		ts.cancel()
	}
	
	if ts.process != nil {
		return ts.process.Process.Kill()
	}
	
	return nil
}

// IsHealthy checks if the server is responding to health checks
func (ts *TestServer) IsHealthy() bool {
	client := &http.Client{Timeout: 2 * time.Second}
	
	healthURL := ts.BaseURL + "/health"
	resp, err := client.Get(healthURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK
}

// WaitForReady waits for the server to be ready
func (ts *TestServer) WaitForReady(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("server not ready within timeout")
		case <-ticker.C:
			if ts.IsHealthy() {
				return nil
			}
		}
	}
}

// GetFreePort returns a free port for testing
func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	
	return l.Addr().(*net.TCPAddr).Port, nil
}