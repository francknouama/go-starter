package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	startTime time.Time
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// Health returns the health status of the service
func (h *HealthHandler) Health(c *gin.Context) {
	uptime := time.Since(h.startTime)
	
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    uptime.String(),
		"service":   "go-starter-web",
	})
}

// Version returns version and build information
func (h *HealthHandler) Version(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	c.JSON(http.StatusOK, gin.H{
		"version":     "1.0.0-dev",
		"go_version":  runtime.Version(),
		"build_time":  "development", // Would be set during build
		"git_commit":  "development", // Would be set during build
		"uptime":      time.Since(h.startTime).String(),
		"memory": gin.H{
			"alloc_mb":      memStats.Alloc / 1024 / 1024,
			"total_alloc_mb": memStats.TotalAlloc / 1024 / 1024,
			"sys_mb":        memStats.Sys / 1024 / 1024,
			"gc_runs":       memStats.NumGC,
		},
	})
}