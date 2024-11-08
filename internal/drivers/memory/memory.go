package memory

import (
	"context"
	"errors"

	telemetry "github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

type Memory map[string]string

func New() Memory {
	return make(Memory)
}

func (mem Memory) Close() error {
	for k := range mem {
		delete(mem, k)
	}

	return nil
}

func (mem Memory) GetString(ctx context.Context, key string) (string, error) {
	_, span := telemetry.Start(ctx)
	defer span.End()

	if mem == nil {
		return "", errors.ErrUnsupported
	}

	value, ok := mem[key]
	if !ok {
		return "", usecases.ErrNotFound
	}

	return value, nil
}

func (mem Memory) SetString(ctx context.Context, key, value string, options ...adapters.CacheOption) error {
	_, span := telemetry.Start(ctx)
	defer span.End()

	if mem == nil {
		return errors.ErrUnsupported
	}

	mem[key] = value

	return nil
}

func (mem Memory) GetMedia(ctx context.Context, id string) (*entities.Media, error) {
	_, span := telemetry.Start(ctx)
	defer span.End()

	if mem == nil {
		return nil, errors.ErrUnsupported
	}

	targetID, ok := mem[id]
	if !ok {
		return nil, usecases.ErrNotFound
	}

	return &entities.Media{
		SourceID: id,
		TargetID: targetID,
	}, nil
}

func (mem Memory) GetMediaBulk(ctx context.Context, ids []string) ([]*entities.Media, error) {
	_, span := telemetry.Start(ctx)
	defer span.End()

	if mem == nil {
		return nil, errors.ErrUnsupported
	}

	medias := make([]*entities.Media, 0, len(ids))

	for _, sourceID := range ids {
		targetID, ok := mem[sourceID]
		if !ok {
			continue
		}

		medias = append(medias, &entities.Media{
			SourceID: sourceID,
			TargetID: targetID,
		})
	}

	return medias, nil
}

func (mem Memory) PutMedia(ctx context.Context, media *entities.Media) error {
	_, span := telemetry.Start(ctx)
	defer span.End()

	if mem == nil {
		return errors.ErrUnsupported
	}

	if media.SourceID == "" {
		return usecases.ErrBadRequest
	}

	mem[media.SourceID] = media.TargetID

	return nil
}

func (mem Memory) PutMediaBulk(ctx context.Context, medias []*entities.Media) error {
	_, span := telemetry.Start(ctx)
	defer span.End()

	if mem == nil {
		return errors.ErrUnsupported
	}

	// we do double loop to simulate a "transaction"
	for _, media := range medias {
		if media.SourceID == "" {
			return usecases.ErrBadRequest
		}
	}

	for _, media := range medias {
		mem[media.SourceID] = media.TargetID
	}

	return nil
}
