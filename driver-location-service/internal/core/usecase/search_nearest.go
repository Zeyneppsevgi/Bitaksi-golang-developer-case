package usecase

import (
	"context"
	"fmt"

	"github.com/driver-location-service/internal/core/domain"
	"github.com/driver-location-service/internal/core/ports"
	"github.com/driver-location-service/pkg/validation"
)

type SearchNearest struct {
	repo ports.Repository
}

func NewSearchNearest(repo ports.Repository) *SearchNearest {
	return &SearchNearest{repo: repo}
}

func (u *SearchNearest) Execute(ctx context.Context, lon, lat float64, radiusM int64, limit int64) ([]domain.NearestDriver, error) {
	if err := validation.ValidateCoordinates(lon, lat); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrValidation, err)
	}
	if radiusM <= 0 {
		return nil, fmt.Errorf("%w: radius_m must be > 0", domain.ErrValidation)
	}
	if limit <= 0 {
		limit = 100
	}
	items, err := u.repo.SearchNearest(ctx, lon, lat, radiusM, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	return items, nil
}
