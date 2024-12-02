PRAGMA encoding = 'UTF-8';

CREATE TABLE IF NOT EXISTS medias (
	source_id TEXT NOT NULL, -- VARCHAR(64)
	target_id TEXT NOT NULL, -- VARCHAR(64)
	CHECK(source_id <> ''),
	CHECK(target_id <> ''),
	PRIMARY KEY(source_id, target_id)
) WITHOUT ROWID, STRICT;
CREATE UNIQUE INDEX IF NOT EXISTS
	medias_unique ON medias (source_id, target_id);

CREATE INDEX IF NOT EXISTS
	medias_source_id ON medias (source_id);

CREATE TABLE IF NOT EXISTS users (
	id   TEXT NOT NULL, -- VARCHAR(64)
	name TEXT NOT NULL, -- VARCHAR(64)
	CHECK(id <> ''),
	CHECK(name <> ''),
	PRIMARY KEY (id, name)
) WITHOUT ROWID, STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS
	users_unique ON users (name, id);

CREATE INDEX IF NOT EXISTS
	users_id ON users (id);

CREATE TABLE IF NOT EXISTS cache (
	key   TEXT NOT NULL, -- VARCHAR(64)
	value TEXT NOT NULL,
	CHECK(key <> ''),
	CHECK(value <> ''),
	PRIMARY KEY(key)
) WITHOUT ROWID, STRICT;

CREATE UNIQUE INDEX IF NOT EXISTS
	cache_unique ON cache (key);

CREATE INDEX IF NOT EXISTS
	cache_key ON cache (key);
