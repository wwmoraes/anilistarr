package usecases

// Logger provides methods to emit informational and error messages. A simple
// logger will print former calls to stdout and latter ones to stderr.
//
// TODO move this interface to gotell
type Logger interface {
	Error(err error, msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	// WithValues(keysAndValues ...any) L
}
