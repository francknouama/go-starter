package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version,omitempty"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// HealthCheck handles GET /health
// @Summary Health check endpoint
// @Description Returns the health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   "1.0.0", // This should come from build info
	}

	c.JSON(http.StatusOK, response)
}

// ReadinessCheck handles GET /ready
// @Summary Readiness check endpoint
// @Description Returns whether the service is ready to accept requests
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Failure 503 {object} HealthResponse
// @Router /ready [get]
func ReadinessCheck(c *gin.Context) {
	checks := make(map[string]string)
	allHealthy := true

	// Add more checks as needed (Redis, external APIs, etc.)

	status := "ready"
	httpStatus := http.StatusOK

	if !allHealthy {
		status = "not ready"
		httpStatus = http.StatusServiceUnavailable
	}

	response := HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Checks:    checks,
	}

	c.JSON(httpStatus, response)
}