package handlers

import (
	"marketplace/internal/auth"
	"marketplace/internal/logger"
	"marketplace/internal/model"
	"marketplace/internal/service"
	"marketplace/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	s      service.IAuthService
	logger *logger.Logger
}

func NewAuthHandler(
	s service.IAuthService,
	logger *logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		s:      s,
		logger: logger}
}

func (h *AuthHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		r, err := utils.BindJSON[model.RegisterReq](c)
		if err != nil {
			h.logger.Error("Ошибка при работе с JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Запрос на регистрацию: %v", r)

		userID, err := h.s.RegisterUser(c.Request.Context(), r.Email, r.Password, r.Role, r.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("ID созданого юзера: %v", userID)

		token, err := auth.GenerateToken(userID, r.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}

func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		r, err := utils.BindJSON[model.LoginReq](c)
		if err != nil {
			h.logger.Error("Ошибка при работе с JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Залогиновшийся юзер: %v", r)

		user, err := h.s.LoginUser(c.Request.Context(), r.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
			return
		}
		h.logger.Info("Залогиновшийся юзер: %v", user)

		if err := bcrypt.CompareHashAndPassword(
			[]byte(user.PasswordHash),
			[]byte(r.Password),
		); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
