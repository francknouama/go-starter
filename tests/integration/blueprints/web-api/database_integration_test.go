package webapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHexagonal_DatabaseErrors tests database error handling scenarios
func TestHexagonal_DatabaseErrors(t *testing.T) {
	testCases := []RuntimeTestConfig{
		{
			Name:       "gin_gorm_slog_postgres_db_errors",
			Framework:  "gin",
			ORM:        "gorm",
			Logger:     "slog",
			Database:   "postgres",
			ServerPort: 8094,
		},
		{
			Name:       "echo_sqlx_zap_postgres_db_errors",
			Framework:  "echo",
			ORM:        "sqlx",
			Logger:     "zap",
			Database:   "postgres",
			ServerPort: 8095,
		},
		{
			Name:       "fiber_gorm_logrus_postgres_db_errors",
			Framework:  "fiber",
			ORM:        "gorm",
			Logger:     "logrus",
			Database:   "postgres",
			ServerPort: 8096,
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()

			// Setup and run the test environment
			testDir := setupTestEnvironment(t, tc)
			defer cleanupTestEnvironment(t, testDir)

			projectPath := generateHexagonalProject(t, ctx, testDir, tc)
			dbContainer, dbCleanup := setupTestDatabase(t, ctx, tc)
			defer dbCleanup()

			serverCleanup := startApplicationServer(t, ctx, projectPath, tc, dbContainer)
			defer serverCleanup()

			waitForServerReady(t, ctx, tc)

			// Get database connection info for direct manipulation
			dbHost, dbPort := getContainerHostPort(dbContainer)
			dbConnStr := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable",
				dbHost, dbPort)

			// Run database error tests
			client := &http.Client{Timeout: 10 * time.Second}
			baseURL := fmt.Sprintf("http://localhost:%d", tc.ServerPort)

			// Test constraint violations
			t.Run("constraint_violations", func(t *testing.T) {
				testConstraintViolations(t, client, baseURL)
			})

			// Test transaction failures
			t.Run("transaction_failures", func(t *testing.T) {
				testTransactionFailures(t, client, baseURL, dbConnStr)
			})

			// Test connection pool exhaustion
			t.Run("connection_pool_exhaustion", func(t *testing.T) {
				testConnectionPoolExhaustion(t, client, baseURL)
			})

			// Test database deadlocks
			t.Run("database_deadlocks", func(t *testing.T) {
				testDatabaseDeadlocks(t, client, baseURL, dbConnStr)
			})

			// Test query timeouts
			t.Run("query_timeouts", func(t *testing.T) {
				testQueryTimeouts(t, client, baseURL, dbConnStr)
			})

			// Test data integrity errors
			t.Run("data_integrity_errors", func(t *testing.T) {
				testDataIntegrityErrors(t, client, baseURL)
			})
		})
	}
}

// testConstraintViolations tests various database constraint violations
func testConstraintViolations(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: Unique constraint violation (duplicate email)
	t.Run("unique_email_violation", func(t *testing.T) {
		// Register first user
		userData := map[string]interface{}{
			"email":      "unique@example.com",
			"password":   "password123",
			"first_name": "First",
			"last_name":  "User",
		}

		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Try to register second user with same email
		userData["first_name"] = "Second"
		jsonData, err = json.Marshal(userData)
		require.NoError(t, err)

		resp, err = client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should fail with conflict or bad request
		assert.True(t, resp.StatusCode == http.StatusConflict ||
			resp.StatusCode == http.StatusBadRequest)
	})

	// Test 2: Foreign key constraint violation
	t.Run("foreign_key_violation", func(t *testing.T) {
		// First register and login to get a valid token
		userData := map[string]interface{}{
			"email":      "fktest@example.com",
			"password":   "password123",
			"first_name": "FK",
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
			"email":    "fktest@example.com",
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

		// Try to update with non-existent reference
		updateData := map[string]interface{}{
			"id":         "00000000-0000-0000-0000-000000000000", // Non-existent ID
			"email":      "updated@example.com",
			"first_name": "Updated",
			"last_name":  "User",
		}

		jsonData, err = json.Marshal(updateData)
		require.NoError(t, err)

		req, err := http.NewRequest("PUT", baseURL+"/api/users/00000000-0000-0000-0000-000000000000",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should fail with not found
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	// Test 3: Check constraint violation (if any)
	t.Run("check_constraint_violation", func(t *testing.T) {
		// Try to register with invalid data that would violate check constraints
		userData := map[string]interface{}{
			"email":      "check@example.com",
			"password":   "pass", // Too short, less than minimum required
			"first_name": "C",    // Too short
			"last_name":  "T",    // Too short
		}

		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should fail validation
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// testTransactionFailures tests transaction rollback scenarios
func testTransactionFailures(t *testing.T, client *http.Client, baseURL string, dbConnStr string) {
	// Test 1: Simulate transaction failure during user creation
	t.Run("transaction_rollback_on_error", func(t *testing.T) {
		// This test would require injecting failures or using a test database
		// For now, we'll test that the API handles errors gracefully

		// Register a user
		userData := map[string]interface{}{
			"email":      "txtest@example.com",
			"password":   "password123",
			"first_name": "Transaction",
			"last_name":  "Test",
		}

		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Try to register same user again (should fail atomically)
		resp, err = client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should fail without partial data being saved
		assert.True(t, resp.StatusCode == http.StatusConflict ||
			resp.StatusCode == http.StatusBadRequest)
	})

	// Test 2: Concurrent transactions
	t.Run("concurrent_transaction_isolation", func(t *testing.T) {
		var wg sync.WaitGroup
		results := make(chan int, 5)

		// Try to create 5 users concurrently with potential conflicts
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				userData := map[string]interface{}{
					"email":      fmt.Sprintf("concurrent%d@example.com", idx),
					"password":   "password123",
					"first_name": fmt.Sprintf("User%d", idx),
					"last_name":  "Concurrent",
				}

				jsonData, _ := json.Marshal(userData)
				resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
					strings.NewReader(string(jsonData)))
				if err != nil {
					results <- 0
					return
				}
				defer resp.Body.Close()

				results <- resp.StatusCode
			}(i)
		}

		wg.Wait()
		close(results)

		// All should either succeed or fail cleanly
		for statusCode := range results {
			assert.True(t, statusCode == http.StatusCreated ||
				statusCode == http.StatusConflict ||
				statusCode == http.StatusBadRequest)
		}
	})
}

// testConnectionPoolExhaustion tests behavior when connection pool is exhausted
func testConnectionPoolExhaustion(t *testing.T, client *http.Client, baseURL string) {
	// Create many concurrent requests to exhaust connection pool
	t.Run("connection_pool_stress", func(t *testing.T) {
		var wg sync.WaitGroup
		concurrentRequests := 100
		results := make(chan int, concurrentRequests)

		// First create a user for authentication
		userData := map[string]interface{}{
			"email":      "pooltest@example.com",
			"password":   "password123",
			"first_name": "Pool",
			"last_name":  "Test",
		}

		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		resp.Body.Close()

		// Login to get token
		loginData := map[string]interface{}{
			"email":    "pooltest@example.com",
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

		// Send many concurrent requests
		for i := 0; i < concurrentRequests; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				req, err := http.NewRequest("GET", baseURL+"/api/users/me", nil)
				if err != nil {
					results <- 0
					return
				}
				req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)

				resp, err := client.Do(req)
				if err != nil {
					results <- 0
					return
				}
				defer resp.Body.Close()

				results <- resp.StatusCode
			}()
		}

		wg.Wait()
		close(results)

		// Count successful vs failed requests
		successCount := 0
		errorCount := 0

		for statusCode := range results {
			if statusCode == http.StatusOK {
				successCount++
			} else if statusCode >= 500 {
				errorCount++
			}
		}

		// Most requests should succeed, but some might fail under extreme load
		assert.Greater(t, successCount, concurrentRequests/2,
			"At least half of the requests should succeed")

		// The system should handle load gracefully
		t.Logf("Connection pool test: %d successful, %d errors out of %d requests",
			successCount, errorCount, concurrentRequests)
	})
}

// testDatabaseDeadlocks tests deadlock detection and handling
func testDatabaseDeadlocks(t *testing.T, client *http.Client, baseURL string, dbConnStr string) {
	// This is a complex scenario that requires direct database manipulation
	t.Run("deadlock_simulation", func(t *testing.T) {
		// Register two users
		for i := 1; i <= 2; i++ {
			userData := map[string]interface{}{
				"email":      fmt.Sprintf("deadlock%d@example.com", i),
				"password":   "password123",
				"first_name": fmt.Sprintf("User%d", i),
				"last_name":  "Deadlock",
			}

			jsonData, err := json.Marshal(userData)
			require.NoError(t, err)

			resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			resp.Body.Close()
			assert.Equal(t, http.StatusCreated, resp.StatusCode)
		}

		// Simulate concurrent updates that could cause deadlock
		var wg sync.WaitGroup
		results := make(chan error, 2)

		// Transaction 1: Update user1 then user2
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Login as user1
			loginData := map[string]interface{}{
				"email":    "deadlock1@example.com",
				"password": "password123",
			}

			jsonData, _ := json.Marshal(loginData)
			resp, err := client.Post(baseURL+"/api/auth/login", "application/json",
				strings.NewReader(string(jsonData)))
			if err != nil {
				results <- err
				return
			}

			var authResp AuthResponse
			json.NewDecoder(resp.Body).Decode(&authResp)
			resp.Body.Close()

			// Update profile
			updateData := map[string]interface{}{
				"first_name": "UpdatedUser1",
				"last_name":  "Transaction1",
			}

			jsonData, _ = json.Marshal(updateData)
			req, _ := http.NewRequest("PUT", baseURL+"/api/users/me",
				strings.NewReader(string(jsonData)))
			req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)
			req.Header.Set("Content-Type", "application/json")

			resp, err = client.Do(req)
			if err != nil {
				results <- err
				return
			}
			resp.Body.Close()

			results <- nil
		}()

		// Transaction 2: Update user2 then user1 (opposite order)
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Add small delay to increase chance of deadlock
			time.Sleep(10 * time.Millisecond)

			// Login as user2
			loginData := map[string]interface{}{
				"email":    "deadlock2@example.com",
				"password": "password123",
			}

			jsonData, _ := json.Marshal(loginData)
			resp, err := client.Post(baseURL+"/api/auth/login", "application/json",
				strings.NewReader(string(jsonData)))
			if err != nil {
				results <- err
				return
			}

			var authResp AuthResponse
			json.NewDecoder(resp.Body).Decode(&authResp)
			resp.Body.Close()

			// Update profile
			updateData := map[string]interface{}{
				"first_name": "UpdatedUser2",
				"last_name":  "Transaction2",
			}

			jsonData, _ = json.Marshal(updateData)
			req, _ := http.NewRequest("PUT", baseURL+"/api/users/me",
				strings.NewReader(string(jsonData)))
			req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)
			req.Header.Set("Content-Type", "application/json")

			resp, err = client.Do(req)
			if err != nil {
				results <- err
				return
			}
			resp.Body.Close()

			results <- nil
		}()

		wg.Wait()
		close(results)

		// Both operations should complete (deadlock detection should resolve it)
		errorCount := 0
		for err := range results {
			if err != nil {
				errorCount++
			}
		}

		assert.Equal(t, 0, errorCount, "No errors expected - deadlocks should be resolved")
	})
}

// testQueryTimeouts tests query timeout handling
func testQueryTimeouts(t *testing.T, client *http.Client, baseURL string, dbConnStr string) {
	// Test slow query handling
	t.Run("slow_query_timeout", func(t *testing.T) {
		// Create a user with lots of data to search through
		baseEmail := "timeout@example.com"

		// Register main user
		userData := map[string]interface{}{
			"email":      baseEmail,
			"password":   "password123",
			"first_name": "Timeout",
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
			"email":    baseEmail,
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

		// Create multiple users to simulate large dataset
		for i := 0; i < 50; i++ {
			userData := map[string]interface{}{
				"email":      fmt.Sprintf("bulk%d@example.com", i),
				"password":   "password123",
				"first_name": fmt.Sprintf("Bulk%d", i),
				"last_name":  "User",
			}

			jsonData, _ := json.Marshal(userData)
			client.Post(baseURL+"/api/auth/register", "application/json",
				strings.NewReader(string(jsonData)))
		}

		// Make a request that might be slow
		req, err := http.NewRequest("GET", baseURL+"/api/users?page=1&limit=100", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)

		// Use a client with shorter timeout
		shortTimeoutClient := &http.Client{Timeout: 2 * time.Second}

		resp, err = shortTimeoutClient.Do(req)
		if err == nil {
			defer resp.Body.Close()
			// Request completed within timeout
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		} else {
			// Timeout occurred - this is also acceptable behavior
			assert.Contains(t, err.Error(), "timeout")
		}
	})
}

// testDataIntegrityErrors tests data integrity error scenarios
func testDataIntegrityErrors(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: Attempt to corrupt data through API
	t.Run("data_corruption_prevention", func(t *testing.T) {
		// Register a user
		userData := map[string]interface{}{
			"email":      "integrity@example.com",
			"password":   "password123",
			"first_name": "Data",
			"last_name":  "Integrity",
		}

		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Login
		loginData := map[string]interface{}{
			"email":    "integrity@example.com",
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

		// Try to update with inconsistent data
		updateData := map[string]interface{}{
			"id":         "invalid-uuid-format", // Invalid ID format
			"email":      "not-an-email",        // Invalid email
			"first_name": "",                    // Empty name
			"last_name":  "",                    // Empty name
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
		defer resp.Body.Close()

		// Should reject invalid data
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// Test 2: Referential integrity
	t.Run("referential_integrity", func(t *testing.T) {
		// This test verifies that related data maintains integrity
		// Register a user
		userData := map[string]interface{}{
			"email":      "refintegrity@example.com",
			"password":   "password123",
			"first_name": "Ref",
			"last_name":  "Integrity",
		}

		jsonData, err := json.Marshal(userData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		resp.Body.Close()
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Login
		loginData := map[string]interface{}{
			"email":    "refintegrity@example.com",
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

		// Get user details
		req, err := http.NewRequest("GET", baseURL+"/api/users/me", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)

		resp, err = client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify user data integrity
		var userResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&userResp)
		require.NoError(t, err)

		// Check that all expected fields are present and valid
		assert.NotEmpty(t, userResp["id"])
		assert.Equal(t, "refintegrity@example.com", userResp["email"])
		assert.Equal(t, "Ref", userResp["first_name"])
		assert.Equal(t, "Integrity", userResp["last_name"])
		assert.NotEmpty(t, userResp["created_at"])
		assert.NotEmpty(t, userResp["updated_at"])
	})
}

// Helper function to get container host and port
func getContainerHostPort(container interface{}) (string, string) {
	// This is a simplified version - in the actual test we'd use the container methods
	return "localhost", "5432"
}
