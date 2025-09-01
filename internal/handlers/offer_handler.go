package handlers

import (
	"marketplace/internal/logger"
	"marketplace/internal/model"
	"marketplace/internal/service"
	"marketplace/internal/utils"

	"github.com/gin-gonic/gin"
)

type OfferHandler struct {
	s      service.IOfferService
	logger *logger.Logger
}

func NewOfferHandler(
	s service.IOfferService,
	logger *logger.Logger) *OfferHandler {
	return &OfferHandler{
		s:      s,
		logger: logger,
	}
}

func (h *OfferHandler) CreateOffer() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID, err := utils.CheckCustomerRole(c)
		if err != nil {
			h.logger.Error("Ошибка проверки роли: %v", err)
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Запрос на создание оффера человека с ID: %v", customerID)

		r, err := utils.BindJSON[model.OfferCreateReq](c)
		if err != nil {
			h.logger.Error("Ошибка при работе с JSON: %v", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Оффер который пришёл: %v", r)

		offer, err := h.s.CreateOffer(c.Request.Context(), customerID, r.Title, r.Description, r.Price)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Созданный оффер %v", offer)

		c.JSON(200, offer)
	}
}

func (h *OfferHandler) UpdateOffer() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID, err := utils.CheckCustomerRole(c)
		if err != nil {
			h.logger.Error("Ошибка проверки роли: %v", err)
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Запрос на обновление оффера человека с ID: %v", customerID)

		r, err := utils.BindJSON[model.OfferUpdateReq](c)
		if err != nil {
			h.logger.Error("Ошибка при работе с JSON: %v", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Оффер который пришёл: %v", r)

		offer, err := h.s.UpdateOffer(c.Request.Context(), r.OfferID, customerID, r.Title, r.Description, r.Price)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Обновлённый оффер %v", offer)

		c.JSON(200, offer)
	}
}

func (h *OfferHandler) DeleteOffer() gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID, err := utils.CheckCustomerRole(c)
		if err != nil {
			h.logger.Error("Ошибка проверки роли: %v", err)
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Запрос на удаление оффера человека с ID: %v", customerID)

		r, err := utils.BindJSON[model.OfferDeleteReq](c)
		if err != nil {
			h.logger.Error("Ошибка при работе с JSON: %v", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		deleted, err := h.s.DeleteOffer(c.Request.Context(), r.OfferID, customerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if !deleted {
			c.JSON(404, gin.H{"error": "Оффер не найден или вам не принадлежит!"})
			return
		}

		c.JSON(200, gin.H{"success": "Успешно!"})
	}
}

func (h *OfferHandler) ListOffers() gin.HandlerFunc {
	return func(c *gin.Context) {
		offers, err := h.s.ListOffers(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if len(offers) == 0 {
			c.JSON(404, gin.H{"error": "Офферов пока нет:("})
			return
		}

		c.JSON(200, offers)
	}
}
