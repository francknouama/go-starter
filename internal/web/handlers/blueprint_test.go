package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/francknouama/go-starter/internal/web/models"
)

func setupBlueprintTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	handler := NewBlueprintHandler()
	
	// Blueprint API routes
	api := router.Group("/api/v1")
	{
		api.GET("/blueprints", handler.ListBlueprints)
		api.GET("/blueprints/:id", handler.GetBlueprint)
		api.GET("/blueprints/:type/defaults", handler.GetBlueprintDefaults)
		api.GET("/blueprints/:type/options", handler.GetBlueprintOptions)
		api.POST("/blueprints/:type/validate", handler.ValidateBlueprintConfig)
	}
	
	return router
}

func TestBlueprintHandler_ListBlueprints(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		validateFunc   func(t *testing.T, response *models.BlueprintListResponse)
	}{
		{
			name:           "successfully lists all blueprints",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *models.BlueprintListResponse) {
				assert.NotEmpty(t, response.Blueprints, "Should return non-empty blueprint list")
				
				// Verify we have the expected production blueprints
				blueprintIDs := make(map[string]bool)
				for _, bp := range response.Blueprints {
					blueprintIDs[bp.ID] = true
					
					// Validate basic blueprint structure
					assert.NotEmpty(t, bp.Name, "Blueprint name should not be empty")
					assert.NotEmpty(t, bp.Type, "Blueprint type should not be empty")
					assert.NotEmpty(t, bp.Complexity, "Blueprint complexity should not be empty")
					assert.Greater(t, bp.FileCount, 0, "Blueprint should have at least one file")
				}
				
				// Verify key production blueprints exist
				expectedBlueprints := []string{
					"cli-simple", "cli-standard", "web-api-standard", 
					"web-api-clean", "grpc-gateway", "microservice-standard",
				}
				for _, expected := range expectedBlueprints {
					assert.True(t, blueprintIDs[expected], 
						"Expected blueprint %s should be present", expected)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupBlueprintTestRouter()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/v1/blueprints", nil)
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if tt.expectedStatus == http.StatusOK {
				var response models.BlueprintListResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err, "Response should be valid JSON")
				
				if tt.validateFunc != nil {
					tt.validateFunc(t, &response)
				}
			}
		})
	}
}

func TestBlueprintHandler_GetBlueprint(t *testing.T) {
	tests := []struct {
		name           string
		blueprintID    string
		expectedStatus int
		validateFunc   func(t *testing.T, response *models.BlueprintDetailResponse)
		validateError  func(t *testing.T, errorResponse map[string]interface{})
	}{
		{
			name:           "successfully get cli-simple blueprint",
			blueprintID:    "cli-simple",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *models.BlueprintDetailResponse) {
				bp := response.Blueprint
				assert.Equal(t, "cli-simple", bp.ID)
				assert.Equal(t, "simple", bp.Complexity)
				assert.NotEmpty(t, bp.Name)
				assert.NotEmpty(t, bp.Description)
				assert.Greater(t, bp.FileCount, 0)
				
				// Verify files are included in detailed response
				assert.NotEmpty(t, response.Files, "Detailed response should include files")
				assert.NotEmpty(t, response.Variables, "Detailed response should include variables")
			},
		},
		{
			name:           "successfully get grpc-gateway blueprint",
			blueprintID:    "grpc-gateway",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *models.BlueprintDetailResponse) {
				bp := response.Blueprint
				assert.Equal(t, "grpc-gateway", bp.ID)
				assert.Equal(t, "expert", bp.Complexity)
				assert.NotEmpty(t, bp.Features, "gRPC Gateway should have features")
				assert.NotEmpty(t, bp.Dependencies, "gRPC Gateway should have dependencies")
				
				// Verify advanced blueprint has more files
				assert.Greater(t, bp.FileCount, 20, "gRPC Gateway should have many files")
			},
		},
		{
			name:           "returns 404 for non-existent blueprint",
			blueprintID:    "non-existent-blueprint",
			expectedStatus: http.StatusNotFound,
			validateError: func(t *testing.T, errorResponse map[string]interface{}) {
				assert.Equal(t, "Blueprint not found", errorResponse["error"])
				assert.Equal(t, "BLUEPRINT_NOT_FOUND", errorResponse["code"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupBlueprintTestRouter()
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/v1/blueprints/%s", tt.blueprintID)
			req, _ := http.NewRequest("GET", url, nil)
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if tt.expectedStatus == http.StatusOK && tt.validateFunc != nil {
				var response models.BlueprintDetailResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err, "Response should be valid JSON")
				tt.validateFunc(t, &response)
			} else if tt.expectedStatus != http.StatusOK && tt.validateError != nil {
				var errorResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				require.NoError(t, err, "Error response should be valid JSON")
				tt.validateError(t, errorResponse)
			}
		})
	}
}

func TestBlueprintHandler_GetBlueprintDefaults(t *testing.T) {
	tests := []struct {
		name           string
		blueprintType  string
		expectedStatus int
		validateFunc   func(t *testing.T, response map[string]interface{})
	}{
		{
			name:           "get defaults for web-api blueprint",
			blueprintType:  "web-api-standard",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response map[string]interface{}) {
				config, ok := response["config"].(map[string]interface{})
				require.True(t, ok, "Response should contain config object")
				
				assert.Equal(t, "web-api-standard", config["project_type"])
				assert.Equal(t, "gin", config["framework"])
				assert.Equal(t, "standard", config["architecture"])
				assert.Equal(t, "slog", config["logger"])
				assert.NotEmpty(t, config["go_version"])
			},
		},
		{
			name:           "get defaults for cli blueprint",
			blueprintType:  "cli-standard",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response map[string]interface{}) {
				config, ok := response["config"].(map[string]interface{})
				require.True(t, ok, "Response should contain config object")
				
				assert.Equal(t, "cli-standard", config["project_type"])
				assert.Equal(t, "cobra", config["framework"])
				assert.Equal(t, "slog", config["logger"])
			},
		},
		{
			name:           "returns 404 for invalid blueprint type",
			blueprintType:  "invalid-blueprint",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupBlueprintTestRouter()
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/v1/blueprints/%s/defaults", tt.blueprintType)
			req, _ := http.NewRequest("GET", url, nil)
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if tt.expectedStatus == http.StatusOK && tt.validateFunc != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err, "Response should be valid JSON")
				tt.validateFunc(t, response)
			}
		})
	}
}

func TestBlueprintHandler_GetBlueprintOptions(t *testing.T) {
	tests := []struct {
		name           string
		blueprintType  string
		expectedStatus int
		validateFunc   func(t *testing.T, response map[string]interface{})
	}{
		{
			name:           "get options for web-api blueprint",
			blueprintType:  "web-api-standard",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response map[string]interface{}) {
				options, ok := response["options"].(map[string]interface{})
				require.True(t, ok, "Response should contain options object")
				
				// Verify common options
				assert.Contains(t, options, "go_versions")
				assert.Contains(t, options, "loggers")
				
				// Verify web-api specific options
				assert.Contains(t, options, "frameworks")
				assert.Contains(t, options, "architectures")
				assert.Contains(t, options, "databases")
				
				frameworks, ok := options["frameworks"].([]interface{})
				require.True(t, ok, "frameworks should be an array")
				assert.Contains(t, frameworks, "gin")
				assert.Contains(t, frameworks, "echo")
			},
		},
		{
			name:           "get options for grpc blueprint",
			blueprintType:  "grpc-gateway",
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response map[string]interface{}) {
				options, ok := response["options"].(map[string]interface{})
				require.True(t, ok, "Response should contain options object")
				
				// Verify gRPC specific options
				assert.Contains(t, options, "architectures")
				assert.Contains(t, options, "features")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupBlueprintTestRouter()
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/v1/blueprints/%s/options", tt.blueprintType)
			req, _ := http.NewRequest("GET", url, nil)
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if tt.expectedStatus == http.StatusOK && tt.validateFunc != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err, "Response should be valid JSON")
				tt.validateFunc(t, response)
			}
		})
	}
}

func TestBlueprintHandler_ValidateBlueprintConfig(t *testing.T) {
	tests := []struct {
		name           string
		blueprintType  string
		requestBody    models.ValidateConfigRequest
		expectedStatus int
		validateFunc   func(t *testing.T, response *models.ValidateConfigResponse)
	}{
		{
			name:          "valid configuration for web-api blueprint",
			blueprintType: "web-api-standard",
			requestBody: models.ValidateConfigRequest{
				Config: models.ProjectConfig{
					ProjectName: "my-api",
					ModuleURL:   "github.com/user/my-api",
					GoVersion:   "1.21",
					ProjectType: "web-api-standard",
					Framework:   "gin",
					Logger:      "slog",
				},
			},
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *models.ValidateConfigResponse) {
				assert.True(t, response.Valid, "Valid configuration should pass validation")
				assert.Empty(t, response.Errors, "Valid configuration should have no errors")
			},
		},
		{
			name:          "invalid configuration missing required fields",
			blueprintType: "web-api-standard",
			requestBody: models.ValidateConfigRequest{
				Config: models.ProjectConfig{
					ProjectName: "", // Missing required field
					ModuleURL:   "",
					GoVersion:   "",
				},
			},
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *models.ValidateConfigResponse) {
				assert.False(t, response.Valid, "Invalid configuration should fail validation")
				assert.NotEmpty(t, response.Errors, "Invalid configuration should have errors")
				
				// Check for specific validation errors
				errorFields := make(map[string]bool)
				for _, err := range response.Errors {
					errorFields[err.Field] = true
				}
				
				assert.True(t, errorFields["project_name"], "Should have project_name error")
				assert.True(t, errorFields["module_url"], "Should have module_url error")
				assert.True(t, errorFields["go_version"], "Should have go_version error")
			},
		},
		{
			name:          "web-api blueprint missing framework",
			blueprintType: "web-api-standard",
			requestBody: models.ValidateConfigRequest{
				Config: models.ProjectConfig{
					ProjectName: "my-api",
					ModuleURL:   "github.com/user/my-api",
					GoVersion:   "1.21",
					ProjectType: "web-api-standard",
					Framework:   "", // Missing required framework
					Logger:      "slog",
				},
			},
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *models.ValidateConfigResponse) {
				assert.False(t, response.Valid, "Web API without framework should fail validation")
				
				// Should have framework error
				frameworkError := false
				for _, err := range response.Errors {
					if err.Field == "framework" {
						frameworkError = true
						assert.Equal(t, "error", err.Severity)
						break
					}
				}
				assert.True(t, frameworkError, "Should have framework validation error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupBlueprintTestRouter()
			w := httptest.NewRecorder()
			
			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err, "Should be able to marshal request body")
			
			url := fmt.Sprintf("/api/v1/blueprints/%s/validate", tt.blueprintType)
			req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if tt.expectedStatus == http.StatusOK && tt.validateFunc != nil {
				var response models.ValidateConfigResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err, "Response should be valid JSON")
				tt.validateFunc(t, &response)
			}
		})
	}
}

func TestBlueprintHandler_ErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		body           interface{}
		expectedStatus int
		expectedError  string
		expectedCode   string
	}{
		{
			name:           "invalid JSON in validation request",
			method:         "POST",
			url:            "/api/v1/blueprints/web-api-standard/validate",
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request format",
			expectedCode:   "INVALID_REQUEST",
		},
		{
			name:           "blueprint not found in validation",
			method:         "POST",
			url:            "/api/v1/blueprints/invalid-type/validate",
			body: models.ValidateConfigRequest{
				Config: models.ProjectConfig{
					ProjectName: "test",
					ModuleURL:   "github.com/user/test",
					GoVersion:   "1.21",
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Blueprint not found",
			expectedCode:   "BLUEPRINT_NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupBlueprintTestRouter()
			w := httptest.NewRecorder()
			
			var body []byte
			var err error
			if tt.body != nil {
				if str, ok := tt.body.(string); ok {
					body = []byte(str)
				} else {
					body, err = json.Marshal(tt.body)
					require.NoError(t, err, "Should be able to marshal request body")
				}
			}
			
			req, _ := http.NewRequest(tt.method, tt.url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tt.expectedStatus, w.Code)
			
			var errorResponse map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &errorResponse)
			require.NoError(t, err, "Error response should be valid JSON")
			
			assert.Equal(t, tt.expectedError, errorResponse["error"])
			assert.Equal(t, tt.expectedCode, errorResponse["code"])
		})
	}
}

// Benchmark tests for performance validation
func BenchmarkBlueprintHandler_ListBlueprints(b *testing.B) {
	router := setupBlueprintTestRouter()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/blueprints", nil)
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}

func BenchmarkBlueprintHandler_GetBlueprint(b *testing.B) {
	router := setupBlueprintTestRouter()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/blueprints/cli-simple", nil)
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
}