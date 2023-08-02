package caches

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
)

type RedisOptions = redis.Options

func NewRedis(options *RedisOptions) (adapters.Cache, error) {
	rdb := redis.NewClient(options)

	err := redisotel.InstrumentTracing(rdb)
	if err != nil {
		return nil, fmt.Errorf("failed to instrument tracing for Redis: %w", err)
	}

	err = redisotel.InstrumentMetrics(rdb)
	if err != nil {
		return nil, fmt.Errorf("failed to instrument metrics for Redis: %w", err)
	}

	return &redisCache{rdb}, nil
}

type redisCache struct {
	*redis.Client
}

func (c *redisCache) GetString(ctx context.Context, key string) (string, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	res, err := c.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", span.Assert(nil)
	}

	if err != nil {
		return "", span.Assert(fmt.Errorf("failed to get string: %w", err))
	}

	return res, span.Assert(nil)
}

func (c *redisCache) SetString(ctx context.Context, key, value string) error {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	return span.Assert(c.Set(ctx, key, value, 0).Err())
}
