package tests

import (
	"context"
	"log"
	"net/http"
	"song-service/internal/infrastructure/database/postgres"
	"song-service/internal/pkg/config"

	"github.com/google/uuid"
	"github.com/hardfinhq/go-date"
)

const (
	configPath = "../config/test-config.yaml"
	baseURL    = "http://localhost:9090"
)

var (
	songServiceDB     SongServiceDatabase
	musicService      *MockMusicService
	songServiceClient *SongServiceClient
)

type Config struct {
	Postgres config.Postgres
}

func init() {
	cfg, err := config.NewConfig[Config](configPath)
	if err != nil {
		log.Fatal(err)
	}

	postgresCfg := postgres.DatabaseConfig{
		Address:  cfg.Postgres.Address,
		Port:     cfg.Postgres.Port,
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		DB:       cfg.Postgres.DB,
		SSL:      cfg.Postgres.SSL,
	}

	database, err := postgres.NewDatabase(context.Background(), postgresCfg)
	if err != nil {
		log.Fatal(err)
	}

	songServiceDB = *NewSongServiceDatabase(database)
	songServiceClient = NewSongClient(http.DefaultClient, baseURL)
	musicService = NewMockMusicService()

	go musicService.Run()
}

func SetUpEmpty() error {
	if err := Erase(); err != nil {
		return err
	}

	return nil
}

func SetUpCreateTest() error {
	if err := Erase(); err != nil {
		return err
	}

	musicService.AddSong(defaultSong)

	return nil
}

func SetUpDefault() error {
	if err := Erase(); err != nil {
		return err
	}

	musicService.AddSong(defaultSong)

	if _, err := songServiceDB.CreateSong(defaultSong); err != nil {
		return err
	}

	return nil
}

func SetUp(musicServiceSongs []Song, songServiceSongs []Song) error {
	if err := Erase(); err != nil {
		return err
	}

	for _, song := range musicServiceSongs {
		musicService.AddSong(song)
	}

	for _, song := range songServiceSongs {
		if _, err := songServiceDB.CreateSong(song); err != nil {
			return err
		}
	}

	return nil
}

func Erase() error {
	musicService.ClearStorage()

	if err := songServiceDB.Truncate(context.Background()); err != nil {
		return err
	}

	return nil
}

var defaultSong = Song{
	ID:          uuid.New(),
	Group:       "existent-song-group",
	Name:        "existent-song-name",
	ReleaseDate: date.NewDate(2025, 1, 1),
	Text:        "existent-song-text",
	Link:        "existent-song-link",
}

var nonExistentSong = Song{
	ID:          uuid.New(),
	Group:       "non-existent-song-group",
	Name:        "non-existent-song-name",
	ReleaseDate: date.NewDate(2025, 1, 1),
	Text:        "non-existent-song-text",
	Link:        "non-existent-song-link",
}
