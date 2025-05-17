package warehouse

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/redis/go-redis/v9"
)

type warehouseRepository struct {
	redis      *redis.Client
	httpClient *http.Client
	baseURL    string
	ttl        time.Duration
}

func NewWarehouseRepository(redis *redis.Client, baseURL string, ttl time.Duration) *warehouseRepository {
	return &warehouseRepository{
		redis:      redis,
		baseURL:    baseURL,
		ttl:        ttl,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (r *warehouseRepository) GetStock(ctx context.Context, productID uuid.UUID) (int, error) {
	stock, err := r.redis.Get(ctx, r.key(productID)).Int()
	if err != nil {
		slog.WarnContext(ctx, "[GetStock] Cache miss or error retrieving stock", "productID", productID, "error", err)
		return 0, fmt.Errorf("cache miss: %w", err)
	}
	return stock, nil
}

func (r *warehouseRepository) FetchStockFromService(ctx context.Context, productID uuid.UUID) (int, error) {
	url := fmt.Sprintf("%s/warehouse/stock/%s", r.baseURL, productID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		slog.ErrorContext(ctx, "[FetchStockFromService] Failed to create HTTP request", "productID", productID, "error", err)
		return 0, err
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "[FetchStockFromService] HTTP request failed", "productID", productID, "error", err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.WarnContext(ctx, "[FetchStockFromService] Non-OK response from warehouse service", "productID", productID, "statusCode", resp.StatusCode)
		return 0, fmt.Errorf("warehouse-service returned %d", resp.StatusCode)
	}

	var data []ProductStockResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		slog.ErrorContext(ctx, "[FetchStockFromService] Failed to decode response body", "productID", productID, "error", err)
		return 0, err
	}

	var stock int
	for _, item := range data {
		stock += item.Quantity - item.Reserved
	}
	slog.InfoContext(ctx, "[FetchStockFromService] Stock fetched from warehouse service", "productID", productID, "stock", stock)
	return stock, nil
}

func (r *warehouseRepository) CacheStock(ctx context.Context, productID uuid.UUID, stock int) error {
	err := r.redis.Set(ctx, r.key(productID), stock, r.ttl).Err()
	if err != nil {
		slog.ErrorContext(ctx, "[CacheStock] Failed to cache stock", "productID", productID, "stock", stock, "error", err)
		return err
	}
	slog.InfoContext(ctx, "[CacheStock] Stock cached successfully", "productID", productID, "stock", stock)
	return nil
}

func (r *warehouseRepository) key(productID uuid.UUID) string {
	return fmt.Sprintf("stock:product:%s", productID.String())
}
