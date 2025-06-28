package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/verify/zap/internal/logger"
)

// Logger creates a logger middleware using our logger interface
func Logger(log logger.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code and response size
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		// Build the full path
		if raw != "" {
			path = path + "?" + raw
		}

		// Prepare log fields
		fields := logger.Fields{
			"method":     c.Request.Method,
			"path":       path,
			"status":     statusCode,
			"latency_ms": latency.Milliseconds(),
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"body_size":  bodySize,
		}

		// Add error field if there's an error
		if len(c.Errors) > 0 {
			fields["error"] = c.Errors.String()
		}

		// Choose log level based on status code
		switch {
		case statusCode >= 400 && statusCode < 500:
			log.WarnWith("HTTP request", fields)
		case statusCode >= 500:
			log.ErrorWith("HTTP request", fields)
		default:
			log.InfoWith("HTTP request", fields)
		}
	})
}