package usecases

import (
	"context"
	"io"
)

type Tracker interface {
	io.Closer

	GetUserID(ctx context.Context, name string) (string, error)
	GetMediaListIDs(ctx context.Context, userId string) ([]string, error)
}
