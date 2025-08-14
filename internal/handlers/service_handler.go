package handlers

import (
	"marketplace/internal/model"
	repository "marketplace/internal/repo"
	"marketplace/internal/utils"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	repo *repository.ServiceRepository
}

func NewServiceHandler(repo *repository.ServiceRepository) *ServiceHandler {
	return &ServiceHandler{repo: repo}
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

		r, err := utils.BindJSON[model.GeneralServiceIdReq](c)
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
