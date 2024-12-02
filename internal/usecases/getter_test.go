package usecases_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestHTTPGetter(t *testing.T) {
	t.Parallel()

	uri := "http://example.com/foo.txt"
	want := []byte("bar")

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		uri,
		http.NoBody,
	)
	require.NoError(t, err)

	resWriter := httptest.NewRecorder()
	_, err = resWriter.Write(want)
	require.NoError(t, err)

	res := resWriter.Result()
	defer res.Body.Close()

	doer := testdata.MockDoer{}

	doer.On("Do", req).
		Return(res, nil).Once()

	getter := usecases.HTTPGetter(&doer)

	got, err := getter.Get(context.TODO(), uri)
	require.NoError(t, err)

	assert.Equal(t, want, got)
	doer.AssertExpectations(t)
}

func TestHTTPGetter_NewRequestWithContext_error(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		uri string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "nil context",
			args: args{
				ctx: nil,
				uri: "",
			},
		},
		{
			name: "invalid URI",
			args: args{
				ctx: context.TODO(),
				uri: "\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			doer := testdata.MockDoer{}

			getter := usecases.HTTPGetter(&doer)

			got, err := getter.Get(tt.args.ctx, tt.args.uri)
			require.ErrorIs(t, err, usecases.ErrStatusInvalidArgument)

			assert.Nil(t, got)
			doer.AssertExpectations(t)
		})
	}
}

func TestHTTPGetter_Do_error(t *testing.T) {
	t.Parallel()

	var res *http.Response

	uri := "http://example.com/foo.txt"
	doErr := errors.New("bar")

	req, err := http.NewRequestWithContext(
		context.TODO(),
		http.MethodGet,
		uri,
		http.NoBody,
	)
	require.NoError(t, err)

	doer := testdata.MockDoer{}

	doer.On("Do", req).
		Return(res, doErr).Once()

	getter := usecases.HTTPGetter(&doer)

	got, err := getter.Get(context.TODO(), uri)
	require.ErrorIs(t, err, usecases.ErrStatusUnknown)

	assert.Nil(t, got)
	doer.AssertExpectations(t)
}

func TestHTTPGetter_non200(t *testing.T) {
	t.Parallel()

	codes := []int{
		http.StatusContinue,                      // 100
		http.StatusSwitchingProtocols,            // 101
		http.StatusProcessing,                    // 102
		http.StatusEarlyHints,                    // 103
		http.StatusCreated,                       // 201
		http.StatusAccepted,                      // 202
		http.StatusNonAuthoritativeInfo,          // 203
		http.StatusNoContent,                     // 204
		http.StatusResetContent,                  // 205
		http.StatusPartialContent,                // 206
		http.StatusMultiStatus,                   // 207
		http.StatusAlreadyReported,               // 208
		http.StatusIMUsed,                        // 226
		http.StatusMultipleChoices,               // 300
		http.StatusMovedPermanently,              // 301
		http.StatusFound,                         // 302
		http.StatusSeeOther,                      // 303
		http.StatusNotModified,                   // 304
		http.StatusUseProxy,                      // 305
		http.StatusTemporaryRedirect,             // 307
		http.StatusPermanentRedirect,             // 308
		http.StatusBadRequest,                    // 400
		http.StatusUnauthorized,                  // 401
		http.StatusPaymentRequired,               // 402
		http.StatusForbidden,                     // 403
		http.StatusNotFound,                      // 404
		http.StatusMethodNotAllowed,              // 405
		http.StatusNotAcceptable,                 // 406
		http.StatusProxyAuthRequired,             // 407
		http.StatusRequestTimeout,                // 408
		http.StatusConflict,                      // 409
		http.StatusGone,                          // 410
		http.StatusLengthRequired,                // 411
		http.StatusPreconditionFailed,            // 412
		http.StatusRequestEntityTooLarge,         // 413
		http.StatusRequestURITooLong,             // 414
		http.StatusUnsupportedMediaType,          // 415
		http.StatusRequestedRangeNotSatisfiable,  // 416
		http.StatusExpectationFailed,             // 417
		http.StatusTeapot,                        // 418
		http.StatusMisdirectedRequest,            // 421
		http.StatusUnprocessableEntity,           // 422
		http.StatusLocked,                        // 423
		http.StatusFailedDependency,              // 424
		http.StatusTooEarly,                      // 425
		http.StatusUpgradeRequired,               // 426
		http.StatusPreconditionRequired,          // 428
		http.StatusTooManyRequests,               // 429
		http.StatusRequestHeaderFieldsTooLarge,   // 431
		http.StatusUnavailableForLegalReasons,    // 451
		http.StatusInternalServerError,           // 500
		http.StatusNotImplemented,                // 501
		http.StatusBadGateway,                    // 502
		http.StatusServiceUnavailable,            // 503
		http.StatusGatewayTimeout,                // 504
		http.StatusHTTPVersionNotSupported,       // 505
		http.StatusVariantAlsoNegotiates,         // 506
		http.StatusInsufficientStorage,           // 507
		http.StatusLoopDetected,                  // 508
		http.StatusNotExtended,                   // 510
		http.StatusNetworkAuthenticationRequired, // 511
		999,
	}

	for _, code := range codes {
		t.Run(fmt.Sprintf("%d %s", code, http.StatusText(code)), func(t *testing.T) {
			t.Parallel()

			uri := "http://example.com/foo.txt"

			req, err := http.NewRequestWithContext(
				context.TODO(),
				http.MethodGet,
				uri,
				http.NoBody,
			)
			require.NoError(t, err)

			resWriter := httptest.NewRecorder()
			resWriter.WriteHeader(code)

			res := resWriter.Result()
			defer res.Body.Close()

			doer := testdata.MockDoer{}

			doer.On("Do", req).
				Return(res, nil).Once()

			getter := usecases.HTTPGetter(&doer)

			got, err := getter.Get(context.TODO(), uri)
			require.Error(t, err)

			assert.Nil(t, got)
			doer.AssertExpectations(t)
		})
	}
}

func TestFSGetter(t *testing.T) {
	t.Parallel()

	name := "foo.txt"
	data := []byte("bar")

	root := testdata.MockFS{}
	file := testdata.MockFile{}
	fileInfo := testdata.MockFileInfo{}

	root.On("Open", name).
		Return(&file, nil).Once()
	file.On("Stat").
		Return(&fileInfo, nil).Once()
	file.On("Close").
		Return(nil).Once()
	fileInfo.On("Size").
		Return(int64(len(data))).Once()

	// first read returns data up to the file size
	fileReadCall := file.On("Read", mock.AnythingOfType("[]uint8")).
		Return(len(data), nil).Once()
	fileReadCall.RunFn = testdata.FileReadWith(data)

	// second read returns EOF and breaks the loop
	file.On("Read", mock.AnythingOfType("[]uint8")).
		Return(0, io.EOF).Once().NotBefore(fileReadCall)

	getter := usecases.FSGetter(&root)

	got, err := getter.Get(context.TODO(), name)
	require.NoError(t, err)

	assert.Equal(t, data, got)
	root.AssertExpectations(t)
	file.AssertExpectations(t)
	fileInfo.AssertExpectations(t)
}
