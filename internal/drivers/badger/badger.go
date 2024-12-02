// Package badger provides a BadgerDB-backed driver that implements use-cases.
package badger

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	telemetry "github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/usecases"
	"github.com/wwmoraes/anilistarr/pkg/with"
)

var (
	_ usecases.Cache = (*Badger)(nil)
	_ usecases.Store = (*Badger)(nil)
)

// Options re-exports the upstream [badger.Options] so consumers don't
// need an extra import.
type Options = badger.Options

// Badger provides a BadgerDB-backed driver that implements [adapters.Cache]
// and [usecases.Store], making it easier to store and retrieve typed entries.
//
// Its store part uses [entities.Media.SourceID] as key. For its cache part it
// makes no assumptions about keys, using whatever the caller passes as the key
// parameter. Thus a single instance may serve as both cache and store as long
// as the caller prevents key conflicts.
type Badger struct {
	db *badger.DB
}

// WithLogger provides a custom logger.
func WithLogger(logger badger.Logger) with.Functor[Options] {
	return with.Functor[Options](func(options *Options) {
		options.WithLogger(logger)
	})
}

// WithInMemory toggles in-memory storage.
func WithInMemory(enable bool) with.Functor[Options] {
	return with.Functor[Options](func(options *Options) {
		options.WithInMemory(enable)
	})
}

// New opens an underlying BadgerDB handler. It uses as options the
// default, recommended ones + error-level logging + given options. The caller
// options have the highest precedence.
func New(path string, options ...with.Option[Options]) (*Badger, error) {
	opts := with.Apply(
		badger.DefaultOptions(path).WithLoggingLevel(badger.ERROR),
		options...,
	)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", convertError(err), "failed to open badger")
	}

	return &Badger{db}, nil
}

// Close closes the underlying BadgerDB handler. This must be done at least once
// to ensure it flushes all updates to stable storage.
func (client *Badger) Close() error {
	//nolint:wrapcheck // passthrough
	return client.db.Close()
}

// GetString retrieves value for key if it is set. Returns
// [usecases.ErrStatusNotFound] otherwise.
//
// May rarely return [usecases.ErrStatusInternal] if the underlying badgerDB has
// any problems.
func (client *Badger) GetString(ctx context.Context, key string) (string, error) {
	_, span := telemetry.Start(ctx)
	defer span.End()

	var value string

	err := client.db.View(func(txn *badger.Txn) error {
		data, err := txn.Get([]byte(key))
		if err := convertError(err); err != nil {
			return err
		}

		value = itemValueAsString(data)

		return nil
	})

	return value, span.Assert(err)
}

// SetString stores value for key. It overrides any previously stored value.
func (client *Badger) SetString(ctx context.Context, key, value string, options ...usecases.CacheOption) error {
	_, span := telemetry.Start(ctx)
	defer span.End()

	params := usecases.NewCacheOptions(options...)

	return span.Assert(client.db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(key), []byte(value))

		if params.TTL > 0 {
			entry = entry.WithTTL(params.TTL)
		}

		return txn.SetEntry(entry)
	}))
}

// GetMedia retrieves a media entry from the cache.
func (client *Badger) GetMedia(ctx context.Context, id string) (*entities.Media, error) {
	_, span := telemetry.Start(ctx)
	defer span.End()

	var media entities.Media

	err := client.db.View(mediaGetter(id, &media))
	if err = convertError(err); err != nil {
		return nil, span.Assert(err)
	}

	return &media, span.Assert(nil)
}

// GetMediaBulk retrieves a set of media entries from the cache. It returns a
// slice with only matched entries. This means results have a length between
// zero (no entries found) and len(ids).
func (client *Badger) GetMediaBulk(ctx context.Context, ids []string) ([]*entities.Media, error) {
	_, span := telemetry.Start(ctx)
	defer span.End()

	medias := make([]*entities.Media, 0, len(ids))

	err := client.db.View(func(txn *badger.Txn) error {
		var err error

		for _, id := range ids {
			var media entities.Media

			err = mediaGetter(id, &media)(txn)
			if err != nil {
				return err
			}

			medias = append(medias, &media)
		}

		return nil
	})

	return medias, span.Assert(err)
}

// PutMedia stores a media in the cache. It uses the media source ID as key in
// the cache. It errors if the source ID is empty.
func (client *Badger) PutMedia(ctx context.Context, media *entities.Media) error {
	_, span := telemetry.Start(ctx)
	defer span.End()

	if !media.Valid() {
		return usecases.ErrStatusInvalidArgument
	}

	return span.Assert(client.db.Update(func(txn *badger.Txn) error {
		// TODO join with usecases errors
		return txn.Set([]byte(media.SourceID), []byte(media.TargetID))
	}))
}

// PutMediaBulk stores multiple media entries in the cache. This happens within
// a transaction where the same validations as [Badger.PutMedia] take place.
// An error means no changes were made to the data.
func (client *Badger) PutMediaBulk(ctx context.Context, medias []*entities.Media) error {
	_, span := telemetry.Start(ctx)
	defer span.End()

	return span.Assert(client.db.Update(func(txn *badger.Txn) error {
		var err error

		for _, media := range medias {
			if !media.Valid() {
				return usecases.ErrStatusInvalidArgument
			}

			err = txn.Set([]byte(media.SourceID), []byte(media.TargetID))
			if err := convertError(err); err != nil {
				return err
			}
		}

		return nil
	}))
}

func mediaGetter(id string, media *entities.Media) func(txn *badger.Txn) error {
	return func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err := convertError(err); err != nil {
			return err
		}

		media.SourceID = id
		media.TargetID = itemValueAsString(item)

		return nil
	}
}

func itemValueAsString(item *badger.Item) string {
	var value string

	//nolint:errcheck // no possible error
	_ = item.Value(func(val []byte) error {
		value = string(val)

		return nil
	})

	return value
}

func convertError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, badger.ErrKeyNotFound) {
		return errors.Join(usecases.ErrStatusNotFound, err)
	}

	if usecases.ErrorIn(err, badger.ErrEmptyKey, badger.ErrBannedKey, badger.ErrInvalidKey) {
		return errors.Join(usecases.ErrStatusInvalidArgument, err)
	}

	return errors.Join(usecases.ErrStatusInternal, err)
}
