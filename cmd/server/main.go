package main

import (
	"marketplace/internal/db"
	"marketplace/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	gormDB := db.Connect()

	r := gin.Default()
	handlers.RegisterRoutes(r, gormDB)

	r.Run(":8080")
}
