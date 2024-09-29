package usecases

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	telemetry "github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

var (
	ErrNoTracker = errors.New("no tracker set up")
	ErrNoMapper  = errors.New("no mapper set up")
	ErrNoCache   = errors.New("no cache set up")
)

type MediaLister interface {
	// Generate fetches the user media list from the Tracker and transform the IDs
	// found to the target service through the Mapper
	Generate(ctx context.Context, name string) (entities.CustomList, error)

	// GetUserID searches the Tracker for the user ID by their name/handle
	GetUserID(ctx context.Context, name string) (string, error)

	// Close closes both the Tracker and Mapper
	Close() error

	// Refresh requests the Mapper to update its mapping definitions
	Refresh(ctx context.Context, client Getter) error
}

func NewMediaLister(tracker Tracker, mapper Mapper) (MediaLister, error) {
	if tracker == nil {
		return nil, ErrNoTracker
	}

	if mapper == nil {
		return nil, ErrNoMapper
	}

	return &mediaLister{
		Tracker: tracker,
		Mapper:  mapper,
	}, nil
}

// mediaLister handles both a tracker to fetch user data and a mapper that can
// transform the media IDs from the tracker to another service
type mediaLister struct {
	Tracker Tracker
	Mapper  Mapper
}

// Generate fetches the user media list from the Tracker and transform the IDs
// found to the target service through the Mapper
func (lister *mediaLister) Generate(ctx context.Context, name string) (entities.CustomList, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	log := telemetry.Logr(ctx).WithValues("username", name)

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

// GetUserID searches the Tracker for the user ID by their name/handle
func (lister *mediaLister) GetUserID(ctx context.Context, name string) (string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := lister.Tracker.GetUserID(ctx, name)

	return res, span.Assert(err)
}

// Close closes both the Tracker and Mapper
func (lister *mediaLister) Close() error {
	return errors.Join(lister.Tracker.Close(), lister.Mapper.Close())
}

// Refresh requests the Mapper to update its mapping definitions
func (lister *mediaLister) Refresh(ctx context.Context, client Getter) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	return span.Assert(lister.Mapper.Refresh(ctx, client))
}
