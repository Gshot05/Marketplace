package repository

import (
	"context"
	"errors"
	errors2 "marketplace/internal/error"
	"marketplace/internal/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepo struct {
	pool *pgxpool.Pool
}

func NewAuthRepo(pool *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{pool: pool}
}

func (r *AuthRepo) RegisterUser(ctx context.Context, email, password, role, name string) (uint, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	query := sq.Insert("users").
		Columns("email", "password_hash", "role", "name").
		Values(email, string(hash), role, name).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}

	var id uint
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&id)
	return id, err
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := sq.Select("id", "email", "password_hash", "role").
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var user model.User
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors2.ErrWrongPassOrLog
		}
		return nil, err
	}

	return &user, nil
}
