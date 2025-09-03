package handlers

import (
	"marketplace/internal/logger"
	"marketplace/internal/model"
	"marketplace/internal/service"
	"marketplace/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	s      service.IFavoriteService
	logger *logger.Logger
}

func NewFavoriteHandler(
	s service.IFavoriteService,
	logger *logger.Logger,
) *FavoriteHandler {
	return &FavoriteHandler{
		s:      s,
		logger: logger}
}

func (h *FavoriteHandler) AddFavorite() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID, err := utils.CheckCustomerRole(c)
		if err != nil {
			h.logger.Error("Ошибка проверки роли: %v", customerID)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Запрос на добавление в избранное от человека с ID: %v", customerID)

		r, err := utils.BindJSON[model.GeneralServiceIdReq](c)
		if err != nil {
			h.logger.Error("Ошибка при работе с JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Запрос на добавление: %v", r)

		fav, err := h.s.AddFavorite(c.Request.Context(), customerID, r.ServiceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Добавленное избранное: %v", fav)

		c.JSON(http.StatusOK, fav)
	}
}

func (h *FavoriteHandler) DeleteFavorite() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID, err := utils.CheckCustomerRole(c)
		if err != nil {
			h.logger.Error("Ошибка проверки роли: %v", err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Удаление из избранного от человека с ID: %v", customerID)

		r, err := utils.BindJSON[model.GeneralServiceIdReq](c)
		if err != nil {
			h.logger.Error("Ошибка при работе с JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		deleted, err := h.s.DeleteFavorite(c.Request.Context(), customerID, r.ServiceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !deleted {
			c.JSON(http.StatusNotFound, gin.H{"error": "Избранное не найдено или вам не принадлежит!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": "Успешно!"})
	}
}

func (h *FavoriteHandler) ListFavorites() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID := c.GetUint("user_id")

		favorites, err := h.s.ListFavorites(c.Request.Context(), customerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(favorites) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Избранное пусто:("})
			return
		}
		c.JSON(http.StatusOK, favorites)
	}
}
