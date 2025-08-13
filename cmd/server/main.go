package main

import (
	"marketplace/internal/db"
	"marketplace/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	pool := db.Connect()

	r := gin.Default()
	router.RegisterRoutes(r, pool)

	r.Run(":8080")
}
