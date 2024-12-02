package usecases_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestErrorIn(t *testing.T) {
	t.Parallel()

	type args struct {
		err     error
		targets []error
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "match plain",
			args: args{
				err:     usecases.ErrStatusUnimplemented,
				targets: []error{usecases.ErrStatusUnimplemented},
			},
			want: true,
		},
		{
			name: "match wrapped",
			args: args{
				err:     errors.Join(usecases.ErrStatusUnimplemented, errors.New("foo")),
				targets: []error{usecases.ErrStatusUnimplemented},
			},
			want: true,
		},
		{
			name: "mismatch plain",
			args: args{
				err:     usecases.ErrStatusUnimplemented,
				targets: []error{usecases.ErrStatusUnknown},
			},
			want: false,
		},
		{
			name: "mismatch wrapped",
			args: args{
				err:     errors.Join(usecases.ErrStatusUnimplemented, errors.New("foo")),
				targets: []error{usecases.ErrStatusUnknown},
			},
			want: false,
		},
		{
			name: "empty targets",
			args: args{
				err:     usecases.ErrStatusUnimplemented,
				targets: []error{},
			},
			want: false,
		},
		{
			name: "nil",
			args: args{
				err:     nil,
				targets: []error{nil},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, usecases.ErrorIn(tt.args.err, tt.args.targets...))
		})
	}
}

func TestErrorFromHTTPStatus(t *testing.T) {
	t.Parallel()

	tests := map[error][]int{
		nil: {
			http.StatusContinue,             // 100 // RFC 9110, 15.2.1
			http.StatusSwitchingProtocols,   // 101 // RFC 9110, 15.2.2
			http.StatusProcessing,           // 102 // RFC 2518, 10.1
			http.StatusEarlyHints,           // 103 // RFC 8297
			http.StatusOK,                   // 200 // RFC 9110, 15.3.1
			http.StatusCreated,              // 201 // RFC 9110, 15.3.2
			http.StatusAccepted,             // 202 // RFC 9110, 15.3.3
			http.StatusNonAuthoritativeInfo, // 203 // RFC 9110, 15.3.4
			http.StatusNoContent,            // 204 // RFC 9110, 15.3.5
			http.StatusResetContent,         // 205 // RFC 9110, 15.3.6
			http.StatusPartialContent,       // 206 // RFC 9110, 15.3.7
			http.StatusMultiStatus,          // 207 // RFC 4918, 11.1
			http.StatusAlreadyReported,      // 208 // RFC 5842, 7.1
			http.StatusIMUsed,               // 226 // RFC 3229, 10.4.1
			http.StatusMultipleChoices,      // 300 // RFC 9110, 15.4.1
			http.StatusMovedPermanently,     // 301 // RFC 9110, 15.4.2
			http.StatusFound,                // 302 // RFC 9110, 15.4.3
			http.StatusSeeOther,             // 303 // RFC 9110, 15.4.4
			http.StatusNotModified,          // 304 // RFC 9110, 15.4.5
			http.StatusUseProxy,             // 305 // RFC 9110, 15.4.6
			http.StatusTemporaryRedirect,    // 307 // RFC 9110, 15.4.8
			http.StatusPermanentRedirect,    // 308 // RFC 9110, 15.4.9
		},
		usecases.ErrStatusNotFound: {
			http.StatusNotFound, // 404
			http.StatusGone,     // 410
		},
		usecases.ErrStatusUnauthenticated: {
			http.StatusUnauthorized,                  // 401
			http.StatusNetworkAuthenticationRequired, // 511
		},
		usecases.ErrStatusPermissionDenied: {
			http.StatusForbidden, // 403
		},
		usecases.ErrStatusResourceExhausted: {
			http.StatusInsufficientStorage, // 507
		},
		usecases.ErrStatusUnavailable: {
			http.StatusLocked,                     // 423
			http.StatusFailedDependency,           // 424
			http.StatusTooManyRequests,            // 429
			http.StatusUnavailableForLegalReasons, // 451
			http.StatusServiceUnavailable,         // 503
		},
		usecases.ErrStatusInvalidArgument: {
			http.StatusBadRequest,                  // 400
			http.StatusMethodNotAllowed,            // 405
			http.StatusNotAcceptable,               // 406
			http.StatusLengthRequired,              // 411
			http.StatusRequestURITooLong,           // 414
			http.StatusUnsupportedMediaType,        // 415
			http.StatusRequestHeaderFieldsTooLarge, // 431
			http.StatusNotExtended,                 // 510
		},
		usecases.ErrStatusUnimplemented: {
			http.StatusNotImplemented, // 501
		},
		usecases.ErrStatusAborted: {
			http.StatusConflict,              // 409
			http.StatusMisdirectedRequest,    // 421
			http.StatusTooEarly,              // 425
			http.StatusPreconditionRequired,  // 428
			http.StatusVariantAlsoNegotiates, // 506
			http.StatusLoopDetected,          // 508
		},
		usecases.ErrStatusFailedPrecondition: {
			http.StatusPaymentRequired,         // 402
			http.StatusProxyAuthRequired,       // 407
			http.StatusPreconditionFailed,      // 412
			http.StatusExpectationFailed,       // 417
			http.StatusUnprocessableEntity,     // 422
			http.StatusUpgradeRequired,         // 426
			http.StatusBadGateway,              // 502
			http.StatusHTTPVersionNotSupported, // 505
		},
		usecases.ErrStatusInternal: {
			http.StatusInternalServerError, // 500
		},
		usecases.ErrStatusDeadlineExceeded: {
			http.StatusRequestTimeout, // 408
			http.StatusGatewayTimeout, // 504
		},
		usecases.ErrStatusOutOfRange: {
			http.StatusRequestEntityTooLarge,        // 413
			http.StatusRequestedRangeNotSatisfiable, // 416
		},
		usecases.ErrStatusUnknown: {
			http.StatusTeapot, // 418
			999,
		},
	}

	for err, codes := range tests {
		for _, code := range codes {
			t.Run(fmt.Sprintf("%d %s", code, http.StatusText(code)), func(t *testing.T) {
				t.Parallel()

				require.ErrorIs(t, usecases.ErrorFromHTTPStatus(code), err)
			})
		}
	}
}

func TestErrorJoinIf(t *testing.T) {
	t.Parallel()

	errFoo := errors.New("foo")

	type args struct {
		status error
		target error
	}

	tests := []struct {
		args       args
		name       string
		wantErrors []error
	}{
		{
			name: "nil target",
			args: args{
				status: usecases.ErrStatusUnknown,
				target: nil,
			},
			wantErrors: []error{nil},
		},
		{
			name: "nil status",
			args: args{
				status: nil,
				target: errFoo,
			},
			wantErrors: []error{errFoo},
		},
		{
			name: "both non-nil",
			args: args{
				status: usecases.ErrStatusInternal,
				target: errFoo,
			},
			wantErrors: []error{usecases.ErrStatusInternal, errFoo},
		},
		{
			name: "both nil",
			args: args{
				status: nil,
				target: nil,
			},
			wantErrors: []error{nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := usecases.ErrorJoinIf(tt.args.status, tt.args.target)
			for _, want := range tt.wantErrors {
				require.ErrorIs(t, err, want)
			}
		})
	}
}
