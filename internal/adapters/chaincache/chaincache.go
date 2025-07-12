// Package chaincache provides a serial cache driver that uses other caches in
// order.
package chaincache

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var _ usecases.Cache = ChainCache(nil)

// ChainCache is a cache that chains multiple others. It resolves all requests
// serially i.e. it stops at the first underlying cache to return a nil error.
type ChainCache []usecases.Cache

// Close closes all underlying caches concurrently, returning either no error
// if they all close successfully or a joined error from those caches that fail.
func (chain ChainCache) Close() error {
	outErr := make(chan error)

	// close each cache in a goroutine, sending its potential error to the channel
	go func() {
		var wg sync.WaitGroup

		wg.Add(len(chain))

		for _, cache := range chain {
			go func() {
				outErr <- cache.Close()

				wg.Done()
			}()
		}

		wg.Wait()

		// close the channel to signal no more jobs
		close(outErr)
	}()

	// collects all errors
	errs := make([]error, 0, len(chain))
	for err := range outErr {
		errs = append(errs, err)
	}

	//nolint:wrapcheck // caches are internal
	return multierror.Append(nil, errs...).ErrorOrNil()
}

// GetString retrieves the value for key from the first cache that has it set.
// It attempts all caches in order. Returns [usecases.ErrStatusNotFound] if no
// cache has it. May short-circuit and return an [usecases.ErrStatusUnknown] if
// a cache produces an unexpected error.
func (chain ChainCache) GetString(ctx context.Context, key string) (string, error) {
	for _, cache := range chain {
		value, err := cache.GetString(ctx, key)
		if errors.Is(err, usecases.ErrStatusNotFound) {
			continue
		}

		if err != nil {
			return value, errors.Join(usecases.ErrStatusUnknown, err)
		}

		return value, nil
	}

	return "", usecases.ErrStatusNotFound
}

// SetString sets they key-value pair to the first cache in the chain.
func (chain ChainCache) SetString(
	ctx context.Context,
	key, value string,
	options ...usecases.CacheOption,
) error {
	if len(chain) == 0 {
		return fmt.Errorf("%w: %s", usecases.ErrStatusFailedPrecondition, "no caches in the chain")
	}

	//nolint:wrapcheck // passthrough
	return chain[0].SetString(ctx, key, value, options...)
}
