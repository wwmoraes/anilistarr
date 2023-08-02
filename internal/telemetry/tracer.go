package telemetry

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"runtime"

	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	fnNameRE = regexp.MustCompile(`.*?(?:\(\*)?([^\./\(\)\[\]]+)(?:[\[\]\.\)]*)?\.([^\./\(\)]+)$`)
)

type TracerOption = trace.TracerOption
type SpanStartOption = trace.SpanStartOption

type Tracer interface {
	StartFunction(ctx context.Context, opts ...SpanStartOption) (context.Context, Span)
	StartHTTPResponse(req *http.Request, opts ...SpanStartOption) (context.Context, HTTPSpan)

	// custom span
	// from trace.Tracer
	Start(ctx context.Context, spanName string, opts ...SpanStartOption) (context.Context, Span)
}

type tracer struct {
	upstream trace.Tracer
}

func newTracer(opts ...TracerOption) Tracer {
	return &tracer{
		upstream: otel.Tracer(NAME, opts...),
	}
}

func DefaultTracer() Tracer {
	return globalTracer
}

func (t *tracer) Start(ctx context.Context, spanName string, opts ...SpanStartOption) (context.Context, Span) {
	ctx, upstreamSpan := t.upstream.Start(ctx, spanName, opts...)
	return ctx, &span{upstreamSpan}
}

func (t *tracer) StartFunction(ctx context.Context, opts ...SpanStartOption) (context.Context, Span) {
	name, opt := functionInfo(1)
	opts = append(opts, opt)

	return t.Start(ctx, name, opts...)
}

func (t *tracer) StartHTTPResponse(req *http.Request, opts ...SpanStartOption) (context.Context, HTTPSpan) {
	opts = append(opts, trace.WithAttributes(
		semconv.HTTPMethod(req.Method),
		semconv.HTTPURL(req.URL.String()),
		semconv.HTTPRoute(req.URL.Path),
		semconv.HTTPScheme(req.URL.Scheme),
		semconv.HTTPClientIP(req.RemoteAddr),
		semconv.HTTPRequestContentLength(int(req.ContentLength)),
	))

	upstreamCtx, upstreamSpan := t.upstream.Start(req.Context(), fmt.Sprintf("%s %s", req.Method, req.URL.Path), opts...)

	return upstreamCtx, &httpSpan{span{upstreamSpan}}
}

func functionInfo(skip int) (string, trace.SpanStartOption) {
	name, fullName := "", ""
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		details := runtime.FuncForPC(pc)
		if details != nil {
			fullName = details.Name()
			name = fnNameRE.ReplaceAllString(details.Name(), "$1.$2")
			file, line = details.FileLine(pc)
		}
	} else {
		name = file
		fullName = file
	}

	return name, trace.WithAttributes(
		semconv.CodeFunction(fullName),
		semconv.CodeFilepath(file),
		semconv.CodeLineNumber(line),
	)
}
