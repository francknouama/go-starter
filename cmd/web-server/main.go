package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/internal/web/handlers"
	"github.com/francknouama/go-starter/internal/web/middleware"
	"github.com/francknouama/go-starter/internal/web/websocket"
)

func main() {
	// Initialize the templates filesystem for development
	// Using os.DirFS to access blueprints from the filesystem
	templatesFS := os.DirFS("blueprints")
	templates.SetTemplatesFS(templatesFS)

	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Create Gin router
	router := gin.New()

	// CORS configuration for development
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // React dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))

	// Middleware
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestID())

	// Initialize handlers
	blueprintHandler := handlers.NewBlueprintHandler()
	generatorHandler := handlers.NewGeneratorHandler()
	healthHandler := handlers.NewHealthHandler()

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	wsHandler := handlers.NewWebSocketHandler(wsHub)

	// Routes
	v1 := router.Group("/api/v1")
	{
		// System endpoints
		v1.GET("/health", healthHandler.Health)
		v1.GET("/version", healthHandler.Version)

		// Blueprint endpoints
		v1.GET("/blueprints", blueprintHandler.ListBlueprints)
		v1.GET("/blueprints/:id", blueprintHandler.GetBlueprint)

		// Generation endpoints
		v1.POST("/validate", generatorHandler.ValidateConfig)
		v1.POST("/generate", generatorHandler.GenerateProject)
		v1.GET("/download/:id", generatorHandler.DownloadProject)
		v1.DELETE("/projects/:id", generatorHandler.CleanupProject)

		// WebSocket endpoints
		v1.GET("/ws/generate", wsHandler.HandleGenerateWS)
		v1.GET("/ws/preview", wsHandler.HandlePreviewWS)
	}

	// Serve static files (for development, serve from filesystem)
	// TODO: In production, this should use embedded files
	router.Static("/static", "web/dist")
	router.Static("/assets", "web/dist/assets")
	
	// Serve index.html for SPA routing
	router.NoRoute(func(c *gin.Context) {
		c.File("web/dist/index.html")
	})

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		slog.Info("Starting web server", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped")
}