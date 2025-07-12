package testdata

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var _ usecases.Cache = (*MockCache)(nil)

type MockCache struct {
	mock.Mock
}

func (cache *MockCache) Close() error {
	args := cache.Called()

	return args.Error(0)
}

func (cache *MockCache) GetString(ctx context.Context, key string) (string, error) {
	args := cache.Called(ctx, key)

	return args.String(0), args.Error(1)
}

func (cache *MockCache) SetString(
	ctx context.Context,
	key, value string,
	options ...usecases.CacheOption,
) error {
	args := cache.Called(ctx, key, value, options)

	return args.Error(0)
}
