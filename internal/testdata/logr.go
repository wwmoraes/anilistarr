package testdata

import (
	"github.com/go-logr/logr"
	"github.com/stretchr/testify/mock"
)

var _ logr.LogSink = (*MockLogrLogSink)(nil)

type MockLogrLogSink struct {
	mock.Mock
}

func (sink *MockLogrLogSink) Init(info logr.RuntimeInfo) {
	sink.Called(info)
}

func (sink *MockLogrLogSink) Enabled(level int) bool {
	args := sink.Called(level)

	return args.Bool(0)
}

func (sink *MockLogrLogSink) Info(level int, msg string, keysAndValues ...any) {
	sink.Called(level, msg, keysAndValues)
}

func (sink *MockLogrLogSink) Error(err error, msg string, keysAndValues ...any) {
	sink.Called(err, msg, keysAndValues)
}

func (sink *MockLogrLogSink) WithValues(keysAndValues ...any) logr.LogSink {
	args := sink.Called(keysAndValues...)

	return args.Get(0).(logr.LogSink)
}

func (sink *MockLogrLogSink) WithName(name string) logr.LogSink {
	args := sink.Called(name)

	return args.Get(0).(logr.LogSink)
}
