package middleware

import "github.com/gin-gonic/gin"

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement panic recovery
		c.Next()
	}
}
