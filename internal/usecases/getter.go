package usecases

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
)

// Doer clients support execution of HTTP requests.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Getter supports fetching data from a given URI. It is client-dependant which
// URI components they support and for what.
type Getter interface {
	Get(ctx context.Context, uri string) ([]byte, error)
}

// GetterFn retrieves raw bytes from a given URI.
type GetterFn func(ctx context.Context, uri string) ([]byte, error)

// Get retrieves raw bytes from a given URI.
func (fn GetterFn) Get(ctx context.Context, uri string) ([]byte, error) {
	return fn(ctx, uri)
}

// HTTPGetter converts a [Doer] to a [Getter]. It considers only status code 200
// OK as successful; any other codes will result in an error.
func HTTPGetter(doer Doer) GetterFn {
	return GetterFn(func(ctx context.Context, uri string) ([]byte, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
		if err != nil {
			return nil, errors.Join(ErrStatusInvalidArgument, err)
		}

		res, err := doer.Do(req)
		if err != nil {
			return nil, errors.Join(ErrStatusUnknown, err)
		}
		defer res.Body.Close()

		err = ErrorFromHTTPStatus(res.StatusCode)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", err, "failed to execute HTTP request")
		}

		if res.StatusCode != http.StatusOK {
			return nil, ErrStatusUnknown
		}

		return io.ReadAll(res.Body)
	})
}

// FSGetter wraps a filesystem handler as a [Getter].
func FSGetter(root fs.FS) GetterFn {
	return GetterFn(func(_ context.Context, uri string) ([]byte, error) {
		return fs.ReadFile(root, uri)
	})
}
