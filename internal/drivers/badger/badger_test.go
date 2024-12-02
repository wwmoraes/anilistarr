package badger_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/anilistarr/internal/drivers/badger"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/usecases"
	"github.com/wwmoraes/anilistarr/pkg/with"
)

func TestBadger(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	mediaA := entities.Media{
		SourceID: "foo",
		TargetID: "bar",
	}
	mediaB := entities.Media{
		SourceID: "baz",
		TargetID: "qux",
	}
	mediaC := entities.Media{
		SourceID: "quux",
		TargetID: "corge",
	}
	bulkMedia := []*entities.Media{
		&mediaB,
		&mediaC,
	}
	bulkIDs := []string{
		mediaB.SourceID,
		mediaC.SourceID,
	}
	cacheKey := "grault"
	cacheValue := "garply"

	client, err := badger.New(
		filepath.Join(t.TempDir(), "badger"),
		badger.WithInMemory(true),
		badger.WithLogger(&badger.Logr{
			Logger: logr.Discard(),
		}),
	)
	require.NoError(t, err)

	// get non-existing media
	gotMedia, err := client.GetMedia(ctx, mediaA.SourceID)
	require.ErrorIs(t, err, usecases.ErrStatusNotFound)

	assert.Nil(t, gotMedia)

	// put media
	err = client.PutMedia(ctx, &mediaA)
	require.NoError(t, err)

	// get existing media
	gotMedia, err = client.GetMedia(ctx, mediaA.SourceID)
	require.NoError(t, err)

	assert.Equal(t, &mediaA, gotMedia)

	// get non-existing bulk media
	gotMedias, err := client.GetMediaBulk(ctx, bulkIDs)
	require.ErrorIs(t, err, usecases.ErrStatusNotFound)

	assert.Empty(t, gotMedias)

	// put bulk media
	err = client.PutMediaBulk(ctx, bulkMedia)
	require.NoError(t, err)

	// get existing bulk media
	gotMedias, err = client.GetMediaBulk(ctx, bulkIDs)
	require.NoError(t, err)

	assert.Equal(t, bulkMedia, gotMedias)

	// get non-existing cache string
	gotString, err := client.GetString(ctx, cacheKey)
	require.ErrorIs(t, err, usecases.ErrStatusNotFound)

	assert.Empty(t, gotString)

	// set cache string
	err = client.SetString(ctx, cacheKey, cacheValue)
	require.NoError(t, err)

	// get existing cache string
	gotString, err = client.GetString(ctx, cacheKey)
	require.NoError(t, err)

	assert.Equal(t, cacheValue, gotString)

	err = client.Close()
	require.NoError(t, err)
}

func TestNew_error(t *testing.T) {
	t.Parallel()

	invalidOption := with.Functor[badger.Options](func(options *badger.Options) {
		// BadgerDB needs at least 2 compactors :)
		options.NumCompactors = 1
	})

	client, err := badger.New(
		t.TempDir(),
		badger.WithInMemory(true),
		invalidOption,
	)
	require.Error(t, err)

	assert.Nil(t, client)
}

func TestBadger_GetString_error(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	invalidKey := ""
	deletedKey := "foo"
	value := "bar"

	client, err := badger.New(
		t.TempDir(),
		badger.WithInMemory(true),
	)
	require.NoError(t, err)

	// try to get invalid key
	got, err := client.GetString(ctx, invalidKey)
	require.ErrorIs(t, err, usecases.ErrStatusInvalidArgument)

	assert.Empty(t, got)

	// set an already dead value
	err = client.SetString(ctx, deletedKey, value, usecases.WithTTL(time.Nanosecond))
	require.NoError(t, err)

	<-time.After(time.Nanosecond)

	// get deleted value
	got, err = client.GetString(ctx, deletedKey)
	require.ErrorIs(t, err, usecases.ErrStatusNotFound)

	assert.Empty(t, got)
}

func TestBadger_PutMedia_error(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	media := entities.Media{
		SourceID: "1",
		TargetID: "",
	}

	client, err := badger.New(
		t.TempDir(),
		badger.WithInMemory(true),
	)
	require.NoError(t, err)

	err = client.PutMedia(ctx, &media)
	require.ErrorIs(t, err, usecases.ErrStatusInvalidArgument)
}

func TestBadger_PutMediaBulk_error(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx    context.Context
		medias []*entities.Media
	}

	tests := []struct {
		wantError error
		name      string
		args      args
	}{
		{
			name: "invalid argument",
			args: args{
				ctx: context.TODO(),
				medias: []*entities.Media{
					{
						SourceID: "1",
						TargetID: "",
					},
				},
			},
			wantError: usecases.ErrStatusInvalidArgument,
		},
		{
			name: "invalid key",
			args: args{
				ctx: context.TODO(),
				medias: []*entities.Media{
					{
						SourceID: "!badger!1",
						TargetID: "91",
					},
				},
			},
			wantError: usecases.ErrStatusInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := badger.New(
				t.TempDir(),
				badger.WithInMemory(true),
			)
			require.NoError(t, err)

			err = client.PutMediaBulk(tt.args.ctx, tt.args.medias)
			require.ErrorIs(t, err, tt.wantError)
		})
	}
}

func TestBadger_GetMedia(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  string
	}

	tests := []struct {
		wantError error
		args      args
		name      string
	}{
		{
			name: "empty ID",
			args: args{
				ctx: context.TODO(),
				id:  "",
			},
			wantError: usecases.ErrStatusInvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, err := badger.New(
				t.TempDir(),
				badger.WithInMemory(true),
			)
			require.NoError(t, err)

			got, err := client.GetMedia(tt.args.ctx, tt.args.id)
			require.ErrorIs(t, err, tt.wantError)

			assert.Nil(t, got)
		})
	}
}
