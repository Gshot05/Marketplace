package handlers

import (
	"marketplace/internal/model"
	repository "marketplace/internal/repo"
	"marketplace/internal/utils"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	repo *repository.FavoriteRepository
}

func NewFavoriteHandler(repo *repository.FavoriteRepository) *FavoriteHandler {
	return &FavoriteHandler{repo: repo}
}

func (h *FavoriteHandler) AddFavorite() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID, ok := utils.CheckCustomerRole(c)
		if !ok {
			return
		}

		r, err := utils.BindJSON[model.GeneralServiceIdReq](c)
		if err != nil {
			c.JSON(400, gin.H{"error": "Неверный формат запроса: " + err.Error()})
			return
		}

		fav, err := h.repo.Add(c.Request.Context(), customerID, r.ServiceID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, fav)
	}
}

func (h *FavoriteHandler) DeleteFavorite() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID, ok := utils.CheckCustomerRole(c)
		if !ok {
			return
		}

		r, err := utils.BindJSON[model.GeneralServiceIdReq](c)
		if err != nil {
			c.JSON(400, gin.H{"error": "Неверный формат запроса: " + err.Error()})
			return
		}

		deleted, err := h.repo.Delete(c.Request.Context(), customerID, r.ServiceID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if !deleted {
			c.String(404, "Избранное не найдено:(")
			return
		}

		c.String(200, "Успешно!")
	}
}

func (h *FavoriteHandler) ListFavorites() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID := c.GetUint("user_id")

		favorites, err := h.repo.List(c.Request.Context(), customerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if len(favorites) == 0 {
			c.String(404, "Избранное пока пусто :(")
			return
		}

		c.JSON(200, favorites)
	}
}
