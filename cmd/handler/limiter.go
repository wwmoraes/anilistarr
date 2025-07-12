package main

import (
	"math"
	"net/http"
	"strconv"

	telemetry "github.com/wwmoraes/gotell"
	"golang.org/x/time/rate"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

// Limiter provides an HTTP middleware that limits requests forwarded to
// the next handler.
//
// Over-limit requests get an immediate response with a status 429 Too Many
// Requests and a Retry-After header instead of hanging the connection. This
// prevents slow HTTP attacks.
//
// TODO rename to LimitWith
func Limiter(limiter *rate.Limiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := telemetry.SpanFromContext(r.Context())

			reservation := limiter.Reserve()
			if !reservation.OK() {
				http.Error(w, usecases.ErrStatusInternal.Error(), http.StatusInternalServerError)
				span.RecordError(usecases.ErrStatusInternal)

				return
			}

			if reservation.Delay() > 0 {
				reservation.Cancel()
				w.Header().
					Add("Retry-After", strconv.FormatFloat(math.Ceil(reservation.Delay().Seconds()), 'f', 0, 64))
				w.WriteHeader(http.StatusTooManyRequests)
				http.Error(
					w,
					usecases.ErrStatusResourceExhausted.Error(),
					http.StatusTooManyRequests,
				)
				span.RecordError(usecases.ErrStatusResourceExhausted)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
