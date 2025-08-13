package handlers

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

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

func (h *OfferHandler) CreateOffer() gin.HandlerFunc {
	type req struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	return func(c *gin.Context) {
		uid, ok := utils.CheckCustomerRole(c)
		if !ok {
			return
		}

		r, ok := utils.BindJSON[req](c)
		if !ok {
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
	type req struct {
		OfferID     uint    `json:"offerID"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	return func(c *gin.Context) {
		customerID, ok := utils.CheckCustomerRole(c)
		if !ok {
			return
		}

		r, ok := utils.BindJSON[req](c)
		if !ok {
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

		r, ok := utils.BindJSON[req](c)
		if !ok {
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
	type req struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	return func(c *gin.Context) {
		uid, ok := utils.CheckPerformerRole(c)
		if !ok {
			return
		}

		r, ok := utils.BindJSON[req](c)
		if !ok {
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
	type req struct {
		ServiceID   uint    `json:"serviceID"`
		Title       string  `json:"title,omitempty"`
		Description string  `json:"description,omitempty"`
		Price       float64 `json:"price,omitempty"`
	}

	return func(c *gin.Context) {
		performerID, ok := utils.CheckPerformerRole(c)
		if !ok {
			return
		}

		r, ok := utils.BindJSON[req](c)
		if !ok {
			return
		}

		service, err := h.repo.Update(c.Request.Context(), r.ServiceID, performerID, r.Title, r.Description, r.Price)
		if err != nil {
			if err.Error() == "service not found or access denied" {
				c.JSON(404, gin.H{"error": err.Error()})
				return
			}
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, service)
	}
}

func (h *ServiceHandler) DeleteService() gin.HandlerFunc {
	type req struct {
		ServiceID uint `json:"serviceID"`
	}

	return func(c *gin.Context) {
		performerID, ok := utils.CheckPerformerRole(c)
		if !ok {
			return
		}

		r, ok := utils.BindJSON[req](c)
		if !ok {
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

func addFavorite(pool *pgxpool.Pool) gin.HandlerFunc {
	type req struct{ ServiceID uint }

	return func(c *gin.Context) {
		uid, ok := utils.CheckCustomerRole(c)
		if !ok {
			return
		}

		r, ok := utils.BindJSON[req](c)
		if !ok {
			return
		}

		ctx := c.Request.Context()

		var exists bool
		err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM services WHERE id = $1)", r.ServiceID).Scan(&exists)
		if err != nil {
			c.JSON(500, gin.H{"error": "database error"})
			return
		}
		if !exists {
			c.JSON(404, gin.H{"error": "Услуга не найдена"})
			return
		}

		queryBuilder := sq.Insert("favorites").
			Columns("customer_id", "service_id").
			Values(uid, r.ServiceID).
			Suffix("RETURNING id").
			PlaceholderFormat(sq.Dollar)

		sql, args, err := queryBuilder.ToSql()
		if err != nil {
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}

		var newID uint
		err = pool.QueryRow(ctx, sql, args...).Scan(&newID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		fav := model.Favorite{
			ID:         newID,
			CustomerID: uid,
			ServiceID:  r.ServiceID,
		}

		c.JSON(200, fav)
	}
}

func listFavorites(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid := c.GetUint("user_id")
		ctx := c.Request.Context()

		queryBuilder := sq.Select(
			"f.id",
			"u.name AS customer_name",
			"s.title AS service_title",
			"s.description AS service_description",
		).
			From("favorites f").
			Join("users u ON f.customer_id = u.id").
			Join("services s ON f.service_id = s.id").
			Where(sq.Eq{"f.customer_id": uid}).
			PlaceholderFormat(sq.Dollar)

		sql, args, err := queryBuilder.ToSql()
		if err != nil {
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}

		rows, err := pool.Query(ctx, sql, args...)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var favs []model.FavoriteInfo
		for rows.Next() {
			var f model.FavoriteInfo
			err := rows.Scan(&f.ID, &f.CustomerName, &f.ServiceTitle, &f.ServiceDescription)
			if err != nil {
				c.JSON(500, gin.H{"error": "error scanning row"})
				return
			}
			favs = append(favs, f)
		}

		c.JSON(200, favs)
	}
}

func deleteFavorite(pool *pgxpool.Pool) gin.HandlerFunc {
	type req struct {
		ServiceID uint `json:"serviceID"`
	}

	return func(c *gin.Context) {
		customerID, ok := utils.CheckCustomerRole(c)
		if !ok {
			return
		}

		req, ok := utils.BindJSON[req](c)
		if !ok {
			return
		}

		ctx := c.Request.Context()

		queryBuilder := sq.Delete("favorites").
			Where(sq.Eq{
				"customer_id": customerID,
				"service_id":  req.ServiceID,
			}).
			PlaceholderFormat(sq.Dollar)

		sql, args, err := queryBuilder.ToSql()
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to build SQL query: " + err.Error()})
			return
		}

		result, err := pool.Exec(ctx, sql, args...)
		if err != nil {
			c.JSON(500, gin.H{"error": "database error: " + err.Error()})
			return
		}

		if result.RowsAffected() == 0 {
			c.JSON(404, gin.H{"error": "Избранное не найдено"})
			return
		}

		c.String(200, "Успешно!")
	}
}
