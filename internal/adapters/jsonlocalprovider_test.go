package adapters_test

import (
	"context"
	"reflect"
	"strconv"
	"testing"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/drivers/stores"
	"github.com/wwmoraes/anilistarr/internal/test"
)

var (
	sampleData = `[
		{"anilist_id": 1, "thetvdb_id": 91},
		{"anilist_id": 2, "thetvdb_id": 92}
	]`
	testProvider adapters.JSONLocalProvider[memoryMetadata] = adapters.JSONLocalProvider[memoryMetadata]{
		Fs: &test.MemoryFS{
			"test.json": []byte(sampleData),
		},
		Name: "test.json",
	}
)

type memoryMetadata struct {
	AnilistID uint64 `json:"anilist_id"`
	TheTvdbID uint64 `json:"thetvdb_id"`
}

func (entry memoryMetadata) GetSourceID() string {
	return strconv.FormatUint(entry.AnilistID, 10)
}

func (entry memoryMetadata) GetTargetID() string {
	return strconv.FormatUint(entry.TheTvdbID, 10)
}

func TestJSONLocalProvider(t *testing.T) {
	store, err := stores.NewBadger("", &stores.BadgerOptions{
		InMemory: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	mapper := &adapters.Mapper{
		Provider: testProvider,
		Store:    store,
	}

	ctx := context.Background()
	err = mapper.Refresh(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}

	wantedTargetIds := []string{"91", "92"}
	targetIds, err := mapper.MapIDs(ctx, []string{"1", "2", "3"})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(targetIds, wantedTargetIds) {
		t.Fatalf("mapped IDs don't match: got %v, wanted %v", targetIds, wantedTargetIds)
	}
}
