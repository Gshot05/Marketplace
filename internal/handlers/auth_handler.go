package handlers

import (
	"marketplace/internal/logger"
	"marketplace/internal/model"
	"marketplace/internal/service"
	"marketplace/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
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
		logger: logger,
	}
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

		err = h.s.RegisterUser(c.Request.Context(), r.Email, r.Password, r.Role, r.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Подтвердите почту"})
	}
}

func (h *AuthHandler) ConfirmEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		r, err := utils.BindJSON[model.VerifyUser](c)

		if err != nil {
			h.logger.Error("Ошибка при работе с JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := h.s.ConfirmEmail(c.Request.Context(), r.Email, r.VerifyCode)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

		token, err := h.s.LoginUser(c.Request.Context(), r.Email, r.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
