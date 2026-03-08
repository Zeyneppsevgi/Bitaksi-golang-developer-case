package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/driver-location-service/internal/core/domain"
	"github.com/driver-location-service/internal/core/ports"
)

type upsertRepo struct{}

func (u *upsertRepo) BulkUpsertLocations(context.Context, []domain.DriverLocation) (ports.UpsertResult, error) {
	return ports.UpsertResult{Upserted: 1, Updated: 0}, nil
}

func (u *upsertRepo) SearchNearest(context.Context, float64, float64, int64, int64) ([]domain.NearestDriver, error) {
	return nil, nil
}
func (u *upsertRepo) CreateIndexes(context.Context) error { return nil }
func (u *upsertRepo) Ping(context.Context) error          { return nil }

func TestUpsertLocationsExecute(t *testing.T) {
	uc := NewUpsertLocations(&upsertRepo{}, func() time.Time { return time.Unix(0, 0) })
	_, err := uc.Execute(context.Background(), []domain.DriverLocation{{
		DriverID: "d1",
		Location: domain.Point{Type: "Point", Coordinates: []float64{29, 41}},
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpsertLocationsValidation(t *testing.T) {
	uc := NewUpsertLocations(&upsertRepo{}, time.Now)
	_, err := uc.Execute(context.Background(), []domain.DriverLocation{{
		DriverID: "",
		Location: domain.Point{Type: "Point", Coordinates: []float64{29, 41}},
	}})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}
