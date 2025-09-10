package router

import (
	"marketplace/internal/handlers"
	"marketplace/internal/logger"
	"marketplace/internal/middleware"
	"marketplace/internal/notifications"
	repository "marketplace/internal/repo"
	"marketplace/internal/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.Engine, pool *pgxpool.Pool) {
	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Accept", "Origin"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Logger
	logger := logger.NewLogger(pool)

	//notifications
	notifications := notifications.NewEmailNotifier()

	// Auth
	auth := r.Group("/auth")
	authRepo := repository.NewAuthRepo(pool)
	authService := service.NewAuthService(authRepo, notifications)
	authHandler := handlers.NewAuthHandler(authService, logger)
	auth.POST("/register", authHandler.Register())
	auth.POST("/login", authHandler.Login())

	// Handlers group
	v1 := r.Group("/api")
	v1.Use(middleware.AuthMiddleware)

	// Offers
	offerRepo := repository.NewOfferRepository(pool)
	offerService := service.NewOfferService(offerRepo)
	offerHandler := handlers.NewOfferHandler(offerService, logger)
	v1.POST("/offers", offerHandler.CreateOffer())
	v1.PATCH("/offers", offerHandler.UpdateOffer())
	v1.DELETE("/offers", offerHandler.DeleteOffer())
	v1.GET("/offers", offerHandler.ListOffers())

	// Services
	serviceRepo := repository.NewServiceRepository(pool)
	serviceService := service.NewServiceService(serviceRepo)
	serviceHandler := handlers.NewServiceHandler(serviceService, logger)
	v1.POST("/services", serviceHandler.CreateService())
	v1.PATCH("/services", serviceHandler.UpdateService())
	v1.DELETE("/services", serviceHandler.DeleteService())
	v1.GET("/services", serviceHandler.ListServices())

	// Favorites
	favoriteRepo := repository.NewFavoriteRepository(pool)
	favoriteService := service.NewFavoriteService(favoriteRepo)
	favoriteHandler := handlers.NewFavoriteHandler(favoriteService, logger)
	v1.POST("/favorites", favoriteHandler.AddFavorite())
	v1.DELETE("/favorites", favoriteHandler.DeleteFavorite())
	v1.GET("/favorites", favoriteHandler.ListFavorites())
}
