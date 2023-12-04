package usecases_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/drivers/stores"
	"github.com/wwmoraes/anilistarr/internal/test"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

const (
	testUsername = "test"
	testUserId   = 1234
)

var (
	testSourceIds = []string{"1", "2", "3", "5", "8", "13"}
	testTargetIds = []string{"91", "92", "93", "95", "98", "913"}

	// drivers

	testClient = test.HTTPClient{
		Data: map[string]string{
			test.Provider.String(): `[
				{"anilist_id": 1, "thetvdb_id": 91},
				{"anilist_id": 2, "thetvdb_id": 92},
				{"anilist_id": 3, "thetvdb_id": 93},
				{"anilist_id": 5, "thetvdb_id": 95},
				{"anilist_id": 8, "thetvdb_id": 98},
				{"anilist_id": 13, "thetvdb_id": 913}
			]`,
		},
	}
	testTracker = &test.Tracker{
		UserIds: map[string]int{
			testUsername: testUserId,
		},
		MediaLists: map[int][]string{
			testUserId: testSourceIds,
		},
	}
)

func TestMediaBridge(t *testing.T) {
	store, err := stores.NewBadger("", &stores.BadgerOptions{
		InMemory: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	bridge, err := usecases.NewMediaLister(testTracker,
		&adapters.Mapper{
			Provider: test.Provider,
			Store:    store,
		},
	)
	if err != nil {
		t.Error("unexpected error when creating MediaLister:", err)
	}

	wantedCustomList := test.SonarrCustomListFromIDs(t, "91", "92", "93", "95", "98", "913")

	ctx := context.Background()
	err = bridge.Refresh(ctx, &testClient)
	if err != nil {
		t.Error("unexpected error on Refresh:", err)
	}

	customList, err := bridge.Generate(ctx, testUsername)
	if err != nil {
		t.Error("unexpected error on GetUserID:", err)
	}
	if !reflect.DeepEqual(customList, wantedCustomList) {
		t.Errorf("custom list does not match: got '%v', expected '%v'", customList, wantedCustomList)
	}

	err = bridge.Close()
	if err != nil {
		t.Error("unexpected error on Close:", err)
	}
}

func TestNewMediaListerNoTracker(t *testing.T) {
	_, err := usecases.NewMediaLister(nil, &adapters.Mapper{})
	if !errors.Is(err, usecases.ErrNoTracker) {
		t.Errorf("expected %q, got %q", usecases.ErrNoTracker, err)
	}
}

func TestNewMediaListerNoMapper(t *testing.T) {
	_, err := usecases.NewMediaLister(&test.Tracker{}, nil)
	if !errors.Is(err, usecases.ErrNoMapper) {
		t.Errorf("expected %q, got %q", usecases.ErrNoMapper, err)
	}
}

// func TestMediaBridge_GetUserID(t *testing.T) {
// 	tracker := &test.Tracker{
// 		UserIds: map[string]int{
// 			testUsername: testUserId,
// 		},
// 		MediaLists: map[int][]string{
// 			testUserId: testSourceIds,
// 		},
// 	}

// 	bridge := usecases.MediaBridge{
// 		Tracker: tracker,
// 		Mapper: &adapters.Mapper{
// 			Provider: testProvider,
// 			Store:    &test.Store{},
// 		},
// 	}

// 	ctx := context.Background()

// 	userId, err := bridge.GetUserID(ctx, "test")
// 	if err != nil {
// 		t.Error("unexpected error:", err)
// 	}

// 	if userId != strconv.Itoa(testUserId) {
// 		t.Errorf("user id does not match: got '%s', expected '%d'", userId, testUserId)
// 	}
// }

// func TestMediaBridge_GenerateCustomList(t *testing.T) {
// 	wantedIds := test.SonarrCustomListFromIDs(t, testTargetIds...)

// 	bridge := usecases.MediaBridge{
// 		Tracker: &test.Tracker{
// 			UserIds: map[string]int{
// 				testUsername: testUserId,
// 			},
// 			MediaLists: map[int][]string{
// 				testUserId: testSourceIds,
// 			},
// 		},
// 		Mapper: &adapters.Mapper{
// 			Provider: test.Provider[test.Metadata]{
// 				test.Metadata{
// 					SourceID: "1",
// 					TargetID: "91",
// 				},
// 				test.Metadata{
// 					SourceID: "13",
// 					TargetID: "913",
// 				},
// 			},
// 			Store: &test.Store{},
// 		},
// 	}

// 	ctx := context.Background()

// 	mediaIds, err := bridge.GenerateCustomList(ctx, testUsername)
// 	if err != nil {
// 		t.Error("unexpected error:", err)
// 	}

// 	if !reflect.DeepEqual(mediaIds, wantedIds) {
// 		t.Errorf("generated list does not match: got '%v', expected '%v'", mediaIds, wantedIds)
// 	}
// }

// func TestMediaBridge_Close(t *testing.T) {
// 	bridge := usecases.MediaBridge{
// 		Tracker: &test.Tracker{},
// 		Mapper: &adapters.Mapper{
// 			Provider: testProvider,
// 			Store:    &test.Store{},
// 		},
// 	}

// 	err := bridge.Close()
// 	if err != nil {
// 		t.Error("unexpected error:", err)
// 	}
// }

// func TestMediaBridge_Refresh(t *testing.T) {
// 	bridge := usecases.MediaBridge{
// 		Tracker: &test.Tracker{},
// 		Mapper: &adapters.Mapper{
// 			Provider: testProvider,
// 			Store:    &test.Store{},
// 		},
// 	}

// 	ctx := context.Background()
// 	err := bridge.Refresh(ctx)
// 	if err != nil {
// 		t.Error("unexpected error:", err)
// 	}
// }
