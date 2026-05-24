package chaincache_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/anilistarr/internal/adapters/chaincache"
	"github.com/wwmoraes/anilistarr/internal/test"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestChainCache_Close(t *testing.T) {
	t.Parallel()

	cache := test.NewMockCache(t)
	errorCache := test.NewMockCache(t)

	cache.EXPECT().Close().Return(nil).Maybe()
	errorCache.EXPECT().Close().Return(errors.New("foo")).Maybe()

	tests := []struct {
		assertError require.ErrorAssertionFunc
		name        string
		chain       chaincache.ChainCache
	}{
		{
			name:        "empty",
			chain:       chaincache.ChainCache{},
			assertError: require.NoError,
		},
		{
			name: "single",
			chain: chaincache.ChainCache{
				cache,
			},
			assertError: require.NoError,
		},
		{
			name: "multi",
			chain: chaincache.ChainCache{
				cache,
				cache,
			},
			assertError: require.NoError,
		},
		{
			name: "error",
			chain: chaincache.ChainCache{
				cache,
				errorCache,
			},
			assertError: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.chain.Close()
			tt.assertError(t, err)
		})
	}
}

func TestChainCache_GetString(t *testing.T) {
	t.Parallel()

	key := "foo"
	value := "bar"

	missCache := test.NewMockCache(t)
	hitCache := test.NewMockCache(t)
	skippedCache := test.NewMockCache(t)
	errorCache := test.NewMockCache(t)

	missCache.EXPECT().GetString(mock.Anything, key).
		Return("", usecases.ErrStatusNotFound).Maybe()
	hitCache.EXPECT().GetString(mock.Anything, key).
		Return(value, nil).Maybe()
	errorCache.EXPECT().GetString(mock.Anything, key).
		Return("", errors.New("qux")).Maybe()

	tests := []struct {
		name      string
		want      string
		wantError error
		chain     chaincache.ChainCache
	}{
		{
			name:      "no providers",
			chain:     chaincache.ChainCache{},
			want:      "",
			wantError: usecases.ErrStatusNotFound,
		},
		{
			name: "single empty",
			chain: chaincache.ChainCache{
				missCache,
			},
			want:      "",
			wantError: usecases.ErrStatusNotFound,
		},
		{
			name: "multi empty",
			chain: chaincache.ChainCache{
				missCache,
				missCache,
			},
			want:      "",
			wantError: usecases.ErrStatusNotFound,
		},
		{
			name: "single match",
			chain: chaincache.ChainCache{
				hitCache,
				skippedCache,
			},
			want:      "bar",
			wantError: nil,
		},
		{
			name: "multi match second",
			chain: chaincache.ChainCache{
				missCache,
				hitCache,
				skippedCache,
			},
			want:      "bar",
			wantError: nil,
		},
		{
			name: "cache error",
			chain: chaincache.ChainCache{
				missCache,
				errorCache,
				skippedCache,
			},
			want:      "",
			wantError: usecases.ErrStatusUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.chain.GetString(t.Context(), key)
			require.ErrorIs(t, err, tt.wantError)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestChainCache_SetString(t *testing.T) {
	t.Parallel()

	key := "foo"
	value := "bar"
	options := []usecases.CacheOption{
		usecases.WithTTL(time.Millisecond),
	}

	hitCache := test.NewMockCache(t)
	skippedCache := test.NewMockCache(t)
	errorCache := test.NewMockCache(t)

	hitCache.EXPECT().SetString(mock.Anything, key, value, options).
		Return(nil).Maybe()
	errorCache.EXPECT().SetString(mock.Anything, key, value, options).
		Return(errors.New("qux")).Maybe()

	tests := []struct {
		assertError require.ErrorAssertionFunc
		name        string
		chain       chaincache.ChainCache
	}{
		{
			name:        "empty",
			chain:       chaincache.ChainCache{},
			assertError: require.Error,
		},
		{
			name: "single",
			chain: chaincache.ChainCache{
				hitCache,
			},
			assertError: require.NoError,
		},
		{
			name: "multi",
			chain: chaincache.ChainCache{
				hitCache,
				skippedCache,
			},
			assertError: require.NoError,
		},
		{
			name: "error",
			chain: chaincache.ChainCache{
				errorCache,
				skippedCache,
			},
			assertError: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.chain.SetString(
				t.Context(),
				key,
				value,
				options...,
			)
			tt.assertError(t, err)
		})
	}
}
