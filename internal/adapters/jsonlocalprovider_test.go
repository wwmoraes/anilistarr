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

//nolint:tagliatelle // JSON tags must match the upstream naming convention
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
	t.Parallel()

	store, err := stores.NewBadger("", &stores.BadgerOptions{
		InMemory: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	mapper := &adapters.Mapper{
		Provider: newJSONLocalProvider(t),
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

func TestJSONLocalProvider_error(t *testing.T) {
	t.Parallel()

	provider := adapters.JSONLocalProvider[memoryMetadata]{
		Fs:   test.MemoryFS{},
		Name: "test.json",
	}

	gotMetadata, err := provider.Fetch(context.TODO(), nil)
	if gotMetadata != nil {
		t.Errorf("JSONLocalProvider.Fetch() = %v, want %v", gotMetadata, nil)
	}

	if err == nil {
		t.Errorf("JSONLocalProvider.Fetch() error = %v, wantErr %v", err, true)
	}
}
