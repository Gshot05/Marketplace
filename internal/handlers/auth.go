package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"marketplace/internal/auth"
	"marketplace/internal/model"
	"marketplace/internal/utils"

	sq "github.com/Masterminds/squirrel"
)

func register(pool *pgxpool.Pool) gin.HandlerFunc {
	type req struct {
		Email,
		Password,
		Role,
		Name string
	}
	return func(c *gin.Context) {
		var r req
		if err := c.BindJSON(&r); err != nil {
			c.JSON(400, gin.H{"error": "bad json"})
			return
		}

		if err := utils.ValidateEmail(r.Email); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := utils.ValidateName(r.Name); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if r.Role != "customer" && r.Role != "performer" {
			c.JSON(400, gin.H{"error": "bad role"})
			return
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
		ctx := c.Request.Context()

		queryBuilder := sq.Insert("users").
			Columns("email", "password_hash", "role", "name").
			Values(r.Email, string(hash), r.Role, r.Name).
			Suffix("RETURNING id").
			PlaceholderFormat(sq.Dollar)

		sql, args, err := queryBuilder.ToSql()
		if err != nil {
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}

		var newUserID int64
		err = pool.QueryRow(ctx, sql, args...).Scan(&newUserID)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		token, err := auth.GenerateToken(uint(newUserID), r.Role)
		if err != nil {
			c.JSON(500, gin.H{"error": "token generation failed"})
			return
		}

		c.JSON(200, gin.H{"token": token})
	}
}

func login(pool *pgxpool.Pool) gin.HandlerFunc {
	type req struct {
		Email,
		Password string
	}

	return func(c *gin.Context) {
		r, ok := utils.BindJSON[req](c)
		if !ok {
			return
		}

		ctx := c.Request.Context()

		queryBuilder := sq.Select("id", "email", "password_hash", "role").
			From("users").
			Where(sq.Eq{"email": r.Email}).
			PlaceholderFormat(sq.Dollar)

		sql, args, err := queryBuilder.ToSql()
		if err != nil {
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}

		var u model.User
		err = pool.QueryRow(ctx, sql, args...).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role)
		if err != nil {
			c.JSON(401, gin.H{"error": "Неверный логин или пароль"})
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(r.Password)) != nil {
			c.JSON(401, gin.H{"error": "Неверный логин или пароль"})
			return
		}

		token, err := auth.GenerateToken(u.ID, u.Role)
		if err != nil {
			c.JSON(500, gin.H{"error": "token generation failed"})
			return
		}

		c.JSON(200, gin.H{"token": token})
	}
}
