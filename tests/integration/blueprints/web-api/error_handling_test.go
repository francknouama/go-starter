package webapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHexagonal_ErrorHandling tests comprehensive error handling scenarios
func TestHexagonal_ErrorHandling(t *testing.T) {
	testCases := []RuntimeTestConfig{
		{
			Name:       "gin_gorm_slog_postgres_errors",
			Framework:  "gin",
			ORM:        "gorm",
			Logger:     "slog",
			Database:   "postgres",
			ServerPort: 8091,
		},
		{
			Name:       "echo_sqlx_zap_postgres_errors",
			Framework:  "echo",
			ORM:        "sqlx",
			Logger:     "zap",
			Database:   "postgres",
			ServerPort: 8092,
		},
		{
			Name:       "fiber_gorm_logrus_postgres_errors",
			Framework:  "fiber",
			ORM:        "gorm",
			Logger:     "logrus",
			Database:   "postgres",
			ServerPort: 8093,
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

			// Run comprehensive error handling tests
			client := &http.Client{Timeout: 10 * time.Second}
			baseURL := fmt.Sprintf("http://localhost:%d", tc.ServerPort)

			// Test malformed request handling
			t.Run("malformed_requests", func(t *testing.T) {
				testMalformedRequests(t, client, baseURL)
			})

			// Test validation errors
			t.Run("validation_errors", func(t *testing.T) {
				testValidationErrors(t, client, baseURL)
			})

			// Test authentication errors
			t.Run("authentication_errors", func(t *testing.T) {
				testAuthenticationErrors(t, client, baseURL)
			})

			// Test resource not found errors
			t.Run("not_found_errors", func(t *testing.T) {
				testNotFoundErrors(t, client, baseURL)
			})

			// Test method not allowed errors
			t.Run("method_not_allowed", func(t *testing.T) {
				testMethodNotAllowed(t, client, baseURL)
			})

			// Test content type errors
			t.Run("content_type_errors", func(t *testing.T) {
				testContentTypeErrors(t, client, baseURL)
			})

			// Test concurrent request errors
			t.Run("concurrent_errors", func(t *testing.T) {
				testConcurrentErrors(t, client, baseURL)
			})
		})
	}
}

// testMalformedRequests tests various malformed request scenarios
func testMalformedRequests(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: Invalid JSON syntax
	t.Run("invalid_json_syntax", func(t *testing.T) {
		resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(`{"email": "test@example.com", "password": "pass123", invalid json`))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// Test 2: Empty request body
	t.Run("empty_request_body", func(t *testing.T) {
		resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(""))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// Test 3: Wrong data types
	t.Run("wrong_data_types", func(t *testing.T) {
		wrongTypeData := map[string]interface{}{
			"email":      123,                 // Should be string
			"password":   true,                // Should be string
			"first_name": []string{"John"},    // Should be string
			"last_name":  map[string]string{}, // Should be string
		}

		jsonData, err := json.Marshal(wrongTypeData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// Test 4: Oversized request body
	t.Run("oversized_request_body", func(t *testing.T) {
		// Create a very large string (10MB)
		largeString := strings.Repeat("a", 10*1024*1024)
		oversizedData := map[string]interface{}{
			"email":      "test@example.com",
			"password":   "password123",
			"first_name": largeString,
			"last_name":  "Doe",
		}

		jsonData, err := json.Marshal(oversizedData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should reject oversized payloads
		assert.True(t, resp.StatusCode == http.StatusBadRequest ||
			resp.StatusCode == http.StatusRequestEntityTooLarge)
	})

	// Test 5: Null values
	t.Run("null_values", func(t *testing.T) {
		nullData := map[string]interface{}{
			"email":      nil,
			"password":   "password123",
			"first_name": "John",
			"last_name":  nil,
		}

		jsonData, err := json.Marshal(nullData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

// testValidationErrors tests field validation error handling
func testValidationErrors(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: Missing required fields
	t.Run("missing_required_fields", func(t *testing.T) {
		incompleteData := map[string]interface{}{
			"email": "test@example.com",
			// Missing password, first_name, last_name
		}

		jsonData, err := json.Marshal(incompleteData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// Test 2: Invalid email format
	t.Run("invalid_email_format", func(t *testing.T) {
		invalidEmails := []string{
			"notanemail",
			"missing@domain",
			"@example.com",
			"user@",
			"user..name@example.com",
			"user@example..com",
		}

		for _, email := range invalidEmails {
			data := map[string]interface{}{
				"email":      email,
				"password":   "password123",
				"first_name": "John",
				"last_name":  "Doe",
			}

			jsonData, err := json.Marshal(data)
			require.NoError(t, err)

			resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
				"Expected bad request for email: %s", email)
		}
	})

	// Test 3: Password validation
	t.Run("password_validation", func(t *testing.T) {
		invalidPasswords := []string{
			"",        // Empty
			"short",   // Too short (less than 8 chars)
			"       ", // Only spaces
		}

		for _, password := range invalidPasswords {
			data := map[string]interface{}{
				"email":      "test@example.com",
				"password":   password,
				"first_name": "John",
				"last_name":  "Doe",
			}

			jsonData, err := json.Marshal(data)
			require.NoError(t, err)

			resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
				"Expected bad request for password: '%s'", password)
		}
	})

	// Test 4: Name validation
	t.Run("name_validation", func(t *testing.T) {
		invalidNames := []struct {
			firstName string
			lastName  string
		}{
			{"", "Doe"},                       // Empty first name
			{"John", ""},                      // Empty last name
			{"J", "Doe"},                      // Too short first name
			{"John", "D"},                     // Too short last name
			{strings.Repeat("A", 51), "Doe"},  // Too long first name
			{"John", strings.Repeat("B", 51)}, // Too long last name
		}

		for _, names := range invalidNames {
			data := map[string]interface{}{
				"email":      "test@example.com",
				"password":   "password123",
				"first_name": names.firstName,
				"last_name":  names.lastName,
			}

			jsonData, err := json.Marshal(data)
			require.NoError(t, err)

			resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
				strings.NewReader(string(jsonData)))
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusBadRequest, resp.StatusCode,
				"Expected bad request for names: '%s' '%s'", names.firstName, names.lastName)
		}
	})
}

// testAuthenticationErrors tests authentication-related error scenarios
func testAuthenticationErrors(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: Login with non-existent user
	t.Run("login_non_existent_user", func(t *testing.T) {
		loginData := map[string]interface{}{
			"email":    "nonexistent@example.com",
			"password": "password123",
		}

		jsonData, err := json.Marshal(loginData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/login", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	// Test 2: Login with wrong password
	t.Run("login_wrong_password", func(t *testing.T) {
		// First register a user
		registerData := map[string]interface{}{
			"email":      "wrongpass@example.com",
			"password":   "correctpassword",
			"first_name": "Wrong",
			"last_name":  "Pass",
		}

		jsonData, err := json.Marshal(registerData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		resp.Body.Close()
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		// Try to login with wrong password
		loginData := map[string]interface{}{
			"email":    "wrongpass@example.com",
			"password": "wrongpassword",
		}

		jsonData, err = json.Marshal(loginData)
		require.NoError(t, err)

		resp, err = client.Post(baseURL+"/api/auth/login", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	// Test 3: Access protected endpoint without token
	t.Run("protected_endpoint_no_token", func(t *testing.T) {
		resp, err := client.Get(baseURL + "/api/users/me")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	// Test 4: Access protected endpoint with invalid token
	t.Run("protected_endpoint_invalid_token", func(t *testing.T) {
		req, err := http.NewRequest("GET", baseURL+"/api/users/me", nil)
		require.NoError(t, err)

		req.Header.Set("Authorization", "Bearer invalid_token_12345")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	// Test 5: Malformed authorization header
	t.Run("malformed_auth_header", func(t *testing.T) {
		malformedHeaders := []string{
			"",                    // Empty
			"Bearer",              // Missing token
			"InvalidScheme token", // Wrong scheme
			"Bearer ",             // Empty token after Bearer
			" Bearer token",       // Leading space
		}

		for _, header := range malformedHeaders {
			req, err := http.NewRequest("GET", baseURL+"/api/users/me", nil)
			require.NoError(t, err)

			if header != "" {
				req.Header.Set("Authorization", header)
			}

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode,
				"Expected unauthorized for header: '%s'", header)
		}
	})
}

// testNotFoundErrors tests 404 error scenarios
func testNotFoundErrors(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: Non-existent endpoints
	t.Run("non_existent_endpoints", func(t *testing.T) {
		nonExistentPaths := []string{
			"/api/nonexistent",
			"/api/v1/doesnotexist",
			"/api/users/invalid-uuid-format",
			"/api/auth/nonexistent-action",
		}

		for _, path := range nonExistentPaths {
			resp, err := client.Get(baseURL + path)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusNotFound, resp.StatusCode,
				"Expected not found for path: %s", path)
		}
	})

	// Test 2: Get non-existent user by ID
	t.Run("get_non_existent_user", func(t *testing.T) {
		// First register and login to get a valid token
		registerData := map[string]interface{}{
			"email":      "getuser@example.com",
			"password":   "password123",
			"first_name": "Get",
			"last_name":  "User",
		}

		jsonData, err := json.Marshal(registerData)
		require.NoError(t, err)

		resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
			strings.NewReader(string(jsonData)))
		require.NoError(t, err)
		resp.Body.Close()

		// Login
		loginData := map[string]interface{}{
			"email":    "getuser@example.com",
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

		// Try to get non-existent user
		req, err := http.NewRequest("GET", baseURL+"/api/users/00000000-0000-0000-0000-000000000000", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+authResp.AccessToken)

		resp, err = client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// testMethodNotAllowed tests HTTP method validation
func testMethodNotAllowed(t *testing.T, client *http.Client, baseURL string) {
	// Test wrong HTTP methods on various endpoints
	wrongMethods := []struct {
		method   string
		endpoint string
	}{
		{"GET", "/api/auth/register"},    // Should be POST
		{"GET", "/api/auth/login"},       // Should be POST
		{"POST", "/api/users/me"},        // Should be GET
		{"DELETE", "/api/auth/register"}, // Should be POST
		{"PUT", "/api/auth/login"},       // Should be POST
		{"PATCH", "/api/health"},         // Should be GET
	}

	for _, test := range wrongMethods {
		req, err := http.NewRequest(test.method, baseURL+test.endpoint, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode,
			"Expected method not allowed for %s %s", test.method, test.endpoint)
	}
}

// testContentTypeErrors tests content type validation
func testContentTypeErrors(t *testing.T, client *http.Client, baseURL string) {
	// Test 1: Wrong content type
	t.Run("wrong_content_type", func(t *testing.T) {
		data := `{"email":"test@example.com","password":"password123"}`

		wrongContentTypes := []string{
			"text/plain",
			"application/xml",
			"application/x-www-form-urlencoded",
			"multipart/form-data",
			"", // Empty content type
		}

		for _, contentType := range wrongContentTypes {
			req, err := http.NewRequest("POST", baseURL+"/api/auth/login",
				strings.NewReader(data))
			require.NoError(t, err)

			if contentType != "" {
				req.Header.Set("Content-Type", contentType)
			}

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Some frameworks are more lenient with content types
			assert.True(t, resp.StatusCode == http.StatusBadRequest ||
				resp.StatusCode == http.StatusUnsupportedMediaType,
				"Expected bad request or unsupported media type for content type: '%s'", contentType)
		}
	})

	// Test 2: Missing content type header
	t.Run("missing_content_type", func(t *testing.T) {
		data := `{"email":"test@example.com","password":"password123"}`

		req, err := http.NewRequest("POST", baseURL+"/api/auth/login",
			strings.NewReader(data))
		require.NoError(t, err)

		// Explicitly remove content type
		req.Header.Del("Content-Type")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Should reject or handle gracefully
		assert.True(t, resp.StatusCode >= 400)
	})
}

// testConcurrentErrors tests error handling under concurrent load
func testConcurrentErrors(t *testing.T, client *http.Client, baseURL string) {
	// Test concurrent registration attempts with same email
	t.Run("concurrent_duplicate_registration", func(t *testing.T) {
		email := "concurrent@example.com"

		// Function to attempt registration
		attemptRegistration := func() int {
			data := map[string]interface{}{
				"email":      email,
				"password":   "password123",
				"first_name": "Concurrent",
				"last_name":  "User",
			}

			jsonData, _ := json.Marshal(data)
			resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
				strings.NewReader(string(jsonData)))
			if err != nil {
				return 0
			}
			defer resp.Body.Close()

			return resp.StatusCode
		}

		// Run 10 concurrent registration attempts
		results := make(chan int, 10)
		for i := 0; i < 10; i++ {
			go func() {
				results <- attemptRegistration()
			}()
		}

		// Collect results
		successCount := 0
		conflictCount := 0

		for i := 0; i < 10; i++ {
			statusCode := <-results
			if statusCode == http.StatusCreated {
				successCount++
			} else if statusCode == http.StatusConflict || statusCode == http.StatusBadRequest {
				conflictCount++
			}
		}

		// Exactly one should succeed, others should fail with conflict
		assert.Equal(t, 1, successCount, "Exactly one registration should succeed")
		assert.Equal(t, 9, conflictCount, "Nine registrations should fail with conflict")
	})

	// Test concurrent invalid requests
	t.Run("concurrent_invalid_requests", func(t *testing.T) {
		// Function to send invalid request
		sendInvalidRequest := func() int {
			resp, err := client.Post(baseURL+"/api/auth/register", "application/json",
				strings.NewReader("invalid json"))
			if err != nil {
				return 0
			}
			defer resp.Body.Close()

			return resp.StatusCode
		}

		// Run 20 concurrent invalid requests
		results := make(chan int, 20)
		for i := 0; i < 20; i++ {
			go func() {
				results <- sendInvalidRequest()
			}()
		}

		// All should return bad request
		for i := 0; i < 20; i++ {
			statusCode := <-results
			assert.Equal(t, http.StatusBadRequest, statusCode)
		}
	})
}
