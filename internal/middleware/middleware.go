package middleware

import (
	"marketplace/internal/auth"
	errors2 "marketplace/internal/error"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	const bearerPrefix = "Bearer "

	token := strings.TrimSpace(c.GetHeader("Authorization"))
	if token == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors2.ErrNoAuth})
		return
	}

	if !strings.HasPrefix(token, bearerPrefix) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors2.ErrBadToken})
		return
	}

	token = strings.TrimPrefix(token, bearerPrefix)

	claims, err := auth.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors2.ErrBadToken})
		return
	}

	c.Set("user_id", claims.UserID)
	c.Set("role", claims.Role)
	c.Next()
}
