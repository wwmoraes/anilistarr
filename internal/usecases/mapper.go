package usecases

import (
	"context"
	"io"
)

type Mapper interface {
	io.Closer

	MapIDs(context.Context, []string) ([]string, error)
	MapID(context.Context, string) (string, error)
	Refresh(context.Context) error
}
