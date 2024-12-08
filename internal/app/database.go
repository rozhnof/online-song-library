package app

import (
	"context"

	"song-service/internal/infrastructure/database/postgres"
	"song-service/internal/pkg/config"
)

func NewPostgresDatabase(ctx context.Context, cfg config.Postgres) (postgres.Database, error) {
	postgresConfig := postgres.DatabaseConfig{
		Address:  cfg.Address,
		Port:     cfg.Port,
		User:     cfg.User,
		Password: cfg.Password,
		DB:       cfg.DB,
		SSL:      cfg.SSL,
	}

	return postgres.NewDatabase(ctx, postgresConfig)
}
