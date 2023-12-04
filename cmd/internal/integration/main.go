package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"github.com/wwmoraes/anilistarr/internal/test"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

const (
	coverageUsername = "coverage"
	coverageUserId   = 9000
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log := telemetry.DefaultLogger()

	var tracker usecases.Tracker = &test.Tracker{
		UserIds: map[string]int{
			coverageUsername: coverageUserId,
		},
		MediaLists: map[int][]string{
			coverageUserId: {"1", "2", "3", "5", "8", "13"},
		},
	}

	bridge, err := usecases.NewMediaLister(
		&adapters.CachedTracker{
			Cache: &test.Cache{},
			Tracker: tracker,
		},
		&adapters.Mapper{
			Provider: test.Provider,
			Store:    &test.Store{},
		},
	)
	assert(err)
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
	assert(err)

	userId, err := bridge.GetUserID(ctx, coverageUsername)
	assert(err)

	log.Info("GetUserID", "username", coverageUsername, "userID", userId)

	customList, err := bridge.Generate(ctx, coverageUsername)
	assert(err)

	log.Info("GenerateCustomList", "username", coverageUsername, "list", customList)
}

func assert(err error) {
	if err == nil {
		return
	}

	log := telemetry.DefaultLogger()

	log.Error(err, "assertion failed")
	os.Exit(1)
}
