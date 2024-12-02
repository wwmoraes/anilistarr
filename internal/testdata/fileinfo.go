package testdata

import (
	"io/fs"
	"time"

	"github.com/stretchr/testify/mock"
)

var _ fs.FileInfo = (*MockFileInfo)(nil)

type MockFileInfo struct {
	mock.Mock
}

func (info *MockFileInfo) Name() string {
	args := info.Called()

	return args.String(0)
}

func (info *MockFileInfo) Size() int64 {
	args := info.Called()

	return args.Get(0).(int64)
}

func (info *MockFileInfo) Mode() fs.FileMode {
	args := info.Called()

	return args.Get(0).(fs.FileMode)
}

func (info *MockFileInfo) ModTime() time.Time {
	args := info.Called()

	return args.Get(0).(time.Time)
}

func (info *MockFileInfo) IsDir() bool {
	args := info.Called()

	return args.Bool(0)
}

func (info *MockFileInfo) Sys() any {
	args := info.Called()

	return args.Get(0)
}
