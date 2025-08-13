package repository

import (
	"context"
	"errors"
	"marketplace/internal/model"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
}

func NewServiceRepository(pool *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *ServiceRepository) Create(ctx context.Context, performerID uint, title, description string, price float64) (*model.Service, error) {
	query := r.sb.Insert("services").
		Columns("performer_id", "title", "description", "price").
		Values(performerID, title, description, price).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var id uint
	err = r.pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &model.Service{
		ID:          id,
		PerformerID: performerID,
		Title:       title,
		Description: description,
		Price:       price,
	}, nil
}

func (r *ServiceRepository) Update(ctx context.Context, serviceID, performerID uint, title, description string, price float64) (*model.Service, error) {
	query := r.sb.Update("services").
		SetMap(map[string]interface{}{
			"title":       title,
			"description": description,
			"price":       price,
		}).
		Where(sq.Eq{
			"id":           serviceID,
			"performer_id": performerID,
		}).
		Suffix("RETURNING id, performer_id, title, description, price")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var service model.Service
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&service.ID, &service.PerformerID, &service.Title, &service.Description, &service.Price,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("service not found or access denied")
		}
		return nil, err
	}

	return &service, nil
}

func (r *ServiceRepository) Delete(ctx context.Context, serviceID, performerID uint) (bool, error) {
	query := r.sb.Delete("services").
		Where(sq.Eq{
			"id":           serviceID,
			"performer_id": performerID,
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

func (r *ServiceRepository) List(ctx context.Context) ([]model.Service, error) {
	query := r.sb.Select("id", "performer_id", "title", "description", "price").
		From("services")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []model.Service
	for rows.Next() {
		var s model.Service
		if err := rows.Scan(&s.ID, &s.PerformerID, &s.Title, &s.Description, &s.Price); err != nil {
			return nil, err
		}
		services = append(services, s)
	}

	return services, nil
}
