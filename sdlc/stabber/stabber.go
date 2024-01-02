package stabber

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"dagger.io/dagger"
)

type Stabber struct {
	Client *dagger.Client
}

func TryUnwrapExecError(err error) error {
	execError := &dagger.ExecError{}

	if !errors.As(err, &execError) {
		return err
	}

	fmt.Fprintln(os.Stdout, strings.Trim(execError.Stdout, "\r\n"))
	fmt.Fprintln(os.Stderr, strings.Trim(execError.Stderr, "\r\n"))

	return fmt.Errorf(execError.Message())
}
