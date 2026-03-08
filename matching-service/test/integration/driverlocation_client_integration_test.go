//go:build integration

package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	driverlocationadapter "github.com/matching-service/internal/adapters/driverlocation"
)

func TestDriverLocationClientSendsInternalAPIKey(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Internal-Api-Key"); got != "some-secret" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"ok":false}`))
			return
		}
		_, _ = w.Write([]byte(`{"ok":true,"data":{"items":[{"driverId":"driver-1","distanceM":12.5}]}}`))
	}))
	defer ts.Close()

	client := driverlocationadapter.NewClient(ts.URL, "some-secret")
	items, err := client.SearchNearest(context.Background(), 29, 41, 3000, "req-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0].DriverID != "driver-1" {
		t.Fatalf("unexpected items: %+v", items)
	}
}
