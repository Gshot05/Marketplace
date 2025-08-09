package handlers

import (
	"marketplace/internal/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/auth/register", register(db))
	r.POST("/auth/login", login(db))

	authG := r.Group("/api")
	authG.Use(func(c *gin.Context) {
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
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	})

	authG.POST("/offers", createOffer(db))
	authG.GET("/offers", listOffers(db))

	authG.POST("/services", createService(db))
	authG.GET("/services", listServices(db))

	authG.POST("/favorites", addFavorite(db))
	authG.GET("/favorites", listFavorites(db))
}
