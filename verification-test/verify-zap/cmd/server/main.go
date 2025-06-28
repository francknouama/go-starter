package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/verify/zap/internal/config"
	"github.com/verify/zap/internal/handlers"
	"github.com/verify/zap/internal/logger"
	"github.com/verify/zap/internal/middleware"
)

func main() {
	// Load configuration first
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger factory
	loggerFactory := logger.NewFactory()
	
	// Create logger from configuration
	appLogger, err := loggerFactory.CreateFromProjectConfig(
		"zap",
		cfg.Logging.Level,
		cfg.Logging.Format,
		cfg.Logging.Structured,
	)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	appLogger.InfoWith("Application starting", logger.Fields{
		"logger": "zap",
		"environment": cfg.Environment,
	})

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.New()

	// Add middleware
	router.Use(middleware.Logger(appLogger))
	router.Use(middleware.Recovery(appLogger))
	router.Use(middleware.CORS())

	// Health check routes
	router.GET("/health", handlers.HealthCheck)
	router.GET("/ready", handlers.ReadinessCheck)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Add your API routes here
		_ = v1 // Placeholder to avoid unused variable error
	}

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	// Start server in a goroutine
	go func() {
		appLogger.InfoWith("Starting server", logger.Fields{
			"address": server.Addr,
			"environment": cfg.Environment,
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.ErrorWith("Failed to start server", logger.Fields{"error": err})
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		appLogger.ErrorWith("Server forced to shutdown", logger.Fields{"error": err})
		os.Exit(1)
	}

	appLogger.Info("Server exited")
}