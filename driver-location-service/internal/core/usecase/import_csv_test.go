package usecase

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/driver-location-service/internal/core/domain"
	"github.com/driver-location-service/internal/core/ports"
)

type mockRepository struct {
	items []domain.DriverLocation
}

func (f *mockRepository) BulkUpsertLocations(_ context.Context, items []domain.DriverLocation) (ports.UpsertResult, error) {
	f.items = append(f.items, items...)
	return ports.UpsertResult{Upserted: int64(len(items)), Updated: 0}, nil
}

func (f *mockRepository) SearchNearest(context.Context, float64, float64, int64, int64) ([]domain.NearestDriver, error) {
	return nil, nil
}

func (f *mockRepository) CreateIndexes(context.Context) error { return nil }
func (f *mockRepository) Ping(context.Context) error          { return nil }

func TestImportCSVExecute(t *testing.T) {
	csv := "driverId,longitude,latitude\n" +
		"d1,29.0,41.0\n" +
		"d2,bad,41.0\n" +
		"d3,29.1,41.1\n"

	repo := &mockRepository{}
	uc := NewImportCSV(repo, 2, func() time.Time { return time.Unix(0, 0) })
	imported, failed, err := uc.Execute(context.Background(), strings.NewReader(csv))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if imported != 2 {
		t.Fatalf("imported=%d want=2", imported)
	}
	if failed != 1 {
		t.Fatalf("failed=%d want=1", failed)
	}
}
