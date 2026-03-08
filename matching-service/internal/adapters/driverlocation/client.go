package driverlocation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/matching-service/internal/core/domain"
	"github.com/matching-service/pkg/httpx"
	"github.com/sony/gobreaker"
)

type breaker interface {
	Execute(func() (interface{}, error)) (interface{}, error)
}

type Client struct {
	baseURL     string
	internalKey string
	httpClient  *httpx.Client
	breaker     breaker
}

func NewClient(baseURL, internalKey string) *Client {
	return &Client{
		baseURL:     baseURL,
		internalKey: internalKey,
		httpClient:  httpx.NewClient(2*time.Second, 2, 100*time.Millisecond),
		breaker:     newCircuitBreaker(),
	}
}

func (c *Client) SearchNearest(ctx context.Context, lon, lat float64, radiusM int64, requestID string) ([]domain.MatchResult, error) {
	res, err := c.breaker.Execute(func() (interface{}, error) {
		return c.search(ctx, lon, lat, radiusM, requestID)
	})
	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
			return nil, domain.ErrUpstreamUnavailable
		}
		return nil, err
	}
	items, _ := res.([]domain.MatchResult)
	return items, nil
}

func (c *Client) search(ctx context.Context, lon, lat float64, radiusM int64, requestID string) ([]domain.MatchResult, error) {
	u, err := url.Parse(c.baseURL + "/v1/driver-locations/search")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("lon", strconv.FormatFloat(lon, 'f', -1, 64))
	q.Set("lat", strconv.FormatFloat(lat, 'f', -1, 64))
	q.Set("radius_m", strconv.FormatInt(radiusM, 10))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Internal-Api-Key", c.internalKey)
	if requestID != "" {
		req.Header.Set("X-Request-Id", requestID)
	}

	resp, err := c.httpClient.Do(ctx, req)
	if err != nil {
		return nil, domain.ErrUpstreamUnavailable
	}
	body, err := httpx.ReadAndClose(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 500 || resp.StatusCode == http.StatusUnauthorized {
		return nil, domain.ErrUpstreamUnavailable
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var out struct {
		Data struct {
			Items []struct {
				DriverID  string  `json:"driverId"`
				DistanceM float64 `json:"distanceM"`
			} `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	res := make([]domain.MatchResult, 0, len(out.Data.Items))
	for _, item := range out.Data.Items {
		res = append(res, domain.MatchResult{DriverID: item.DriverID, DistanceM: item.DistanceM})
	}
	return res, nil
}
