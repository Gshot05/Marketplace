package repository

import (
	"context"
	errors "marketplace/internal/error"
	"marketplace/internal/model"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FavoriteRepository struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
}

func NewFavoriteRepository(pool *pgxpool.Pool) *FavoriteRepository {
	return &FavoriteRepository{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *FavoriteRepository) Add(ctx context.Context, customerID, serviceID uint) (*model.FavoriteReq, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM services WHERE id = $1)",
		serviceID,
	).Scan(&exists)

	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, errors.ErrNotFindService
	}

	query := r.sb.Insert("favorites").
		Columns("customer_id", "service_id").
		Values(customerID, serviceID).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var fav model.FavoriteReq
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&fav.ID)
	if err != nil {
		return nil, err
	}

	fav.CustomerID = customerID
	fav.ServiceID = serviceID

	return &fav, nil
}

func (r *FavoriteRepository) Delete(ctx context.Context, customerID, serviceID uint) (bool, error) {
	query := r.sb.Delete("favorites").
		Where(sq.Eq{
			"customer_id": customerID,
			"service_id":  serviceID,
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return false, err
	}

	result, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return false, err
	}

	return result.RowsAffected() > 0, nil
}

func (r *FavoriteRepository) List(ctx context.Context, customerID uint) ([]model.FavoriteInfoReq, error) {
	query := r.sb.Select(
		"f.id",
		"u.name AS customer_name",
		"s.title AS service_title",
		"s.id AS serviceID",
	).
		From("favorites f").
		Join("users u ON f.customer_id = u.id").
		Join("services s ON f.service_id = s.id").
		Where(sq.Eq{"f.customer_id": customerID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []model.FavoriteInfoReq
	for rows.Next() {
		var f model.FavoriteInfoReq
		if err := rows.Scan(
			&f.ID,
			&f.CustomerName,
			&f.ServiceTitle,
			&f.ServiceID,
		); err != nil {
			return nil, err
		}
		favorites = append(favorites, f)
	}

	return favorites, nil
}
