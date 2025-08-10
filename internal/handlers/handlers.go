package handlers

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"marketplace/internal/model"
)

func createOffer(pool *pgxpool.Pool) gin.HandlerFunc {
	type req struct {
		Title, Description string
		Price              float64
	}
	return func(c *gin.Context) {
		uid := c.GetUint("user_id")
		role := c.GetString("role")
		if role != "customer" {
			c.JSON(403, gin.H{"error": "only customers can create offers"})
			return
		}

		var r req
		if err := c.BindJSON(&r); err != nil {
			c.JSON(400, gin.H{"error": "bad json"})
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

		c.JSON(200, offers)
	}
}

func createService(pool *pgxpool.Pool) gin.HandlerFunc {
	type req struct {
		Title, Description string
		Price              float64
	}
	return func(c *gin.Context) {
		uid := c.GetUint("user_id")
		role := c.GetString("role")
		if role != "performer" {
			c.JSON(403, gin.H{"error": "only performers can create services"})
			return
		}

		var r req
		if err := c.BindJSON(&r); err != nil {
			c.JSON(400, gin.H{"error": "bad json"})
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

		c.JSON(200, services)
	}
}

func addFavorite(pool *pgxpool.Pool) gin.HandlerFunc {
	type req struct{ ServiceID uint }

	return func(c *gin.Context) {
		uid := c.GetUint("user_id")
		role := c.GetString("role")
		if role != "customer" {
			c.JSON(403, gin.H{"error": "only customers can favorite"})
			return
		}

		var r req
		if err := c.BindJSON(&r); err != nil {
			c.JSON(400, gin.H{"error": "bad json"})
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
			c.JSON(404, gin.H{"error": "service not found"})
			return
		}

		// Вставляем новый favorite
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

		queryBuilder := sq.Select("id", "customer_id", "service_id").
			From("favorites").
			Where(sq.Eq{"customer_id": uid}).
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

		var favs []model.Favorite
		for rows.Next() {
			var f model.Favorite
			if err := rows.Scan(&f.ID,
				&f.CustomerID,
				&f.ServiceID); err != nil {
				c.JSON(500, gin.H{"error": "error scanning row"})
				return
			}
			favs = append(favs, f)
		}

		c.JSON(200, favs)
	}
}
