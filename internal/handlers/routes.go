package handlers

import (
	"marketplace/internal/middleware"
	repository "marketplace/internal/repo"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.Engine, pool *pgxpool.Pool) {
	authRepo := repository.NewAuthRepo(pool)
	offerRepo := repository.NewOfferRepository(pool)
	serviceRepo := repository.NewServiceRepository(pool)
	favoriteRepo := repository.NewFavoriteRepository(pool)

	authHandler := NewAuthHandler(authRepo)
	offerHandler := NewOfferHandler(offerRepo)
	serviceHandler := NewServiceHandler(serviceRepo)
	favoriteHandler := NewFavoriteHandler(favoriteRepo)

	r.POST("/auth/register", authHandler.Register())
	r.POST("/auth/login", authHandler.Login())

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

	authG.POST("/favorites", favoriteHandler.AddFavorite())
	authG.DELETE("/favorites", favoriteHandler.DeleteFavorite())
	authG.GET("/favorites", favoriteHandler.ListFavorites())
}
