package main

import (
	"marketplace/internal/db"
	"marketplace/internal/logger"
	"marketplace/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	pool := db.Connect()

	r := gin.Default()
	router.RegisterRoutes(r, pool)
	logger := logger.NewLogger(pool)

	r.Run(":8080")
	logger.Info("Сервер стартовал на порту 8080")
}
