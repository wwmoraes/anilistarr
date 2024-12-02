// Package process allows to assert errors and exit an application in a
// non-disruptive way. Contrary to plain [os.Exit], it calls deferred functions.
//
// It atomically stores the first non-zero exit code to terminate the
// application with it. Assert functions prints errors to [os.Stderr].
package process

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync/atomic"
)

var exitCode atomic.Int64

// HandleExit terminates the program by calling [os.Exit] with an atomically
// stored exit code. The idiomatic way is to call it as the first deferred
// function within a main function body.
//
// If never changed, the exit code is zero (success). Calls to [process.Exit],
// [process.Assert] and its variants may atomically change the exit code once to
// a non-zero value.
func HandleExit() {
	//nolint:revive // this runs within main so not a deep-exit
	os.Exit(int(exitCode.Load()))
}

// Exit atomically stores the exit code and exits the current goroutine. If
// this routine is main then the program will terminate after unwinding deferred
// calls until [process.HandleExit] executes.
func Exit(code int) {
	exitCode.CompareAndSwap(0, int64(code))
	runtime.Goexit()
}

// Assert prints error to the standard error and exits the current goroutine
// with code 1.
//
// It does nothing if error is nil.
func Assert(err error) {
	AssertWith(err, "assertion failure")
}

// AssertWith prints error to the standard error with prefix and exits the
// current goroutine with code 1.
//
// It does nothing if error is nil.
func AssertWith(err error, prefix string) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "%s: %s\n", prefix, err.Error())
	Exit(1)
}

// AssertClose calls [io.Closer.Close] from closer and asserts its error with
// [process.AssertWith].
func AssertClose(closer io.Closer, prefix string) {
	AssertWith(closer.Close(), prefix)
}
