package test

import (
	"context"
)

type Cache map[string]string

func (cache *Cache) Close() error {
	*cache = Cache{}

	return nil
}

func (cache Cache) GetString(ctx context.Context, key string) (string, error) {
	return cache[key], nil
}

func (cache Cache) SetString(ctx context.Context, key, value string) error {
	cache[key] = value

	return nil
}
