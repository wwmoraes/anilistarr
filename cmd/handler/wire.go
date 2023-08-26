package main

import (
	"fmt"
	"path"
	"time"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/drivers/caches"
	"github.com/wwmoraes/anilistarr/internal/drivers/providers"
	"github.com/wwmoraes/anilistarr/internal/drivers/stores"
	"github.com/wwmoraes/anilistarr/internal/drivers/trackers/anilist"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func NewAnilistBridge(anilistEndpoint string, dataPath string) (*usecases.MediaLister, error) {
	tracker := anilist.New(anilistEndpoint, 50)

	// cache, err := caches.NewRedis(cacheOptions)
	// cache, err := caches.NewBolt(path.Join(dataPath, "bolt-cache.db"), nil)
	cache, err := caches.NewBadger(path.Join(dataPath, "badger", "cache"), caches.WithLogger(telemetry.DefaultLogger()))
	if err != nil {
		return nil, fmt.Errorf("bolt cache initialization failed: %w", err)
	}

	tracker = &adapters.CachedTracker{
		Cache:   cache,
		Tracker: tracker,
		TTL: adapters.CachedTrackerTTL{
			UserID:       24 * time.Hour,
			MediaListIDs: 60 * time.Minute,
		},
	}

	// store, err := stores.NewSQL("sqlite", path.Join(dataPath, "media2.db?loc=auto"))
	store, err := stores.NewBadger(path.Join(dataPath, "badger", "store"))
	if err != nil {
		return nil, fmt.Errorf("store initialization failed: %w", err)
	}

	return &usecases.MediaLister{
		Tracker: tracker,
		Mapper: &adapters.Mapper{
			Provider: providers.AnilistFribbsProvider,
			Store:    store,
		},
	}, nil
}
