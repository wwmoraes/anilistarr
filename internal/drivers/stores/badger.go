package stores

import (
	"context"
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/wwmoraes/anilistarr/internal/adapters"
	"github.com/wwmoraes/anilistarr/internal/drivers/caches"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/telemetry"
)

type badgerStore struct {
	*badger.DB
}

// TODO add constructor options support
func NewBadger(path string) (adapters.Store, error) {
	options := badger.DefaultOptions(path)

	// TODO move badger options to a generic package
	options = caches.WithLogger(telemetry.DefaultLogger()).Apply(options)

	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}

	return &badgerStore{db}, nil
}

func (db *badgerStore) GetMedia(ctx context.Context, id string) (*entities.Media, error) {
	_, span := telemetry.StartFunction(ctx)
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

func (db *badgerStore) GetMediaBulk(ctx context.Context, ids []string) ([]*entities.Media, error) {
	_, span := telemetry.StartFunction(ctx)
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

func (db *badgerStore) PutMedia(ctx context.Context, media *entities.Media) error {
	_, span := telemetry.StartFunction(ctx)
	defer span.End()

	return span.Assert(db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(media.SourceID), []byte(media.TargetID))
	}))
}

func (db *badgerStore) PutMediaBulk(ctx context.Context, medias []*entities.Media) error {
	_, span := telemetry.StartFunction(ctx)
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
