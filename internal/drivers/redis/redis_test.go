// Package redis implements a Redis-backed driver to fit use-cases needs.
package redis_test

import (
	"context"
	"net"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/anilistarr/internal/drivers/redis"
	"github.com/wwmoraes/anilistarr/internal/testdata"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func TestRedis(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	options := redis.Options{
		DisableIndentity: true,
	}

	if testing.Short() {
		listener, err := net.Listen("unix", filepath.Join(t.TempDir(), "socket"))
		require.NoError(t, err)

		defer listener.Close()

		server := testdata.NewCommandServer(
			listener,
			testdata.Command{
				Request:  testdata.RESP("*2", "$5", "hello", "$1", "3"),
				Response: testdata.RESP("%1", "+server", "+test"),
			},
			testdata.Command{
				Request:  testdata.RESP("*1", "$4", "ping"),
				Response: testdata.RESP("+PONG"),
			},
			testdata.Command{
				Request:  testdata.RESP("*2", "$3", "get", "$3", "foo"),
				Response: testdata.RESP("_"),
			},
			testdata.Command{
				Request:  testdata.RESP("*3", "$3", "set", "$3", "foo", "$3", "bar"),
				Response: testdata.RESP("+OK"),
			},
			testdata.Command{
				Request:  testdata.RESP("*2", "$3", "get", "$3", "foo"),
				Response: testdata.RESP("+bar"),
			},
		)

		server.Start(ctx)

		options.Addr = listener.Addr().String()
		options.Network = listener.Addr().Network()
	} else {
		t.Skip("TODO implement external Redis server test")
	}

	key, value := "foo", "bar"

	cache, err := redis.New(&options)
	require.NoError(t, err)

	got, err := cache.GetString(ctx, key)
	require.ErrorIs(t, err, usecases.ErrStatusNotFound)

	assert.Empty(t, got)

	err = cache.SetString(ctx, key, value)
	require.NoError(t, err)

	got, err = cache.GetString(ctx, key)
	require.NoError(t, err)

	assert.Equal(t, value, got)

	err = cache.Close()
	require.NoError(t, err)
}

func TestNew(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	options := redis.Options{
		DisableIndentity: true,
	}

	if testing.Short() {
		listener, err := net.Listen("unix", filepath.Join(t.TempDir(), "socket"))
		require.NoError(t, err)
		defer listener.Close()

		server := testdata.NewCommandServer(listener)
		server.Start(ctx)

		options.Addr = listener.Addr().String()
		options.Network = listener.Addr().Network()
	} else {
		t.Skip("TODO implement external Redis server test")
	}

	cache, err := redis.New(&options)
	require.ErrorIs(t, err, usecases.ErrStatusUnavailable)

	assert.Nil(t, cache)
}
