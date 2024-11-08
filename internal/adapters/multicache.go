package adapters

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

type MultiCache []Cache

func (chain MultiCache) Close() error {
	errs := errgroup.Group{}

	for _, cache := range chain {
		errs.Go(cache.Close)
	}

	return errs.Wait()
}

func (chain MultiCache) GetString(ctx context.Context, key string) (string, error) {
	for _, cache := range chain {
		value, err := cache.GetString(ctx, key)
		if errors.Is(err, usecases.ErrNotFound) {
			continue
		}

		if err != nil || value != "" {
			return value, err
		}
	}

	return "", usecases.ErrNotFound
}

func (chain MultiCache) SetString(ctx context.Context, key, value string, options ...CacheOption) error {
	if len(chain) == 0 {
		return usecases.ErrNoCache
	}

	return chain[0].SetString(ctx, key, value)
}
