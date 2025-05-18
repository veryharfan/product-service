package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"product-service/app/domain"
	"strings"
)

type productReadRepository struct {
	conn *sql.DB
}

func NewProductReadRepository(db *sql.DB) domain.ProductReadRepository {
	return &productReadRepository{db}
}

func (r *productReadRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `SELECT id, name, description, price, category, image_url, shop_id, active, created_at, updated_at FROM products WHERE id = $1 AND active = true`
	row := r.conn.QueryRowContext(ctx, query, id)

	var product domain.Product
	if err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.ImageURL, &product.ShopID, &product.Active, &product.CreatedAt, &product.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		slog.ErrorContext(ctx, "[productReadRepository] GetByID", "query", err)
		return nil, domain.ErrInternal
	}

	return &product, nil
}

func (r *productReadRepository) GetListByQuery(ctx context.Context, query domain.ProductQuery) ([]*domain.Product, error) {
	sqlQuery := `SELECT id, name, description, price, category, image_url, shop_id, active, created_at, updated_at FROM products WHERE active = true`
	args := []any{}

	placeholderIndex := 1 // Start placeholder index

	if query.ShopID > 0 {
		sqlQuery += fmt.Sprintf(" AND shop_id = $%d", placeholderIndex)
		args = append(args, query.ShopID)
		placeholderIndex++
	}
	if query.Category != "" {
		sqlQuery += fmt.Sprintf(" AND category = $%d", placeholderIndex)
		args = append(args, strings.ToLower(query.Category))
		placeholderIndex++
	}
	if query.MinPrice > 0 {
		sqlQuery += fmt.Sprintf(" AND price >= $%d", placeholderIndex)
		args = append(args, query.MinPrice)
		placeholderIndex++
	}
	if query.MaxPrice > 0 {
		sqlQuery += fmt.Sprintf(" AND price <= $%d", placeholderIndex)
		args = append(args, query.MaxPrice)
		placeholderIndex++
	}
	if query.Keyword != "" {
		sqlQuery += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", placeholderIndex, placeholderIndex)
		args = append(args, "%"+query.Keyword+"%")
		placeholderIndex++
	}

	if query.SortBy == "" {
		query.SortBy = "created_at"
	}
	if query.SortOrder == "" {
		query.SortOrder = "asc"
	}

	sqlQuery += " ORDER BY " + query.SortBy + " " + query.SortOrder

	if query.Limit > 0 {
		sqlQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", placeholderIndex, placeholderIndex+1)
		args = append(args, query.Limit, (query.Page-1)*query.Limit)
	}

	rows, err := r.conn.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		slog.ErrorContext(ctx, "[productReadRepository] GetListByQuery", "query", err)
		return nil, domain.ErrInternal
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Category, &product.ImageURL, &product.ShopID, &product.Active, &product.CreatedAt, &product.UpdatedAt); err != nil {
			slog.ErrorContext(ctx, "[productReadRepository] GetListByQuery", "scan", err)
			return nil, domain.ErrInternal
		}
		products = append(products, &product)
	}

	return products, nil
}
