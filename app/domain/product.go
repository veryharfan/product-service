package domain

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"
)

type Product struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"image_url"`
	ShopID      uuid.UUID `json:"shop_id"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductQuery struct {
	ShopID    string `query:"shop_id"`
	Category  string `query:"category"`
	MinPrice  int64  `query:"min_price"`
	MaxPrice  int64  `query:"max_price"`
	Keyword   string `query:"keyword"`
	SortBy    string `query:"sort_by"`
	SortOrder string `query:"sort_order"`
	Page      int    `query:"page"`
	Limit     int    `query:"limit"`
}

type ProductResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"image_url"`
	ShopID      uuid.UUID `json:"shop_id"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductReadRepository interface {
	GetByID(ctx context.Context, id string) (*Product, error)
	GetListByQuery(ctx context.Context, query ProductQuery) ([]*Product, error)
}

type ProductReadUsecase interface {
	GetByID(ctx context.Context, id string) (*ProductResponse, error)
	GetListByQuery(ctx context.Context, query ProductQuery) ([]*Product, error)
}
