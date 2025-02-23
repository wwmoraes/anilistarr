// Package cachedtracker provides a transparent [usecases.Tracker] wrapper that
// enables cache-first responses.
package cachedtracker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	telemetry "github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

const (
	cacheKeyUserID     string = "anilist:user:%s:id"
	cacheKeyUserMedia  string = "anilist:user:%s:media"
	mediaListSeparator string = "|"
)

var _ usecases.Tracker = (*CachedTracker)(nil)

// CachedTracker is a meta-tracker that provides cached responses.
//
// It is an [usecases.Tracker] drop-in replacement that wraps another
// tracker with an [adapters.Cache]. It tries to use cached information first,
// falling back to the original tracker on miss.
//
// Misses automatically update the cache with the tracker response.
//
// It sets a TTL for each cached response for implementations that supports it.
type CachedTracker struct {
	Cache   usecases.Cache
	Tracker usecases.Tracker
	TTL     TTLs
}

// TTLs contains the time-to-live for entries of each type handled by trackers.
//
// TODO refactor to use [adapters.CacheParams] instead of [time.Duration]
type TTLs struct {
	UserID       time.Duration
	MediaListIDs time.Duration
}

// GetUserID retrieves an user ID for a given name. It returns a cache value if
// available; otherwise it queries the tracker and caches it for future use.
func (wrapper *CachedTracker) GetUserID(ctx context.Context, name string) (string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	key := fmt.Sprintf(cacheKeyUserID, name)

	span.AddEvent("try cache")

	userID, err := wrapper.Cache.GetString(ctx, key)
	if err != nil && !errors.Is(err, usecases.ErrStatusNotFound) {
		return "", span.Assert(errors.Join(usecases.ErrStatusUnknown, err))
	}

	if userID != "" {
		span.AddEvent("cache hit")

		return userID, span.Assert(nil)
	}

	span.AddEvent("cache miss")

	userID, err = wrapper.Tracker.GetUserID(ctx, name)
	if err != nil {
		return "", span.Assert(errors.Join(usecases.ErrStatusUnknown, err))
	}

	err = wrapper.Cache.SetString(
		ctx,
		key,
		userID,
		usecases.WithTTL(wrapper.TTL.UserID),
	)

	return userID, span.Assert(err)
}

// GetMediaListIDs retrieves the list of medias for an user ID. It returns a
// cache value if available; otherwise it requests the tracker and caches it for
// future use.
func (wrapper *CachedTracker) GetMediaListIDs(ctx context.Context, userID string) ([]string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	key := fmt.Sprintf(cacheKeyUserMedia, userID)

	span.AddEvent("try cache")

	cachedIDs, err := wrapper.Cache.GetString(ctx, key)
	if err != nil && !errors.Is(err, usecases.ErrStatusNotFound) {
		return nil, span.Assert(errors.Join(usecases.ErrStatusUnknown, err))
	}

	if cachedIDs != "" {
		span.AddEvent("cache hit")

		return strings.Split(cachedIDs, mediaListSeparator), span.Assert(err)
	}

	span.AddEvent("cache miss")

	ids, err := wrapper.Tracker.GetMediaListIDs(ctx, userID)
	if err != nil {
		return nil, span.Assert(errors.Join(usecases.ErrStatusUnknown, err))
	}

	err = wrapper.Cache.SetString(
		ctx,
		key,
		strings.Join(ids, mediaListSeparator),
		usecases.WithTTL(wrapper.TTL.MediaListIDs),
	)

	return ids, span.Assert(err)
}

// Close terminates the client and its connection to the cache.
func (wrapper *CachedTracker) Close() error {
	closers := [...]io.Closer{
		wrapper.Cache,
		wrapper.Tracker,
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

	//nolint:wrapcheck // components are internal
	return multierror.Append(nil, errs...).ErrorOrNil()
}
