package telemetry

import (
	"log"
	"os"
	"time"

	"github.com/go-logr/logr"
)

type stdLogSink struct {
	stdout *log.Logger
	stderr *log.Logger
	values map[interface{}]interface{}
}

func NewStdLogSink() *stdLogSink {
	return &stdLogSink{
		stdout: log.New(os.Stdout, "", 0),
		stderr: log.New(os.Stderr, "", 0),
	}
}

func (sink *stdLogSink) Enabled(level int) bool {
	return true
}

func (sink *stdLogSink) Error(err error, msg string, keysAndValues ...interface{}) {
	values := mergeMaps(sink.values, kv2Map(keysAndValues...))

	sink.stderr.Printf("%s: %s [%s]", msg, err.Error(), mapString(values))
}

func (sink *stdLogSink) Info(level int, msg string, keysAndValues ...interface{}) {
	values := mergeMaps(sink.values, kv2Map(keysAndValues...))

	sink.stdout.Printf("%s [%s] %s", time.Now().Format(time.Stamp), mapString(values), msg)
}

func (sink *stdLogSink) Init(info logr.RuntimeInfo) {}

func (sink *stdLogSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	return &stdLogSink{
		stdout: log.New(sink.stdout.Writer(), sink.stdout.Prefix(), sink.stdout.Flags()),
		stderr: log.New(sink.stderr.Writer(), sink.stderr.Prefix(), sink.stderr.Flags()),
		values: mergeMaps(sink.values, kv2Map(keysAndValues...)),
	}
}

func (sink *stdLogSink) WithName(name string) logr.LogSink {
	return &stdLogSink{
		stdout: log.New(sink.stdout.Writer(), sink.stdout.Prefix()+name, sink.stderr.Flags()),
		stderr: log.New(sink.stderr.Writer(), sink.stderr.Prefix()+name, sink.stderr.Flags()),
		values: sink.values,
	}
}
