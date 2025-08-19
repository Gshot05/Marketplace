package repository

import (
	"context"
	errors2 "marketplace/internal/error"
	"marketplace/internal/model"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OfferRepository struct {
	pool *pgxpool.Pool
	sb   squirrel.StatementBuilderType
}

func NewOfferRepository(pool *pgxpool.Pool) *OfferRepository {
	return &OfferRepository{
		pool: pool,
		sb:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *OfferRepository) Create(ctx context.Context, customerID uint, title, description string, price float64) (*model.Offer, error) {
	query := r.sb.Insert("offers").
		Columns("customer_id", "title", "description", "price").
		Values(customerID, title, description, price).
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

	return &model.Offer{
		ID:          id,
		CustomerID:  customerID,
		Title:       title,
		Description: description,
		Price:       price,
	}, nil
}

func (r *OfferRepository) Update(ctx context.Context, offerID, customerID uint, title, description string, price float64) (*model.Offer, error) {
	query := r.sb.Update("offers").
		SetMap(map[string]interface{}{
			"title":       title,
			"description": description,
			"price":       price,
		}).
		Where(sq.Eq{
			"id":          offerID,
			"customer_id": customerID,
		}).
		Suffix("RETURNING id, customer_id, title, description, price")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var offer model.Offer
	err = r.pool.QueryRow(ctx, sql, args...).Scan(
		&offer.ID, &offer.CustomerID, &offer.Title, &offer.Description, &offer.Price,
	)
	if err != nil {
		return nil, errors2.ErrWrongUpdateOffer
	}

	return &offer, nil
}

func (r *OfferRepository) Delete(ctx context.Context, offerID, customerID uint) (bool, error) {
	query := r.sb.Delete("offers").
		Where(sq.Eq{
			"id":          offerID,
			"customer_id": customerID,
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

func (r *OfferRepository) List(ctx context.Context) ([]model.Offer, error) {
	query := r.sb.Select("id", "customer_id", "title", "description", "price").
		From("offers")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offers []model.Offer
	for rows.Next() {
		var o model.Offer
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.Title, &o.Description, &o.Price); err != nil {
			return nil, err
		}
		offers = append(offers, o)
	}

	return offers, nil
}
