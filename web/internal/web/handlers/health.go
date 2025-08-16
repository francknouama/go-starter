package handlers

import (
	"fmt"
	"net/http"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/francknouama/go-starter/web/internal/web/models"
)

var (
	startTime = time.Now()
	version   = "1.0.0" // This would normally be set during build
)

// HealthCheck returns the health status of the API
func HealthCheck(c *gin.Context) {
	// Get system information
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	buildInfo, _ := debug.ReadBuildInfo()
	
	response := &models.HealthResponse{
		Status:    "healthy",
		Version:   version,
		Timestamp: time.Now(),
		Services: map[string]string{
			"api":       "healthy",
			"generator": "healthy",
			"templates": "healthy",
			"websocket": "healthy",
		},
		System: &models.SystemInfo{
			GoVersion:      runtime.Version(),
			OS:             runtime.GOOS,
			Arch:           runtime.GOARCH,
			NumCPU:         runtime.NumCPU(),
			MemoryUsage:    int64(memStats.Alloc),
			GoroutineCount: runtime.NumGoroutine(),
			Uptime:         time.Since(startTime).String(),
		},
	}
	
	// Add build information if available
	if buildInfo != nil {
		if buildInfo.Main.Version != "" && buildInfo.Main.Version != "(devel)" {
			response.Version = buildInfo.Main.Version
		}
	}
	
	apiResponse := models.NewSuccessResponse(response)
	c.JSON(http.StatusOK, apiResponse)
}

// HealthCheckSimple returns a simple health check for load balancers
func HealthCheckSimple(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Unix(),
	})
}

// ReadinessCheck checks if the service is ready to serve requests
func ReadinessCheck(c *gin.Context) {
	// Check if critical services are available
	errors := []string{}
	
	// Check template registry (this would be a real check in production)
	if !checkTemplateRegistry() {
		errors = append(errors, "template registry unavailable")
	}
	
	// Check generator service
	if !checkGeneratorService() {
		errors = append(errors, "generator service unavailable")
	}
	
	if len(errors) > 0 {
		response := models.NewErrorResponse(
			"SERVICE_UNAVAILABLE",
			"Service not ready",
			fmt.Sprintf("Issues: %v", errors),
		)
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"time":   time.Now().Unix(),
	})
}

// LivenessCheck checks if the service is alive
func LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
		"time":   time.Now().Unix(),
		"uptime": time.Since(startTime).Seconds(),
	})
}

// MetricsHandler returns basic metrics about the service
func MetricsHandler(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	metrics := map[string]interface{}{
		"uptime_seconds":       time.Since(startTime).Seconds(),
		"memory_alloc_bytes":   memStats.Alloc,
		"memory_total_alloc_bytes": memStats.TotalAlloc,
		"memory_sys_bytes":     memStats.Sys,
		"memory_heap_alloc_bytes": memStats.HeapAlloc,
		"memory_heap_sys_bytes": memStats.HeapSys,
		"memory_heap_idle_bytes": memStats.HeapIdle,
		"memory_heap_inuse_bytes": memStats.HeapInuse,
		"memory_stack_inuse_bytes": memStats.StackInuse,
		"memory_stack_sys_bytes": memStats.StackSys,
		"goroutines_count":     runtime.NumGoroutine(),
		"gc_runs":             memStats.NumGC,
		"gc_pause_total_ns":   memStats.PauseTotalNs,
		"next_gc_bytes":       memStats.NextGC,
		"cpu_count":           runtime.NumCPU(),
		"version":             version,
		"go_version":          runtime.Version(),
		"os":                  runtime.GOOS,
		"arch":                runtime.GOARCH,
	}
	
	response := models.NewSuccessResponse(metrics)
	c.JSON(http.StatusOK, response)
}

// StatusHandler returns detailed status information
func StatusHandler(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	status := map[string]interface{}{
		"service": map[string]interface{}{
			"name":      "go-starter-web-api",
			"version":   version,
			"status":    "running",
			"startTime": startTime,
			"uptime":    time.Since(startTime).String(),
		},
		"system": map[string]interface{}{
			"go_version":      runtime.Version(),
			"os":              runtime.GOOS,
			"arch":            runtime.GOARCH,
			"cpu_count":       runtime.NumCPU(),
			"goroutine_count": runtime.NumGoroutine(),
			"memory": map[string]interface{}{
				"alloc_mb":      float64(memStats.Alloc) / 1024 / 1024,
				"total_alloc_mb": float64(memStats.TotalAlloc) / 1024 / 1024,
				"sys_mb":        float64(memStats.Sys) / 1024 / 1024,
				"heap_alloc_mb": float64(memStats.HeapAlloc) / 1024 / 1024,
				"heap_sys_mb":   float64(memStats.HeapSys) / 1024 / 1024,
			},
			"gc": map[string]interface{}{
				"num_gc":           memStats.NumGC,
				"pause_total_ms":   float64(memStats.PauseTotalNs) / 1000000,
				"next_gc_mb":       float64(memStats.NextGC) / 1024 / 1024,
			},
		},
		"dependencies": map[string]interface{}{
			"template_registry": checkTemplateRegistry(),
			"generator_service": checkGeneratorService(),
		},
	}
	
	response := models.NewSuccessResponse(status)
	c.JSON(http.StatusOK, response)
}

// VersionHandler returns version information
func VersionHandler(c *gin.Context) {
	buildInfo, _ := debug.ReadBuildInfo()
	
	versionInfo := map[string]interface{}{
		"version":    version,
		"go_version": runtime.Version(),
		"os":         runtime.GOOS,
		"arch":       runtime.GOARCH,
		"build_time": startTime.Format(time.RFC3339),
	}
	
	if buildInfo != nil {
		if buildInfo.Main.Version != "" && buildInfo.Main.Version != "(devel)" {
			versionInfo["version"] = buildInfo.Main.Version
		}
		
		// Add build settings
		settings := make(map[string]string)
		for _, setting := range buildInfo.Settings {
			settings[setting.Key] = setting.Value
		}
		versionInfo["build_settings"] = settings
	}
	
	response := models.NewSuccessResponse(versionInfo)
	c.JSON(http.StatusOK, response)
}

// checkTemplateRegistry checks if the template registry is available
func checkTemplateRegistry() bool {
	// This would be a real check against the template registry
	// For now, we'll assume it's always available
	return true
}

// checkGeneratorService checks if the generator service is available
func checkGeneratorService() bool {
	// This would be a real check against the generator service
	// For now, we'll assume it's always available
	return true
}