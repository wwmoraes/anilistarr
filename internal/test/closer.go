package test

import "io"

var _ io.Closer = CloserFn(nil)

// CloserFn represents a stateless [io.Closer].
type CloserFn func() error

// Close implements [io.Closer.Close].
func (fn CloserFn) Close() error {
	return fn()
}
