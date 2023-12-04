package adapters

import (
	"context"
	"fmt"

	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

// Mapper handles the mapping of media data. It also handles the retrieval of
// data from a provider and its persistent storage
type Mapper struct {
	Provider Provider[Metadata]
	Store    Store
}

func (mapper *Mapper) Close() error {
	return mapper.Store.Close()
}

func (mapper *Mapper) MapIDs(ctx context.Context, ids []entities.SourceID) ([]entities.TargetID, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	records, err := mapper.Store.GetMediaBulk(ctx, ids)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to map IDs: %w", err))
	}

	targetIds := make([]string, len(records))
	for index, record := range records {
		targetIds[index] = record.TargetID
	}

	// for _, sourceId := range anilistIds {
	// 	targetId, err := mapper.MapID(ctx, sourceId)
	// 	if err != nil {
	// 		return nil, span.Assert(fmt.Errorf("failed to map IDs: %w", err))
	// 	}

	// 	ids = append(ids, targetId)
	// }

	return targetIds, span.Assert(nil)
}

func (mapper *Mapper) Refresh(ctx context.Context, client usecases.Getter) error {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	data, err := mapper.Provider.Fetch(ctx, client)
	if err != nil {
		return span.Assert(fmt.Errorf("failed to refresh anilist mapper: %w", err))
	}

	medias := make([]*entities.Media, 0, len(data))
	for _, entry := range data {
		if entry.GetTargetID() == "0" || entry.GetSourceID() == "0" {
			continue
		}

		medias = append(medias, &entities.Media{
			SourceID: entry.GetSourceID(),
			TargetID: entry.GetTargetID(),
		})
	}

	err = mapper.Store.PutMediaBulk(ctx, medias)
	if err != nil {
		return span.Assert(fmt.Errorf("failed to store media during refresh: %w", err))
	}

	return span.Assert(nil)
}
