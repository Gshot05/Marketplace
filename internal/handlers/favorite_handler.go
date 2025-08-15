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
		customerID, err := utils.CheckCustomerRole(c)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		r, err := utils.BindJSON[model.GeneralServiceIdReq](c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
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
		customerID, err := utils.CheckCustomerRole(c)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}

		r, err := utils.BindJSON[model.GeneralServiceIdReq](c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		deleted, err := h.repo.Delete(c.Request.Context(), customerID, r.ServiceID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if !deleted {
			c.JSON(404, gin.H{"error": "Избранное не найдено или вам не принадлежит!"})
			return
		}

		c.JSON(200, gin.H{"success": "Успешно!"})
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
			c.JSON(404, gin.H{"error": "Избранное пусто:("})
			return
		}

		c.JSON(200, favorites)
	}
}
