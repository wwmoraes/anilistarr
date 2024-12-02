package process_test

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/pkg/process"
)

func ExampleHandleExit_success() {
	defer process.HandleExit()

	process.Assert(nil)

	// Output:
	// os.Exit called
}

func ExampleHandleExit_error() {
	stderr := os.Stderr
	defer func() {
		//nolint:reassign // needed for example output validation
		os.Stderr = stderr
	}()

	//nolint:reassign // needed for example output validation
	os.Stderr = os.Stdout

	defer process.HandleExit()

	process.Assert(errors.New("foo"))

	// Output:
	// assertion failure: foo
	// os.Exit called
}

func ExampleAssertClose() {
	stderr := os.Stderr
	defer func() {
		//nolint:reassign // needed for example output validation
		os.Stderr = stderr
	}()

	//nolint:reassign // needed for example output validation
	os.Stderr = os.Stdout

	defer process.HandleExit()

	process.AssertClose(testdata.CloserFn(func() error {
		return errors.New("bar")
	}), "failed to close")

	// Output:
	// failed to close: bar
	// os.Exit called
}

func TestMain(m *testing.M) {
	patchOSExit := gomonkey.ApplyFunc(os.Exit, func(_ int) {
		fmt.Fprintln(os.Stdout, "os.Exit called")
	})
	defer patchOSExit.Reset()

	patchGoexit := gomonkey.ApplyFunc(runtime.Goexit, func() {})
	defer patchGoexit.Reset()

	m.Run()
}
