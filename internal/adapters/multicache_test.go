package adapters_test

import (
	"context"
	"fmt"
	"path"
	"testing"
	"time"

	"github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/drivers/caches"
	"github.com/wwmoraes/anilistarr/internal/drivers/stores"
)

func NewFileCache(tb testing.TB) adapters.Cache {
	tb.Helper()

	filePath := path.Join(tb.TempDir(), tb.Name(), fmt.Sprint(time.Now().UnixNano()))

	store, err := stores.NewBadger(filePath, &stores.BadgerOptions{
		Logger: &caches.BadgerLogr{Logger: gotell.Logr(context.TODO())},
	})
	if err != nil {
		tb.Fatal(err)
	}

	return store
}

func TestMultiCache_Close(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		chain   adapters.MultiCache
		wantErr bool
	}{
		{
			name:    "empty",
			chain:   adapters.MultiCache{},
			wantErr: false,
		},
		{
			name: "single",
			chain: adapters.MultiCache{
				NewFileCache(t),
			},
			wantErr: false,
		},
		{
			name: "multi",
			chain: adapters.MultiCache{
				NewFileCache(t),
				NewFileCache(t),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.chain.Close(); (err != nil) != tt.wantErr {
				t.Errorf("MultiCache.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMultiCache_GetString(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		key string
	}

	tests := []struct {
		name    string
		chain   adapters.MultiCache
		args    args
		want    string
		wantErr bool
	}{
		{
			name:  "empty",
			chain: adapters.MultiCache{},
			args: args{
				ctx: context.TODO(),
				key: "foo",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "single",
			chain: adapters.MultiCache{
				NewFileCache(t),
			},
			args: args{
				ctx: context.TODO(),
				key: "foo",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "multi",
			chain: adapters.MultiCache{
				NewFileCache(t),
				NewFileCache(t),
			},
			args: args{
				ctx: context.TODO(),
				key: "foo",
			},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.chain.GetString(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("MultiCache.GetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MultiCache.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiCache_SetString(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx     context.Context
		key     string
		value   string
		options []adapters.CacheOption
	}

	tests := []struct {
		name    string
		chain   adapters.MultiCache
		args    args
		wantErr bool
	}{
		{
			name:  "empty",
			chain: adapters.MultiCache{},
			args: args{
				ctx:   context.TODO(),
				key:   "foo",
				value: "bar",
				options: []adapters.CacheOption{
					adapters.WithTTL(time.Millisecond),
				},
			},
			wantErr: true,
		},
		{
			name: "single",
			chain: adapters.MultiCache{
				NewFileCache(t),
			},
			args: args{
				ctx:   context.TODO(),
				key:   "foo",
				value: "bar",
				options: []adapters.CacheOption{
					adapters.WithTTL(time.Millisecond),
				},
			},
			wantErr: false,
		},
		{
			name: "multi",
			chain: adapters.MultiCache{
				NewFileCache(t),
				NewFileCache(t),
			},
			args: args{
				ctx:   context.TODO(),
				key:   "foo",
				value: "bar",
				options: []adapters.CacheOption{
					adapters.WithTTL(time.Millisecond),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.chain.SetString(tt.args.ctx, tt.args.key, tt.args.value, tt.args.options...); (err != nil) != tt.wantErr {
				t.Errorf("MultiCache.SetString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
