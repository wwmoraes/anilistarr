package functional

func Unwrap[T any](value T, err error) T {
	return value
}
