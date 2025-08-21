package middleware

import (
	"marketplace/internal/auth"
	errors2 "marketplace/internal/error"
	"marketplace/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	token, err := utils.ValidateBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
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
