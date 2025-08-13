package handlers

import (
	"marketplace/internal/auth"
	repository "marketplace/internal/repo"
	"marketplace/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo *repository.AuthRepo
}

func NewAuthHandler(repo *repository.AuthRepo) *AuthHandler {
	return &AuthHandler{repo: repo}
}

func (h *AuthHandler) Register() gin.HandlerFunc {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
		Name     string `json:"name"`
	}

	return func(c *gin.Context) {
		r, err := utils.BindJSON[req](c)
		if err != nil {
			c.JSON(400, gin.H{"error": "Неверный формат запроса: " + err.Error()})
			return
		}

		if err := utils.ValidateEmail(r.Email); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := utils.ValidateName(r.Name); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := utils.CheckRole(r.Role); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		userID, err := h.repo.RegisterUser(c.Request.Context(), r.Email, r.Password, r.Role, r.Name)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		token, err := auth.GenerateToken(userID, r.Role)
		if err != nil {
			c.JSON(500, gin.H{"error": "token generation failed"})
			return
		}

		c.JSON(200, gin.H{"token": token})
	}
}

func (h *AuthHandler) Login() gin.HandlerFunc {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(c *gin.Context) {
		r, err := utils.BindJSON[req](c)
		if err != nil {
			c.JSON(400, gin.H{"error": "Неверный формат запроса: " + err.Error()})
			return
		}
		user, err := h.repo.GetUserByEmail(c.Request.Context(), r.Email)
		if err != nil {
			c.JSON(401, gin.H{"error": "Неверный логин или пароль"})
			return
		}

		if err := bcrypt.CompareHashAndPassword(
			[]byte(user.PasswordHash),
			[]byte(r.Password),
		); err != nil {
			c.JSON(401, gin.H{"error": "Неверный логин или пароль"})
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Role)
		if err != nil {
			c.JSON(500, gin.H{"error": "Ошибка генерации токена"})
			return
		}

		c.JSON(200, gin.H{"token": token})
	}
}
