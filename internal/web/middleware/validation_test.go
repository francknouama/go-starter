package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMiddlewareTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestValidationMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		method           string
		url              string
		expectedHeaders  map[string]string
	}{
		{
			name:   "adds security headers",
			method: "GET",
			url:    "/test",
			expectedHeaders: map[string]string{
				"X-Content-Type-Options": "nosniff",
				"X-Frame-Options":        "DENY",
				"X-XSS-Protection":       "1; mode=block",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupMiddlewareTestRouter()
			router.Use(ValidationMiddleware())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "test"})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.url, nil)
			router.ServeHTTP(w, req)

			// Verify security headers are added
			for header, expectedValue := range tt.expectedHeaders {
				assert.Equal(t, expectedValue, w.Header().Get(header), 
					"Header %s should be %s", header, expectedValue)
			}

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestContentTypeValidation(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		contentType    string
		expectedStatus int
		expectError    bool
		expectedError  string
		expectedCode   string
	}{
		{
			name:           "GET request bypasses content type validation",
			method:         "GET",
			url:            "/test",
			contentType:    "",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "POST with valid JSON content type",
			method:         "POST",
			url:            "/test",
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "POST with JSON charset content type",
			method:         "POST",
			url:            "/test",
			contentType:    "application/json; charset=utf-8",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "POST with invalid content type",
			method:         "POST",
			url:            "/test",
			contentType:    "text/plain",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			expectedError:  "Content-Type must be application/json",
			expectedCode:   "INVALID_CONTENT_TYPE",
		},
		{
			name:           "POST with no content type",
			method:         "POST",
			url:            "/test",
			contentType:    "",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			expectedError:  "Content-Type must be application/json",
			expectedCode:   "INVALID_CONTENT_TYPE",
		},
		{
			name:           "PUT with valid JSON content type",
			method:         "PUT",
			url:            "/test",
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "PUT with invalid content type",
			method:         "PUT",
			url:            "/test",
			contentType:    "application/xml",
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			expectedError:  "Content-Type must be application/json",
			expectedCode:   "INVALID_CONTENT_TYPE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupMiddlewareTestRouter()
			router.Use(ContentTypeValidation())
			
			// Add handlers for different methods
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "test"})
			})
			router.POST("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "test"})
			})
			router.PUT("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "test"})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.url, bytes.NewBuffer([]byte(`{"test": "data"}`)))
			
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectError {
				var errorResponse ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				require.NoError(t, err, "Error response should be valid JSON")
				
				assert.Equal(t, tt.expectedError, errorResponse.Error)
				assert.Equal(t, tt.expectedCode, errorResponse.Code)
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err, "Success response should be valid JSON")
				assert.Equal(t, "test", response["message"])
			}
		})
	}
}

func TestRateLimiting(t *testing.T) {
	tests := []struct {
		name            string
		expectedHeaders map[string]string
	}{
		{
			name: "adds rate limiting headers",
			expectedHeaders: map[string]string{
				"X-RateLimit-Limit":     "100",
				"X-RateLimit-Remaining": "99",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupMiddlewareTestRouter()
			router.Use(RateLimiting())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "test"})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			// Verify rate limiting headers are added
			for header, expectedValue := range tt.expectedHeaders {
				assert.Equal(t, expectedValue, w.Header().Get(header), 
					"Header %s should be %s", header, expectedValue)
			}

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestCORS(t *testing.T) {
	tests := []struct {
		name                    string
		method                  string
		origin                  string
		expectedStatus          int
		expectCORSHeaders       bool
		expectedAllowedOrigin   string
	}{
		{
			name:                  "allowed origin localhost:3000",
			method:                "GET",
			origin:                "http://localhost:3000",
			expectedStatus:        http.StatusOK,
			expectCORSHeaders:     true,
			expectedAllowedOrigin: "http://localhost:3000",
		},
		{
			name:                  "allowed origin localhost:5173 (Vite)",
			method:                "GET",
			origin:                "http://localhost:5173",
			expectedStatus:        http.StatusOK,
			expectCORSHeaders:     true,
			expectedAllowedOrigin: "http://localhost:5173",
		},
		{
			name:                  "allowed origin localhost:8080",
			method:                "GET",
			origin:                "http://localhost:8080",
			expectedStatus:        http.StatusOK,
			expectCORSHeaders:     true,
			expectedAllowedOrigin: "http://localhost:8080",
		},
		{
			name:              "disallowed origin",
			method:            "GET",
			origin:            "http://malicious-site.com",
			expectedStatus:    http.StatusOK,
			expectCORSHeaders: false,
		},
		{
			name:              "no origin header",
			method:            "GET",
			origin:            "",
			expectedStatus:    http.StatusOK,
			expectCORSHeaders: false,
		},
		{
			name:           "OPTIONS preflight request with allowed origin",
			method:         "OPTIONS",
			origin:         "http://localhost:3000",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "OPTIONS preflight request with disallowed origin",
			method:         "OPTIONS",
			origin:         "http://malicious-site.com",
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupMiddlewareTestRouter()
			router.Use(CORS())
			router.GET("/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "test"})
			})
			router.OPTIONS("/test", func(c *gin.Context) {
				// This should not be reached due to middleware handling
				c.JSON(http.StatusOK, gin.H{"message": "options"})
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, "/test", nil)
			
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectCORSHeaders {
				assert.Equal(t, tt.expectedAllowedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
				assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
				assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
			} else {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"), 
					"Should not set CORS headers for disallowed origin")
			}
		})
	}
}

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name           string
		setupError     func(*gin.Context)
		expectedStatus int
		expectedCode   string
		expectResponse bool
	}{
		{
			name: "binding error",
			setupError: func(c *gin.Context) {
				c.Error(gin.Error{
					Err:  assert.AnError,
					Type: gin.ErrorTypeBind,
				})
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "BINDING_ERROR",
			expectResponse: true,
		},
		{
			name: "public error",
			setupError: func(c *gin.Context) {
				c.Error(gin.Error{
					Err:  assert.AnError,
					Type: gin.ErrorTypePublic,
				})
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "VALIDATION_ERROR",
			expectResponse: true,
		},
		{
			name: "internal error",
			setupError: func(c *gin.Context) {
				c.Error(gin.Error{
					Err:  assert.AnError,
					Type: gin.ErrorTypePrivate,
				})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
			expectResponse: true,
		},
		{
			name: "no error - passes through",
			setupError: func(c *gin.Context) {
				// No error added
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			},
			expectedStatus: http.StatusOK,
			expectResponse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupMiddlewareTestRouter()
			router.Use(ErrorHandler())
			router.GET("/test", func(c *gin.Context) {
				tt.setupError(c)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectResponse {
				var errorResponse ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				require.NoError(t, err, "Error response should be valid JSON")
				
				assert.Equal(t, tt.expectedCode, errorResponse.Code)
				assert.NotEmpty(t, errorResponse.Error)
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err, "Success response should be valid JSON")
				assert.Equal(t, "success", response["message"])
			}
		})
	}
}

func TestMiddlewareChaining(t *testing.T) {
	t.Run("multiple middleware work together", func(t *testing.T) {
		router := setupMiddlewareTestRouter()
		router.Use(ValidationMiddleware())
		router.Use(ContentTypeValidation())
		router.Use(RateLimiting())
		router.Use(CORS())
		router.Use(ErrorHandler())

		router.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer([]byte(`{"test": "data"}`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "http://localhost:3000")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Check that all middleware added their headers
		assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
		assert.Equal(t, "100", w.Header().Get("X-RateLimit-Limit"))
		assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Response should be valid JSON")
		assert.Equal(t, "success", response["message"])
	})

	t.Run("middleware stops chain on validation failure", func(t *testing.T) {
		router := setupMiddlewareTestRouter()
		router.Use(ValidationMiddleware())
		router.Use(ContentTypeValidation())
		router.Use(RateLimiting())

		handlerCalled := false
		router.POST("/test", func(c *gin.Context) {
			handlerCalled = true
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer([]byte(`{"test": "data"}`)))
		req.Header.Set("Content-Type", "text/plain") // Invalid content type

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.False(t, handlerCalled, "Handler should not be called when middleware validation fails")

		var errorResponse ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
		require.NoError(t, err, "Error response should be valid JSON")
		assert.Equal(t, "INVALID_CONTENT_TYPE", errorResponse.Code)
	})
}

// Benchmark tests for middleware performance
func BenchmarkValidationMiddleware(b *testing.B) {
	router := setupMiddlewareTestRouter()
	router.Use(ValidationMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkMiddlewareChain(b *testing.B) {
	router := setupMiddlewareTestRouter()
	router.Use(ValidationMiddleware())
	router.Use(ContentTypeValidation())
	router.Use(RateLimiting())
	router.Use(CORS())
	router.Use(ErrorHandler())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}