package anilist

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/Khan/genqlient/graphql"
	telemetry "github.com/wwmoraes/gotell"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	"golang.org/x/time/rate"
)

type ratedClient struct {
	client *http.Client
	rater  *rate.Limiter
}

// NewRatedClient creates a GraphQL-compatible HTTP client with rate-limiting
// restrictions. Useful to avoid blacklisting on upstream services.
func NewRatedClient(interval time.Duration, requests int, base *http.Client) graphql.Doer {
	if base == nil {
		base = http.DefaultClient
	}

	return &ratedClient{
		client: base,
		rater:  rate.NewLimiter(rate.Every(interval), requests),
	}
}

// Do executes a HTTP request right away if its within the limits. Otherwise it
// returns a 429 + Retry-After header with the seconds to wait for
func (c *ratedClient) Do(req *http.Request) (*http.Response, error) {
	span := telemetry.SpanFromContext(req.Context())

	re := c.rater.Reserve()
	if !re.OK() {
		err := fmt.Errorf("misconfigured rate limiter on the API, cannot act")
		span.RecordError(err)

		return nil, err
	}

	if re.Delay() > 0 {
		re.Cancel()

		return &http.Response{
			Status:     http.StatusText(http.StatusTooManyRequests),
			StatusCode: http.StatusTooManyRequests,
			Body:       io.NopCloser(bytes.NewBuffer([]byte{})),
			Header: http.Header{
				"Retry-After": []string{strconv.FormatFloat(math.Ceil(re.Delay().Seconds()), 'f', 0, 64)},
			},
		}, nil
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	span.SetAttributes(httpconv.ResponseHeader(telemetry.FilterHeaders(
		resp.Header,
		"X-RateLimit-Remaining",
		"X-RateLimit-Limit",
		"Retry-After",
	))...)

	if resp.StatusCode != http.StatusTooManyRequests {
		return resp, nil
	}

	remaining, err := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	if err != nil {
		return nil, err
	}

	c.rater.SetBurst(remaining)

	burst, err := strconv.Atoi(resp.Header.Get("X-RateLimit-Limit"))
	if err != nil {
		return nil, err
	}

	reset, err := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 0)
	if err != nil {
		return nil, err
	}

	c.rater.SetBurstAt(time.Unix(reset, 0), burst)

	return resp, nil
}
