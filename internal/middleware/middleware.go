package middleware

import (
	"marketplace/internal/auth"
	errors2 "marketplace/internal/error"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors2.ErrNoAuth})
		return
	}
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	claims, err := auth.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors2.ErrBadToken})
		return
	}

	c.Set("user_id", claims.UserID)
	c.Set("role", claims.Role)
	c.Next()
}
