package telemetry

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/XSAM/otelsql"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

func StartFunction(ctx context.Context, opts ...trace.SpanStartOption) (context.Context, Span) {
	name, opt := functionInfo(2)
	opts = append(opts, opt)

	return globalTracer.Start(ctx, name, opts...)
}

func Start(ctx context.Context, spanName string, opts ...SpanStartOption) (context.Context, Span) {
	return globalTracer.Start(ctx, spanName, opts...)
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func WantedRequestHeaders(h http.Header, keys ...string) http.Header {
	target := http.Header{}

	for _, key := range keys {
		target[key] = h.Values(key)
	}

	return target
}

func OpenSQL(driverName, dataSourceName string) (*sql.DB, error) {
	attributes := otelsql.WithAttributes(
		attribute.String(string(semconv.DBSystemKey), driverName),
	)

	db, err := otelsql.Open(driverName, dataSourceName, attributes)
	if err != nil {
		return nil, err
	}

	err = otelsql.RegisterDBStatsMetrics(db, attributes)
	if err != nil {
		return nil, err
	}

	return db, nil
}
