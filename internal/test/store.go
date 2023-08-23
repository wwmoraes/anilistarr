package test

import (
	"context"
	"fmt"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

type Store map[string]*entities.Media

func (store *Store) Close() error {
	*store = Store{}

	return nil
}

func (store Store) GetMedia(ctx context.Context, id string) (*entities.Media, error) {
	value, ok := store[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	return value, nil
}

func (store Store) GetMediaBulk(ctx context.Context, ids []string) ([]*entities.Media, error) {
	values := make([]*entities.Media, 0, len(ids))

	for _, id := range ids {
		value, ok := store[id]
		if !ok {
			continue
		}

		values = append(values, value)
	}

	return values, nil
}

func (store Store) PutMedia(ctx context.Context, media *entities.Media) error {
	store[media.SourceID] = media

	return nil
}

func (store Store) PutMediaBulk(ctx context.Context, medias []*entities.Media) error {
	for _, media := range medias {
		store[media.SourceID] = media
	}

	return nil
}
