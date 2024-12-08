package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	repo "song-service/internal/application/repository"
	"song-service/internal/domain/models"

	"github.com/gin-gonic/gin"
)

type CreateSongRequest struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song"  binding:"required"`
}

type CreateSongResponse struct {
	Song models.Song `json:"song"`
}

// CreateSong godoc
// @Summary      Add song from Music Service
// @Description  Добавление песни в библиотеку из Music Service
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        request body     CreateSongRequest  true  "Song details"
// @Success      200    {object}  CreateSongResponse
// @Failure      400    {string}  string             "Invalid input data"
// @Failure      409    {string}  string             "Song already exists"
// @Failure      500    {string}  string             "Internal Server Error"
// @Router       /songs [post]
func (h *SongHandler) CreateSong(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "SongHandler.CreateSong")
	defer span.End()

	var request CreateSongRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Debug("failed to parse request body", slog.String("error", err.Error()))

		c.String(http.StatusBadRequest, err.Error())
		return
	}

	songInfo, code, err := h.client.Info(ctx, request.Group, request.Song)
	if err != nil {
		h.logger.Warn("failed to get song info from music service", slog.String("error", err.Error()))

		if code != 0 {
			c.String(code, err.Error())
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	song := models.Song{
		Name:        request.Song,
		Group:       request.Group,
		ReleaseDate: songInfo.ReleaseDate,
		Text:        songInfo.Text,
		Link:        songInfo.Link,
	}

	createdSong, err := h.songService.CreateSong(ctx, song)
	if err != nil {
		if errors.Is(err, repo.ErrDuplicate) {
			c.String(http.StatusConflict, err.Error())
			return
		}

		h.logger.Warn("failed to create song", slog.String("error", err.Error()))

		c.Status(http.StatusInternalServerError)
		return
	}

	response := CreateSongResponse{
		Song: createdSong,
	}

	c.JSON(http.StatusOK, response)
}
