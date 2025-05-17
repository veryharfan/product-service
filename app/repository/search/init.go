package search

import (
	"context"
	"log/slog"

	"github.com/meilisearch/meilisearch-go"
)

// InitMeiliSearch initializes the MeiliSearch client
func InitMeiliSearch(host, apiKey string) (*meilisearch.ServiceManager, error) {
	meiliClient := meilisearch.New(host, meilisearch.WithAPIKey(apiKey))

	// Test the connection
	_, err := meiliClient.Health()
	if err != nil {
		slog.WarnContext(context.Background(), "[InitMeiliSearch] Health check failed", "error", err)
		return nil, err
	}

	slog.InfoContext(context.Background(), "[InitMeiliSearch] MeiliSearch client initialized successfully")
	return &meiliClient, nil
}
