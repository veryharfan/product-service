package domain

import (
	"context"

	"github.com/gofrs/uuid/v5"
)

type WarehouseRepository interface {
	GetStock(ctx context.Context, productID uuid.UUID) (int, error)
	FetchStockFromService(ctx context.Context, productID uuid.UUID) (int, error)
	CacheStock(ctx context.Context, productID uuid.UUID, stock int) error
}
