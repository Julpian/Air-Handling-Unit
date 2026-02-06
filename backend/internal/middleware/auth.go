package middleware

import (
	"strings"

	"ahu-backend/internal/security"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := security.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		// 🔥🔥🔥 INI YANG PALING PENTING
		c.Set("user_id", claims.UserID) // ✅ UUID
		c.Set("role", claims.Role)
		c.Set("is_active", claims.IsActive)

		c.Next()
	}
}
