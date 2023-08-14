package usecases

import (
	"context"
	"fmt"
	"strconv"

	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
)

type MediaBridge struct {
	Tracker Tracker
	Mapper  Mapper
}

func (linker *MediaBridge) GenerateCustomList(ctx context.Context, name string) (entities.SonarrCustomList, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()
	log := telemetry.LoggerFromContext(ctx).WithValues("username", name)

	log.Info("retrieving user ID")
	userId, err := linker.GetUserID(ctx, name)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to get user ID: %w", err))
	}

	log.Info("retrieving media list IDs", "userID", userId)
	sourceIds, err := linker.Tracker.GetMediaListIDs(ctx, userId)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to get media list IDs: %w", err))
	}

	targetIds, err := linker.Mapper.MapIDs(ctx, sourceIds)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to get mapped IDs: %w", err))
	}

	customList := make(entities.SonarrCustomList, 0, len(targetIds))
	for index, entry := range targetIds {
		if entry == "" {
			log.Info("no TVDB ID registered for source ID", "sourceID", sourceIds[index])
			continue
		}

		tvdbID, err := strconv.ParseUint(entry, 10, 0)
		if err != nil {
			return nil, span.Assert(fmt.Errorf("failed to parse TVDB ID: %w", err))
		}

		customList = append(customList, entities.SonarrCustomEntry{
			TvdbID: tvdbID,
		})
	}

	return customList, span.Assert(nil)
}

func (linker *MediaBridge) GetUserID(ctx context.Context, name string) (string, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	res, err := linker.Tracker.GetUserID(ctx, name)
	return res, span.Assert(err)
}

func (linker *MediaBridge) Close() error {
	errT := linker.Tracker.Close()
	errR := linker.Mapper.Close()

	if errT != nil || errR != nil {
		return fmt.Errorf("failed to close mapper dependencies: %v", []error{errT, errR})
	}

	return nil
}

func (linker *MediaBridge) Refresh(ctx context.Context) error {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	return span.Assert(linker.Mapper.Refresh(ctx))
}
