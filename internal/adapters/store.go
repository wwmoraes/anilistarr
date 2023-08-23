package adapters

import (
	"context"
	"io"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

type Store interface {
	io.Closer

	GetMedia(ctx context.Context, id string) (*entities.Media, error)
	GetMediaBulk(ctx context.Context, ids []string) ([]*entities.Media, error)
	PutMedia(ctx context.Context, media *entities.Media) error
	PutMediaBulk(ctx context.Context, medias []*entities.Media) error
}
