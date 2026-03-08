package usecase

import (
	"context"
	"fmt"

	"github.com/matching-service/internal/core/domain"
	"github.com/matching-service/internal/core/ports"
	"github.com/matching-service/pkg/validation"
)

type FindDriver struct {
	client ports.DriverLocationClient
}

func NewFindDriver(client ports.DriverLocationClient) *FindDriver {
	return &FindDriver{client: client}
}

// En yakınları getir
func (u *FindDriver) Execute(ctx context.Context, lon, lat float64, radiusM int64, requestID string) (domain.MatchResult, error) {
	if err := validation.ValidateCoordinates(lon, lat); err != nil {
		return domain.MatchResult{}, fmt.Errorf("%w: %v", domain.ErrValidation, err)
	}
	if radiusM <= 0 {
		return domain.MatchResult{}, fmt.Errorf("%w: radius_m must be > 0", domain.ErrValidation)
	}
	items, err := u.client.SearchNearest(ctx, lon, lat, radiusM, requestID)
	if err != nil {
		if err == domain.ErrUpstreamUnavailable {
			return domain.MatchResult{}, err
		}
		return domain.MatchResult{}, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	if len(items) == 0 {
		return domain.MatchResult{}, domain.ErrNotFound
	}
	return items[0], nil
}
