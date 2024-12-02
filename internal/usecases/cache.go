package usecases

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

// CacheOptions contains optional parameters to use when setting cache entries.
type CacheOptions struct {
	TTL time.Duration
}

// CacheOption provides a method to apply changes to an existing set of cache
// options.
type CacheOption interface {
	Apply(params *CacheOptions)
}

// CacheOptionFn applies changes to existing cache parameters.
type CacheOptionFn func(params *CacheOptions)

// Apply changes cache parameters by applying itself.
func (fn CacheOptionFn) Apply(params *CacheOptions) {
	fn(params)
}

// WithTTL defines a time-to-live (TTL) for a cache entry.
func WithTTL(duration time.Duration) CacheOption {
	return CacheOptionFn(func(params *CacheOptions) {
		params.TTL = duration
	})
}

// NewCacheOptions creates a default set of parameters and applies all provided
// option in order.
func NewCacheOptions(options ...CacheOption) *CacheOptions {
	params := &CacheOptions{
		TTL: 0,
	}

	for _, option := range options {
		option.Apply(params)
	}

	return params
}
