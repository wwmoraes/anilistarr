package usecases

import (
	"context"
	"io"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

// Tracker provides access to user metadata such as ID and media list from an
// upstream media tracking service
type Tracker interface {
	io.Closer

	GetUserID(ctx context.Context, name string) (string, error)
	GetMediaListIDs(ctx context.Context, userID string) ([]entities.SourceID, error)
}
