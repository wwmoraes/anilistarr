// Package redis implements a Redis-backed driver to fit use-cases needs.
package redis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	telemetry "github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var _ usecases.Cache = (*Redis)(nil)

// Options re-exports [redis.Options] so consumers avoid an extra import
type Options = redis.Options

// Redis provides a Redis-backed cache driver. It augments the vanilla client
// with OpenTelemetry tracing and metrics.
type Redis struct {
	client *redis.Client
}

// New creates the underlying Redis client, then configures tracing and
// metrics on top of it. Returns an error if it fails to instrument the client.
func New(options *Options) (*Redis, error) {
	client := redis.NewClient(options)

	err := client.Ping(context.TODO()).Err()
	if err != nil {
		return nil, errors.Join(usecases.ErrStatusUnavailable, err)
	}

	//nolint:errcheck // will never error
	redisotel.InstrumentTracing(client)

	//nolint:errcheck // will never error
	redisotel.InstrumentMetrics(client)

	return &Redis{client}, nil
}

// Close terminates the underlying Redis client, freeing up resources. It is
// undefined behavior to use this driver after closing it.
func (cache *Redis) Close() error {
	err := cache.client.Close()

	return usecases.ErrorJoinIf(usecases.ErrStatusInternal, err)
}

// GetString queries the value for key. It returns [usecases.ErrStatusNotFound]
// if key is not in the cache.
func (cache *Redis) GetString(ctx context.Context, key string) (string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := cache.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", span.Assert(usecases.ErrStatusNotFound)
	}

	return res, span.Assert(usecases.ErrorJoinIf(usecases.ErrStatusUnknown, err))
}

// SetString stores value for key in the cache. It supports entries with a set
// expiration time by using [adapters.WithTTL] option.
func (cache *Redis) SetString(ctx context.Context, key, value string, options ...usecases.CacheOption) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	params := usecases.NewCacheOptions(options...)

	err := cache.client.Set(ctx, key, value, params.TTL).Err()

	return span.Assert(usecases.ErrorJoinIf(usecases.ErrStatusInternal, err))
}
