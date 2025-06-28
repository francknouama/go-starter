package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/verify/web-api/internal/logger"
)

// Recovery creates a recovery middleware using our logger interface
func Recovery(log logger.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				log.ErrorWith("Panic recovered", logger.Fields{
					"error":      err,
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
					"ip":         c.ClientIP(),
					"user_agent": c.Request.UserAgent(),
					"stack":      string(debug.Stack()),
				})

				// Return internal server error
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
					"code":  "INTERNAL_ERROR",
				})
				c.Abort()
			}
		}()

		c.Next()
	})
}