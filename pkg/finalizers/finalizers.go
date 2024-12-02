// Package finalizers provides functions to use with [runtime.SetFinalizer].
package finalizers

import (
	"fmt"
	"io"
	"os"
)

// Close calls [io.Closer.Close], printing its error to [os.Stderr] if any.
func Close(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
