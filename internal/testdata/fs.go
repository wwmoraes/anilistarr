package testdata

import (
	"io/fs"

	"github.com/stretchr/testify/mock"
)

var _ fs.FS = (*MockFS)(nil)

type MockFS struct {
	mock.Mock
}

func (root *MockFS) Open(name string) (fs.File, error) {
	args := root.Called(name)

	return args.Get(0).(fs.File), args.Error(1)
}
