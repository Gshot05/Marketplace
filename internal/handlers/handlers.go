package handlers

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace/internal/model"
	"marketplace/internal/utils"
)

func createOffer(pool *pgxpool.Pool) gin.HandlerFunc {
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

		ctx := c.Request.Context()

		queryBuilder := sq.Insert("offers").
			Columns("customer_id", "title", "description", "price").
			Values(uid, r.Title, r.Description, r.Price).
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

		o := model.Offer{
			ID:          newID,
			CustomerID:  uid,
			Title:       r.Title,
			Description: r.Description,
			Price:       r.Price,
		}

		c.JSON(200, o)
	}
}

func updateOffer(pool *pgxpool.Pool) gin.HandlerFunc {
	type request struct {
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

		req, ok := utils.BindJSON[request](c)
		if !ok {
			return
		}

		ctx := c.Request.Context()

		queryBuilder := sq.Update("offers").
			SetMap(map[string]interface{}{
				"title":       req.Title,
				"description": req.Description,
				"price":       req.Price,
			}).
			Where(sq.Eq{
				"id":          req.OfferID,
				"customer_id": customerID,
			}).
			Suffix("RETURNING id, customer_id, title, description, price").
			PlaceholderFormat(sq.Dollar)

		sql, args, err := queryBuilder.ToSql()
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to build query"})
			return
		}

		var updated model.Offer
		if err := pool.QueryRow(ctx, sql, args...).Scan(
			&updated.ID, &updated.CustomerID, &updated.Title,
			&updated.Description, &updated.Price,
		); err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(404, gin.H{"error": "Оффер не найден или принадлежит не вам!"})
				return
			}
			c.JSON(500, gin.H{"error": "database error"})
			return
		}

		c.JSON(200, updated)
	}
}

func deleteOffer(pool *pgxpool.Pool) gin.HandlerFunc {
	type request struct {
		OfferID uint `json:"offerID"`
	}

	return func(c *gin.Context) {
		userID, ok := utils.CheckCustomerRole(c)
		if !ok {
			return
		}

		req, ok := utils.BindJSON[request](c)
		if !ok {
			return
		}

		ctx := c.Request.Context()

		queryBuilder := sq.Delete("offers").
			Where(sq.Eq{
				"id":          req.OfferID,
				"customer_id": userID,
			}).
			PlaceholderFormat(sq.Dollar)

		sql, args, err := queryBuilder.ToSql()
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to build SQL query"})
			return
		}

		result, err := pool.Exec(ctx, sql, args...)
		if err != nil {
			c.JSON(500, gin.H{"error": "database error"})
			return
		}

		if result.RowsAffected() == 0 {
			c.String(404, "Оффер не найден или вам не принадлежит")
			return
		}

		c.String(200, "Успешно!")
	}
}

func listOffers(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		queryBuilder := sq.Select("id", "customer_id", "title", "description", "price").
			From("offers").
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

		var offers []model.Offer
		for rows.Next() {
			var o model.Offer
			err := rows.Scan(&o.ID, &o.CustomerID, &o.Title, &o.Description, &o.Price)
			if err != nil {
				c.JSON(500, gin.H{"error": "error scanning row"})
				return
			}
			offers = append(offers, o)
		}

		if offers == nil {
			c.String(404, "Офферов пока нет:(")
			return
		}

		c.JSON(200, offers)
	}
}

func createService(pool *pgxpool.Pool) gin.HandlerFunc {
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

		ctx := c.Request.Context()

		queryBuilder := sq.Insert("services").
			Columns("performer_id", "title", "description", "price").
			Values(uid, r.Title, r.Description, r.Price).
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

		s := model.Service{
			ID:          newID,
			PerformerID: uid,
			Title:       r.Title,
			Description: r.Description,
			Price:       r.Price,
		}

		c.JSON(200, s)
	}
}

func updateService(pool *pgxpool.Pool) gin.HandlerFunc {
	type request struct {
		ServiceID   uint    `json:"serviceID" binding:"required"`
		Title       string  `json:"title,omitempty"`
		Description string  `json:"description,omitempty"`
		Price       float64 `json:"price,omitempty"`
	}

	return func(c *gin.Context) {
		performerID, ok := utils.CheckPerformerRole(c)
		if !ok {
			return
		}

		req, ok := utils.BindJSON[request](c)
		if !ok {
			return
		}

		ctx := c.Request.Context()

		queryBuilder := sq.Update("services").
			SetMap(map[string]interface{}{
				"title":       req.Title,
				"description": req.Description,
				"price":       req.Price,
			}).
			Where(sq.Eq{
				"id":           req.ServiceID,
				"performer_id": performerID,
			}).
			Suffix("RETURNING id, performer_id, title, description, price").
			PlaceholderFormat(sq.Dollar)

		sql, args, err := queryBuilder.ToSql()
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to build query"})
			return
		}

		var updated model.Service
		if err := pool.QueryRow(ctx, sql, args...).Scan(
			&updated.ID,
			&updated.PerformerID,
			&updated.Title,
			&updated.Description,
			&updated.Price,
		); err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(404, gin.H{"error": "service not found or access denied"})
				return
			}
			c.JSON(500, gin.H{"error": "database error"})
			return
		}

		c.JSON(200, updated)
	}
}

func deleteService(pool *pgxpool.Pool) gin.HandlerFunc {
	type request struct {
		ServiceID uint `json:"ServiceID"`
	}

	return func(c *gin.Context) {
		userID, ok := utils.CheckPerformerRole(c)
		if !ok {
			return
		}

		req, ok := utils.BindJSON[request](c)
		if !ok {
			return
		}

		ctx := c.Request.Context()

		queryBuilder := sq.Delete("services").
			Where(sq.Eq{
				"id":           req.ServiceID,
				"performer_id": userID,
			}).
			PlaceholderFormat(sq.Dollar)

		sql, args, err := queryBuilder.ToSql()
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to build SQL query"})
			return
		}

		result, err := pool.Exec(ctx, sql, args...)
		if err != nil {
			c.JSON(500, gin.H{"error": "database error"})
			return
		}

		if result.RowsAffected() == 0 {
			c.String(404, "Услуга не найдена или вам не принадлежит")
			return
		}

		c.String(200, "Успешно!")
	}
}

func listServices(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		queryBuilder := sq.Select("id", "performer_id", "title", "description", "price").
			From("services").
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

		var services []model.Service
		for rows.Next() {
			var s model.Service
			if err := rows.Scan(&s.ID, &s.PerformerID, &s.Title, &s.Description, &s.Price); err != nil {
				c.JSON(500, gin.H{"error": "error scanning row"})
				return
			}
			services = append(services, s)
		}
		if services == nil {
			c.String(400, "Услуг пока нет:(")
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
	type request struct {
		ServiceID uint `json:"serviceID"`
	}

	return func(c *gin.Context) {
		customerID, ok := utils.CheckCustomerRole(c)
		if !ok {
			return
		}

		req, ok := utils.BindJSON[request](c)
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
