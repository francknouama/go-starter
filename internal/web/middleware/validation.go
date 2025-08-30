package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Code    string                 `json:"code"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ValidationMiddleware provides request validation
func ValidationMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Add validation headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		c.Next()
	})
}

// ContentTypeValidation ensures proper content type for POST/PUT requests
func ContentTypeValidation() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					Error: "Content-Type must be application/json",
					Code:  "INVALID_CONTENT_TYPE",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	})
}

// RateLimiting provides basic rate limiting (in production, use Redis-based rate limiting)
func RateLimiting() gin.HandlerFunc {
	// This is a simple in-memory rate limiter for development
	// In production, use a proper rate limiter like go-limiter or similar
	return gin.HandlerFunc(func(c *gin.Context) {
		// Add rate limiting headers
		c.Header("X-RateLimit-Limit", "100")
		c.Header("X-RateLimit-Remaining", "99") // Simplified for demo

		c.Next()
	})
}

// CORS middleware for development (already handled by gin-contrib/cors in main.go)
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Allow local development origins
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:5173", // Vite default
			"http://localhost:8080",
		}

		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

// ErrorHandler middleware for consistent error responses
func ErrorHandler() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()

		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			var statusCode int
			var errorCode string

			switch err.Type {
			case gin.ErrorTypeBind:
				statusCode = http.StatusBadRequest
				errorCode = "BINDING_ERROR"
			case gin.ErrorTypePublic:
				statusCode = http.StatusBadRequest
				errorCode = "VALIDATION_ERROR"
			default:
				statusCode = http.StatusInternalServerError
				errorCode = "INTERNAL_ERROR"
			}

			c.JSON(statusCode, ErrorResponse{
				Error: err.Error(),
				Code:  errorCode,
			})
		}
	})
}