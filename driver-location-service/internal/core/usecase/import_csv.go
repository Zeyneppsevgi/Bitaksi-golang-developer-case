package usecase

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/driver-location-service/internal/core/domain"
	"github.com/driver-location-service/internal/core/ports"
)

type ImportCSV struct {
	repo      ports.Repository
	chunkSize int
	now       func() time.Time
}

func NewImportCSV(repo ports.Repository, chunkSize int, now func() time.Time) *ImportCSV {
	if chunkSize <= 0 {
		chunkSize = 1000
	}
	if now == nil {
		now = time.Now
	}
	return &ImportCSV{repo: repo, chunkSize: chunkSize, now: now}
}

func (u *ImportCSV) Execute(ctx context.Context, r io.Reader) (int, int, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		return 0, 0, fmt.Errorf("%w: read header: %v", domain.ErrValidation, err)
	}

	if len(headers) == 1 && strings.Contains(headers[0], ";") {
		reader.Comma = ';'
		headers = strings.Split(headers[0], ";")
	}

	idx := mapHeaders(headers)
	for _, col := range []string{"driverid", "longitude", "latitude"} {
		if _, ok := idx[col]; !ok {
			return 0, 0, fmt.Errorf("%w: missing %s column", domain.ErrValidation, col)
		}
	}

	imported := 0
	failed := 0
	batch := make([]domain.DriverLocation, 0, u.chunkSize)

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}
		res, upsertErr := u.repo.BulkUpsertLocations(ctx, batch)
		if upsertErr != nil {
			return fmt.Errorf("%w: %v", domain.ErrInternal, upsertErr)
		}
		imported += int(res.Upserted + res.Updated)
		batch = batch[:0]
		return nil
	}

	for {
		rec, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			failed++
			continue
		}

		driverID := strings.TrimSpace(rec[idx["driverid"]])
		lon, lonErr := strconv.ParseFloat(strings.TrimSpace(rec[idx["longitude"]]), 64)
		lat, latErr := strconv.ParseFloat(strings.TrimSpace(rec[idx["latitude"]]), 64)

		if driverID == "" || lonErr != nil || latErr != nil {
			failed++
			continue
		}

		item := domain.DriverLocation{
			DriverID:  driverID,
			Location:  domain.Point{Type: "Point", Coordinates: []float64{lon, lat}},
			UpdatedAt: u.now().UTC(),
		}

		if err := domain.ValidateLocation(item); err != nil {
			failed++
			continue
		}

		batch = append(batch, item)

		if len(batch) >= u.chunkSize {
			if err := flush(); err != nil {
				return imported, failed, err
			}
		}
	}

	if err := flush(); err != nil {
		return imported, failed, err
	}

	return imported, failed, nil
}

func mapHeaders(headers []string) map[string]int {
	idx := make(map[string]int, len(headers))
	for i, h := range headers {
		clean := strings.ToLower(strings.TrimSpace(strings.Trim(h, "\ufeff")))
		idx[clean] = i
	}
	return idx
}
