package middleware

import "github.com/gin-gonic/gin"

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement request logging
		c.Next()
	}
}
