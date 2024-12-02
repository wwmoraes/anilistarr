package usecases

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/hashicorp/go-multierror"
	telemetry "github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/entities"
)

var _ MediaLister = (*MediaList)(nil)

// MediaList handles both a tracker to fetch user data and a mapper that can
// transform the media IDs from the tracker to another service
type MediaList struct {
	Tracker Tracker
	Source  Source
	Store   Store
}

// Generate fetches the user media list from the Tracker and transform the IDs
// found to the target service through the Mapper
func (lister *MediaList) Generate(ctx context.Context, name string) (entities.CustomList, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	if lister.Tracker == nil {
		return nil, ErrStatusFailedPrecondition
	}

	log := telemetry.Logr(ctx).WithValues("username", name)

	log.Info("retrieving user ID")

	userID, err := lister.GetUserID(ctx, name)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to get user ID: %w", err))
	}

	log.Info("retrieving media list IDs", "userID", userID)

	sourceIDs, err := lister.Tracker.GetMediaListIDs(ctx, userID)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to get media list IDs: %w", err))
	}

	targetIDs, err := lister.MapIDs(ctx, sourceIDs)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to get mapped IDs: %w", err))
	}

	customList := make(entities.CustomList, 0, len(targetIDs))

	for _, entry := range targetIDs {
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
func (lister *MediaList) GetUserID(ctx context.Context, name string) (string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	if lister.Tracker == nil {
		return "", ErrStatusFailedPrecondition
	}

	res, err := lister.Tracker.GetUserID(ctx, name)

	return res, span.Assert(err)
}

// Close closes both the Tracker and Mapper
func (lister *MediaList) Close() error {
	closers := [...]io.Closer{
		lister.Store,
		lister.Tracker,
	}
	outErr := make(chan error)

	go func() {
		var wg sync.WaitGroup

		wg.Add(len(closers))

		for _, closer := range closers {
			go func() {
				if closer != nil {
					outErr <- closer.Close()
				}

				wg.Done()
			}()
		}

		wg.Wait()

		// close the channel to signal no more jobs
		close(outErr)
	}()

	// collects all errors
	errs := make([]error, 0, len(closers))
	for err := range outErr {
		errs = append(errs, err)
	}

	//nolint:wrapcheck // caches are internal
	return multierror.Append(nil, errs...).ErrorOrNil()
}

// Refresh requests the Mapper to update its mapping definitions
func (lister *MediaList) Refresh(ctx context.Context, client Getter) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	if lister.Source == nil || lister.Store == nil {
		return ErrStatusFailedPrecondition
	}

	data, err := lister.Source.Fetch(ctx, client)
	if err != nil {
		return span.Assert(fmt.Errorf("failed to refresh anilist mapper: %w", err))
	}

	medias := make([]*entities.Media, 0, len(data))

	for _, entry := range data {
		if !entry.Valid() {
			continue
		}

		medias = append(medias, &entities.Media{
			SourceID: entry.GetSourceID(),
			TargetID: entry.GetTargetID(),
		})
	}

	err = lister.Store.PutMediaBulk(ctx, medias)
	if err != nil {
		return span.Assert(fmt.Errorf("failed to store media during refresh: %w", err))
	}

	return span.Assert(nil)
}

// MapIDs converts IDs between a source tracker and a target reference. Returns
// all IDs that were found, or an empty slice if no matches were found.
func (lister *MediaList) MapIDs(ctx context.Context, ids []entities.SourceID) ([]entities.TargetID, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	if lister.Store == nil {
		return nil, ErrStatusFailedPrecondition
	}

	records, err := lister.Store.GetMediaBulk(ctx, ids)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to map IDs: %w", err))
	}

	targetIDs := make([]string, 0, len(records))
	for _, record := range records {
		targetIDs = append(targetIDs, record.TargetID)
	}

	return targetIDs, span.Assert(nil)
}
