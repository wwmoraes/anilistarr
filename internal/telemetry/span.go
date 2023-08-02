package telemetry

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Span interface {
	trace.Span

	EndWith(err error, options ...trace.SpanEndOption)
	Assert(error) error
	Int(k string, v int)
}

type span struct {
	trace.Span
}

func (s *span) Assert(err error) error {
	if err == nil {
		s.SetStatus(codes.Ok, "")
	} else {
		s.SetStatus(codes.Error, err.Error())
		s.RecordError(err)
	}

	return err
}

func (s *span) EndWith(err error, options ...trace.SpanEndOption) {
	if err == nil {
		s.SetStatus(codes.Ok, "")
	} else {
		s.SetStatus(codes.Error, err.Error())
		s.RecordError(err)
	}

	s.End(options...)
}

func (s *span) Int(k string, v int) {
	s.SetAttributes(attribute.Int(k, v))
}

func WithInt(k string, v int) trace.SpanStartEventOption {
	return trace.WithAttributes(
		attribute.Int(k, v),
	)
}

func WithSpanKindClient() trace.SpanStartOption {
	return trace.WithSpanKind(trace.SpanKindClient)
}
