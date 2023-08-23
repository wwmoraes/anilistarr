CREATE TABLE mapping (
  target_id TEXT UNIQUE
            PRIMARY KEY
            NOT NULL,
  source_id TEXT UNIQUE
            NOT NULL
)
WITHOUT ROWID,
STRICT;
