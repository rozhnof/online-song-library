package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	repo "song-service/internal/application/repository"
	"song-service/internal/domain/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hardfinhq/go-date"
)

type UpdateSongRequest struct {
	Name        string    `json:"song"         binding:"required"`
	Group       string    `json:"group"        binding:"required"`
	ReleaseDate date.Date `json:"release_date" binding:"required" swaggertype:"primitive,string"`
	Text        string    `json:"text"         binding:"required"`
	Link        string    `json:"link"         binding:"required"`
}

type UpdateSongResponse struct {
	Song models.Song `json:"song"`
}

// UpdateSong godoc
// @Summary      Update song by ID
// @Description  Полное обновление информации о песне в библиотеке
// @Tags         songs
// @Accept       json
// @Produce      json
// @Param        id       path     string             true   "Song ID"
// @Param        request  body     UpdateSongRequest  true   "Song details to update"
// @Success      200      {object} UpdateSongResponse
// @Failure      400      {string} string             "Invalid input data"
// @Failure      404      {string} string             "Song not found"
// @Failure      409      {string} string             "Name conflict"
// @Failure      500      {string} string             "Internal Server Error"
// @Router       /songs/{id} [put]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "SongHandler.UpdateSong")
	defer span.End()

	id, err := uuid.Parse(c.Param(pathParamID))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ID format")
		return
	}

	var request UpdateSongRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Debug("failed to parse request body", slog.String("error", err.Error()))

		c.String(http.StatusBadRequest, err.Error())
		return
	}

	song := models.Song{
		ID:          id,
		Name:        request.Name,
		Group:       request.Group,
		ReleaseDate: request.ReleaseDate,
		Text:        request.Text,
		Link:        request.Link,
	}

	updatedSong, err := h.songService.UpdateSong(ctx, song)
	if err != nil {
		if errors.Is(err, repo.ErrObjectNotFound) {
			c.String(http.StatusNotFound, err.Error())
			return
		}

		if errors.Is(err, repo.ErrDuplicate) {
			c.String(http.StatusConflict, err.Error())
			return
		}

		h.logger.Warn("failed to create song", slog.String("error", err.Error()))

		c.Status(http.StatusInternalServerError)
		return
	}

	response := SongResponse{
		Song: updatedSong,
	}

	c.JSON(http.StatusOK, response)
}
