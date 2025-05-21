package stockrepo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"product-service/app/domain"
	"product-service/pkg"
	"time"

	"github.com/redis/go-redis/v9"
)

type stockRepository struct {
	redis              *redis.Client
	ttl                time.Duration
	httpClient         *http.Client
	baseURL            string
	internalAuthHeader string
}

func NewStockRepository(redis *redis.Client, ttl time.Duration, baseURL string, internalAuthHeader string) domain.StockRepository {
	return &stockRepository{
		redis:              redis,
		ttl:                ttl,
		httpClient:         &http.Client{Timeout: 30 * time.Second},
		baseURL:            baseURL,
		internalAuthHeader: internalAuthHeader,
	}
}

func (r *stockRepository) GetStock(ctx context.Context, productID int64) (int, error) {
	stock, err := r.redis.Get(ctx, r.key(productID)).Int()
	if err != nil {
		slog.WarnContext(ctx, "[GetStock] Cache miss or error retrieving stock", "productID", productID, "error", err)
		return 0, fmt.Errorf("cache miss: %w", err)
	}
	return stock, nil
}

func (r *stockRepository) FetchStockFromService(ctx context.Context, productID int64) (int, error) {
	url := fmt.Sprintf("%s/internal/warehouse-service/products/%d/stocks", r.baseURL, productID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		slog.ErrorContext(ctx, "[FetchStockFromService] Failed to create HTTP request", "productID", productID, "error", err)
		return 0, err
	}

	pkg.AddRequestHeader(ctx, r.internalAuthHeader, req)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "[FetchStockFromService] HTTP request failed", "productID", productID, "error", err)
		return 0, err
	}
	defer resp.Body.Close()

	var data []ProductStockResponse
	if err := pkg.DecodeResponseBody(resp, &data); err != nil {
		slog.ErrorContext(ctx, "[FetchStockFromService] Failed to decode response body", "productID", productID, "error", err)
		return 0, err
	}

	var stock int64
	for _, item := range data {
		stock += item.Available
	}
	slog.InfoContext(ctx, "[FetchStockFromService] Stock fetched from warehouse service", "productID", productID, "stock", stock)
	return int(stock), nil
}

func (r *stockRepository) CacheStock(ctx context.Context, productID int64, stock int) error {
	err := r.redis.Set(ctx, r.key(productID), stock, r.ttl).Err()
	if err != nil {
		slog.ErrorContext(ctx, "[CacheStock] Failed to cache stock", "productID", productID, "stock", stock, "error", err)
		return err
	}
	slog.InfoContext(ctx, "[CacheStock] Stock cached successfully", "productID", productID, "stock", stock)
	return nil
}

func (r *stockRepository) key(productID int64) string {
	return fmt.Sprintf("stock:product:%d", productID)
}

func (r *stockRepository) InitStockToWarehouse(ctx context.Context, warehouse domain.InitStockRequest) error {
	url := fmt.Sprintf("%s/internal/warehouse-service/stocks", r.baseURL)
	reqBody, err := json.Marshal(warehouse)
	if err != nil {
		slog.ErrorContext(ctx, "[stockRepository] InitStockToWarehouse", "json.Marshal", err)
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		slog.ErrorContext(ctx, "[stockRepository] InitStockToWarehouse", "http.NewRequestWithContext", err)
		return err
	}

	pkg.AddRequestHeader(ctx, r.internalAuthHeader, httpReq)

	resp, err := r.httpClient.Do(httpReq)
	if err != nil {
		slog.ErrorContext(ctx, "[stockRepository] InitStockToWarehouse", "httpClient.Do", err)
		return err
	}
	defer resp.Body.Close()

	var res any
	if err := pkg.DecodeResponseBody(resp, &res); err != nil {
		slog.ErrorContext(ctx, "[stockRepository] InitStockToWarehouse", "DecodeResponseBody", err)
		return err
	}

	return nil
}
