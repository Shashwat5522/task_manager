package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware returns a gin middleware for panic recovery
func RecoveryMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error("Panic recovered",
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.Any("error", err),
				)
				c.JSON(500, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}
