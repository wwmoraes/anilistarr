package telemetry

import "go.opentelemetry.io/otel/attribute"

type Attributable interface {
	SetAttributes(kv ...attribute.KeyValue)
}

func Int(element Attributable, k string, v int) {
	element.SetAttributes(attribute.Int(k, v))
}
