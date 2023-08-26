package adapters

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

const (
	cacheKeyUserID     string = "anilist:user:%s:id"
	cacheKeyUserMedia  string = "anilist:user:%s:media"
	mediaListSeparator string = "|"
)

// CachedTracker wraps a Tracker with a passthrough cache
type CachedTracker struct {
	Cache   Cache
	Tracker usecases.Tracker
	TTL     CachedTrackerTTL
}

type CachedTrackerTTL struct {
	UserID       time.Duration
	MediaListIDs time.Duration
}

func (wrapper *CachedTracker) GetUserID(ctx context.Context, name string) (string, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	key := fmt.Sprintf(cacheKeyUserID, name)

	span.AddEvent("try cache")
	userId, err := wrapper.Cache.GetString(ctx, key)
	if err != nil {
		return "", span.Assert(fmt.Errorf("failed to get user ID: %w", err))
	}

	if userId != "" {
		span.AddEvent("cache hit")
		return userId, span.Assert(err)
	}

	span.AddEvent("cache miss")
	userId, err = wrapper.Tracker.GetUserID(ctx, name)
	if err != nil {
		return "", span.Assert(fmt.Errorf("failed to get user ID: %w", err))
	}

	err = wrapper.Cache.SetString(
		ctx,
		key,
		userId,
		WithTTL(wrapper.TTL.UserID),
	)

	return userId, span.Assert(err)
}

func (wrapper *CachedTracker) GetMediaListIDs(ctx context.Context, userId string) ([]string, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	key := fmt.Sprintf(cacheKeyUserMedia, userId)

	span.AddEvent("try cache")
	cachedIds, err := wrapper.Cache.GetString(ctx, key)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to get media list: %w", err))
	}

	if cachedIds != "" {
		span.AddEvent("cache hit")
		return strings.Split(cachedIds, mediaListSeparator), span.Assert(err)
	}

	span.AddEvent("cache miss")
	ids, err := wrapper.Tracker.GetMediaListIDs(ctx, userId)
	if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to get user ID: %w", err))
	}

	err = wrapper.Cache.SetString(
		ctx,
		key,
		strings.Join(ids, mediaListSeparator),
		WithTTL(wrapper.TTL.MediaListIDs),
	)

	return ids, span.Assert(err)
}

func (wrapper *CachedTracker) Close() error {
	return wrapper.Cache.Close()
}
