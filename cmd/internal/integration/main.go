/*
Integration is a self-contained system test program. It runs both client and
server code to simulate the entire solution and its use-cases.
*/
package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/go-logr/logr"
	"github.com/wwmoraes/gotell/logging"
	_ "go.uber.org/automaxprocs"

	"github.com/wwmoraes/anilistarr/internal/adapters/cachedtracker"
	"github.com/wwmoraes/anilistarr/internal/adapters/sources"
	"github.com/wwmoraes/anilistarr/internal/drivers/badger"
	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/internal/usecases"
	"github.com/wwmoraes/anilistarr/pkg/process"
)

const (
	coverageUsername = "coverage"
	coverageUserID   = 9000
)

//nolint:funlen // TODO refactor integration main func
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
		Source:  sources.JSON[testdata.Metadata](`memory:///test`),
		Store:   store,
	}
	defer process.AssertClose(&mediaLister, "failed to close media lister")

	err = mediaLister.Refresh(ctx, usecases.HTTPGetter(&httpClient{
		Data: map[string]string{
			"memory:///test": `[
				{"source_id": "1", "target_id": "101"},
				{"source_id": "2", "target_id": "102"},
				{"source_id": "3", "target_id": "103"},
				{"source_id": "5", "target_id": "105"},
				{"source_id": "8", "target_id": "108"},
				{"source_id": "13", "target_id": "113"}
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
}
