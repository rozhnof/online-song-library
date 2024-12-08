package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"song-service/internal/application/services"
	"song-service/internal/infrastructure/database/postgres"
	pgrepo "song-service/internal/infrastructure/repository"
	"song-service/internal/pkg/config"
	"song-service/internal/pkg/server"
	"song-service/internal/presentation/client"
	handlers "song-service/internal/presentation/handlers"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
)

const (
	ServiceName = "Song Service"
)

type Config struct {
	Mode         string              `yaml:"mode"          env-required:"true"`
	Server       config.Server       `yaml:"server"        env-required:"true"`
	Logger       config.Logger       `yaml:"logging"       env-required:"true"`
	Tracing      config.Tracing      `yaml:"tracing"       env-required:"true"`
	MusicService config.MusicService `yaml:"music_service" env-required:"true"`
	Postgres     config.Postgres
}

type SongApp struct {
	logger     *slog.Logger
	httpServer *server.HTTPServer
}

func NewSongApp(ctx context.Context, cfg *Config, logger *slog.Logger, postgresDatabase postgres.Database, tracer trace.Tracer) *SongApp {
	var (
		txManager = postgres.NewTransactionManager(postgresDatabase.Pool)
	)

	musicServiceClient := client.NewMusicServiceClient(&http.Client{
		Timeout: cfg.MusicService.Timeout,
	}, cfg.MusicService.Address)

	var (
		songRepository = pgrepo.NewSongRepository(txManager, logger, tracer)
		songService    = services.NewSongService(songRepository, tracer)
		songHandler    = handlers.NewSongHandler(songService, logger, tracer, musicServiceClient)
	)

	gin.SetMode(cfg.Mode)

	var (
		router = gin.New()
	)

	router.Use(
		gin.Recovery(),
		otelgin.Middleware(ServiceName),
		LogMiddleware(logger),
	)

	InitRoutes(router, songHandler)

	var (
		httpServer = server.NewHTTPServer(ctx, cfg.Server.Address, router)
	)

	return &SongApp{
		logger:     logger,
		httpServer: httpServer,
	}
}

func (a *SongApp) Run(ctx context.Context) error {
	errChan := make(chan error)
	defer close(errChan)

	go func() {
		if err := a.httpServer.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	go func() {
		<-ctx.Done()

		errChan <- a.httpServer.Shutdown()
	}()

	return <-errChan
}
