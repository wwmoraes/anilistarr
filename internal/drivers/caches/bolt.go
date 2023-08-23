package caches

import (
	"context"
	"fmt"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
	"go.etcd.io/bbolt"
)

const bucketName = "anilistarr"

type BoltOptions = bbolt.Options

type boltCache struct {
	*bbolt.DB
}

func NewBolt(path string, options *BoltOptions) (adapters.Cache, error) {
	db, err := bbolt.Open(path, 0640, options)
	if err != nil {
		return nil, fmt.Errorf("failed to open bolt database: %w", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return fmt.Errorf("failed to create/get bucket: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize bolt: %w", err)
	}

	return &boltCache{db}, nil
}

func (c *boltCache) GetString(ctx context.Context, key string) (string, error) {
	_, span := telemetry.StartFunction(ctx)
	defer span.End()

	var value string

	err := c.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", bucketName)
		}

		data := bucket.Get([]byte(key))
		value = string(data)

		return nil
	})
	if err != nil {
		return "", span.Assert(fmt.Errorf("failed to get string: %w", err))
	}

	return value, span.Assert(nil)
}

func (c *boltCache) SetString(ctx context.Context, key, value string) error {
	_, span := telemetry.StartFunction(ctx)
	defer span.End()

	return span.Assert(c.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		if bucket == nil {
			return fmt.Errorf("bucket %s does not exist", bucketName)
		}

		return bucket.Put([]byte(key), []byte(value))
	}))
}
