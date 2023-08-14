package telemetry

import (
	"os"

	"go.opentelemetry.io/otel/metric"
)

var WithDescription = metric.WithDescription
var WithUnit = metric.WithUnit

func Must[M any](metric M, err error) M {
	if err != nil {
		globalLogger.Error(err, "failed to create metric")
		os.Exit(1)
	}

	return metric
}

func Int64Counter(name string, options ...metric.Int64CounterOption) metric.Int64Counter {
	return Must(globalMeter.Int64Counter(name, options...))
}

func Float64Histogram(name string, options ...metric.Float64HistogramOption) metric.Float64Histogram {
	return Must(globalMeter.Float64Histogram(name, options...))
}
