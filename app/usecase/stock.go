package usecase

import (
	"context"
	"log/slog"
	"product-service/app/domain"
	"product-service/config"
)

type stockUsecase struct {
	stockRepository domain.StockRepository
	cfg             *config.Config
}

func NewStockUsecase(stockRepository domain.StockRepository, cfg *config.Config) domain.StockUsecase {
	return &stockUsecase{
		stockRepository: stockRepository,
		cfg:             cfg,
	}
}

func (u *stockUsecase) UpdateStock(ctx context.Context, msg domain.StockMessage) error {
	// Cache the stock
	if err := u.stockRepository.CacheStock(ctx, msg.ProductID, msg.Available); err != nil {
		slog.ErrorContext(ctx, "[stockUsecase] UpdateStock", "cacheStock", err)
		return err
	}

	return nil
}
