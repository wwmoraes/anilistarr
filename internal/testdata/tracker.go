package testdata

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var _ usecases.Tracker = (*MockTracker)(nil)

type MockTracker struct {
	mock.Mock
}

func (tracker *MockTracker) Close() error {
	args := tracker.Called()

	return args.Error(0)
}

func (tracker *MockTracker) GetUserID(ctx context.Context, name string) (string, error) {
	args := tracker.Called(ctx, name)

	return args.String(0), args.Error(1)
}

func (tracker *MockTracker) GetMediaListIDs(
	ctx context.Context,
	userID string,
) ([]entities.SourceID, error) {
	args := tracker.Called(ctx, userID)

	return args.Get(0).([]entities.SourceID), args.Error(1)
}
