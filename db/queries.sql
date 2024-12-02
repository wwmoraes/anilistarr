-- name: GetCacheString :one
SELECT value
FROM cache
WHERE key = @key
LIMIT 1;

-- name: PutCacheString :exec
REPLACE INTO cache (key, value)
VALUES (@key, @value);

-- name: GetMedia :one
SELECT * FROM medias
WHERE source_id = @id
LIMIT 1;

-- name: GetMediaBulk :many
SELECT * FROM medias
WHERE source_id
IN (sqlc.slice(source_ids));

-- name: PutMedia :exec
REPLACE INTO medias (source_id, target_id)
VALUES (@source_id, @target_id);
