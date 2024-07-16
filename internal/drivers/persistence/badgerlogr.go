package persistence

import (
	"fmt"

	"github.com/go-logr/logr"
)

// BadgerLogr wraps a logr.Logger and implements the badger.Logger interface
type BadgerLogr struct {
	logr.Logger
}

func (log *BadgerLogr) Errorf(format string, a ...any) {
	log.Error(fmt.Errorf(format, a...), "Badger Error")
}

func (log *BadgerLogr) Warningf(format string, a ...any) {
	log.Info(fmt.Sprintf(format, a...), "Badger Warning")
}

func (log *BadgerLogr) Infof(format string, a ...any) {
	log.Info(fmt.Sprintf(format, a...), "Badger Info")
}

func (log *BadgerLogr) Debugf(format string, a ...any) {
	log.Info(fmt.Sprintf(format, a...), "Badger Debug")
}
