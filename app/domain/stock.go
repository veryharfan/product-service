package domain

import (
	"context"
)

type StockMessage struct {
	ProductID int64 `json:"product_id"`
	Available int   `json:"available"`
}

type StockRepository interface {
	GetStock(ctx context.Context, productID int64) (int, error)
	FetchStockFromService(ctx context.Context, productID int64) (int, error)
	CacheStock(ctx context.Context, productID int64, stock int) error
}

type StockUsecase interface {
	UpdateStock(ctx context.Context, msg StockMessage) error
}
