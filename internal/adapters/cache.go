package adapters

import (
	"context"
	"io"
)

type Cache interface {
	io.Closer

	GetString(ctx context.Context, key string) (string, error)
	SetString(ctx context.Context, key, value string) error
}
