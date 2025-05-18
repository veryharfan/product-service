package domain

import (
	"context"
)

type WarehouseRepository interface {
	GetStock(ctx context.Context, productID int64) (int, error)
	FetchStockFromService(ctx context.Context, productID int64) (int, error)
	CacheStock(ctx context.Context, productID int64, stock int) error
}
