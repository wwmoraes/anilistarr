package adapters

import "context"

type MetadataSource[T Metadata] interface {
	Fetch(ctx context.Context) ([]T, error)
}
