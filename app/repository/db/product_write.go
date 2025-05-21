package db

import (
	"context"
	"database/sql"
	"log/slog"
	"product-service/app/domain"
)

type productWriteRepository struct {
	conn *sql.DB
}

func NewProductWriteRepository(db *sql.DB) domain.ProductWriteRepository {
	return &productWriteRepository{db}
}

func (r *productWriteRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `INSERT INTO products (name, description, price, category, image_url, shop_id, active) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id , created_at, updated_at`

	err := r.conn.QueryRowContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Category,
		product.ImageURL,
		product.ShopID,
		product.Active).
		Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "[productWriteRepository] Create", "scan", err)
		return domain.ErrInternal
	}

	return nil
}

func (r *productWriteRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `UPDATE products SET name = $1, description = $2, price = $3, category = $4, image_url = $5, shop_id = $6, active = $7, updated_at = $8 WHERE id = $9`
	_, err := r.conn.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Category,
		product.ImageURL,
		product.ShopID,
		product.Active,
		product.UpdatedAt,
		product.ID)
	if err != nil {
		slog.ErrorContext(ctx, "[productWriteRepository] Update", "exec", err)
		return domain.ErrInternal
	}

	return nil
}

func (r *productWriteRepository) SetActiveStatus(ctx context.Context, id int64, active bool) error {
	query := `UPDATE products SET active = $1 WHERE id = $2`
	_, err := r.conn.ExecContext(ctx, query, active, id)
	if err != nil {
		slog.ErrorContext(ctx, "[productWriteRepository] SetActiveStatus", "exec", err)
		return domain.ErrInternal
	}

	return nil
}

func (r *productWriteRepository) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		slog.ErrorContext(ctx, "[productWriteRepository] BeginTransaction", "beginTx", err)
		return nil, err
	}
	return tx, nil
}

func (r *productWriteRepository) WithTransaction(ctx context.Context, tx *sql.Tx, fn func(context.Context, *sql.Tx) error) error {
	if err := fn(ctx, tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			slog.ErrorContext(ctx, "[productWriteRepository] WithTransaction", "rollback", rollbackErr)
			return rollbackErr
		}
		return err
	}
	return tx.Commit()
}
