package main

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"golang.org/x/time/rate"
)

func Limiter(limiter *rate.Limiter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := telemetry.SpanFromContext(r.Context())

			re := limiter.Reserve()
			if !re.OK() {
				err := fmt.Errorf("misconfigured rate limiter on the API, cannot act")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				span.RecordError(err)
				return
			}

			if re.Delay() > 0 {
				re.Cancel()
				w.Header().Add("Retry-After", strconv.FormatFloat(math.Ceil(re.Delay().Seconds()), 'f', 0, 64))
				w.WriteHeader(http.StatusTooManyRequests)
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
