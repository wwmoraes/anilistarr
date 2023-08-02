package telemetry

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

func NewHandler(handler http.Handler, operation string) http.Handler {
	return otelhttp.NewHandler(handler, operation)
}

func NewHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := globalTracer.StartHTTPResponse(r)
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r.WithContext(ctx))
		span.EndWithStatus(ww.Status())
	})
}

func NewHandleFunc(fn http.HandlerFunc, operation string) http.Handler {
	return NewHandler(fn, operation)
}

type responseWriterSnooper struct {
	w          http.ResponseWriter
	statusCode int
}

func (ws *responseWriterSnooper) WriteHeader(statusCode int) {
	ws.statusCode = statusCode
	ws.w.WriteHeader(statusCode)
}

func (ws *responseWriterSnooper) Header() http.Header {
	return ws.w.Header()
}
func (ws *responseWriterSnooper) Write(data []byte) (int, error) {
	return ws.w.Write(data)
}

func HandlerFunc(fn http.HandlerFunc, startOptions []trace.SpanStartOption, endOptions []trace.SpanEndOption) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := globalTracer.StartHTTPResponse(r, startOptions...)

		res := responseWriterSnooper{
			w:          w,
			statusCode: http.StatusOK,
		}

		fn(&res, r.WithContext(ctx))

		span.EndWithStatus(res.statusCode, endOptions...)
	}
}
