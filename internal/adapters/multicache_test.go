package adapters_test

import (
	"context"
	"testing"
	"time"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/drivers/memory"
)

func NewMemCache(tb testing.TB, mutations ...func(adapters.Cache) error) memory.Memory {
	tb.Helper()

	mem := memory.New()

	var err error
	for _, mutation := range mutations {
		err = mutation(mem)
		if err != nil {
			tb.Fatal(err)
		}
	}

	return mem
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
				NewMemCache(t),
			},
			wantErr: false,
		},
		{
			name: "multi",
			chain: adapters.MultiCache{
				NewMemCache(t),
				NewMemCache(t),
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
			name:  "no providers",
			chain: adapters.MultiCache{},
			args: args{
				ctx: context.TODO(),
				key: "foo",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "single empty",
			chain: adapters.MultiCache{
				NewMemCache(t),
			},
			args: args{
				ctx: context.TODO(),
				key: "foo",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "multi empty",
			chain: adapters.MultiCache{
				NewMemCache(t),
				NewMemCache(t),
			},
			args: args{
				ctx: context.TODO(),
				key: "foo",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "single match",
			chain: adapters.MultiCache{
				NewMemCache(t, func(c adapters.Cache) error {
					return c.SetString(context.TODO(), "foo", "bar")
				}),
			},
			args: args{
				ctx: context.TODO(),
				key: "foo",
			},
			want:    "bar",
			wantErr: false,
		},
		{
			name: "multi match second",
			chain: adapters.MultiCache{
				NewMemCache(t),
				NewMemCache(t, func(c adapters.Cache) error {
					return c.SetString(context.TODO(), "foo", "bar")
				}),
			},
			args: args{
				ctx: context.TODO(),
				key: "foo",
			},
			want:    "bar",
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
				NewMemCache(t),
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
				NewMemCache(t),
				NewMemCache(t),
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
