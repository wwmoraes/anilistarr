package testdata

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

var _ context.Context = (*MockContext)(nil)

type MockContext struct {
	mock.Mock
}

func (ctx *MockContext) Deadline() (deadline time.Time, ok bool) {
	args := ctx.Called()

	return args.Get(0).(time.Time), args.Bool(1)
}

func (ctx *MockContext) Done() <-chan struct{} {
	args := ctx.Called()

	return args.Get(0).(<-chan struct{})
}

func (ctx *MockContext) Err() error {
	args := ctx.Called()

	return args.Error(0)
}

func (ctx *MockContext) Value(key any) any {
	args := ctx.Called()

	return args.Get(0)
}
