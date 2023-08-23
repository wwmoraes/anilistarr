package adapters

import (
	"context"
	"fmt"

	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

const (
	cacheKeyUserID string = "anilist:user:%s:id"
)

type cachedTracker struct {
	cache   Cache
	tracker usecases.Tracker
}

func NewCachedTracker(tracker usecases.Tracker, cache Cache) (usecases.Tracker, error) {
	return &cachedTracker{
		cache:   cache,
		tracker: tracker,
	}, nil
}

func (wrapper *cachedTracker) GetUserID(ctx context.Context, name string) (string, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	key := fmt.Sprintf(cacheKeyUserID, name)

	span.AddEvent("try cache")
	userId, err := wrapper.cache.GetString(ctx, key)
	if err != nil {
		return "", span.Assert(fmt.Errorf("failed to get user ID: %w", err))
	}

	if userId != "" {
		span.AddEvent("cache hit")
		return userId, span.Assert(err)
	}

	span.AddEvent("cache miss")
	userId, err = wrapper.tracker.GetUserID(ctx, name)
	if err != nil {
		return "", span.Assert(fmt.Errorf("failed to get user ID: %w", err))
	}

	return userId, span.Assert(wrapper.cache.SetString(ctx, key, userId))
}

func (wrapper *cachedTracker) GetMediaListIDs(ctx context.Context, userId string) ([]string, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	// TODO use cache to avoid DB DDoS/increased costs
	ids, err := wrapper.tracker.GetMediaListIDs(ctx, userId)

	return ids, span.Assert(err)
}

func (wrapper *cachedTracker) Close() error {
	return wrapper.cache.Close()
}
