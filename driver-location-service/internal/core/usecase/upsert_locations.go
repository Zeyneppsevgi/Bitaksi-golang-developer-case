package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/driver-location-service/internal/core/domain"
	"github.com/driver-location-service/internal/core/ports"
)

type UpsertLocations struct {
	repo ports.Repository
	now  func() time.Time
}

func NewUpsertLocations(repo ports.Repository, now func() time.Time) *UpsertLocations {
	if now == nil {
		now = time.Now
	}
	return &UpsertLocations{repo: repo, now: now}
}

func (u *UpsertLocations) Execute(ctx context.Context, items []domain.DriverLocation) (ports.UpsertResult, error) {
	if len(items) == 0 {
		return ports.UpsertResult{}, fmt.Errorf("%w: items cannot be empty", domain.ErrValidation)
	}
	for i := range items {
		if err := domain.ValidateLocation(items[i]); err != nil {
			return ports.UpsertResult{}, err
		}
		items[i].UpdatedAt = u.now().UTC()
	}
	res, err := u.repo.BulkUpsertLocations(ctx, items)
	if err != nil {
		return ports.UpsertResult{}, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}
	return res, nil
}
