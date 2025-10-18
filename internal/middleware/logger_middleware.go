package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware returns a gin middleware for request logging
func LoggerMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.Info("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("duration", duration),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
		)
	}
}
