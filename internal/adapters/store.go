package adapters

import (
	"context"
	"io"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

type Store interface {
	io.Closer

	PutMedia(ctx context.Context, media *entities.Media) error
	PutMediaBulk(ctx context.Context, medias []*entities.Media) error
}

type AnilistStore interface {
	Store

	MappingByAnilistID(ctx context.Context, anilistId string) (*entities.Media, error)
	MappingByAnilistIDBulk(ctx context.Context, anilistIds []string) ([]*entities.Media, error)
}
