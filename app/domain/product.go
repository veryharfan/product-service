package domain

import (
	"context"
	"database/sql"
	"time"
)

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"image_url"`
	ShopID      int64     `json:"shop_id"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductQuery struct {
	ShopID    int64  `query:"shop_id"`
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
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"image_url"`
	ShopID      int64     `json:"shop_id"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Price       int64  `json:"price" validate:"required"`
	Category    string `json:"category" validate:"required"`
	ImageURL    string `json:"image_url" validate:"required"`
}

type CreateProductResponse struct {
	ID int64 `json:"id"`
}

type UpdateProductRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Price       int64  `json:"price" validate:"required"`
	Category    string `json:"category" validate:"required"`
	ImageURL    string `json:"image_url" validate:"required"`
	ShopID      int64  `json:"shop_id" validate:"required"`
	Active      bool   `json:"active" validate:"required"`
}

type SetActiveStatusRequest struct {
	Active bool `json:"active"`
}

type ProductReadRepository interface {
	GetByID(ctx context.Context, id int64) (*Product, error)
	GetListByQuery(ctx context.Context, query ProductQuery) ([]*Product, error)
}

type ProductReadUsecase interface {
	GetByID(ctx context.Context, id int64) (*ProductResponse, error)
	GetListByQuery(ctx context.Context, query ProductQuery) ([]*Product, error)
}

type ProductWriteRepository interface {
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	SetActiveStatus(ctx context.Context, id int64, active bool) error

	BeginTransaction(ctx context.Context) (*sql.Tx, error)
	WithTransaction(ctx context.Context, tx *sql.Tx, fn func(context.Context, *sql.Tx) error) error
}

type ProductWriteUsecase interface {
	Create(ctx context.Context, shopID int64, product *CreateProductRequest) (*CreateProductResponse, error)
	Update(ctx context.Context, id int64, product *UpdateProductRequest) (*Product, error)
	SetActiveStatus(ctx context.Context, id int64, active bool) error
}
