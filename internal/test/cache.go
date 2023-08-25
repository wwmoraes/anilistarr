package test

import (
	"context"

	"github.com/wwmoraes/anilistarr/internal/adapters"
)

type Cache map[string]string

func (cache *Cache) Close() error {
	*cache = Cache{}

	return nil
}

func (cache Cache) GetString(ctx context.Context, key string) (string, error) {
	return cache[key], nil
}

func (cache Cache) SetString(ctx context.Context, key, value string, options ...adapters.CacheOption) error {
	cache[key] = value

	return nil
}
