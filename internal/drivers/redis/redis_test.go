package redis_test

import (
	"context"
	"flag"
	"io"
	"net"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wwmoraes/anilistarr/internal/drivers/redis"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var verbose bool

type WriterFn func([]byte) (int, error)

func (fn WriterFn) Write(p []byte) (int, error) {
	return fn(p)
}

func stdoutWriter(tb testing.TB) io.Writer {
	tb.Helper()

	return WriterFn(func(p []byte) (int, error) {
		tb.Helper()

		tb.Log(string(p))

		return len(p), nil
	})
}

func stderrWriter(tb testing.TB) io.Writer {
	tb.Helper()

	return WriterFn(func(p []byte) (int, error) {
		tb.Helper()

		tb.Error(string(p))

		return len(p), nil
	})
}

func runValkey(tb testing.TB) string {
	tb.Helper()

	socketPath := filepath.Join(tb.TempDir(), "valkey.sock")
	if verbose {
		tb.Log("runValkey socket path:", socketPath)
	}

	// XXX: Darwin (and some CIs that set GOTMPDIR) have long base temp paths e.g.
	// on Darwin: /var/folders/7x/\w{30}/T/TestName\d{10}/\d{3}/valkey.sock
	// in this case, 64 characters are out of our control; the difference is split
	// between test name and file descriptor name.
	//
	// Solution: set GOTMPDIR to /tmp to enforce that POSIX path, which also works
	// on Darwin. Done through the .env file.
	if len(socketPath) >= 104 {
		tb.Fatal("unix socket path too long for Valkey, must be under 104")
	}

	//nolint:forbidigo,gosec // fine for test purposes
	cmd := exec.CommandContext(
		tb.Context(),
		"valkey-server",
		"--maxmemory", "64mb",
		"--port", "0",
		"--unixsocket", socketPath,
	)

	if verbose {
		//nolint:forbidigo // fine for test purposes
		cmd.Stdout = stdoutWriter(tb)
		//nolint:forbidigo // fine for test purposes
		cmd.Stderr = stderrWriter(tb)
	}

	//nolint:forbidigo // fine for test purposes
	err := cmd.Start()
	if err != nil {
		tb.Fatal("failed to start valkey", err)
	}

	context.AfterFunc(tb.Context(), func() {
		//nolint:forbidigo // fine for test purposes
		err := cmd.Wait()
		if verbose {
			tb.Log(err)
		}
	})

	return socketPath
}

func runDummySocket(tb testing.TB) string {
	tb.Helper()

	sockPath := filepath.Join(tb.TempDir(), "dummy.sock")

	var lc net.ListenConfig

	listener, err := lc.Listen(tb.Context(), "unix", sockPath)
	require.NoError(tb, err)

	context.AfterFunc(tb.Context(), func() {
		listener.Close()
	})

	go func() {
		var conn net.Conn

		var err error

		for {
			conn, err = listener.Accept()
			if err != nil {
				break
			}

			conn.Close()
		}
	}()

	return sockPath
}

func TestRedis(t *testing.T) {
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

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		options              *redis.Options
		assertValue          require.ValueAssertionFunc
		assertError          require.ErrorAssertionFunc
		socketAddressBuilder func(tb testing.TB) string
		name                 string
	}{
		{
			name: "succeeds with a valid cache instance",
			options: &redis.Options{
				DisableIndentity: true,
				Network:          "unix",
			},
			socketAddressBuilder: runValkey,
			assertValue:          require.NotNil,
			assertError:          require.NoError,
		},
		{
			name: "fails with an invalid cache instance",
			options: &redis.Options{
				DisableIndentity: true,
				Network:          "unix",
			},
			socketAddressBuilder: runDummySocket,
			assertValue:          require.Nil,
			assertError:          require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.options.Addr = tt.socketAddressBuilder(t)

			got, err := redis.New(t.Context(), tt.options)
			tt.assertError(t, err)
			tt.assertValue(t, got)
		})
	}
}

func TestMain(m *testing.M) {
	flag.BoolVar(&verbose, "v", false, "verbose logging")
	flag.Parse()

	m.Run()
}
