package main

import (
	"marketplace/internal/db"
	"marketplace/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	pool := db.Connect()

	r := gin.Default()
	handlers.RegisterRoutes(r, pool)

	r.Run(":8080")
}
