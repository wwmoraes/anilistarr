// Package sqlite provides a SQLite-backed persistence store driver that
// implements both [adapters.Cache] and [usecases.Store].
package sqlite

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	telemetry "github.com/wwmoraes/gotell"

	"github.com/wwmoraes/anilistarr/internal/drivers/sqlite/model"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/usecases"
)

var (
	//go:embed schema.sql
	schema string

	_ usecases.Cache = (*SQLite)(nil)
	_ usecases.Store = (*SQLite)(nil)
)

// SQLite provides a SQLite-backed cache and store driver.
type SQLite struct {
	handler *sql.DB
	queries *model.Queries
}

// New creates a SQLite database handler and tests connection to it. It does NOT
// import a driver; the caller is responsible for importing one.
func New(dataSourceName string) (*SQLite, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// database handlers are lazy and only connect on request. Thus we need to
	// force a connection here to ensure the data source is valid and prevent a
	// misplaced error elsewhere on code
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to execute schema queries: %w", err)
	}

	return &SQLite{
		handler: db,
		queries: model.New(db),
	}, nil
}

// Close terminates the connection to the underlying database.
func (db *SQLite) Close() error {
	return usecases.ErrorJoinIf(
		usecases.ErrStatusInternal,
		db.handler.Close(),
	)
}

// GetString retrieves value for key if it is set. Returns
// [usecases.ErrStatusNotFound] otherwise.
func (db *SQLite) GetString(ctx context.Context, key string) (string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	value, err := db.queries.GetCacheString(ctx, key)
	if err != nil {
		return "", span.Assert(errors.Join(usecases.ErrStatusNotFound, err))
	}

	return value, span.Assert(nil)
}

// SetString stores value for key. It overrides any previously existing value.
func (db *SQLite) SetString(
	ctx context.Context,
	key, value string,
	_ ...usecases.CacheOption,
) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := db.queries.PutCacheString(ctx, model.PutCacheStringParams{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return span.Assert(errors.Join(usecases.ErrStatusUnknown, err))
	}

	return span.Assert(nil)
}

// GetMedia retrieves a media entry from the cache.
func (db *SQLite) GetMedia(ctx context.Context, id string) (*entities.Media, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := db.queries.GetMedia(ctx, id)
	if err != nil {
		return nil, span.Assert(errors.Join(usecases.ErrStatusNotFound, err))
	}

	return &entities.Media{
		SourceID: res.SourceID,
		TargetID: res.TargetID,
	}, nil
}

// GetMediaBulk retrieves a set of media entries from the cache. It returns a
// slice with only matched entries. This means results have a length between
// zero (no entries found) and len(ids).
func (db *SQLite) GetMediaBulk(ctx context.Context, ids []string) ([]*entities.Media, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := db.queries.GetMediaBulk(ctx, ids)
	if err != nil {
		return nil, errors.Join(usecases.ErrStatusUnknown, err)
	}

	medias := make([]*entities.Media, 0, len(res))

	for _, entry := range res {
		medias = append(medias, &entities.Media{
			SourceID: entry.SourceID,
			TargetID: entry.TargetID,
		})
	}

	return medias, nil
}

// PutMedia stores a media in the cache. It uses the media source ID as key in
// the cache. It errors if the source ID is empty.
func (db *SQLite) PutMedia(ctx context.Context, media *entities.Media) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	if !media.Valid() {
		return usecases.ErrStatusInvalidArgument
	}

	err := db.queries.PutMedia(ctx, model.PutMediaParams{
		SourceID: media.SourceID,
		TargetID: media.TargetID,
	})
	if err != nil {
		return errors.Join(usecases.ErrStatusFailedPrecondition, err)
	}

	return nil
}

// PutMediaBulk stores multiple media entries in the cache. This happens within
// a transaction where the same validations as [SQLite.PutMedia] take place.
// An error means no changes were made to the data.
func (db *SQLite) PutMediaBulk(ctx context.Context, medias []*entities.Media) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	if len(medias) == 0 {
		return nil
	}

	// SQLC does not support batch queries for SQLite, so we have to loop it...
	tx, err := db.handler.BeginTx(ctx, nil)
	if err != nil {
		return errors.Join(usecases.ErrStatusFailedPrecondition, err)
	}
	defer tx.Rollback()

	qtx := db.queries.WithTx(tx)

	for _, media := range medias {
		if !media.Valid() {
			return usecases.ErrStatusInvalidArgument
		}

		err = qtx.PutMedia(ctx, model.PutMediaParams{
			SourceID: media.SourceID,
			TargetID: media.TargetID,
		})
		if err != nil {
			return errors.Join(usecases.ErrStatusAborted, err)
		}
	}

	return span.Assert(usecases.ErrorJoinIf(
		usecases.ErrStatusFailedPrecondition,
		tx.Commit(),
	))
}
