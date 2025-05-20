package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"product-service/app/domain"

	"github.com/nats-io/nats.go/jetstream"
)

type StockConsumerHandler struct {
	stockUsecase domain.StockUsecase
}

func NewStockConsumerHandler(stockUsecase domain.StockUsecase) *StockConsumerHandler {
	return &StockConsumerHandler{
		stockUsecase: stockUsecase,
	}
}

func (h *StockConsumerHandler) UpdateStock(msg jetstream.Msg) {
	// Unmarshal the message
	var stockMsg domain.StockMessage
	if err := json.Unmarshal(msg.Data(), &stockMsg); err != nil {
		slog.Error("[HandleStockMessage] Unmarshal failed", "error", err)
		return
	}

	if err := h.stockUsecase.UpdateStock(context.Background(), stockMsg); err != nil {
		slog.Error("[HandleStockMessage] UpdateStock failed", "error", err)
		return
	}

	msg.Ack()

	slog.Info("[HandleStockMessage] Stock updated successfully", "stock", stockMsg)

	return
}
