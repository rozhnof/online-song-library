package services

import (
	"context"
	repo "song-service/internal/application/repository"
	"song-service/internal/domain/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type SongService struct {
	repository repo.SongRepository
	tracer     trace.Tracer
}

func NewSongService(repository repo.SongRepository, tracer trace.Tracer) *SongService {
	return &SongService{
		repository: repository,
		tracer:     tracer,
	}
}

func (s *SongService) CreateSong(ctx context.Context, song models.Song) (models.Song, error) {
	ctx, span := s.tracer.Start(ctx, "SongService.CreateSong")
	defer span.End()

	createdSong, err := s.repository.Create(ctx, song)
	if err != nil {
		return models.Song{}, err
	}

	return createdSong, nil
}

func (s *SongService) UpdateSong(ctx context.Context, song models.Song) (models.Song, error) {
	ctx, span := s.tracer.Start(ctx, "SongService.UpdateSong")
	defer span.End()

	updatedSong, err := s.repository.Update(ctx, song)
	if err != nil {
		return models.Song{}, err
	}

	return updatedSong, nil
}

func (s *SongService) Song(ctx context.Context, id uuid.UUID, pagination *repo.Pagination) (models.Song, error) {
	ctx, span := s.tracer.Start(ctx, "SongService.Song")
	defer span.End()

	song, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return models.Song{}, err
	}

	if pagination != nil {
		song = addTextPagination(song, *pagination)
	}

	return song, nil
}

func (s *SongService) SongList(ctx context.Context, filter *repo.SongFilter, pagination *repo.Pagination) ([]models.Song, error) {
	ctx, span := s.tracer.Start(ctx, "SongService.SongList")
	defer span.End()

	songList, err := s.repository.List(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}

	return songList, nil
}

func (s *SongService) DeleteSong(ctx context.Context, id uuid.UUID) (*time.Time, error) {
	ctx, span := s.tracer.Start(ctx, "SongService.DeleteSong")
	defer span.End()

	deletedTime, err := s.repository.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return deletedTime, nil
}

func addTextPagination(song models.Song, pagination repo.Pagination) models.Song {
	const delim = "\n\n"

	verseList := strings.Split(song.Text, delim)

	if pagination.Offset >= int32(len(verseList)) {
		song.Text = ""
		return song
	}

	if pagination.Offset > 0 {
		verseList = verseList[pagination.Offset:]
	}

	if pagination.Limit <= 0 || pagination.Limit >= int32(len(verseList)) {
		song.Text = strings.Join(verseList, delim)
		return song
	}

	song.Text = strings.Join(verseList[:pagination.Limit], delim)
	return song
}
