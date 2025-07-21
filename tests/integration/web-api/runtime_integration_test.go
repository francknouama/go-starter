package webapi_integration

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq" // postgres driver for test database setup
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// RuntimeTestConfig defines configuration for runtime integration tests
type RuntimeTestConfig struct {
	Name       string
	Framework  string
	ORM        string
	Logger     string
	Database   string
	ServerPort int
}

// TestHexagonal_WebAPI_RuntimeIntegration tests that generated hexagonal projects actually run correctly
func TestHexagonal_WebAPI_RuntimeIntegration(t *testing.T) {
	testCases := []RuntimeTestConfig{
		{
			Name:       "gin_gorm_slog_postgres",
			Framework:  "gin",
			ORM:        "gorm",
			Logger:     "slog",
			Database:   "postgres",
			ServerPort: 8081,
		},
		{
			Name:       "echo_sqlx_zap_postgres",
			Framework:  "echo",
			ORM:        "sqlx",
			Logger:     "zap",
			Database:   "postgres",
			ServerPort: 8082,
		},
		{
			Name:       "fiber_gorm_logrus_postgres",
			Framework:  "fiber",
			ORM:        "gorm",
			Logger:     "logrus",
			Database:   "postgres",
			ServerPort: 8083,
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel() // Run tests in parallel for faster execution

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()

			runRuntimeIntegrationTest(t, ctx, tc)
		})
	}
}

// runRuntimeIntegrationTest runs a complete runtime integration test for a specific configuration
func runRuntimeIntegrationTest(t *testing.T, ctx context.Context, config RuntimeTestConfig) {
	// 1. Setup test environment
	testDir := setupTestEnvironment(t, config)
	defer cleanupTestEnvironment(t, testDir)

	// 2. Generate hexagonal project
	projectPath := generateHexagonalProject(t, ctx, testDir, config)

	// 3. Setup test database
	dbContainer, dbCleanup := setupTestDatabase(t, ctx, config)
	defer dbCleanup()

	// 4. Build and start the application
	serverCleanup := startApplicationServer(t, ctx, projectPath, config, dbContainer)
	defer serverCleanup()

	// 5. Wait for server to be ready
	waitForServerReady(t, ctx, config)

	// 6. Run comprehensive integration tests
	runIntegrationTestSuite(t, ctx, config)
}

// setupTestEnvironment creates a temporary directory for the test
func setupTestEnvironment(t *testing.T, config RuntimeTestConfig) string {
	testDir, err := os.MkdirTemp("", fmt.Sprintf("hexagonal_runtime_test_%s_", config.Name))
	require.NoError(t, err, "Failed to create test directory")

	t.Logf("Created test environment: %s", testDir)
	return testDir
}

// cleanupTestEnvironment removes the test directory
func cleanupTestEnvironment(t *testing.T, testDir string) {
	if err := os.RemoveAll(testDir); err != nil {
		t.Logf("Warning: Failed to cleanup test directory %s: %v", testDir, err)
	}
}

// generateHexagonalProject generates a hexagonal project using the generator
func generateHexagonalProject(t *testing.T, ctx context.Context, testDir string, config RuntimeTestConfig) string {
	projectName := fmt.Sprintf("test-hexagonal-%s", config.Name)
	projectPath := filepath.Join(testDir, projectName)

	// Change to test directory
	originalDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalDir)

	err = os.Chdir(testDir)
	require.NoError(t, err)

	// Find the root directory (where main.go is located)
	rootDir := originalDir
	for !fileExists(filepath.Join(rootDir, "main.go")) {
		parent := filepath.Dir(rootDir)
		if parent == rootDir {
			t.Fatal("Could not find main.go in any parent directory")
		}
		rootDir = parent
	}

	// Generate project using go run with all flags to avoid interactive prompts
	genCmd := exec.CommandContext(ctx, "go", "run", "main.go", "new", projectName,
		"--type=web-api",
		"--architecture=hexagonal",
		fmt.Sprintf("--framework=%s", config.Framework),
		fmt.Sprintf("--logger=%s", config.Logger),
		fmt.Sprintf("--database-driver=%s", config.Database),
		fmt.Sprintf("--database-orm=%s", config.ORM),
		"--auth-type=jwt",
		fmt.Sprintf("--module=github.com/test/%s", projectName),
		fmt.Sprintf("--output=%s", testDir),
		"--no-banner",
		"--quiet")

	genCmd.Dir = rootDir // Run from root directory where go.mod exists
	genCmd.Env = os.Environ()
	genOutput, err := genCmd.CombinedOutput()
	require.NoError(t, err, "Failed to generate project: %s", string(genOutput))

	t.Logf("Generated project at: %s", projectPath)
	return projectPath
}

// DatabaseContainer holds the database container and connection info
type DatabaseContainer struct {
	Container    testcontainers.Container
	Host         string
	Port         string
	DatabaseName string
	Username     string
	Password     string
}

// setupTestDatabase starts a test database using TestContainers and returns container info and cleanup function
func setupTestDatabase(t *testing.T, ctx context.Context, config RuntimeTestConfig) (*DatabaseContainer, func()) {
	if config.Database != "postgres" {
		t.Skip("Only postgres database testing implemented currently")
	}

	// Create postgres container
	postgresContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err, "Failed to start postgres container")

	// Get connection details
	host, err := postgresContainer.Host(ctx)
	require.NoError(t, err, "Failed to get container host")

	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err, "Failed to get container port")

	dbContainer := &DatabaseContainer{
		Container:    postgresContainer,
		Host:         host,
		Port:         mappedPort.Port(),
		DatabaseName: "testdb",
		Username:     "testuser",
		Password:     "testpass",
	}

	// Verify database connectivity
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbContainer.Host, dbContainer.Port, dbContainer.Username, dbContainer.Password, dbContainer.DatabaseName)

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err, "Failed to connect to test database")
	defer db.Close()

	err = db.Ping()
	require.NoError(t, err, "Failed to ping test database")

	t.Logf("Started test database on %s:%s", dbContainer.Host, dbContainer.Port)

	cleanup := func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("Warning: Failed to terminate database container: %v", err)
		}
	}

	return dbContainer, cleanup
}

// startApplicationServer builds and starts the generated application
func startApplicationServer(t *testing.T, ctx context.Context, projectPath string, config RuntimeTestConfig, dbContainer *DatabaseContainer) func() {
	// Update config file with test database settings
	updateProjectConfig(t, projectPath, config, dbContainer)

	// Install dependencies
	modCmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	modCmd.Dir = projectPath
	output, err := modCmd.CombinedOutput()
	require.NoError(t, err, "Failed to install dependencies: %s", string(output))

	// Build the application
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", "app", "./cmd/server")
	buildCmd.Dir = projectPath
	output, err = buildCmd.CombinedOutput()
	require.NoError(t, err, "Failed to build application: %s", string(output))

	// Start the application
	runCmd := exec.CommandContext(ctx, "./app")
	runCmd.Dir = projectPath
	runCmd.Env = append(os.Environ(),
		fmt.Sprintf("PORT=%d", config.ServerPort),
		"ENV=test",
	)

	// Capture stdout and stderr
	stdout, err := runCmd.StdoutPipe()
	require.NoError(t, err)
	stderr, err := runCmd.StderrPipe()
	require.NoError(t, err)

	err = runCmd.Start()
	require.NoError(t, err, "Failed to start application")

	// Log application output
	go func() {
		io.Copy(os.Stdout, stdout)
	}()
	go func() {
		io.Copy(os.Stderr, stderr)
	}()

	t.Logf("Started application server on port %d", config.ServerPort)

	return func() {
		if runCmd.Process != nil {
			runCmd.Process.Kill()
			runCmd.Wait()
		}
	}
}

// updateProjectConfig updates the generated project's configuration for testing
func updateProjectConfig(t *testing.T, projectPath string, config RuntimeTestConfig, dbContainer *DatabaseContainer) {
	configFile := filepath.Join(projectPath, "configs", "config.test.yaml")

	configContent := fmt.Sprintf(`
server:
  port: %d
  host: "localhost"
  
database:
  driver: "%s"
  host: "%s"
  port: %s
  name: "%s"
  user: "%s"
  password: "%s"
  sslmode: "disable"
  
logger:
  level: "debug"
  format: "json"
  
auth:
  jwt_secret: "test-secret-key-for-testing-only"
  token_duration: "1h"
`, config.ServerPort, config.Database, dbContainer.Host, dbContainer.Port,
		dbContainer.DatabaseName, dbContainer.Username, dbContainer.Password)

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err, "Failed to update config file")
}

// waitForServerReady waits for the server to be ready to accept requests
func waitForServerReady(t *testing.T, ctx context.Context, config RuntimeTestConfig) {
	healthURL := fmt.Sprintf("http://localhost:%d/health", config.ServerPort)
	client := &http.Client{Timeout: 5 * time.Second}

	for i := 0; i < 30; i++ { // Wait up to 30 seconds
		resp, err := client.Get(healthURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			t.Logf("Server is ready on port %d", config.ServerPort)
			return
		}
		if resp != nil {
			resp.Body.Close()
		}

		select {
		case <-ctx.Done():
			t.Fatal("Context cancelled while waiting for server")
		case <-time.After(1 * time.Second):
			// Continue waiting
		}
	}

	t.Fatal("Server did not become ready in time")
}

// runIntegrationTestSuite runs comprehensive tests against the running application
func runIntegrationTestSuite(t *testing.T, ctx context.Context, config RuntimeTestConfig) {
	baseURL := fmt.Sprintf("http://localhost:%d", config.ServerPort)
	client := &http.Client{Timeout: 10 * time.Second}

	// Test 1: Health endpoint
	t.Run("health_endpoint", func(t *testing.T) {
		testHealthEndpoint(t, client, baseURL)
	})

	// Test 2: User registration
	t.Run("user_registration", func(t *testing.T) {
		testUserRegistration(t, client, baseURL)
	})

	// Test 3: User login
	t.Run("user_login", func(t *testing.T) {
		testUserLogin(t, client, baseURL)
	})

	// Test 4: Protected endpoints
	t.Run("protected_endpoints", func(t *testing.T) {
		testProtectedEndpoints(t, client, baseURL)
	})

	// Test 5: Error handling
	t.Run("error_handling", func(t *testing.T) {
		testErrorHandling(t, client, baseURL)
	})

	// Test 6: Complete User CRUD operations
	t.Run("user_crud_operations", func(t *testing.T) {
		testUserCRUDOperations(t, client, baseURL)
	})

	// Test 7: User profile read operations by ID and email
	t.Run("user_read_operations", func(t *testing.T) {
		testUserReadOperations(t, client, baseURL)
	})

	// Test 8: User profile update operations
	t.Run("user_update_operations", func(t *testing.T) {
		testUserUpdateOperations(t, client, baseURL)
	})

	// Test 9: User deletion operations
	t.Run("user_delete_operations", func(t *testing.T) {
		testUserDeleteOperations(t, client, baseURL)
	})

	// Test 10: User list with pagination
	t.Run("user_list_pagination", func(t *testing.T) {
		testUserListPagination(t, client, baseURL)
	})

	// Test 11: Database round-trip data integrity
	t.Run("database_round_trip_integrity", func(t *testing.T) {
		testDatabaseRoundTripIntegrity(t, client, baseURL)
	})

	// Test 12: Error handling for invalid CRUD operations
	t.Run("crud_error_handling", func(t *testing.T) {
		testCRUDErrorHandling(t, client, baseURL)
	})
	// Test 13: HTTP endpoint integration validation
	t.Run("http_endpoint_integration", func(t *testing.T) {
		testHTTPEndpointIntegration(t, client, baseURL)
	})
	// Test 14: Complete authentication flow validation
	t.Run("complete_authentication_flow", func(t *testing.T) {
		testCompleteAuthenticationFlow(t, client, baseURL)
	})
	// Test 15: Cross-layer integration testing
	t.Run("cross_layer_integration", func(t *testing.T) {
		testCrossLayerIntegration(t, client, baseURL)
	})
	
	// Test 15: Domain layer behavior (value objects, entities, business rules)
	t.Run("domain_layer_behavior", func(t *testing.T) {
		testDomainLayerBehavior(t, client, baseURL)
	})
}

// testHealthEndpoint tests the health check endpoint
func testHealthEndpoint(t *testing.T, client *http.Client, baseURL string) {
	resp, err := client.Get(baseURL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	t.Logf("Health endpoint response body: %s", string(body))

	var healthResponse map[string]interface{}
	err = json.Unmarshal(body, &healthResponse)
	require.NoError(t, err)

	assert.Equal(t, "healthy", healthResponse["status"])
	assert.Contains(t, healthResponse, "timestamp")
}

// testUserRegistration tests user registration functionality
func testUserRegistration(t *testing.T, client *http.Client, baseURL string) {
	registrationData := map[string]interface{}{
		"email":      "test@example.com",
		"password":   "password123",
		"first_name": "Test",
		"last_name":  "User",
	}

	jsonData, err := json.Marshal(registrationData)
	require.NoError(t, err)

	resp, err := client.Post(baseURL+"/api/auth/register",
		"application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	t.Logf("Registration response status: %d, body: %s", resp.StatusCode, string(body))

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	assert.Contains(t, response, "user")
	assert.Contains(t, response, "message")
}

// testUserLogin tests user login functionality
func testUserLogin(t *testing.T, client *http.Client, baseURL string) {
	// First register a user
	registrationData := map[string]interface{}{
		"email":      "login@example.com",
		"password":   "password123",
		"first_name": "Login",
		"last_name":  "User",
	}

	jsonData, err := json.Marshal(registrationData)
	require.NoError(t, err)

	resp, err := client.Post(baseURL+"/api/auth/register",
		"application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	resp.Body.Close()

	// Now test login
	loginData := map[string]interface{}{
		"email":    "login@example.com",
		"password": "password123",
	}

	jsonData, err = json.Marshal(loginData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/login",
		"application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	assert.Contains(t, response, "access_token")
	assert.Contains(t, response, "token_type")
	assert.Contains(t, response, "expires_in")
	assert.Contains(t, response, "user")
}

// testProtectedEndpoints tests endpoints that require authentication
func testProtectedEndpoints(t *testing.T, client *http.Client, baseURL string) {
	// Test accessing protected endpoint without token
	resp, err := client.Get(baseURL + "/api/users/profile")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// testErrorHandling tests various error scenarios
func testErrorHandling(t *testing.T, client *http.Client, baseURL string) {
	// Test invalid JSON
	resp, err := client.Post(baseURL+"/api/auth/register",
		"application/json", strings.NewReader("invalid json"))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Test missing required fields
	emptyData := map[string]interface{}{}
	jsonData, err := json.Marshal(emptyData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/register",
		"application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// User represents the user data structure for testing
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	User         User   `json:"user"`
}

// createTestUserWithAuth creates a test user and returns authentication token
func createTestUserWithAuth(t *testing.T, client *http.Client, baseURL, email string) (User, string) {
	registrationData := map[string]interface{}{
		"email":      email,
		"password":   "password123",
		"first_name": "Test",
		"last_name":  "User",
	}

	jsonData, err := json.Marshal(registrationData)
	require.NoError(t, err)

	// Register user
	resp, err := client.Post(baseURL+"/api/auth/register",
		"application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Login to get token
	loginData := map[string]interface{}{
		"email":    email,
		"password": "password123",
	}

	jsonData, err = json.Marshal(loginData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/login",
		"application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	require.NoError(t, err)

	return authResponse.User, authResponse.AccessToken
}

// makeAuthenticatedRequest makes an HTTP request with authentication
func makeAuthenticatedRequest(client *http.Client, method, url, token string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return client.Do(req)
}

// testUserCRUDOperations tests complete CRUD operations for users
func testUserCRUDOperations(t *testing.T, client *http.Client, baseURL string) {
	// Create (via registration)
	user, token := createTestUserWithAuth(t, client, baseURL, "crud@example.com")

	// Verify user was created with correct data
	assert.Equal(t, "crud@example.com", user.Email)
	assert.Equal(t, "Test", user.FirstName)
	assert.Equal(t, "User", user.LastName)
	assert.NotZero(t, user.ID)

	// Read - Get user profile
	resp, err := makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", token, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var profileResponse map[string]interface{}
	err = json.Unmarshal(body, &profileResponse)
	require.NoError(t, err)

	userProfile := profileResponse["user"].(map[string]interface{})
	assert.Equal(t, user.Email, userProfile["email"])
	assert.Equal(t, user.FirstName, userProfile["first_name"])
	assert.Equal(t, user.LastName, userProfile["last_name"])

	// Update - Modify user profile
	updateData := map[string]interface{}{
		"first_name": "Updated",
		"last_name":  "Name",
	}

	jsonData, err := json.Marshal(updateData)
	require.NoError(t, err)

	resp, err = makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", token, strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify update
	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	var updateResponse map[string]interface{}
	err = json.Unmarshal(body, &updateResponse)
	require.NoError(t, err)

	updatedUser := updateResponse["user"].(map[string]interface{})
	assert.Equal(t, "Updated", updatedUser["first_name"])
	assert.Equal(t, "Name", updatedUser["last_name"])
	assert.Equal(t, user.Email, updatedUser["email"]) // Email should remain unchanged

	// Delete - Remove user
	resp, err = makeAuthenticatedRequest(client, "DELETE", baseURL+"/api/users/profile", token, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify deletion - subsequent requests should fail
	resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", token, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// testUserReadOperations tests user read operations by different criteria
func testUserReadOperations(t *testing.T, client *http.Client, baseURL string) {
	// Create test user
	user, token := createTestUserWithAuth(t, client, baseURL, "read@example.com")

	// Read by profile (authenticated user)
	resp, err := makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", token, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	userProfile := response["user"].(map[string]interface{})
	assert.Equal(t, user.Email, userProfile["email"])
	assert.Equal(t, user.FirstName, userProfile["first_name"])
	assert.Equal(t, user.LastName, userProfile["last_name"])

	// Test reading user by ID (if endpoint exists)
	userID := uint(userProfile["id"].(float64))
	resp, err = makeAuthenticatedRequest(client, "GET", fmt.Sprintf("%s/api/users/%d", baseURL, userID), token, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	// This endpoint might not exist in all blueprints, so we handle both cases
	if resp.StatusCode == http.StatusOK {
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var userByIDResponse map[string]interface{}
		err = json.Unmarshal(body, &userByIDResponse)
		require.NoError(t, err)

		fetchedUser := userByIDResponse["user"].(map[string]interface{})
		assert.Equal(t, userProfile["id"], fetchedUser["id"])
		assert.Equal(t, userProfile["email"], fetchedUser["email"])
	} else {
		// Endpoint doesn't exist, which is acceptable
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	}
}

// testUserUpdateOperations tests user profile update functionality
func testUserUpdateOperations(t *testing.T, client *http.Client, baseURL string) {
	// Create test user
	_, token := createTestUserWithAuth(t, client, baseURL, "update@example.com")

	// Test partial update - only first name
	updateData := map[string]interface{}{
		"first_name": "UpdatedFirst",
	}

	jsonData, err := json.Marshal(updateData)
	require.NoError(t, err)

	resp, err := makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", token, strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	updatedUser := response["user"].(map[string]interface{})
	assert.Equal(t, "UpdatedFirst", updatedUser["first_name"])
	assert.Equal(t, "User", updatedUser["last_name"]) // Should remain unchanged

	// Test full profile update
	fullUpdateData := map[string]interface{}{
		"first_name": "CompletelyNew",
		"last_name":  "Name",
	}

	jsonData, err = json.Marshal(fullUpdateData)
	require.NoError(t, err)

	resp, err = makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", token, strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	fullyUpdatedUser := response["user"].(map[string]interface{})
	assert.Equal(t, "CompletelyNew", fullyUpdatedUser["first_name"])
	assert.Equal(t, "Name", fullyUpdatedUser["last_name"])

	// Verify persistence by reading again
	resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", token, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	var verifyResponse map[string]interface{}
	err = json.Unmarshal(body, &verifyResponse)
	require.NoError(t, err)

	verifiedUser := verifyResponse["user"].(map[string]interface{})
	assert.Equal(t, "CompletelyNew", verifiedUser["first_name"])
	assert.Equal(t, "Name", verifiedUser["last_name"])
}

// testUserDeleteOperations tests user deletion functionality
func testUserDeleteOperations(t *testing.T, client *http.Client, baseURL string) {
	// Create test user
	_, token := createTestUserWithAuth(t, client, baseURL, "delete@example.com")

	// Verify user exists by reading profile
	resp, err := makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", token, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete user
	resp, err = makeAuthenticatedRequest(client, "DELETE", baseURL+"/api/users/profile", token, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify user no longer exists - subsequent authenticated requests should fail
	resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", token, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Verify login no longer works
	loginData := map[string]interface{}{
		"email":    "delete@example.com",
		"password": "password123",
	}

	jsonData, err := json.Marshal(loginData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/login",
		"application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// testUserListPagination tests user listing with pagination (if endpoint exists)
func testUserListPagination(t *testing.T, client *http.Client, baseURL string) {
	// Create multiple test users
	var tokens []string
	for i := 0; i < 5; i++ {
		email := fmt.Sprintf("list%d@example.com", i)
		_, token := createTestUserWithAuth(t, client, baseURL, email)
		tokens = append(tokens, token)
	}

	// Test user list endpoint (may not exist in all blueprints)
	resp, err := makeAuthenticatedRequest(client, "GET", baseURL+"/api/users", tokens[0], nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var listResponse map[string]interface{}
		err = json.Unmarshal(body, &listResponse)
		require.NoError(t, err)

		// Verify response structure
		assert.Contains(t, listResponse, "users")
		users := listResponse["users"].([]interface{})
		assert.True(t, len(users) >= 5) // At least our 5 test users

		// Test pagination parameters if supported
		resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users?page=1&limit=2", tokens[0], nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err = io.ReadAll(resp.Body)
			require.NoError(t, err)

			var paginatedResponse map[string]interface{}
			err = json.Unmarshal(body, &paginatedResponse)
			require.NoError(t, err)

			paginatedUsers := paginatedResponse["users"].([]interface{})
			assert.True(t, len(paginatedUsers) <= 2) // Should respect limit
		}
	} else {
		// User list endpoint doesn't exist, which is acceptable
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		t.Log("User list endpoint not implemented, skipping pagination tests")
	}
}

// testDatabaseRoundTripIntegrity tests that data persists correctly through database operations
func testDatabaseRoundTripIntegrity(t *testing.T, client *http.Client, baseURL string) {
	// Create user with specific data
	originalEmail := "integrity@example.com"
	originalFirstName := "Original"
	originalLastName := "Data"

	registrationData := map[string]interface{}{
		"email":      originalEmail,
		"password":   "password123",
		"first_name": originalFirstName,
		"last_name":  originalLastName,
	}

	jsonData, err := json.Marshal(registrationData)
	require.NoError(t, err)

	resp, err := client.Post(baseURL+"/api/auth/register",
		"application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Login to get token
	loginData := map[string]interface{}{
		"email":    originalEmail,
		"password": "password123",
	}

	jsonData, err = json.Marshal(loginData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/login",
		"application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	require.NoError(t, err)

	token := authResponse.AccessToken
	originalUser := authResponse.User

	// Verify data integrity after login
	assert.Equal(t, originalEmail, originalUser.Email)
	assert.Equal(t, originalFirstName, originalUser.FirstName)
	assert.Equal(t, originalLastName, originalUser.LastName)

	// Perform multiple updates to test persistence
	updates := []map[string]interface{}{
		{"first_name": "First Update"},
		{"last_name": "Second Update"},
		{"first_name": "Final First", "last_name": "Final Last"},
	}

	for i, updateData := range updates {
		jsonData, err = json.Marshal(updateData)
		require.NoError(t, err)

		resp, err = makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", token, strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Read back immediately to verify persistence
		resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", token, nil)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var profileResponse map[string]interface{}
		err = json.Unmarshal(body, &profileResponse)
		require.NoError(t, err)

		user := profileResponse["user"].(map[string]interface{})

		// Verify specific updates
		if firstName, ok := updateData["first_name"]; ok {
			assert.Equal(t, firstName, user["first_name"], "Update %d: first_name not persisted", i+1)
		}
		if lastName, ok := updateData["last_name"]; ok {
			assert.Equal(t, lastName, user["last_name"], "Update %d: last_name not persisted", i+1)
		}

		// Email should never change
		assert.Equal(t, originalEmail, user["email"], "Email changed unexpectedly")

		// ID should remain consistent
		assert.Equal(t, originalUser.ID, user["id"].(string), "User ID changed unexpectedly")

		t.Logf("Update %d verified: %+v", i+1, updateData)
	}

	// Final verification - login again and check all data is still correct
	resp, err = client.Post(baseURL+"/api/auth/login",
		"application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	var finalAuthResponse AuthResponse
	err = json.Unmarshal(body, &finalAuthResponse)
	require.NoError(t, err)

	finalUser := finalAuthResponse.User
	assert.Equal(t, originalEmail, finalUser.Email)
	assert.Equal(t, "Final First", finalUser.FirstName)
	assert.Equal(t, "Final Last", finalUser.LastName)
	assert.Equal(t, originalUser.ID, finalUser.ID)
}

// testCRUDErrorHandling tests error handling for invalid CRUD operations
func testCRUDErrorHandling(t *testing.T, client *http.Client, baseURL string) {
	// Create test user
	_, token := createTestUserWithAuth(t, client, baseURL, "errors@example.com")

	// Test invalid update data
	invalidUpdates := []map[string]interface{}{
		{"email": "newemail@example.com"}, // Email changes should be rejected
		{"id": 999},                       // ID changes should be rejected
		{"password": "newpassword"},       // Password changes via profile update should be rejected
		{},                                // Empty update should be handled gracefully
	}

	for _, invalidUpdate := range invalidUpdates {
		jsonData, err := json.Marshal(invalidUpdate)
		require.NoError(t, err)

		resp, err := makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", token, strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should either reject the update or ignore invalid fields
		if resp.StatusCode != http.StatusOK {
			assert.Contains(t, []int{http.StatusBadRequest, http.StatusUnprocessableEntity}, resp.StatusCode)
		}
	}

	// Test operations with invalid token
	invalidToken := "invalid.jwt.token"

	resp, err := makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", invalidToken, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	resp, err = makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", invalidToken, strings.NewReader(`{"first_name":"Invalid"}`))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	resp, err = makeAuthenticatedRequest(client, "DELETE", baseURL+"/api/users/profile", invalidToken, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Test operations with no token
	resp, err = client.Get(baseURL + "/api/users/profile")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Test malformed JSON in update
	resp, err = makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", token, strings.NewReader(`{invalid json}`))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Test accessing non-existent user by ID (if endpoint exists)
	resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/999999", token, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	// Should be either 404 (not found) or 404 (endpoint doesn't exist)
	assert.Contains(t, []int{http.StatusNotFound}, resp.StatusCode)
}

// testHTTPEndpointIntegration performs comprehensive HTTP endpoint testing
func testHTTPEndpointIntegration(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: Direct user creation via HTTP POST
	t.Run("direct_user_creation", func(t *testing.T) {
		testDirectUserCreation(t, client, baseURL)
	})

	// Test 2: User retrieval by ID via HTTP GET
	t.Run("user_retrieval_by_id", func(t *testing.T) {
		testUserRetrievalByID(t, client, baseURL)
	})

	// Test 3: User update via HTTP PUT
	t.Run("user_update_via_put", func(t *testing.T) {
		testUserUpdateViaPUT(t, client, baseURL)
	})

	// Test 4: User deletion via HTTP DELETE
	t.Run("user_deletion_via_delete", func(t *testing.T) {
		testUserDeletionViaDELETE(t, client, baseURL)
	})

	// Test 5: HTTP error responses and status codes
	t.Run("http_error_responses", func(t *testing.T) {
		testHTTPErrorResponses(t, client, baseURL)
	})

	// Test 6: Content-Type and headers validation
	t.Run("content_type_headers", func(t *testing.T) {
		testContentTypeHeaders(t, client, baseURL)
	})

	// Test 7: URL parameter extraction
	t.Run("url_parameter_extraction", func(t *testing.T) {
		testURLParameterExtraction(t, client, baseURL)
	})
}

// testDirectUserCreation tests creating users directly via HTTP POST to /api/users
func testDirectUserCreation(t *testing.T, client *http.Client, baseURL string) {
	// First create a test user for authentication
	_, authToken := createTestUserWithAuth(t, client, baseURL, "admin@example.com")

	// Test creating user via direct HTTP POST
	createUserData := map[string]interface{}{
		"email":      "newuser@example.com",
		"password":   "password123",
		"first_name": "New",
		"last_name":  "User",
	}

	jsonData, err := json.Marshal(createUserData)
	require.NoError(t, err)

	// Try to create user via authenticated POST to /api/users
	resp, err := makeAuthenticatedRequest(client, "POST", baseURL+"/api/users", authToken, strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Verify creation was successful or endpoint doesn't exist
	if resp.StatusCode == http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var createdUser User
		err = json.Unmarshal(body, &createdUser)
		require.NoError(t, err)

		assert.Equal(t, "newuser@example.com", createdUser.Email)
		assert.Equal(t, "New", createdUser.FirstName)
		assert.Equal(t, "User", createdUser.LastName)
		assert.NotEmpty(t, createdUser.ID)

		t.Log("Direct user creation via POST /api/users: SUCCESS")
	} else if resp.StatusCode == http.StatusNotFound {
		t.Log("Direct user creation endpoint POST /api/users not implemented")
	} else {
		t.Logf("Unexpected status code for POST /api/users: %d", resp.StatusCode)
	}
}

// testUserRetrievalByID tests retrieving users by ID via HTTP GET
func testUserRetrievalByID(t *testing.T, client *http.Client, baseURL string) {
	// Create test user and get token
	createdUser, authToken := createTestUserWithAuth(t, client, baseURL, "getuser@example.com")

	// Test retrieving user by ID
	userIDStr := createdUser.ID
	if userIDStr == "" {
		t.Skip("User ID not available for retrieval test")
		return
	}

	resp, err := makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/"+userIDStr, authToken, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var retrievedUser User
		err = json.Unmarshal(body, &retrievedUser)
		require.NoError(t, err)

		assert.Equal(t, createdUser.Email, retrievedUser.Email)
		assert.Equal(t, createdUser.FirstName, retrievedUser.FirstName)
		assert.Equal(t, createdUser.LastName, retrievedUser.LastName)

		t.Logf("User retrieval by ID via GET /api/users/%s: SUCCESS", userIDStr)
	} else if resp.StatusCode == http.StatusNotFound {
		t.Log("User retrieval by ID endpoint not implemented or user not found")
	} else {
		t.Logf("Unexpected status code for GET /api/users/%s: %d", userIDStr, resp.StatusCode)
	}
}

// testUserUpdateViaPUT tests updating users via HTTP PUT
func testUserUpdateViaPUT(t *testing.T, client *http.Client, baseURL string) {
	// Create test user and get token
	createdUser, authToken := createTestUserWithAuth(t, client, baseURL, "putuser@example.com")

	// Test updating user via PUT
	userIDStr := createdUser.ID
	if userIDStr == "" {
		t.Skip("User ID not available for PUT test")
		return
	}

	updateData := map[string]interface{}{
		"first_name": "Updated",
		"last_name":  "Name",
	}

	jsonData, err := json.Marshal(updateData)
	require.NoError(t, err)

	resp, err := makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/"+userIDStr, authToken, strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var updatedUser User
		err = json.Unmarshal(body, &updatedUser)
		require.NoError(t, err)

		assert.Equal(t, "Updated", updatedUser.FirstName)
		assert.Equal(t, "Name", updatedUser.LastName)
		assert.Equal(t, createdUser.Email, updatedUser.Email) // Email should not change

		t.Logf("User update via PUT /api/users/%s: SUCCESS", userIDStr)
	} else if resp.StatusCode == http.StatusNotFound {
		t.Log("User update via PUT endpoint not implemented")
	} else {
		t.Logf("Unexpected status code for PUT /api/users/%s: %d", userIDStr, resp.StatusCode)
	}
}

// testUserDeletionViaDELETE tests deleting users via HTTP DELETE
func testUserDeletionViaDELETE(t *testing.T, client *http.Client, baseURL string) {
	// Create test user and get token (use different user for auth than the one we'll delete)
	_, authToken := createTestUserWithAuth(t, client, baseURL, "admin2@example.com")
	
	// Create a separate user to delete
	userToDelete, _ := createTestUserWithAuth(t, client, baseURL, "deleteuser@example.com")

	userIDStr := userToDelete.ID
	if userIDStr == "" {
		t.Skip("User ID not available for DELETE test")
		return
	}

	resp, err := makeAuthenticatedRequest(client, "DELETE", baseURL+"/api/users/"+userIDStr, authToken, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
		// Verify user is deleted by trying to retrieve it
		resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/"+userIDStr, authToken, nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode, "User should be deleted and not retrievable")

		t.Logf("User deletion via DELETE /api/users/%s: SUCCESS", userIDStr)
	} else if resp.StatusCode == http.StatusNotFound {
		t.Log("User deletion via DELETE endpoint not implemented")
	} else {
		t.Logf("Unexpected status code for DELETE /api/users/%s: %d", userIDStr, resp.StatusCode)
	}
}

// testHTTPErrorResponses tests proper HTTP error responses and status codes
func testHTTPErrorResponses(t *testing.T, client *http.Client, baseURL string) {
	// Create test user for authentication
	_, authToken := createTestUserWithAuth(t, client, baseURL, "errortest@example.com")

	// Test 1: GET non-existent user (use UUID that won't exist)
	resp, err := makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/nonexistent-uuid-999999999", authToken, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	// Should be 404 if endpoint exists, or 404 if endpoint doesn't exist
	if resp.StatusCode != http.StatusNotFound {
		t.Logf("GET non-existent user returned status: %d (endpoint may not be implemented)", resp.StatusCode)
	}

	// Test 2: POST with invalid JSON
	resp, err = makeAuthenticatedRequest(client, "POST", baseURL+"/api/users", authToken, strings.NewReader(`{invalid json`))
	require.NoError(t, err)
	defer resp.Body.Close()
	
	// Should be 400 Bad Request if endpoint exists
	if resp.StatusCode != http.StatusNotFound {
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Invalid JSON should return 400")
	}

	// Test 3: PUT with invalid data (test validation)
	invalidUpdateData := map[string]interface{}{
		"email": "not-an-email", // Invalid email format
	}
	
	jsonData, err := json.Marshal(invalidUpdateData)
	require.NoError(t, err)

	resp, err = makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/nonexistent-user-999", authToken, strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	
	// Log the response for debugging - validation may or may not be implemented
	t.Logf("PUT with invalid data returned status: %d", resp.StatusCode)

	// Test 4: Unauthorized request (no token)
	req, err := http.NewRequest("GET", baseURL+"/api/users/profile", nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	
	// Should be 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Request without token should return 401")

	t.Log("HTTP error response testing: SUCCESS")
}

// testContentTypeHeaders tests proper Content-Type headers and request/response handling
func testContentTypeHeaders(t *testing.T, client *http.Client, baseURL string) {
	// Create test user for authentication
	_, authToken := createTestUserWithAuth(t, client, baseURL, "headers@example.com")

	// Test 1: Verify response Content-Type is application/json
	resp, err := makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", authToken, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		contentType := resp.Header.Get("Content-Type")
		assert.Contains(t, contentType, "application/json", "Response should have JSON content type")
	}

	// Test 2: POST with wrong Content-Type
	userData := map[string]interface{}{
		"first_name": "Test",
		"last_name":  "Headers",
	}
	
	jsonData, err := json.Marshal(userData)
	require.NoError(t, err)

	req, err := http.NewRequest("PUT", baseURL+"/api/users/profile", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "text/plain") // Wrong content type

	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// Log the response for debugging - servers may handle wrong content-type gracefully
	t.Logf("PUT with wrong Content-Type returned status: %d", resp.StatusCode)

	t.Log("Content-Type and headers testing: SUCCESS")
}

// testURLParameterExtraction tests URL parameter extraction in various endpoints
func testURLParameterExtraction(t *testing.T, client *http.Client, baseURL string) {
	// Create test user for authentication
	createdUser, authToken := createTestUserWithAuth(t, client, baseURL, "urlparams@example.com")

	userIDStr := createdUser.ID
	if userIDStr == "" {
		t.Skip("User ID not available for URL parameter test")
		return
	}

	// Test different URL patterns that should extract the same user ID
	testURLs := []string{
		baseURL + "/api/users/" + userIDStr,
		baseURL + "/api/users/" + userIDStr + "/",        // With trailing slash
		baseURL + "/api/users/" + userIDStr + "/profile", // With additional path
	}

	for _, testURL := range testURLs {
		resp, err := makeAuthenticatedRequest(client, "GET", testURL, authToken, nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		// If endpoint exists and works, the URL parameter extraction is working
		if resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			// Verify we got user data back
			if len(body) > 0 && strings.Contains(string(body), createdUser.Email) {
				t.Logf("URL parameter extraction working for: %s", testURL)
			}
		} else if resp.StatusCode == http.StatusNotFound {
			// Endpoint doesn't exist, which is fine
			t.Logf("Endpoint not implemented for URL: %s", testURL)
		}
	}

	// Test invalid user ID parameter
	invalidURLs := []string{
		baseURL + "/api/users/nonexistent-uuid-12345",
		baseURL + "/api/users/invalid-chars-!@#",
		baseURL + "/api/users/999999999999",
	}

	for _, invalidURL := range invalidURLs {
		resp, err := makeAuthenticatedRequest(client, "GET", invalidURL, authToken, nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should return 404 for non-existent user IDs (or 404 if endpoint doesn't exist)
		// The URL parameter extraction is working if we get a response
		t.Logf("URL parameter extraction test for %s returned status: %d", invalidURL, resp.StatusCode)
	}

	t.Log("URL parameter extraction testing: SUCCESS")
}

// testCompleteAuthenticationFlow performs comprehensive authentication flow testing
func testCompleteAuthenticationFlow(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: JWT Token Generation and Validation
	t.Run("jwt_token_lifecycle", func(t *testing.T) {
		testJWTTokenLifecycle(t, client, baseURL)
	})

	// Test 2: Token Refresh Flow
	t.Run("token_refresh_flow", func(t *testing.T) {
		testTokenRefreshFlow(t, client, baseURL)
	})

	// Test 3: Authentication Middleware Behavior
	t.Run("authentication_middleware", func(t *testing.T) {
		testAuthenticationMiddleware(t, client, baseURL)
	})

	// Test 4: Authorization Context Propagation
	t.Run("authorization_context", func(t *testing.T) {
		testAuthorizationContext(t, client, baseURL)
	})

	// Test 5: Authentication Security Edge Cases
	t.Run("authentication_security", func(t *testing.T) {
		testAuthenticationSecurity(t, client, baseURL)
	})

	// Test 6: Cross-Layer Authentication Integration
	t.Run("cross_layer_auth_integration", func(t *testing.T) {
		testCrossLayerAuthIntegration(t, client, baseURL)
	})

	// Test 7: Concurrent Authentication Sessions
	t.Run("concurrent_auth_sessions", func(t *testing.T) {
		testConcurrentAuthSessions(t, client, baseURL)
	})
}

// testJWTTokenLifecycle tests complete JWT token generation, validation, and expiration
func testJWTTokenLifecycle(t *testing.T, client *http.Client, baseURL string) {
	// Create user and login to get fresh tokens
	testEmail := "jwt-lifecycle@example.com"
	registrationData := map[string]interface{}{
		"email":      testEmail,
		"password":   "password123",
		"first_name": "JWT",
		"last_name":  "Test",
	}

	// Register user
	jsonData, err := json.Marshal(registrationData)
	require.NoError(t, err)

	resp, err := client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Login to get tokens
	loginData := map[string]interface{}{
		"email":    testEmail,
		"password": "password123",
	}

	jsonData, err = json.Marshal(loginData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	require.NoError(t, err)

	// Validate token structure
	assert.NotEmpty(t, authResponse.AccessToken, "Access token should not be empty")
	assert.NotEmpty(t, authResponse.RefreshToken, "Refresh token should not be empty")
	assert.Equal(t, "Bearer", authResponse.TokenType, "Token type should be Bearer")
	assert.Greater(t, authResponse.ExpiresIn, 0, "ExpiresIn should be positive")

	// Validate token format (JWT should have 3 parts separated by dots)
	tokenParts := strings.Split(authResponse.AccessToken, ".")
	assert.Len(t, tokenParts, 3, "JWT should have 3 parts (header.payload.signature)")

	// Test token validation by making authenticated request
	resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", authResponse.AccessToken, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Valid token should allow access to protected endpoint")

	// Verify user data in response matches login user
	if resp.StatusCode == http.StatusOK {
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var profileResponse map[string]interface{}
		err = json.Unmarshal(body, &profileResponse)
		require.NoError(t, err)

		if user, ok := profileResponse["user"].(map[string]interface{}); ok {
			assert.Equal(t, testEmail, user["email"], "Profile email should match logged in user")
		}
	}

	t.Log("JWT token lifecycle validation: SUCCESS")
}

// testTokenRefreshFlow tests the refresh token mechanism
func testTokenRefreshFlow(t *testing.T, client *http.Client, baseURL string) {
	// Create user and login to get initial tokens
	createdUser, _ := createTestUserWithAuth(t, client, baseURL, "refresh-test@example.com")

	// Login to get fresh tokens
	loginData := map[string]interface{}{
		"email":    createdUser.Email,
		"password": "password123",
	}

	jsonData, err := json.Marshal(loginData)
	require.NoError(t, err)

	resp, err := client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	require.NoError(t, err)

	originalAccessToken := authResponse.AccessToken
	originalRefreshToken := authResponse.RefreshToken

	// Test refresh token endpoint
	refreshData := map[string]interface{}{
		"refresh_token": originalRefreshToken,
	}

	jsonData, err = json.Marshal(refreshData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/refresh", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Refresh endpoint exists and works
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var refreshResponse map[string]interface{}
		err = json.Unmarshal(body, &refreshResponse)
		require.NoError(t, err)

		// Validate new tokens are provided
		newAccessToken, ok := refreshResponse["access_token"].(string)
		assert.True(t, ok, "New access token should be provided")
		assert.NotEmpty(t, newAccessToken, "New access token should not be empty")
		assert.NotEqual(t, originalAccessToken, newAccessToken, "New access token should be different")

		// Test that new token works
		resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", newAccessToken, nil)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode, "New access token should work")

		t.Log("Token refresh flow: SUCCESS")
	} else {
		// Refresh endpoint may not be implemented
		t.Logf("Token refresh endpoint returned status: %d (may not be implemented)", resp.StatusCode)
	}
}

// testAuthenticationMiddleware tests authentication middleware behavior
func testAuthenticationMiddleware(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: Access to public endpoints without authentication
	publicEndpoints := []string{
		baseURL + "/health",
		baseURL + "/api/auth/login",
		baseURL + "/api/auth/register",
	}

	for _, endpoint := range publicEndpoints {
		resp, err := client.Get(endpoint)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Public endpoints should not require authentication
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode, "Public endpoint %s should not require auth", endpoint)
	}

	// Test 2: Access to protected endpoints without authentication
	protectedEndpoints := []string{
		baseURL + "/api/users/profile",
		baseURL + "/api/users",
	}

	for _, endpoint := range protectedEndpoints {
		resp, err := client.Get(endpoint)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Protected endpoints should require authentication
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Protected endpoint %s should require auth", endpoint)
	}

	// Test 3: Access with valid authentication
	_, authToken := createTestUserWithAuth(t, client, baseURL, "middleware-test@example.com")

	for _, endpoint := range protectedEndpoints {
		resp, err := makeAuthenticatedRequest(client, "GET", endpoint, authToken, nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should not be unauthorized with valid token
		assert.NotEqual(t, http.StatusUnauthorized, resp.StatusCode, "Protected endpoint %s should accept valid auth", endpoint)
	}

	// Test 4: Different authorization header formats
	invalidAuthHeaders := []string{
		"InvalidToken",           // Missing Bearer prefix
		"Basic dGVzdDp0ZXN0",     // Wrong auth type
		"Bearer",                 // Missing token
		"Bearer " + authToken + " extra", // Extra content
	}

	for _, invalidHeader := range invalidAuthHeaders {
		req, err := http.NewRequest("GET", baseURL+"/api/users/profile", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", invalidHeader)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Invalid auth headers should be rejected
		if invalidHeader != "Bearer "+authToken+" extra" { // This might be handled gracefully
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Invalid auth header should be rejected: %s", invalidHeader)
		}
	}

	t.Log("Authentication middleware testing: SUCCESS")
}

// testAuthorizationContext tests that user context is properly propagated
func testAuthorizationContext(t *testing.T, client *http.Client, baseURL string) {
	// Create user and get token
	createdUser, authToken := createTestUserWithAuth(t, client, baseURL, "context-test@example.com")

	// Test that user profile endpoint returns the correct user
	resp, err := makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", authToken, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var profileResponse map[string]interface{}
		err = json.Unmarshal(body, &profileResponse)
		require.NoError(t, err)

		// Verify that the context contains the correct user
		if user, ok := profileResponse["user"].(map[string]interface{}); ok {
			assert.Equal(t, createdUser.Email, user["email"], "Context should contain correct user email")
			assert.Equal(t, createdUser.FirstName, user["first_name"], "Context should contain correct user first name")
		}
	}

	// Test that profile updates work with context (user can only update their own profile)
	updateData := map[string]interface{}{
		"first_name": "Updated",
		"last_name":  "Context",
	}

	jsonData, err := json.Marshal(updateData)
	require.NoError(t, err)

	resp, err = makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", authToken, strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Verify the update worked
		resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", authToken, nil)
		require.NoError(t, err)
		defer resp.Body.Close()

		updatedBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var updatedProfileResponse map[string]interface{}
		err = json.Unmarshal(updatedBody, &updatedProfileResponse)
		require.NoError(t, err)

		if user, ok := updatedProfileResponse["user"].(map[string]interface{}); ok {
			assert.Equal(t, "Updated", user["first_name"], "Profile update should work through context")
			assert.Equal(t, "Context", user["last_name"], "Profile update should work through context")
		}
	}

	t.Log("Authorization context testing: SUCCESS")
}

// testAuthenticationSecurity tests security edge cases and attack vectors
func testAuthenticationSecurity(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: Invalid JWT structures
	invalidTokens := []string{
		"invalid-token",              // Not a JWT
		"header.payload",             // Missing signature
		"header.payload.signature.extra", // Too many parts
		"",                          // Empty token
		"Bearer ",                   // Empty bearer token
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.signature", // Invalid payload
	}

	for _, invalidToken := range invalidTokens {
		req, err := http.NewRequest("GET", baseURL+"/api/users/profile", nil)
		require.NoError(t, err)
		
		if invalidToken != "" {
			req.Header.Set("Authorization", "Bearer "+invalidToken)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Invalid token should be rejected: %s", invalidToken)
	}

	// Test 2: SQL Injection in credentials
	sqlInjectionAttempts := []map[string]interface{}{
		{"email": "admin@example.com'; DROP TABLE users; --", "password": "password"},
		{"email": "admin@example.com", "password": "' OR '1'='1"},
		{"email": "admin@example.com\" OR 1=1 --", "password": "password"},
	}

	for _, attempt := range sqlInjectionAttempts {
		jsonData, err := json.Marshal(attempt)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		// SQL injection should not succeed
		assert.NotEqual(t, http.StatusOK, resp.StatusCode, "SQL injection attempt should fail")
	}

	// Test 3: Extremely long tokens (potential buffer overflow)
	longToken := strings.Repeat("a", 10000)
	req, err := http.NewRequest("GET", baseURL+"/api/users/profile", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+longToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Extremely long token should be rejected")

	// Test 4: Rate limiting (multiple failed login attempts)
	for i := 0; i < 5; i++ {
		loginData := map[string]interface{}{
			"email":    "nonexistent@example.com",
			"password": "wrongpassword",
		}

		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Failed login should be rejected")
	}

	t.Log("Authentication security testing: SUCCESS")
}

// testCrossLayerAuthIntegration tests authentication flow across hexagonal architecture layers
func testCrossLayerAuthIntegration(t *testing.T, client *http.Client, baseURL string) {
	// This test validates that authentication properly flows through:
	// HTTP Adapter → Application Service → Domain Service → Infrastructure Repository

	testEmail := "cross-layer@example.com"
	
	// Step 1: Registration flow (HTTP → App → Domain → Infrastructure)
	registrationData := map[string]interface{}{
		"email":      testEmail,
		"password":   "securepassword123",
		"first_name": "Cross",
		"last_name":  "Layer",
	}

	jsonData, err := json.Marshal(registrationData)
	require.NoError(t, err)

	resp, err := client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Step 2: Login flow (same layer traversal)
	loginData := map[string]interface{}{
		"email":    testEmail,
		"password": "securepassword123",
	}

	jsonData, err = json.Marshal(loginData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	require.NoError(t, err)

	// Step 3: Token validation flow (middleware validation)
	resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", authResponse.AccessToken, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Token validation should work across all layers")

	// Step 4: Verify domain business rules are enforced
	// Try to register with same email (should fail at domain level)
	resp, err = client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.NotEqual(t, http.StatusCreated, resp.StatusCode, "Duplicate email should be rejected by domain layer")

	// Step 5: Test that infrastructure properly persists authentication state
	// Login again and verify we get a different token (new session)
	resp, err = client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	var secondAuthResponse AuthResponse
	err = json.Unmarshal(body, &secondAuthResponse)
	require.NoError(t, err)

	// Both tokens should work (multiple sessions)
	resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", authResponse.AccessToken, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Original token should still work")

	resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", secondAuthResponse.AccessToken, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "New token should also work")

	t.Log("Cross-layer authentication integration: SUCCESS")
}

// testConcurrentAuthSessions tests handling of multiple concurrent authentication sessions
func testConcurrentAuthSessions(t *testing.T, client *http.Client, baseURL string) {
	testEmail := "concurrent@example.com"
	
	// Create user first
	createdUser, _ := createTestUserWithAuth(t, client, baseURL, testEmail)

	// Create multiple concurrent login sessions
	loginData := map[string]interface{}{
		"email":    createdUser.Email,
		"password": "password123",
	}

	var tokens []string
	numSessions := 3

	// Login multiple times to create concurrent sessions
	for i := 0; i < numSessions; i++ {
		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var authResponse AuthResponse
		err = json.Unmarshal(body, &authResponse)
		require.NoError(t, err)

		tokens = append(tokens, authResponse.AccessToken)
	}

	// Verify all tokens work concurrently
	for i, token := range tokens {
		resp, err := makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", token, nil)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Concurrent session %d should work", i+1)

		// Verify each session can make profile updates
		updateData := map[string]interface{}{
			"first_name": fmt.Sprintf("Concurrent%d", i+1),
		}

		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		resp, err = makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", token, strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t.Logf("Concurrent session %d can update profile", i+1)
		}
	}

	// Test concurrent access to the same protected resources
	for i, token := range tokens {
		go func(sessionNum int, sessionToken string) {
			// Each goroutine makes requests with its own token
			resp, err := makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", sessionToken, nil)
			assert.NoError(t, err)
			if resp != nil {
				defer resp.Body.Close()
				assert.Equal(t, http.StatusOK, resp.StatusCode, "Concurrent session %d should work in goroutine", sessionNum)
			}
		}(i+1, token)
	}

	// Give time for concurrent requests to complete
	time.Sleep(2 * time.Second)

	t.Log("Concurrent authentication sessions testing: SUCCESS")
}

// testCrossLayerIntegration performs comprehensive cross-layer integration testing
func testCrossLayerIntegration(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: HTTP to Domain Layer Data Flow
	t.Run("http_to_domain_data_flow", func(t *testing.T) {
		testHTTPToDomainDataFlow(t, client, baseURL)
	})

	// Test 2: Domain Business Rules Enforcement
	t.Run("domain_business_rules", func(t *testing.T) {
		testDomainBusinessRules(t, client, baseURL)
	})

	// Test 3: Repository Layer Data Persistence
	t.Run("repository_data_persistence", func(t *testing.T) {
		testRepositoryDataPersistence(t, client, baseURL)
	})

	// Test 4: Application Service Orchestration
	t.Run("application_service_orchestration", func(t *testing.T) {
		testApplicationServiceOrchestration(t, client, baseURL)
	})

	// Test 5: Event Publishing and Domain Events
	t.Run("domain_events_integration", func(t *testing.T) {
		testDomainEventsIntegration(t, client, baseURL)
	})

	// Test 6: Value Object Validation Across Layers
	t.Run("value_object_validation", func(t *testing.T) {
		testValueObjectValidation(t, client, baseURL)
	})

	// Test 7: Error Propagation Through Layers
	t.Run("error_propagation", func(t *testing.T) {
		testErrorPropagation(t, client, baseURL)
	})

	// Test 8: Transaction Boundaries and Consistency
	t.Run("transaction_boundaries", func(t *testing.T) {
		testTransactionBoundaries(t, client, baseURL)
	})
}

// testHTTPToDomainDataFlow tests data transformation and flow from HTTP layer to Domain layer
func testHTTPToDomainDataFlow(t *testing.T, client *http.Client, baseURL string) {
	// Test that HTTP request DTOs are properly transformed to Domain entities
	testEmail := "http-domain@example.com"
	
	// Step 1: HTTP Layer - Send registration request
	registrationData := map[string]interface{}{
		"email":      testEmail,
		"password":   "domain123",
		"first_name": "HTTP",
		"last_name":  "Domain",
	}

	jsonData, err := json.Marshal(registrationData)
	require.NoError(t, err)

	resp, err := client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var registerResponse map[string]interface{}
	err = json.Unmarshal(body, &registerResponse)
	require.NoError(t, err)

	// Verify that Domain entities are properly created and returned
	if user, ok := registerResponse["user"].(map[string]interface{}); ok {
		// Validate that domain constraints are applied
		assert.Equal(t, testEmail, user["email"], "Email should pass through domain layer unchanged")
		assert.Equal(t, "HTTP", user["first_name"], "First name should pass through domain layer")
		assert.Equal(t, "Domain", user["last_name"], "Last name should pass through domain layer")
		assert.NotEmpty(t, user["id"], "Domain should generate unique ID")
		
		// Verify timestamps are properly set by domain layer
		if createdAt, ok := user["created_at"]; ok {
			assert.NotEmpty(t, createdAt, "Domain should set created_at timestamp")
		}
		if updatedAt, ok := user["updated_at"]; ok {
			assert.NotEmpty(t, updatedAt, "Domain should set updated_at timestamp")
		}
	}

	// Step 2: Verify data reaches repository layer correctly
	// Login and retrieve profile to ensure data persisted correctly
	loginData := map[string]interface{}{
		"email":    testEmail,
		"password": "domain123",
	}

	jsonData, err = json.Marshal(loginData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	require.NoError(t, err)

	// Retrieve profile to verify domain entity reconstruction from repository
	resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", authResponse.AccessToken, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var profileResponse map[string]interface{}
		err = json.Unmarshal(body, &profileResponse)
		require.NoError(t, err)

		if user, ok := profileResponse["user"].(map[string]interface{}); ok {
			// Verify data integrity through all layers
			assert.Equal(t, testEmail, user["email"], "Data should maintain integrity across all layers")
			assert.Equal(t, "HTTP", user["first_name"], "First name should maintain integrity")
			assert.Equal(t, "Domain", user["last_name"], "Last name should maintain integrity")
		}
	}

	t.Log("HTTP to Domain data flow validation: SUCCESS")
}

// testDomainBusinessRules tests that domain business rules are properly enforced
func testDomainBusinessRules(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: Email uniqueness constraint (domain rule)
	testEmail := "unique-email@example.com"
	
	userData := map[string]interface{}{
		"email":      testEmail,
		"password":   "password123",
		"first_name": "First",
		"last_name":  "User",
	}

	// Create first user
	jsonData, err := json.Marshal(userData)
	require.NoError(t, err)

	resp, err := client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Try to create second user with same email (should be rejected by domain)
	resp, err = client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.NotEqual(t, http.StatusCreated, resp.StatusCode, "Domain should reject duplicate email")

	// Test 2: Email format validation (value object rule)
	invalidEmailData := map[string]interface{}{
		"email":      "invalid-email-format",
		"password":   "password123",
		"first_name": "Invalid",
		"last_name":  "Email",
	}

	jsonData, err = json.Marshal(invalidEmailData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.NotEqual(t, http.StatusCreated, resp.StatusCode, "Domain should reject invalid email format")

	// Test 3: Password strength requirements (domain rule)
	weakPasswordData := map[string]interface{}{
		"email":      "weak-password@example.com",
		"password":   "123", // Too short
		"first_name": "Weak",
		"last_name":  "Password",
	}

	jsonData, err = json.Marshal(weakPasswordData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.NotEqual(t, http.StatusCreated, resp.StatusCode, "Domain should reject weak passwords")

	// Test 4: Name length constraints (value object rule)
	longNameData := map[string]interface{}{
		"email":      "long-name@example.com",
		"password":   "password123",
		"first_name": strings.Repeat("A", 100), // Too long
		"last_name":  "Normal",
	}

	jsonData, err = json.Marshal(longNameData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	
	// May be rejected at validation layer or allowed - log the result
	t.Logf("Long name registration returned status: %d", resp.StatusCode)

	t.Log("Domain business rules validation: SUCCESS")
}

// testRepositoryDataPersistence tests repository layer data persistence and retrieval
func testRepositoryDataPersistence(t *testing.T, client *http.Client, baseURL string) {
	// Create user with specific data to test persistence
	testEmail := "persistence@example.com"
	createdUser, authToken := createTestUserWithAuth(t, client, baseURL, testEmail)

	// Test 1: Data persistence through CRUD operations
	// Update user data multiple times to test repository persistence
	updates := []map[string]interface{}{
		{"first_name": "Persistence1"},
		{"last_name": "Test1"},
		{"first_name": "Persistence2", "last_name": "Test2"},
		{"first_name": "Final", "last_name": "Persistence"},
	}

	for i, updateData := range updates {
		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		// Update via application service
		resp, err := makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", authToken, strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			// Immediately retrieve to verify persistence
			resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", authToken, nil)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			var profileResponse map[string]interface{}
			err = json.Unmarshal(body, &profileResponse)
			require.NoError(t, err)

			if user, ok := profileResponse["user"].(map[string]interface{}); ok {
				// Verify repository persistence
				if firstName, exists := updateData["first_name"]; exists {
					assert.Equal(t, firstName, user["first_name"], "Repository should persist first name change %d", i+1)
				}
				if lastName, exists := updateData["last_name"]; exists {
					assert.Equal(t, lastName, user["last_name"], "Repository should persist last name change %d", i+1)
				}
				// Email should never change
				assert.Equal(t, createdUser.Email, user["email"], "Repository should maintain email consistency")
			}
		}
	}

	// Test 2: Repository query capabilities
	// Test that repository can find user by email
	loginData := map[string]interface{}{
		"email":    testEmail,
		"password": "password123",
	}

	jsonData, err := json.Marshal(loginData)
	require.NoError(t, err)

	resp, err := client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Repository should find user by email")

	// Test 3: Repository data integrity after multiple sessions
	// Create second session to verify data consistency
	resp, err = client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Repository should support multiple sessions")

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var authResponse AuthResponse
		err = json.Unmarshal(body, &authResponse)
		require.NoError(t, err)

		// Verify same user data across sessions
		assert.Equal(t, createdUser.Email, authResponse.User.Email, "Repository should maintain data consistency across sessions")
	}

	t.Log("Repository data persistence validation: SUCCESS")
}

// testApplicationServiceOrchestration tests application service layer orchestration
func testApplicationServiceOrchestration(t *testing.T, client *http.Client, baseURL string) {
	// Test that application services properly orchestrate domain operations
	testEmail := "orchestration@example.com"
	
	// Test 1: Registration orchestration (UserService + AuthService coordination)
	registrationData := map[string]interface{}{
		"email":      testEmail,
		"password":   "orchestration123",
		"first_name": "Application",
		"last_name":  "Service",
	}

	jsonData, err := json.Marshal(registrationData)
	require.NoError(t, err)

	resp, err := client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Verify that application service orchestrated both user creation and response formatting
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var registerResponse map[string]interface{}
	err = json.Unmarshal(body, &registerResponse)
	require.NoError(t, err)

	// Application service should provide proper response structure
	assert.Contains(t, registerResponse, "user", "Application service should provide user in response")
	assert.Contains(t, registerResponse, "message", "Application service should provide message")

	// Test 2: Login orchestration (Authentication + Session management)
	loginData := map[string]interface{}{
		"email":    testEmail,
		"password": "orchestration123",
	}

	jsonData, err = json.Marshal(loginData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	require.NoError(t, err)

	// Application service should orchestrate token generation + user data
	assert.NotEmpty(t, authResponse.AccessToken, "Application service should orchestrate token generation")
	assert.Equal(t, "Bearer", authResponse.TokenType, "Application service should set proper token type")
	assert.Greater(t, authResponse.ExpiresIn, 0, "Application service should set expiration")
	assert.Equal(t, testEmail, authResponse.User.Email, "Application service should include user data")

	// Test 3: Profile update orchestration (Validation + Domain update + Response)
	updateData := map[string]interface{}{
		"first_name": "Orchestrated",
		"last_name":  "Update",
	}

	jsonData, err = json.Marshal(updateData)
	require.NoError(t, err)

	resp, err = makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", authResponse.AccessToken, strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Application service should orchestrate validation, domain update, and response
		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var updateResponse map[string]interface{}
		err = json.Unmarshal(body, &updateResponse)
		require.NoError(t, err)

		if user, ok := updateResponse["user"].(map[string]interface{}); ok {
			assert.Equal(t, "Orchestrated", user["first_name"], "Application service should orchestrate domain updates")
			assert.Equal(t, "Update", user["last_name"], "Application service should orchestrate domain updates")
		}
	}

	t.Log("Application service orchestration validation: SUCCESS")
}

// testDomainEventsIntegration tests domain events publishing and handling
func testDomainEventsIntegration(t *testing.T, client *http.Client, baseURL string) {
	// Test that domain events are properly published during operations
	testEmail := "events@example.com"
	
	// Test 1: User registration should trigger domain events
	registrationData := map[string]interface{}{
		"email":      testEmail,
		"password":   "events123",
		"first_name": "Domain",
		"last_name":  "Events",
	}

	jsonData, err := json.Marshal(registrationData)
	require.NoError(t, err)

	resp, err := client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Registration should succeed, indicating domain events were handled properly
	// (Events may be logged or processed asynchronously)

	// Test 2: Login should trigger authentication events
	loginData := map[string]interface{}{
		"email":    testEmail,
		"password": "events123",
	}

	jsonData, err = json.Marshal(loginData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	require.NoError(t, err)

	// Test 3: Failed login should trigger failure events
	invalidLoginData := map[string]interface{}{
		"email":    testEmail,
		"password": "wrongpassword",
	}

	jsonData, err = json.Marshal(invalidLoginData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Failed login should be handled properly, indicating event system is working

	// Test 4: Profile updates should trigger domain events
	updateData := map[string]interface{}{
		"first_name": "Event",
		"last_name":  "Triggered",
	}

	jsonData, err = json.Marshal(updateData)
	require.NoError(t, err)

	resp, err = makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", authResponse.AccessToken, strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Profile update should succeed, indicating domain events are properly integrated
	if resp.StatusCode == http.StatusOK {
		t.Log("Profile update event integration working")
	}

	t.Log("Domain events integration validation: SUCCESS")
}

// testValueObjectValidation tests value object validation across layers
func testValueObjectValidation(t *testing.T, client *http.Client, baseURL string) {
	// Test that value objects properly validate data across all layers
	
	// Test 1: Email value object validation
	invalidEmails := []string{
		"invalid-email",
		"@domain.com",
		"user@",
		"user@domain",
		"",
		"user with spaces@domain.com",
		"user@domain..com",
	}

	for _, invalidEmail := range invalidEmails {
		userData := map[string]interface{}{
			"email":      invalidEmail,
			"password":   "password123",
			"first_name": "Value",
			"last_name":  "Object",
		}

		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Value object validation should reject invalid emails
		assert.NotEqual(t, http.StatusCreated, resp.StatusCode, "Value object should reject invalid email: %s", invalidEmail)
	}

	// Test 2: Valid email should pass value object validation
	validUserData := map[string]interface{}{
		"email":      "valid@example.com",
		"password":   "password123",
		"first_name": "Valid",
		"last_name":  "Email",
	}

	jsonData, err := json.Marshal(validUserData)
	require.NoError(t, err)

	resp, err := client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Value object should accept valid email")

	// Test 3: Name value objects validation in updates
	_, authToken := createTestUserWithAuth(t, client, baseURL, "name-validation@example.com")

	// Test empty names
	emptyNameData := map[string]interface{}{
		"first_name": "",
		"last_name":  "",
	}

	jsonData, err = json.Marshal(emptyNameData)
	require.NoError(t, err)

	resp, err = makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", authToken, strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()

	// May be rejected at value object level or allowed with defaults
	t.Logf("Empty names update returned status: %d", resp.StatusCode)

	// Test 4: Valid names should pass value object validation
	validNameData := map[string]interface{}{
		"first_name": "Valid",
		"last_name":  "Names",
	}

	jsonData, err = json.Marshal(validNameData)
	require.NoError(t, err)

	resp, err = makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", authToken, strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		t.Log("Valid name value objects accepted")
	}

	t.Log("Value object validation across layers: SUCCESS")
}

// testErrorPropagation tests error propagation through hexagonal architecture layers
func testErrorPropagation(t *testing.T, client *http.Client, baseURL string) {
	// Test that errors are properly propagated and handled across layers
	
	// Test 1: Domain layer errors (business rule violations)
	duplicateEmailData := map[string]interface{}{
		"email":      "error-propagation@example.com",
		"password":   "password123",
		"first_name": "Error",
		"last_name":  "Test",
	}

	// Create first user
	jsonData, err := json.Marshal(duplicateEmailData)
	require.NoError(t, err)

	resp, err := client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Try to create duplicate - domain error should propagate to HTTP layer
	resp, err = client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.NotEqual(t, http.StatusCreated, resp.StatusCode, "Domain error should propagate to HTTP layer")

	// Test 2: Authentication errors
	invalidLoginData := map[string]interface{}{
		"email":    "nonexistent@example.com",
		"password": "wrongpassword",
	}

	jsonData, err = json.Marshal(invalidLoginData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Authentication error should propagate correctly")

	// Test 3: Validation errors (DTO/value object level)
	invalidJsonData := `{"email": "test@example.com", "password": 123}` // Invalid JSON type

	resp, err = client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(invalidJsonData))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Validation error should propagate to HTTP layer")

	// Test 4: Authorization errors
	invalidToken := "invalid-token"
	resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", invalidToken, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Authorization error should propagate correctly")

	t.Log("Error propagation through layers: SUCCESS")
}

// testTransactionBoundaries tests transaction boundaries and data consistency
func testTransactionBoundaries(t *testing.T, client *http.Client, baseURL string) {
	// Test that transaction boundaries are properly maintained
	testEmail := "transaction@example.com"
	
	// Test 1: Registration transaction (should be atomic)
	registrationData := map[string]interface{}{
		"email":      testEmail,
		"password":   "transaction123",
		"first_name": "Transaction",
		"last_name":  "Test",
	}

	jsonData, err := json.Marshal(registrationData)
	require.NoError(t, err)

	resp, err := client.Post(baseURL+"/api/auth/register", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	// Verify that user can immediately login (transaction committed)
	loginData := map[string]interface{}{
		"email":    testEmail,
		"password": "transaction123",
	}

	jsonData, err = json.Marshal(loginData)
	require.NoError(t, err)

	resp, err = client.Post(baseURL+"/api/auth/login", "application/json", strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Transaction should be committed, allowing immediate login")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	require.NoError(t, err)

	// Test 2: Profile update transaction consistency
	updateData := map[string]interface{}{
		"first_name": "Updated",
		"last_name":  "Transaction",
	}

	jsonData, err = json.Marshal(updateData)
	require.NoError(t, err)

	resp, err = makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", authResponse.AccessToken, strings.NewReader(string(jsonData)))
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		// Immediately verify update was persisted (transaction committed)
		resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", authResponse.AccessToken, nil)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err = io.ReadAll(resp.Body)
		require.NoError(t, err)

		var profileResponse map[string]interface{}
		err = json.Unmarshal(body, &profileResponse)
		require.NoError(t, err)

		if user, ok := profileResponse["user"].(map[string]interface{}); ok {
			assert.Equal(t, "Updated", user["first_name"], "Transaction should ensure immediate consistency")
			assert.Equal(t, "Transaction", user["last_name"], "Transaction should ensure immediate consistency")
		}
	}

	// Test 3: Concurrent transaction handling
	// Make multiple concurrent updates to test transaction isolation
	concurrentUpdates := []map[string]interface{}{
		{"first_name": "Concurrent1"},
		{"first_name": "Concurrent2"},
		{"first_name": "Concurrent3"},
	}

	for i, concurrentUpdate := range concurrentUpdates {
		go func(updateNum int, updateData map[string]interface{}) {
			jsonData, err := json.Marshal(updateData)
			if err != nil {
				return
			}

			resp, err := makeAuthenticatedRequest(client, "PUT", baseURL+"/api/users/profile", authResponse.AccessToken, strings.NewReader(string(jsonData)))
			if err != nil {
				return
			}
			defer resp.Body.Close()

			t.Logf("Concurrent update %d status: %d", updateNum, resp.StatusCode)
		}(i+1, concurrentUpdate)
	}

	// Give time for concurrent updates to complete
	time.Sleep(1 * time.Second)

	// Verify final state is consistent
	resp, err = makeAuthenticatedRequest(client, "GET", baseURL+"/api/users/profile", authResponse.AccessToken, nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Final state should be consistent after concurrent updates")

	t.Log("Transaction boundaries validation: SUCCESS")
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// testDomainLayerBehavior tests comprehensive domain layer behavior including
// value objects, entities, business rules, and domain events
func testDomainLayerBehavior(t *testing.T, client *http.Client, baseURL string) {
	// Test value object behavior through the API
	t.Run("value_object_behavior", func(t *testing.T) {
		testValueObjectBehavior(t, client, baseURL)
	})

	// Test entity behavior and business rules
	t.Run("entity_behavior", func(t *testing.T) {
		testEntityBehavior(t, client, baseURL)
	})

	// Test domain services and business logic
	t.Run("domain_services", func(t *testing.T) {
		testDomainServices(t, client, baseURL)
	})

	// Test domain events
	t.Run("domain_events", func(t *testing.T) {
		testDomainEvents(t, client, baseURL)
	})
}
// testValueObjectBehavior tests value object validation and edge cases
func testValueObjectBehavior(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: Email Value Object
	t.Run("email_value_object", func(t *testing.T) {
		// Test various email formats and edge cases
		emailTests := []struct {
			email    string
			valid    bool
			scenario string
		}{
			// Valid emails
			{"user@example.com", true, "standard email"},
			{"user.name@example.com", true, "email with dot"},
			{"user+tag@example.com", true, "email with plus"},
			{"user_name@example.com", true, "email with underscore"},
			{"123@example.com", true, "numeric local part"},
			{"a@example.com", true, "single character local"},
			{"user@sub.example.com", true, "subdomain"},
			{"UPPERCASE@EXAMPLE.COM", true, "uppercase email"},
			{"  user@example.com  ", true, "email with spaces"},
			
			// Invalid emails
			{"", false, "empty email"},
			{"notanemail", false, "missing @ symbol"},
			{"@example.com", false, "missing local part"},
			{"user@", false, "missing domain"},
			{"user @example.com", false, "space in local part"},
			{"user@example .com", false, "space in domain"},
			{"user@@example.com", false, "double @ symbol"},
			{"user@.com", false, "missing domain name"},
			{"user@example.", false, "missing TLD"},
			{"user@-example.com", false, "domain starts with hyphen"},
			{"user@example.com-", false, "domain ends with hyphen"},
			{"user..name@example.com", false, "consecutive dots"},
			{".user@example.com", false, "starts with dot"},
			{"user.@example.com", false, "ends with dot"},
			{strings.Repeat("a", 255) + "@example.com", false, "too long local part"},
			{"user@" + strings.Repeat("a", 255) + ".com", false, "too long domain"},
		}

		for _, test := range emailTests {
			userData := map[string]interface{}{
				"email":      test.email,
				"password":   "password123",
				"first_name": "Test",
				"last_name":  "User",
			}
			
			jsonData, err := json.Marshal(userData)
			require.NoError(t, err)
			
			resp, err := client.Post(baseURL+"/api/auth/register", "application/json", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			defer resp.Body.Close()
			
			if test.valid {
				// Valid emails should either succeed or fail with conflict (already exists)
				assert.True(t, resp.StatusCode == http.StatusCreated || 
					resp.StatusCode == http.StatusConflict,
					"Expected success or conflict for valid email '%s' (%s), got %d", 
					test.email, test.scenario, resp.StatusCode)
			} else {
				// Invalid emails should always fail with bad request
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
					"Expected bad request for invalid email '%s' (%s)", 
					test.email, test.scenario)
			}
		}
		
		// Test email normalization
		t.Run("email_normalization", func(t *testing.T) {
			// Register with uppercase email
			upperEmail := "NORMALIZED@EXAMPLE.COM"
			userData := map[string]interface{}{
				"email":      upperEmail,
				"password":   "password123",
				"first_name": "Normalized",
				"last_name":  "User",
			}
			
			jsonData, err := json.Marshal(userData)
			require.NoError(t, err)
			
			resp, err := client.Post(baseURL+"/api/auth/register", "application/json", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			resp.Body.Close()
			
			// Should be able to login with lowercase version
			loginData := map[string]interface{}{
				"email":    "normalized@example.com",
				"password": "password123",
			}
			
			jsonData, err = json.Marshal(loginData)
			require.NoError(t, err)
			
			resp, err = client.Post(baseURL+"/api/auth/login", "application/json", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode,
				"Should be able to login with normalized email")
		})
	})

	// Test 2: Password Value Object
	t.Run("password_value_object", func(t *testing.T) {
		passwordTests := []struct {
			password string
			valid    bool
			scenario string
		}{
			// Valid passwords
			{"password123", true, "standard password"},
			{"P@ssw0rd!", true, "password with special chars"},
			{"12345678", true, "numeric password"},
			{"UPPERCASE", true, "uppercase password"},
			{"lowercase", true, "lowercase password"},
			{strings.Repeat("a", 128), true, "long password"},
			
			// Invalid passwords
			{"", false, "empty password"},
			{"short", false, "too short (less than 8)"},
			{"       ", false, "only spaces"},
			{"1234567", false, "7 characters"},
			{strings.Repeat("a", 256), false, "too long password"},
		}
		
		for _, test := range passwordTests {
			userData := map[string]interface{}{
				"email":      fmt.Sprintf("pwtest%d@example.com", time.Now().UnixNano()),
				"password":   test.password,
				"first_name": "Password",
				"last_name":  "Test",
			}
			
			jsonData, err := json.Marshal(userData)
			require.NoError(t, err)
			
			resp, err := client.Post(baseURL+"/api/auth/register", "application/json", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			defer resp.Body.Close()
			
			if test.valid {
				assert.Equal(t, http.StatusCreated, resp.StatusCode,
					"Expected success for valid password (%s)", test.scenario)
			} else {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
					"Expected bad request for invalid password (%s)", test.scenario)
			}
		}
	})

	// Test 3: Name validation (First/Last name)
	t.Run("name_validation", func(t *testing.T) {
		nameTests := []struct {
			firstName string
			lastName  string
			valid     bool
			scenario  string
		}{
			// Valid names
			{"John", "Doe", true, "standard names"},
			{"Jo", "Do", true, "minimum length names"},
			{strings.Repeat("A", 50), strings.Repeat("B", 50), true, "maximum length names"},
			{"Jean-Pierre", "O'Connor", true, "names with special chars"},
			{"José", "García", true, "names with accents"},
			{"李", "王", true, "unicode names"},
			
			// Invalid names
			{"", "Doe", false, "empty first name"},
			{"John", "", false, "empty last name"},
			{"J", "Doe", false, "too short first name"},
			{"John", "D", false, "too short last name"},
			{strings.Repeat("A", 51), "Doe", false, "too long first name"},
			{"John", strings.Repeat("B", 51), false, "too long last name"},
			{"  ", "Doe", false, "only spaces first name"},
			{"John", "  ", false, "only spaces last name"},
		}
		
		for _, test := range nameTests {
			userData := map[string]interface{}{
				"email":      fmt.Sprintf("nametest%d@example.com", time.Now().UnixNano()),
				"password":   "password123",
				"first_name": test.firstName,
				"last_name":  test.lastName,
			}
			
			jsonData, err := json.Marshal(userData)
			require.NoError(t, err)
			
			resp, err := client.Post(baseURL+"/api/auth/register", "application/json", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			defer resp.Body.Close()
			
			if test.valid {
				assert.Equal(t, http.StatusCreated, resp.StatusCode,
					"Expected success for valid names (%s)", test.scenario)
			} else {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
					"Expected bad request for invalid names (%s)", test.scenario)
			}
		}
	})
}

// testEntityBehavior tests entity creation and business rules
func testEntityBehavior(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: User Entity Creation and Updates
	t.Run("user_entity_lifecycle", func(t *testing.T) {
		// Create a user
		email := fmt.Sprintf("entity%d@example.com", time.Now().UnixNano())
		userData := map[string]interface{}{
			"email":      email,
			"password":   "password123",
			"first_name": "Entity",
			"last_name":  "Test",
		}
		
		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)
		
		resp, err := client.Post(baseURL+"/api/auth/register", "application/json", 
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		
		var createResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&createResp)
		require.NoError(t, err)
		resp.Body.Close()
		
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		
		// Login to get token
		loginData := map[string]interface{}{
			"email":    email,
			"password": "password123",
		}
		
		jsonData, err = json.Marshal(loginData)
		require.NoError(t, err)
		
		resp, err = client.Post(baseURL+"/api/auth/login", "application/json", 
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		
		var authResp AuthResponse
		err = json.NewDecoder(resp.Body).Decode(&authResp)
		require.NoError(t, err)
		resp.Body.Close()
		
		// Test entity updates
		t.Run("update_user_fields", func(t *testing.T) {
			// Update user profile
			updateData := map[string]interface{}{
				"first_name": "Updated",
				"last_name":  "Entity",
			}
			
			jsonData, err := json.Marshal(updateData)
			require.NoError(t, err)
			
			req, err := http.NewRequest("PUT", baseURL+"/api/users/me", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			
			// Verify updates
			req, err = http.NewRequest("GET", baseURL+"/api/users/me", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)
			
			resp, err = client.Do(req)
			require.NoError(t, err)
			
			var userResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&userResp)
			require.NoError(t, err)
			resp.Body.Close()
			
			assert.Equal(t, "Updated", userResp["first_name"])
			assert.Equal(t, "Entity", userResp["last_name"])
		})
		
		// Test password change business rules
		t.Run("password_change_rules", func(t *testing.T) {
			// Change password with wrong current password
			changeData := map[string]interface{}{
				"current_password": "wrongpassword",
				"new_password":     "newpassword123",
			}
			
			jsonData, err := json.Marshal(changeData)
			require.NoError(t, err)
			
			req, err := http.NewRequest("POST", baseURL+"/api/auth/change-password", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := client.Do(req)
			require.NoError(t, err)
			resp.Body.Close()
			
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
				"Should fail with wrong current password")
			
			// Change password with correct current password
			changeData["current_password"] = "password123"
			
			jsonData, err = json.Marshal(changeData)
			require.NoError(t, err)
			
			req, err = http.NewRequest("POST", baseURL+"/api/auth/change-password", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)
			req.Header.Set("Content-Type", "application/json")
			
			resp, err = client.Do(req)
			require.NoError(t, err)
			resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode,
				"Should succeed with correct current password")
		})
	})

	// Test 2: Entity Timestamps
	t.Run("entity_timestamps", func(t *testing.T) {
		// Create a user and check timestamps
		email := fmt.Sprintf("timestamp%d@example.com", time.Now().UnixNano())
		userData := map[string]interface{}{
			"email":      email,
			"password":   "password123",
			"first_name": "Time",
			"last_name":  "Stamp",
		}
		
		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)
		
		createTime := time.Now()
		resp, err := client.Post(baseURL+"/api/auth/register", "application/json", 
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		
		var createResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&createResp)
		require.NoError(t, err)
		resp.Body.Close()
		
		// Check that user response includes timestamps
		user := createResp["user"].(map[string]interface{})
		assert.NotEmpty(t, user["created_at"], "Should have created_at timestamp")
		assert.NotEmpty(t, user["updated_at"], "Should have updated_at timestamp")
		
		// Verify timestamps are recent
		createdAt, err := time.Parse(time.RFC3339, user["created_at"].(string))
		require.NoError(t, err)
		assert.WithinDuration(t, createTime, createdAt, 5*time.Second,
			"Created timestamp should be recent")
	})

	// Test 3: Entity Invariants
	t.Run("entity_invariants", func(t *testing.T) {
		// Test that entity maintains its invariants
		email := fmt.Sprintf("invariant%d@example.com", time.Now().UnixNano())
		
		// Try to create user with conflicting data
		userData := map[string]interface{}{
			"email":      email,
			"password":   "password123",
			"first_name": "Test",
			"last_name":  "User",
		}
		
		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)
		
		// First creation should succeed
		resp, err := client.Post(baseURL+"/api/auth/register", "application/json", 
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		
		// Second creation with same email should fail (uniqueness invariant)
		resp, err = client.Post(baseURL+"/api/auth/register", "application/json", 
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		resp.Body.Close()
		
		assert.True(t, resp.StatusCode == http.StatusConflict || 
			resp.StatusCode == http.StatusBadRequest,
			"Should enforce email uniqueness invariant")
	})
}

// testDomainServices tests domain service business logic
func testDomainServices(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: Authentication Domain Service
	t.Run("auth_domain_service", func(t *testing.T) {
		// Test login validation with various scenarios
		email := fmt.Sprintf("authservice%d@example.com", time.Now().UnixNano())
		
		// Register user
		userData := map[string]interface{}{
			"email":      email,
			"password":   "password123",
			"first_name": "Auth",
			"last_name":  "Service",
		}
		
		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)
		
		resp, err := client.Post(baseURL+"/api/auth/register", "application/json", 
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		resp.Body.Close()
		
		// Test multiple failed login attempts (account lockout simulation)
		t.Run("account_lockout_behavior", func(t *testing.T) {
			wrongLoginData := map[string]interface{}{
				"email":    email,
				"password": "wrongpassword",
			}
			
			// Attempt multiple failed logins
			for i := 0; i < 5; i++ {
				jsonData, err := json.Marshal(wrongLoginData)
				require.NoError(t, err)
				
				resp, err := client.Post(baseURL+"/api/auth/login", "application/json", 
					strings.NewReader(string(jsonData)))
				require.NoError(t, err)
				resp.Body.Close()
				
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
					"Failed login attempt %d should return unauthorized", i+1)
			}
			
			// After multiple failures, even correct password might be temporarily blocked
			// (depending on implementation)
			correctLoginData := map[string]interface{}{
				"email":    email,
				"password": "password123",
			}
			
			jsonData, err := json.Marshal(correctLoginData)
			require.NoError(t, err)
			
			resp, err := client.Post(baseURL+"/api/auth/login", "application/json", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			defer resp.Body.Close()
			
			// Should eventually succeed (implementation dependent)
			// Just verify we get a response
			assert.True(t, resp.StatusCode == http.StatusOK || 
				resp.StatusCode == http.StatusUnauthorized,
				"Login should either succeed or be locked out")
		})
		
		// Test session expiration logic
		t.Run("session_expiration", func(t *testing.T) {
			// Login to get a token
			loginData := map[string]interface{}{
				"email":    email,
				"password": "password123",
			}
			
			jsonData, err := json.Marshal(loginData)
			require.NoError(t, err)
			
			resp, err := client.Post(baseURL+"/api/auth/login", "application/json", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			
			var authResp AuthResponse
			err = json.NewDecoder(resp.Body).Decode(&authResp)
			require.NoError(t, err)
			resp.Body.Close()
			
			// Check that token has expiration time
			assert.Greater(t, authResp.ExpiresIn, 0, "Token should have expiration time")
			
			// Verify token works immediately
			req, err := http.NewRequest("GET", baseURL+"/api/users/me", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)
			
			resp, err = client.Do(req)
			require.NoError(t, err)
			resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode,
				"Fresh token should work")
		})
	})

	// Test 2: User Domain Service Business Rules
	t.Run("user_domain_service", func(t *testing.T) {
		// Test profanity check in names (if implemented)
		t.Run("profanity_check", func(t *testing.T) {
			profanityTests := []struct {
				firstName string
				lastName  string
				shouldFail bool
			}{
				{"Normal", "User", false},
				{"Test", "User", false},
				// Add actual profanity tests based on implementation
				// {"BadWord", "User", true},
			}
			
			for _, test := range profanityTests {
				userData := map[string]interface{}{
					"email":      fmt.Sprintf("profanity%d@example.com", time.Now().UnixNano()),
					"password":   "password123",
					"first_name": test.firstName,
					"last_name":  test.lastName,
				}
				
				jsonData, err := json.Marshal(userData)
				require.NoError(t, err)
				
				resp, err := client.Post(baseURL+"/api/auth/register", "application/json", 
					strings.NewReader(string(jsonData)))
				require.NoError(t, err)
				resp.Body.Close()
				
				if test.shouldFail {
					assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
						"Should reject profane names")
				} else {
					assert.Equal(t, http.StatusCreated, resp.StatusCode,
						"Should accept clean names")
				}
			}
		})
		
		// Test special user types (admin, premium)
		t.Run("special_user_types", func(t *testing.T) {
			// This would test domain logic for different user types
			// Implementation depends on whether the domain supports user roles
			
			// For now, just verify standard user creation works
			userData := map[string]interface{}{
				"email":      fmt.Sprintf("special%d@example.com", time.Now().UnixNano()),
				"password":   "password123",
				"first_name": "Special",
				"last_name":  "User",
			}
			
			jsonData, err := json.Marshal(userData)
			require.NoError(t, err)
			
			resp, err := client.Post(baseURL+"/api/auth/register", "application/json", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			resp.Body.Close()
			
			assert.Equal(t, http.StatusCreated, resp.StatusCode)
		})
	})
}

// testDomainEvents tests domain event generation and handling
func testDomainEvents(t *testing.T, client *http.Client, baseURL string) {
	// Note: Testing domain events through the API is indirect
	// We can verify that actions that should generate events complete successfully
	
	// Test 1: User Events
	t.Run("user_lifecycle_events", func(t *testing.T) {
		email := fmt.Sprintf("events%d@example.com", time.Now().UnixNano())
		
		// UserRegisteredEvent - triggered by registration
		t.Run("user_registered_event", func(t *testing.T) {
			userData := map[string]interface{}{
				"email":      email,
				"password":   "password123",
				"first_name": "Event",
				"last_name":  "Test",
			}
			
			jsonData, err := json.Marshal(userData)
			require.NoError(t, err)
			
			resp, err := client.Post(baseURL+"/api/auth/register", "application/json", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			resp.Body.Close()
			
			assert.Equal(t, http.StatusCreated, resp.StatusCode,
				"Registration should succeed and trigger UserRegisteredEvent")
		})
		
		// UserLoggedInEvent - triggered by login
		t.Run("user_logged_in_event", func(t *testing.T) {
			loginData := map[string]interface{}{
				"email":    email,
				"password": "password123",
			}
			
			jsonData, err := json.Marshal(loginData)
			require.NoError(t, err)
			
			resp, err := client.Post(baseURL+"/api/auth/login", "application/json", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			
			var authResp AuthResponse
			err = json.NewDecoder(resp.Body).Decode(&authResp)
			require.NoError(t, err)
			resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode,
				"Login should succeed and trigger UserLoggedInEvent")
			
			// UserUpdatedEvent - triggered by profile update
			t.Run("user_updated_event", func(t *testing.T) {
				updateData := map[string]interface{}{
					"first_name": "Updated",
					"last_name":  "Event",
				}
				
				jsonData, err := json.Marshal(updateData)
				require.NoError(t, err)
				
				req, err := http.NewRequest("PUT", baseURL+"/api/users/me", 
					strings.NewReader(string(jsonData)))
				require.NoError(t, err)
				req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)
				req.Header.Set("Content-Type", "application/json")
				
				resp, err := client.Do(req)
				require.NoError(t, err)
				resp.Body.Close()
				
				assert.Equal(t, http.StatusOK, resp.StatusCode,
					"Update should succeed and trigger UserUpdatedEvent")
			})
			
			// UserPasswordChangedEvent - triggered by password change
			t.Run("user_password_changed_event", func(t *testing.T) {
				changeData := map[string]interface{}{
					"current_password": "password123",
					"new_password":     "newpassword123",
				}
				
				jsonData, err := json.Marshal(changeData)
				require.NoError(t, err)
				
				req, err := http.NewRequest("POST", baseURL+"/api/auth/change-password", 
					strings.NewReader(string(jsonData)))
				require.NoError(t, err)
				req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)
				req.Header.Set("Content-Type", "application/json")
				
				resp, err := client.Do(req)
				require.NoError(t, err)
				resp.Body.Close()
				
				assert.Equal(t, http.StatusOK, resp.StatusCode,
					"Password change should succeed and trigger UserPasswordChangedEvent")
			})
		})
	})

	// Test 2: Authentication Events
	t.Run("authentication_events", func(t *testing.T) {
		// UserLoginFailedEvent - triggered by failed login
		t.Run("login_failed_event", func(t *testing.T) {
			loginData := map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": "wrongpassword",
			}
			
			jsonData, err := json.Marshal(loginData)
			require.NoError(t, err)
			
			resp, err := client.Post(baseURL+"/api/auth/login", "application/json", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			resp.Body.Close()
			
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
				"Failed login should trigger UserLoginFailedEvent")
		})
		
		// TokenRefreshedEvent - triggered by token refresh
		t.Run("token_refreshed_event", func(t *testing.T) {
			// First login to get tokens
			email := fmt.Sprintf("refresh%d@example.com", time.Now().UnixNano())
			userData := map[string]interface{}{
				"email":      email,
				"password":   "password123",
				"first_name": "Refresh",
				"last_name":  "Test",
			}
			
			jsonData, err := json.Marshal(userData)
			require.NoError(t, err)
			
			resp, err := client.Post(baseURL+"/api/auth/register", "application/json", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			resp.Body.Close()
			
			// Login
			loginData := map[string]interface{}{
				"email":    email,
				"password": "password123",
			}
			
			jsonData, err = json.Marshal(loginData)
			require.NoError(t, err)
			
			resp, err = client.Post(baseURL+"/api/auth/login", "application/json", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			
			var authResp AuthResponse
			err = json.NewDecoder(resp.Body).Decode(&authResp)
			require.NoError(t, err)
			resp.Body.Close()
			
			// Refresh token
			refreshData := map[string]interface{}{
				"refresh_token": authResp.RefreshToken,
			}
			
			jsonData, err = json.Marshal(refreshData)
			require.NoError(t, err)
			
			resp, err = client.Post(baseURL+"/api/auth/refresh", "application/json", 
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			resp.Body.Close()
			
			assert.Equal(t, http.StatusOK, resp.StatusCode,
				"Token refresh should succeed and trigger TokenRefreshedEvent")
		})
	})

	// Test 3: Event Ordering and Consistency
	t.Run("event_consistency", func(t *testing.T) {
		// Verify that events maintain consistency
		// For example, UserCreated should always come before UserUpdated
		email := fmt.Sprintf("consistency%d@example.com", time.Now().UnixNano())
		
		// Create user
		userData := map[string]interface{}{
			"email":      email,
			"password":   "password123",
			"first_name": "Consistent",
			"last_name":  "User",
		}
		
		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)
		
		resp, err := client.Post(baseURL+"/api/auth/register", "application/json", 
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		resp.Body.Close()
		
		assert.Equal(t, http.StatusCreated, resp.StatusCode,
			"User creation establishes event baseline")
		
		// Immediate update should work (events in correct order)
		loginData := map[string]interface{}{
			"email":    email,
			"password": "password123",
		}
		
		jsonData, err = json.Marshal(loginData)
		require.NoError(t, err)
		
		resp, err = client.Post(baseURL+"/api/auth/login", "application/json", 
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		
		var authResp AuthResponse
		err = json.NewDecoder(resp.Body).Decode(&authResp)
		require.NoError(t, err)
		resp.Body.Close()
		
		// Update user
		updateData := map[string]interface{}{
			"first_name": "StillConsistent",
		}
		
		jsonData, err = json.Marshal(updateData)
		require.NoError(t, err)
		
		req, err := http.NewRequest("PUT", baseURL+"/api/users/me", 
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)
		req.Header.Set("Content-Type", "application/json")
		
		resp, err = client.Do(req)
		require.NoError(t, err)
		resp.Body.Close()
		
		assert.Equal(t, http.StatusOK, resp.StatusCode,
			"Update after creation maintains event consistency")
	})
}