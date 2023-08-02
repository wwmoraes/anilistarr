package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"
	"database/sql"
)

// Mapping represents a row from 'mapping'.
type Mapping struct {
	TvdbID    string         `json:"tvdb_id"`    // tvdb_id
	AnilistID sql.NullString `json:"anilist_id"` // anilist_id
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the Mapping exists in the database.
func (m *Mapping) Exists() bool {
	return m._exists
}

// Deleted returns true when the Mapping has been marked for deletion from
// the database.
func (m *Mapping) Deleted() bool {
	return m._deleted
}

// Insert inserts the Mapping to the database.
func (m *Mapping) Insert(ctx context.Context, db DB) error {
	switch {
	case m._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case m._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (manual)
	const sqlstr = `INSERT INTO mapping (` +
		`tvdb_id, anilist_id` +
		`) VALUES (` +
		`$1, $2` +
		`)`
	// run
	logf(sqlstr, m.TvdbID, m.AnilistID)
	if _, err := db.ExecContext(ctx, sqlstr, m.TvdbID, m.AnilistID); err != nil {
		return logerror(err)
	}
	// set exists
	m._exists = true
	return nil
}

// Update updates a Mapping in the database.
func (m *Mapping) Update(ctx context.Context, db DB) error {
	switch {
	case !m._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case m._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE mapping SET ` +
		`anilist_id = $1 ` +
		`WHERE tvdb_id = $2`
	// run
	logf(sqlstr, m.AnilistID, m.TvdbID)
	if _, err := db.ExecContext(ctx, sqlstr, m.AnilistID, m.TvdbID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the Mapping to the database.
func (m *Mapping) Save(ctx context.Context, db DB) error {
	if m.Exists() {
		return m.Update(ctx, db)
	}
	return m.Insert(ctx, db)
}

// Upsert performs an upsert for Mapping.
func (m *Mapping) Upsert(ctx context.Context, db DB) error {
	switch {
	case m._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO mapping (` +
		`tvdb_id, anilist_id` +
		`) VALUES (` +
		`$1, $2` +
		`)` +
		` ON CONFLICT (tvdb_id) DO ` +
		`UPDATE SET ` +
		`anilist_id = EXCLUDED.anilist_id `
	// run
	logf(sqlstr, m.TvdbID, m.AnilistID)
	if _, err := db.ExecContext(ctx, sqlstr, m.TvdbID, m.AnilistID); err != nil {
		return logerror(err)
	}
	// set exists
	m._exists = true
	return nil
}

// Delete deletes the Mapping from the database.
func (m *Mapping) Delete(ctx context.Context, db DB) error {
	switch {
	case !m._exists: // doesn't exist
		return nil
	case m._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM mapping ` +
		`WHERE tvdb_id = $1`
	// run
	logf(sqlstr, m.TvdbID)
	if _, err := db.ExecContext(ctx, sqlstr, m.TvdbID); err != nil {
		return logerror(err)
	}
	// set deleted
	m._deleted = true
	return nil
}

// MappingByTvdbID retrieves a row from 'mapping' as a Mapping.
//
// Generated from index 'sqlite_autoindex_mapping_1'.
func MappingByTvdbID(ctx context.Context, db DB, tvdbID string) (*Mapping, error) {
	// query
	const sqlstr = `SELECT ` +
		`tvdb_id, anilist_id ` +
		`FROM mapping ` +
		`WHERE tvdb_id = $1`
	// run
	logf(sqlstr, tvdbID)
	m := Mapping{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, tvdbID).Scan(&m.TvdbID, &m.AnilistID); err != nil {
		return nil, logerror(err)
	}
	return &m, nil
}

// MappingByAnilistID retrieves a row from 'mapping' as a Mapping.
//
// Generated from index 'sqlite_autoindex_mapping_2'.
func MappingByAnilistID(ctx context.Context, db DB, anilistID sql.NullString) (*Mapping, error) {
	// query
	const sqlstr = `SELECT ` +
		`tvdb_id, anilist_id ` +
		`FROM mapping ` +
		`WHERE anilist_id = $1`
	// run
	logf(sqlstr, anilistID)
	m := Mapping{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, anilistID).Scan(&m.TvdbID, &m.AnilistID); err != nil {
		return nil, logerror(err)
	}
	return &m, nil
}
