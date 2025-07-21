package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

// Logger returns a Gin middleware for structured logging
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		slog.Info("HTTP request",
			"method", param.Method,
			"path", param.Path,
			"status", param.StatusCode,
			"latency", param.Latency,
			"client_ip", param.ClientIP,
			"user_agent", param.Request.UserAgent(),
			"error", param.ErrorMessage,
		)
		return ""
	})
}