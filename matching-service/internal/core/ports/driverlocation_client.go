package ports

import (
	"context"

	"github.com/matching-service/internal/core/domain"
)

type DriverLocationClient interface {
	SearchNearest(ctx context.Context, lon, lat float64, radiusM int64, requestID string) ([]domain.MatchResult, error)
}
