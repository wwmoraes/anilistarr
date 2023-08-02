CREATE TABLE mapping (
    tvdb_id    TEXT UNIQUE
                    PRIMARY KEY
                    NOT NULL,
    anilist_id TEXT UNIQUE
)
WITHOUT ROWID,
STRICT;
