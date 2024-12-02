package testdata

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var _ usecases.MediaLister = (*MockMediaLister)(nil)

type MockMediaLister struct {
	mock.Mock
}

func (lister *MockMediaLister) Generate(ctx context.Context, name string) (entities.CustomList, error) {
	args := lister.Called(ctx, name)

	return args.Get(0).(entities.CustomList), args.Error(1)
}

func (lister *MockMediaLister) GetUserID(ctx context.Context, name string) (string, error) {
	args := lister.Called(ctx, name)

	return args.String(0), args.Error(1)
}

func (lister *MockMediaLister) Close() error {
	args := lister.Called()

	return args.Error(0)
}

func (lister *MockMediaLister) Refresh(ctx context.Context, client usecases.Getter) error {
	args := lister.Called(ctx, client)

	return args.Error(0)
}

func (lister *MockMediaLister) MapIDs(ctx context.Context, ids []entities.SourceID) ([]entities.TargetID, error) {
	args := lister.Called(ctx, ids)

	return args.Get(0).([]entities.TargetID), args.Error(1)
}
