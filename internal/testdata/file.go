package testdata

import (
	"io/fs"

	"github.com/stretchr/testify/mock"
)

var _ fs.File = (*MockFile)(nil)

type MockFile struct {
	mock.Mock
}

func (file *MockFile) Stat() (fs.FileInfo, error) {
	args := file.Called()

	return args.Get(0).(fs.FileInfo), args.Error(1)
}

func (file *MockFile) Read(dst []byte) (int, error) {
	args := file.Called(dst)

	return args.Int(0), args.Error(1)
}

func (file *MockFile) Close() error {
	args := file.Called()

	return args.Error(0)
}

func FileReadWith(data []byte) func(mock.Arguments) {
	return func(args mock.Arguments) {
		dst := args.Get(0).([]byte)
		copy(dst, data)
	}
}
