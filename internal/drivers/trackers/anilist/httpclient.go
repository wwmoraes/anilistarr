package anilist

import (
	"context"
	"fmt"
	"maps"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/Khan/genqlient/graphql"
	telemetry "github.com/wwmoraes/gotell"
	"go.opentelemetry.io/otel/semconv/v1.20.0/httpconv"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/time/rate"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var _ graphql.Doer = (*RatedClient)(nil)

// RatedClient is a rate-limited HTTP client. This allows consuming upstream
// resources with usage limits in a friendly way.
type RatedClient struct {
	usecases.Doer
	Limiter *rate.Limiter
}

// Do executes a HTTP request right away if its within the limits. Otherwise it
// returns a 429 + Retry-After header with the seconds to wait for
func (client *RatedClient) Do(req *http.Request) (*http.Response, error) {
	span := telemetry.SpanFromContext(req.Context())

	reservation := client.Limiter.Reserve()
	if !reservation.OK() {
		return nil, span.Assert(usecases.ErrStatusInternal)
	}

	if reservation.Delay() > 0 {
		reservation.Cancel()

		res := newResponseFor(req, http.StatusTooManyRequests, nil, http.Header{
			"Retry-After": []string{strconv.FormatFloat(math.Ceil(reservation.Delay().Seconds()), 'f', 0, 64)},
		})

		return res, nil
	}

	resp, err := client.Doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", usecases.ErrStatusUnknown, err)
	}

	span.SetAttributes(httpconv.ResponseHeader(telemetry.FilterHeaders(
		resp.Header,
		"Retry-After",
		"X-Ratelimit-Limit",
		"X-Ratelimit-Remaining",
		"X-Ratelimit-Reset",
	))...)

	if resp.StatusCode == http.StatusTooManyRequests {
		tryUpdateLimiterBurstFromHeaders(req.Context(), client.Limiter, resp.Header)
		tryUpdateLimiterBurstAtFromHeaders(req.Context(), client.Limiter, resp.Header)
	}

	return resp, span.Assert(nil)
}

func newResponseFor(req *http.Request, status int, data []byte, headers http.Header) *http.Response {
	span := telemetry.SpanFromContext(req.Context())

	writer := httptest.NewRecorder()
	maps.Copy(writer.Header(), headers)
	writer.WriteHeader(status)

	_, err := writer.Write(data)
	if err != nil {
		span.RecordError(err, trace.WithStackTrace(true))
	}

	res := writer.Result()
	res.Request = req

	return res
}

func tryUpdateLimiterBurstFromHeaders(ctx context.Context, limiter *rate.Limiter, headers http.Header) {
	span := telemetry.SpanFromContext(ctx)

	value := headers.Get("X-Ratelimit-Remaining")
	if value == "" {
		return
	}

	remaining, err := strconv.Atoi(value)
	if err != nil {
		span.RecordError(fmt.Errorf("failed to parse remaining limit: %w", err))

		return
	}

	limiter.SetBurst(remaining)
}

func tryUpdateLimiterBurstAtFromHeaders(ctx context.Context, limiter *rate.Limiter, headers http.Header) {
	span := telemetry.SpanFromContext(ctx)

	burstValue := headers.Get("X-Ratelimit-Limit")
	if burstValue == "" {
		return
	}

	resetValue := headers.Get("X-Ratelimit-Reset")
	if resetValue == "" {
		return
	}

	burst, err := strconv.Atoi(burstValue)
	if err != nil {
		span.RecordError(fmt.Errorf("failed to parse rate limit burst: %w", err))

		return
	}

	reset, err := strconv.ParseInt(resetValue, 10, 0)
	if err != nil {
		span.RecordError(fmt.Errorf("failed to parse rate limit reset time: %w", err))

		return
	}

	limiter.SetBurstAt(time.Unix(reset, 0), burst)
}
