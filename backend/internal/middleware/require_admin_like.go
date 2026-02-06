// middleware/require_admin_like.go
package middleware

import (
	"ahu-backend/internal/domain"

	"github.com/gin-gonic/gin"
)

func RequireAdminLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		if !domain.IsAdminLike(role) {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
