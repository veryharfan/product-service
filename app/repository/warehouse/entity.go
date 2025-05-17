package warehouse

import "github.com/gofrs/uuid/v5"

type ProductStockResponse struct {
	ProductID   uuid.UUID `json:"product_id"`
	WarehouseID uuid.UUID `json:"warehouse_id"`
	Quantity    int       `json:"quantity"`
	Reserved    int       `json:"reserved"`
}
