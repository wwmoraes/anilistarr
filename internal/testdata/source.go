package testdata

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var _ usecases.Source = (*MockSource)(nil)

type MockSource struct {
	mock.Mock
}

func (source *MockSource) String() string {
	args := source.Called()

	return args.String(0)
}

func (source *MockSource) Fetch(
	ctx context.Context,
	client usecases.Getter,
) ([]usecases.Metadata, error) {
	args := source.Called(ctx, client)

	return args.Get(0).([]usecases.Metadata), args.Error(1)
}
