package usecases

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
)

var (
	NoTrackerError = errors.New("no tracker set up")
)

type MediaLister struct {
	Tracker Tracker
	Mapper  Mapper
}

func (lister *MediaLister) Generate(ctx context.Context, name string) (entities.CustomList, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()
	log := telemetry.LoggerFromContext(ctx).WithValues("username", name)

	log.Info("retrieving user ID")
	userId, err := lister.GetUserID(ctx, name)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to get user ID: %w", err))
	}

	log.Info("retrieving media list IDs", "userID", userId)
	sourceIds, err := lister.Tracker.GetMediaListIDs(ctx, userId)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to get media list IDs: %w", err))
	}

	targetIds, err := lister.Mapper.MapIDs(ctx, sourceIds)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to get mapped IDs: %w", err))
	}

	customList := make(entities.CustomList, 0, len(targetIds))
	for index, entry := range targetIds {
		if entry == "" {
			log.Info("no TVDB ID registered for source ID", "sourceID", sourceIds[index])
			continue
		}

		tvdbID, err := strconv.ParseUint(entry, 10, 0)
		if err != nil {
			return nil, span.Assert(fmt.Errorf("failed to parse TVDB ID: %w", err))
		}

		customList = append(customList, entities.CustomEntry{
			TvdbID: tvdbID,
		})
	}

	return customList, span.Assert(nil)
}

func (lister *MediaLister) GetUserID(ctx context.Context, name string) (string, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	if lister.Tracker == nil {
		return "", NoTrackerError
	}

	res, err := lister.Tracker.GetUserID(ctx, name)
	return res, span.Assert(err)
}

func (lister *MediaLister) Close() error {
	errT := lister.Tracker.Close()
	errR := lister.Mapper.Close()

	if errT != nil || errR != nil {
		return fmt.Errorf("failed to close dependencies: %v", []error{errT, errR})
	}

	return nil
}

func (lister *MediaLister) Refresh(ctx context.Context, client Getter) error {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	return span.Assert(lister.Mapper.Refresh(ctx, client))
}