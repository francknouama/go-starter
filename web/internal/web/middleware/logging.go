package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Logger returns a gin.HandlerFunc (middleware) that logs requests using the standard library logger
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// RequestID middleware adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID is already set in headers
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate new UUID for request ID
			requestID = uuid.New().String()
		}

		// Set request ID in context and response header
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// StructuredLogger returns a structured logging middleware
func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request ID
		requestID, exists := c.Get("RequestID")
		if !exists {
			requestID = "unknown"
		}

		// Prepare log fields
		fields := map[string]interface{}{
			"timestamp":  start.Format(time.RFC3339),
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       path,
			"query":      raw,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"status":     c.Writer.Status(),
			"latency":    latency.String(),
			"latency_ms": float64(latency.Nanoseconds()) / 1000000.0,
			"bytes_in":   c.Request.ContentLength,
			"bytes_out":  c.Writer.Size(),
		}

		// Add error information if present
		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		// Log based on status code
		statusCode := c.Writer.Status()
		message := fmt.Sprintf("%s %s - %d", c.Request.Method, path, statusCode)

		if statusCode >= 500 {
			fmt.Printf("ERROR: %s %+v\n", message, fields)
		} else if statusCode >= 400 {
			fmt.Printf("WARN: %s %+v\n", message, fields)
		} else {
			fmt.Printf("INFO: %s %+v\n", message, fields)
		}
	}
}

// APIMetrics middleware collects basic API metrics
func APIMetrics() gin.HandlerFunc {
	var (
		requestCount    = make(map[string]int64)
		responseTimeSum = make(map[string]time.Duration)
	)

	return func(c *gin.Context) {
		start := time.Now()
		
		c.Next()
		
		// Calculate metrics
		duration := time.Since(start)
		endpoint := fmt.Sprintf("%s %s", c.Request.Method, c.FullPath())
		
		// Update metrics (in a real application, you'd want to use proper atomic operations)
		requestCount[endpoint]++
		responseTimeSum[endpoint] += duration
		
		// Set metrics in response headers for debugging
		c.Header("X-Response-Time", duration.String())
		
		// Log metrics every 100 requests (simple example)
		if requestCount[endpoint]%100 == 0 {
			avgResponseTime := responseTimeSum[endpoint] / time.Duration(requestCount[endpoint])
			fmt.Printf("METRICS: %s - Count: %d, Avg Response Time: %s\n", 
				endpoint, requestCount[endpoint], avgResponseTime)
		}
	}
}

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one
func Recovery() gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(500, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(500)
	})
}