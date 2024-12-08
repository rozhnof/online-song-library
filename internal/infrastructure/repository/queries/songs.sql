-- songs.sql

-- name: CreateSong :one
INSERT INTO songs (
    name,
    group_id,
    release_date,
    text,
    link
)
VALUES (
    $1, $2, $3, $4, $5
)
ON CONFLICT (name, group_id) 
DO UPDATE 
SET
    release_date = EXCLUDED.release_date,
    text = EXCLUDED.text,
    link = EXCLUDED.link
WHERE 
    songs.release_date = EXCLUDED.release_date
    AND songs.text = EXCLUDED.text
    AND songs.link = EXCLUDED.link
RETURNING id;


-- name: UpdateSong :exec
UPDATE 
    songs 
SET 
  name = $2,
  group_id = $3,
  release_date = $4,
  text = $5,
  link = $6
WHERE
    id = $1;


-- name: DeleteSong :one
UPDATE
    songs
SET
    deleted_at = COALESCE(deleted_at, NOW())
WHERE
    id = $1
RETURNING
    deleted_at;


-- name: GetSongByID :one
SELECT
    sqlc.embed(s),
    sqlc.embed(g)
FROM 
    songs s
JOIN
    groups g ON s.group_id = g.id
WHERE
    s.id = $1
    AND s.deleted_at IS NULL
    AND g.deleted_at IS NULL;


-- name: ListSong :many
SELECT
    sqlc.embed(s),
    sqlc.embed(g)
FROM 
    songs s
JOIN
    groups g ON s.group_id = g.id
WHERE
    s.deleted_at IS NULL
    AND g.deleted_at IS NULL
    AND (sqlc.narg('name')::VARCHAR(255)[] IS NULL OR s.name = ANY(sqlc.narg('name')::VARCHAR(255)[]))
    AND (sqlc.narg('group')::VARCHAR(255)[] IS NULL OR g.name = ANY(sqlc.narg('group')::VARCHAR(255)[]))
    AND (sqlc.narg('release_date_from')::DATE IS NULL OR s.release_date >= sqlc.narg('release_date_from')::DATE)
    AND (sqlc.narg('release_date_to')::DATE IS NULL OR s.release_date <= sqlc.narg('release_date_to')::DATE)
    AND (sqlc.narg('text')::TEXT[] IS NULL OR s.text = ANY(sqlc.narg('text')::TEXT[]))
    AND (sqlc.narg('link')::TEXT[] IS NULL OR s.link = ANY(sqlc.narg('link')::TEXT[]))
LIMIT 
    sqlc.narg('limit')
OFFSET 
    sqlc.arg('offset');
