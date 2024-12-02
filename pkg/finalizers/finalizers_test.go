package finalizers_test

import (
	"errors"
	"os"

	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/pkg/finalizers"
)

func ExampleClose() {
	//nolint:reassign // needed as examples ignore stderr
	os.Stderr = os.Stdout

	finalizers.Close(testdata.CloserFn(func() error {
		return nil
	}))

	finalizers.Close(testdata.CloserFn(func() error {
		return errors.New("foo")
	}))

	// Output:
	// foo
}
