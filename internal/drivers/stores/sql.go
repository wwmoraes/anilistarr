package stores

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/wwmoraes/anilistarr/internal/drivers/stores/models"
	"github.com/wwmoraes/anilistarr/internal/entities"
	"github.com/wwmoraes/anilistarr/internal/telemetry"

	_ "modernc.org/sqlite"
)

type Sql struct {
	db *sql.DB
}

func NewSQL(driverName, dataSourceName string) (*Sql, error) {
	db, err := telemetry.OpenSQL(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQL database: %w", err)
	}

	return &Sql{
		db: db,
	}, nil
}

func (s *Sql) PutMedia(ctx context.Context, media *entities.Media) error {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	record := models.MappingFromMedia(media)

	return span.Assert(record.Upsert(ctx, s.db))
}

func (s *Sql) PutMediaBulk(ctx context.Context, medias []*entities.Media) error {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	records := make(models.MappingList, len(medias))
	for index, media := range medias {
		records[index] = &models.Mapping{
			TvdbID: media.TvdbID,
			AnilistID: sql.NullString{
				String: media.AnilistID,
				Valid:  len(media.AnilistID) > 0,
			},
		}
	}

	return span.Assert(records.Upsert(ctx, s.db))
}

func (s *Sql) MappingByAnilistID(ctx context.Context, anilistId string) (*entities.Media, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	record, err := models.MappingByAnilistID(ctx, s.db, sql.NullString{
		String: anilistId,
		Valid:  len(anilistId) > 0,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil, span.Assert(nil)
	} else if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to get mapping by anilist ID: %w", err))
	}

	return record.ToMedia(), span.Assert(nil)
}

func (s *Sql) MappingByAnilistIDBulk(ctx context.Context, anilistIds []string) ([]*entities.Media, error) {
	ctx, span := telemetry.StartFunction(ctx)
	defer span.End()

	ids := make([]sql.NullString, len(anilistIds))
	for index, id := range anilistIds {
		ids[index] = sql.NullString{
			String: id,
			Valid:  len(id) > 0,
		}
	}

	records, err := models.MappingByAnilistIDBulk(ctx, s.db, ids)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, span.Assert(nil)
	} else if err != nil {
		return nil, span.Assert(fmt.Errorf("failed to get mapping by anilist ID: %w", err))
	}

	results := make([]*entities.Media, len(records))
	for index, entry := range records {
		results[index] = entry.ToMedia()
	}

	return results, span.Assert(nil)
}

func (s *Sql) Close() error {
	return s.db.Close()
}
