package handlers

import (
	"marketplace/internal/middleware"
	repository "marketplace/internal/repo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.Engine, pool *pgxpool.Pool) {
	offerRepo := repository.NewOfferRepository(pool)
	serviceRepo := repository.NewServiceRepository(pool)

	offerHandler := NewOfferHandler(offerRepo)
	serviceHandler := NewServiceHandler(serviceRepo)

	r.POST("/auth/register", register(pool))
	r.POST("/auth/login", login(pool))

	authG := r.Group("/api")
	authG.Use(middleware.AuthMiddleware)

	authG.POST("/offers", offerHandler.CreateOffer())
	authG.PATCH("/offers/", offerHandler.UpdateOffer())
	authG.DELETE("/offers", offerHandler.DeleteOffer())
	authG.GET("/offers", offerHandler.ListOffers())

	authG.POST("/services", serviceHandler.CreateService())
	authG.PATCH("/services/", serviceHandler.UpdateService())
	authG.DELETE("/services", serviceHandler.DeleteService())
	authG.GET("/services", serviceHandler.ListServices())

	authG.POST("/favorites", addFavorite(pool))
	authG.GET("/favorites", listFavorites(pool))
	authG.DELETE("/favorites", deleteFavorite(pool))
}
