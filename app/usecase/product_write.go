package usecase

import (
	"context"
	"database/sql"
	"log/slog"
	"product-service/app/domain"
	"product-service/config"
)

type productWriteUsecase struct {
	productReadRepo  domain.ProductReadRepository
	productWriteRepo domain.ProductWriteRepository
	stockRepo        domain.StockRepository
	cfg              *config.Config
}

func NewProductWriteUsecase(productReadRepo domain.ProductReadRepository, productWriteRepo domain.ProductWriteRepository, stockRepo domain.StockRepository, cfg *config.Config) domain.ProductWriteUsecase {
	return &productWriteUsecase{productReadRepo, productWriteRepo, stockRepo, cfg}
}

func (u *productWriteUsecase) Create(ctx context.Context, req *domain.CreateProductRequest) (*domain.CreateProductResponse, error) {
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		ShopID:      req.ShopID,
		Active:      true,
	}

	// Use transaction
	tx, err := u.productWriteRepo.BeginTransaction(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "[productWriteUsecase] Create", "BeginTransaction", err)
		return nil, err
	}

	err = u.productWriteRepo.WithTransaction(ctx, tx, func(ctx context.Context, tx *sql.Tx) error {
		// Create product
		if err := u.productWriteRepo.Create(ctx, product); err != nil {
			slog.ErrorContext(ctx, "[productWriteUsecase] Create", "repository", err)
			return err
		}

		// init stock
		err = u.stockRepo.InitStockToWarehouse(ctx, domain.InitStockRequest{
			ShopID:    product.ShopID,
			ProductID: product.ID,
		})
		if err != nil {
			slog.ErrorContext(ctx, "[productWriteUsecase] Create", "InitStockToWarehouse", err)
			return err
		}
		return nil
	})

	// cache stock
	_ = u.stockRepo.CacheStock(ctx, product.ID, 0)

	slog.InfoContext(ctx, "[productWriteUsecase] success Create", "product_id", product.ID)
	return &domain.CreateProductResponse{
		ID: product.ID,
	}, nil
}

func (u *productWriteUsecase) Update(ctx context.Context, id int64, req *domain.UpdateProductRequest) (*domain.Product, error) {
	product, err := u.productReadRepo.GetByID(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "[productWriteUsecase] Update", "GetByID", err)
		return nil, err
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Category = req.Category
	product.ImageURL = req.ImageURL
	product.ShopID = req.ShopID

	if err := u.productWriteRepo.Update(ctx, product); err != nil {
		slog.ErrorContext(ctx, "[productWriteUsecase] Update", "Update", err)
		return nil, err
	}

	slog.InfoContext(ctx, "[productWriteUsecase] success Update", "product_id", product.ID)
	return product, nil
}

func (u *productWriteUsecase) SetActiveStatus(ctx context.Context, id int64, active bool) error {
	if err := u.productWriteRepo.SetActiveStatus(ctx, id, active); err != nil {
		slog.ErrorContext(ctx, "[productWriteUsecase] SetActiveStatus", "SetActiveStatus", err)
		return err
	}
	slog.InfoContext(ctx, "[productWriteUsecase] success SetActiveStatus", "product_id", id)

	return nil
}
