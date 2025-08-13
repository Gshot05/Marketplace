package handlers

import (
	"github.com/gin-gonic/gin"

	"marketplace/internal/model"
	repository "marketplace/internal/repo"
	"marketplace/internal/utils"
)

type OfferHandler struct {
	repo *repository.OfferRepository
}

func NewOfferHandler(repo *repository.OfferRepository) *OfferHandler {
	return &OfferHandler{repo: repo}
}

type ServiceHandler struct {
	repo *repository.ServiceRepository
}

func NewServiceHandler(repo *repository.ServiceRepository) *ServiceHandler {
	return &ServiceHandler{repo: repo}
}

type FavoriteHandler struct {
	repo *repository.FavoriteRepository
}

func NewFavoriteHandler(repo *repository.FavoriteRepository) *FavoriteHandler {
	return &FavoriteHandler{repo: repo}
}

func (h *OfferHandler) CreateOffer() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := utils.CheckCustomerRole(c)
		if !ok {
			return
		}

		r, err := utils.BindJSON[model.OfferCreateReq](c)
		if err != nil {
			c.JSON(400, gin.H{"error": "Неверный формат запроса: " + err.Error()})
			return
		}

		offer, err := h.repo.Create(c.Request.Context(), uid, r.Title, r.Description, r.Price)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, offer)
	}
}

func (h *OfferHandler) UpdateOffer() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID, ok := utils.CheckCustomerRole(c)
		if !ok {
			return
		}

		r, err := utils.BindJSON[model.OfferUpdateReq](c)
		if err != nil {
			c.JSON(400, gin.H{"error": "Неверный формат запроса: " + err.Error()})
			return
		}

		offer, err := h.repo.Update(c.Request.Context(), r.OfferID, customerID, r.Title, r.Description, r.Price)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, offer)
	}
}

func (h *OfferHandler) DeleteOffer() gin.HandlerFunc {
	type req struct {
		OfferID uint `json:"offerID"`
	}

	return func(c *gin.Context) {
		customerID, ok := utils.CheckCustomerRole(c)
		if !ok {
			return
		}

		r, err := utils.BindJSON[req](c)
		if err != nil {
			c.JSON(400, gin.H{"error": "Неверный формат запроса: " + err.Error()})
			return
		}

		deleted, err := h.repo.Delete(c.Request.Context(), r.OfferID, customerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if !deleted {
			c.String(404, "Оффер не найден или вам не принадлежит")
			return
		}

		c.String(200, "Успешно!")
	}
}

func (h *OfferHandler) ListOffers() gin.HandlerFunc {
	return func(c *gin.Context) {
		offers, err := h.repo.List(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if len(offers) == 0 {
			c.String(404, "Офферов пока нет:(")
			return
		}

		c.JSON(200, offers)
	}
}

func (h *ServiceHandler) CreateService() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := utils.CheckPerformerRole(c)
		if !ok {
			return
		}

		r, err := utils.BindJSON[model.ServiceCreateReq](c)
		if err != nil {
			c.JSON(400, gin.H{"error": "Неверный формат запроса: " + err.Error()})
			return
		}

		service, err := h.repo.Create(c.Request.Context(), uid, r.Title, r.Description, r.Price)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, service)
	}
}

func (h *ServiceHandler) UpdateService() gin.HandlerFunc {
	return func(c *gin.Context) {
		performerID, ok := utils.CheckPerformerRole(c)
		if !ok {
			return
		}

		r, err := utils.BindJSON[model.ServiceUpdateReq](c)
		if err != nil {
			c.JSON(400, gin.H{"error": "Неверный формат запроса: " + err.Error()})
			return
		}

		service, err := h.repo.Update(c.Request.Context(), r.ServiceID, performerID, r.Title, r.Description, r.Price)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, service)
	}
}

func (h *ServiceHandler) DeleteService() gin.HandlerFunc {
	return func(c *gin.Context) {
		performerID, ok := utils.CheckPerformerRole(c)
		if !ok {
			return
		}

		r, err := utils.BindJSON[model.ServiceDeleteReq](c)
		if err != nil {
			c.JSON(400, gin.H{"error": "Неверный формат запроса: " + err.Error()})
			return
		}

		deleted, err := h.repo.Delete(c.Request.Context(), r.ServiceID, performerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if !deleted {
			c.String(404, "Услуга не найдена или вам не принадлежит")
			return
		}

		c.String(200, "Успешно!")
	}
}

func (h *ServiceHandler) ListServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		services, err := h.repo.List(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if len(services) == 0 {
			c.String(404, "Услуг пока нет:(")
			return
		}

		c.JSON(200, services)
	}
}

func (h *FavoriteHandler) AddFavorite() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID, ok := utils.CheckCustomerRole(c)
		if !ok {
			return
		}

		r, err := utils.BindJSON[model.FavoriteAddReq](c)
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

		r, err := utils.BindJSON[model.FavoriteDeleteReq](c)
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
			c.String(404, "Избранное не найдено")
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
			c.String(404, "Избранное пока пусто")
			return
		}

		c.JSON(200, favorites)
	}
}
