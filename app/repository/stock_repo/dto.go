package stockrepo

type ProductStockResponse struct {
	ProductID   int64 `json:"product_id"`
	WarehouseID int64 `json:"warehouse_id"`
	Quantity    int64 `json:"quantity"`
	Reserved    int64 `json:"reserved"`
	Available   int64 `json:"available"`
}
