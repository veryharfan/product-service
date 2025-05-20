package handler

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"
)

func SetupConsumer(ctx context.Context, stream jetstream.Stream, stockConsumerHandler *StockConsumerHandler) {

	cons, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable: "processor",
	})
	if err != nil {
		slog.ErrorContext(ctx, "[SetupConsumer] CreateOrUpdateConsumer", "error", err)
		return
	}

	cons.Consume(func(msg jetstream.Msg) {
		if msg.Subject() != "stock.available" {
			return
		}
		stockConsumerHandler.UpdateStock(msg)
	})
	slog.InfoContext(ctx, "[SetupConsumer] Consumer setup successfully", "subject", "stock.available")
}
