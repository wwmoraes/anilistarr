package telemetry

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
)

const (
	keyValueSeparator = ": "
	entrySeparator    = " | "
)

type Logger = logr.Logger

func DefaultLogger() logr.Logger {
	return globalLogger
}

func ContextWithLogger(ctx context.Context) context.Context {
	return logr.NewContext(ctx, globalLogger)
}

func LoggerFromContext(ctx context.Context) logr.Logger {
	return logr.FromContextOrDiscard(ctx)
}

func kv2Map(keysAndValues ...interface{}) map[interface{}]interface{} {
	values := make(map[interface{}]interface{}, len(keysAndValues)/2)

	for i := 0; i+1 < len(keysAndValues); i = i + 2 {
		values[keysAndValues[i]] = keysAndValues[i+1]
	}

	return values
}

func mergeMaps(maps ...map[interface{}]interface{}) map[interface{}]interface{} {
	values := make(map[interface{}]interface{})

	for _, entry := range maps {
		for k, v := range entry {
			values[k] = v
		}
	}

	return values
}

func mapString(m map[interface{}]interface{}) string {
	entries := make([]string, 0, len(m))

	for k, v := range m {
		if k == nil || v == nil {
			continue
		}

		entries = append(entries, fmt.Sprintf("%v%s%v", k, keyValueSeparator, v))
	}

	return strings.Join(entries, entrySeparator)
}
