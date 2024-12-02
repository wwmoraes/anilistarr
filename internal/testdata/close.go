package testdata

import "io"

var _ io.Closer = CloserFn(nil)

type CloserFn func() error

func (fn CloserFn) Close() error {
	return fn()
}
