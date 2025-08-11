package handlers

import (
	"marketplace/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.Engine, pool *pgxpool.Pool) {
	r.POST("/auth/register", register(pool))
	r.POST("/auth/login", login(pool))

	authG := r.Group("/api")
	authG.Use(middleware.AuthMiddleware)

	authG.POST("/offers", createOffer(pool))
	authG.GET("/offers", listOffers(pool))
	authG.PATCH("/offers/", updateOffer(pool))
	authG.DELETE("/offers", deleteOffer(pool))

	authG.POST("/services", createService(pool))
	authG.GET("/services", listServices(pool))
	authG.PATCH("/services/", updateService(pool))
	authG.DELETE("/services", deleteService(pool))

	authG.POST("/favorites", addFavorite(pool))
	authG.GET("/favorites", listFavorites(pool))
	authG.DELETE("/favorites", deleteFavorite(pool))
}
