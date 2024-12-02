package testdata

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var _ usecases.Getter = (*MockGetter)(nil)

type MockGetter struct {
	mock.Mock
}

func (getter *MockGetter) Get(ctx context.Context, uri string) ([]byte, error) {
	args := getter.Called(ctx, uri)

	return args.Get(0).([]byte), args.Error(1)
}
