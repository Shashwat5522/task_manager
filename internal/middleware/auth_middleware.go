package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vedologic/task-manager/pkg/utils"
)

// AuthMiddleware returns a gin middleware for JWT authentication
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(token, jwtSecret)
		if err != nil {
			c.JSON(401, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}
