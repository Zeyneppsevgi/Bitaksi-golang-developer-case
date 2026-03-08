package httpx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"time"
)

var ErrUpstreamUnavailable = errors.New("upstream unavailable")

type Client struct {
	httpClient *http.Client
	retries    int
	baseDelay  time.Duration
}

func NewClient(timeout time.Duration, retries int, baseDelay time.Duration) *Client {
	return &Client{httpClient: &http.Client{Timeout: timeout}, retries: retries, baseDelay: baseDelay}
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	var last error
	for i := 0; i <= c.retries; i++ {
		r := req.Clone(ctx)
		resp, err := c.httpClient.Do(r)
		if err == nil {
			return resp, nil
		}
		last = err
		if i == c.retries || !isRetryable(err) {
			break
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(c.backoff(i)):
		}
	}
	return nil, fmt.Errorf("%w: %v", ErrUpstreamUnavailable, last)
}

func ReadAndClose(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	return io.ReadAll(body)
}

func isRetryable(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) || errors.Is(err, context.DeadlineExceeded)
}

func (c *Client) backoff(attempt int) time.Duration {
	return time.Duration(float64(c.baseDelay) * math.Pow(2, float64(attempt)))
}
