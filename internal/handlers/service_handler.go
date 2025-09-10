package handlers

import (
	"marketplace/internal/logger"
	"marketplace/internal/model"
	"marketplace/internal/service"
	"marketplace/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	s      service.IServiceService
	logger *logger.Logger
}

func NewServiceHandler(
	s service.IServiceService,
	logger *logger.Logger,
) *ServiceHandler {
	return &ServiceHandler{
		s:      s,
		logger: logger,
	}
}

func (h *ServiceHandler) CreateService() gin.HandlerFunc {
	return func(c *gin.Context) {
		performerID, err := utils.CheckPerformerRole(c)
		if err != nil {
			h.logger.Error("Ошибка проверки роли: %v", err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Запрос на создание сервиса человека с ID: %v", performerID)

		r, err := utils.BindJSON[model.ServiceCreateReq](c)
		if err != nil {
			h.logger.Error("Ошибка при работе с JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Сервис который пришёл: %v", r)

		service, err := h.s.CreateService(c.Request.Context(), performerID, r.Title, r.Description, r.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Созданный сервис %v", service)

		c.JSON(http.StatusOK, service)
	}
}

func (h *ServiceHandler) UpdateService() gin.HandlerFunc {
	return func(c *gin.Context) {
		performerID, err := utils.CheckPerformerRole(c)
		if err != nil {
			h.logger.Error("Ошибка проверки роли: %v", err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Запрос на обновление оффера человека с ID: %v", performerID)

		r, err := utils.BindJSON[model.ServiceUpdateReq](c)
		if err != nil {
			h.logger.Error("Ошибка при работе с JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Сервис который пришёл: %v", r)

		service, err := h.s.UpdateService(c.Request.Context(), r.ServiceID, performerID, r.Title, r.Description, r.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Обновлённый оффер %v", service)

		c.JSON(http.StatusOK, service)
	}
}

func (h *ServiceHandler) DeleteService() gin.HandlerFunc {
	return func(c *gin.Context) {
		performerID, err := utils.CheckPerformerRole(c)
		if err != nil {
			h.logger.Error("Ошибка проверки роли: %v", err)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		h.logger.Info("Запрос на удаление оффера человека с ID: %v", performerID)

		r, err := utils.BindJSON[model.GeneralServiceIdReq](c)
		if err != nil {
			h.logger.Error("Ошибка при работе с JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		deleted, err := h.s.DeleteService(c.Request.Context(), r.ServiceID, performerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if !deleted {
			c.JSON(http.StatusNotFound, gin.H{"error": "Услуга не найдена или вам не принадлежит!"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": "Успешно!"})
	}
}

func (h *ServiceHandler) ListServices() gin.HandlerFunc {
	return func(c *gin.Context) {
		services, err := h.s.ListServices(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, services)
	}
}
