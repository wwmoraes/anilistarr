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

type Tracker interface {
	io.Closer

	GetUserID(ctx context.Context, name string) (string, error)
	GetMediaListIDs(ctx context.Context, userId string) ([]entities.SourceID, error)
}
