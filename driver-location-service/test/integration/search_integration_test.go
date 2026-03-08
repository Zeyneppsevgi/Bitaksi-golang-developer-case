//go:build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	mongoadapter "github.com/driver-location-service/internal/adapters/mongo"
	"github.com/driver-location-service/internal/core/domain"
	"github.com/driver-location-service/internal/core/usecase"
)

func TestSearchNearestWithMongo(t *testing.T) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongoadapter.NewClient(ctx, uri)
	if err != nil {
		t.Skipf("mongo unavailable: %v", err)
	}
	defer client.Disconnect(ctx)

	repo := mongoadapter.NewRepository(client, "driver_location_it", "driver_locations")
	if err := repo.CreateIndexes(ctx); err != nil {
		t.Fatalf("create indexes: %v", err)
	}

	upsert := usecase.NewUpsertLocations(repo, time.Now)
	_, err = upsert.Execute(ctx, []domain.DriverLocation{
		{DriverID: "near", Location: domain.Point{Type: "Point", Coordinates: []float64{29.0001, 41.0001}}},
		{DriverID: "far", Location: domain.Point{Type: "Point", Coordinates: []float64{29.01, 41.01}}},
	})
	if err != nil {
		t.Fatalf("upsert: %v", err)
	}

	search := usecase.NewSearchNearest(repo)
	items, err := search.Execute(ctx, 29, 41, 3000, 10)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(items) < 2 {
		t.Fatalf("expected at least 2 items")
	}
	if items[0].DriverID != "near" {
		t.Fatalf("expected near first, got %s", items[0].DriverID)
	}
	if items[0].DistanceM <= 0 {
		t.Fatalf("expected positive distance")
	}
}
