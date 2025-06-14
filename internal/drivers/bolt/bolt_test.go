package bolt_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/bbolt"
	bolterrors "go.etcd.io/bbolt/errors"

	"github.com/wwmoraes/anilistarr/internal/drivers/bolt"
	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		options *bolt.Options
		path    string
	}

	tests := []struct {
		args        args
		assertValue assert.ValueAssertionFunc
		assertError require.ErrorAssertionFunc
		name        string
	}{
		{
			name: "success",
			args: args{
				path:    path.Join(t.TempDir(), "success"),
				options: bbolt.DefaultOptions,
			},
			assertValue: assert.NotNil,
			assertError: require.NoError,
		},
		{
			name: "open error",
			args: args{
				path:    path.Join(t.TempDir(), "foo", "bar"),
				options: bbolt.DefaultOptions,
			},
			assertValue: assert.Nil,
			assertError: require.Error,
		},
		{
			name: "update error",
			args: args{
				path: path.Join(t.TempDir(), "update"),
				options: &bbolt.Options{
					Timeout:      0,
					NoGrowSync:   false,
					FreelistType: bbolt.FreelistArrayType,
					ReadOnly:     true,
				},
			},
			assertValue: assert.Nil,
			assertError: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := bolt.New(tt.args.path, tt.args.options)
			tt.assertError(t, err)

			if got != nil {
				testdata.Close(t, got)
			}

			tt.assertValue(t, got)
		})
	}
}

func TestBolt_GetString(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx     context.Context
		options *bolt.Options
		key     string
	}

	tests := []struct {
		args      args
		assertion require.ErrorAssertionFunc
		name      string
		want      string
	}{
		{
			name: "hit",
			args: args{
				ctx:     context.TODO(),
				key:     "foo",
				options: bbolt.DefaultOptions,
			},
			want:      "bar",
			assertion: require.NoError,
		},
		{
			name: "miss",
			args: args{
				ctx:     context.TODO(),
				key:     "bar",
				options: bbolt.DefaultOptions,
			},
			want:      "",
			assertion: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dirhash := sha256.New().Sum([]byte(tt.name))
			dirname := hex.EncodeToString(dirhash)
			cachePath := path.Join(t.TempDir(), dirname)

			cache, err := bolt.New(cachePath, tt.args.options)
			require.NoError(t, err)
			defer testdata.Close(t, cache)

			err = cache.SetString(tt.args.ctx, "foo", "bar")
			require.NoError(t, err)

			got, err := cache.GetString(tt.args.ctx, tt.args.key)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBolt_SetString(t *testing.T) {
	t.Parallel()

	type args struct {
		options *bolt.Options
		ctx     context.Context
		key     string
		value   string
	}

	tests := []struct {
		wantError error
		name      string
		args      args
	}{
		{
			name: "success",
			args: args{
				ctx:     context.TODO(),
				key:     "foo",
				value:   "bar",
				options: nil,
			},
			wantError: nil,
		},
		{
			name: "empty key",
			args: args{
				ctx:     context.TODO(),
				key:     "",
				value:   "bar",
				options: nil,
			},
			wantError: usecases.ErrStatusInvalidArgument,
		},
		{
			name: "readonly",
			args: args{
				ctx:   context.TODO(),
				key:   "foo",
				value: "bar",
				options: &bbolt.Options{
					Timeout:      0,
					NoGrowSync:   false,
					FreelistType: bbolt.FreelistArrayType,
					ReadOnly:     true,
				},
			},
			wantError: bolterrors.ErrDatabaseReadOnly,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dirhash := sha256.New().Sum([]byte(tt.name))
			dirname := hex.EncodeToString(dirhash)
			cachePath := path.Join(t.TempDir(), dirname)

			// create DB first
			cache, err := bolt.New(cachePath, nil)
			require.NoError(t, err)

			// close it
			err = cache.Close()
			require.NoError(t, err)

			// start test
			cache, err = bolt.New(cachePath, tt.args.options)
			require.NoError(t, err)

			defer testdata.Close(t, cache)

			err = cache.SetString(
				tt.args.ctx,
				tt.args.key,
				tt.args.value,
			)
			require.ErrorIs(t, err, tt.wantError)
		})
	}
}
