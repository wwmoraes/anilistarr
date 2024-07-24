package process

import (
	"fmt"
	"os"
	"runtime"
	"sync/atomic"
)

var exitCode atomic.Int32

func HandleExit() {
	os.Exit(int(exitCode.Load()))
}

func Exit(code int) {
	exitCode.CompareAndSwap(0, int32(code))
	runtime.Goexit()
}

func Assert(err error) {
	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr, err)
	Exit(1)
}
