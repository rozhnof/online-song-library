package tests

import (
	"context"
	"song-service/internal/infrastructure/database/postgres"

	"github.com/google/uuid"
)

type SongServiceDatabase struct {
	db postgres.Database
}

func NewSongServiceDatabase(db postgres.Database) *SongServiceDatabase {
	return &SongServiceDatabase{
		db: db,
	}
}

func (d *SongServiceDatabase) CreateGroup(name string) (uuid.UUID, error) {
	const query = `
		INSERT INTO groups (name)
		VALUES ($1)
		ON CONFLICT (name) 
		DO UPDATE 
		SET name = EXCLUDED.name
		RETURNING id;
	`

	var groupID uuid.UUID

	if err := d.db.QueryRow(context.Background(), query, name).Scan(&groupID); err != nil {
		return uuid.UUID{}, err
	}

	return groupID, nil
}

func (d *SongServiceDatabase) CreateSong(song Song) (Song, error) {
	groupID, err := d.CreateGroup(song.Group)
	if err != nil {
		return Song{}, err
	}

	const query = `
		INSERT INTO songs (id, name, group_id, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	if err := d.db.QueryRow(context.Background(), query, song.ID, song.Name, groupID, song.ReleaseDate, song.Text, song.Link).Scan(&song.ID); err != nil {
		return Song{}, err
	}

	return song, nil
}

func (d *SongServiceDatabase) GetSongByID(songID uuid.UUID) (Song, error) {
	const query = `
		SELECT s.id, s.name, g.name, s.release_date, s.text, s.link
		FROM songs s
		JOIN groups g ON s.group_id = g.id
		WHERE s.id = $1 AND s.deleted_at IS NULL AND g.deleted_at IS NULL;
	`

	row := d.db.QueryRow(context.Background(), query, songID)

	var song Song

	if err := row.Scan(
		&song.ID,
		&song.Name,
		&song.Group,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
	); err != nil {
		return Song{}, err
	}

	return song, nil
}

func (d *SongServiceDatabase) Truncate(ctx context.Context) error {
	query := `
		DO $$ DECLARE
			table_name TEXT;
		BEGIN
			FOR table_name IN 
				SELECT tablename 
				FROM pg_tables 
				WHERE schemaname = 'public'
			LOOP
				EXECUTE format('TRUNCATE TABLE %I CASCADE', table_name);
			END LOOP;
		END $$;
	`
	if _, err := d.db.Exec(ctx, query); err != nil {
		return err
	}

	return nil
}
