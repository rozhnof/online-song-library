package pgrepo

import (
	"cmp"
	"context"
	"log/slog"
	repo "song-service/internal/application/repository"
	"song-service/internal/domain/models"
	"song-service/internal/infrastructure/database/postgres"
	"song-service/internal/infrastructure/repository/queries"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

type SongRepository struct {
	txManager postgres.TransactionManager
	logger    *slog.Logger
	tracer    trace.Tracer
}

func NewSongRepository(txManager postgres.TransactionManager, logger *slog.Logger, tracer trace.Tracer) *SongRepository {
	return &SongRepository{
		txManager: txManager,
		logger:    logger,
		tracer:    tracer,
	}
}

func (s *SongRepository) Create(ctx context.Context, song models.Song) (models.Song, error) {
	ctx, span := s.tracer.Start(ctx, "SongRepository.Create")
	defer span.End()

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		db := s.txManager.TxOrDB(ctx)
		querier := queries.New(db)

		groupID, err := querier.CreateGroup(ctx, song.Group)
		if err != nil {
			s.logger.Warn("execute query failed", slog.String("error", err.Error()))

			return err
		}

		songArgs := queries.CreateSongParams{
			Name:        song.Name,
			GroupID:     groupID,
			ReleaseDate: song.ReleaseDate,
			Text:        song.Text,
			Link:        song.Link,
		}

		songID, err := querier.CreateSong(ctx, songArgs)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.Wrapf(repo.ErrDuplicate, "song with name = %s and group = %s already exists", song.Name, song.Group)
			}

			s.logger.Warn("execute query failed", slog.String("error", err.Error()))

			return err
		}

		song.ID = songID

		return nil
	}); err != nil {
		return models.Song{}, err
	}

	return song, nil
}

func (s *SongRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Song, error) {
	ctx, span := s.tracer.Start(ctx, "SongRepository.GetByID")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)
	querier := queries.New(db)

	row, err := querier.GetSongByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Song{}, errors.Wrapf(repo.ErrObjectNotFound, "song with id = %s not found", id.String())
		}

		s.logger.Warn("execute query failed", slog.String("error", err.Error()))

		return models.Song{}, err
	}

	song := models.Song{
		ID:          row.Song.ID,
		Name:        row.Song.Name,
		Group:       row.Group.Name,
		ReleaseDate: row.Song.ReleaseDate,
		Text:        row.Song.Text,
		Link:        row.Song.Link,
	}

	return song, nil
}

func (s *SongRepository) List(ctx context.Context, filter *repo.SongFilter, pagination *repo.Pagination) ([]models.Song, error) {
	ctx, span := s.tracer.Start(ctx, "SongRepository.List")
	defer span.End()

	db := s.txManager.TxOrDB(ctx)
	querier := queries.New(db)

	var args queries.ListSongParams

	if filter != nil {
		args.Name = filter.Name
		args.Group = filter.Group
		args.Text = filter.Text
		args.Link = filter.Link
		args.ReleaseDateFrom = filter.ReleaseDateFrom
		args.ReleaseDateTo = filter.ReleaseDateTo
	}

	if pagination != nil {
		if pagination.Limit > 0 {
			args.Limit = &pagination.Limit
		}

		args.Offset = pagination.Offset
	}

	rows, err := querier.ListSong(ctx, args)
	if err != nil {
		s.logger.Warn("execute query failed", slog.String("error", err.Error()))

		return nil, err
	}

	songList := make([]models.Song, 0, len(rows))
	for _, row := range rows {
		song := models.Song{
			ID:          row.Song.ID,
			Name:        row.Song.Name,
			Group:       row.Group.Name,
			ReleaseDate: row.Song.ReleaseDate,
			Text:        row.Song.Text,
			Link:        row.Song.Link,
		}

		songList = append(songList, song)
	}

	return songList, nil
}

func (s *SongRepository) Update(ctx context.Context, song models.Song) (models.Song, error) {
	ctx, span := s.tracer.Start(ctx, "SongRepository.Update")
	defer span.End()

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		db := s.txManager.TxOrDB(ctx)
		querier := queries.New(db)

		row, err := querier.GetSongByID(ctx, song.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.Wrapf(repo.ErrObjectNotFound, "song with id = %s not found", song.ID)
			}

			s.logger.Warn("execute query failed", slog.String("error", err.Error()))

			return err
		}

		song.Name = cmp.Or(song.Name, row.Song.Name)
		song.Group = cmp.Or(song.Group, row.Group.Name)
		song.ReleaseDate = cmp.Or(song.ReleaseDate, row.Song.ReleaseDate)
		song.Text = cmp.Or(song.Text, row.Song.Text)
		song.Link = cmp.Or(song.Link, row.Song.Link)

		songArgs := queries.UpdateSongParams{
			ID:          song.ID,
			Name:        song.Name,
			ReleaseDate: song.ReleaseDate,
			Text:        song.Text,
			Link:        song.Link,
		}

		groupID, err := querier.CreateGroup(ctx, song.Group)
		if err != nil {
			s.logger.Warn("execute query failed", slog.String("error", err.Error()))

			return err
		}

		songArgs.GroupID = groupID

		if err := querier.UpdateSong(ctx, songArgs); err != nil {
			s.logger.Warn("execute query failed", slog.String("error", err.Error()))

			var pgErr *pgconn.PgError

			if errors.As(err, &pgErr) {
				if pgErr.Code == pgerrcode.UniqueViolation {
					return errors.Wrapf(repo.ErrDuplicate, "song with name = %s and group = %s already exists", song.Name, song.Group)
				}
			}

			return err
		}

		return nil
	}); err != nil {
		return models.Song{}, err
	}

	return song, nil
}

func (s *SongRepository) Delete(ctx context.Context, id uuid.UUID) (*time.Time, error) {
	ctx, span := s.tracer.Start(ctx, "SongRepository.Delete")
	defer span.End()

	var deletedTime *time.Time

	if err := s.txManager.WithTransaction(ctx, func(ctx context.Context) error {
		db := s.txManager.TxOrDB(ctx)
		querier := queries.New(db)

		deletedAt, err := querier.DeleteSong(ctx, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return errors.Wrapf(repo.ErrObjectNotFound, "song with id = %s not found", id.String())
			}

			s.logger.Warn("execute query failed", slog.String("error", err.Error()))

			return err
		}

		deletedTime = deletedAt

		return nil
	}); err != nil {
		return nil, err
	}

	return deletedTime, nil
}
