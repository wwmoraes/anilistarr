package telemetry

import (
	"net/http"

	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

type HTTPSpan interface {
	Span

	HTTPStatus(status int)
	EndWithStatus(status int, options ...trace.SpanEndOption)
}

type httpSpan struct {
	span
}

func (s *httpSpan) HTTPStatus(status int) {
	code := codes.Ok
	if status >= 400 {
		code = codes.Error
	}

	s.SetAttributes(semconv.HTTPStatusCode(status))
	s.SetStatus(code, http.StatusText(status))
}

func (s *httpSpan) EndWithStatus(status int, options ...trace.SpanEndOption) {
	code := codes.Ok
	if status >= 400 {
		code = codes.Error
	}

	s.SetAttributes(semconv.HTTPStatusCode(status))
	s.SetStatus(code, http.StatusText(status))
	s.End(options...)
}
