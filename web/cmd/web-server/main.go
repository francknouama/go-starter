package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/francknouama/go-starter/web/internal/web/handlers"
	"github.com/francknouama/go-starter/web/internal/web/middleware"
	"github.com/francknouama/go-starter/web/internal/web/websocket"
)

//go:embed dist
var distFS embed.FS

func main() {
	// Set Gin mode
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.ErrorHandler())

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",  // React dev server
		"http://localhost:5173",  // Vite dev server
		"http://localhost:8080",  // Production server
	}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	corsConfig.AllowCredentials = true
	
	// Allow additional origins from environment
	if additionalOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); additionalOrigins != "" {
		corsConfig.AllowOrigins = append(corsConfig.AllowOrigins, additionalOrigins)
	}
	
	router.Use(cors.New(corsConfig))

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// API routes
	apiV1 := router.Group("/api/v1")
	{
		// Health and system endpoints
		apiV1.GET("/health", handlers.HealthCheck)
		apiV1.GET("/health/simple", handlers.HealthCheckSimple)
		apiV1.GET("/health/readiness", handlers.ReadinessCheck)
		apiV1.GET("/health/liveness", handlers.LivenessCheck)
		apiV1.GET("/metrics", handlers.MetricsHandler)
		apiV1.GET("/status", handlers.StatusHandler)
		apiV1.GET("/version", handlers.VersionHandler)
		
		// Configuration endpoints
		apiV1.GET("/config", handlers.GetConfig)
		apiV1.GET("/config/types/:type", handlers.GetProjectTypeDetails)
		apiV1.GET("/config/frameworks", handlers.GetFrameworksForType)
		apiV1.GET("/config/architectures", handlers.GetArchitecturesForType)
		
		// Blueprint endpoints
		apiV1.GET("/blueprints", handlers.ListBlueprints)
		apiV1.GET("/blueprints/:id", handlers.GetBlueprint)
		apiV1.GET("/blueprints/category/:category", handlers.GetBlueprintsByCategory)
		apiV1.POST("/blueprints/:id/validate", handlers.GetBlueprintValidation)
		
		// Project generation endpoints
		apiV1.POST("/preview", handlers.PreviewProject)
		apiV1.POST("/generate", handlers.GenerateProject)
		apiV1.POST("/generate/download", handlers.GenerateAndDownloadProject)
		
		// Download endpoints
		apiV1.GET("/download/:token", handlers.DownloadZip)
		apiV1.GET("/download/:token/status", handlers.GetDownloadStatus)
		
		// WebSocket endpoint
		apiV1.GET("/ws", func(c *gin.Context) {
			websocket.HandleWebSocket(hub, c.Writer, c.Request)
		})
		
		// WebSocket management endpoints (for debugging/admin)
		wsHandler := websocket.NewWSHandler(hub)
		hub.SetWSHandler(wsHandler) // Connect the handler to the hub
		apiV1.GET("/ws/info", wsHandler.ServeWebSocketInfo)
		apiV1.GET("/ws/clients", wsHandler.GetConnectedClients)
		apiV1.GET("/ws/stats", wsHandler.GetHubStats)
		apiV1.GET("/ws/test", wsHandler.WebSocketTestPage)
	}

	// Legacy API routes (for backward compatibility)
	api := router.Group("/api")
	{
		api.GET("/health", handlers.HealthCheck)
		api.GET("/config", handlers.GetConfig)
		api.GET("/blueprints", handlers.ListBlueprints)
		api.POST("/preview", handlers.PreviewProject)
		api.POST("/generate", handlers.GenerateProject)
		api.GET("/download/:token", handlers.DownloadZip)
	}

	// Serve static files from embedded dist directory
	distSubFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		log.Fatal("Failed to create dist sub filesystem:", err)
	}

	// Serve React app for all non-API routes
	router.NoRoute(func(c *gin.Context) {
		// Try to serve the file from dist
		path := c.Request.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		file, err := distSubFS.Open(path[1:]) // Remove leading slash
		if err != nil {
			// File not found, serve index.html for client-side routing
			indexFile, indexErr := distSubFS.Open("index.html")
			if indexErr != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Frontend assets not found"})
				return
			}
			defer indexFile.Close()

			stat, statErr := indexFile.Stat()
			if statErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stat index.html"})
				return
			}

			c.Header("Content-Type", "text/html")
			http.ServeContent(c.Writer, c.Request, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stat file"})
			return
		}

		// Set appropriate content type based on file extension
		contentType := getContentType(path)
		if contentType != "" {
			c.Header("Content-Type", contentType)
		}

		// Add cache headers for static assets
		if path != "/index.html" {
			c.Header("Cache-Control", "public, max-age=31536000") // 1 year
		}

		http.ServeContent(c.Writer, c.Request, path, stat.ModTime(), file.(io.ReadSeeker))
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("🚀 Web server starting on http://localhost:%s\n", port)
		fmt.Printf("📚 API documentation available at http://localhost:%s/api/health\n", port)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("✅ Server exited gracefully")
}

// getContentType returns the appropriate content type for common file extensions
func getContentType(path string) string {
	switch {
	case len(path) >= 4 && path[len(path)-4:] == ".css":
		return "text/css"
	case len(path) >= 3 && path[len(path)-3:] == ".js":
		return "application/javascript"
	case len(path) >= 5 && path[len(path)-5:] == ".json":
		return "application/json"
	case len(path) >= 4 && path[len(path)-4:] == ".svg":
		return "image/svg+xml"
	case len(path) >= 4 && path[len(path)-4:] == ".png":
		return "image/png"
	case len(path) >= 4 && path[len(path)-4:] == ".jpg" || len(path) >= 5 && path[len(path)-5:] == ".jpeg":
		return "image/jpeg"
	case len(path) >= 4 && path[len(path)-4:] == ".ico":
		return "image/x-icon"
	case len(path) >= 5 && path[len(path)-5:] == ".html":
		return "text/html"
	default:
		return ""
	}
}