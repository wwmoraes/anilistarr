package main

import (
	"context"
	"fmt"
	"path"
	"time"

	telemetry "github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/drivers/badger"
	"github.com/wwmoraes/anilistarr/internal/drivers/bolt"
	"github.com/wwmoraes/anilistarr/internal/drivers/redis"
	"github.com/wwmoraes/anilistarr/internal/drivers/sqlite"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

const (
	cacheUserTTL      = 24 * time.Hour
	cacheMediaListTTL = time.Hour
	anilistPageSize   = 50
	storeType         = StoreBadger
	cacheType         = CacheBadger
)

func newStore(ctx context.Context, dataPath string) (usecases.Store, error) {
	var store usecases.Store

	var err error

	// TODO kill this switch...
	switch storeType {
	case StoreBadger:
		store, err = badger.New(
			path.Join(dataPath, "badger", "store"),
			badger.WithLogger(&badger.Logr{
				Logger: telemetry.Logr(ctx),
			}),
		)
	case StoreSQL:
		store, err = sqlite.New(ctx, path.Join(dataPath, "sqlite-store.db"))
	default:
		return nil, usecases.ErrStatusUnimplemented
	}

	if err != nil {
		return nil, fmt.Errorf("store initialization failed: %w", err)
	}

	return store, nil
}

func newCache(ctx context.Context, dataPath string) (usecases.Cache, error) {
	var cache usecases.Cache

	var err error

	// TODO kill this switch...
	switch cacheType {
	case CacheBadger:
		cache, err = badger.New(
			path.Join(dataPath, "badger", "cache"),
			badger.WithLogger(&badger.Logr{
				Logger: telemetry.Logr(ctx),
			}),
		)
	case CacheBolt:
		cache, err = bolt.New(path.Join(dataPath, "bolt-cache.db"), nil)
	case CacheRedis:
		cache, err = redis.New(ctx, nil)
	case CacheSQL:
		cache, err = sqlite.New(ctx, path.Join(dataPath, "sqlite-cache.db"))
	default:
		return nil, usecases.ErrStatusUnimplemented
	}

	if err != nil {
		return nil, fmt.Errorf("cache initialization failed: %w", err)
	}

	return cache, nil
}
