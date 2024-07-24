package adapters_test

import (
	"context"
	"reflect"
	"slices"
	"testing"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/drivers/stores"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestMapper(t *testing.T) {
	t.Parallel()

	store, err := stores.NewBadger("", &stores.BadgerOptions{
		InMemory: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	mapper := adapters.Mapper{
		Provider: newJSONLocalProvider(t),
		Store:    store,
	}

	defer mapper.Close()

	gotIDs, gotErr := mapper.MapIDs(context.Background(), testSourceIDs)
	if gotErr != nil {
		t.Errorf("unexpected error %v", gotErr)
	}

	if !slices.Equal(gotIDs, []string{}) {
		t.Errorf("got %v, want %v", gotIDs, []string{})
	}

	// refresh to get some data into the store
	err = mapper.Refresh(context.Background(), usecases.HTTPGetterAsGetter(newMemoryGetter(t)))
	if err != nil {
		t.Error(err)
	}

	gotIDs, gotErr = mapper.MapIDs(context.Background(), testSourceIDs)
	if gotErr != nil {
		t.Errorf("unexpected error %v", gotErr)
	}

	if !slices.Equal(gotIDs, testTargetIDs) {
		t.Errorf("got %v, want %v", gotIDs, testTargetIDs)
	}
}

func TestMapper_MapIDs(t *testing.T) {
	type fields struct {
		Provider adapters.Provider[adapters.Metadata]
		Store    adapters.Store
	}

	type args struct {
		ctx context.Context
		ids []entities.SourceID
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entities.TargetID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapper := &adapters.Mapper{
				Provider: tt.fields.Provider,
				Store:    tt.fields.Store,
			}

			got, err := mapper.MapIDs(tt.args.ctx, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mapper.MapIDs() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mapper.MapIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}
