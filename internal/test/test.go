// Package test provides utilities for testing. Do NOT use in production code!
package test

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// FileReadWith mocks a file read operation by copying the given data into the
// destination buffer of an [io.Reader.Read].
func FileReadWith(tb testing.TB, data []byte) func(mock.Arguments) {
	tb.Helper()

	return func(args mock.Arguments) {
		tb.Helper()

		dst, ok := args.Get(0).([]byte)
		if !ok {
			tb.Fatal("file read destination isn't a byte slice")

			return
		}

		copy(dst, data)
	}
}
