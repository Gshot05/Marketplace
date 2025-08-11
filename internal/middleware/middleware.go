package middleware

import (
	"log"
	"marketplace/internal/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no auth"})
		return
	}
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	claims, err := auth.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Delete потом
	log.Printf("AuthMiddleware: user_id=%d, role=%s", claims.UserID, claims.Role)

	c.Set("user_id", claims.UserID)
	c.Set("role", claims.Role)
	c.Next()
}
