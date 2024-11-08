package persistence

import (
	"context"
	"errors"
	"fmt"

	"dario.cat/mergo"
	"github.com/dgraph-io/badger/v4"
	telemetry "github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

type BadgerOptions = badger.Options

type badgerPersistence struct {
	*badger.DB
}

func NewBadger(path string, options *BadgerOptions) (Persistence, error) {
	opt := badger.DefaultOptions(path).WithLoggingLevel(badger.ERROR)

	if options != nil {
		err := mergo.Merge(&opt, *options, mergo.WithOverride)
		if err != nil {
			return nil, fmt.Errorf("failed to merge Badger options: %w", err)
		}
	}

	db, err := badger.Open(opt)
	if err != nil {
		return nil, err
	}

	return &badgerPersistence{db}, nil
}

func (c *badgerPersistence) GetString(ctx context.Context, key string) (string, error) {
	_, span := telemetry.Start(ctx)
	defer span.End()

	var value string

	err := c.View(func(txn *badger.Txn) error {
		data, err := txn.Get([]byte(key))
		if errors.Is(err, badger.ErrKeyNotFound) {
			return usecases.ErrNotFound
		}

		if err != nil {
			return err
		}

		if data.IsDeletedOrExpired() {
			return nil
		}

		return data.Value(func(val []byte) error {
			value = string(val)

			return nil
		})
	})

	return value, span.Assert(err)
}

func (c *badgerPersistence) SetString(ctx context.Context, key, value string, options ...adapters.CacheOption) error {
	_, span := telemetry.Start(ctx)
	defer span.End()

	params, err := adapters.NewCacheParams(options...)
	if err != nil {
		return err
	}

	return span.Assert(c.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry([]byte(key), []byte(value))

		return txn.SetEntry(entry.WithTTL(params.TTL))
	}))
}

func (db *badgerPersistence) GetMedia(ctx context.Context, id string) (*entities.Media, error) {
	_, span := telemetry.Start(ctx)
	defer span.End()

	var media *entities.Media

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil
		}

		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			media = &entities.Media{
				SourceID: id,
				TargetID: string(val),
			}

			return nil
		})
	})

	return media, span.Assert(err)
}

func (db *badgerPersistence) GetMediaBulk(ctx context.Context, ids []string) ([]*entities.Media, error) {
	_, span := telemetry.Start(ctx)
	defer span.End()

	medias := make([]*entities.Media, 0, len(ids))

	err := db.View(func(txn *badger.Txn) error {
		var item *badger.Item

		var err error

		for _, id := range ids {
			if len(id) == 0 {
				continue
			}

			item, err = txn.Get([]byte(id))
			if errors.Is(err, badger.ErrKeyNotFound) {
				continue
			}

			if err != nil {
				return err
			}

			err = item.Value(func(val []byte) error {
				medias = append(medias, &entities.Media{
					SourceID: id,
					TargetID: string(val),
				})

				return nil
			})
			if err != nil {
				return err
			}
		}

		return nil
	})

	return medias, span.Assert(err)
}

func (db *badgerPersistence) PutMedia(ctx context.Context, media *entities.Media) error {
	_, span := telemetry.Start(ctx)
	defer span.End()

	return span.Assert(db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(media.SourceID), []byte(media.TargetID))
	}))
}

func (db *badgerPersistence) PutMediaBulk(ctx context.Context, medias []*entities.Media) error {
	_, span := telemetry.Start(ctx)
	defer span.End()

	return span.Assert(db.Update(func(txn *badger.Txn) error {
		var err error

		for _, media := range medias {
			err = txn.Set([]byte(media.SourceID), []byte(media.TargetID))
			if err != nil {
				return err
			}
		}

		return nil
	}))
}
