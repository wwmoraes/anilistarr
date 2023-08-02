package models

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type MappingList []*Mapping

func (m MappingList) Upsert(ctx context.Context, db DB) error {
	rows := make([]string, len(m))
	for index, entry := range m {
		if entry._deleted {
			return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
		}

		rows[index] = fmt.Sprintf("(%s, %s)", entry.TvdbID, nullableString(entry.AnilistID))
	}

	const baseSqlstr = `INSERT INTO mapping (` +
		`tvdb_id, anilist_id` +
		`) VALUES %s` +
		` ON CONFLICT (tvdb_id) DO ` +
		`UPDATE SET ` +
		`anilist_id = EXCLUDED.anilist_id `

	sqlstr := fmt.Sprintf(baseSqlstr, strings.Join(rows, ","))

	logf(sqlstr)
	if _, err := db.ExecContext(ctx, sqlstr); err != nil {
		return logerror(err)
	}

	for _, entry := range m {
		entry._exists = true
	}

	return nil
}

func MappingByAnilistIDBulk(ctx context.Context, db DB, anilistIDs []sql.NullString) ([]*Mapping, error) {
	ids := make([]string, 0, len(anilistIDs))
	for _, id := range anilistIDs {
		if !id.Valid {
			continue
		}

		ids = append(ids, id.String)
	}

	// query
	const baseSqlstr = `SELECT ` +
		`tvdb_id, anilist_id ` +
		`FROM mapping ` +
		`WHERE anilist_id IN (%s)`
	sqlstr := fmt.Sprintf(baseSqlstr, strings.Join(ids, ","))

	// run
	logf(sqlstr)
	rows, err := db.QueryContext(ctx, sqlstr)
	if err != nil {
		return nil, logerror(err)
	}
	defer rows.Close()

	results := make([]*Mapping, 0, len(ids))
	var m *Mapping
	for rows.Next() {
		m = &Mapping{
			_exists: true,
		}

		err = rows.Scan(&m.TvdbID, &m.AnilistID)
		if err != nil {
			return nil, logerror(err)
		}

		results = append(results, m)
	}

	return results, nil
}

func nullableString(v sql.NullString) string {
	if !v.Valid {
		return "NULL"
	}

	return v.String
}
