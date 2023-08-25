package anilist

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	"golang.org/x/time/rate"
)

type ratedClient struct {
	client *http.Client
	rater  *rate.Limiter
}

// NewRatedClient creates a GraphQL-compatible HTTP client with rate-limiting
// restrictions. Its useful to avoid blacklisting and a high rate of errors.
//
// Note: rate limiting is a blocking action. Requests that exceed the limit will
// block and wait for the given interval to proceed. This may cause metric
// distortions.
func NewRatedClient(interval time.Duration, requests int, base *http.Client) graphql.Doer {
	if base == nil {
		base = http.DefaultClient
	}

	return &ratedClient{
		client: base,
		rater:  rate.NewLimiter(rate.Every(interval), requests),
	}
}

// Do executes a HTTP request right away if its within the limits, or blocks
// and waits for the limiter to allow it. See `NewRatedClient` for more info
func (c *ratedClient) Do(req *http.Request) (*http.Response, error) {
	span := telemetry.SpanFromContext(req.Context())
	err := c.rater.Wait(req.Context())
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(httpconv.ResponseHeader(telemetry.WantedRequestHeaders(
		resp.Header,
		"X-RateLimit-Remaining",
		"X-RateLimit-Limit",
		"Retry-After",
	))...)

	if resp.StatusCode != http.StatusTooManyRequests {
		return resp, nil
	}

	//// exceptional case: we wait and retry
	//// first we make sure the rate limiter has the right burst
	remaining, err := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	if err != nil {
		return nil, err
	}
	c.rater.SetBurst(remaining)

	//// this should never happen if the rate limiter is properly set and the API
	//// lives by its documentation
	if remaining != 0 {
		return nil, fmt.Errorf("WARNING inconsistent upstream API - try increasing the rate interval")
	}

	//// update future bursts
	burst, err := strconv.Atoi(resp.Header.Get("X-RateLimit-Limit"))
	if err != nil {
		return nil, err
	}

	reset, err := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 0)
	if err != nil {
		return nil, err
	}

	c.rater.SetBurstAt(time.Unix(reset, 0), burst)

	//// respect the wait time proposed by the API
	after, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 0)
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Duration(after) * time.Second)

	return c.Do(req)
}
