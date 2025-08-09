package handlers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"marketplace/internal/model"
)

func createOffer(db *gorm.DB) gin.HandlerFunc {
	type req struct {
		Title, Description string
		Price              float64
	}
	return func(c *gin.Context) {
		uid := c.GetUint("user_id")
		role := c.GetString("role")
		if role != "customer" {
			c.JSON(403, gin.H{"error": "only customers can create offers"})
			return
		}
		var r req
		if err := c.BindJSON(&r); err != nil {
			c.JSON(400, gin.H{"error": "bad json"})
			return
		}
		o := model.Offer{CustomerID: uid, Title: r.Title, Description: r.Description, Price: r.Price}
		if err := db.Create(&o).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, o)
	}
}

func listOffers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var list []model.Offer
		db.Find(&list)
		c.JSON(200, list)
	}
}

func createService(db *gorm.DB) gin.HandlerFunc {
	type req struct {
		Title, Description string
		Price              float64
	}
	return func(c *gin.Context) {
		uid := c.GetUint("user_id")
		role := c.GetString("role")
		if role != "performer" {
			c.JSON(403, gin.H{"error": "only performers can create services"})
			return
		}
		var r req
		if err := c.BindJSON(&r); err != nil {
			c.JSON(400, gin.H{"error": "bad json"})
			return
		}
		s := model.Service{PerformerID: uid, Title: r.Title, Description: r.Description, Price: r.Price}
		if err := db.Create(&s).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, s)
	}
}

func listServices(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var list []model.Service
		db.Find(&list)
		c.JSON(200, list)
	}
}

func addFavorite(db *gorm.DB) gin.HandlerFunc {
	type req struct{ ServiceID uint }
	return func(c *gin.Context) {
		uid := c.GetUint("user_id")
		role := c.GetString("role")
		if role != "customer" {
			c.JSON(403, gin.H{"error": "only customers can favorite"})
			return
		}
		var r req
		if err := c.BindJSON(&r); err != nil {
			c.JSON(400, gin.H{"error": "bad json"})
			return
		}

		var service model.Service
		if err := db.First(&service, r.ServiceID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(404, gin.H{"error": "service not found"})
			} else {
				c.JSON(500, gin.H{"error": err.Error()})
			}
			return
		}

		fav := model.Favorite{
			CustomerID: uid,
			ServiceID:  r.ServiceID}
		if err := db.Create(&fav).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, fav)
	}
}

func listFavorites(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetUint("user_id")
		var favs []model.Favorite
		db.Where("customer_id = ?", uid).Find(&favs)
		c.JSON(200, favs)
	}
}
