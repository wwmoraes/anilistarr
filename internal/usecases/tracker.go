package usecases

import (
	"context"
	"io"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

const (
	FailedGetUserErrorTemplate = "failed to get user ID: %w"
	ConvertUserIDErrorTemplate = "failed to convert user ID to integer: %w"
	FailedMediaErrorTemplate   = "failed to fetch media list: %w"
)

// Tracker provides access to user metadata such as ID and media list from an
// upstream media tracking service
type Tracker interface {
	io.Closer

	GetUserID(ctx context.Context, name string) (string, error)
	GetMediaListIDs(ctx context.Context, userId string) ([]entities.SourceID, error)
}
