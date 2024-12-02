package testdata

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var _ usecases.Store = (*MockStore)(nil)

type MockStore struct {
	mock.Mock
}

func (store *MockStore) Close() error {
	args := store.Called()

	return args.Error(0)
}

func (store *MockStore) GetMedia(ctx context.Context, id string) (*entities.Media, error) {
	args := store.Called(ctx, id)

	return args.Get(0).(*entities.Media), args.Error(1)
}

func (store *MockStore) GetMediaBulk(ctx context.Context, ids []string) ([]*entities.Media, error) {
	args := store.Called(ctx, ids)

	return args.Get(0).([]*entities.Media), args.Error(1)
}

func (store *MockStore) PutMedia(ctx context.Context, media *entities.Media) error {
	args := store.Called(ctx, media)

	return args.Error(0)
}

func (store *MockStore) PutMediaBulk(ctx context.Context, medias []*entities.Media) error {
	args := store.Called(ctx, medias)

	return args.Error(0)
}
