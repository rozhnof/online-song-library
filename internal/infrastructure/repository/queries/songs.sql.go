// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: songs.sql

package queries

import (
	"context"
	"time"

	"github.com/google/uuid"
	date "github.com/hardfinhq/go-date"
)

const createSong = `-- name: CreateSong :one

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
RETURNING id
`

type CreateSongParams struct {
	Name        string
	GroupID     uuid.UUID
	ReleaseDate date.Date
	Text        string
	Link        string
}

// songs.sql
func (q *Queries) CreateSong(ctx context.Context, arg CreateSongParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createSong,
		arg.Name,
		arg.GroupID,
		arg.ReleaseDate,
		arg.Text,
		arg.Link,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteSong = `-- name: DeleteSong :one
UPDATE
    songs
SET
    deleted_at = COALESCE(deleted_at, NOW())
WHERE
    id = $1
RETURNING
    deleted_at
`

func (q *Queries) DeleteSong(ctx context.Context, id uuid.UUID) (*time.Time, error) {
	row := q.db.QueryRow(ctx, deleteSong, id)
	var deleted_at *time.Time
	err := row.Scan(&deleted_at)
	return deleted_at, err
}

const getSongByID = `-- name: GetSongByID :one
SELECT
    s.id, s.name, s.group_id, s.release_date, s.text, s.link, s.deleted_at,
    g.id, g.name, g.deleted_at
FROM 
    songs s
JOIN
    groups g ON s.group_id = g.id
WHERE
    s.id = $1
    AND s.deleted_at IS NULL
    AND g.deleted_at IS NULL
`

type GetSongByIDRow struct {
	Song  Song
	Group Group
}

func (q *Queries) GetSongByID(ctx context.Context, id uuid.UUID) (GetSongByIDRow, error) {
	row := q.db.QueryRow(ctx, getSongByID, id)
	var i GetSongByIDRow
	err := row.Scan(
		&i.Song.ID,
		&i.Song.Name,
		&i.Song.GroupID,
		&i.Song.ReleaseDate,
		&i.Song.Text,
		&i.Song.Link,
		&i.Song.DeletedAt,
		&i.Group.ID,
		&i.Group.Name,
		&i.Group.DeletedAt,
	)
	return i, err
}

const listSong = `-- name: ListSong :many
SELECT
    s.id, s.name, s.group_id, s.release_date, s.text, s.link, s.deleted_at,
    g.id, g.name, g.deleted_at
FROM 
    songs s
JOIN
    groups g ON s.group_id = g.id
WHERE
    s.deleted_at IS NULL
    AND g.deleted_at IS NULL
    AND ($1::VARCHAR(255)[] IS NULL OR s.name = ANY($1::VARCHAR(255)[]))
    AND ($2::VARCHAR(255)[] IS NULL OR g.name = ANY($2::VARCHAR(255)[]))
    AND ($3::DATE IS NULL OR s.release_date >= $3::DATE)
    AND ($4::DATE IS NULL OR s.release_date <= $4::DATE)
    AND ($5::TEXT[] IS NULL OR s.text = ANY($5::TEXT[]))
    AND ($6::TEXT[] IS NULL OR s.link = ANY($6::TEXT[]))
LIMIT 
    $8
OFFSET 
    $7
`

type ListSongParams struct {
	Name            []string
	Group           []string
	ReleaseDateFrom *date.Date
	ReleaseDateTo   *date.Date
	Text            []string
	Link            []string
	Offset          int32
	Limit           *int32
}

type ListSongRow struct {
	Song  Song
	Group Group
}

func (q *Queries) ListSong(ctx context.Context, arg ListSongParams) ([]ListSongRow, error) {
	rows, err := q.db.Query(ctx, listSong,
		arg.Name,
		arg.Group,
		arg.ReleaseDateFrom,
		arg.ReleaseDateTo,
		arg.Text,
		arg.Link,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListSongRow{}
	for rows.Next() {
		var i ListSongRow
		if err := rows.Scan(
			&i.Song.ID,
			&i.Song.Name,
			&i.Song.GroupID,
			&i.Song.ReleaseDate,
			&i.Song.Text,
			&i.Song.Link,
			&i.Song.DeletedAt,
			&i.Group.ID,
			&i.Group.Name,
			&i.Group.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSong = `-- name: UpdateSong :exec
UPDATE 
    songs 
SET 
  name = $2,
  group_id = $3,
  release_date = $4,
  text = $5,
  link = $6
WHERE
    id = $1
`

type UpdateSongParams struct {
	ID          uuid.UUID
	Name        string
	GroupID     uuid.UUID
	ReleaseDate date.Date
	Text        string
	Link        string
}

func (q *Queries) UpdateSong(ctx context.Context, arg UpdateSongParams) error {
	_, err := q.db.Exec(ctx, updateSong,
		arg.ID,
		arg.Name,
		arg.GroupID,
		arg.ReleaseDate,
		arg.Text,
		arg.Link,
	)
	return err
}
