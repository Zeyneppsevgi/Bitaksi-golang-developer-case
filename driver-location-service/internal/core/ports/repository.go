package ports

import (
	"context"

	"github.com/driver-location-service/internal/core/domain"
)

type UpsertResult struct {
	Upserted int64
	Updated  int64
}

type Repository interface {
	BulkUpsertLocations(ctx context.Context, items []domain.DriverLocation) (UpsertResult, error)
	SearchNearest(ctx context.Context, lon, lat float64, radiusM int64, limit int64) ([]domain.NearestDriver, error)
	CreateIndexes(ctx context.Context) error
	Ping(ctx context.Context) error
}
