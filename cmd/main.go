package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	_ "song-service/docs"
	"song-service/internal/app"
	"song-service/internal/pkg/config"

	"syscall"
)

const (
	EnvConfigPath = "CONFIG_PATH"
)

// @title     Song Service API
// @version   1.0
// @host      localhost:8080
// @BasePath  /
func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.NewConfig[app.Config](os.Getenv(EnvConfigPath))
	if err != nil {
		log.Fatal(err)
	}

	logger, err := app.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("init logger success")

	postgresDatabase, err := app.NewPostgresDatabase(ctx, cfg.Postgres)
	if err != nil {
		logger.Error("init postgres failed", slog.String("error", err.Error()))
		return
	}
	defer postgresDatabase.Close()
	logger.Info("init postgres success")

	tracer, shutdown, err := app.NewTracer(ctx, cfg.Tracing, app.ServiceName)
	if err != nil {
		logger.Error("init tracer failed", slog.String("error", err.Error()))
		return
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			logger.Error("close tracer failed", slog.String("error", err.Error()))
		}
	}()
	logger.Info("init tracer success")

	authApp := app.NewSongApp(ctx, cfg, logger, postgresDatabase, tracer)
	logger.Info("init app success")

	logger.Info("run app")
	if err := authApp.Run(ctx); err != nil {
		logger.Error("run app error", slog.String("error", err.Error()))
		return
	}

	logger.Error("shutdown app")
}
