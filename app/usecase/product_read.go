package usecase

import (
	"context"
	"log/slog"
	"product-service/app/domain"
	"product-service/config"
)

type productReadUsecase struct {
	productReadRepo domain.ProductReadRepository
	warehouseRepo   domain.WarehouseRepository
	cfg             *config.Config
}

func NewProductReadUsecase(productReadRepo domain.ProductReadRepository, warehouseRepo domain.WarehouseRepository, cfg *config.Config) domain.ProductReadUsecase {
	return &productReadUsecase{productReadRepo, warehouseRepo, cfg}
}

func (u *productReadUsecase) GetByID(ctx context.Context, id int64) (*domain.ProductResponse, error) {
	product, err := u.productReadRepo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "[productReadUsecase] GetByID", "error", err)
		return nil, err
	}
	if product == nil {
		slog.ErrorContext(ctx, "[productReadUsecase] GetByID", "error", domain.ErrNotFound)
		return nil, domain.ErrNotFound
	}

	stock, err := u.warehouseRepo.GetStock(ctx, product.ID)
	if err != nil {
		slog.WarnContext(ctx, "[productReadUsecase] GetStock", "error", err)

		stock, err = u.warehouseRepo.FetchStockFromService(ctx, product.ID)
		if err != nil {
			slog.ErrorContext(ctx, "[productReadUsecase] FetchStockFromService", "error", err)
			return nil, err
		}
		if err := u.warehouseRepo.CacheStock(ctx, product.ID, stock); err != nil {
			slog.ErrorContext(ctx, "[productReadUsecase] CacheStock", "error", err)
			return nil, err
		}
	}

	return &domain.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		ImageURL:    product.ImageURL,
		ShopID:      product.ShopID,
		Stock:       stock,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

func (u *productReadUsecase) GetListByQuery(ctx context.Context, query domain.ProductQuery) ([]*domain.Product, error) {
	products, err := u.productReadRepo.GetListByQuery(ctx, query)
	if err != nil {
		slog.ErrorContext(ctx, "[productReadUsecase] GetListByQuery", "error", err)
		return nil, err
	}
	if products == nil {
		slog.ErrorContext(ctx, "[productReadUsecase] GetListByQuery", "error", domain.ErrNotFound)
		return nil, domain.ErrNotFound
	}

	return products, nil
}
