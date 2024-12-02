package testdata

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type Constructor[T any] func(testing.TB) T
type Functor[T any] func(testing.TB, T) T

func Compose[T any](tb testing.TB, value T, functors ...Functor[T]) T {
	tb.Helper()

	for _, decorate := range functors {
		value = decorate(tb, value)
	}

	return value
}

func Close(tb testing.TB, closer io.Closer) {
	tb.Helper()

	err := closer.Close()
	require.NoError(tb, err)
}

func Implements[T any](tb testing.TB) any {
	tb.Helper()

	return implements[T]()
}

func implements[T any]() any {
	return mock.MatchedBy(func(value any) bool {
		_, ok := value.(T)

		return ok
	})
}

func HTTPRequestWithJSONBody(tb testing.TB, body any) any {
	tb.Helper()

	want, err := json.Marshal(body)
	if err != nil {
		tb.Fatal(err)
	}

	return httpRequestWithBody(want)
}

func httpRequestWithBody(want []byte) any {
	return mock.MatchedBy(func(req *http.Request) bool {
		var body io.ReadCloser
		var data bytes.Buffer

		body, req.Body = req.Body, io.NopCloser(&data)

		_, err := io.Copy(&data, body)
		if err != nil {
			panic(err)
		}

		err = body.Close()
		if err != nil {
			panic(err)
		}

		return bytes.Equal(data.Bytes(), want)
	})
}

func HTTPResponseWithJSONBody(tb testing.TB, body any) *http.Response {
	w := httptest.NewRecorder()

	data, err := json.Marshal(body)
	if err != nil {
		tb.Fatal(err)
	}

	_, err = w.Write(data)
	if err != nil {
		tb.Fatal(err)
	}

	return w.Result()
}

func RESP(elements ...string) string {
	return RESPCommand(elements).String()
}

type RESPCommand []string

func NewRESPCommand(elements ...string) RESPCommand {
	return RESPCommand(elements)
}

func (cmd RESPCommand) String() string {
	return strings.Join(cmd, "\r\n") + "\r\n"
}

func CallerInfoLogger(tb testing.TB) {
	tb.Helper()

	calledPc, _, _, ok := runtime.Caller(4)
	if !ok {
		tb.Fatal("failed to retrieve caller program counter")
	}

	calledFuncInfo := runtime.FuncForPC(calledPc)
	if calledFuncInfo == nil {
		tb.Fatal("failed to retrieve caller function info")
	}

	pc, file, line, ok := runtime.Caller(5)
	if !ok {
		tb.Fatal("failed to retrieve caller program counter")
	}

	funcInfo := runtime.FuncForPC(pc)
	if funcInfo == nil {
		tb.Fatal("failed to retrieve caller function info")
	}

	tb.Logf("%s:%d:\ncaller: %s\ncalled: %s", file, line, funcInfo.Name(), calledFuncInfo.Name())
}

func CallerMatches(tb testing.TB, target string) bool {
	tb.Helper()

	pc, _, _, ok := runtime.Caller(5)
	if !ok {
		tb.Fatal("failed to retrieve caller program counter")
	}

	funcInfo := runtime.FuncForPC(pc)
	if funcInfo == nil {
		tb.Fatal("failed to retrieve caller function info")
	}

	return funcInfo.Name() == target
}
