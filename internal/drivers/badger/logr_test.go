package badger_test

import (
	"errors"
	"testing"

	"github.com/go-logr/logr"

	"github.com/wwmoraes/anilistarr/internal/drivers/badger"
	"github.com/wwmoraes/anilistarr/internal/testdata"
)

func TestLogr(t *testing.T) {
	t.Parallel()

	var keysAndValues []any

	err := errors.New("foo: bar")
	sink := testdata.MockLogrLogSink{}

	sink.On("Init", logr.RuntimeInfo{CallDepth: 1}).Once()
	sink.On("Enabled", 0).Return(true).Times(3)
	sink.On("Error", err, "Badger Error", keysAndValues).Once()
	sink.On("Info", 0, "baz: qux", keysAndValues).Once()
	sink.On("Info", 0, "quux: corge", keysAndValues).Once()
	sink.On("Info", 0, "grault: garply", keysAndValues).Once()

	log := badger.Logr{
		Logger: logr.New(&sink),
	}

	log.Errorf("foo: %s", "bar")
	log.Warningf("baz: %s", "qux")
	log.Infof("quux: %s", "corge")
	log.Debugf("grault: %s", "garply")

	sink.AssertExpectations(t)
}
