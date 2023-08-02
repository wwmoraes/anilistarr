package adapters

import (
	"context"
	"fmt"

	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
)

type AnilistMapper struct {
	Source MetadataSource[Metadata]
	Store  AnilistStore
}

func (mapper *AnilistMapper) Close() error {
	return mapper.Store.Close()
}

func (mapper *AnilistMapper) MapID(ctx context.Context, anilistId string) (string, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	media, err := mapper.Store.MappingByAnilistID(ctx, anilistId)

	if err != nil {
		return "", span.Assert(fmt.Errorf("failed to map ID %s: %w", anilistId, err))
	}

	if media == nil {
		return "", span.Assert(nil)
	}

	return media.TvdbID, span.Assert(nil)
}

func (mapper *AnilistMapper) MapIDs(ctx context.Context, anilistIds []string) ([]string, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	records, err := mapper.Store.MappingByAnilistIDBulk(ctx, anilistIds)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to map IDs: %w", err))
	}

	ids := make([]string, len(records))
	for index, record := range records {
		ids[index] = record.TvdbID
	}

	// for _, sourceId := range anilistIds {
	// 	targetId, err := mapper.MapID(ctx, sourceId)
	// 	if err != nil {
	// 		return nil, span.Assert(fmt.Errorf("failed to map IDs: %w", err))
	// 	}

	// 	ids = append(ids, targetId)
	// }

	return ids, span.Assert(nil)
}

func (mapper *AnilistMapper) Refresh(ctx context.Context) error {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	data, err := mapper.Source.Fetch(ctx, nil)
	if err != nil {
		return span.Assert(fmt.Errorf("failed to refresh anilist mapper: %w", err))
	}

	medias := make([]*entities.Media, 0, len(data))
	for _, entry := range data {
		if entry.GetTvdbID() == "0" || entry.GetAnilistID() == "0" {
			continue
		}

		medias = append(medias, &entities.Media{
			AnilistID: entry.GetAnilistID(),
			TvdbID:    entry.GetTvdbID(),
		})
	}

	err = mapper.Store.PutMediaBulk(ctx, medias)
	if err != nil {
		return span.Assert(fmt.Errorf("failed to store media during refresh: %w", err))
	}

	return span.Assert(nil)
}
