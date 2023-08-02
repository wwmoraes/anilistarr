package main

import (
	"fmt"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/drivers/anilist"
	"github.com/wwmoraes/anilistarr/internal/drivers/caches"
	"github.com/wwmoraes/anilistarr/internal/drivers/providers"
	"github.com/wwmoraes/anilistarr/internal/drivers/stores"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func NewAnilistLinker(anilistEndpoint string, cacheOptions *caches.RedisOptions) (*usecases.MediaLinker, error) {
	tracker, err := anilist.New(anilistEndpoint, 50)
	if err != nil {
		return nil, fmt.Errorf("anilist tracker initialization failed: %w", err)
	}

	if cacheOptions != nil {
		// cache, err := caches.NewRedis(cacheOptions)
		cache, err := caches.NewBolt("tmp/cache.db", nil)
		if err != nil {
			return nil, fmt.Errorf("bolt cache initialization failed: %w", err)
		}

		tracker, err = adapters.NewCachedTracker(tracker, cache)
		if err != nil {
			return nil, fmt.Errorf("cached adapter initialization failed: %w", err)
		}
	}

	store, err := stores.NewSQL("sqlite", "tmp/media.db?loc=auto")
	if err != nil {
		return nil, fmt.Errorf("sql store initialization failed: %w", err)
	}

	linker, err := adapters.NewAnilistMapper(providers.FribbsSource, store)
	if err != nil {
		return nil, fmt.Errorf("anilist linker initialization failed: %w", err)
	}

	return usecases.NewMediaMapper(
		tracker,
		linker,
	)
}
