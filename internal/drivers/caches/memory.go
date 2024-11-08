package caches

import (
	"context"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var _ adapters.Cache = &Memory{}

type Memory map[string]string

func NewMemory() Memory {
	return make(Memory)
}

func (mem Memory) Close() error {
	for k := range mem {
		delete(mem, k)
	}

	return nil
}

func (mem Memory) GetString(ctx context.Context, key string) (string, error) {
	value, ok := mem[key]
	if !ok {
		return "", usecases.ErrNotFound
	}

	return value, nil
}

func (mem Memory) SetString(ctx context.Context, key, value string, options ...adapters.CacheOption) error {
	mem[key] = value

	return nil
}
