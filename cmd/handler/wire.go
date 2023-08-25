package main

import (
	"fmt"
	"time"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/drivers/caches"
	"github.com/wwmoraes/anilistarr/internal/drivers/providers"
	"github.com/wwmoraes/anilistarr/internal/drivers/stores"
	"github.com/wwmoraes/anilistarr/internal/drivers/trackers/anilist"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func NewAnilistBridge(anilistEndpoint string, cacheOptions *caches.RedisOptions) (*usecases.MediaLister, error) {
	tracker := anilist.New(anilistEndpoint, 50)
	if cacheOptions != nil {
		// cache, err := caches.NewRedis(cacheOptions)
		// cache, err := caches.NewBolt("tmp/cache.db", nil)
		cache, err := caches.NewBadger("tmp/badger")
		if err != nil {
			return nil, fmt.Errorf("bolt cache initialization failed: %w", err)
		}

		tracker = &adapters.CachedTracker{
			Cache:   cache,
			Tracker: tracker,
			TTL: adapters.CachedTrackerTTL{
				UserID:       24 * time.Hour,
				MediaListIDs: 5 * time.Minute,
			},
		}
	}

	store, err := stores.NewSQL("sqlite", "tmp/media.db?loc=auto")
	if err != nil {
		return nil, fmt.Errorf("sql store initialization failed: %w", err)
	}

	return &usecases.MediaLister{
		Tracker: tracker,
		Mapper: &adapters.Mapper{
			Provider: providers.AnilistFribbsProvider,
			Store:    store,
		},
	}, nil
}
