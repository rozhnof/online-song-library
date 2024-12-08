package repo

import (
	"context"
	"song-service/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type SongRepository interface {
	Create(ctx context.Context, song models.Song) (models.Song, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Song, error)
	List(ctx context.Context, filter *SongFilter, pagination *Pagination) ([]models.Song, error)
	Update(ctx context.Context, song models.Song) (models.Song, error)
	Delete(ctx context.Context, id uuid.UUID) (*time.Time, error)
}
