/*
Integration is a self-contained system test program. It runs both client and
server code to simulate the entire solution and its use-cases.
*/
package main

import (
	"context"
	"os"
	"os/signal"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	"github.com/wwmoraes/gotell/logging"

	"github.com/wwmoraes/anilistarr/internal/adapters/cachedtracker"
	"github.com/wwmoraes/anilistarr/internal/adapters/sources"
	"github.com/wwmoraes/anilistarr/internal/drivers/animelists"
	"github.com/wwmoraes/anilistarr/internal/drivers/badger"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/usecases"
	"github.com/wwmoraes/anilistarr/pkg/process"
)

const (
	coverageUsername = "coverage"
	coverageUserID   = 9000
)

//nolint:funlen,maintidx // TODO refactor integration main func
func main() {
	defer process.HandleExit()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log := logr.New(logging.NewStandardLogSink())

	ctx = logr.NewContext(ctx, log)

	tracker := memoryTracker{
		UserIDs: map[string]int{
			coverageUsername: coverageUserID,
		},
		MediaLists: map[int][]string{
			coverageUserID: {"1", "2", "3", "5", "8", "13"},
		},
	}

	cachePath, err := os.MkdirTemp("", "anilistarr-integration-badger-cache-*")
	process.AssertWith(err, "failed to create temporary directory")
	defer os.RemoveAll(cachePath)

	cache, err := badger.New(
		cachePath,
		badger.WithInMemory(true),
		// badger.WithLogger(&badger.Logr{Logger: log}),
	)
	process.AssertWith(err, "failed to create badger driver")

	storePath, err := os.MkdirTemp("", "anilistarr-integration-badger-store-*")
	process.AssertWith(err, "failed to create temporary directory")
	defer os.RemoveAll(storePath)

	store, err := badger.New(
		storePath,
		badger.WithInMemory(true),
		// badger.WithLogger(&badger.Logr{Logger: log}),
	)
	process.Assert(err)

	cachedTracker := cachedtracker.CachedTracker{
		Cache:   cache,
		Tracker: &tracker,
		TTL: cachedtracker.TTLs{
			UserID:       time.Hour,
			MediaListIDs: time.Hour,
		},
	}

	mediaLister := usecases.MediaList{
		Tracker: &cachedTracker,
		Source:  sources.JSON[animelists.Anilist2TVDBMetadata](`memory:///test`),
		Store:   store,
	}
	defer process.AssertClose(&mediaLister, "failed to close media lister")

	err = mediaLister.Refresh(ctx, usecases.HTTPGetter(&httpClient{
		Data: map[string]string{
			"memory:///test": `[
				{"anilist_id": 1, "thetvdb_id": 101},
				{"anilist_id": 2, "thetvdb_id": 102},
				{"anilist_id": 3, "thetvdb_id": 103},
				{"anilist_id": 5, "thetvdb_id": 105},
				{"anilist_id": 8, "thetvdb_id": 108},
				{"anilist_id": 13, "thetvdb_id": 113}
			]`,
		},
	}))
	process.Assert(err)

	userID, err := mediaLister.GetUserID(ctx, coverageUsername)
	process.Assert(err)

	log.Info("GetUserID", "username", coverageUsername, "userID", userID)

	customList, err := mediaLister.Generate(ctx, coverageUsername)
	process.Assert(err)

	log.Info("GenerateCustomList", "username", coverageUsername, "list", customList)

	//nolint:mnd // test data
	wantedCustomList := entities.CustomList{
		entities.CustomEntry{TvdbID: 101},
		entities.CustomEntry{TvdbID: 102},
		entities.CustomEntry{TvdbID: 103},
		entities.CustomEntry{TvdbID: 105},
		entities.CustomEntry{TvdbID: 108},
		entities.CustomEntry{TvdbID: 113},
	}

	if !reflect.DeepEqual(customList, wantedCustomList) {
		process.AssertWith(usecases.ErrStatusUnknown, "custom list does not matches expectations")
	}
}
