package handlers

import (
	"marketplace/internal/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.Engine, pool *pgxpool.Pool) {
	r.POST("/auth/register", register(pool))
	r.POST("/auth/login", login(pool))

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

	authG.POST("/offers", createOffer(pool))
	authG.GET("/offers", listOffers(pool))

	authG.POST("/services", createService(pool))
	authG.GET("/services", listServices(pool))

	authG.POST("/favorites", addFavorite(pool))
	authG.GET("/favorites", listFavorites(pool))
}
