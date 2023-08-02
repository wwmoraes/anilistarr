package telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type Meter = metric.Meter
type MeterOption = metric.MeterOption

func newMeter(opts ...MeterOption) Meter {
	return otel.Meter(NAME, opts...)
}

func DefaultMeter() Meter {
	return globalMeter
}
