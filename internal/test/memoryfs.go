package test

import (
	"bytes"
	"fmt"
	"io/fs"
	"path/filepath"
	"time"
)

type MemoryFS map[string][]byte

func (mfs MemoryFS) Open(path string) (fs.File, error) {
	data, ok := mfs[path]
	if !ok {
		return nil, fmt.Errorf("file not found")
	}

	return &MemoryFile{
		name:    filepath.Base(path),
		content: bytes.NewBuffer(data),
	}, nil
}

type MemoryFile struct {
	closed  bool
	mode    fs.FileMode
	modTime time.Time

	name    string
	content *bytes.Buffer
}

func (fd *MemoryFile) Stat() (fs.FileInfo, error) {
	return fd, nil
}

func (fd *MemoryFile) Read(b []byte) (int, error) {
	if fd.closed {
		return 0, fs.ErrClosed
	}

	if fd.content == nil {
		return 0, fs.ErrInvalid
	}

	return fd.content.Read(b)
}

func (fd *MemoryFile) Close() error {
	if fd.closed {
		return fs.ErrClosed
	}

	fd.closed = true

	return nil
}

func (fd *MemoryFile) Name() string {
	return fd.name
}

func (fd *MemoryFile) Size() int64 {
	return int64(fd.content.Len())
}

func (fd *MemoryFile) Mode() fs.FileMode {
	return fd.mode
}

func (fd *MemoryFile) ModTime() time.Time {
	return fd.modTime
}

func (fd *MemoryFile) IsDir() bool {
	return false
}

func (fd *MemoryFile) Sys() any {
	return nil
}
