package redis_test

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/anilistarr/internal/drivers/redis"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

func runValkey(tb testing.TB) string {
	tb.Helper()

	socketPath := filepath.Join(tb.TempDir(), "valkey.sock")
	//nolint:forbidigo,gosec // fine for test purposes
	cmd := exec.CommandContext(
		tb.Context(),
		"valkey-server",
		"--maxmemory", "64mb",
		"--port", "0",
		"--unixsocket", socketPath,
	)

	//nolint:forbidigo // fine for test purposes
	err := cmd.Start()
	if err != nil {
		tb.Fatal("failed to start valkey", err)
	}

	return socketPath
}

func TestRedis_works(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	options := redis.Options{
		Addr:             runValkey(t),
		DisableIndentity: true,
		Network:          "unix",
	}

	key, value := "foo", "bar"

	cache, err := redis.New(ctx, &options)
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

func TestNew_createInstanceSuccessfully(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	options := redis.Options{
		Addr:             runValkey(t),
		DisableIndentity: true,
		Network:          "unix",
	}

	cache, err := redis.New(ctx, &options)
	require.ErrorIs(t, err, usecases.ErrStatusUnavailable)

	assert.Nil(t, cache)
}
