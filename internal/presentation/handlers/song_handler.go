package handlers

import (
	"log/slog"
	"song-service/internal/application/services"
	"song-service/internal/presentation/client"

	"go.opentelemetry.io/otel/trace"
)

const (
	pathParamID = "id"
)

type SongHandler struct {
	songService *services.SongService
	logger      *slog.Logger
	tracer      trace.Tracer
	client      *client.MusicServiceClient
}

func NewSongHandler(songService *services.SongService, logger *slog.Logger, tracer trace.Tracer, client *client.MusicServiceClient) *SongHandler {
	return &SongHandler{
		songService: songService,
		logger:      logger,
		tracer:      tracer,
		client:      client,
	}
}
