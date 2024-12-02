// Package bolt implements a BoltDB-backed driver that implements use-cases.
package bolt

import (
	"context"
	"errors"
	"fmt"

	telemetry "github.com/wwmoraes/gotell"
	"go.etcd.io/bbolt"

	"github.com/wwmoraes/anilistarr/internal/usecases"
)

// BucketName is the internal bucket name that the driver uses. Consumers should
// avoid re-using the same bucket name.
const BucketName = "anilistarr"

var _ usecases.Cache = (*Bolt)(nil)

// Options re-rexports the upstream [bbolt.Options] so consumers don't need
// to have an extra import.
type Options = bbolt.Options

// Bolt provides a BoltDB-backed cache driver. It uses a single constant bucket
// name defined by [BucketName].
type Bolt struct {
	db *bbolt.DB
}

// New creates a Bolt-based Cache
func New(path string, options *Options) (*Bolt, error) {
	db, err := bbolt.Open(path, 0o640, options)
	if err != nil {
		return nil, fmt.Errorf("failed to open bolt database: %w", err)
	}

	// initialize DB if RW
	if !db.IsReadOnly() {
		//nolint:errcheck // no error to check
		db.Update(func(tx *bbolt.Tx) error {
			//nolint:errcheck // known bucket name + open and RW transaction
			tx.CreateBucketIfNotExists([]byte(BucketName))

			return nil
		})
	}

	return &Bolt{db}, nil
}

// Close closes the underlying BoltDB handler, finishing up transactions and
// cleaning up resources.
func (cache *Bolt) Close() error {
	return usecases.ErrorJoinIf(
		usecases.ErrStatusUnknown,
		cache.db.Close(),
	)
}

// GetString retrieves a value stored in the underlying BoltDB. It returns
// [usecases.ErrStatusNotFound] if key isn't in the cache.
func (cache *Bolt) GetString(ctx context.Context, key string) (string, error) {
	_, span := telemetry.Start(ctx)
	defer span.End()

	var value string

	err := cache.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		data := bucket.Get([]byte(key))
		if data == nil {
			return usecases.ErrStatusNotFound
		}

		value = string(data)

		return nil
	})
	if err != nil {
		return "", span.Assert(fmt.Errorf("failed to get string: %w", err))
	}

	return value, span.Assert(nil)
}

// SetString stores a key and value in the underlying BoltDB. It returns
func (cache *Bolt) SetString(ctx context.Context, key, value string, _ ...usecases.CacheOption) error {
	_, span := telemetry.Start(ctx)
	defer span.End()

	return span.Assert(cache.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		err := bucket.Put([]byte(key), []byte(value))
		if err != nil {
			return errors.Join(usecases.ErrStatusInvalidArgument, err)
		}

		return nil
	}))
}
