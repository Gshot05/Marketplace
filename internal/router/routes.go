package router

import (
	"marketplace/internal/handlers"
	"marketplace/internal/middleware"
	repository "marketplace/internal/repo"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.Engine, pool *pgxpool.Pool) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authRepo := repository.NewAuthRepo(pool)
	offerRepo := repository.NewOfferRepository(pool)
	serviceRepo := repository.NewServiceRepository(pool)
	favoriteRepo := repository.NewFavoriteRepository(pool)

	authHandler := handlers.NewAuthHandler(authRepo)
	offerHandler := handlers.NewOfferHandler(offerRepo)
	serviceHandler := handlers.NewServiceHandler(serviceRepo)
	favoriteHandler := handlers.NewFavoriteHandler(favoriteRepo)

	auth := r.Group("/auth")
	auth.POST("/register", authHandler.Register())
	auth.POST("/login", authHandler.Login())

	v1 := r.Group("/api")
	v1.Use(middleware.AuthMiddleware)
	v1.POST("/offers", offerHandler.CreateOffer())
	v1.PATCH("/offers", offerHandler.UpdateOffer())
	v1.DELETE("/offers", offerHandler.DeleteOffer())
	v1.GET("/offers", offerHandler.ListOffers())

	v1.POST("/services", serviceHandler.CreateService())
	v1.PATCH("/services", serviceHandler.UpdateService())
	v1.DELETE("/services", serviceHandler.DeleteService())
	v1.GET("/services", serviceHandler.ListServices())

	v1.POST("/favorites", favoriteHandler.AddFavorite())
	v1.DELETE("/favorites", favoriteHandler.DeleteFavorite())
	v1.GET("/favorites", favoriteHandler.ListFavorites())
}
