package telemetry

import "github.com/go-logr/logr"

type teeSink []logr.LogSink

func TeeSink(sinks ...logr.LogSink) logr.LogSink {
	return teeSink(sinks)
}

func (sinks teeSink) Enabled(level int) bool {
	for _, sink := range sinks {
		if !sink.Enabled(level) {
			return false
		}
	}

	return true
}

func (sinks teeSink) Error(err error, msg string, keysAndValues ...interface{}) {
	for _, sink := range sinks {
		sink.Error(err, msg, keysAndValues...)
	}
}

func (sinks teeSink) Info(level int, msg string, keysAndValues ...interface{}) {
	for _, sink := range sinks {
		sink.Info(level, msg, keysAndValues...)
	}
}

func (sinks teeSink) Init(info logr.RuntimeInfo) {
	for _, sink := range sinks {
		sink.Init(info)
	}
}

func (sinks teeSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	newSinks := make(teeSink, len(sinks))

	for index, sink := range sinks {
		newSinks[index] = sink.WithValues(keysAndValues...)
	}

	return newSinks
}

func (sinks teeSink) WithName(name string) logr.LogSink {
	newSinks := make(teeSink, len(sinks))

	for index, sink := range sinks {
		newSinks[index] = sink.WithName(name)
	}

	return newSinks
}
