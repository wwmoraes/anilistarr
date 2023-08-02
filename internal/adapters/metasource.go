package adapters

import (
	"context"
	"net/http"
)

type Getter interface {
	Get(string) (*http.Response, error)
}

type MetadataSource[T Metadata] interface {
	Fetch(ctx context.Context, client Getter) ([]T, error)
}
