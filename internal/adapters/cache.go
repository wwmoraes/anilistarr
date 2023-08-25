package adapters

import (
	"context"
	"io"
	"time"
)

// Cache is a solution for fast storage and retrieval of ephemeral data
type Cache interface {
	io.Closer

	GetString(ctx context.Context, key string) (string, error)
	SetString(ctx context.Context, key, value string, options ...CacheOption) error
}

type CacheParams struct {
	TTL time.Duration
}

type CacheOption interface {
	Apply(params *CacheParams) error
}

type CacheOptionFn func(params *CacheParams) error

func (fn CacheOptionFn) Apply(params *CacheParams) error {
	return fn(params)
}

func WithTTL(duration time.Duration) CacheOption {
	return CacheOptionFn(func(params *CacheParams) error {
		params.TTL = duration

		return nil
	})
}

func NewCacheParams(options ...CacheOption) (*CacheParams, error) {
	params := &CacheParams{}

	for _, option := range options {
		err := option.Apply(params)
		if err != nil {
			return nil, err
		}
	}

	return params, nil
}
