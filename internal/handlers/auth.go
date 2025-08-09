package handlers

import (
	"net/mail"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"marketplace/internal/auth"
	"marketplace/internal/model"
)

func register(db *gorm.DB) gin.HandlerFunc {
	type req struct{ Email, Password, Role, Name string }
	return func(c *gin.Context) {
		var r req
		if err := c.BindJSON(&r); err != nil {
			c.JSON(400, gin.H{"error": "bad json"})
			return
		}

		if _, err := mail.ParseAddress(r.Email); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if r.Role != "customer" && r.Role != "performer" {
			c.JSON(400, gin.H{"error": "bad role"})
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
		u := model.User{Email: r.Email,
			PasswordHash: string(hash),
			Role:         r.Role,
			Name:         r.Name}

		if err := db.Create(&u).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		token, _ := auth.GenerateToken(u.ID, u.Role)
		c.JSON(200, gin.H{"token": token})
	}
}

func login(db *gorm.DB) gin.HandlerFunc {
	type req struct {
		Email,
		Password string
	}
	return func(c *gin.Context) {
		var r req
		if err := c.BindJSON(&r); err != nil {
			c.JSON(400, gin.H{"error": "bad json"})
			return
		}

		var u model.User
		if err := db.Where("email = ?", r.Email).First(&u).Error; err != nil {
			c.JSON(401, gin.H{"error": "no user"})
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(r.Password)) != nil {
			c.JSON(401, gin.H{"error": "bad creds"})
			return
		}
		token, _ := auth.GenerateToken(u.ID, u.Role)
		c.JSON(200, gin.H{"token": token})
	}
}
