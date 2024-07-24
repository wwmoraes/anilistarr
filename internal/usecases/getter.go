package usecases

import (
	"fmt"
	"io"
	"net/http"
)

// HTTPGetter provides a way to retrieve data through HTTP GET requests
type HTTPGetter interface {
	Get(uri string) (*http.Response, error)
}

// Getter provides a way to retrieve data from a URI
type Getter interface {
	Get(uri string) ([]byte, error)
}

type GetterFn func(uri string) ([]byte, error)

func (fn GetterFn) Get(uri string) ([]byte, error) {
	return fn(uri)
}

func HTTPGetterAsGetter(getter HTTPGetter) Getter {
	return GetterFn(func(uri string) ([]byte, error) {
		res, err := getter.Get(uri)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch remote JSON: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("provider data not found")
		}

		return io.ReadAll(res.Body)
	})
}
