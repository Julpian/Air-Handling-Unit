package middleware

import "github.com/gin-gonic/gin"

type AuthUser struct {
	ID   string
	Role string
}

func GetUser(c *gin.Context) *AuthUser {
	return &AuthUser{
		ID:   c.GetString("user_id"),
		Role: c.GetString("role"),
	}
}
