package telemetry

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/MrAlias/otlpr"
	"github.com/go-logr/logr"
	otelruntime "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	NoEndpointError = errors.New("OTLP endpoint not provided")

	otlpConnHandler sync.Once
	otlpConn        *grpc.ClientConn
	otlpConnErr     error

	otlpResource *resource.Resource

	globalTracer Tracer
	globalMeter  Meter
	globalLogger Logger
)

func init() {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	var err error
	otlpResource, err = resource.Merge(resource.Empty(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.CodeNamespace(MODULE),
		semconv.ServiceName(NAME),
		semconv.ServiceNamespace(NAMESPACE),
		semconv.ServiceVersion(VERSION),
	))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create OTLP resource: %s", err.Error())
	}

	globalTracer = newTracer()
	globalMeter = newMeter()
	globalLogger = logr.New(NewStdLogSink())

	initMetrics(globalMeter)
}

func getOTLPConnGRPC(ctx context.Context, otlpEndpoint string) (*grpc.ClientConn, error) {
	otlpConnHandler.Do(func() {
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		otlpConn, otlpConnErr = grpc.DialContext(ctx, otlpEndpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
		)
		if otlpConnErr != nil {
			otlpConnErr = fmt.Errorf("failed to connect to the OTLP endpoint: %w", otlpConnErr)
		}
	})

	return otlpConn, otlpConnErr
}

func providerShutdown(shutdown func(context.Context) error) func(context.Context) {
	return func(ctx context.Context) {
		if err := shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}
}

func InstrumentTracing(ctx context.Context, otlpEndpoint string) (func(context.Context), error) {
	conn, err := getOTLPConnGRPC(ctx, otlpEndpoint)
	if err != nil {
		return nil, err
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create an OTLP exporter: %w", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(otlpResource),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(traceProvider)

	return providerShutdown(traceProvider.Shutdown), nil
}

func InstrumentMetrics(ctx context.Context, otlpEndpoint string) (func(context.Context), error) {
	conn, err := getOTLPConnGRPC(ctx, otlpEndpoint)
	if err != nil {
		return nil, err
	}

	meterExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create an OTLP exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(otlpResource),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(meterExporter)),
	)

	err = otelruntime.Start(
		otelruntime.WithMeterProvider(meterProvider),
		otelruntime.WithMinimumReadMemStatsInterval(time.Second),
	)
	if err != nil {
		providerShutdown(meterProvider.Shutdown)(ctx)
		return nil, err
	}

	otel.SetMeterProvider(meterProvider)

	return providerShutdown(meterProvider.Shutdown), nil
}

func InstrumentLogging(ctx context.Context, otlpEndpoint string) error {
	conn, err := getOTLPConnGRPC(ctx, otlpEndpoint)
	if err != nil {
		return err
	}

	logger := otlpr.NewWithOptions(conn, otlpr.Options{
		LogCaller:     otlpr.All,
		LogCallerFunc: true,
		Batcher: otlpr.Batcher{
			Messages: 10,
			Timeout:  10 * time.Second,
		},
	})
	logger = otlpr.WithResource(logger, otlpResource)
	otlpSink := logger.GetSink()

	globalLogger = logger.WithSink(TeeSink(globalLogger.GetSink(), otlpSink))

	otel.SetLogger(globalLogger)

	return nil
}

// func InstrumentProfiling(ctx context.Context) (func() error, error) {
// 	runtime.SetMutexProfileFraction(5)
// 	runtime.SetBlockProfileRate(5)

// 	// TODO or https://pkg.go.dev/net/http/pprof
// 	profiler, err := pyroscope.Start(pyroscope.Config{
// 		ApplicationName: NAME,
// 		ServerAddress:   "",
// 		AuthToken:       "",
// 		Logger:          pyroscope.StandardLogger,
// 		ProfileTypes: []pyroscope.ProfileType{
// 			pyroscope.ProfileCPU,
// 			pyroscope.ProfileAllocObjects,
// 		},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return profiler.Stop, nil
// }

func InstrumentAll(ctx context.Context, otlpEndpoint string) (func(context.Context), error) {
	if len(otlpEndpoint) <= 0 {
		return nil, NoEndpointError
	}

	tracingShutdown, err := InstrumentTracing(ctx, otlpEndpoint)
	if err != nil {
		return nil, err
	}

	metricsShutdown, err := InstrumentMetrics(ctx, otlpEndpoint)
	if err != nil {
		return nil, err
	}

	err = InstrumentLogging(ctx, otlpEndpoint)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context) {
		tracingShutdown(ctx)
		metricsShutdown(ctx)
	}, nil
}
