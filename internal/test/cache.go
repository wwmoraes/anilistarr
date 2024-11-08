package test

import (
	"context"
	"errors"
	"sync"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

type Cache struct {
	mu sync.RWMutex

	Data map[string]string
}

func (cache *Cache) Close() error {
	return nil
}

func (cache *Cache) GetString(ctx context.Context, key string) (string, error) {
	if cache == nil || cache.Data == nil {
		return "", errors.New("cache data is nil")
	}

	cache.mu.RLock()
	defer cache.mu.RUnlock()

	value, ok := cache.Data[key]
	if !ok {
		return "", usecases.ErrNotFound
	}

	return value, nil
}

func (cache *Cache) SetString(ctx context.Context, key, value string, options ...adapters.CacheOption) error {
	if cache == nil || cache.Data == nil {
		return errors.New("cache data is nil")
	}

	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.Data[key] = value

	return nil
}
