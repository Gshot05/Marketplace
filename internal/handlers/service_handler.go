package handlers

import (
	"marketplace/internal/logger"
	"marketplace/internal/model"
	repository "marketplace/internal/repo"
	"marketplace/internal/utils"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	repo   *repository.ServiceRepository
	logger *logger.Logger
}

func NewServiceHandler(
	repo *repository.ServiceRepository,
	logger *logger.Logger,
) *ServiceHandler {
	return &ServiceHandler{
		repo:   repo,
		logger: logger}
}

func (h *ServiceHandler) CreateService() gin.HandlerFunc {
	return func(c *gin.Context) {
		performerID, err := utils.CheckPerformerRole(c)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Запрос на создание сервиса человека с ID: %v", performerID)

		r, err := utils.BindJSON[model.ServiceCreateReq](c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Сервис который пришёл: %v", r)

		service, err := h.repo.Create(c.Request.Context(), performerID, r.Title, r.Description, r.Price)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Созданный сервис %v", service)

		c.JSON(200, service)
	}
}

func (h *ServiceHandler) UpdateService() gin.HandlerFunc {
	return func(c *gin.Context) {
		performerID, err := utils.CheckPerformerRole(c)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Запрос на обновление оффера человека с ID: %v", performerID)

		r, err := utils.BindJSON[model.ServiceUpdateReq](c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Сервис который пришёл: %v", r)

		service, err := h.repo.Update(c.Request.Context(), r.ServiceID, performerID, r.Title, r.Description, r.Price)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Обновлённый оффер %v", service)

		c.JSON(200, service)
	}
}

func (h *ServiceHandler) DeleteService() gin.HandlerFunc {
	return func(c *gin.Context) {
		performerID, err := utils.CheckPerformerRole(c)
		if err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Запрос на удаление оффера человека с ID: %v", performerID)

		r, err := utils.BindJSON[model.GeneralServiceIdReq](c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		deleted, err := h.repo.Delete(c.Request.Context(), r.ServiceID, performerID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if !deleted {
			c.JSON(404, gin.H{"error": "Услуга не найдена или вам не принадлежит!"})
			return
		}

		c.JSON(200, gin.H{"success": "Успешно!"})
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
			c.JSON(404, gin.H{"error": "Услуг пока нет:("})
			return
		}

		c.JSON(200, services)
	}
}
