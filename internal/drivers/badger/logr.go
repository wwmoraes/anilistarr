package badger

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-logr/logr"
)

var _ badger.Logger = (*Logr)(nil)

// Logr is an adapter to log Badger DB internal problems using a [logr.Logger]
// instance.
type Logr struct {
	logr.Logger
}

// Errorf formats a new error and sends it to the configured log sink as an
// error event.
func (log *Logr) Errorf(format string, a ...any) {
	//nolint:err113 // upstream interface does not provide errors unfortunately
	log.Error(fmt.Errorf(format, a...), "Badger Error")
}

// Warningf formats a message and sends it to the configured log sink as an
// info event.
func (log *Logr) Warningf(format string, a ...any) {
	log.Info(fmt.Sprintf(format, a...))
}

// Infof formats a message and sends it to the configured log sink as an info
// event.
func (log *Logr) Infof(format string, a ...any) {
	log.Info(fmt.Sprintf(format, a...))
}

// Debugf formats a message and sends it to the configured log sink as an info
// event.
func (log *Logr) Debugf(format string, a ...any) {
	log.Info(fmt.Sprintf(format, a...))
}
