package main

import (
	"context"
	"fmt"
	"path"
	"time"

	telemetry "github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/drivers/caches"
	"github.com/wwmoraes/anilistarr/internal/drivers/providers"
	"github.com/wwmoraes/anilistarr/internal/drivers/stores"
	"github.com/wwmoraes/anilistarr/internal/drivers/trackers/anilist"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

const (
	cacheUserTTL      = 24 * time.Hour
	cacheMediaListTTL = time.Hour
	anilistPageSize   = 50
)

func NewAnilistMediaLister(anilistEndpoint string, store adapters.Store, cache adapters.Cache) (usecases.MediaLister, error) {
	tracker := anilist.New(anilistEndpoint, anilistPageSize)

	if cache != nil {
		tracker = &adapters.CachedTracker{
			Cache:   cache,
			Tracker: tracker,
			TTL: adapters.CachedTrackerTTL{
				UserID:       cacheUserTTL,
				MediaListIDs: cacheMediaListTTL,
			},
		}
	}

	return usecases.NewMediaLister(tracker, &adapters.Mapper{
		Provider: providers.AnilistFribbsProvider,
		Store:    store,
	})
}

func NewStore(dataPath string) (adapters.Store, error) {
	// store, err := stores.NewSQL("sqlite", path.Join(dataPath, "media2.db?loc=auto"))
	store, err := stores.NewBadger(path.Join(dataPath, "badger", "store"), &stores.BadgerOptions{
		Logger: &caches.BadgerLogr{Logger: telemetry.Logr(context.TODO())},
	})
	if err != nil {
		return nil, fmt.Errorf("store initialization failed: %w", err)
	}

	return store, nil
}

func NewCache(dataPath string) (adapters.Cache, error) {
	// cache, err := caches.NewRedis(cacheOptions)
	// cache, err := caches.NewBolt(path.Join(dataPath, "bolt-cache.db"), nil)
	// cache, err := caches.NewFile("tmp/cache.txt")
	cache, err := caches.NewBadger(path.Join(dataPath, "badger", "cache"), &caches.BadgerOptions{
		Logger: &caches.BadgerLogr{Logger: telemetry.Logr(context.TODO())},
	})
	if err != nil {
		return nil, fmt.Errorf("cache initialization failed: %w", err)
	}

	return cache, nil
}
