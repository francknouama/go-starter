package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/verify/zap/internal/handlers"
	"github.com/verify/zap/internal/logger"
	"github.com/verify/zap/internal/middleware"
)

type APITestSuite struct {
	suite.Suite
	router    *gin.Engine
	logger    logger.Logger
}

func (suite *APITestSuite) SetupSuite() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize logger for tests
	loggerFactory := logger.NewFactory()
	var err error
	suite.logger, err = loggerFactory.CreateFromProjectConfig("zap", "error", "text", false)
	suite.Require().NoError(err)

	// Setup router
	suite.router = gin.New()
	
	// Add middleware
	suite.router.Use(middleware.CORS())
	suite.router.Use(middleware.Logger(suite.logger))
	suite.router.Use(middleware.Recovery(suite.logger))

	// Setup routes
	suite.setupRoutes()
}

func (suite *APITestSuite) TearDownSuite() {
}

func (suite *APITestSuite) setupRoutes() {
	// Health routes
	suite.router.GET("/health", handlers.HealthCheck)
	suite.router.GET("/ready", handlers.ReadinessCheck)
}

// Test health endpoints
func (suite *APITestSuite) TestHealthEndpoints() {
	tests := []struct {
		name           string
		endpoint       string
		expectedStatus int
		expectedValue  string
	}{
		{
			name:           "Health check should return 200",
			endpoint:       "/health",
			expectedStatus: http.StatusOK,
			expectedValue:  "healthy",
		},
		{
			name:           "Ready check should return 200",
			endpoint:       "/ready",
			expectedStatus: http.StatusOK,
			expectedValue:  "ready",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			req, _ := http.NewRequest("GET", tt.endpoint, nil)
			w := httptest.NewRecorder()
			suite.router.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)
			
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)
			assert.Equal(suite.T(), tt.expectedValue, response["status"])
		})
	}
}

// Run the test suite
func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}