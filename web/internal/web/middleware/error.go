package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/francknouama/go-starter/web/internal/web/models"
)

// ErrorHandler is a middleware that handles errors and panics
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with stack trace
				stack := debug.Stack()
				fmt.Printf("PANIC: %v\nStack trace:\n%s\n", err, stack)

				// Return internal server error
				response := models.NewInternalErrorResponse("Internal server error occurred")
				c.JSON(http.StatusInternalServerError, response)
				c.Abort()
			}
		}()

		c.Next()

		// Handle errors that were set during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Determine status code based on error type
			statusCode := http.StatusInternalServerError
			var response *models.APIResponse

			switch err.Type {
			case gin.ErrorTypeBind:
				statusCode = http.StatusBadRequest
				response = models.NewValidationErrorResponse(err.Error())
			case gin.ErrorTypePublic:
				statusCode = http.StatusBadRequest
				response = models.NewErrorResponse("BAD_REQUEST", "Bad request", err.Error())
			case gin.ErrorTypeRender:
				statusCode = http.StatusInternalServerError
				response = models.NewInternalErrorResponse("Rendering error")
			default:
				response = models.NewInternalErrorResponse("Internal server error")
			}

			c.JSON(statusCode, response)
			return
		}

		// Handle 404 for API routes
		if c.Writer.Status() == http.StatusNotFound && 
		   (gin.Mode() != gin.TestMode && c.Request.URL.Path[:4] == "/api") {
			response := models.NewNotFoundErrorResponse("API endpoint not found")
			c.JSON(http.StatusNotFound, response)
		}
	}
}

// RequestSizeLimit middleware limits the size of request bodies
func RequestSizeLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			response := models.NewErrorResponse(
				"REQUEST_TOO_LARGE",
				"Request body too large",
				fmt.Sprintf("Maximum allowed size is %d bytes", maxBytes),
			)
			c.JSON(http.StatusRequestEntityTooLarge, response)
			c.Abort()
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}

// Timeout middleware sets a timeout for request processing
func Timeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// For now, just set a header - timeout handling can be implemented later
		c.Header("X-Request-Timeout", "30s")
		c.Next()
	}
}