package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/francknouama/go-starter/internal/templates"
)

type SystemHandler struct {
	registry  *templates.Registry
	startTime time.Time
	stats     *SystemStats
}

type SystemStats struct {
	ProjectsGenerated int64     `json:"projects_generated"`
	LastGeneration    time.Time `json:"last_generation,omitempty"`
	TotalUptime       string    `json:"total_uptime"`
}

type VersionInfo struct {
	Version    string `json:"version"`
	BuildTime  string `json:"build_time"`
	GoVersion  string `json:"go_version"`
	GitCommit  string `json:"git_commit,omitempty"`
}

type SystemInfo struct {
	Version        VersionInfo  `json:"version"`
	Runtime        RuntimeInfo  `json:"runtime"`
	Stats          SystemStats  `json:"stats"`
	Blueprints     BlueprintStats `json:"blueprints"`
}

type RuntimeInfo struct {
	GoVersion    string `json:"go_version"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
}

type BlueprintStats struct {
	Total            int            `json:"total"`
	ProductionReady  int            `json:"production_ready"`
	ByType          map[string]int  `json:"by_type"`
	ByComplexity    map[string]int  `json:"by_complexity"`
}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{
		registry:  templates.NewRegistry(),
		startTime: time.Now(),
		stats: &SystemStats{
			ProjectsGenerated: 0,
		},
	}
}

// GetVersion returns version information
func (h *SystemHandler) GetVersion(c *gin.Context) {
	versionInfo := VersionInfo{
		Version:   "v1.0.0", // TODO: Get from build flags
		BuildTime: "2025-08-29T00:00:00Z", // TODO: Get from build flags  
		GoVersion: runtime.Version(),
		GitCommit: "unknown", // TODO: Get from build flags
	}

	c.JSON(http.StatusOK, versionInfo)
}

// GetSystemStats returns system statistics
func (h *SystemHandler) GetSystemStats(c *gin.Context) {
	templates := h.registry.List()

	// Calculate blueprint statistics
	blueprintStats := BlueprintStats{
		Total:           len(templates),
		ProductionReady: len(templates), // All 12 blueprints are production ready
		ByType:          make(map[string]int),
		ByComplexity:    make(map[string]int),
	}

	for _, template := range templates {
		blueprintStats.ByType[template.Type]++
		complexity := getBlueprintComplexity(template.ID)
		blueprintStats.ByComplexity[complexity]++
	}

	// Update stats uptime
	h.stats.TotalUptime = time.Since(h.startTime).String()

	systemInfo := SystemInfo{
		Version: VersionInfo{
			Version:   "v1.0.0",
			BuildTime: "2025-08-29T00:00:00Z",
			GoVersion: runtime.Version(),
			GitCommit: "unknown",
		},
		Runtime: RuntimeInfo{
			GoVersion:    runtime.Version(),
			OS:           runtime.GOOS,
			Architecture: runtime.GOARCH,
			NumCPU:       runtime.NumCPU(),
			NumGoroutine: runtime.NumGoroutine(),
		},
		Stats:      *h.stats,
		Blueprints: blueprintStats,
	}

	c.JSON(http.StatusOK, systemInfo)
}

// IncrementProjectCount increments the project generation counter
func (h *SystemHandler) IncrementProjectCount() {
	h.stats.ProjectsGenerated++
	h.stats.LastGeneration = time.Now()
}

// getBlueprintComplexity determines the complexity level based on blueprint name
func getBlueprintComplexity(name string) string {
	switch name {
	case "cli-simple":
		return "simple"
	case "cli-standard", "library-standard", "lambda-standard", "lambda-proxy":
		return "standard"
	case "web-api-standard", "web-api-echo", "web-api-fiber", "monolith":
		return "intermediate"
	case "web-api-clean", "web-api-ddd", "microservice-standard":
		return "advanced"
	case "web-api-hexagonal", "grpc-gateway":
		return "expert"
	default:
		return "standard"
	}
}