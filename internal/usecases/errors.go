package usecases

import (
	"errors"
	"net/http"
)

// TODO refactor errors, those feel quite bland
var (
	// these "error codes" simulate RPC statuses
	// See https://grpc.io/docs/guides/status-codes/#the-full-list-of-status-codes

	// ErrStatusCancelled The operation was cancelled, typically by the caller.
	ErrStatusCancelled = errors.New("operation cancelled")

	// ErrStatusUnknown Unknown error. For example, this error may be returned when
	// a Status value received from another address space belongs to an error space
	// that is not known in this address space. Also errors raised by APIs that do
	// not return enough error information may be converted to this error.
	ErrStatusUnknown = errors.New("unknown error")

	// ErrStatusInvalidArgument The client specified an invalid argument. Note that
	// this differs from FAILED_PRECONDITION. INVALID_ARGUMENT indicates arguments
	// that are problematic regardless of the state of the system (e.g., a malformed
	// file name).
	ErrStatusInvalidArgument = errors.New("invalid argument")

	// ErrStatusDeadlineExceeded The deadline expired before the operation could
	// complete. For operations that change the state of the system, this error may
	// be returned even if the operation has completed successfully. For example, a
	// successful response from a server could have been delayed long enough for the
	// deadline to expire.
	ErrStatusDeadlineExceeded = errors.New("deadline exceeded")

	// ErrStatusNotFound Some requested entity (e.g., file or directory) was not
	// found. Note to server developers: if a request is denied for an entire
	// class of users, such as gradual feature rollout or undocumented allowlist,
	// NOT_FOUND may be used. If a request is denied for some users within a class
	// of users, such as user-based access control, PERMISSION_DENIED must be used.
	ErrStatusNotFound = errors.New("not found")

	// ErrStatusAlreadyExists The entity that a client attempted to create (e.g.,
	// file or directory) already exists.
	ErrStatusAlreadyExists = errors.New("already exists")

	// ErrStatusPermissionDenied The caller does not have permission to execute
	// the specified operation. PERMISSION_DENIED must not be used for rejections
	// caused by exhausting some resource (use RESOURCE_EXHAUSTED instead for those
	// errors). PERMISSION_DENIED must not be used if the caller can not be identified
	// (use UNAUTHENTICATED instead for those errors). This error code does not
	// imply the request is valid or the requested entity exists or satisfies other
	// pre-conditions.
	ErrStatusPermissionDenied = errors.New("permission denied")

	// ErrStatusResourceExhausted Some resource has been exhausted, perhaps a
	// per-user quota, or perhaps the entire file system is out of space.
	ErrStatusResourceExhausted = errors.New("resource exhausted")

	// ErrStatusFailedPrecondition The operation was rejected because the system
	// is not in a state required for the operation’s execution. For example, the
	// directory to be deleted is non-empty, an rmdir operation is applied to a
	// non-directory, etc. Service implementors can use the following guidelines
	// to decide between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE: (a) Use
	// UNAVAILABLE if the client can retry just the failing call. (b) Use ABORTED
	// if the client should retry at a higher level (e.g., when a client-specified
	// test-and-set fails, indicating the client should restart a read-modify-write
	// sequence). (c) Use FAILED_PRECONDITION if the client should not retry until the
	// system state has been explicitly fixed. E.g., if an “rmdir” fails because the
	// directory is non-empty, FAILED_PRECONDITION should be returned since the client
	// should not retry unless the files are deleted from the directory.
	ErrStatusFailedPrecondition = errors.New("failed precondition")

	// ErrStatusAborted The operation was aborted, typically due to a concurrency
	// issue such as a sequencer check failure or transaction abort. See the guidelines
	// above for deciding between FAILED_PRECONDITION, ABORTED, and UNAVAILABLE.
	ErrStatusAborted = errors.New("aborted")

	// ErrStatusOutOfRange The operation was attempted past the valid range.
	// E.g., seeking or reading past end-of-file. Unlike INVALID_ARGUMENT, this error
	// indicates a problem that may be fixed if the system state changes. For example,
	// a 32-bit file system will generate INVALID_ARGUMENT if asked to read at an
	// offset that is not in the range [0,2^32-1], but it will generate OUT_OF_RANGE
	// if asked to read from an offset past the current file size. There is a fair
	// bit of overlap between FAILED_PRECONDITION and OUT_OF_RANGE. We recommend using
	// OUT_OF_RANGE (the more specific error) when it applies so that callers who are
	// iterating through a space can easily look for an OUT_OF_RANGE error to detect
	// when they are done.
	ErrStatusOutOfRange = errors.New("out of range")

	// ErrStatusUnimplemented The operation is not implemented or is not supported/
	// enabled in this service.
	ErrStatusUnimplemented = errors.New("unimplemented")

	// ErrStatusInternal Internal errors. This means that some invariants expected
	// by the underlying system have been broken. This error code is reserved for
	// serious errors.
	ErrStatusInternal = errors.New("internal")

	// ErrStatusUnavailable The service is currently unavailable. This is most
	// likely a transient condition, which can be corrected by retrying with a backoff.
	// Note that it is not always safe to retry non-idempotent operations.
	ErrStatusUnavailable = errors.New("unavailable")

	// ErrStatusDataLoss Unrecoverable data loss or corruption.
	ErrStatusDataLoss = errors.New("data loss")

	// ErrStatusUnauthenticated The request does not have valid authentication
	// credentials for the operation.
	ErrStatusUnauthenticated = errors.New("unauthenticated")
)

// ErrorIn checks if an error matches any of the target errors. It uses
// [errors.Is] to check each target in order. Returns true if err matches any
// target, false otherwise.
func ErrorIn(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}

	return false
}

// ErrorJoinIf joins errors if the target error is not nil. Otherwise it returns
// nil.
func ErrorJoinIf(status, target error) error {
	if target == nil {
		return nil
	}

	return errors.Join(status, target)
}

// ErrorFromHTTPStatus maps HTTP status codes to use-case errors.
//
// Codes below 400 return nil as they're not errors in essence; callers must
// handle such cases directly, e.g. whether to follow a 3xx redirection or error
// is implementation-dependant.
//
//nolint:funlen,gocyclo,cyclop,maintidx // better than a global map
func ErrorFromHTTPStatus(code int) error {
	if code < http.StatusBadRequest {
		return nil
	}

	switch code {
	case http.StatusBadRequest: // 400 // RFC 9110, 15.5.1
		return ErrStatusInvalidArgument
	case http.StatusUnauthorized: // 401 // RFC 9110, 15.5.2
		return ErrStatusUnauthenticated
	case http.StatusPaymentRequired: // 402 // RFC 9110, 15.5.3
		return ErrStatusFailedPrecondition
	case http.StatusForbidden: // 403 // RFC 9110, 15.5.4
		return ErrStatusPermissionDenied
	case http.StatusNotFound: // 404 // RFC 9110, 15.5.5
		return ErrStatusNotFound
	case http.StatusMethodNotAllowed: // 405 // RFC 9110, 15.5.6
		return ErrStatusInvalidArgument
	case http.StatusNotAcceptable: // 406 // RFC 9110, 15.5.7
		return ErrStatusInvalidArgument
	case http.StatusProxyAuthRequired: // 407 // RFC 9110, 15.5.8
		return ErrStatusFailedPrecondition
	case http.StatusRequestTimeout: // 408 // RFC 9110, 15.5.9
		return ErrStatusDeadlineExceeded
	case http.StatusConflict: // 409 // RFC 9110, 15.5.10
		return ErrStatusAborted
	case http.StatusGone: // 410 // RFC 9110, 15.5.11
		return ErrStatusNotFound
	case http.StatusLengthRequired: // 411 // RFC 9110, 15.5.12
		return ErrStatusInvalidArgument
	case http.StatusPreconditionFailed: // 412 // RFC 9110, 15.5.13
		return ErrStatusFailedPrecondition
	case http.StatusRequestEntityTooLarge: // 413 // RFC 9110, 15.5.14
		return ErrStatusOutOfRange
	case http.StatusRequestURITooLong: // 414 // RFC 9110, 15.5.15
		return ErrStatusInvalidArgument
	case http.StatusUnsupportedMediaType: // 415 // RFC 9110, 15.5.16
		return ErrStatusInvalidArgument
	case http.StatusRequestedRangeNotSatisfiable: // 416 // RFC 9110, 15.5.17
		return ErrStatusOutOfRange
	case http.StatusExpectationFailed: // 417 // RFC 9110, 15.5.18
		return ErrStatusFailedPrecondition
	case http.StatusMisdirectedRequest: // 421 // RFC 9110, 15.5.20
		return ErrStatusAborted
	case http.StatusUnprocessableEntity: // 422 // RFC 9110, 15.5.21
		return ErrStatusFailedPrecondition
	case http.StatusLocked: // 423 // RFC 4918, 11.3
		return ErrStatusUnavailable
	case http.StatusFailedDependency: // 424 // RFC 4918, 11.4
		return ErrStatusUnavailable
	case http.StatusTooEarly: // 425 // RFC 8470, 5.2.
		return ErrStatusAborted
	case http.StatusUpgradeRequired: // 426 // RFC 9110, 15.5.22
		return ErrStatusFailedPrecondition
	case http.StatusPreconditionRequired: // 428 // RFC 6585, 3
		return ErrStatusAborted
	case http.StatusTooManyRequests: // 429 // RFC 6585, 4
		return ErrStatusUnavailable
	case http.StatusRequestHeaderFieldsTooLarge: // 431 // RFC 6585, 5
		return ErrStatusInvalidArgument
	case http.StatusUnavailableForLegalReasons: // 451 // RFC 7725, 3
		return ErrStatusUnavailable
	case http.StatusInternalServerError: // 500 // RFC 9110, 15.6.1
		return ErrStatusInternal
	case http.StatusNotImplemented: // 501 // RFC 9110, 15.6.2
		return ErrStatusUnimplemented
	case http.StatusBadGateway: // 502 // RFC 9110, 15.6.3
		return ErrStatusFailedPrecondition
	case http.StatusServiceUnavailable: // 503 // RFC 9110, 15.6.4
		return ErrStatusUnavailable
	case http.StatusGatewayTimeout: // 504 // RFC 9110, 15.6.5
		return ErrStatusDeadlineExceeded
	case http.StatusHTTPVersionNotSupported: // 505 // RFC 9110, 15.6.6
		return ErrStatusFailedPrecondition
	case http.StatusVariantAlsoNegotiates: // 506
		return ErrStatusAborted
	case http.StatusInsufficientStorage: // 507 // RFC 4918, 11.5
		return ErrStatusResourceExhausted
	case http.StatusLoopDetected: // 508 // RFC 5842, 7.2
		return ErrStatusAborted
	case http.StatusNotExtended: // 510 // RFC 2774, 7
		return ErrStatusInvalidArgument
	case http.StatusNetworkAuthenticationRequired: // 511 // RFC 6585, 6
		return ErrStatusUnauthenticated
	default:
		return ErrStatusUnknown
	}
}
