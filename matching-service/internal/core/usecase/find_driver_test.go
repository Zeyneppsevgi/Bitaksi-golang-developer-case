package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/matching-service/internal/core/domain"
)

type fakeClient struct {
	items []domain.MatchResult
	err   error
}

func (f *fakeClient) SearchNearest(context.Context, float64, float64, int64, string) ([]domain.MatchResult, error) {
	return f.items, f.err
}

func TestFindDriverExecute(t *testing.T) {
	uc := NewFindDriver(&fakeClient{items: []domain.MatchResult{{DriverID: "d1", DistanceM: 1.2}}})
	res, err := uc.Execute(context.Background(), 29, 41, 3000, "req-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.DriverID != "d1" {
		t.Fatalf("unexpected result: %+v", res)
	}
}

func TestFindDriverNotFound(t *testing.T) {
	uc := NewFindDriver(&fakeClient{items: []domain.MatchResult{}})
	_, err := uc.Execute(context.Background(), 29, 41, 3000, "req-1")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected not found, got %v", err)
	}
}
