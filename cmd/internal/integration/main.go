package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/go-logr/logr"
	"github.com/wwmoraes/gotell/logging"
	_ "go.uber.org/automaxprocs"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/drivers/caches"
	"github.com/wwmoraes/anilistarr/internal/drivers/stores"
	"github.com/wwmoraes/anilistarr/internal/test"
	"github.com/wwmoraes/anilistarr/internal/usecases"
	"github.com/wwmoraes/anilistarr/pkg/process"
)

const (
	coverageUsername = "coverage"
	coverageUserId   = 9000
)

//nolint:funlen // TODO refactor integration main func
func main() {
	defer process.HandleExit()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log := logr.New(logging.NewStandardLogSink())

	ctx = logr.NewContext(ctx, log)

	var tracker usecases.Tracker = &test.Tracker{
		UserIds: map[string]int{
			coverageUsername: coverageUserId,
		},
		MediaLists: map[int][]string{
			coverageUserId: {"1", "2", "3", "5", "8", "13"},
		},
	}

	cache, err := caches.NewBadger("", &caches.BadgerOptions{
		InMemory: true,
		// Logger:   &caches.BadgerLogr{Logger: log},
	})
	process.Assert(err)

	store, err := stores.NewBadger("", &stores.BadgerOptions{
		InMemory: true,
		// Logger:   &stores.BadgerLogr{Logger: log},
	})
	process.Assert(err)

	bridge, err := usecases.NewMediaLister(
		&adapters.CachedTracker{
			Cache:   cache,
			Tracker: tracker,
		},
		&adapters.Mapper{
			Provider: test.Provider,
			Store:    store,
		},
	)
	process.Assert(err)
	defer bridge.Close()

	err = bridge.Refresh(ctx, &test.HTTPClient{
		Data: map[string]string{
			test.Provider.String(): `[
				{"anilist_id": 1, "thetvdb_id": 101},
				{"anilist_id": 2, "thetvdb_id": 102},
				{"anilist_id": 3, "thetvdb_id": 103},
				{"anilist_id": 5, "thetvdb_id": 105},
				{"anilist_id": 8, "thetvdb_id": 108},
				{"anilist_id": 13, "thetvdb_id": 113}
			]`,
		},
	})
	process.Assert(err)

	userId, err := bridge.GetUserID(ctx, coverageUsername)
	process.Assert(err)

	log.Info("GetUserID", "username", coverageUsername, "userID", userId)

	customList, err := bridge.Generate(ctx, coverageUsername)
	process.Assert(err)

	log.Info("GenerateCustomList", "username", coverageUsername, "list", customList)
}
